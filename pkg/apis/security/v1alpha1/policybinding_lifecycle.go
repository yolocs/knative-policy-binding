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
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/logging"
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

// Do implements psbinding.Bindable
func (pb *PolicyBinding) Do(ctx context.Context, ps *duckv1.WithPod) {
	// First undo.
	pb.Undo(ctx, ps)

	pstatus := GetPolicyStatus(ctx)
	if pstatus == nil {
		logging.FromContext(ctx).Error(fmt.Sprintf("No policy status associated with context for %+v", pb))
		return
	}

	spec := ps.Spec.Template.Spec
	for i := range spec.Containers {
		spec.Containers[i].Env = append(spec.Containers[i].Env, corev1.EnvVar{
			Name:  "K_POLICY_DECIDER",
			Value: pstatus.DeciderURI.String(),
		})
	}

	if pstatus.AgentSpec != nil {
		spec.Volumes = append(spec.Volumes, pstatus.AgentSpec.Volumes...)
		spec.Containers = append(spec.Containers, pstatus.AgentSpec.Container)
	}

	ps.Spec.Template.Spec = spec
}

// Undo implements psbinding.Bindable
func (pb *PolicyBinding) Undo(ctx context.Context, ps *duckv1.WithPod) {
	pstatus := GetPolicyStatus(ctx)
	if pstatus == nil {
		logging.FromContext(ctx).Error(fmt.Sprintf("No policy status associated with context for %+v", pb))
		return
	}

	spec := ps.Spec.Template.Spec
	for i, c := range spec.Containers {
		for j, ev := range c.Env {
			if ev.Name == "K_POLICY_DECIDER" {
				spec.Containers[i].Env = append(spec.Containers[i].Env[:j], spec.Containers[i].Env[j+1:]...)
				break
			}
		}
	}

	if pstatus.AgentSpec == nil {
		return
	}

	// Remove volumes.
	for _, v := range pstatus.AgentSpec.Volumes {
		for j, ev := range spec.Volumes {
			if ev.Name == v.Name {
				spec.Volumes = append(spec.Volumes[:j], spec.Volumes[j+1:]...)
			}
		}
	}

	// Remove agent sidecar.
	for i, c := range spec.Containers {
		if c.Name == pstatus.AgentSpec.Container.Name {
			spec.Containers = append(spec.Containers[:i], spec.Containers[i+1:]...)
		}
	}

	ps.Spec.Template.Spec = spec
}
