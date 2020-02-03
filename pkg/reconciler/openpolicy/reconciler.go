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
	"context"
	"fmt"
	"reflect"
	"time"

	security "github.com/yolocs/knative-policy-binding/pkg/apis/security/v1alpha1"
	openpolicylisters "github.com/yolocs/knative-policy-binding/pkg/client/listers/security/v1alpha1"
	"github.com/yolocs/knative-policy-binding/pkg/opa"
	"github.com/yolocs/knative-policy-binding/pkg/reconciler"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
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

	openpolicyLister openpolicylisters.OpenPolicyLister
	configmapLister  corev1listers.ConfigMapLister

	agentImage string
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
	original, err := r.openpolicyLister.OpenPolicies(namespace).Get(name)
	if apierrs.IsNotFound(err) {
		// The resource may no longer exist, in which case we stop processing.
		logger.Errorf("openpolicy %q no longer exists", key)
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
		logging.FromContext(ctx).Warn("Error reconciling OpenPolicy", zap.Error(reconcileErr))
		r.Recorder.Eventf(policy, corev1.EventTypeWarning, policyReconcileError, fmt.Sprintf("OpenPolicy reconcile error: %v", reconcileErr))
	} else {
		logging.FromContext(ctx).Debug("OpenPolicy reconciled")
	}

	// Since the reconciler took a crack at this, make sure it's reflected
	// in the status correctly.
	policy.Status.ObservedGeneration = original.Generation

	if _, updateStatusErr := r.updateStatus(ctx, policy); updateStatusErr != nil {
		logging.FromContext(ctx).Warn("Failed to update the OpenPolicy status", zap.Error(updateStatusErr))
		r.Recorder.Eventf(policy, corev1.EventTypeWarning, policyUpdateStatusFailed, "Failed to update OpenPolicy's status: %v", updateStatusErr)
		return updateStatusErr
	}

	// Requeue if the resource is not ready:
	return reconcileErr
}

func (r *Reconciler) reconcile(ctx context.Context, p *security.OpenPolicy) error {
	p.Status.InitializeConditions()

	if p.DeletionTimestamp != nil {
		// Everything is cleaned up by the garbage collector.
		return nil
	}

	if err := r.reconcileConfigMap(ctx, p); err != nil {
		logging.FromContext(ctx).Error("Problem reconcile ConfigMap", zap.Error(err))
		p.Status.MarkConfigMapFailed("ConfigMapFailure", "%v", err)
		return err
	}

	p.Status.MarkConfigMapReady(p.Name)
	p.Status.SetDeciderURI(MakeDeciderURL())
	p.Status.SetAgentSpec(MakeAgentSpec(r.agentImage, p.Name, p.ObjectMeta.Labels["agenttest"]))
	p.Status.MarkReady()
	return nil
}

func (r *Reconciler) reconcileConfigMap(ctx context.Context, p *security.OpenPolicy) error {
	m := opa.GenerateFromTemplate(p.Spec.Rule)
	desired := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            p.Name,
			Namespace:       p.Namespace,
			OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(p)},
		},
		Data: map[string]string{
			policyFileName: m,
		},
	}
	cm, err := r.configmapLister.ConfigMaps(p.Namespace).Get(p.Name)
	if apierrs.IsNotFound(err) {
		cm, err = r.KubeClientSet.CoreV1().ConfigMaps(p.Namespace).Create(desired)
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

func (r *Reconciler) updateStatus(ctx context.Context, desired *security.OpenPolicy) (*security.OpenPolicy, error) {
	policy, err := r.openpolicyLister.OpenPolicies(desired.Namespace).Get(desired.Name)
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

	p, err := r.SecurityClientSet.SecurityV1alpha1().OpenPolicies(desired.Namespace).UpdateStatus(existing)
	if err == nil && becomesReady {
		duration := time.Since(p.ObjectMeta.CreationTimestamp.Time)
		logging.FromContext(ctx).Infof("OpenPolicy %q became ready after %v", policy.Name, duration)
		r.Recorder.Event(policy, corev1.EventTypeNormal, policyReadinessChanged, fmt.Sprintf("OpenPolicy %q became ready", policy.Name))
	}

	return p, err
}
