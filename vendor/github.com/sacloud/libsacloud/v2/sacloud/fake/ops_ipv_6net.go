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

// List is fake implementation
func (o *IPv6NetOp) List(ctx context.Context, zone string) (*sacloud.IPv6NetListResult, error) {
	results, _ := find(o.key, zone, nil)
	var values []*sacloud.IPv6Net
	for _, res := range results {
		dest := &sacloud.IPv6Net{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &sacloud.IPv6NetListResult{
		Total:    len(results),
		Count:    len(results),
		From:     0,
		IPv6Nets: values,
	}, nil
}

// Read is fake implementation
func (o *IPv6NetOp) Read(ctx context.Context, zone string, id types.ID) (*sacloud.IPv6Net, error) {
	value := getIPv6NetByID(zone, id)
	if value == nil {
		return nil, newErrorNotFound(o.key, id)
	}
	dest := &sacloud.IPv6Net{}
	copySameNameField(value, dest)
	return dest, nil
}
