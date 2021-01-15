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
	"errors"
	"fmt"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

type nicState struct {
	upstreamType   types.EUpstreamNetworkType
	switchID       types.ID
	packetFilterID types.ID
	displayIP      string
}

// NICSettingHolder NIC設定を保持するためのインターフェース
type NICSettingHolder interface {
	GetConnectedSwitchParam() *sacloud.ConnectedSwitch
	GetDisplayIPAddress() string

	GetPacketFilterID() types.ID
	Validate(ctx context.Context, client *APIClient, zone string) error
	state() *nicState
}

// AdditionalNICSettingHolder 追加NIC設定を保持するためのインターフェース
type AdditionalNICSettingHolder interface {
	GetSwitchID() types.ID

	GetDisplayIPAddress() string
	GetPacketFilterID() types.ID
	Validate(ctx context.Context, client *APIClient, zone string) error
	state() *nicState
}

// SharedNICSetting サーバ作成時に共有セグメントに接続するためのパラメータ
//
// NICSettingHolderを実装し、Builder.NICに利用できる。
type SharedNICSetting struct {
	PacketFilterID types.ID
}

// GetConnectedSwitchParam サーバ作成時の接続先指定パラメータを作成して返す
func (c *SharedNICSetting) GetConnectedSwitchParam() *sacloud.ConnectedSwitch {
	return &sacloud.ConnectedSwitch{Scope: types.Scopes.Shared}
}

// GetPacketFilterID このNICに接続するパケットフィルタのIDを返す
func (c *SharedNICSetting) GetPacketFilterID() types.ID {
	return c.PacketFilterID
}

// Validate 設定値の検証
func (c *SharedNICSetting) Validate(ctx context.Context, client *APIClient, zone string) error {
	if !c.PacketFilterID.IsEmpty() {
		if _, err := client.PacketFilter.Read(ctx, zone, c.PacketFilterID); err != nil {
			return fmt.Errorf("reading packet filter info(id:%d) is failed: %s", c.PacketFilterID, err)
		}
	}
	return nil
}

// GetDisplayIPAddress 表示用IPアドレスを返す
func (c *SharedNICSetting) GetDisplayIPAddress() string {
	return ""
}

func (c *SharedNICSetting) state() *nicState {
	return &nicState{
		upstreamType:   types.UpstreamNetworkTypes.Shared,
		switchID:       types.ID(0),
		packetFilterID: c.PacketFilterID,
		displayIP:      "",
	}
}

// ConnectedNICSetting サーバ作成時にスイッチに接続するためのパラメータ
//
// NICSettingHolderとAdditionalNICSettingHolderを実装し、Builder.NIC/Builder.AdditionalNICsに利用できる。
type ConnectedNICSetting struct {
	SwitchID         types.ID
	DisplayIPAddress string
	PacketFilterID   types.ID
}

// GetConnectedSwitchParam サーバ作成時の接続先指定パラメータを作成して返す
func (c *ConnectedNICSetting) GetConnectedSwitchParam() *sacloud.ConnectedSwitch {
	return &sacloud.ConnectedSwitch{ID: c.SwitchID}
}

// GetSwitchID このNICが接続するスイッチのIDを返す
func (c *ConnectedNICSetting) GetSwitchID() types.ID {
	return c.SwitchID
}

// GetDisplayIPAddress 表示用IPアドレスを返す
func (c *ConnectedNICSetting) GetDisplayIPAddress() string {
	return c.DisplayIPAddress
}

// GetPacketFilterID このNICに接続するパケットフィルタのIDを返す
func (c *ConnectedNICSetting) GetPacketFilterID() types.ID {
	return c.PacketFilterID
}

// Validate 設定値の検証
func (c *ConnectedNICSetting) Validate(ctx context.Context, client *APIClient, zone string) error {
	if c.SwitchID.IsEmpty() {
		return errors.New("ConnectedNICSetting: SwitchID is required")
	}

	if _, err := client.Switch.Read(ctx, zone, c.SwitchID); err != nil {
		return fmt.Errorf("reading switch info(id:%d) is failed: %s", c.SwitchID, err)
	}

	if !c.PacketFilterID.IsEmpty() {
		if _, err := client.PacketFilter.Read(ctx, zone, c.PacketFilterID); err != nil {
			return fmt.Errorf("reading packet filter info(id:%d) is failed: %s", c.PacketFilterID, err)
		}
	}

	return nil
}

func (c *ConnectedNICSetting) state() *nicState {
	return &nicState{
		upstreamType:   types.UpstreamNetworkTypes.Switch,
		switchID:       c.SwitchID,
		packetFilterID: c.PacketFilterID,
		displayIP:      c.DisplayIPAddress,
	}
}

// DisconnectedNICSetting 切断状態のNICを作成するためのパラメータ
//
// NICSettingHolderとAdditionalNICSettingHolderを実装し、Builder.NIC/Builder.AdditionalNICsに利用できる。
type DisconnectedNICSetting struct{}

// GetConnectedSwitchParam サーバ作成時の接続先指定パラメータを作成して返す
func (d *DisconnectedNICSetting) GetConnectedSwitchParam() *sacloud.ConnectedSwitch {
	return nil
}

// GetSwitchID このNICが接続するスイッチのIDを返す
func (d *DisconnectedNICSetting) GetSwitchID() types.ID {
	return types.ID(0)
}

// GetDisplayIPAddress 表示用IPアドレスを返す
func (d *DisconnectedNICSetting) GetDisplayIPAddress() string {
	return ""
}

// GetPacketFilterID このNICに接続するパケットフィルタのIDを返す
func (d *DisconnectedNICSetting) GetPacketFilterID() types.ID {
	return types.ID(0)
}

// Validate 設定値の検証
func (d *DisconnectedNICSetting) Validate(ctx context.Context, client *APIClient, zone string) error {
	return nil
}

func (d *DisconnectedNICSetting) state() *nicState {
	return &nicState{
		upstreamType:   types.UpstreamNetworkTypes.None,
		switchID:       types.ID(0),
		packetFilterID: types.ID(0),
		displayIP:      "",
	}
}
