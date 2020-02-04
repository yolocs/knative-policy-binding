/*
Copyright 2020 The Knative Authors.

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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/tracker"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AuthorizableBinding defines the binding to any authorization object.
type AuthorizableBinding struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AuthorizableBindingSpec   `json:"spec"`
	Status AuthorizableBindingStatus `json:"status"`
}

var (
	_ apis.Validatable   = (*AuthorizableBinding)(nil)
	_ apis.Defaultable   = (*AuthorizableBinding)(nil)
	_ apis.HasSpec       = (*AuthorizableBinding)(nil)
	_ runtime.Object     = (*AuthorizableBinding)(nil)
	_ kmeta.OwnerRefable = (*AuthorizableBinding)(nil)
)

// AuthorizableBindingSpec is the spec for an authorizable binding.
type AuthorizableBindingSpec struct {
	// Subject is the authorizable object to bind the policy to.
	Subject *corev1.ObjectReference `json:"subject"`

	// Policy is the policy reference to bind.
	Policy *corev1.ObjectReference `json:"policy"`
}

// AuthorizableBindingStatus is the status of the binding.
type AuthorizableBindingStatus struct {
	// inherits duck/v1 Status, which currently provides:
	// * ObservedGeneration - the 'Generation' of the Service that was last processed by the controller.
	// * Conditions - the latest available observations of a resource's current state.
	duckv1.Status `json:",inline"`

	// ResolvedSubject is resolved policy subject.
	ResolvedSubject *tracker.Reference `json:"resolvedSubject,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AuthorizableBindingList is a collection of AuthorizableBindings.
type AuthorizableBindingList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AuthorizableBinding `json:"items"`
}
