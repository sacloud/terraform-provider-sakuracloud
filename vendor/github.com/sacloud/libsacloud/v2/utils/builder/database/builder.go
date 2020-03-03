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

package database

import (
	"context"
	"fmt"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/builder"
	"github.com/sacloud/libsacloud/v2/utils/power"
	"github.com/sacloud/libsacloud/v2/utils/setup"
)

// Builder データベースの構築を行う
type Builder struct {
	PlanID             types.ID
	SwitchID           types.ID
	IPAddresses        []string
	NetworkMaskLen     int
	DefaultRoute       string
	Conf               *sacloud.DatabaseRemarkDBConfCommon
	SourceID           types.ID
	CommonSetting      *sacloud.DatabaseSettingCommon
	BackupSetting      *sacloud.DatabaseSettingBackup
	ReplicationSetting *sacloud.DatabaseReplicationSetting
	Name               string
	Description        string
	Tags               types.Tags
	IconID             types.ID

	SetupOptions *builder.RetryableSetupParameter
	Client       *APIClient
}

func (b *Builder) init() {
	if b.SetupOptions == nil {
		b.SetupOptions = builder.DefaultSetupOptions()
	}
}

// Validate 設定値の検証
func (b *Builder) Validate(ctx context.Context, zone string) error {
	requiredValues := map[string]bool{
		"PlanID":         b.PlanID.IsEmpty(),
		"SwitchID":       b.SwitchID.IsEmpty(),
		"IPAddresses":    len(b.IPAddresses) == 0,
		"NetworkMaskLen": b.NetworkMaskLen == 0,
		"Conf":           b.Conf == nil,
		"CommonSetting":  b.CommonSetting == nil,
	}
	for key, empty := range requiredValues {
		if empty {
			return fmt.Errorf("%s is required", key)
		}
	}
	return nil
}

// Build データベースの作成や設定をまとめて行う
func (b *Builder) Build(ctx context.Context, zone string) (*sacloud.Database, error) {
	b.init()

	if err := b.Validate(ctx, zone); err != nil {
		return nil, err
	}

	builder := &setup.RetryableSetup{
		Create: func(ctx context.Context, zone string) (accessor.ID, error) {
			return b.Client.Database.Create(ctx, zone, &sacloud.DatabaseCreateRequest{
				PlanID:             b.PlanID,
				SwitchID:           b.SwitchID,
				IPAddresses:        b.IPAddresses,
				NetworkMaskLen:     b.NetworkMaskLen,
				DefaultRoute:       b.DefaultRoute,
				Conf:               b.Conf,
				SourceID:           b.SourceID,
				CommonSetting:      b.CommonSetting,
				BackupSetting:      b.BackupSetting,
				ReplicationSetting: b.ReplicationSetting,
				Name:               b.Name,
				Description:        b.Description,
				Tags:               b.Tags,
				IconID:             b.IconID,
			})
		},
		Delete: func(ctx context.Context, zone string, id types.ID) error {
			return b.Client.Database.Delete(ctx, zone, id)
		},
		Read: func(ctx context.Context, zone string, id types.ID) (interface{}, error) {
			return b.Client.Database.Read(ctx, zone, id)
		},
		ProvisionBeforeUp: func(ctx context.Context, zone string, id types.ID, _ interface{}) error {
			return b.Client.Database.Config(ctx, zone, id)
		},
		IsWaitForCopy:       true,
		IsWaitForUp:         true,
		RetryCount:          b.SetupOptions.RetryCount,
		DeleteRetryCount:    b.SetupOptions.DeleteRetryCount,
		DeleteRetryInterval: b.SetupOptions.DeleteRetryInterval,
		PollingInterval:     b.SetupOptions.PollingInterval,
	}

	result, err := builder.Setup(ctx, zone)
	var db *sacloud.Database
	if result != nil {
		db = result.(*sacloud.Database)
	}
	if err != nil {
		return db, err
	}

	// refresh
	db, err = b.Client.Database.Read(ctx, zone, db.ID)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Update データベースの更新
func (b *Builder) Update(ctx context.Context, zone string, id types.ID) (*sacloud.Database, error) {
	b.init()

	if err := b.Validate(ctx, zone); err != nil {
		return nil, err
	}

	// check Database is exists
	db, err := b.Client.Database.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}

	isNeedShutdown, err := b.collectUpdateInfo(db)
	if err != nil {
		return nil, err
	}

	isNeedRestart := false
	if db.InstanceStatus.IsUp() && isNeedShutdown {
		isNeedRestart = true
		if err := power.ShutdownDatabase(ctx, b.Client.Database, zone, id, false); err != nil {
			return nil, err
		}
	}

	_, err = b.Client.Database.Update(ctx, zone, id, &sacloud.DatabaseUpdateRequest{
		Name:               b.Name,
		Description:        b.Description,
		Tags:               b.Tags,
		IconID:             b.IconID,
		CommonSetting:      b.CommonSetting,
		BackupSetting:      b.BackupSetting,
		ReplicationSetting: b.ReplicationSetting,
		SettingsHash:       db.SettingsHash,
	})
	if err != nil {
		return nil, err
	}
	if err := b.Client.Database.Config(ctx, zone, id); err != nil {
		return nil, err
	}
	if isNeedRestart {
		if err := power.BootDatabase(ctx, b.Client.Database, zone, id); err != nil {
			return nil, err
		}
	}

	// refresh
	db, err = b.Client.Database.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}
	return db, err
}

func (b *Builder) collectUpdateInfo(db *sacloud.Database) (isNeedShutdown bool, err error) {
	isNeedShutdown = b.CommonSetting.ReplicaPassword != db.CommonSetting.ReplicaPassword
	return
}
