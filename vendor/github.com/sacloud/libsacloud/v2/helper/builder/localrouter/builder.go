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

package localrouter

import (
	"context"
	"fmt"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// Builder ローカルルータの構築を行う
type Builder struct {
	Name        string
	Description string
	Tags        types.Tags
	IconID      types.ID

	Switch       *sacloud.LocalRouterSwitch
	Interface    *sacloud.LocalRouterInterface
	Peers        []*sacloud.LocalRouterPeer
	StaticRoutes []*sacloud.LocalRouterStaticRoute

	Client *APIClient
}

// Validate 設定値の検証
func (b *Builder) Validate(ctx context.Context) error {
	requiredValues := map[string]bool{
		"Switch":    b.Switch == nil,
		"Interface": b.Interface == nil,
	}
	for key, empty := range requiredValues {
		if empty {
			return fmt.Errorf("%s is required", key)
		}
	}
	return nil
}

// Build ローカルルータの作成や設定をまとめて行う
func (b *Builder) Build(ctx context.Context) (*sacloud.LocalRouter, error) {
	if err := b.Validate(ctx); err != nil {
		return nil, err
	}

	localRouter, err := b.Client.LocalRouter.Create(ctx, &sacloud.LocalRouterCreateRequest{
		Name:        b.Name,
		Description: b.Description,
		Tags:        b.Tags,
		IconID:      b.IconID,
	})
	if err != nil {
		return nil, err
	}

	lr, err := b.Client.LocalRouter.UpdateSettings(ctx, localRouter.ID, &sacloud.LocalRouterUpdateSettingsRequest{
		Switch:       b.Switch,
		Interface:    b.Interface,
		StaticRoutes: b.StaticRoutes,
		SettingsHash: localRouter.SettingsHash,
	})
	if err != nil {
		return localRouter, err
	}
	localRouter = lr

	if len(b.Peers) > 0 {
		lr, err := b.Client.LocalRouter.UpdateSettings(ctx, localRouter.ID, &sacloud.LocalRouterUpdateSettingsRequest{
			Switch:       localRouter.Switch,
			Interface:    localRouter.Interface,
			StaticRoutes: localRouter.StaticRoutes,
			Peers:        b.Peers,
			SettingsHash: localRouter.SettingsHash,
		})
		if err != nil {
			return localRouter, err
		}
		localRouter = lr
	}

	return localRouter, nil
}

// Update ローカルルータの更新
func (b *Builder) Update(ctx context.Context, id types.ID) (*sacloud.LocalRouter, error) {
	if err := b.Validate(ctx); err != nil {
		return nil, err
	}

	// check Internet is exists
	localRouter, err := b.Client.LocalRouter.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	localRouter, err = b.Client.LocalRouter.Update(ctx, id, &sacloud.LocalRouterUpdateRequest{
		Switch:       b.Switch,
		Interface:    b.Interface,
		Peers:        b.Peers,
		StaticRoutes: b.StaticRoutes,
		SettingsHash: localRouter.SettingsHash,
		Name:         b.Name,
		Description:  b.Description,
		Tags:         b.Tags,
		IconID:       b.IconID,
	})
	if err != nil {
		return nil, err
	}

	return localRouter, nil
}
