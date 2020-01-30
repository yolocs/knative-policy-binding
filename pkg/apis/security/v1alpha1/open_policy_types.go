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
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenPolicy is a policy defined in OPA rego language.
type OpenPolicy struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec of the policy.
	Spec OpenPolicySpec `json:"spec,omitempty"`

	// Status of the policy.
	// +optional
	Status OpenPolicyStatus `json:"status,omitempty"`
}

// OpenPolicySpec is open policy spec.
type OpenPolicySpec struct {
	// Duck spec.
	policyduck.PolicyableSpec `json:",inline"`

	// Rule is the rule defined in rego.
	Rule string `json:"rule,omitempty"`
}

// OpenPolicyStatus is open policy status.
type OpenPolicyStatus struct {
	// Duck status.
	policyduck.PolicyableStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenPolicyList is a collection of OpenPolicies.
type OpenPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenPolicy `json:"items"`
}

// GetGroupVersionKind returns GroupVersionKind for Triggers
func (p *OpenPolicy) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("OpenPolicy")
}

// GetUntypedSpec returns the spec of the Trigger.
func (p *OpenPolicy) GetUntypedSpec() interface{} {
	return p.Spec
}
