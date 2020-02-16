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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck"
	"knative.dev/pkg/tracker"
)

var PolicyPodspecableBindingCondSet = apis.NewLivingConditionSet(
	PolicyPodspecableBindingConditionReady,
)

const (
	PolicyPodspecableBindingConditionReady = apis.ConditionReady
)

// GetGroupVersionKind returns GroupVersionKind for Triggers
func (pb *PolicyPodspecableBinding) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("PolicyPodspecableBinding")
}

// GetUntypedSpec returns the spec of the Trigger.
func (pb *PolicyPodspecableBinding) GetUntypedSpec() interface{} {
	return pb.Spec
}

// GetSubject implements psbinding.Bindable
func (pb *PolicyPodspecableBinding) GetSubject() tracker.Reference {
	return pb.Spec.Subject
}

// GetBindingStatus implements psbinding.Bindable
func (pb *PolicyPodspecableBinding) GetBindingStatus() duck.BindableStatus {
	return &pb.Status
}

// IsReady returns true if the resource is ready overall.
func (pbs *PolicyPodspecableBindingStatus) IsReady() bool {
	return PolicyPodspecableBindingCondSet.Manage(pbs).IsHappy()
}

// GetTopLevelCondition returns the top level Condition.
func (pbs *PolicyPodspecableBindingStatus) GetTopLevelCondition() *apis.Condition {
	return PolicyPodspecableBindingCondSet.Manage(pbs).GetTopLevelCondition()
}

// SetObservedGeneration implements psbinding.BindableStatus
func (pbs *PolicyPodspecableBindingStatus) SetObservedGeneration(gen int64) {
	pbs.ObservedGeneration = gen
}

// InitializeConditions populates the PolicyPodspecableBindingStatus's conditions field
// with all of its conditions configured to Unknown.
func (pbs *PolicyPodspecableBindingStatus) InitializeConditions() {
	PolicyPodspecableBindingCondSet.Manage(pbs).InitializeConditions()
}

// MarkBindingUnavailable marks the SinkBinding's Ready condition to False with
// the provided reason and message.
func (pbs *PolicyPodspecableBindingStatus) MarkBindingUnavailable(reason, message string) {
	PolicyPodspecableBindingCondSet.Manage(pbs).MarkFalse(PolicyPodspecableBindingConditionReady, reason, message)
}

// MarkBindingAvailable marks the SinkBinding's Ready condition to True.
func (pbs *PolicyPodspecableBindingStatus) MarkBindingAvailable() {
	PolicyPodspecableBindingCondSet.Manage(pbs).MarkTrue(PolicyPodspecableBindingConditionReady)
}
