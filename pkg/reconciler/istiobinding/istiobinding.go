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

package istiobinding

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
	"knative.dev/pkg/tracker"

	"github.com/yolocs/knative-policy-binding/pkg/apis/security/v1alpha2"
	istioclientset "github.com/yolocs/knative-policy-binding/pkg/client/istio/clientset/versioned"
	istiolisters "github.com/yolocs/knative-policy-binding/pkg/client/istio/listers/security/v1beta1"
	securitylisters "github.com/yolocs/knative-policy-binding/pkg/client/listers/security/v1alpha2"
	"github.com/yolocs/knative-policy-binding/pkg/reconciler"
	"github.com/yolocs/knative-policy-binding/pkg/reconciler/internal/resolver"
	istiosecurityv1beta1 "istio.io/api/security/v1beta1"
	istiotypev1beta1 "istio.io/api/type/v1beta1"
	istiov1beta1 "istio.io/client-go/pkg/apis/security/v1beta1"
)

const (
	bindingReconcileError     = "HTTPPolicyBindingReconcileError"
	bindingReconciled         = "HTTPPolicyBindingReconciled"
	bindingClassAnnotationKey = "security.knative.dev/binding.class"
	bindingClass              = "istio"
)

type Reconciler struct {
	*reconciler.Base

	policybindingLister securitylisters.HTTPPolicyBindingLister
	policyLister        securitylisters.HTTPPolicyLister
	istioauthzLister    istiolisters.AuthorizationPolicyLister

	istioClientSet istioclientset.Interface

	subjectResolver *resolver.SubjectResolver
	policyTracker   tracker.Interface
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

	if err := r.reconcileIstioAuthzPolicies(ctx, b, sub, p); err != nil {
		logging.FromContext(ctx).Error("Problem reconciling Istio AuthorizationPolicy", zap.Error(err))
		b.Status.MarkBindingSubjectResolvingFaiulre("IstioAuthorizationPolicyFailure", "%v", err)
		return err
	}

	b.Status.MarkBindingAvailable()
	return nil
}

func (r *Reconciler) reconcileIstioAuthzPolicies(
	ctx context.Context,
	b *v1alpha2.HTTPPolicyBinding,
	sub *tracker.Reference,
	policy *v1alpha2.HTTPPolicy,
) pkgreconciler.Event {

	// No need: https://istio.io/docs/reference/config/security/authorization-policy/#AuthorizationPolicy
	// rejectPolicy := &istiov1beta1.AuthorizationPolicy{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:            fmt.Sprintf("%s-denyall", b.Name),
	// 		Namespace:       sub.Namespace,
	// 		OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(b)},
	// 	},
	// 	Spec: istiosecurityv1beta1.AuthorizationPolicy{
	// 		Selector: &istiotypev1beta1.WorkloadSelector{
	// 			MatchLabels: sub.Selector.MatchLabels,
	// 		},
	// 	},
	// }
	// if err := r.reconcileIstioAuthz(ctx, rejectPolicy); err != nil {
	// 	return err
	// }

	allowPolicy := &istiov1beta1.AuthorizationPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:            b.Name,
			Namespace:       sub.Namespace,
			OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(b)},
		},
		Spec: istiosecurityv1beta1.AuthorizationPolicy{
			Selector: &istiotypev1beta1.WorkloadSelector{
				MatchLabels: sub.Selector.MatchLabels,
			},
			Rules: istioAuthzRulesFromPolicy(policy),
		},
	}
	return r.reconcileIstioAuthz(ctx, allowPolicy)
}

func (r *Reconciler) reconcileIstioAuthz(ctx context.Context, desired *istiov1beta1.AuthorizationPolicy) pkgreconciler.Event {
	existing, err := r.istioauthzLister.AuthorizationPolicies(desired.Namespace).Get(desired.Name)
	if apierrs.IsNotFound(err) {
		existing, err = r.istioClientSet.SecurityV1beta1().AuthorizationPolicies(desired.Namespace).Create(desired)
		if err != nil {
			return fmt.Errorf("Failed to create Istio AuthorizationPolicy: %w", err)
		} else if err != nil {
			return fmt.Errorf("Failed to get Istio AuthorizationPolicy: %w", err)
		}
	}

	if !equality.Semantic.DeepDerivative(desired.Spec, existing.Spec) {
		// Don't modify the informers copy.
		cp := existing.DeepCopy()
		cp.Spec = desired.Spec
		existing, err = r.istioClientSet.SecurityV1beta1().AuthorizationPolicies(desired.Namespace).Update(cp)
		if err != nil {
			return fmt.Errorf("Failed to update Istio AuthorizationPolicy: %w", err)
		}
	}

	return nil
}

func istioAuthzRulesFromPolicy(policy *v1alpha2.HTTPPolicy) []*istiosecurityv1beta1.Rule {
	var ret []*istiosecurityv1beta1.Rule
	for _, r := range policy.Spec.Rules {
		ir := &istiosecurityv1beta1.Rule{}
		ir.From = []*istiosecurityv1beta1.Rule_From{
			{Source: &istiosecurityv1beta1.Source{RequestPrincipals: r.Auth.Principals}},
		}
		for _, cl := range r.Auth.Claims {
			ir.When = append(ir.When, &istiosecurityv1beta1.Condition{
				Key:    fmt.Sprintf("request.auth.cliams[%s]", cl.Key),
				Values: cl.Values,
			})
		}
		for _, h := range r.Headers {
			ir.When = append(ir.When, &istiosecurityv1beta1.Condition{
				Key:    fmt.Sprintf("request.headers[%s]", h.Key),
				Values: h.Values,
			})
		}
		for _, op := range r.Operations {
			ir.To = append(ir.To, &istiosecurityv1beta1.Rule_To{
				Operation: &istiosecurityv1beta1.Operation{
					Hosts:   op.Hosts,
					Methods: op.Methods,
					Paths:   op.Paths,
				},
			})
		}
		ret = append(ret, ir)
	}
	return ret
}
