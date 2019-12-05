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

package server

import (
	"context"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// BuildersAPIClient builderが利用するAPIクライアント群
type BuildersAPIClient struct {
	Archive      ArchiveFinder
	Disk         DiskHandler
	DiskPlan     DiskPlanReader
	Interface    InterfaceHandler
	PacketFilter PacketFilterReader
	Server       CreateServerHandler
	ServerPlan   PlanFinder
	Switch       SwitchReader
	Note         NoteHandler
	SSHKey       SSHKeyHandler
}

// ArchiveFinder アーカイブ検索のためのインターフェース
type ArchiveFinder interface {
	Find(ctx context.Context, zone string, conditions *sacloud.FindCondition) (*sacloud.ArchiveFindResult, error)
	Read(ctx context.Context, zone string, id types.ID) (*sacloud.Archive, error)
}

// DiskHandler ディスク操作のためのインターフェース
type DiskHandler interface {
	Create(ctx context.Context, zone string, createParam *sacloud.DiskCreateRequest, distantFrom []types.ID) (*sacloud.Disk, error)
	CreateWithConfig(
		ctx context.Context,
		zone string,
		createParam *sacloud.DiskCreateRequest,
		editParam *sacloud.DiskEditRequest,
		bootAtAvailable bool,
		distantFrom []types.ID,
	) (*sacloud.Disk, error)
	Read(ctx context.Context, zone string, id types.ID) (*sacloud.Disk, error)
}

// DiskPlanReader ディスクプラン取得のためのインターフェース
type DiskPlanReader interface {
	Read(ctx context.Context, zone string, id types.ID) (*sacloud.DiskPlan, error)
}

// SwitchReader スイッチ参照のためのインターフェース
type SwitchReader interface {
	Read(ctx context.Context, zone string, id types.ID) (*sacloud.Switch, error)
}

// InterfaceHandler NIC操作のためのインターフェース
type InterfaceHandler interface {
	ConnectToPacketFilter(ctx context.Context, zone string, id types.ID, packetFilterID types.ID) error
	Update(ctx context.Context, zone string, id types.ID, param *sacloud.InterfaceUpdateRequest) (*sacloud.Interface, error)
}

// PacketFilterReader パケットフィルタ参照のためのインターフェース
type PacketFilterReader interface {
	Read(ctx context.Context, zone string, id types.ID) (*sacloud.PacketFilter, error)
}

// CreateServerHandler サーバ操作のためのインターフェース
type CreateServerHandler interface {
	Create(ctx context.Context, zone string, param *sacloud.ServerCreateRequest) (*sacloud.Server, error)
	Read(ctx context.Context, zone string, id types.ID) (*sacloud.Server, error)
	InsertCDROM(ctx context.Context, zone string, id types.ID, insertParam *sacloud.InsertCDROMRequest) error
	Boot(ctx context.Context, zone string, id types.ID) error
}

// NoteHandler スタートアップスクリプト参照のためのインターフェース
type NoteHandler interface {
	Read(ctx context.Context, id types.ID) (*sacloud.Note, error)
	Create(ctx context.Context, param *sacloud.NoteCreateRequest) (*sacloud.Note, error)
	Delete(ctx context.Context, id types.ID) error
}

// SSHKeyHandler SSHKey参照のためのインターフェース
type SSHKeyHandler interface {
	Read(ctx context.Context, id types.ID) (*sacloud.SSHKey, error)
	Generate(ctx context.Context, param *sacloud.SSHKeyGenerateRequest) (*sacloud.SSHKeyGenerated, error)
	Delete(ctx context.Context, id types.ID) error
}

// NewBuildersAPIClient APIクライアントの作成
func NewBuildersAPIClient(caller sacloud.APICaller) *BuildersAPIClient {
	return &BuildersAPIClient{
		Archive:      sacloud.NewArchiveOp(caller),
		Disk:         sacloud.NewDiskOp(caller),
		DiskPlan:     sacloud.NewDiskPlanOp(caller),
		Interface:    sacloud.NewInterfaceOp(caller),
		PacketFilter: sacloud.NewPacketFilterOp(caller),
		Server:       sacloud.NewServerOp(caller),
		ServerPlan:   sacloud.NewServerPlanOp(caller),
		Switch:       sacloud.NewSwitchOp(caller),
		Note:         sacloud.NewNoteOp(caller),
		SSHKey:       sacloud.NewSSHKeyOp(caller),
	}
}
