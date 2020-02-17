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

package opabinding

import (
	"context"
	"fmt"

	"github.com/yolocs/knative-policy-binding/pkg/apis/security/v1alpha2"
	securitylisters "github.com/yolocs/knative-policy-binding/pkg/client/listers/security/v1alpha2"
	"github.com/yolocs/knative-policy-binding/pkg/opa"
	"github.com/yolocs/knative-policy-binding/pkg/reconciler"
	"github.com/yolocs/knative-policy-binding/pkg/reconciler/internal/resolver"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	duckv1alpha1 "knative.dev/pkg/apis/duck/v1alpha1"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
	"knative.dev/pkg/tracker"
)

const (
	bindingReconcileError     = "HTTPPolicyBindingReconcileError"
	bindingReconciled         = "HTTPPolicyBindingReconciled"
	bindingClassAnnotationKey = "security.knative.dev/binding.class"
	bindingClass              = "opa"
)

type Reconciler struct {
	*reconciler.Base

	policybindingLister securitylisters.HTTPPolicyBindingLister
	policyLister        securitylisters.HTTPPolicyLister
	psbindingLister     securitylisters.PolicyPodspecableBindingLister
	configmapLister     corev1listers.ConfigMapLister

	subjectResolver *resolver.SubjectResolver
	policyTracker   tracker.Interface

	agentImage string
}

func (r *Reconciler) ReconcileKind(ctx context.Context, b *v1alpha2.HTTPPolicyBinding) pkgreconciler.Event {
	if b.GetAnnotations()[bindingClassAnnotationKey] != bindingClass {
		logging.FromContext(ctx).Info("Not reconciling binding, cause it's not mine", zap.String("HTTPPolicyBinding", b.Name))
		return nil
	}

	logging.FromContext(ctx).Debug("Reconciling", zap.Any("HTTPPolicyBinding", b))
	b.Status.InitializeConditions()
	b.Status.ObservedGeneration = b.Generation

	sub, err := r.subjectResolver.ResolveFromRef(b.Spec.Subject, b)
	if err != nil {
		logging.FromContext(ctx).Error("Problem resolving binding subject", zap.Error(err))
		b.Status.MarkBindingSubjectResolvingFaiulre("SubjectResolvingFailure", "%v", err)
		return fmt.Errorf("Failed to reconcile HTTP policy binding: %w", err)
	}
	if sub.Selector == nil || len(sub.Selector.MatchLabels) == 0 {
		logging.FromContext(ctx).Error("Resolved binding target is not a label selector")
		b.Status.MarkBindingSubjectResolvingFaiulre("SubjectNotLabelSelector", "Istio AuthorizationPolicy can only select workload by labels")
		return fmt.Errorf("Failed to reconcile HTTP policy binding: %w", err)
	}
	b.Status.MarkBindingSubjectResolved(sub)

	p, err := r.policyLister.HTTPPolicies(b.Spec.Policy.Namespace).Get(b.Spec.Policy.Name)
	if err != nil {
		logging.FromContext(ctx).Error("Problem getting policy", zap.Error(err))
		b.Status.MarkBindingUnavailable("GetPolicyFailure", err.Error())
		return fmt.Errorf("Failed to get the referencing policy: %w", err)
	}
	r.policyTracker.Track(*b.Spec.Policy, b)

	if err := r.reconcileConfigMap(ctx, b, p); err != nil {
		logging.FromContext(ctx).Error("Problem reconciling OPA policy configmap", zap.Error(err))
		b.Status.MarkBindingUnavailable("ConfigMapFailure", err.Error())
		return fmt.Errorf("Failed to reconcile OPA policy configmap: %w", err)
	}

	pb, err := r.reconcilePodspecableBinding(ctx, sub, p, b)
	if err != nil {
		logging.FromContext(ctx).Error("Problem reconciling policy podspecable binding", zap.Error(err))
		b.Status.MarkBindingUnavailable("PolicyPodspecableBindingFailure", err.Error())
		return fmt.Errorf("Failed to reconcile PolicyPodspecablebinding: %w", err)
	}

	if pb.Status.IsReady() {
		b.Status.MarkBindingAvailable()
	} else {
		cond := pb.Status.GetTopLevelCondition()
		if cond != nil {
			b.Status.MarkBindingUnavailable(cond.Reason, cond.Message)
		} else {
			b.Status.MarkBindingUnavailable("PolicyPodspecableBindingNotReady", "")
		}
	}

	b.Status.MarkBindingAvailable()
	return nil
}

func (r *Reconciler) reconcileConfigMap(ctx context.Context, b *v1alpha2.HTTPPolicyBinding, p *v1alpha2.HTTPPolicy) pkgreconciler.Event {
	m := policyToRego(&p.Spec)
	desired := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            b.Name,
			Namespace:       b.Namespace,
			OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(b)},
		},
		Data: map[string]string{
			"policy.rego": m,
		},
	}
	cm, err := r.configmapLister.ConfigMaps(b.Namespace).Get(b.Name)
	if apierrs.IsNotFound(err) {
		cm, err = r.KubeClientSet.CoreV1().ConfigMaps(b.Namespace).Create(desired)
		if err != nil {
			return fmt.Errorf("failed to create configmap: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to get configmap: %w", err)
	}

	if !equality.Semantic.DeepDerivative(desired.Data, cm.Data) {
		// Don't modify the informers copy.
		cp := cm.DeepCopy()
		cp.Data = desired.Data
		cm, err = r.KubeClientSet.CoreV1().ConfigMaps(cp.Namespace).Update(cp)
		if err != nil {
			return fmt.Errorf("failed to update configmap: %w", err)
		}
	}

	return nil
}

func (r *Reconciler) reconcilePodspecableBinding(
	ctx context.Context, sub *tracker.Reference, p *v1alpha2.HTTPPolicy, b *v1alpha2.HTTPPolicyBinding) (*v1alpha2.PolicyPodspecableBinding, pkgreconciler.Event) {
	desired := &v1alpha2.PolicyPodspecableBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:            b.Name,
			Namespace:       b.Namespace,
			OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(b)},
			Annotations:     map[string]string{"security.knative.dev/policyGeneration": fmt.Sprintf("%d", p.ObjectMeta.Generation)},
		},
		Spec: v1alpha2.PolicyPodspecableBindingSpec{
			BindingSpec: duckv1alpha1.BindingSpec{
				Subject: *sub,
			},
			DeciderURI: "http://localhost:8090",
			AgentSpec:  r.genAgentSpec(b),
		},
	}
	pb, err := r.psbindingLister.PolicyPodspecableBindings(desired.Namespace).Get(desired.Name)
	if apierrs.IsNotFound(err) {
		pb, err = r.SecurityClientSet.SecurityV1alpha2().PolicyPodspecableBindings(desired.Namespace).Create(desired)
		if err != nil {
			return nil, fmt.Errorf("Failed to create PolicyPodspecableBinding: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("Failed to get PolicyPodspecableBinding: %w", err)
	}

	if !equality.Semantic.DeepDerivative(desired.Spec, pb.Spec) || desired.GetAnnotations()["security.knative.dev/policyGeneration"] != pb.GetAnnotations()["security.knative.dev/policyGeneration"] {
		// Don't modify the informers copy.
		cp := pb.DeepCopy()
		cp.Spec = desired.Spec
		cp.GetAnnotations()["security.knative.dev/policyGeneration"] = desired.GetAnnotations()["security.knative.dev/policyGeneration"]
		pb, err = r.SecurityClientSet.SecurityV1alpha2().PolicyPodspecableBindings(cp.Namespace).Update(cp)
		if err != nil {
			return nil, fmt.Errorf("Failed to update PolicyPodspecableBinding: %w", err)
		}
	}
	return pb, nil
}

func policyToRego(spec *v1alpha2.HTTPPolicySpec) string {
	pbuilder := opa.NewPolicyBuilder()
	for _, rule := range spec.Rules {
		// ignore others for now
		rbuilder := pbuilder.NewRule()
		for _, h := range rule.Headers {
			rbuilder.AppendOneOf(fmt.Sprintf("input.httpRequest.header[%q][_]", h.Key), h.Values)
		}
		for _, op := range rule.Operations {
			rbuilder.AppendOneOf("input.httpRequest.method", op.Methods)
			rbuilder.AppendOneOf("input.httpRequest.host", op.Hosts)
			rbuilder.AppendOneOf("input.httpRequest.path", op.Paths)
		}
	}
	return pbuilder.String()
}

func (r *Reconciler) genAgentSpec(b *v1alpha2.HTTPPolicyBinding) *v1alpha2.PolicyAgentSpec {
	cfg, _ := logging.NewConfigFromMap(nil)
	return &v1alpha2.PolicyAgentSpec{
		Volumes: []corev1.Volume{
			{
				Name: "open-policy",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: b.Name,
						},
					},
				},
			},
		},
		Container: corev1.Container{
			Name:  "kn-policy-agent",
			Image: r.agentImage,
			Env: []corev1.EnvVar{
				{
					Name:  "POLICY_PATH",
					Value: "/var/run/knative/security/policy.rego",
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
