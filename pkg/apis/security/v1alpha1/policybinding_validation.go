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

	"knative.dev/pkg/apis"
)

// Validate implements apis.Validatable
func (pb *PolicyBinding) Validate(ctx context.Context) *apis.FieldError {
	errs := pb.Spec.Validate(ctx).ViaField("spec")
	if pb.Spec.Subject.Namespace != "" && pb.Namespace != pb.Spec.Subject.Namespace {
		errs = errs.Also(apis.ErrInvalidValue(pb.Spec.Subject.Namespace, "spec.subject.namespace"))
	}
	if pb.Spec.Policy.Namespace != "" && pb.Namespace != pb.Spec.Policy.Namespace {
		errs = errs.Also(apis.ErrInvalidValue(pb.Spec.Policy.Namespace, "spec.policy.namespace"))
	}
	return errs
}

// Validate implements apis.Validatable
func (pbs *PolicyBindingSpec) Validate(ctx context.Context) *apis.FieldError {
	return pbs.Subject.Validate(ctx).ViaField("subject")
}
