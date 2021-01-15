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

package cleanup

import (
	"context"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// DeleteInternet スイッチ+ルータの削除 IPv6の無効化やサブネットの削除を一括して行う
func DeleteInternet(ctx context.Context, client sacloud.InternetAPI, zone string, id types.ID) error {
	internet, err := client.Read(ctx, zone, id)
	if err != nil {
		return err
	}

	// Disable IPv6
	if len(internet.Switch.IPv6Nets) > 0 {
		if err := client.DisableIPv6(ctx, zone, id, internet.Switch.IPv6Nets[0].ID); err != nil {
			return err
		}
	}

	// Delete Subnets
	for _, subnet := range internet.Switch.Subnets {
		if subnet.NextHop != "" {
			if err := client.DeleteSubnet(ctx, zone, internet.ID, subnet.ID); err != nil {
				return err
			}
		}
	}
	return client.Delete(ctx, zone, id)
}
