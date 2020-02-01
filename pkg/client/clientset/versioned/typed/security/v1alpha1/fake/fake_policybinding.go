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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/yolocs/knative-policy-binding/pkg/apis/security/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakePolicyBindings implements PolicyBindingInterface
type FakePolicyBindings struct {
	Fake *FakeSecurityV1alpha1
	ns   string
}

var policybindingsResource = schema.GroupVersionResource{Group: "security.knative.dev", Version: "v1alpha1", Resource: "policybindings"}

var policybindingsKind = schema.GroupVersionKind{Group: "security.knative.dev", Version: "v1alpha1", Kind: "PolicyBinding"}

// Get takes name of the policyBinding, and returns the corresponding policyBinding object, and an error if there is any.
func (c *FakePolicyBindings) Get(name string, options v1.GetOptions) (result *v1alpha1.PolicyBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(policybindingsResource, c.ns, name), &v1alpha1.PolicyBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PolicyBinding), err
}

// List takes label and field selectors, and returns the list of PolicyBindings that match those selectors.
func (c *FakePolicyBindings) List(opts v1.ListOptions) (result *v1alpha1.PolicyBindingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(policybindingsResource, policybindingsKind, c.ns, opts), &v1alpha1.PolicyBindingList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.PolicyBindingList{ListMeta: obj.(*v1alpha1.PolicyBindingList).ListMeta}
	for _, item := range obj.(*v1alpha1.PolicyBindingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested policyBindings.
func (c *FakePolicyBindings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(policybindingsResource, c.ns, opts))

}

// Create takes the representation of a policyBinding and creates it.  Returns the server's representation of the policyBinding, and an error, if there is any.
func (c *FakePolicyBindings) Create(policyBinding *v1alpha1.PolicyBinding) (result *v1alpha1.PolicyBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(policybindingsResource, c.ns, policyBinding), &v1alpha1.PolicyBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PolicyBinding), err
}

// Update takes the representation of a policyBinding and updates it. Returns the server's representation of the policyBinding, and an error, if there is any.
func (c *FakePolicyBindings) Update(policyBinding *v1alpha1.PolicyBinding) (result *v1alpha1.PolicyBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(policybindingsResource, c.ns, policyBinding), &v1alpha1.PolicyBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PolicyBinding), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakePolicyBindings) UpdateStatus(policyBinding *v1alpha1.PolicyBinding) (*v1alpha1.PolicyBinding, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(policybindingsResource, "status", c.ns, policyBinding), &v1alpha1.PolicyBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PolicyBinding), err
}

// Delete takes name of the policyBinding and deletes it. Returns an error if one occurs.
func (c *FakePolicyBindings) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(policybindingsResource, c.ns, name), &v1alpha1.PolicyBinding{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePolicyBindings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(policybindingsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.PolicyBindingList{})
	return err
}

// Patch applies the patch and returns the patched policyBinding.
func (c *FakePolicyBindings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PolicyBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(policybindingsResource, c.ns, name, pt, data, subresources...), &v1alpha1.PolicyBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PolicyBinding), err
}