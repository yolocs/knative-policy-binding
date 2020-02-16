// +build !ignore_autogenerated

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha2

import (
	v1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	tracker "knative.dev/pkg/tracker"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPPolicy) DeepCopyInto(out *HTTPPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPPolicy.
func (in *HTTPPolicy) DeepCopy() *HTTPPolicy {
	if in == nil {
		return nil
	}
	out := new(HTTPPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HTTPPolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPPolicyBinding) DeepCopyInto(out *HTTPPolicyBinding) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPPolicyBinding.
func (in *HTTPPolicyBinding) DeepCopy() *HTTPPolicyBinding {
	if in == nil {
		return nil
	}
	out := new(HTTPPolicyBinding)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HTTPPolicyBinding) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPPolicyBindingList) DeepCopyInto(out *HTTPPolicyBindingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HTTPPolicyBinding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPPolicyBindingList.
func (in *HTTPPolicyBindingList) DeepCopy() *HTTPPolicyBindingList {
	if in == nil {
		return nil
	}
	out := new(HTTPPolicyBindingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HTTPPolicyBindingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPPolicyBindingSpec) DeepCopyInto(out *HTTPPolicyBindingSpec) {
	*out = *in
	if in.Subject != nil {
		in, out := &in.Subject, &out.Subject
		*out = new(v1.ObjectReference)
		**out = **in
	}
	if in.Policy != nil {
		in, out := &in.Policy, &out.Policy
		*out = new(v1.ObjectReference)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPPolicyBindingSpec.
func (in *HTTPPolicyBindingSpec) DeepCopy() *HTTPPolicyBindingSpec {
	if in == nil {
		return nil
	}
	out := new(HTTPPolicyBindingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPPolicyBindingStatus) DeepCopyInto(out *HTTPPolicyBindingStatus) {
	*out = *in
	in.Status.DeepCopyInto(&out.Status)
	if in.ResolvedSubject != nil {
		in, out := &in.ResolvedSubject, &out.ResolvedSubject
		*out = new(tracker.Reference)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPPolicyBindingStatus.
func (in *HTTPPolicyBindingStatus) DeepCopy() *HTTPPolicyBindingStatus {
	if in == nil {
		return nil
	}
	out := new(HTTPPolicyBindingStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPPolicyList) DeepCopyInto(out *HTTPPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HTTPPolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPPolicyList.
func (in *HTTPPolicyList) DeepCopy() *HTTPPolicyList {
	if in == nil {
		return nil
	}
	out := new(HTTPPolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HTTPPolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPPolicySpec) DeepCopyInto(out *HTTPPolicySpec) {
	*out = *in
	in.JWT.DeepCopyInto(&out.JWT)
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]RuleSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPPolicySpec.
func (in *HTTPPolicySpec) DeepCopy() *HTTPPolicySpec {
	if in == nil {
		return nil
	}
	out := new(HTTPPolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JWTSpec) DeepCopyInto(out *JWTSpec) {
	*out = *in
	if in.TriggerRules != nil {
		in, out := &in.TriggerRules, &out.TriggerRules
		*out = make([]TriggerRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JWTSpec.
func (in *JWTSpec) DeepCopy() *JWTSpec {
	if in == nil {
		return nil
	}
	out := new(JWTSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeyValueMatch) DeepCopyInto(out *KeyValueMatch) {
	*out = *in
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeyValueMatch.
func (in *KeyValueMatch) DeepCopy() *KeyValueMatch {
	if in == nil {
		return nil
	}
	out := new(KeyValueMatch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Operation) DeepCopyInto(out *Operation) {
	*out = *in
	if in.Hosts != nil {
		in, out := &in.Hosts, &out.Hosts
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Paths != nil {
		in, out := &in.Paths, &out.Paths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Methods != nil {
		in, out := &in.Methods, &out.Methods
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Operation.
func (in *Operation) DeepCopy() *Operation {
	if in == nil {
		return nil
	}
	out := new(Operation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicyAgentSpec) DeepCopyInto(out *PolicyAgentSpec) {
	*out = *in
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Container.DeepCopyInto(&out.Container)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyAgentSpec.
func (in *PolicyAgentSpec) DeepCopy() *PolicyAgentSpec {
	if in == nil {
		return nil
	}
	out := new(PolicyAgentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicyPodspecableBinding) DeepCopyInto(out *PolicyPodspecableBinding) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyPodspecableBinding.
func (in *PolicyPodspecableBinding) DeepCopy() *PolicyPodspecableBinding {
	if in == nil {
		return nil
	}
	out := new(PolicyPodspecableBinding)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PolicyPodspecableBinding) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicyPodspecableBindingList) DeepCopyInto(out *PolicyPodspecableBindingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PolicyPodspecableBinding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyPodspecableBindingList.
func (in *PolicyPodspecableBindingList) DeepCopy() *PolicyPodspecableBindingList {
	if in == nil {
		return nil
	}
	out := new(PolicyPodspecableBindingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PolicyPodspecableBindingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicyPodspecableBindingSpec) DeepCopyInto(out *PolicyPodspecableBindingSpec) {
	*out = *in
	in.BindingSpec.DeepCopyInto(&out.BindingSpec)
	if in.AgentSpec != nil {
		in, out := &in.AgentSpec, &out.AgentSpec
		*out = new(PolicyAgentSpec)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyPodspecableBindingSpec.
func (in *PolicyPodspecableBindingSpec) DeepCopy() *PolicyPodspecableBindingSpec {
	if in == nil {
		return nil
	}
	out := new(PolicyPodspecableBindingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicyPodspecableBindingStatus) DeepCopyInto(out *PolicyPodspecableBindingStatus) {
	*out = *in
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyPodspecableBindingStatus.
func (in *PolicyPodspecableBindingStatus) DeepCopy() *PolicyPodspecableBindingStatus {
	if in == nil {
		return nil
	}
	out := new(PolicyPodspecableBindingStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RequestAuth) DeepCopyInto(out *RequestAuth) {
	*out = *in
	if in.Principals != nil {
		in, out := &in.Principals, &out.Principals
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Claims != nil {
		in, out := &in.Claims, &out.Claims
		*out = make([]KeyValueMatch, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RequestAuth.
func (in *RequestAuth) DeepCopy() *RequestAuth {
	if in == nil {
		return nil
	}
	out := new(RequestAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RuleSpec) DeepCopyInto(out *RuleSpec) {
	*out = *in
	in.Auth.DeepCopyInto(&out.Auth)
	if in.Headers != nil {
		in, out := &in.Headers, &out.Headers
		*out = make([]KeyValueMatch, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Operations != nil {
		in, out := &in.Operations, &out.Operations
		*out = make([]Operation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RuleSpec.
func (in *RuleSpec) DeepCopy() *RuleSpec {
	if in == nil {
		return nil
	}
	out := new(RuleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TriggerRule) DeepCopyInto(out *TriggerRule) {
	*out = *in
	if in.ExcludePaths != nil {
		in, out := &in.ExcludePaths, &out.ExcludePaths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.IncludePaths != nil {
		in, out := &in.IncludePaths, &out.IncludePaths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TriggerRule.
func (in *TriggerRule) DeepCopy() *TriggerRule {
	if in == nil {
		return nil
	}
	out := new(TriggerRule)
	in.DeepCopyInto(out)
	return out
}