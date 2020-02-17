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

package opabinding

import (
	"context"
	"log"

	"go.uber.org/zap"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/tracker"

	"github.com/kelseyhightower/envconfig"
	"github.com/yolocs/knative-policy-binding/pkg/apis/security/v1alpha2"
	policyinformer "github.com/yolocs/knative-policy-binding/pkg/client/injection/informers/security/v1alpha2/httppolicy"
	bindinginformer "github.com/yolocs/knative-policy-binding/pkg/client/injection/informers/security/v1alpha2/httppolicybinding"
	policypsbindinginformer "github.com/yolocs/knative-policy-binding/pkg/client/injection/informers/security/v1alpha2/policypodspecablebinding"
	bindingreconciler "github.com/yolocs/knative-policy-binding/pkg/client/injection/reconciler/security/v1alpha2/httppolicybinding"
	"github.com/yolocs/knative-policy-binding/pkg/reconciler"
	"github.com/yolocs/knative-policy-binding/pkg/reconciler/internal/resolver"
	configmapinformer "knative.dev/pkg/client/injection/kube/informers/core/v1/configmap"
)

const (
	// ReconcilerName is the name of the reconciler
	ReconcilerName      = "OPABindings"
	controllerAgentName = "opabinding-controller"
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

	bindingInformer := bindinginformer.Get(ctx)
	policyInformer := policyinformer.Get(ctx)
	configmapInformer := configmapinformer.Get(ctx)
	psbindingInformer := policypsbindinginformer.Get(ctx)

	r := &Reconciler{
		Base:                reconciler.NewBase(ctx, controllerAgentName, cmw),
		policybindingLister: bindingInformer.Lister(),
		policyLister:        policyInformer.Lister(),
		psbindingLister:     psbindingInformer.Lister(),
		configmapLister:     configmapInformer.Lister(),
		agentImage:          env.AgentImage,
	}
	impl := bindingreconciler.NewImpl(ctx, r)

	r.Logger.Info("Setting up event handlers")

	bindingInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	psbindingInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.Filter(v1alpha2.SchemeGroupVersion.WithKind("HTTPPolicyBinding")),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	r.subjectResolver = resolver.NewSubjectResolver(ctx, impl.EnqueueKey)
	r.policyTracker = tracker.New(impl.EnqueueKey, controller.GetTrackerLease(ctx))

	policyInformer.Informer().AddEventHandler(controller.HandleAll(r.policyTracker.OnChanged))

	return impl
}
