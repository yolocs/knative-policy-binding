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

package resolver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	pkgapisduck "knative.dev/pkg/apis/duck"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/client/injection/ducks/duck/v1/conditions"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/tracker"
)

// SubjectResolver resolves policy subject from Authorizables.
type SubjectResolver struct {
	tracker         tracker.Interface
	informerFactory pkgapisduck.InformerFactory
}

// NewSubjectResolver constructs a new PodspecableResolver with context and a callback
// for a given listableType (Listable) passed to the PodspecableResolver's tracker.
func NewSubjectResolver(ctx context.Context, callback func(types.NamespacedName)) *SubjectResolver {
	ret := &SubjectResolver{}

	ret.tracker = tracker.New(callback, controller.GetTrackerLease(ctx))
	ret.informerFactory = &pkgapisduck.CachedInformerFactory{
		Delegate: &pkgapisduck.EnqueueInformerFactory{
			Delegate:     conditions.Get(ctx),
			EventHandler: controller.HandleAll(ret.tracker.OnChanged),
		},
	}

	return ret
}

// ResolveFromRef resolves podspecable from the reference.
func (r *SubjectResolver) ResolveFromRef(ref *corev1.ObjectReference, parent interface{}) (*tracker.Reference, error) {
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

	kr, ok := obj.(*duckv1.KResource)
	if !ok {
		return nil, fmt.Errorf("%+v (%T) is not a KResource", ref, ref)
	}

	// Parse the annotation?
	subRaw, ok := kr.Annotations["security.knative.dev/authorizableOn"]
	if !ok {
		return nil, errors.New("the reference is not an authorizable; expecting annotation 'security.knative.dev/authorizableOn'")
	}
	// Handle this special case where the object itself is already podspecable.
	if subRaw == "self" {
		return &tracker.Reference{
			APIVersion: kr.APIVersion,
			Kind:       kr.Kind,
			Name:       kr.Name,
			Namespace:  kr.Namespace,
			// Populate labels in case the binding implementation only supports labels.
			Selector: &metav1.LabelSelector{
				MatchLabels: kr.GetLabels(),
			},
		}, nil
	}

	var t tracker.Reference
	if err := json.Unmarshal([]byte(subRaw), &t); err != nil {
		return nil, fmt.Errorf("the reference doesn't provide a valid podspecable subject in annotation 'security.knative.dev/authorizableOn': %w", err)
	}
	// Some defaulting.
	if t.Namespace == "" {
		t.Namespace = ref.Namespace
	}

	return &t, nil
}
