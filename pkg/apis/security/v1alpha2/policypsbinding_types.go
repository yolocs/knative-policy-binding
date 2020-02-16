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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	duckv1alpha1 "knative.dev/pkg/apis/duck/v1alpha1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PolicyPodspecableBinding is policy binding on podspecable.
type PolicyPodspecableBinding struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PolicyPodspecableBindingSpec   `json:"spec"`
	Status PolicyPodspecableBindingStatus `json:"status"`
}

type PolicyPodspecableBindingSpec struct {
	duckv1alpha1.BindingSpec `json:",inline"`

	// DeciderURI is the decision service endpoint.
	DeciderURI string `json:"deciderURI,omitempty"`

	// AgentSpec is the podspec if an agent is required to inject to the
	// user pod.
	// +optional
	AgentSpec *PolicyAgentSpec `json:"agentSpec,omitempty"`
}

type PolicyAgentSpec struct {
	// Volumes to mount.
	// +optional
	Volumes []corev1.Volume `json:"volumes,omitempty"`

	// Container to inject as the agent sidecar.
	Container corev1.Container `json:"container,omitempty"`
}

// PolicyPodspecableBindingStatus is the status of the binding.
type PolicyPodspecableBindingStatus struct {
	// inherits duck/v1 Status, which currently provides:
	// * ObservedGeneration - the 'Generation' of the Service that was last processed by the controller.
	// * Conditions - the latest available observations of a resource's current state.
	duckv1.Status `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PolicyPodspecableBindingList is a collection of PolicyPodspecableBindings.
type PolicyPodspecableBindingList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PolicyPodspecableBinding `json:"items"`
}
