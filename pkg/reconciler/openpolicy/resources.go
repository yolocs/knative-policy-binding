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

package openpolicy

import (
	policyduck "github.com/yolocs/knative-policy-binding/pkg/apis/duck/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/logging"
)

const (
	policyFileName = "policy.rego"
)

func MakeDeciderURL() *apis.URL {
	return apis.HTTP("localhost:8090")
}

func MakeAgentSpec(agentImage, cmName, toRemove string) *policyduck.PolicyableAgentSpec {
	cfg, _ := logging.NewConfigFromMap(nil)
	return &policyduck.PolicyableAgentSpec{
		Volumes: []corev1.Volume{
			{
				Name: "open-policy",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: cmName,
						},
					},
				},
			},
		},
		Container: corev1.Container{
			Name:  "kn-policy-agent",
			Image: agentImage,
			Env: []corev1.EnvVar{
				{
					Name:  "POLICY_PATH",
					Value: "/var/run/knative/security/" + policyFileName,
				},
				{
					Name:  "AGENT_PORT",
					Value: "8090",
				},
				{
					Name:  "AGENT_LOGGING_CONFIG",
					Value: cfg.LoggingConfig,
				},
				{
					Name:  "AGENT_LOGGING_LEVEL",
					Value: "debug",
				},
				{
					Name:  "AGENT_TEST",
					Value: toRemove,
				},
			},
			Ports: []corev1.ContainerPort{
				{
					Name:          "http",
					ContainerPort: 8090,
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "open-policy",
					MountPath: "/var/run/knative/security",
					ReadOnly:  true,
				},
			},
		},
	}
}
