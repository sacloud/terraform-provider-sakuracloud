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

package cleanup

import (
	"context"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/query"
)

// DeleteSIM SIMの無効化&削除
func DeleteSIM(ctx context.Context, client sacloud.SIMAPI, id types.ID) error {
	sim, err := query.FindSIMByID(ctx, client, id)
	if err != nil {
		return err
	}
	if sim.Info.Activated {
		if err := client.Deactivate(ctx, id); err != nil {
			return err
		}
	}
	return client.Delete(ctx, id)
}
