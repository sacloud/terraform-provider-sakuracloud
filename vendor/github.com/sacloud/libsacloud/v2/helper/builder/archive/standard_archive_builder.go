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

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// StandardArchiveBuilder 同一アカウント/同一ゾーンのディスク/アーカイブからアーカイブの作成を行う
type StandardArchiveBuilder struct {
	Name            string
	Description     string
	Tags            types.Tags
	IconID          types.ID
	SourceDiskID    types.ID
	SourceArchiveID types.ID

	NoWait bool
	Client *APIClient
}

// Validate 設定値の検証
func (b *StandardArchiveBuilder) Validate(ctx context.Context, zone string) error {
	requiredValues := map[string]bool{
		"Name":                            b.Name == "",
		"SourceDiskID or SourceArchiveID": b.SourceArchiveID.IsEmpty() && b.SourceDiskID.IsEmpty(),
	}
	for key, empty := range requiredValues {
		if empty {
			return fmt.Errorf("%s is required", key)
		}
	}
	return nil
}

// Build 同一アカウント/同一ゾーンのディスク/アーカイブからアーカイブの作成を行う
func (b *StandardArchiveBuilder) Build(ctx context.Context, zone string) (*sacloud.Archive, error) {
	if err := b.Validate(ctx, zone); err != nil {
		return nil, err
	}

	archive, err := b.Client.Archive.Create(ctx, zone,
		&sacloud.ArchiveCreateRequest{
			Name:            b.Name,
			Description:     b.Description,
			Tags:            b.Tags,
			IconID:          b.IconID,
			SourceDiskID:    b.SourceDiskID,
			SourceArchiveID: b.SourceArchiveID,
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
