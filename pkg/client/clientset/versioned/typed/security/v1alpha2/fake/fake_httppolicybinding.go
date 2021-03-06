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
	v1alpha2 "github.com/yolocs/knative-policy-binding/pkg/apis/security/v1alpha2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeHTTPPolicyBindings implements HTTPPolicyBindingInterface
type FakeHTTPPolicyBindings struct {
	Fake *FakeSecurityV1alpha2
	ns   string
}

var httppolicybindingsResource = schema.GroupVersionResource{Group: "security.knative.dev", Version: "v1alpha2", Resource: "httppolicybindings"}

var httppolicybindingsKind = schema.GroupVersionKind{Group: "security.knative.dev", Version: "v1alpha2", Kind: "HTTPPolicyBinding"}

// Get takes name of the hTTPPolicyBinding, and returns the corresponding hTTPPolicyBinding object, and an error if there is any.
func (c *FakeHTTPPolicyBindings) Get(name string, options v1.GetOptions) (result *v1alpha2.HTTPPolicyBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(httppolicybindingsResource, c.ns, name), &v1alpha2.HTTPPolicyBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.HTTPPolicyBinding), err
}

// List takes label and field selectors, and returns the list of HTTPPolicyBindings that match those selectors.
func (c *FakeHTTPPolicyBindings) List(opts v1.ListOptions) (result *v1alpha2.HTTPPolicyBindingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(httppolicybindingsResource, httppolicybindingsKind, c.ns, opts), &v1alpha2.HTTPPolicyBindingList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha2.HTTPPolicyBindingList{ListMeta: obj.(*v1alpha2.HTTPPolicyBindingList).ListMeta}
	for _, item := range obj.(*v1alpha2.HTTPPolicyBindingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested hTTPPolicyBindings.
func (c *FakeHTTPPolicyBindings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(httppolicybindingsResource, c.ns, opts))

}

// Create takes the representation of a hTTPPolicyBinding and creates it.  Returns the server's representation of the hTTPPolicyBinding, and an error, if there is any.
func (c *FakeHTTPPolicyBindings) Create(hTTPPolicyBinding *v1alpha2.HTTPPolicyBinding) (result *v1alpha2.HTTPPolicyBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(httppolicybindingsResource, c.ns, hTTPPolicyBinding), &v1alpha2.HTTPPolicyBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.HTTPPolicyBinding), err
}

// Update takes the representation of a hTTPPolicyBinding and updates it. Returns the server's representation of the hTTPPolicyBinding, and an error, if there is any.
func (c *FakeHTTPPolicyBindings) Update(hTTPPolicyBinding *v1alpha2.HTTPPolicyBinding) (result *v1alpha2.HTTPPolicyBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(httppolicybindingsResource, c.ns, hTTPPolicyBinding), &v1alpha2.HTTPPolicyBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.HTTPPolicyBinding), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeHTTPPolicyBindings) UpdateStatus(hTTPPolicyBinding *v1alpha2.HTTPPolicyBinding) (*v1alpha2.HTTPPolicyBinding, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(httppolicybindingsResource, "status", c.ns, hTTPPolicyBinding), &v1alpha2.HTTPPolicyBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.HTTPPolicyBinding), err
}

// Delete takes name of the hTTPPolicyBinding and deletes it. Returns an error if one occurs.
func (c *FakeHTTPPolicyBindings) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(httppolicybindingsResource, c.ns, name), &v1alpha2.HTTPPolicyBinding{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeHTTPPolicyBindings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(httppolicybindingsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha2.HTTPPolicyBindingList{})
	return err
}

// Patch applies the patch and returns the patched hTTPPolicyBinding.
func (c *FakeHTTPPolicyBindings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha2.HTTPPolicyBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(httppolicybindingsResource, c.ns, name, pt, data, subresources...), &v1alpha2.HTTPPolicyBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.HTTPPolicyBinding), err
}
