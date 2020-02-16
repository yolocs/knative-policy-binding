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
	"context"
	"fmt"
	"strings"

	jsonpatch "gomodules.xyz/jsonpatch/v2"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis/duck"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/logging"
)

// Do implements psbinding.Bindable
func (pb *PolicyPodspecableBinding) Do(ctx context.Context, ps *duckv1.WithPod) duck.JSONPatch {
	patch := pb.Undo(ctx, ps)

	binding := GetBinding(ctx)
	if binding == nil {
		logging.FromContext(ctx).Error(fmt.Sprintf("No binding associated with context for %+v", pb))
		return nil
	}

	envs := []corev1.EnvVar{
		{
			Name:  "K_POLICY_DECIDER",
			Value: binding.Spec.DeciderURI,
		},
	}
	patch = append(patch, addEnvs(ps, envs)...)
	// patch = append(patch, addAnnotation(ps, "security.knative.dev/policyGeneration", fmt.Sprintf("%d", policy.Generation))...)

	if binding.Spec.AgentSpec == nil {
		return patch
	}

	patch = append(patch, addVolumes(ps, binding.Spec.AgentSpec.Volumes)...)
	patch = append(patch, addContainer(ps, binding.Spec.AgentSpec.Container)...)
	return patch
}

// Undo implements psbinding.Bindable
func (pb *PolicyPodspecableBinding) Undo(ctx context.Context, ps *duckv1.WithPod) duck.JSONPatch {
	binding := GetBinding(ctx)
	if binding == nil {
		logging.FromContext(ctx).Error(fmt.Sprintf("No binding associated with context for %+v", pb))
		return nil
	}

	var patch duck.JSONPatch
	// if ps.Spec.Template.Annotations != nil {
	// 	if _, ok := ps.Spec.Template.Annotations["security.knative.dev/policyGeneration"]; ok {
	// 		patch = append(patch, removeAnnotation(ps, "security.knative.dev/policyGeneration")...)
	// 	}
	// }

	envs := []string{"K_POLICY_DECIDER"}

	// This has problem when previously there is agent spec and then removed.
	if binding.Spec.AgentSpec == nil {
		return append(patch, removeEnvs(ps, envs)...)
	}

	patch = append(patch, removeVolumes(ps, binding.Spec.AgentSpec.Volumes)...)
	patch = append(patch, removeContainer(ps, binding.Spec.AgentSpec.Container.Name)...)
	return append(patch, removeEnvs(ps, envs)...)
}

func removeAnnotation(ps *duckv1.WithPod, key string) (patch duck.JSONPatch) {
	delete(ps.Spec.Template.Annotations, key)
	patch = append(patch, jsonpatch.Operation{
		Operation: "remove",
		Path:      "/spec/template/metadata/annotations/" + strings.ReplaceAll(key, "/", "~1"),
	})
	return patch
}

func addAnnotation(ps *duckv1.WithPod, key, value string) (patch duck.JSONPatch) {
	if ps.Spec.Template.Annotations == nil {
		ps.Spec.Template.Annotations = map[string]string{key: value}
		patch = append(patch, jsonpatch.Operation{
			Operation: "add",
			Path:      "/spec/template/metadata/annotations",
			Value:     ps.Spec.Template.Annotations,
		})
		return patch
	}

	ps.Spec.Template.Annotations[key] = value
	patch = append(patch, jsonpatch.Operation{
		Operation: "add",
		Path:      "/spec/template/metadata/annotations/" + strings.ReplaceAll(key, "/", "~1"),
		Value:     value,
	})
	return patch
}

func removeContainer(ps *duckv1.WithPod, containerName string) (patch duck.JSONPatch) {
	spec := ps.Spec.Template.Spec
	for i, c := range spec.Containers {
		if c.Name == containerName {
			patch = append(patch, jsonpatch.Operation{
				Operation: "remove",
				Path:      fmt.Sprintf("%v/%v", "/spec/template/spec/containers", i),
			})
			spec.Containers = append(spec.Containers[:i], spec.Containers[i+1:]...)
			break
		}
	}
	ps.Spec.Template.Spec = spec
	return patch
}

func removeVolumes(ps *duckv1.WithPod, vols []corev1.Volume) (patch duck.JSONPatch) {
	spec := ps.Spec.Template.Spec
	for _, vn := range vols {
		for i, v := range spec.Volumes {
			if v.Name == vn.Name {
				spec.Volumes = append(spec.Volumes[:i], spec.Volumes[i+1:]...)
				patch = append(patch, jsonpatch.Operation{
					Operation: "remove",
					Path:      fmt.Sprintf("/spec/template/spec/volumes/%d", i),
				})
			}
		}
	}
	ps.Spec.Template.Spec = spec
	return patch
}

func removeEnvs(ps *duckv1.WithPod, envNames []string) (patch duck.JSONPatch) {
	spec := ps.Spec.Template.Spec
	changedIndex := []int{}
	for _, en := range envNames {
		for i, c := range spec.Containers {
			startLen := len(spec.Containers[i].Env)
			for j, ev := range c.Env {
				if ev.Name == en {
					spec.Containers[i].Env = append(spec.Containers[i].Env[:j], spec.Containers[i].Env[j+1:]...)
				}
			}
			if startLen != len(spec.Containers[i].Env) {
				changedIndex = append(changedIndex, i)
			}
		}
	}
	for _, i := range changedIndex {
		patch = append(patch, jsonpatch.Operation{
			Operation: "replace",
			Path:      fmt.Sprintf("/spec/template/spec/containers/%v/env", i),
			Value:     spec.Containers[i].Env,
		})
	}
	ps.Spec.Template.Spec = spec
	return patch
}

func addContainer(ps *duckv1.WithPod, c corev1.Container) (patch duck.JSONPatch) {
	var value interface{}
	value = c
	path := "/spec/template/spec/containers"
	if len(ps.Spec.Template.Spec.Containers) == 0 {
		value = []corev1.Container{c}
	} else {
		path += "/-"
	}
	patch = append(patch, jsonpatch.Operation{
		Operation: "add",
		Path:      path,
		Value:     value,
	})
	ps.Spec.Template.Spec.Containers = append(ps.Spec.Template.Spec.Containers, c)
	return patch
}

func addVolumes(ps *duckv1.WithPod, vols []corev1.Volume) (patch duck.JSONPatch) {
	first := len(ps.Spec.Template.Spec.Volumes) == 0
	var value interface{}
	for _, v := range vols {
		value = v
		path := "/spec/template/spec/volumes"
		if first {
			first = false
			value = []corev1.Volume{v}
		} else {
			path += "/-"
		}
		patch = append(patch, jsonpatch.Operation{
			Operation: "add",
			Path:      path,
			Value:     value,
		})
	}
	ps.Spec.Template.Spec.Volumes = append(ps.Spec.Template.Spec.Volumes, vols...)
	return patch
}

func addEnvs(ps *duckv1.WithPod, envs []corev1.EnvVar) (patch duck.JSONPatch) {
	for i, c := range ps.Spec.Template.Spec.Containers {
		ps.Spec.Template.Spec.Containers[i].Env = append(c.Env, envs...)
		patch = append(patch, jsonpatch.Operation{
			Operation: "replace",
			Path:      fmt.Sprintf("/spec/template/spec/containers/%v/env", i),
			Value:     ps.Spec.Template.Spec.Containers[i].Env,
		})
	}
	return patch
}
