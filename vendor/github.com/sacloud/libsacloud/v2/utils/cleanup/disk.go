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

package cleanup

import (
	"context"
	"fmt"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/query"
)

// DeleteDisk 他のリソースから参照されていないかを確認した上で削除する
func DeleteDisk(ctx context.Context, caller sacloud.APICaller, zone string, id types.ID, option query.CheckReferencedOption) error {
	if err := query.WaitWhileDiskIsReferenced(ctx, caller, zone, id, option); err != nil {
		return fmt.Errorf("Disk[%s] is still being used by other resources: %s", id, err)
	}
	return sacloud.NewDiskOp(caller).Delete(ctx, zone, id)
}
