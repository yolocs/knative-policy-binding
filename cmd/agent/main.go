/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/kelseyhightower/envconfig"
	"github.com/yolocs/knative-policy-binding/pkg/agent"
	"go.uber.org/zap"
	pkglogging "knative.dev/pkg/logging"
	pkgnet "knative.dev/pkg/network"
	"knative.dev/pkg/signals"
)

type config struct {
	AgentPort          int    `envconfig:"AGENT_PORT" required:"true"`
	PolicyPath         string `envconfig:"POLICY_PATH" required:"true"`
	AgentLoggingConfig string `envconfig:"AGENT_LOGGING_CONFIG" required:"true"`
	AgentLoggingLevel  string `envconfig:"AGENT_LOGGING_LEVEL" required:"true"`
}

var logger *zap.SugaredLogger

func main() {
	// Parse the environment.
	var env config
	if err := envconfig.Process("", &env); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	logger, _ = pkglogging.NewLogger(env.AgentLoggingConfig, env.AgentLoggingLevel)
	logger = logger.Named("knative-policy-agent")
	defer flush(logger)

	ctx, cancel := context.WithCancel(pkglogging.WithLogger(context.Background(), logger))
	decider, err := agent.NewDecider(ctx, env.PolicyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create decider: %v", err)
		os.Exit(1)
	}

	servers := map[string]*http.Server{
		"decision-server": pkgnet.NewServer(":"+strconv.Itoa(env.AgentPort), decider),
	}

	errCh := make(chan error, len(servers))
	for name, server := range servers {
		go func(name string, s *http.Server) {
			logger.Infof("Starting server %q...", name)
			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errCh <- fmt.Errorf("%s server failed: %w", name, err)
			}
		}(name, server)
	}

	// Blocks until we actually receive a TERM signal or one of the servers
	// exit unexpectedly. We fold both signals together because we only want
	// to act on the first of those to reach here.
	select {
	case err := <-errCh:
		logger.Errorw("Failed to bring up kn-proxy, shutting down.", zap.Error(err))
		flush(logger)
		os.Exit(1)
	case <-signals.SetupSignalHandler():
		logger.Info("Received TERM signal, attempting to gracefully shutdown servers.")
		flush(logger)
		cancel()
		for serverName, srv := range servers {
			if err := srv.Shutdown(context.Background()); err != nil {
				logger.Errorw("Failed to shutdown server", zap.String("server", serverName), zap.Error(err))
			}
		}
	}
}

func flush(logger *zap.SugaredLogger) {
	logger.Sync()
	os.Stdout.Sync()
	os.Stderr.Sync()
}
