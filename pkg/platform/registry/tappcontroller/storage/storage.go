/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

package storage

import (
	"context"

	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/platform"
	apiserverutil "tkestack.io/tke/pkg/apiserver/util"
	"tkestack.io/tke/pkg/platform/registry/tappcontroller"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
)

// Storage includes storage for tapp controller and all sub resources.
type Storage struct {
	TappController *REST
	Status         *StatusREST
}

// NewStorage returns a Storage object that will work against tapp controller.
func NewStorage(optsGetter genericregistry.RESTOptionsGetter, privilegedUsername string) *Storage {
	strategy := tappcontroller.NewStrategy()
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &platform.TappController{} },
		NewListFunc:              func() runtime.Object { return &platform.TappControllerList{} },
		DefaultQualifiedResource: platform.Resource("tappcontrollers"),
		PredicateFunc:            tappcontroller.MatchTappController,

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
	}
	store.TableConvertor = rest.NewDefaultTableConvertor(store.DefaultQualifiedResource)
	options := &genericregistry.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    tappcontroller.GetAttrs,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create tapp controller etcd rest storage", log.Err(err))
	}

	statusStore := *store
	statusStore.UpdateStrategy = tappcontroller.NewStatusStrategy(strategy)

	return &Storage{
		TappController: &REST{store, privilegedUsername},
		Status:         &StatusREST{&statusStore},
	}
}

// ValidateGetObjectAndTenantID validate name and tenantID, if success return TappController
func ValidateGetObjectAndTenantID(ctx context.Context, store *registry.Store, name string, options *metav1.GetOptions) (runtime.Object, error) {
	obj, err := store.Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	tc := obj.(*platform.TappController)
	if err := util.FilterTappController(ctx, tc); err != nil {
		return nil, err
	}
	return tc, nil
}

// REST implements a RESTStorage for tapp controller against etcd.
type REST struct {
	*registry.Store
	privilegedUsername string
}

var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of
// short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"tc"}
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	wrappedOptions := apiserverutil.PredicateListOptions(ctx, options)
	return r.Store.List(ctx, wrappedOptions)
}

// Get finds a resource in the storage by name and returns it.
func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return ValidateGetObjectAndTenantID(ctx, r.Store, name, options)
}

// Update finds a resource in the storage and updates it.
func (r *REST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// sub resources should never allow create on update.
	_, err := ValidateGetObjectAndTenantID(ctx, r.Store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	return r.Store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}

// Delete enforces life-cycle rules for cluster termination
func (r *REST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	_, err := ValidateGetObjectAndTenantID(ctx, r.Store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	return r.Store.Delete(ctx, name, deleteValidation, options)
}

// Watch selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) Watch(ctx context.Context, options *metainternal.ListOptions) (watch.Interface, error) {
	wrappedOptions := apiserverutil.PredicateListOptions(ctx, options)
	return r.Store.Watch(ctx, wrappedOptions)
}

// DeleteCollection selects all resources in the storage matching given 'listOptions'
// and deletes them.
func (r *REST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternal.ListOptions) (runtime.Object, error) {
	wrappedOptions := apiserverutil.PredicateListOptions(ctx, listOptions)
	return r.Store.DeleteCollection(ctx, deleteValidation, options, wrappedOptions)
}

// StatusREST implements the REST endpoint for changing the status of a tapp controller.
type StatusREST struct {
	store *registry.Store
}

// StatusREST implements Patcher.
var _ = rest.Patcher(&StatusREST{})

// New returns an empty object that can be used with Create and Update after
// request data has been put into it.
func (r *StatusREST) New() runtime.Object {
	return r.store.New()
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *StatusREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return ValidateGetObjectAndTenantID(ctx, r.store, name, options)
}

// Update alters the status subset of an object.
func (r *StatusREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// sub resources should never allow create on update.
	_, err := ValidateGetObjectAndTenantID(ctx, r.store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	return r.store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}
