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

package query

import (
	"context"
	"fmt"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/ostype"
)

// FindArchiveByOSType OS種別ごとの最新安定板のアーカイブを取得
func FindArchiveByOSType(ctx context.Context, api ArchiveFinder, zone string, os ostype.ArchiveOSType) (*sacloud.Archive, error) {
	filter, ok := ostype.ArchiveCriteria[os]
	if !ok {
		return nil, fmt.Errorf("unsupported ostype.ArchiveOSType: %v", os)
	}

	searched, err := api.Find(ctx, zone, &sacloud.FindCondition{Filter: filter})
	if err != nil {
		return nil, err
	}
	if searched.Count == 0 {
		return nil, fmt.Errorf("archive not found with ostype.ArchiveOSType: %v", os)
	}
	return searched.Archives[0], nil
}
