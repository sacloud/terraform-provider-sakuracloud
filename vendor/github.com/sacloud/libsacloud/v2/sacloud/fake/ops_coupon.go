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
func (o *CouponOp) Find(ctx context.Context, accountID types.ID) (*sacloud.CouponFindResult, error) {
	results, _ := find(o.key, sacloud.APIDefaultZone, nil)
	var values []*sacloud.Coupon
	for _, res := range results {
		dest := &sacloud.Coupon{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &sacloud.CouponFindResult{
		Total:   len(results),
		Count:   len(results),
		From:    0,
		Coupons: values,
	}, nil
}
