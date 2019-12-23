// Copyright 2016-2019 The Libsacloud Authors
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

package registry

import (
	"context"
	"fmt"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// Builder SIMのセットアップを行う
type Builder struct {
	Name        string
	Description string
	Tags        types.Tags
	IconID      types.ID
	Visibility  types.EContainerRegistryVisibility
	NamePrefix  string
	Users       []*User

	SettingsHash string
	Client       *APIClient
}

// User ユーザー
type User struct {
	UserName string
	Password string
}

// Validate 値の検証
func (b *Builder) Validate(ctx context.Context) error {
	if b.NamePrefix == "" {
		return fmt.Errorf("name prefix is required")
	}
	return nil
}

// Build コンテナレジストリの作成
func (b *Builder) Build(ctx context.Context) (*sacloud.ContainerRegistry, error) {
	if err := b.Validate(ctx); err != nil {
		return nil, err
	}

	reg, err := b.Client.ContainerRegistry.Create(ctx, &sacloud.ContainerRegistryCreateRequest{
		Name:        b.Name,
		Description: b.Description,
		Tags:        b.Tags,
		IconID:      b.IconID,
		Visibility:  b.Visibility,
		NamePrefix:  b.NamePrefix,
	})
	if err != nil {
		return nil, err
	}

	// add users
	for _, user := range b.Users {
		u := &sacloud.ContainerRegistryUserCreateRequest{
			UserName: user.UserName,
			Password: user.Password,
		}
		if err := b.Client.ContainerRegistry.AddUser(ctx, reg.ID, u); err != nil {
			return nil, err
		}
	}

	return reg, nil
}

// Update SIMの更新
func (b *Builder) Update(ctx context.Context, id types.ID) (*sacloud.ContainerRegistry, error) {
	if err := b.Validate(ctx); err != nil {
		return nil, err
	}

	// check exists
	_, err := b.Client.ContainerRegistry.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = b.Client.ContainerRegistry.Update(ctx, id, &sacloud.ContainerRegistryUpdateRequest{
		Name:         b.Name,
		Description:  b.Description,
		Tags:         b.Tags,
		IconID:       b.IconID,
		Visibility:   b.Visibility,
		SettingsHash: b.SettingsHash,
	})
	if err != nil {
		return nil, err
	}

	// reconcile user
	added, updated, deleted, err := b.collectUserUpdates(ctx, id)
	if err != nil {
		return nil, err
	}
	// added
	for _, user := range added {
		u := &sacloud.ContainerRegistryUserCreateRequest{
			UserName: user.UserName,
			Password: user.Password,
		}
		if err := b.Client.ContainerRegistry.AddUser(ctx, id, u); err != nil {
			return nil, err
		}
	}
	// updated
	for _, user := range updated {
		u := &sacloud.ContainerRegistryUserUpdateRequest{
			Password: user.Password,
		}
		err := b.Client.ContainerRegistry.UpdateUser(ctx, id, user.UserName, u)
		if err != nil {
			return nil, err
		}
	}
	// deleted
	for _, u := range deleted {
		if err := b.Client.ContainerRegistry.DeleteUser(ctx, id, u.UserName); err != nil {
			return nil, err
		}
	}

	// reload
	reg, err := b.Client.ContainerRegistry.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	return reg, nil
}

func (b *Builder) collectUserUpdates(ctx context.Context, id types.ID) (added, updated, deleted []*User, e error) {
	searched, err := b.Client.ContainerRegistry.ListUsers(ctx, id)
	if err != nil {
		e = err
		return
	}
	users := searched.Users

	// added/updated
	for _, desired := range b.Users {
		isExists := false
		for _, current := range users {
			if desired.UserName == current.UserName {
				updated = append(updated, desired)
				isExists = true
				break
			}
		}
		if !isExists {
			added = append(added, desired)
		}
	}
	// deleted
	for _, current := range users {
		isExists := false
		for _, desired := range b.Users {
			if desired.UserName == current.UserName {
				isExists = true
				break
			}
		}
		if !isExists {
			deleted = append(deleted, &User{
				UserName: current.UserName,
			})
		}
	}
	return added, updated, deleted, nil
}
