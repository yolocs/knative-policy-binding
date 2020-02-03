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
	policyduck "github.com/yolocs/knative-policy-binding/pkg/apis/duck/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventPolicy is a policy defined for eventing.
type EventPolicy struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec of the policy.
	Spec EventPolicySpec `json:"spec,omitempty"`

	// Status of the policy.
	// +optional
	Status EventPolicyStatus `json:"status,omitempty"`
}

var (
	_ apis.Validatable   = (*EventPolicy)(nil)
	_ apis.Defaultable   = (*EventPolicy)(nil)
	_ apis.HasSpec       = (*EventPolicy)(nil)
	_ runtime.Object     = (*EventPolicy)(nil)
	_ kmeta.OwnerRefable = (*EventPolicy)(nil)
)

// EventPolicySpec is the spec for an event policy.
type EventPolicySpec struct {
	// Duck spec.
	policyduck.PolicyableSpec `json:",inline"`

	// Rules define the rules of the policy.
	Rules [][]EventPolicyRule `json:"rules,omitempty"`
}

// EventPolicyRule is a single event policy rule.
type EventPolicyRule struct {
	Name          string `json:"name,omitempty"`
	ExactMatch    string `json:"exactMatch,omitempty"`
	PrefixMatch   string `json:"prefixMatch,omitempty"`
	SuffixMatch   string `json:"suffixMatch,omitempty"`
	ContainsMatch string `json:"containsMatch,omitempty"`
}

// EventPolicyStatus is the status of an EventPolicy.
type EventPolicyStatus struct {
	// inherits duck/v1 Status, which currently provides:
	// * ObservedGeneration - the 'Generation' of the Service that was last processed by the controller.
	// * Conditions - the latest available observations of a resource's current state.
	duckv1.Status `json:",inline"`

	// Duck status.
	policyduck.PolicyableStatus `json:",inline"`

	// The name of underlying open policy.
	OpenPolicyName string `json:"openpolicyName,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventPolicyList is a collection of EventPolicies.
type EventPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EventPolicy `json:"items"`
}

// GetGroupVersionKind returns GroupVersionKind for Triggers
func (p *EventPolicy) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("EventPolicy")
}

// GetUntypedSpec returns the spec of the Trigger.
func (p *EventPolicy) GetUntypedSpec() interface{} {
	return p.Spec
}
