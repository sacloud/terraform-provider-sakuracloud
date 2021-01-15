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

package disk

import (
	"context"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// APIClient builderが利用するAPIクライアント群
type APIClient struct {
	Archive  ArchiveFinder
	Disk     CreateDiskHandler
	DiskPlan PlanReader
	Note     NoteHandler
	SSHKey   SSHKeyHandler
}

// ArchiveFinder アーカイブ検索のためのインターフェース
type ArchiveFinder interface {
	Find(ctx context.Context, zone string, conditions *sacloud.FindCondition) (*sacloud.ArchiveFindResult, error)
	Read(ctx context.Context, zone string, id types.ID) (*sacloud.Archive, error)
}

// CreateDiskHandler ディスク操作のためのインターフェース
type CreateDiskHandler interface {
	Create(ctx context.Context, zone string, createParam *sacloud.DiskCreateRequest, distantFrom []types.ID) (*sacloud.Disk, error)
	CreateWithConfig(
		ctx context.Context,
		zone string,
		createParam *sacloud.DiskCreateRequest,
		editParam *sacloud.DiskEditRequest,
		bootAtAvailable bool,
		distantFrom []types.ID,
	) (*sacloud.Disk, error)
	Update(ctx context.Context, zone string, id types.ID, updateParam *sacloud.DiskUpdateRequest) (*sacloud.Disk, error)
	Config(ctx context.Context, zone string, id types.ID, editParam *sacloud.DiskEditRequest) error
	Read(ctx context.Context, zone string, id types.ID) (*sacloud.Disk, error)
	ConnectToServer(ctx context.Context, zone string, id types.ID, serverID types.ID) error
}

// PlanReader ディスクプラン取得のためのインターフェース
type PlanReader interface {
	Read(ctx context.Context, zone string, id types.ID) (*sacloud.DiskPlan, error)
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
func NewBuildersAPIClient(caller sacloud.APICaller) *APIClient {
	return &APIClient{
		Archive:  sacloud.NewArchiveOp(caller),
		Disk:     sacloud.NewDiskOp(caller),
		DiskPlan: sacloud.NewDiskPlanOp(caller),
		Note:     sacloud.NewNoteOp(caller),
		SSHKey:   sacloud.NewSSHKeyOp(caller),
	}
}
