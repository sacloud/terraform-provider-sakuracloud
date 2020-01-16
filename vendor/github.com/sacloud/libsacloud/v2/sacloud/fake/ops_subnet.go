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
func (o *SubnetOp) Find(ctx context.Context, zone string, conditions *sacloud.FindCondition) (*sacloud.SubnetFindResult, error) {
	results, _ := find(o.key, zone, conditions)
	var values []*sacloud.Subnet
	for _, res := range results {
		dest := &sacloud.Subnet{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &sacloud.SubnetFindResult{
		Total:   len(results),
		Count:   len(results),
		From:    0,
		Subnets: values,
	}, nil
}

// Read is fake implementation
func (o *SubnetOp) Read(ctx context.Context, zone string, id types.ID) (*sacloud.Subnet, error) {
	value := getSubnetByID(zone, id)
	if value == nil {
		return nil, newErrorNotFound(o.key, id)
	}
	dest := &sacloud.Subnet{}
	copySameNameField(value, dest)
	return dest, nil
}
