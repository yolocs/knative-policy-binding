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

package authbinding

import (
	"context"
	"fmt"
	"reflect"
	"time"

	security "github.com/yolocs/knative-policy-binding/pkg/apis/security/v1alpha1"
	policylisters "github.com/yolocs/knative-policy-binding/pkg/client/listers/security/v1alpha1"
	"github.com/yolocs/knative-policy-binding/pkg/reconciler"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	duckv1alpha1 "knative.dev/pkg/apis/duck/v1alpha1"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/tracker"
)

const (
	// Name of the corev1.Events emitted from the reconciliation process.
	authorizableBindingReconcileError     = "AuthorizableBindingReconcileError"
	authorizableBindingUpdateStatusFailed = "AuthorizableBindingUpdateStatusFailed"
	authorizableBindingReadinessChanged   = "AuthorizableBindingReadinessChanged"
)

// Reconciler reconciles the authorizable bindings.
type Reconciler struct {
	*reconciler.Base

	authbindingLister   policylisters.AuthorizableBindingLister
	policybindingLister policylisters.PolicyBindingLister

	podspecableResolver *PodspecableResolver
}

// Check that our Reconciler implements controller.Reconciler
var _ controller.Reconciler = (*Reconciler)(nil)

// Reconcile compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Broker resource
// with the current status of the resource.
func (r *Reconciler) Reconcile(ctx context.Context, key string) error {
	logger := logging.FromContext(ctx)

	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		logger.Errorf("invalid resource key: %s", key)
		return nil
	}

	// Get the resource with this namespace/name.
	original, err := r.authbindingLister.AuthorizableBindings(namespace).Get(name)
	if apierrs.IsNotFound(err) {
		// The resource may no longer exist, in which case we stop processing.
		logger.Errorf("authorizablebinding %q no longer exists", key)
		return nil
	} else if err != nil {
		return err
	}
	// Don't modify the informers copy.
	binding := original.DeepCopy()

	// Reconcile this copy of the Broker and then write back any status
	// updates regardless of whether the reconcile error out.
	reconcileErr := r.reconcile(ctx, binding)
	if reconcileErr != nil {
		logging.FromContext(ctx).Warn("Error reconciling EventPolicy", zap.Error(reconcileErr))
		r.Recorder.Eventf(binding, corev1.EventTypeWarning, authorizableBindingReconcileError, fmt.Sprintf("AuthorizableBinding reconcile error: %v", reconcileErr))
	} else {
		logging.FromContext(ctx).Debug("EventPolicy reconciled")
	}

	// Since the reconciler took a crack at this, make sure it's reflected
	// in the status correctly.
	binding.Status.ObservedGeneration = original.Generation

	if _, updateStatusErr := r.updateStatus(ctx, binding); updateStatusErr != nil {
		logging.FromContext(ctx).Warn("Failed to update the EventPolicy status", zap.Error(updateStatusErr))
		r.Recorder.Eventf(binding, corev1.EventTypeWarning, authorizableBindingUpdateStatusFailed, "Failed to update AuthorizableBinding's status: %v", updateStatusErr)
		return updateStatusErr
	}

	// Requeue if the resource is not ready:
	return reconcileErr
}

func (r *Reconciler) reconcile(ctx context.Context, binding *security.AuthorizableBinding) error {
	binding.Status.InitializeConditions()

	if binding.DeletionTimestamp != nil {
		// Everything is cleaned up by the garbage collector.
		return nil
	}

	sub, err := r.podspecableResolver.PodspecableFromRef(binding.Spec.Subject, binding)
	if err != nil {
		logging.FromContext(ctx).Error("Problem resolving authorizable subject", zap.Error(err))
		binding.Status.MarkBindingSubjectResolvingFaiulre("SubjectResolvingFailure", "%v", err)
		return err
	}
	binding.Status.MarkBindingSubjectResolved(sub)

	pb, err := r.reconcilePolicyBinding(ctx, binding, sub)
	if err != nil {
		logging.FromContext(ctx).Error("Problem reconciling policybinding", zap.Error(err))
		binding.Status.MarkBindingPolicyFaiulre("PolicyBindingFailure", "%v", err)
		return err
	}

	if pb.Status.IsReady() {
		binding.Status.MarkBindingPolicyReady()
		binding.Status.MarkBindingAvailable()
	} else {
		c := pb.Status.GetTopLevelCondition()
		if c != nil {
			binding.Status.MarkBindingPolicyFaiulre(c.Reason, c.Message)
		}
	}

	return nil
}

func (r *Reconciler) reconcilePolicyBinding(ctx context.Context, binding *security.AuthorizableBinding, sub *tracker.Reference) (*security.PolicyBinding, error) {
	desired := &security.PolicyBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:            binding.Name,
			Namespace:       binding.Namespace,
			OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(binding)},
		},
		Spec: security.PolicyBindingSpec{
			BindingSpec: duckv1alpha1.BindingSpec{
				Subject: *sub,
			},
			Policy: binding.Spec.Policy,
		},
	}
	b, err := r.policybindingLister.PolicyBindings(desired.Namespace).Get(desired.Name)
	if apierrs.IsNotFound(err) {
		b, err = r.SecurityClientSet.SecurityV1alpha1().PolicyBindings(desired.Namespace).Create(desired)
		if err != nil {
			return nil, fmt.Errorf("failed to create policybinding: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to get policybinding: %w", err)
	}

	if !equality.Semantic.DeepDerivative(desired.Spec, b.Spec) {
		// Don't modify the informers copy.
		cp := b.DeepCopy()
		cp.Spec = desired.Spec
		b, err = r.SecurityClientSet.SecurityV1alpha1().PolicyBindings(desired.Namespace).Update(cp)
		if err != nil {
			return nil, fmt.Errorf("failed to update policybinding: %w", err)
		}
	}

	return b, nil
}

func (r *Reconciler) updateStatus(ctx context.Context, desired *security.AuthorizableBinding) (*security.AuthorizableBinding, error) {
	binding, err := r.authbindingLister.AuthorizableBindings(desired.Namespace).Get(desired.Name)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(binding.Status, desired.Status) {
		return binding, nil
	}

	becomesReady := desired.Status.IsReady() && !binding.Status.IsReady()

	// Don't modify the informers copy.
	existing := binding.DeepCopy()
	existing.Status = desired.Status

	b, err := r.SecurityClientSet.SecurityV1alpha1().AuthorizableBindings(desired.Namespace).UpdateStatus(existing)
	if err == nil && becomesReady {
		duration := time.Since(b.ObjectMeta.CreationTimestamp.Time)
		logging.FromContext(ctx).Infof("EventPolicy %q became ready after %v", binding.Name, duration)
		r.Recorder.Event(binding, corev1.EventTypeNormal, authorizableBindingReadinessChanged, fmt.Sprintf("AuthorizableBinding %q became ready", binding.Name))
	}

	return b, err
}
