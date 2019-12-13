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

package sim

import (
	"context"
	"fmt"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/query"
)

// Builder SIMのセットアップを行う
type Builder struct {
	Name        string
	Description string
	Tags        types.Tags
	IconID      types.ID
	ICCID       string
	PassCode    string

	Activate bool
	IMEI     string
	Carrier  []*sacloud.SIMNetworkOperatorConfig

	Client *APIClient
}

// Validate 値の検証
func (b *Builder) Validate(ctx context.Context) error {
	if b.ICCID == "" {
		return fmt.Errorf("iccid is required")
	}
	if b.PassCode == "" {
		return fmt.Errorf("iccid is required")
	}
	if len(b.Carrier) == 0 {
		return fmt.Errorf("carrier is required")
	}
	return nil
}

// Build SIMの作成
func (b *Builder) Build(ctx context.Context) (*sacloud.SIM, error) {
	if err := b.Validate(ctx); err != nil {
		return nil, err
	}

	sim, err := b.Client.SIM.Create(ctx, &sacloud.SIMCreateRequest{
		Name:        b.Name,
		Description: b.Description,
		Tags:        b.Tags,
		IconID:      b.IconID,
		ICCID:       b.ICCID,
		PassCode:    b.PassCode,
	})
	if err != nil {
		return nil, err
	}

	if err := b.Client.SIM.SetNetworkOperator(ctx, sim.ID, b.Carrier); err != nil {
		return nil, err
	}

	if b.Activate {
		if err := b.Client.SIM.Activate(ctx, sim.ID); err != nil {
			return nil, err
		}
	}

	if b.IMEI != "" {
		if err := b.Client.SIM.IMEILock(ctx, sim.ID, &sacloud.SIMIMEILockRequest{IMEI: b.IMEI}); err != nil {
			return nil, err
		}
	}

	// reload
	sim, err = query.FindSIMByID(ctx, b.Client.SIM, sim.ID)
	if err != nil {
		return nil, err
	}
	return sim, nil
}

// Update SIMの更新
func (b *Builder) Update(ctx context.Context, id types.ID) (*sacloud.SIM, error) {
	if err := b.Validate(ctx); err != nil {
		return nil, err
	}

	sim, err := query.FindSIMByID(ctx, b.Client.SIM, id)
	if err != nil {
		return nil, err
	}

	_, err = b.Client.SIM.Update(ctx, id, &sacloud.SIMUpdateRequest{
		Name:        b.Name,
		Description: b.Description,
		Tags:        b.Tags,
		IconID:      b.IconID,
	})
	if err != nil {
		return nil, err
	}

	if err := b.Client.SIM.SetNetworkOperator(ctx, sim.ID, b.Carrier); err != nil {
		return nil, err
	}

	if !b.Activate && sim.Info.Activated {
		if err := b.Client.SIM.Deactivate(ctx, sim.ID); err != nil {
			return nil, err
		}
	}
	if b.Activate && !sim.Info.Activated {
		if err := b.Client.SIM.Activate(ctx, sim.ID); err != nil {
			return nil, err
		}
	}

	if b.IMEI == "" || sim.Info.IMEILock {
		if err := b.Client.SIM.IMEIUnlock(ctx, sim.ID); err != nil {
			return nil, err
		}
	}
	if b.IMEI != "" {
		if err := b.Client.SIM.IMEILock(ctx, sim.ID, &sacloud.SIMIMEILockRequest{IMEI: b.IMEI}); err != nil {
			return nil, err
		}
	}

	// reload
	sim, err = query.FindSIMByID(ctx, b.Client.SIM, id)
	if err != nil {
		return nil, err
	}
	return sim, nil
}
