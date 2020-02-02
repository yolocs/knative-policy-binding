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

	policyduck "github.com/yolocs/knative-policy-binding/pkg/apis/duck/v1alpha1"
)

type policyKey struct{}

// WithPolicy notes the context for the policy duck type.
func WithPolicy(ctx context.Context, policy *policyduck.Policyable) context.Context {
	return context.WithValue(ctx, policyKey{}, policy)
}

// GetPolicy accesses the policy in the context.
func GetPolicy(ctx context.Context) *policyduck.Policyable {
	value := ctx.Value(policyKey{})
	if value == nil {
		return nil
	}
	return value.(*policyduck.Policyable)
}
