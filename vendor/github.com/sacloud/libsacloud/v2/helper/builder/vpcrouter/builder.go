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

package vpcrouter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sacloud/libsacloud/v2/helper/builder"
	"github.com/sacloud/libsacloud/v2/helper/power"
	"github.com/sacloud/libsacloud/v2/helper/setup"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// Builder VPCルータの構築を行う
type Builder struct {
	Name                  string
	Description           string
	Tags                  types.Tags
	IconID                types.ID
	PlanID                types.ID
	NICSetting            NICSettingHolder
	AdditionalNICSettings []AdditionalNICSettingHolder
	RouterSetting         *RouterSetting

	SetupOptions *builder.RetryableSetupParameter
	Client       sacloud.VPCRouterAPI
	NoWait       bool
}

// RouterSetting VPCルータの設定
type RouterSetting struct {
	VRID                      int
	InternetConnectionEnabled types.StringFlag
	StaticNAT                 []*sacloud.VPCRouterStaticNAT
	PortForwarding            []*sacloud.VPCRouterPortForwarding
	Firewall                  []*sacloud.VPCRouterFirewall
	DHCPServer                []*sacloud.VPCRouterDHCPServer
	DHCPStaticMapping         []*sacloud.VPCRouterDHCPStaticMapping
	PPTPServer                *sacloud.VPCRouterPPTPServer
	L2TPIPsecServer           *sacloud.VPCRouterL2TPIPsecServer
	RemoteAccessUsers         []*sacloud.VPCRouterRemoteAccessUser
	SiteToSiteIPsecVPN        []*sacloud.VPCRouterSiteToSiteIPsecVPN
	StaticRoute               []*sacloud.VPCRouterStaticRoute
	SyslogHost                string
}

func (b *Builder) init() {
	if b.SetupOptions == nil {
		b.SetupOptions = builder.DefaultSetupOptions()
	}
	if b.RouterSetting == nil {
		b.RouterSetting = &RouterSetting{
			InternetConnectionEnabled: true,
		}
	}
}

func (b *Builder) getInitInterfaceSettings() []*sacloud.VPCRouterInterfaceSetting {
	s := b.NICSetting.getInterfaceSetting()
	if s != nil {
		return []*sacloud.VPCRouterInterfaceSetting{s}
	}
	return nil
}

func (b *Builder) getInterfaceSettings() []*sacloud.VPCRouterInterfaceSetting {
	var settings []*sacloud.VPCRouterInterfaceSetting
	if s := b.NICSetting.getInterfaceSetting(); s != nil {
		settings = append(settings, s)
	}
	for _, additionalNIC := range b.AdditionalNICSettings {
		settings = append(settings, additionalNIC.getInterfaceSetting())
	}
	return settings
}

// Validate 設定値の検証
func (b *Builder) Validate(ctx context.Context, zone string) error {
	if err := b.validateCommon(ctx, zone); err != nil {
		return err
	}

	if b.NoWait {
		if len(b.AdditionalNICSettings) > 0 || b.RouterSetting != nil {
			return errors.New("NoWait=true is not supported with AdditionalNICSettings and RouterSetting")
		}
		if b.SetupOptions != nil && b.SetupOptions.BootAfterBuild {
			return errors.New("NoWait=true is not supported with SetupOptions.BootAfterBuild")
		}
	}

	switch b.PlanID {
	case types.VPCRouterPlans.Standard:
		return b.validateForStandard(ctx, zone)
	default:
		return b.validateForPremium(ctx, zone)
	}
}

func (b *Builder) validateCommon(ctx context.Context, zone string) error {
	if b.NICSetting == nil {
		return errors.New("required field is missing: NICSetting")
	}
	switch b.PlanID {
	case types.VPCRouterPlans.Standard, types.VPCRouterPlans.Premium, types.VPCRouterPlans.HighSpec, types.VPCRouterPlans.HighSpec4000:
		// noop
	default:
		return fmt.Errorf("invalid plan: PlanID: %s", b.PlanID.String())
	}

	for i, nic := range b.AdditionalNICSettings {
		switchID, index := nic.getSwitchInfo()
		if switchID.IsEmpty() {
			return fmt.Errorf("invalid SwitchID is specified: AdditionalNICSettings[%d].SwitchID is empty", i)
		}
		if index == 0 {
			return fmt.Errorf("invalid SwitchID is specified: AdditionalNICSettings[%d].Index is Zero", i)
		}
	}

	return nil
}

func (b *Builder) validateForStandard(ctx context.Context, zone string) error {
	if _, ok := b.NICSetting.(*StandardNICSetting); !ok {
		return fmt.Errorf("invalid NICSetting is specified: %v", b.NICSetting)
	}
	for i, nic := range b.AdditionalNICSettings {
		if _, ok := nic.(*AdditionalStandardNICSetting); !ok {
			return fmt.Errorf("invalid AdditionalNICSettings is specified: AdditionalNICSettings[%d]:%v", i, nic)
		}
	}

	// Static NAT is only for Premium+
	if b.RouterSetting.StaticNAT != nil {
		return errors.New("invalid RouterSetting is specified: StaticNAT is only for Premium+ plan")
	}
	return nil
}

func (b *Builder) validateForPremium(ctx context.Context, zone string) error {
	if _, ok := b.NICSetting.(*PremiumNICSetting); !ok {
		return fmt.Errorf("invalid NICSetting is specified: %v", b.NICSetting)
	}
	for i, nic := range b.AdditionalNICSettings {
		if _, ok := nic.(*AdditionalPremiumNICSetting); !ok {
			return fmt.Errorf("invalid AdditionalNICSettings is specified: AdditionalNICSettings[%d]:%v", i, nic)
		}
	}
	return nil
}

// Build VPCルータの作成、スイッチの接続をまとめて行う
func (b *Builder) Build(ctx context.Context, zone string) (*sacloud.VPCRouter, error) {
	b.init()

	if err := b.Validate(ctx, zone); err != nil {
		return nil, err
	}

	builder := &setup.RetryableSetup{
		Create: func(ctx context.Context, zone string) (accessor.ID, error) {
			return b.Client.Create(ctx, zone, &sacloud.VPCRouterCreateRequest{
				Name:        b.Name,
				Description: b.Description,
				Tags:        b.Tags,
				IconID:      b.IconID,
				PlanID:      b.PlanID,
				Switch:      b.NICSetting.getConnectedSwitch(),
				IPAddresses: b.NICSetting.getIPAddresses(),
				Settings: &sacloud.VPCRouterSetting{
					VRID:                      b.RouterSetting.VRID,
					InternetConnectionEnabled: b.RouterSetting.InternetConnectionEnabled,
					Interfaces:                b.getInitInterfaceSettings(),
					SyslogHost:                b.RouterSetting.SyslogHost,
				},
			})
		},
		ProvisionBeforeUp: func(ctx context.Context, zone string, id types.ID, target interface{}) error {
			if b.NoWait {
				return nil
			}
			vpcRouter := target.(*sacloud.VPCRouter)

			// スイッチの接続
			for _, additionalNIC := range b.AdditionalNICSettings {
				switchID, index := additionalNIC.getSwitchInfo()
				if err := b.Client.ConnectToSwitch(ctx, zone, id, index, switchID); err != nil {
					return err
				}
			}

			// [HACK] スイッチ接続直後だとエラーになることがあるため数秒待つ
			time.Sleep(b.SetupOptions.NICUpdateWaitDuration)

			// 残りの設定の投入
			_, err := b.Client.UpdateSettings(ctx, zone, id, &sacloud.VPCRouterUpdateSettingsRequest{
				Settings: &sacloud.VPCRouterSetting{
					VRID:                      b.RouterSetting.VRID,
					InternetConnectionEnabled: b.RouterSetting.InternetConnectionEnabled,
					Interfaces:                b.getInterfaceSettings(),
					StaticNAT:                 b.RouterSetting.StaticNAT,
					PortForwarding:            b.RouterSetting.PortForwarding,
					Firewall:                  b.RouterSetting.Firewall,
					DHCPServer:                b.RouterSetting.DHCPServer,
					DHCPStaticMapping:         b.RouterSetting.DHCPStaticMapping,
					PPTPServer:                b.RouterSetting.PPTPServer,
					PPTPServerEnabled:         b.RouterSetting.PPTPServer != nil,
					L2TPIPsecServer:           b.RouterSetting.L2TPIPsecServer,
					L2TPIPsecServerEnabled:    b.RouterSetting.L2TPIPsecServer != nil,
					RemoteAccessUsers:         b.RouterSetting.RemoteAccessUsers,
					SiteToSiteIPsecVPN:        b.RouterSetting.SiteToSiteIPsecVPN,
					StaticRoute:               b.RouterSetting.StaticRoute,
					SyslogHost:                b.RouterSetting.SyslogHost,
				},
				SettingsHash: vpcRouter.SettingsHash,
			})
			if err != nil {
				return err
			}
			if err := b.Client.Config(ctx, zone, id); err != nil {
				return err
			}

			if b.SetupOptions.BootAfterBuild {
				return power.BootVPCRouter(ctx, b.Client, zone, id)
			}
			return nil
		},
		Delete: func(ctx context.Context, zone string, id types.ID) error {
			return b.Client.Delete(ctx, zone, id)
		},
		Read: func(ctx context.Context, zone string, id types.ID) (interface{}, error) {
			return b.Client.Read(ctx, zone, id)
		},
		IsWaitForCopy:             !b.NoWait,
		IsWaitForUp:               !b.NoWait && b.SetupOptions.BootAfterBuild,
		RetryCount:                b.SetupOptions.RetryCount,
		ProvisioningRetryCount:    1,
		ProvisioningRetryInterval: b.SetupOptions.ProvisioningRetryInterval,
		DeleteRetryCount:          b.SetupOptions.DeleteRetryCount,
		DeleteRetryInterval:       b.SetupOptions.DeleteRetryInterval,
		PollingInterval:           b.SetupOptions.PollingInterval,
	}

	result, err := builder.Setup(ctx, zone)
	var vpcRouter *sacloud.VPCRouter
	if result != nil {
		vpcRouter = result.(*sacloud.VPCRouter)
	}
	if err != nil {
		return vpcRouter, err
	}

	// refresh
	refreshed, err := b.Client.Read(ctx, zone, vpcRouter.ID)
	if err != nil {
		return vpcRouter, err
	}
	return refreshed, nil
}

// Update VPCルータの更新(再起動を伴う場合あり)
//
// 接続先スイッチが変更されていた場合、VPCルータの再起動が行われます。
func (b *Builder) Update(ctx context.Context, zone string, id types.ID) (*sacloud.VPCRouter, error) {
	b.init()

	if err := b.Validate(ctx, zone); err != nil {
		return nil, err
	}

	// check VPCRouter is exists
	vpcRouter, err := b.Client.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}

	isNeedShutdown, err := b.collectUpdateInfo(vpcRouter)
	if err != nil {
		return nil, err
	}

	isNeedRestart := false
	if vpcRouter.InstanceStatus.IsUp() && isNeedShutdown {
		if b.NoWait {
			return nil, errors.New("NoWait option is not available due to the need to shut down")
		}

		isNeedRestart = true
		if err := power.ShutdownVPCRouter(ctx, b.Client, zone, id, false); err != nil {
			return nil, err
		}
	}

	// NICの切断/変更(変更分のみ)
	for _, iface := range vpcRouter.Interfaces {
		if iface.Index == 0 {
			continue
		}

		newSwitchID := b.findAdditionalSwitchSettingByIndex(iface.Index) // 削除されていた場合types.ID(0)が返る
		if iface.SwitchID != newSwitchID {
			// disconnect
			if err := b.Client.DisconnectFromSwitch(ctx, zone, id, iface.Index); err != nil {
				return nil, err
			}
			// connect
			if !newSwitchID.IsEmpty() {
				if err := b.Client.ConnectToSwitch(ctx, zone, id, iface.Index, newSwitchID); err != nil {
					return nil, err
				}
			}
		}
	}

	// 追加されたNICの接続
	for _, nicSetting := range b.AdditionalNICSettings {
		switchID, index := nicSetting.getSwitchInfo()
		iface := b.findInterfaceByIndex(vpcRouter, index)
		if iface == nil {
			if err := b.Client.ConnectToSwitch(ctx, zone, id, index, switchID); err != nil {
				return nil, err
			}
		}
	}
	// [HACK] スイッチ接続直後だとエラーになることがあるため数秒待つ
	time.Sleep(b.SetupOptions.NICUpdateWaitDuration)

	_, err = b.Client.Update(ctx, zone, id, &sacloud.VPCRouterUpdateRequest{
		Name:        b.Name,
		Description: b.Description,
		Tags:        b.Tags,
		IconID:      b.IconID,
		Settings: &sacloud.VPCRouterSetting{
			VRID:                      b.RouterSetting.VRID,
			InternetConnectionEnabled: b.RouterSetting.InternetConnectionEnabled,
			Interfaces:                b.getInterfaceSettings(),
			StaticNAT:                 b.RouterSetting.StaticNAT,
			PortForwarding:            b.RouterSetting.PortForwarding,
			Firewall:                  b.RouterSetting.Firewall,
			DHCPServer:                b.RouterSetting.DHCPServer,
			DHCPStaticMapping:         b.RouterSetting.DHCPStaticMapping,
			PPTPServer:                b.RouterSetting.PPTPServer,
			PPTPServerEnabled:         b.RouterSetting.PPTPServer != nil,
			L2TPIPsecServer:           b.RouterSetting.L2TPIPsecServer,
			L2TPIPsecServerEnabled:    b.RouterSetting.L2TPIPsecServer != nil,
			RemoteAccessUsers:         b.RouterSetting.RemoteAccessUsers,
			SiteToSiteIPsecVPN:        b.RouterSetting.SiteToSiteIPsecVPN,
			StaticRoute:               b.RouterSetting.StaticRoute,
			SyslogHost:                b.RouterSetting.SyslogHost,
		},
		SettingsHash: vpcRouter.SettingsHash,
	})
	if err != nil {
		return nil, err
	}

	if err := b.Client.Config(ctx, zone, id); err != nil {
		return nil, err
	}

	if isNeedRestart {
		if err := power.BootVPCRouter(ctx, b.Client, zone, id); err != nil {
			return nil, err
		}
	}
	// refresh
	vpcRouter, err = b.Client.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}
	return vpcRouter, err
}

func (b *Builder) collectUpdateInfo(vpcRouter *sacloud.VPCRouter) (isNeedShutdown bool, err error) {
	// プランの変更はエラーとする
	if vpcRouter.PlanID != b.PlanID {
		err = fmt.Errorf("unsupported operation: VPCRouter is not allowd changing Plan: currentPlan: %s", vpcRouter.PlanID.String())
		return
	}

	// スイッチの変更/削除は再起動が必要
	for _, iface := range vpcRouter.Interfaces {
		if iface.Index == 0 {
			continue
		}
		newSwitchID := b.findAdditionalSwitchSettingByIndex(iface.Index) // 削除された場合はtypes.ID(0)が返る
		isNeedShutdown = iface.SwitchID != newSwitchID
	}
	if isNeedShutdown {
		return
	}

	// スイッチの増設は再起動が必要
	if len(vpcRouter.Interfaces)-1 != len(b.AdditionalNICSettings) {
		isNeedShutdown = true
	}
	return
}

func (b *Builder) findInterfaceByIndex(vpcRouter *sacloud.VPCRouter, ifIndex int) *sacloud.VPCRouterInterface {
	for _, iface := range vpcRouter.Interfaces {
		if iface.Index == ifIndex {
			return iface
		}
	}
	return nil
}

func (b *Builder) findAdditionalSwitchSettingByIndex(ifIndex int) types.ID {
	for _, nic := range b.AdditionalNICSettings {
		switchID, index := nic.getSwitchInfo()
		if index == ifIndex {
			return switchID
		}
	}
	return types.ID(0)
}
