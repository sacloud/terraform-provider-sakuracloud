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
func (o *GSLBOp) Find(ctx context.Context, conditions *sacloud.FindCondition) (*sacloud.GSLBFindResult, error) {
	results, _ := find(o.key, sacloud.APIDefaultZone, conditions)
	var values []*sacloud.GSLB
	for _, res := range results {
		dest := &sacloud.GSLB{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &sacloud.GSLBFindResult{
		Total: len(results),
		Count: len(results),
		From:  0,
		GSLBs: values,
	}, nil
}

// Create is fake implementation
func (o *GSLBOp) Create(ctx context.Context, param *sacloud.GSLBCreateRequest) (*sacloud.GSLB, error) {
	result := &sacloud.GSLB{}
	copySameNameField(param, result)
	fill(result, fillID, fillCreatedAt, fillAvailability)

	result.FQDN = fmt.Sprintf("site-%d.gslb7.example.ne.jp", result.ID)
	result.SettingsHash = "settingshash"

	putGSLB(sacloud.APIDefaultZone, result)
	return result, nil
}

// Read is fake implementation
func (o *GSLBOp) Read(ctx context.Context, id types.ID) (*sacloud.GSLB, error) {
	value := getGSLBByID(sacloud.APIDefaultZone, id)
	if value == nil {
		return nil, newErrorNotFound(o.key, id)
	}

	dest := &sacloud.GSLB{}
	copySameNameField(value, dest)
	return dest, nil
}

// Update is fake implementation
func (o *GSLBOp) Update(ctx context.Context, id types.ID, param *sacloud.GSLBUpdateRequest) (*sacloud.GSLB, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	if param.DelayLoop == 0 {
		param.DelayLoop = 10 // default value
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)

	putGSLB(sacloud.APIDefaultZone, value)
	return value, nil
}

// Patch is fake implementation
func (o *GSLBOp) Patch(ctx context.Context, id types.ID, param *sacloud.GSLBPatchRequest) (*sacloud.GSLB, error) {
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
	if param.PatchEmptyToHealthCheck {
		value.HealthCheck = nil
	}
	if param.PatchEmptyToDelayLoop {
		value.DelayLoop = 0
	}
	if param.PatchEmptyToWeighted {
		value.Weighted = types.StringFlag(false)
	}
	if param.PatchEmptyToSorryServer {
		value.SorryServer = ""
	}
	if param.PatchEmptyToDestinationServers {
		value.DestinationServers = nil
	}

	putGSLB(sacloud.APIDefaultZone, value)
	return value, nil
}

// UpdateSettings is fake implementation
func (o *GSLBOp) UpdateSettings(ctx context.Context, id types.ID, param *sacloud.GSLBUpdateSettingsRequest) (*sacloud.GSLB, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)

	putGSLB(sacloud.APIDefaultZone, value)
	return value, nil
}

// PatchSettings is fake implementation
func (o *GSLBOp) PatchSettings(ctx context.Context, id types.ID, param *sacloud.GSLBPatchSettingsRequest) (*sacloud.GSLB, error) {
	patchParam := &sacloud.GSLBPatchRequest{}
	copySameNameField(param, patchParam)
	return o.Patch(ctx, id, patchParam)
}

// Delete is fake implementation
func (o *GSLBOp) Delete(ctx context.Context, id types.ID) error {
	_, err := o.Read(ctx, id)
	if err != nil {
		return err
	}
	ds().Delete(o.key, sacloud.APIDefaultZone, id)
	return nil
}
