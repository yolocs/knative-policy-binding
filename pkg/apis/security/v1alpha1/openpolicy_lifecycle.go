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
	policyduck "github.com/yolocs/knative-policy-binding/pkg/apis/duck/v1alpha1"
	"knative.dev/pkg/apis"
)

var openPolicyCondSet = apis.NewLivingConditionSet(
	OpenPolicyConditionReady,
	OpenPolicyConditionConfigMap,
)

const (
	OpenPolicyConditionReady                        = apis.ConditionReady
	OpenPolicyConditionConfigMap apis.ConditionType = "ConfigMapReady"
)

// GetCondition returns the condition currently associated with the given type, or nil.
func (ps *OpenPolicyStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return openPolicyCondSet.Manage(ps).GetCondition(t)
}

// GetTopLevelCondition returns the top level Condition.
func (ps *OpenPolicyStatus) GetTopLevelCondition() *apis.Condition {
	return openPolicyCondSet.Manage(ps).GetTopLevelCondition()
}

// IsReady returns true if the resource is ready overall.
func (ps *OpenPolicyStatus) IsReady() bool {
	return openPolicyCondSet.Manage(ps).IsHappy()
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (ps *OpenPolicyStatus) InitializeConditions() {
	openPolicyCondSet.Manage(ps).InitializeConditions()
}

// MarkReady marks the status ready.
func (ps *OpenPolicyStatus) MarkReady() {
	openPolicyCondSet.Manage(ps).MarkTrue(OpenPolicyConditionReady)
}

// MarkFailed marks the status ready.
func (ps *OpenPolicyStatus) MarkFailed(reason, messageFormat string, messageA ...interface{}) {
	openPolicyCondSet.Manage(ps).MarkFalse(OpenPolicyConditionReady, reason, messageFormat, messageA...)
}

// MarkConfigMapReady marks the config map is ready.
func (ps *OpenPolicyStatus) MarkConfigMapReady(cmName string) {
	ps.ConfigMapName = cmName
	openPolicyCondSet.Manage(ps).MarkTrue(OpenPolicyConditionConfigMap)
}

// MarkConfigMapFailed marks the config map is ready.
func (ps *OpenPolicyStatus) MarkConfigMapFailed(reason, messageFormat string, messageA ...interface{}) {
	openPolicyCondSet.Manage(ps).MarkFalse(OpenPolicyConditionConfigMap, reason, messageFormat, messageA...)
}

// SetDeciderURI sets the decider URI.
func (ps *OpenPolicyStatus) SetDeciderURI(url *apis.URL) {
	if url != nil {
		ps.DeciderURI = url
	}
}

// SetAgentSpec sets the agent spec.
func (ps *OpenPolicyStatus) SetAgentSpec(agentSpec *policyduck.PolicyableAgentSpec) {
	if agentSpec != nil {
		ps.AgentSpec = agentSpec
	}
}
