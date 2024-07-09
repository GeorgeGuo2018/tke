/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1 "tkestack.io/tke/api/application/v1"
	scheme "tkestack.io/tke/api/client/clientset/versioned/scheme"
)

// UpgradePoliciesGetter has a method to return a UpgradePolicyInterface.
// A group's client should implement this interface.
type UpgradePoliciesGetter interface {
	UpgradePolicies() UpgradePolicyInterface
}

// UpgradePolicyInterface has methods to work with UpgradePolicy resources.
type UpgradePolicyInterface interface {
	Create(ctx context.Context, upgradePolicy *v1.UpgradePolicy, opts metav1.CreateOptions) (*v1.UpgradePolicy, error)
	Update(ctx context.Context, upgradePolicy *v1.UpgradePolicy, opts metav1.UpdateOptions) (*v1.UpgradePolicy, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.UpgradePolicy, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.UpgradePolicyList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.UpgradePolicy, err error)
	UpgradePolicyExpansion
}

// upgradePolicies implements UpgradePolicyInterface
type upgradePolicies struct {
	client rest.Interface
}

// newUpgradePolicies returns a UpgradePolicies
func newUpgradePolicies(c *ApplicationV1Client) *upgradePolicies {
	return &upgradePolicies{
		client: c.RESTClient(),
	}
}

// Get takes name of the upgradePolicy, and returns the corresponding upgradePolicy object, and an error if there is any.
func (c *upgradePolicies) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.UpgradePolicy, err error) {
	result = &v1.UpgradePolicy{}
	err = c.client.Get().
		Resource("upgradepolicies").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of UpgradePolicies that match those selectors.
func (c *upgradePolicies) List(ctx context.Context, opts metav1.ListOptions) (result *v1.UpgradePolicyList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.UpgradePolicyList{}
	err = c.client.Get().
		Resource("upgradepolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested upgradePolicies.
func (c *upgradePolicies) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("upgradepolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a upgradePolicy and creates it.  Returns the server's representation of the upgradePolicy, and an error, if there is any.
func (c *upgradePolicies) Create(ctx context.Context, upgradePolicy *v1.UpgradePolicy, opts metav1.CreateOptions) (result *v1.UpgradePolicy, err error) {
	result = &v1.UpgradePolicy{}
	err = c.client.Post().
		Resource("upgradepolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(upgradePolicy).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a upgradePolicy and updates it. Returns the server's representation of the upgradePolicy, and an error, if there is any.
func (c *upgradePolicies) Update(ctx context.Context, upgradePolicy *v1.UpgradePolicy, opts metav1.UpdateOptions) (result *v1.UpgradePolicy, err error) {
	result = &v1.UpgradePolicy{}
	err = c.client.Put().
		Resource("upgradepolicies").
		Name(upgradePolicy.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(upgradePolicy).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the upgradePolicy and deletes it. Returns an error if one occurs.
func (c *upgradePolicies) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("upgradepolicies").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *upgradePolicies) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("upgradepolicies").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched upgradePolicy.
func (c *upgradePolicies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.UpgradePolicy, err error) {
	result = &v1.UpgradePolicy{}
	err = c.client.Patch(pt).
		Resource("upgradepolicies").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}