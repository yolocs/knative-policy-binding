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

package istiobinding

import (
	"context"

	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/tracker"

	"github.com/yolocs/knative-policy-binding/pkg/apis/security/v1alpha2"
	policyinformer "github.com/yolocs/knative-policy-binding/pkg/client/injection/informers/security/v1alpha2/httppolicy"
	bindinginformer "github.com/yolocs/knative-policy-binding/pkg/client/injection/informers/security/v1alpha2/httppolicybinding"
	bindingreconciler "github.com/yolocs/knative-policy-binding/pkg/client/injection/reconciler/security/v1alpha2/httppolicybinding"
	istioclient "github.com/yolocs/knative-policy-binding/pkg/client/istio/injection/client"
	istioauthzinformer "github.com/yolocs/knative-policy-binding/pkg/client/istio/injection/informers/security/v1beta1/authorizationpolicy"
	"github.com/yolocs/knative-policy-binding/pkg/reconciler"
)

const (
	// ReconcilerName is the name of the reconciler
	ReconcilerName      = "IstioBindings"
	controllerAgentName = "istiobinding-controller"
)

// NewController initializes the controller and is called by the generated code
// Registers event handlers to enqueue events
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {

	bindingInformer := bindinginformer.Get(ctx)
	policyInformer := policyinformer.Get(ctx)
	istioauthzInformer := istioauthzinformer.Get(ctx)

	r := &Reconciler{
		Base:                reconciler.NewBase(ctx, controllerAgentName, cmw),
		policybindingLister: bindingInformer.Lister(),
		policyLister:        policyInformer.Lister(),
		istioauthzLister:    istioauthzInformer.Lister(),
		istioClientSet:      istioclient.Get(ctx),
	}
	impl := bindingreconciler.NewImpl(ctx, r)

	r.Logger.Info("Setting up event handlers")

	bindingInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	istioauthzInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.Filter(v1alpha2.SchemeGroupVersion.WithKind("HTTPPolicyBinding")),
		Handler:    controller.HandleAll(impl.Enqueue),
	})

	r.subjectResolver = NewSubjectResolver(ctx, impl.EnqueueKey)
	r.policyTracker = tracker.New(impl.EnqueueKey, controller.GetTrackerLease(ctx))

	return impl
}
