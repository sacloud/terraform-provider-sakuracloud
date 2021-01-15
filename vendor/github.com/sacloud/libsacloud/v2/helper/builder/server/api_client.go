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

package server

import (
	"context"

	"github.com/sacloud/libsacloud/v2/helper/query"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// APIClient builderが利用するAPIクライアント群
type APIClient struct {
	Disk         DiskHandler
	Interface    InterfaceHandler
	PacketFilter PacketFilterReader
	Server       CreateServerHandler
	ServerPlan   query.ServerPlanFinder
	Switch       SwitchReader
}

// DiskHandler ディスクの接続/切断のためのインターフェース
type DiskHandler interface {
	ConnectToServer(ctx context.Context, zone string, id types.ID, serverID types.ID) error
	DisconnectFromServer(ctx context.Context, zone string, id types.ID) error
}

// SwitchReader スイッチ参照のためのインターフェース
type SwitchReader interface {
	Read(ctx context.Context, zone string, id types.ID) (*sacloud.Switch, error)
}

// InterfaceHandler NIC操作のためのインターフェース
type InterfaceHandler interface {
	Create(ctx context.Context, zone string, param *sacloud.InterfaceCreateRequest) (*sacloud.Interface, error)
	Update(ctx context.Context, zone string, id types.ID, param *sacloud.InterfaceUpdateRequest) (*sacloud.Interface, error)
	Delete(ctx context.Context, zone string, id types.ID) error
	ConnectToSharedSegment(ctx context.Context, zone string, id types.ID) error
	ConnectToSwitch(ctx context.Context, zone string, id types.ID, switchID types.ID) error
	DisconnectFromSwitch(ctx context.Context, zone string, id types.ID) error
	ConnectToPacketFilter(ctx context.Context, zone string, id types.ID, packetFilterID types.ID) error
	DisconnectFromPacketFilter(ctx context.Context, zone string, id types.ID) error
}

// PacketFilterReader パケットフィルタ参照のためのインターフェース
type PacketFilterReader interface {
	Read(ctx context.Context, zone string, id types.ID) (*sacloud.PacketFilter, error)
}

// CreateServerHandler サーバ操作のためのインターフェース
type CreateServerHandler interface {
	Create(ctx context.Context, zone string, param *sacloud.ServerCreateRequest) (*sacloud.Server, error)
	Update(ctx context.Context, zone string, id types.ID, param *sacloud.ServerUpdateRequest) (*sacloud.Server, error)
	Read(ctx context.Context, zone string, id types.ID) (*sacloud.Server, error)
	InsertCDROM(ctx context.Context, zone string, id types.ID, insertParam *sacloud.InsertCDROMRequest) error
	EjectCDROM(ctx context.Context, zone string, id types.ID, ejectParam *sacloud.EjectCDROMRequest) error
	Boot(ctx context.Context, zone string, id types.ID) error
	Shutdown(ctx context.Context, zone string, id types.ID, shutdownOption *sacloud.ShutdownOption) error
	ChangePlan(ctx context.Context, zone string, id types.ID, plan *sacloud.ServerChangePlanRequest) (*sacloud.Server, error)
}

// NewBuildersAPIClient APIクライアントの作成
func NewBuildersAPIClient(caller sacloud.APICaller) *APIClient {
	return &APIClient{
		Disk:         sacloud.NewDiskOp(caller),
		Interface:    sacloud.NewInterfaceOp(caller),
		PacketFilter: sacloud.NewPacketFilterOp(caller),
		Server:       sacloud.NewServerOp(caller),
		ServerPlan:   sacloud.NewServerPlanOp(caller),
		Switch:       sacloud.NewSwitchOp(caller),
	}
}
