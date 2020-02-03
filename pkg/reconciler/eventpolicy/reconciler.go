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

package eventpolicy

import (
	"context"
	"fmt"
	"reflect"
	"time"

	policyduck "github.com/yolocs/knative-policy-binding/pkg/apis/duck/v1alpha1"
	security "github.com/yolocs/knative-policy-binding/pkg/apis/security/v1alpha1"
	policylisters "github.com/yolocs/knative-policy-binding/pkg/client/listers/security/v1alpha1"
	"github.com/yolocs/knative-policy-binding/pkg/reconciler"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/logging"
)

const (
	// Name of the corev1.Events emitted from the reconciliation process.
	policyReconcileError     = "PolicyReconcileError"
	policyUpdateStatusFailed = "PolicyUpdateStatusFailed"
	policyReadinessChanged   = "PolicyReadinessChanged"
)

// Reconciler is the open policy controller
type Reconciler struct {
	*reconciler.Base

	openpolicyLister  policylisters.OpenPolicyLister
	eventpolicyLister policylisters.EventPolicyLister
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
	original, err := r.eventpolicyLister.EventPolicies(namespace).Get(name)
	if apierrs.IsNotFound(err) {
		// The resource may no longer exist, in which case we stop processing.
		logger.Errorf("eventpolicy %q no longer exists", key)
		return nil
	} else if err != nil {
		return err
	}
	// Don't modify the informers copy.
	policy := original.DeepCopy()

	// Reconcile this copy of the Broker and then write back any status
	// updates regardless of whether the reconcile error out.
	reconcileErr := r.reconcile(ctx, policy)
	if reconcileErr != nil {
		logging.FromContext(ctx).Warn("Error reconciling EventPolicy", zap.Error(reconcileErr))
		r.Recorder.Eventf(policy, corev1.EventTypeWarning, policyReconcileError, fmt.Sprintf("EventPolicy reconcile error: %v", reconcileErr))
	} else {
		logging.FromContext(ctx).Debug("EventPolicy reconciled")
	}

	// Since the reconciler took a crack at this, make sure it's reflected
	// in the status correctly.
	policy.Status.ObservedGeneration = original.Generation

	if _, updateStatusErr := r.updateStatus(ctx, policy); updateStatusErr != nil {
		logging.FromContext(ctx).Warn("Failed to update the EventPolicy status", zap.Error(updateStatusErr))
		r.Recorder.Eventf(policy, corev1.EventTypeWarning, policyUpdateStatusFailed, "Failed to update EventPolicy's status: %v", updateStatusErr)
		return updateStatusErr
	}

	// Requeue if the resource is not ready:
	return reconcileErr
}

func (r *Reconciler) reconcile(ctx context.Context, p *security.EventPolicy) error {
	p.Status.InitializeConditions()

	if p.DeletionTimestamp != nil {
		// Everything is cleaned up by the garbage collector.
		return nil
	}

	op, err := r.reconcileOpenPolicy(ctx, p)
	if err != nil {
		logging.FromContext(ctx).Error("Problem reconcile OpenPolicy", zap.Error(err))
		p.Status.MarkOpenPolicyFailed("OpenPolicyFailure", "%v", err)
		return err
	}

	p.Status.PropagateFromOpenPolicyStatus(op.Name, &op.Status)
	return nil
}

func (r *Reconciler) reconcileOpenPolicy(ctx context.Context, p *security.EventPolicy) (*security.OpenPolicy, error) {
	// To fill.
	rule := MakeOpenPolicyRule(p.Spec.Rules)
	desired := &security.OpenPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "ce-" + p.Name,
			Namespace:       p.Namespace,
			OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(p)},
		},
		Spec: security.OpenPolicySpec{
			PolicyableSpec: policyduck.PolicyableSpec{
				CheckPayload: p.Spec.CheckPayload,
			},
			Rule: rule,
		},
	}
	op, err := r.openpolicyLister.OpenPolicies(desired.Namespace).Get(desired.Name)
	if apierrs.IsNotFound(err) {
		op, err = r.SecurityClientSet.SecurityV1alpha1().OpenPolicies(desired.Namespace).Create(desired)
		if err != nil {
			return nil, fmt.Errorf("failed to create openpolicy: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to get openpolicy: %w", err)
	}

	if !equality.Semantic.DeepDerivative(desired.Spec, op.Spec) {
		// Don't modify the informers copy.
		cp := op.DeepCopy()
		cp.Spec = desired.Spec
		op, err = r.SecurityClientSet.SecurityV1alpha1().OpenPolicies(desired.Namespace).Update(cp)
		if err != nil {
			return nil, fmt.Errorf("failed to update openpolicy: %w", err)
		}
	}

	return op, nil
}

func (r *Reconciler) updateStatus(ctx context.Context, desired *security.EventPolicy) (*security.EventPolicy, error) {
	policy, err := r.eventpolicyLister.EventPolicies(desired.Namespace).Get(desired.Name)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(policy.Status, desired.Status) {
		return policy, nil
	}

	becomesReady := desired.Status.IsReady() && !policy.Status.IsReady()

	// Don't modify the informers copy.
	existing := policy.DeepCopy()
	existing.Status = desired.Status

	p, err := r.SecurityClientSet.SecurityV1alpha1().EventPolicies(desired.Namespace).UpdateStatus(existing)
	if err == nil && becomesReady {
		duration := time.Since(p.ObjectMeta.CreationTimestamp.Time)
		logging.FromContext(ctx).Infof("EventPolicy %q became ready after %v", policy.Name, duration)
		r.Recorder.Event(policy, corev1.EventTypeNormal, policyReadinessChanged, fmt.Sprintf("EventPolicy %q became ready", policy.Name))
	}

	return p, err
}
