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

package eventpolicy

import (
	"context"

	security "github.com/yolocs/knative-policy-binding/pkg/apis/security/v1alpha1"
	eventpolicyinformer "github.com/yolocs/knative-policy-binding/pkg/client/injection/informers/security/v1alpha1/eventpolicy"
	openpolicyinformer "github.com/yolocs/knative-policy-binding/pkg/client/injection/informers/security/v1alpha1/openpolicy"
	"github.com/yolocs/knative-policy-binding/pkg/reconciler"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
)

const (
	reconcilerName      = "EventPolicies"
	controllerAgentName = "eventpolicy-controller"
)

// NewController initializes the controller and is called by the generated code
// Registers event handlers to enqueue events
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {

	eventpolicyInformer := eventpolicyinformer.Get(ctx)
	openpolicyInformer := openpolicyinformer.Get(ctx)

	r := &Reconciler{
		Base:              reconciler.NewBase(ctx, controllerAgentName, cmw),
		eventpolicyLister: eventpolicyInformer.Lister(),
		openpolicyLister:  openpolicyInformer.Lister(),
	}
	impl := controller.NewImpl(r, r.Logger, reconcilerName)

	r.Logger.Info("Setting up event handlers")

	eventpolicyInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	openpolicyInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.Filter(security.SchemeGroupVersion.WithKind("EventPolicy")),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	return impl
}
