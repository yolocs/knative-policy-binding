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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck"
)

// +genduck
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Policyable is the duck type policy. This is not a real resource.
type Policyable struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the common policy spec.
	Spec PolicyableSpec `json:"spec,omitempty"`

	// Status is the common policy status.
	Status PolicyableStatus `json:"status,omitempty"`
}

// PolicyableSpec contains the spec of a Policyable object.
type PolicyableSpec struct {
	CheckPayload bool `json:"checkPayload,omitempty"`
}

// PolicyableStatus contains the status of a Policyable object.
type PolicyableStatus struct {
	// DeciderURI is the decision service endpoint.
	DeciderURI *apis.URL `json:"deciderURI,omitempty"`

	// AgentSpec is the podspec if an agent is required to inject to the
	// user pod.
	// +optional
	AgentSpec *PolicyableAgentSpec `json:"agentSpec,omitempty"`
}

type PolicyableAgentSpec struct {
	// Volumes to mount.
	// +optional
	Volumes []corev1.Volume `json:"volumes,omitempty"`

	// Container to inject as the agent sidecar.
	Container corev1.Container `json:"container,omitempty"`
}

var (
	// Verify Policyable resources meet duck contracts.
	_ duck.Populatable   = (*Policyable)(nil)
	_ duck.Implementable = (*Policyable)(nil)
	_ apis.Listable      = (*Policyable)(nil)
)

// Populate implements duck.Populatable
func (p *Policyable) Populate() {
	p.Spec = PolicyableSpec{
		CheckPayload: true,
	}
	p.Status = PolicyableStatus{
		DeciderURI: apis.HTTP("localhost:8085"),
		AgentSpec: &PolicyableAgentSpec{
			Volumes: []corev1.Volume{
				{
					Name: "vol-mount",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			},
			Container: corev1.Container{
				Name:  "agent",
				Image: "image.example.com",
			},
		},
	}
}

// GetFullType implements duck.Implementable
func (s *Policyable) GetFullType() duck.Populatable {
	return &Policyable{}
}

// GetListType implements apis.Listable
func (c *Policyable) GetListType() runtime.Object {
	return &PolicyableList{}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PolicyableList is a list of Policyable resources.
type PolicyableList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Policyable `json:"items"`
}
