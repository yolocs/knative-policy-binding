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

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/kmeta"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HTTPPolicy is a policy defined for HTTP requests.
type HTTPPolicy struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec HTTPPolicySpec `json:"spec"`
}

var (
	_ apis.Validatable   = (*HTTPPolicy)(nil)
	_ apis.Defaultable   = (*HTTPPolicy)(nil)
	_ apis.HasSpec       = (*HTTPPolicy)(nil)
	_ runtime.Object     = (*HTTPPolicy)(nil)
	_ kmeta.OwnerRefable = (*HTTPPolicy)(nil)
)

type HTTPPolicySpec struct {
	JWT   JWTSpec    `json:"jwt,omitempty"`
	Rules []RuleSpec `json:"rules,omitempty"`
}

type JWTSpec struct {
	JwksURI      string        `json:"jwksUri,omitempty"`
	Jwks         string        `json:"jwks,omitempty"`
	JwtHeader    string        `json:"jwtHead,omitempty"`
	TriggerRules []TriggerRule `json:"triggerRules,omitempty"`
}

type RuleSpec struct {
	Auth       RequestAuth     `json:"auth,omitempty"`
	Headers    []KeyValueMatch `json:"headers,omitempty"`
	Operations []Operation     `json:"operations,omitempty"`
}

type Operation struct {
	Hosts   []string `json:"hosts,omitempty"`
	Paths   []string `json:"paths,omitempty"`
	Methods []string `json:"methods,omitempty"`
}

type KeyValueMatch struct {
	Key    string   `json:"key,omitempty"`
	Values []string `json:"values,omitempty"`
}

type RequestAuth struct {
	Principals []string        `json:"principals,omitempty"`
	Claims     []KeyValueMatch `json:"claims,omitempty"`
}

type TriggerRule struct {
	ExcludePaths []string `json:"excludePaths,omitempty"`
	IncludePaths []string `json:"includePaths,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HTTPPolicyList is a collection of OpenPolicies.
type HTTPPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HTTPPolicy `json:"items"`
}

// GetGroupVersionKind returns GroupVersionKind for Triggers
func (p *HTTPPolicy) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("HTTPPolicy")
}

// GetUntypedSpec returns the spec of the Trigger.
func (p *HTTPPolicy) GetUntypedSpec() interface{} {
	return p.Spec
}
