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

package policypsbinding

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"knative.dev/pkg/apis/duck"
	"knative.dev/pkg/client/injection/ducks/duck/v1/podspecable"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/clients/dynamicclient"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/tracker"

	"github.com/yolocs/knative-policy-binding/pkg/apis/security/v1alpha2"
	"github.com/yolocs/knative-policy-binding/pkg/client/clientset/versioned/scheme"
	policybindinginformer "github.com/yolocs/knative-policy-binding/pkg/client/injection/informers/security/v1alpha2/policypodspecablebinding"
	"github.com/yolocs/knative-policy-binding/pkg/webhook/psbinding"
)

const (
	reconcilerName      = "PolicyPodspecableBindings"
	controllerAgentName = "policypsbinding-controller"
)

// NewController initializes the controller and is called by the generated code
// Registers event handlers to enqueue events
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {

	logger := logging.FromContext(ctx)

	policybindinginformer := policybindinginformer.Get(ctx)
	dc := dynamicclient.Get(ctx)
	psInformerFactory := podspecable.Get(ctx)

	r := &psbinding.BaseReconciler{
		GVR: v1alpha2.SchemeGroupVersion.WithResource("policypodspecablebindings"),
		Get: func(namespace string, name string) (psbinding.Bindable, error) {
			return policybindinginformer.Lister().PolicyPodspecableBindings(namespace).Get(name)
		},
		DynamicClient: dc,
		Recorder: record.NewBroadcaster().NewRecorder(
			scheme.Scheme, corev1.EventSource{Component: controllerAgentName}),
	}
	impl := controller.NewImpl(r, logger, reconcilerName)

	logger.Info("Setting up event handlers")

	policybindinginformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	r.WithContext = WithContextFactory(ctx, impl.EnqueueKey)
	r.Tracker = tracker.New(impl.EnqueueKey, controller.GetTrackerLease(ctx))
	r.Factory = &duck.CachedInformerFactory{
		Delegate: &duck.EnqueueInformerFactory{
			Delegate:     psInformerFactory,
			EventHandler: controller.HandleAll(r.Tracker.OnChanged),
		},
	}

	return impl
}

func ListAll(ctx context.Context, handler cache.ResourceEventHandler) psbinding.ListAll {
	policybindinginformer := policybindinginformer.Get(ctx)

	// Whenever a binding changes our webhook programming might change.
	policybindinginformer.Informer().AddEventHandler(handler)

	return func() ([]psbinding.Bindable, error) {
		l, err := policybindinginformer.Lister().List(labels.Everything())
		if err != nil {
			return nil, err
		}
		bl := make([]psbinding.Bindable, 0, len(l))
		for _, elt := range l {
			bl = append(bl, elt)
		}
		return bl, nil
	}

}

func WithContextFactory(ctx context.Context, handler func(types.NamespacedName)) psbinding.BindableContext {
	return func(c context.Context, b psbinding.Bindable) (context.Context, error) {
		pb := b.(*v1alpha2.PolicyPodspecableBinding)
		return v1alpha2.WithBinding(ctx, pb), nil
	}
}
