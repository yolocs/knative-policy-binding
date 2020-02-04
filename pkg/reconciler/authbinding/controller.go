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

package authbinding

import (
	"context"

	security "github.com/yolocs/knative-policy-binding/pkg/apis/security/v1alpha1"
	authbindinginformer "github.com/yolocs/knative-policy-binding/pkg/client/injection/informers/security/v1alpha1/authorizablebinding"
	policybindinginformer "github.com/yolocs/knative-policy-binding/pkg/client/injection/informers/security/v1alpha1/policybinding"
	"github.com/yolocs/knative-policy-binding/pkg/reconciler"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
)

const (
	reconcilerName      = "AuthorizableBinding"
	controllerAgentName = "authbinding-controller"
)

// NewController initializes the controller and is called by the generated code
// Registers event handlers to enqueue events
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {

	authbindingInformer := authbindinginformer.Get(ctx)
	policybindingInformer := policybindinginformer.Get(ctx)

	r := &Reconciler{
		Base:                reconciler.NewBase(ctx, controllerAgentName, cmw),
		authbindingLister:   authbindingInformer.Lister(),
		policybindingLister: policybindingInformer.Lister(),
	}
	impl := controller.NewImpl(r, r.Logger, reconcilerName)
	r.podspecableResolver = NewPodspecableResolver(ctx, impl.EnqueueKey)

	r.Logger.Info("Setting up event handlers")

	authbindingInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	policybindingInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.Filter(security.SchemeGroupVersion.WithKind("AuthorizableBinding")),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	return impl
}
