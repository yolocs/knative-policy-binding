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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck"
	"knative.dev/pkg/tracker"
)

var authorizableBindingCondSet = apis.NewLivingConditionSet(
	AuthorizableBindingConditionReady,
	AuthorizableBindingConditionSubjectResolved,
	AuthorizableBindingConditionPolicyBindingReady,
)

const (
	AuthorizableBindingConditionReady                                 = apis.ConditionReady
	AuthorizableBindingConditionSubjectResolved    apis.ConditionType = "AuthorizableSubjectResolved"
	AuthorizableBindingConditionPolicyBindingReady                    = "PolicyBindingReady"
)

// GetGroupVersionKind returns GroupVersionKind for Triggers
func (ab *AuthorizableBinding) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("AuthorizableBinding")
}

// GetUntypedSpec returns the spec of the Trigger.
func (ab *AuthorizableBinding) GetUntypedSpec() interface{} {
	return ab.Spec
}

// GetBindingStatus implements psbinding.Bindable
func (ab *AuthorizableBinding) GetBindingStatus() duck.BindableStatus {
	return &ab.Status
}

// IsReady returns true if the resource is ready overall.
func (ab *AuthorizableBindingStatus) IsReady() bool {
	return authorizableBindingCondSet.Manage(ab).IsHappy()
}

// SetObservedGeneration implements psbinding.BindableStatus
func (abs *AuthorizableBindingStatus) SetObservedGeneration(gen int64) {
	abs.ObservedGeneration = gen
}

// InitializeConditions populates the AuthorizableBindingStatus's conditions field
// with all of its conditions configured to Unknown.
func (abs *AuthorizableBindingStatus) InitializeConditions() {
	authorizableBindingCondSet.Manage(abs).InitializeConditions()
}

// MarkBindingFailure marks the SinkBinding's Ready condition to False with
// the provided reason and message.
func (abs *AuthorizableBindingStatus) MarkBindingUnavailable(reason, message string) {
	authorizableBindingCondSet.Manage(abs).MarkFalse(AuthorizableBindingConditionReady, reason, message)
}

// MarkBindingReady marks the SinkBinding's Ready condition to True.
func (abs *AuthorizableBindingStatus) MarkBindingAvailable() {
	authorizableBindingCondSet.Manage(abs).MarkTrue(AuthorizableBindingConditionReady)
}

// MarkBindingSubjectResolved marks the subject is resolved.
func (abs *AuthorizableBindingStatus) MarkBindingSubjectResolved(sub *tracker.Reference) {
	authorizableBindingCondSet.Manage(abs).MarkTrue(AuthorizableBindingConditionSubjectResolved)
	abs.ResolvedSubject = sub
}

// MarkBindingSubjectResolvingFaiulre marks subject resolving failure.
func (abs *AuthorizableBindingStatus) MarkBindingSubjectResolvingFaiulre(reason, messageFormat string, messageA ...interface{}) {
	authorizableBindingCondSet.Manage(abs).MarkFalse(AuthorizableBindingConditionSubjectResolved, reason, messageFormat, messageA...)
}

// MarkBindingPolicyReady marks PolicyBinding ready.
func (abs *AuthorizableBindingStatus) MarkBindingPolicyReady() {
	authorizableBindingCondSet.Manage(abs).MarkTrue(AuthorizableBindingConditionPolicyBindingReady)
}

// MarkBindingPolicyFaiulre marks PolicyBinding not ready.
func (abs *AuthorizableBindingStatus) MarkBindingPolicyFaiulre(reason, messageFormat string, messageA ...interface{}) {
	authorizableBindingCondSet.Manage(abs).MarkFalse(AuthorizableBindingConditionPolicyBindingReady, reason, messageFormat, messageA...)
}
