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

package openpolicy

import (
	"context"
	"log"

	"go.uber.org/zap"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"

	"github.com/kelseyhightower/envconfig"
	security "github.com/yolocs/knative-policy-binding/pkg/apis/security/v1alpha1"
	openpolicyinformer "github.com/yolocs/knative-policy-binding/pkg/client/injection/informers/security/v1alpha1/openpolicy"
	"github.com/yolocs/knative-policy-binding/pkg/reconciler"
	configmapinformer "knative.dev/pkg/client/injection/kube/informers/core/v1/configmap"
)

const (
	reconcilerName      = "OpenPolicies"
	controllerAgentName = "openpolicy-controller"
)

type config struct {
	AgentImage string `envconfig:"AGENT_IMAGE" required:"true"`
}

// NewController initializes the controller and is called by the generated code
// Registers event handlers to enqueue events
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {

	var env config
	if err := envconfig.Process("", &env); err != nil {
		log.Fatal("Failed to process env var", zap.Error(err))
	}

	openpolicyInformer := openpolicyinformer.Get(ctx)
	configmapInformer := configmapinformer.Get(ctx)

	r := &Reconciler{
		Base:             reconciler.NewBase(ctx, controllerAgentName, cmw),
		configmapLister:  configmapInformer.Lister(),
		openpolicyLister: openpolicyInformer.Lister(),
		agentImage:       env.AgentImage,
	}
	impl := controller.NewImpl(r, r.Logger, reconcilerName)

	r.Logger.Info("Setting up event handlers")

	openpolicyInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	configmapInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.Filter(security.SchemeGroupVersion.WithKind("OpenPolicy")),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	return impl
}
