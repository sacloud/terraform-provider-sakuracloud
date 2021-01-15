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
func (o *PrivateHostOp) Find(ctx context.Context, zone string, conditions *sacloud.FindCondition) (*sacloud.PrivateHostFindResult, error) {
	results, _ := find(o.key, zone, conditions)
	var values []*sacloud.PrivateHost
	for _, res := range results {
		dest := &sacloud.PrivateHost{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &sacloud.PrivateHostFindResult{
		Total:        len(results),
		Count:        len(results),
		From:         0,
		PrivateHosts: values,
	}, nil
}

// Create is fake implementation
func (o *PrivateHostOp) Create(ctx context.Context, zone string, param *sacloud.PrivateHostCreateRequest) (*sacloud.PrivateHost, error) {
	planOp := NewPrivateHostPlanOp()
	plan, err := planOp.Read(ctx, zone, param.PlanID)
	if err != nil {
		return nil, err
	}

	result := &sacloud.PrivateHost{}
	copySameNameField(param, result)
	fill(result, fillID, fillCreatedAt)

	result.PlanName = plan.Name
	result.PlanClass = plan.Class
	result.CPU = plan.CPU
	result.MemoryMB = plan.MemoryMB
	result.HostName = "sac-zone-svNNN"
	putPrivateHost(zone, result)
	return result, nil
}

// Read is fake implementation
func (o *PrivateHostOp) Read(ctx context.Context, zone string, id types.ID) (*sacloud.PrivateHost, error) {
	value := getPrivateHostByID(zone, id)
	if value == nil {
		return nil, newErrorNotFound(o.key, id)
	}
	dest := &sacloud.PrivateHost{}
	copySameNameField(value, dest)
	return dest, nil
}

// Update is fake implementation
func (o *PrivateHostOp) Update(ctx context.Context, zone string, id types.ID, param *sacloud.PrivateHostUpdateRequest) (*sacloud.PrivateHost, error) {
	value, err := o.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)

	putPrivateHost(zone, value)
	return value, nil
}

// Delete is fake implementation
func (o *PrivateHostOp) Delete(ctx context.Context, zone string, id types.ID) error {
	_, err := o.Read(ctx, zone, id)
	if err != nil {
		return err
	}

	ds().Delete(o.key, zone, id)
	return nil
}
