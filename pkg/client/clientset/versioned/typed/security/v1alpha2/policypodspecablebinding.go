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

package v1alpha2

import (
	"time"

	v1alpha2 "github.com/yolocs/knative-policy-binding/pkg/apis/security/v1alpha2"
	scheme "github.com/yolocs/knative-policy-binding/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PolicyPodspecableBindingsGetter has a method to return a PolicyPodspecableBindingInterface.
// A group's client should implement this interface.
type PolicyPodspecableBindingsGetter interface {
	PolicyPodspecableBindings(namespace string) PolicyPodspecableBindingInterface
}

// PolicyPodspecableBindingInterface has methods to work with PolicyPodspecableBinding resources.
type PolicyPodspecableBindingInterface interface {
	Create(*v1alpha2.PolicyPodspecableBinding) (*v1alpha2.PolicyPodspecableBinding, error)
	Update(*v1alpha2.PolicyPodspecableBinding) (*v1alpha2.PolicyPodspecableBinding, error)
	UpdateStatus(*v1alpha2.PolicyPodspecableBinding) (*v1alpha2.PolicyPodspecableBinding, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha2.PolicyPodspecableBinding, error)
	List(opts v1.ListOptions) (*v1alpha2.PolicyPodspecableBindingList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha2.PolicyPodspecableBinding, err error)
	PolicyPodspecableBindingExpansion
}

// policyPodspecableBindings implements PolicyPodspecableBindingInterface
type policyPodspecableBindings struct {
	client rest.Interface
	ns     string
}

// newPolicyPodspecableBindings returns a PolicyPodspecableBindings
func newPolicyPodspecableBindings(c *SecurityV1alpha2Client, namespace string) *policyPodspecableBindings {
	return &policyPodspecableBindings{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the policyPodspecableBinding, and returns the corresponding policyPodspecableBinding object, and an error if there is any.
func (c *policyPodspecableBindings) Get(name string, options v1.GetOptions) (result *v1alpha2.PolicyPodspecableBinding, err error) {
	result = &v1alpha2.PolicyPodspecableBinding{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("policypodspecablebindings").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of PolicyPodspecableBindings that match those selectors.
func (c *policyPodspecableBindings) List(opts v1.ListOptions) (result *v1alpha2.PolicyPodspecableBindingList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha2.PolicyPodspecableBindingList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("policypodspecablebindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested policyPodspecableBindings.
func (c *policyPodspecableBindings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("policypodspecablebindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a policyPodspecableBinding and creates it.  Returns the server's representation of the policyPodspecableBinding, and an error, if there is any.
func (c *policyPodspecableBindings) Create(policyPodspecableBinding *v1alpha2.PolicyPodspecableBinding) (result *v1alpha2.PolicyPodspecableBinding, err error) {
	result = &v1alpha2.PolicyPodspecableBinding{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("policypodspecablebindings").
		Body(policyPodspecableBinding).
		Do().
		Into(result)
	return
}

// Update takes the representation of a policyPodspecableBinding and updates it. Returns the server's representation of the policyPodspecableBinding, and an error, if there is any.
func (c *policyPodspecableBindings) Update(policyPodspecableBinding *v1alpha2.PolicyPodspecableBinding) (result *v1alpha2.PolicyPodspecableBinding, err error) {
	result = &v1alpha2.PolicyPodspecableBinding{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("policypodspecablebindings").
		Name(policyPodspecableBinding.Name).
		Body(policyPodspecableBinding).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *policyPodspecableBindings) UpdateStatus(policyPodspecableBinding *v1alpha2.PolicyPodspecableBinding) (result *v1alpha2.PolicyPodspecableBinding, err error) {
	result = &v1alpha2.PolicyPodspecableBinding{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("policypodspecablebindings").
		Name(policyPodspecableBinding.Name).
		SubResource("status").
		Body(policyPodspecableBinding).
		Do().
		Into(result)
	return
}

// Delete takes name of the policyPodspecableBinding and deletes it. Returns an error if one occurs.
func (c *policyPodspecableBindings) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("policypodspecablebindings").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *policyPodspecableBindings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("policypodspecablebindings").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched policyPodspecableBinding.
func (c *policyPodspecableBindings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha2.PolicyPodspecableBinding, err error) {
	result = &v1alpha2.PolicyPodspecableBinding{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("policypodspecablebindings").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
