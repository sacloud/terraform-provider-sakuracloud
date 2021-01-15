// Copyright 2016-2021 The Libsacloud Authors
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

package internet

import (
	"context"
	"errors"
	"fmt"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// DefaultNotFoundRetry スイッチ+ルータ作成後のReadで404が返ってこなくなるまでに許容する404エラーの回数
var DefaultNotFoundRetry = 360 // デフォルトの5秒おきリトライの場合30分

// Builder スイッチ+ルータの構築を行う
type Builder struct {
	Name           string
	Description    string
	Tags           types.Tags
	IconID         types.ID
	NetworkMaskLen int
	BandWidthMbps  int
	EnableIPv6     bool

	NotFoundRetry int

	NoWait bool

	Client *APIClient
}

// Validate 設定値の検証
func (b *Builder) Validate(ctx context.Context, zone string) error {
	requiredValues := map[string]bool{
		"NetworkMaskLen": b.NetworkMaskLen == 0,
		"BandWidthMbps":  b.BandWidthMbps == 0,
	}
	for key, empty := range requiredValues {
		if empty {
			return fmt.Errorf("%s is required", key)
		}
	}

	if b.NoWait && b.EnableIPv6 {
		return errors.New("NoWait=true is not supported when EnableIPv6=true")
	}
	return nil
}

// Build ルータ+スイッチの作成や設定をまとめて行う
func (b *Builder) Build(ctx context.Context, zone string) (*sacloud.Internet, error) {
	if b.NotFoundRetry == 0 {
		b.NotFoundRetry = DefaultNotFoundRetry
	}

	if err := b.Validate(ctx, zone); err != nil {
		return nil, err
	}

	internet, err := b.Client.Internet.Create(ctx, zone, &sacloud.InternetCreateRequest{
		Name:           b.Name,
		Description:    b.Description,
		Tags:           b.Tags,
		IconID:         b.IconID,
		NetworkMaskLen: b.NetworkMaskLen,
		BandWidthMbps:  b.BandWidthMbps,
	})
	if err != nil {
		return nil, err
	}

	if b.NoWait {
		return internet, nil
	}

	// [HACK] ルータ作成直後は GET /internet/:id が404を返すことへの対応
	waiter := sacloud.WaiterForApplianceUp(func() (interface{}, error) {
		return b.Client.Internet.Read(ctx, zone, internet.ID)
	}, b.NotFoundRetry)
	if _, err := waiter.WaitForState(ctx); err != nil {
		return internet, err
	}

	if b.EnableIPv6 {
		_, err = b.Client.Internet.EnableIPv6(ctx, zone, internet.ID)
		if err != nil {
			return internet, err
		}
	}

	return b.Client.Internet.Read(ctx, zone, internet.ID)
}

// Update スイッチ+ルータの更新
func (b *Builder) Update(ctx context.Context, zone string, id types.ID) (*sacloud.Internet, error) {
	if err := b.Validate(ctx, zone); err != nil {
		return nil, err
	}

	// check Internet is exists
	internet, err := b.Client.Internet.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}

	if b.NetworkMaskLen != internet.NetworkMaskLen {
		return nil, fmt.Errorf("unsupported operation: NetworkMaskLen is changed: current: %d new: %d", internet.NetworkMaskLen, b.NetworkMaskLen)
	}

	internet, err = b.Client.Internet.Update(ctx, zone, internet.ID, &sacloud.InternetUpdateRequest{
		Name:        b.Name,
		Description: b.Description,
		Tags:        b.Tags,
		IconID:      b.IconID,
	})
	if err != nil {
		return nil, err
	}

	if internet.BandWidthMbps != b.BandWidthMbps {
		// 成功するとIDが変更となる
		internet, err = b.Client.Internet.UpdateBandWidth(ctx, zone, internet.ID, &sacloud.InternetUpdateBandWidthRequest{
			BandWidthMbps: b.BandWidthMbps,
		})
		if err != nil {
			return nil, err
		}
	}

	currentIPv6Enabled := len(internet.Switch.IPv6Nets) > 0
	if b.EnableIPv6 != currentIPv6Enabled {
		if currentIPv6Enabled {
			if err := b.Client.Internet.DisableIPv6(ctx, zone, internet.ID, internet.Switch.IPv6Nets[0].ID); err != nil {
				return nil, err
			}
		} else {
			if _, err := b.Client.Internet.EnableIPv6(ctx, zone, internet.ID); err != nil {
				return nil, err
			}
		}
	}

	return b.Client.Internet.Read(ctx, zone, internet.ID)
}
