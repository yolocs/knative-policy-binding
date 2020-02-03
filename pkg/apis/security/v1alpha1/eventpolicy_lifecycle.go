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
	"knative.dev/pkg/apis"
)

var eventPolicyCondSet = apis.NewLivingConditionSet(
	EventPolicyConditionReady,
	EventPolicyConditionOpenPolicy,
)

const (
	EventPolicyConditionReady                         = apis.ConditionReady
	EventPolicyConditionOpenPolicy apis.ConditionType = "OpenPolicyReady"
)

// GetCondition returns the condition currently associated with the given type, or nil.
func (ps *EventPolicyStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return eventPolicyCondSet.Manage(ps).GetCondition(t)
}

// GetTopLevelCondition returns the top level Condition.
func (ps *EventPolicyStatus) GetTopLevelCondition() *apis.Condition {
	return eventPolicyCondSet.Manage(ps).GetTopLevelCondition()
}

// IsReady returns true if the resource is ready overall.
func (ps *EventPolicyStatus) IsReady() bool {
	return eventPolicyCondSet.Manage(ps).IsHappy()
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (ps *EventPolicyStatus) InitializeConditions() {
	eventPolicyCondSet.Manage(ps).InitializeConditions()
}

// MarkFailed marks the status ready.
func (ps *EventPolicyStatus) MarkFailed(reason, messageFormat string, messageA ...interface{}) {
	eventPolicyCondSet.Manage(ps).MarkFalse(EventPolicyConditionReady, reason, messageFormat, messageA...)
}

// MarkOpenPolicyFailed marks the open policy not ready.
func (ps *EventPolicyStatus) MarkOpenPolicyFailed(reason, messageFormat string, messageA ...interface{}) {
	openPolicyCondSet.Manage(ps).MarkFalse(EventPolicyConditionOpenPolicy, reason, messageFormat, messageA...)
}

// PropagateFromOpenPolicyStatus fills the status with open policy status.
func (ps *EventPolicyStatus) PropagateFromOpenPolicyStatus(opName string, status *OpenPolicyStatus) {
	if status == nil {
		return
	}

	ps.OpenPolicyName = opName
	ps.PolicyableStatus = status.PolicyableStatus
	if status.IsReady() {
		eventPolicyCondSet.Manage(ps).MarkTrue(EventPolicyConditionReady)
		eventPolicyCondSet.Manage(ps).MarkTrue(EventPolicyConditionOpenPolicy)
	} else {
		cond := status.GetTopLevelCondition()
		reason, message := "", ""
		if cond != nil {
			reason = cond.Reason
			message = cond.Message
		}
		eventPolicyCondSet.Manage(ps).MarkFalse(EventPolicyConditionReady, reason, message)
		eventPolicyCondSet.Manage(ps).MarkFalse(EventPolicyConditionOpenPolicy, reason, message)
	}
}
