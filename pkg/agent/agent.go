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

package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"go.uber.org/zap"
	"knative.dev/pkg/logging"
)

type Decider struct {
	cache  *cachedQuery
	logger *zap.SugaredLogger
}

type cachedQuery struct {
	policyPath string
	query      rego.PreparedEvalQuery
}

func (c *cachedQuery) start(ctx context.Context) error {
	if err := c.load(ctx); err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-time.After(5 * time.Second):
				if err := c.load(ctx); err != nil {
					logging.FromContext(ctx).Errorf("%v", err)
				} else {
					logging.FromContext(ctx).Debugf("file %q refreshed", c.policyPath)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (c *cachedQuery) load(ctx context.Context) error {
	b, err := ioutil.ReadFile(c.policyPath)
	if err != nil {
		return fmt.Errorf("failed to refresh policy file %q: %w", c.policyPath, err)
	}

	module := string(b)
	compiler, err := ast.CompileModules(map[string]string{
		"policy": module,
	})
	if err != nil {
		return fmt.Errorf("failed to compile rego module: %w", err)
	}

	r := rego.New(
		rego.Query("data.security.knative.dev.allow"),
		rego.Compiler(compiler),
	)

	pq, err := r.PrepareForEval(ctx)
	if err != nil {
		return fmt.Errorf("failed to prepare for eval: %w", err)
	}

	c.query = pq
	return nil
}

func NewDecider(ctx context.Context, policyPath string) (*Decider, error) {
	c := &cachedQuery{policyPath: policyPath}
	if err := c.start(ctx); err != nil {
		return nil, err
	}
	return &Decider{
		cache:  c,
		logger: logging.FromContext(ctx),
	}, nil
}

func (d *Decider) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	d.logger.Debug("Received decision request")

	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	var dreq DecisionRequest
	if err := json.Unmarshal(b, &dreq); err != nil {
		d.logger.Errorf("Failed to unmarshal decision request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dresp := &DecisionResponse{
		Allow: d.isAllowed(req.Context(), dreq),
	}

	respBytes, err := json.Marshal(dresp)
	if err != nil {
		d.logger.Errorf("Failed to marshal decision response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respBytes)
}

func (d *Decider) isAllowed(ctx context.Context, dr DecisionRequest) bool {
	d.logger.Debugf("Evaluating input: %v", *dr.HTTPRequest)

	rs, err := d.cache.query.Eval(ctx, rego.EvalInput(dr))
	if err != nil {
		d.logger.Warnw("failed to evaluate input and will fail-close", zap.Error(err))
		return false
	}

	if len(rs) == 0 || len(rs[0].Expressions) == 0 {
		return false
	}

	d.logger.Debugf("Eval result: %v", rs[0])

	// Super hacky
	return strings.Contains(fmt.Sprintf("%v", rs[0]), "true")
}
