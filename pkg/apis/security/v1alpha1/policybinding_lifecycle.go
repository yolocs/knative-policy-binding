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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck"
	"knative.dev/pkg/tracker"
)

var policyBindingCondSet = apis.NewLivingConditionSet(
	PolicyBindingConditionReady,
)

const (
	PolicyBindingConditionReady = apis.ConditionReady
)

// GetGroupVersionKind returns GroupVersionKind for Triggers
func (pb *PolicyBinding) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("PolicyBinding")
}

// GetUntypedSpec returns the spec of the Trigger.
func (pb *PolicyBinding) GetUntypedSpec() interface{} {
	return pb.Spec
}

// GetSubject implements psbinding.Bindable
func (pb *PolicyBinding) GetSubject() tracker.Reference {
	return pb.Spec.Subject
}

// GetBindingStatus implements psbinding.Bindable
func (pb *PolicyBinding) GetBindingStatus() duck.BindableStatus {
	return &pb.Status
}

// SetObservedGeneration implements psbinding.BindableStatus
func (pbs *PolicyBindingStatus) SetObservedGeneration(gen int64) {
	pbs.ObservedGeneration = gen
}

// InitializeConditions populates the PolicyBindingStatus's conditions field
// with all of its conditions configured to Unknown.
func (pbs *PolicyBindingStatus) InitializeConditions() {
	policyBindingCondSet.Manage(pbs).InitializeConditions()
}

// MarkBindingUnavailable marks the SinkBinding's Ready condition to False with
// the provided reason and message.
func (pbs *PolicyBindingStatus) MarkBindingUnavailable(reason, message string) {
	policyBindingCondSet.Manage(pbs).MarkFalse(PolicyBindingConditionReady, reason, message)
}

// MarkBindingAvailable marks the SinkBinding's Ready condition to True.
func (pbs *PolicyBindingStatus) MarkBindingAvailable() {
	policyBindingCondSet.Manage(pbs).MarkTrue(PolicyBindingConditionReady)
}
