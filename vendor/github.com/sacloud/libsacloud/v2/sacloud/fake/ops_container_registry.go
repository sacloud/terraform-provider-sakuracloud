// Copyright 2016-2020 The Libsacloud Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fake

import (
	"context"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// Find is fake implementation
func (o *ContainerRegistryOp) Find(ctx context.Context, conditions *sacloud.FindCondition) (*sacloud.ContainerRegistryFindResult, error) {
	results, _ := find(o.key, sacloud.APIDefaultZone, conditions)
	var values []*sacloud.ContainerRegistry
	for _, res := range results {
		dest := &sacloud.ContainerRegistry{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &sacloud.ContainerRegistryFindResult{
		Total:               len(results),
		Count:               len(results),
		From:                0,
		ContainerRegistries: values,
	}, nil
}

// Create is fake implementation
func (o *ContainerRegistryOp) Create(ctx context.Context, param *sacloud.ContainerRegistryCreateRequest) (*sacloud.ContainerRegistry, error) {
	result := &sacloud.ContainerRegistry{}
	copySameNameField(param, result)
	fill(result, fillID, fillCreatedAt)

	result.FQDN = result.SubDomainLabel + ".sakuracr.jp"
	result.Availability = types.Availabilities.Available
	putContainerRegistry(sacloud.APIDefaultZone, result)
	return result, nil
}

// Read is fake implementation
func (o *ContainerRegistryOp) Read(ctx context.Context, id types.ID) (*sacloud.ContainerRegistry, error) {
	value := getContainerRegistryByID(sacloud.APIDefaultZone, id)
	if value == nil {
		return nil, newErrorNotFound(o.key, id)
	}
	dest := &sacloud.ContainerRegistry{}
	copySameNameField(value, dest)
	return dest, nil
}

// Update is fake implementation
func (o *ContainerRegistryOp) Update(ctx context.Context, id types.ID, param *sacloud.ContainerRegistryUpdateRequest) (*sacloud.ContainerRegistry, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)
	putContainerRegistry(sacloud.APIDefaultZone, value)
	return value, nil
}

// UpdateSettings is fake implementation
func (o *ContainerRegistryOp) UpdateSettings(ctx context.Context, id types.ID, param *sacloud.ContainerRegistryUpdateSettingsRequest) (*sacloud.ContainerRegistry, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)
	putContainerRegistry(sacloud.APIDefaultZone, value)
	return value, nil
}

// Delete is fake implementation
func (o *ContainerRegistryOp) Delete(ctx context.Context, id types.ID) error {
	_, err := o.Read(ctx, id)
	if err != nil {
		return err
	}
	ds().Delete(o.key, sacloud.APIDefaultZone, id)
	return nil
}

// ListUsers is fake implementation
func (o *ContainerRegistryOp) ListUsers(ctx context.Context, id types.ID) (*sacloud.ContainerRegistryUsers, error) {
	_, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	v := ds().Get(ResourceContainerRegistry+"Users", sacloud.APIDefaultZone, id)
	if v != nil {
		users := v.([]*sacloud.ContainerRegistryUser)
		return &sacloud.ContainerRegistryUsers{
			Users: users,
		}, nil
	}

	return nil, err
}

// AddUser is fake implementation
func (o *ContainerRegistryOp) AddUser(ctx context.Context, id types.ID, param *sacloud.ContainerRegistryUserCreateRequest) error {
	_, err := o.Read(ctx, id)
	if err != nil {
		return err
	}

	var users []*sacloud.ContainerRegistryUser
	v := ds().Get(ResourceContainerRegistry+"Users", sacloud.APIDefaultZone, id)
	if v != nil {
		users = v.([]*sacloud.ContainerRegistryUser)
	}
	users = append(users, &sacloud.ContainerRegistryUser{
		UserName:   param.UserName,
		Permission: param.Permission,
	})

	ds().Put(ResourceContainerRegistry+"Users", sacloud.APIDefaultZone, id, users)
	return nil
}

// UpdateUser is fake implementation
func (o *ContainerRegistryOp) UpdateUser(ctx context.Context, id types.ID, username string, param *sacloud.ContainerRegistryUserUpdateRequest) error {
	_, err := o.Read(ctx, id)
	if err != nil {
		return err
	}

	v := ds().Get(ResourceContainerRegistry+"Users", sacloud.APIDefaultZone, id)
	if v == nil {
		return newErrorNotFound(ResourceContainerRegistry+"Users", id)
	}
	users := v.([]*sacloud.ContainerRegistryUser)
	for _, u := range users {
		if u.UserName == username {
			u.Permission = param.Permission
		}
	}
	ds().Put(ResourceContainerRegistry+"Users", sacloud.APIDefaultZone, id, users)
	return nil
}

// DeleteUser is fake implementation
func (o *ContainerRegistryOp) DeleteUser(ctx context.Context, id types.ID, username string) error {
	_, err := o.Read(ctx, id)
	if err != nil {
		return err
	}

	v := ds().Get(ResourceContainerRegistry+"Users", sacloud.APIDefaultZone, id)
	if v == nil {
		return newErrorNotFound(ResourceContainerRegistry+"Users", id)
	}
	users := v.([]*sacloud.ContainerRegistryUser)
	var newUsers []*sacloud.ContainerRegistryUser
	for _, u := range users {
		if u.UserName != username {
			newUsers = append(newUsers, u)
		}
	}

	ds().Put(ResourceContainerRegistry+"Users", sacloud.APIDefaultZone, id, newUsers)
	return nil
}
