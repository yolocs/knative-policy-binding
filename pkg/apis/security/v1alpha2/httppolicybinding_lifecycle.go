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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck"
	"knative.dev/pkg/tracker"
)

var httpPolicyBindingCondSet = apis.NewLivingConditionSet(
	HTTPPolicyBindingConditionReady,
)

const (
	HTTPPolicyBindingConditionReady                          = apis.ConditionReady
	HTTPPolicyAuthorizableSubjectResolved apis.ConditionType = "AuthorizableSubjectResolved"
)

// GetGroupVersionKind returns GroupVersionKind for Triggers
func (pb *HTTPPolicyBinding) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("HTTPPolicyBinding")
}

// GetUntypedSpec returns the spec of the Trigger.
func (pb *HTTPPolicyBinding) GetUntypedSpec() interface{} {
	return pb.Spec
}

// GetSubject implements psbinding.Bindable
func (pb *HTTPPolicyBinding) GetSubject() *corev1.ObjectReference {
	return pb.Spec.Subject
}

// GetBindingStatus implements psbinding.Bindable
func (pb *HTTPPolicyBinding) GetBindingStatus() duck.BindableStatus {
	return &pb.Status
}

// IsReady returns true if the resource is ready overall.
func (pbs *HTTPPolicyBindingStatus) IsReady() bool {
	return httpPolicyBindingCondSet.Manage(pbs).IsHappy()
}

// GetTopLevelCondition returns the top level Condition.
func (pbs *HTTPPolicyBindingStatus) GetTopLevelCondition() *apis.Condition {
	return httpPolicyBindingCondSet.Manage(pbs).GetTopLevelCondition()
}

// SetObservedGeneration implements psbinding.BindableStatus
func (pbs *HTTPPolicyBindingStatus) SetObservedGeneration(gen int64) {
	pbs.ObservedGeneration = gen
}

// InitializeConditions populates the HTTPPolicyBindingStatus's conditions field
// with all of its conditions configured to Unknown.
func (pbs *HTTPPolicyBindingStatus) InitializeConditions() {
	httpPolicyBindingCondSet.Manage(pbs).InitializeConditions()
}

// MarkBindingUnavailable marks the SinkBinding's Ready condition to False with
// the provided reason and message.
func (pbs *HTTPPolicyBindingStatus) MarkBindingUnavailable(reason, message string) {
	httpPolicyBindingCondSet.Manage(pbs).MarkFalse(HTTPPolicyBindingConditionReady, reason, message)
}

// MarkBindingAvailable marks the SinkBinding's Ready condition to True.
func (pbs *HTTPPolicyBindingStatus) MarkBindingAvailable() {
	httpPolicyBindingCondSet.Manage(pbs).MarkTrue(HTTPPolicyBindingConditionReady)
}

// MarkBindingSubjectResolved marks the subject is resolved.
func (pbs *HTTPPolicyBindingStatus) MarkBindingSubjectResolved(sub *tracker.Reference) {
	httpPolicyBindingCondSet.Manage(pbs).MarkTrue(HTTPPolicyAuthorizableSubjectResolved)
	pbs.ResolvedSubject = sub
}

// MarkBindingSubjectResolvingFaiulre marks subject resolving failure.
func (pbs *HTTPPolicyBindingStatus) MarkBindingSubjectResolvingFaiulre(reason, messageFormat string, messageA ...interface{}) {
	httpPolicyBindingCondSet.Manage(pbs).MarkFalse(HTTPPolicyAuthorizableSubjectResolved, reason, messageFormat, messageA...)
}
