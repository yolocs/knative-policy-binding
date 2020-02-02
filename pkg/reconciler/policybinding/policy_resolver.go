/*
Copyright 2019 The Knative Authors

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

package policybinding

import (
	"context"
	"errors"
	"fmt"

	policyduck "github.com/yolocs/knative-policy-binding/pkg/apis/duck/v1alpha1"
	"github.com/yolocs/knative-policy-binding/pkg/client/injection/ducks/duck/v1alpha1/policyable"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	pkgapisduck "knative.dev/pkg/apis/duck"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/tracker"
)

// PolicyResolver resolves policies.
type PolicyResolver struct {
	tracker         tracker.Interface
	informerFactory pkgapisduck.InformerFactory
}

// NewPolicyResolver constructs a new PolicyResolver with context and a callback
// for a given listableType (Listable) passed to the PolicyResolver's tracker.
func NewPolicyResolver(ctx context.Context, callback func(types.NamespacedName)) *PolicyResolver {
	ret := &PolicyResolver{}

	ret.tracker = tracker.New(callback, controller.GetTrackerLease(ctx))
	ret.informerFactory = &pkgapisduck.CachedInformerFactory{
		Delegate: &pkgapisduck.EnqueueInformerFactory{
			Delegate:     policyable.Get(ctx),
			EventHandler: controller.HandleAll(ret.tracker.OnChanged),
		},
	}

	return ret
}

// PolicyFromRef resolves policy.
func (r *PolicyResolver) PolicyFromRef(ref *corev1.ObjectReference, parent interface{}) (*policyduck.Policyable, error) {
	if ref == nil {
		return nil, errors.New("ref is nil")
	}

	if err := r.tracker.Track(*ref, parent); err != nil {
		return nil, fmt.Errorf("failed to track %+v: %v", ref, err)
	}

	gvr, _ := meta.UnsafeGuessKindToResource(ref.GroupVersionKind())
	_, lister, err := r.informerFactory.Get(gvr)
	if err != nil {
		return nil, fmt.Errorf("failed to get lister for %+v: %v", gvr, err)
	}

	obj, err := lister.ByNamespace(ref.Namespace).Get(ref.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get ref %+v: %v", ref, err)
	}

	policy, ok := obj.(*policyduck.Policyable)
	if !ok {
		return nil, fmt.Errorf("%+v (%T) is not a Policyable", ref, ref)
	}

	return policy.DeepCopy(), nil
}
