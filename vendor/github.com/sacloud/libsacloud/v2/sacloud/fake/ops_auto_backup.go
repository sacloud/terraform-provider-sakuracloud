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
func (o *AutoBackupOp) Find(ctx context.Context, zone string, conditions *sacloud.FindCondition) (*sacloud.AutoBackupFindResult, error) {
	results, _ := find(o.key, zone, conditions)
	var values []*sacloud.AutoBackup
	for _, res := range results {
		dest := &sacloud.AutoBackup{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &sacloud.AutoBackupFindResult{
		Total:       len(results),
		Count:       len(results),
		From:        0,
		AutoBackups: values,
	}, nil
}

// Create is fake implementation
func (o *AutoBackupOp) Create(ctx context.Context, zone string, param *sacloud.AutoBackupCreateRequest) (*sacloud.AutoBackup, error) {
	result := &sacloud.AutoBackup{}
	copySameNameField(param, result)
	fill(result, fillID, fillCreatedAt)

	result.Availability = types.Availabilities.Available
	result.SettingsHash = "settingshash"
	result.AccountID = accountID
	result.ZoneID = zoneIDs[zone]
	result.ZoneName = zone

	putAutoBackup(zone, result)
	return result, nil
}

// Read is fake implementation
func (o *AutoBackupOp) Read(ctx context.Context, zone string, id types.ID) (*sacloud.AutoBackup, error) {
	value := getAutoBackupByID(zone, id)
	if value == nil {
		return nil, newErrorNotFound(o.key, id)
	}
	dest := &sacloud.AutoBackup{}
	copySameNameField(value, dest)
	return dest, nil
}

// Update is fake implementation
func (o *AutoBackupOp) Update(ctx context.Context, zone string, id types.ID, param *sacloud.AutoBackupUpdateRequest) (*sacloud.AutoBackup, error) {
	value, err := o.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)

	putAutoBackup(zone, value)
	return value, nil
}

// UpdateSettings is fake implementation
func (o *AutoBackupOp) UpdateSettings(ctx context.Context, zone string, id types.ID, param *sacloud.AutoBackupUpdateSettingsRequest) (*sacloud.AutoBackup, error) {
	value, err := o.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)

	putAutoBackup(zone, value)
	return value, nil
}

// Delete is fake implementation
func (o *AutoBackupOp) Delete(ctx context.Context, zone string, id types.ID) error {
	_, err := o.Read(ctx, zone, id)
	if err != nil {
		return err
	}

	ds().Delete(o.key, zone, id)
	return nil
}
