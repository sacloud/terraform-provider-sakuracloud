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

package query

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// FindSIMByID SIM+詳細情報をIDから検索
func FindSIMByID(ctx context.Context, client sacloud.SIMAPI, id types.ID) (*sacloud.SIM, error) {
	var sim *sacloud.SIM
	searched, err := client.Find(ctx, &sacloud.FindCondition{
		Include: []string{"*", "Status.sim"},
	})
	if err != nil {
		return nil, fmt.Errorf("could not find SakuraCloud SIM[%s]: %s", id, err)
	}
	for _, s := range searched.SIMs {
		if s.ID == id {
			sim = s
			break
		}
	}
	if sim == nil {
		return nil, sacloud.NewAPIError(http.MethodGet, nil, "", http.StatusNotFound, nil)
	}
	return sim, nil
}
