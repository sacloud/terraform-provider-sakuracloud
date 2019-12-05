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

package fake

import (
	"context"
	"fmt"

	"github.com/imdario/mergo"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// Find is fake implementation
func (o *DNSOp) Find(ctx context.Context, conditions *sacloud.FindCondition) (*sacloud.DNSFindResult, error) {
	results, _ := find(o.key, sacloud.APIDefaultZone, conditions)
	var values []*sacloud.DNS
	for _, res := range results {
		dest := &sacloud.DNS{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &sacloud.DNSFindResult{
		Total: len(results),
		Count: len(results),
		From:  0,
		DNS:   values,
	}, nil
}

// Create is fake implementation
func (o *DNSOp) Create(ctx context.Context, param *sacloud.DNSCreateRequest) (*sacloud.DNS, error) {
	result := &sacloud.DNS{}
	copySameNameField(param, result)
	fill(result, fillID, fillCreatedAt)

	result.Availability = types.Availabilities.Available
	result.SettingsHash = "settingshash"
	result.DNSZone = param.Name

	putDNS(sacloud.APIDefaultZone, result)
	return result, nil
}

// Read is fake implementation
func (o *DNSOp) Read(ctx context.Context, id types.ID) (*sacloud.DNS, error) {
	value := getDNSByID(sacloud.APIDefaultZone, id)
	if value == nil {
		return nil, newErrorNotFound(o.key, id)
	}
	dest := &sacloud.DNS{}
	copySameNameField(value, dest)
	return dest, nil
}

// Update is fake implementation
func (o *DNSOp) Update(ctx context.Context, id types.ID, param *sacloud.DNSUpdateRequest) (*sacloud.DNS, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)

	putDNS(sacloud.APIDefaultZone, value)
	return value, nil
}

// Patch is fake implementation
func (o *DNSOp) Patch(ctx context.Context, id types.ID, param *sacloud.DNSPatchRequest) (*sacloud.DNS, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	patchParam := make(map[string]interface{})
	if err := mergo.Map(&patchParam, value); err != nil {
		return nil, fmt.Errorf("patch is failed: %s", err)
	}
	if err := mergo.Map(&patchParam, param); err != nil {
		return nil, fmt.Errorf("patch is failed: %s", err)
	}
	if err := mergo.Map(param, &patchParam); err != nil {
		return nil, fmt.Errorf("patch is failed: %s", err)
	}
	copySameNameField(param, value)

	if param.PatchEmptyToDescription {
		value.Description = ""
	}
	if param.PatchEmptyToTags {
		value.Tags = nil
	}
	if param.PatchEmptyToIconID {
		value.IconID = types.ID(int64(0))
	}
	if param.PatchEmptyToRecords {
		value.Records = nil
	}

	putDNS(sacloud.APIDefaultZone, value)
	return value, nil
}

// UpdateSettings is fake implementation
func (o *DNSOp) UpdateSettings(ctx context.Context, id types.ID, param *sacloud.DNSUpdateSettingsRequest) (*sacloud.DNS, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)

	putDNS(sacloud.APIDefaultZone, value)
	return value, nil
}

// PatchSettings is fake implementation
func (o *DNSOp) PatchSettings(ctx context.Context, id types.ID, param *sacloud.DNSPatchSettingsRequest) (*sacloud.DNS, error) {
	patchParam := &sacloud.DNSPatchRequest{}
	copySameNameField(param, patchParam)
	return o.Patch(ctx, id, patchParam)
}

// Delete is fake implementation
func (o *DNSOp) Delete(ctx context.Context, id types.ID) error {
	_, err := o.Read(ctx, id)
	if err != nil {
		return err
	}

	ds().Delete(o.key, sacloud.APIDefaultZone, id)
	return nil
}
