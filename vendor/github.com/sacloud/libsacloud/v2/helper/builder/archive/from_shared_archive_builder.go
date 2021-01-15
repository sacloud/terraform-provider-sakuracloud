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

package archive

import (
	"context"
	"fmt"

	"github.com/sacloud/libsacloud/v2/helper/query"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// FromSharedArchiveBuilder 共有アーカイブからアーカイブの作成を行う
type FromSharedArchiveBuilder struct {
	Name            string
	Description     string
	Tags            types.Tags
	IconID          types.ID
	SourceSharedKey types.ArchiveShareKey

	NoWait bool
	Client *APIClient
}

// Validate 設定値の検証
func (b *FromSharedArchiveBuilder) Validate(ctx context.Context, zone string) error {
	requiredValues := map[string]bool{
		"Name":            b.Name == "",
		"SourceSharedKey": b.SourceSharedKey == "",
	}
	for key, empty := range requiredValues {
		if empty {
			return fmt.Errorf("%s is required", key)
		}
	}
	if !b.SourceSharedKey.ValidFormat() {
		return fmt.Errorf("archive shared key is invalid format: key:%q", b.SourceSharedKey)
	}
	return nil
}

// Build 共有アーカイブからアーカイブの作成を行う
func (b *FromSharedArchiveBuilder) Build(ctx context.Context, zone string) (*sacloud.Archive, error) {
	if err := b.Validate(ctx, zone); err != nil {
		return nil, err
	}

	zoneID, err := query.ZoneIDFromName(ctx, b.Client.Zone, zone)
	if err != nil {
		return nil, err
	}

	archive, err := b.Client.Archive.CreateFromShared(ctx, b.SourceSharedKey.Zone(), b.SourceSharedKey.SourceArchiveID(), zoneID,
		&sacloud.ArchiveCreateRequestFromShared{
			Name:            b.Name,
			Description:     b.Description,
			Tags:            b.Tags,
			IconID:          b.IconID,
			SourceSharedKey: b.SourceSharedKey,
		})
	if err != nil {
		return nil, err
	}

	if b.NoWait {
		return archive, nil
	}

	lastState, err := sacloud.WaiterForReady(func() (interface{}, error) {
		return b.Client.Archive.Read(ctx, zone, archive.ID)
	}).WaitForState(ctx)

	var ret *sacloud.Archive
	if lastState != nil {
		ret = lastState.(*sacloud.Archive)
	}
	return ret, err
}
