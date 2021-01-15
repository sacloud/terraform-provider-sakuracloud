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

package fake

import (
	"context"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// Find is fake implementation
func (o *LicenseOp) Find(ctx context.Context, conditions *sacloud.FindCondition) (*sacloud.LicenseFindResult, error) {
	results, _ := find(o.key, sacloud.APIDefaultZone, conditions)
	var values []*sacloud.License
	for _, res := range results {
		dest := &sacloud.License{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &sacloud.LicenseFindResult{
		Total:    len(results),
		Count:    len(results),
		From:     0,
		Licenses: values,
	}, nil
}

// Create is fake implementation
func (o *LicenseOp) Create(ctx context.Context, param *sacloud.LicenseCreateRequest) (*sacloud.License, error) {
	result := &sacloud.License{}
	copySameNameField(param, result)
	fill(result, fillID, fillCreatedAt, fillModifiedAt)
	result.LicenseInfoName = "Windows RDS SAL"
	putLicense(sacloud.APIDefaultZone, result)
	return result, nil
}

// Read is fake implementation
func (o *LicenseOp) Read(ctx context.Context, id types.ID) (*sacloud.License, error) {
	value := getLicenseByID(sacloud.APIDefaultZone, id)
	if value == nil {
		return nil, newErrorNotFound(o.key, id)
	}
	dest := &sacloud.License{}
	copySameNameField(value, dest)
	return dest, nil
}

// Update is fake implementation
func (o *LicenseOp) Update(ctx context.Context, id types.ID, param *sacloud.LicenseUpdateRequest) (*sacloud.License, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)

	putLicense(sacloud.APIDefaultZone, value)
	return value, nil
}

// Delete is fake implementation
func (o *LicenseOp) Delete(ctx context.Context, id types.ID) error {
	_, err := o.Read(ctx, id)
	if err != nil {
		return err
	}
	ds().Delete(o.key, sacloud.APIDefaultZone, id)
	return nil
}
