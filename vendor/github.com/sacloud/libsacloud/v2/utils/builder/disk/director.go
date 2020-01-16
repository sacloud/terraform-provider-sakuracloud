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

package disk

import (
	"github.com/sacloud/libsacloud/v2/sacloud/ostype"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// Director パラメータに応じて適切なDiskBuilderを構築する
type Director struct {
	OSType ostype.ArchiveOSType

	Name        string
	SizeGB      int
	DistantFrom []types.ID
	PlanID      types.ID
	Connection  types.EDiskConnection
	Description string
	Tags        types.Tags
	IconID      types.ID

	DiskID          types.ID
	SourceDiskID    types.ID
	SourceArchiveID types.ID

	EditParameter *EditRequest
	Client        *APIClient
}

// Builder パラメータに応じて適切なDiskBuilderを返す
func (d *Director) Builder() Builder {
	switch {
	case d.OSType == ostype.Custom:
		switch {
		case !d.DiskID.IsEmpty():
			return &ConnectedDiskBuilder{
				ID:            d.DiskID,
				EditParameter: d.EditParameter.ToUnixDiskEditRequest(),
				Client:        d.Client,
			}
		case !d.SourceDiskID.IsEmpty(), !d.SourceArchiveID.IsEmpty():
			return &FromDiskOrArchiveBuilder{
				SourceDiskID:    d.SourceDiskID,
				SourceArchiveID: d.SourceArchiveID,
				Name:            d.Name,
				SizeGB:          d.SizeGB,
				DistantFrom:     d.DistantFrom,
				PlanID:          d.PlanID,
				Connection:      d.Connection,
				Description:     d.Description,
				Tags:            d.Tags,
				IconID:          d.IconID,
				EditParameter:   d.EditParameter.ToUnixDiskEditRequest(),
				Client:          d.Client,
			}
		default:
			return &BlankBuilder{
				Name:        d.Name,
				SizeGB:      d.SizeGB,
				DistantFrom: d.DistantFrom,
				PlanID:      d.PlanID,
				Connection:  d.Connection,
				Description: d.Description,
				Tags:        d.Tags,
				IconID:      d.IconID,
				Client:      d.Client,
			}
		}
	case d.OSType.IsSupportDiskEdit():
		return &FromUnixBuilder{
			OSType:        d.OSType,
			Name:          d.Name,
			SizeGB:        d.SizeGB,
			DistantFrom:   d.DistantFrom,
			PlanID:        d.PlanID,
			Connection:    d.Connection,
			Description:   d.Description,
			Tags:          d.Tags,
			IconID:        d.IconID,
			EditParameter: d.EditParameter.ToUnixDiskEditRequest(),
			Client:        d.Client,
		}
	case d.OSType.IsWindows():
		return &FromWindowsBuilder{
			OSType:        d.OSType,
			Name:          d.Name,
			SizeGB:        d.SizeGB,
			DistantFrom:   d.DistantFrom,
			PlanID:        d.PlanID,
			Connection:    d.Connection,
			Description:   d.Description,
			Tags:          d.Tags,
			IconID:        d.IconID,
			EditParameter: d.EditParameter.ToWindowsDiskEditRequest(),
			Client:        d.Client,
		}
	default:
		// 現在はOSTypeにディスクの修正不可のアーカイブはないためここには到達しない
		return &FromFixedArchiveBuilder{
			OSType:      d.OSType,
			Name:        d.Name,
			SizeGB:      d.SizeGB,
			DistantFrom: d.DistantFrom,
			PlanID:      d.PlanID,
			Connection:  d.Connection,
			Description: d.Description,
			Tags:        d.Tags,
			IconID:      d.IconID,
			Client:      d.Client,
		}
	}
}
