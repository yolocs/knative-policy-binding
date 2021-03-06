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

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/yolocs/knative-policy-binding/pkg/apis/security/v1alpha1"
	scheme "github.com/yolocs/knative-policy-binding/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PolicyBindingsGetter has a method to return a PolicyBindingInterface.
// A group's client should implement this interface.
type PolicyBindingsGetter interface {
	PolicyBindings(namespace string) PolicyBindingInterface
}

// PolicyBindingInterface has methods to work with PolicyBinding resources.
type PolicyBindingInterface interface {
	Create(*v1alpha1.PolicyBinding) (*v1alpha1.PolicyBinding, error)
	Update(*v1alpha1.PolicyBinding) (*v1alpha1.PolicyBinding, error)
	UpdateStatus(*v1alpha1.PolicyBinding) (*v1alpha1.PolicyBinding, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.PolicyBinding, error)
	List(opts v1.ListOptions) (*v1alpha1.PolicyBindingList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PolicyBinding, err error)
	PolicyBindingExpansion
}

// policyBindings implements PolicyBindingInterface
type policyBindings struct {
	client rest.Interface
	ns     string
}

// newPolicyBindings returns a PolicyBindings
func newPolicyBindings(c *SecurityV1alpha1Client, namespace string) *policyBindings {
	return &policyBindings{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the policyBinding, and returns the corresponding policyBinding object, and an error if there is any.
func (c *policyBindings) Get(name string, options v1.GetOptions) (result *v1alpha1.PolicyBinding, err error) {
	result = &v1alpha1.PolicyBinding{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("policybindings").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of PolicyBindings that match those selectors.
func (c *policyBindings) List(opts v1.ListOptions) (result *v1alpha1.PolicyBindingList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.PolicyBindingList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("policybindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested policyBindings.
func (c *policyBindings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("policybindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a policyBinding and creates it.  Returns the server's representation of the policyBinding, and an error, if there is any.
func (c *policyBindings) Create(policyBinding *v1alpha1.PolicyBinding) (result *v1alpha1.PolicyBinding, err error) {
	result = &v1alpha1.PolicyBinding{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("policybindings").
		Body(policyBinding).
		Do().
		Into(result)
	return
}

// Update takes the representation of a policyBinding and updates it. Returns the server's representation of the policyBinding, and an error, if there is any.
func (c *policyBindings) Update(policyBinding *v1alpha1.PolicyBinding) (result *v1alpha1.PolicyBinding, err error) {
	result = &v1alpha1.PolicyBinding{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("policybindings").
		Name(policyBinding.Name).
		Body(policyBinding).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *policyBindings) UpdateStatus(policyBinding *v1alpha1.PolicyBinding) (result *v1alpha1.PolicyBinding, err error) {
	result = &v1alpha1.PolicyBinding{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("policybindings").
		Name(policyBinding.Name).
		SubResource("status").
		Body(policyBinding).
		Do().
		Into(result)
	return
}

// Delete takes name of the policyBinding and deletes it. Returns an error if one occurs.
func (c *policyBindings) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("policybindings").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *policyBindings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("policybindings").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched policyBinding.
func (c *policyBindings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PolicyBinding, err error) {
	result = &v1alpha1.PolicyBinding{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("policybindings").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
