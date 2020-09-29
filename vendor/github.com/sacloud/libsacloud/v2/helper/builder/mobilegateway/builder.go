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

package mobilegateway

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/sacloud/libsacloud/v2/helper/builder"
	"github.com/sacloud/libsacloud/v2/helper/power"
	"github.com/sacloud/libsacloud/v2/helper/setup"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// Builder モバイルゲートウェイの構築を行う
type Builder struct {
	Name                            string
	Description                     string
	Tags                            types.Tags
	IconID                          types.ID
	PrivateInterface                *PrivateInterfaceSetting
	StaticRoutes                    []*sacloud.MobileGatewayStaticRoute
	SIMRoutes                       []*SIMRouteSetting
	InternetConnectionEnabled       bool
	InterDeviceCommunicationEnabled bool
	DNS                             *sacloud.MobileGatewayDNSSetting
	SIMs                            []*SIMSetting
	TrafficConfig                   *sacloud.MobileGatewayTrafficControl

	SetupOptions *builder.RetryableSetupParameter
	Client       *APIClient
}

// PrivateInterfaceSetting モバイルゲートウェイのプライベート側インターフェース設定
type PrivateInterfaceSetting struct {
	SwitchID       types.ID
	IPAddress      string
	NetworkMaskLen int
}

// SIMSetting モバイルゲートウェイに接続するSIM設定
type SIMSetting struct {
	SIMID     types.ID
	IPAddress string
}

// SIMRouteSetting SIMルート設定
type SIMRouteSetting struct {
	SIMID  types.ID
	Prefix string
}

func (b *Builder) init() {
	if b.SetupOptions == nil {
		b.SetupOptions = builder.DefaultSetupOptions()
	}
}

// Validate 設定値の検証
func (b *Builder) Validate(ctx context.Context, zone string) error {
	if b.PrivateInterface != nil {
		if b.PrivateInterface.SwitchID.IsEmpty() {
			return fmt.Errorf("switch id is required when specified private interface")
		}
		if b.PrivateInterface.IPAddress == "" {
			return fmt.Errorf("ip address is required when specified private interface")
		}
		if b.PrivateInterface.NetworkMaskLen == 0 {
			return fmt.Errorf("ip address is required when specified private interface")
		}
	}
	if len(b.SIMRoutes) > 0 && len(b.SIMs) == 0 {
		return fmt.Errorf("sim settings are required when specified sim routes")
	}
	return nil
}

// Build モバイルゲートウェイの作成や設定をまとめて行う
func (b *Builder) Build(ctx context.Context, zone string) (*sacloud.MobileGateway, error) {
	b.init()

	if err := b.Validate(ctx, zone); err != nil {
		return nil, err
	}

	builder := &setup.RetryableSetup{
		Create: func(ctx context.Context, zone string) (accessor.ID, error) {
			return b.Client.MobileGateway.Create(ctx, zone, &sacloud.MobileGatewayCreateRequest{
				Name:                            b.Name,
				Description:                     b.Description,
				Tags:                            b.Tags,
				IconID:                          b.IconID,
				InternetConnectionEnabled:       types.StringFlag(b.InternetConnectionEnabled),
				InterDeviceCommunicationEnabled: types.StringFlag(b.InterDeviceCommunicationEnabled),
			})
		},
		ProvisionBeforeUp: func(ctx context.Context, zone string, id types.ID, target interface{}) error {
			mgw := target.(*sacloud.MobileGateway)

			// スイッチの接続
			if b.PrivateInterface != nil {
				if err := b.Client.MobileGateway.ConnectToSwitch(ctx, zone, id, b.PrivateInterface.SwitchID); err != nil {
					return err
				}
			}

			// [HACK] スイッチ接続直後だとエラーになることがあるため数秒待つ
			time.Sleep(b.SetupOptions.NICUpdateWaitDuration)

			// Interface設定
			updated, err := b.Client.MobileGateway.UpdateSettings(ctx, zone, id, &sacloud.MobileGatewayUpdateSettingsRequest{
				InterfaceSettings:               b.getInterfaceSettings(),
				InternetConnectionEnabled:       types.StringFlag(b.InternetConnectionEnabled),
				InterDeviceCommunicationEnabled: types.StringFlag(b.InterDeviceCommunicationEnabled),
				SettingsHash:                    mgw.SettingsHash,
			})
			if err != nil {
				return err
			}
			// [HACK] インターフェースの設定をConfigで反映させておかないとエラーになることへの対応
			// see: https://github.com/sacloud/libsacloud/issues/589
			if err := b.Client.MobileGateway.Config(ctx, zone, id); err != nil {
				return err
			}
			mgw = updated

			// traffic config
			if b.TrafficConfig != nil {
				if err := b.Client.MobileGateway.SetTrafficConfig(ctx, zone, id, b.TrafficConfig); err != nil {
					return err
				}
			}

			// dns
			if b.DNS != nil {
				if err := b.Client.MobileGateway.SetDNS(ctx, zone, id, b.DNS); err != nil {
					return err
				}
			}

			// static route
			if len(b.StaticRoutes) > 0 {
				_, err := b.Client.MobileGateway.UpdateSettings(ctx, zone, id, &sacloud.MobileGatewayUpdateSettingsRequest{
					InterfaceSettings:               b.getInterfaceSettings(),
					StaticRoutes:                    b.StaticRoutes,
					InternetConnectionEnabled:       types.StringFlag(b.InternetConnectionEnabled),
					InterDeviceCommunicationEnabled: types.StringFlag(b.InterDeviceCommunicationEnabled),
					SettingsHash:                    mgw.SettingsHash,
				})
				if err != nil {
					return err
				}
			}

			// SIMs
			if len(b.SIMs) > 0 {
				for _, sim := range b.SIMs {
					if err := b.Client.MobileGateway.AddSIM(ctx, zone, id, &sacloud.MobileGatewayAddSIMRequest{SIMID: sim.SIMID.String()}); err != nil {
						return err
					}
					if err := b.Client.SIM.AssignIP(ctx, sim.SIMID, &sacloud.SIMAssignIPRequest{IP: sim.IPAddress}); err != nil {
						return err
					}
				}
			}

			// SIM routes
			if len(b.SIMRoutes) > 0 {
				if err := b.Client.MobileGateway.SetSIMRoutes(ctx, zone, id, b.getSIMRouteSettings()); err != nil {
					return err
				}
			}

			if err := b.Client.MobileGateway.Config(ctx, zone, id); err != nil {
				return err
			}

			if b.SetupOptions.BootAfterBuild {
				return power.BootMobileGateway(ctx, b.Client.MobileGateway, zone, id)
			}
			return nil
		},
		Delete: func(ctx context.Context, zone string, id types.ID) error {
			return b.Client.MobileGateway.Delete(ctx, zone, id)
		},
		Read: func(ctx context.Context, zone string, id types.ID) (interface{}, error) {
			return b.Client.MobileGateway.Read(ctx, zone, id)
		},
		IsWaitForCopy:             true,
		IsWaitForUp:               b.SetupOptions.BootAfterBuild,
		RetryCount:                b.SetupOptions.RetryCount,
		ProvisioningRetryCount:    1,
		ProvisioningRetryInterval: b.SetupOptions.ProvisioningRetryInterval,
		DeleteRetryCount:          b.SetupOptions.DeleteRetryCount,
		DeleteRetryInterval:       b.SetupOptions.DeleteRetryInterval,
		PollingInterval:           b.SetupOptions.PollingInterval,
	}

	result, err := builder.Setup(ctx, zone)
	var mgw *sacloud.MobileGateway
	if result != nil {
		mgw = result.(*sacloud.MobileGateway)
	}
	if err != nil {
		return mgw, err
	}

	// refresh
	refreshed, err := b.Client.MobileGateway.Read(ctx, zone, mgw.ID)
	if err != nil {
		return mgw, err
	}
	return refreshed, nil
}

// Update モバイルゲートウェイの更新
//
// 更新中、SIMルートが一時的にクリアされます。また、接続先スイッチが変更されていた場合は再起動されます。
func (b *Builder) Update(ctx context.Context, zone string, id types.ID) (*sacloud.MobileGateway, error) {
	b.init()

	if err := b.Validate(ctx, zone); err != nil {
		return nil, err
	}

	// check MobileGateway is exists
	mgw, err := b.Client.MobileGateway.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}

	isNeedShutdown, err := b.collectUpdateInfo(mgw)
	if err != nil {
		return nil, err
	}

	isNeedRestart := false
	if mgw.InstanceStatus.IsUp() && isNeedShutdown {
		isNeedRestart = true
		if err := power.ShutdownMobileGateway(ctx, b.Client.MobileGateway, zone, id, false); err != nil {
			return nil, err
		}
	}

	// NICの切断/変更
	if b.isPrivateInterfaceChanged(mgw) {
		if len(mgw.Interfaces) > 1 && !mgw.Interfaces[1].SwitchID.IsEmpty() {
			// 切断
			if err := b.Client.MobileGateway.DisconnectFromSwitch(ctx, zone, id); err != nil {
				return nil, err
			}
			// [HACK] スイッチ接続直後だとエラーになることがあるため数秒待つ
			time.Sleep(b.SetupOptions.NICUpdateWaitDuration)

			updated, err := b.Client.MobileGateway.UpdateSettings(ctx, zone, id, &sacloud.MobileGatewayUpdateSettingsRequest{
				InternetConnectionEnabled:       types.StringFlag(b.InternetConnectionEnabled),
				InterDeviceCommunicationEnabled: types.StringFlag(b.InterDeviceCommunicationEnabled),
				SettingsHash:                    mgw.SettingsHash,
			})
			if err != nil {
				return nil, err
			}
			// [HACK] インターフェースの設定をConfigで反映させておかないとエラーになることへの対応
			// see: https://github.com/sacloud/libsacloud/issues/589
			if err := b.Client.MobileGateway.Config(ctx, zone, id); err != nil {
				return nil, err
			}
			mgw = updated
		}

		// 接続
		if b.PrivateInterface != nil {
			// スイッチの接続
			if err := b.Client.MobileGateway.ConnectToSwitch(ctx, zone, id, b.PrivateInterface.SwitchID); err != nil {
				return nil, err
			}

			// [HACK] スイッチ接続直後だとエラーになることがあるため数秒待つ
			time.Sleep(b.SetupOptions.NICUpdateWaitDuration)

			// Interface設定
			updated, err := b.Client.MobileGateway.UpdateSettings(ctx, zone, id, &sacloud.MobileGatewayUpdateSettingsRequest{
				InterfaceSettings:               b.getInterfaceSettings(),
				InternetConnectionEnabled:       types.StringFlag(b.InternetConnectionEnabled),
				InterDeviceCommunicationEnabled: types.StringFlag(b.InterDeviceCommunicationEnabled),
				SettingsHash:                    mgw.SettingsHash,
			})
			if err != nil {
				return nil, err
			}
			// [HACK] インターフェースの設定をConfigで反映させておかないとエラーになることへの対応
			// see: https://github.com/sacloud/libsacloud/issues/589
			if err := b.Client.MobileGateway.Config(ctx, zone, id); err != nil {
				return nil, err
			}
			mgw = updated
		}
	}

	mgw, err = b.Client.MobileGateway.Update(ctx, zone, id, &sacloud.MobileGatewayUpdateRequest{
		Name:                            b.Name,
		Description:                     b.Description,
		Tags:                            b.Tags,
		IconID:                          b.IconID,
		InterfaceSettings:               b.getInterfaceSettings(),
		InternetConnectionEnabled:       types.StringFlag(b.InternetConnectionEnabled),
		InterDeviceCommunicationEnabled: types.StringFlag(b.InterDeviceCommunicationEnabled),
		SettingsHash:                    mgw.SettingsHash,
	})
	if err != nil {
		return nil, err
	}

	// traffic config
	trafficConfig, err := b.Client.MobileGateway.GetTrafficConfig(ctx, zone, id)
	if err != nil {
		if !sacloud.IsNotFoundError(err) {
			return nil, err
		}
	}
	if !reflect.DeepEqual(trafficConfig, b.TrafficConfig) {
		if trafficConfig != nil && b.TrafficConfig == nil {
			if err := b.Client.MobileGateway.DeleteTrafficConfig(ctx, zone, id); err != nil {
				return nil, err
			}
		} else {
			if err := b.Client.MobileGateway.SetTrafficConfig(ctx, zone, id, b.TrafficConfig); err != nil {
				return nil, err
			}
		}
	}

	// dns
	dns, err := b.Client.MobileGateway.GetDNS(ctx, zone, id)
	if err != nil {
		if !sacloud.IsNotFoundError(err) {
			return nil, err
		}
	}
	if !reflect.DeepEqual(dns, b.DNS) {
		if dns == nil {
			zone, err := b.Client.Zone.Read(ctx, mgw.ZoneID)
			if err != nil {
				return nil, err
			}
			b.DNS = &sacloud.MobileGatewayDNSSetting{
				DNS1: zone.Region.NameServers[0],
				DNS2: zone.Region.NameServers[1],
			}
		}
		if err := b.Client.MobileGateway.SetDNS(ctx, zone, id, b.DNS); err != nil {
			return nil, err
		}
	}

	// static route(
	if len(b.StaticRoutes) > 0 {
		_, err := b.Client.MobileGateway.UpdateSettings(ctx, zone, id, &sacloud.MobileGatewayUpdateSettingsRequest{
			InterfaceSettings:               b.getInterfaceSettings(),
			StaticRoutes:                    b.StaticRoutes,
			InternetConnectionEnabled:       types.StringFlag(b.InternetConnectionEnabled),
			InterDeviceCommunicationEnabled: types.StringFlag(b.InterDeviceCommunicationEnabled),
			SettingsHash:                    mgw.SettingsHash,
		})
		if err != nil {
			return nil, err
		}
	}

	// SIMs and SIMRoutes
	currentSIMs, err := b.currentConnectedSIMs(ctx, zone, id)
	if err != nil {
		return nil, err
	}
	currentSIMRoutes, err := b.currentSIMRoutes(ctx, zone, id)
	if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(currentSIMs, b.SIMs) || !reflect.DeepEqual(currentSIMRoutes, b.SIMRoutes) {
		if len(currentSIMRoutes) > 0 {
			// SIMルートクリア
			if err := b.Client.MobileGateway.SetSIMRoutes(ctx, zone, id, []*sacloud.MobileGatewaySIMRouteParam{}); err != nil {
				return nil, err
			}
		}
		// SIM変更
		added, updated, deleted := b.changedSIMs(currentSIMs, b.SIMs)
		for _, sim := range deleted {
			if err := b.Client.SIM.ClearIP(ctx, sim.SIMID); err != nil {
				return nil, err
			}
			if err := b.Client.MobileGateway.DeleteSIM(ctx, zone, id, sim.SIMID); err != nil {
				return nil, err
			}
		}
		for _, sim := range updated {
			if err := b.Client.SIM.ClearIP(ctx, sim.SIMID); err != nil {
				return nil, err
			}
			if err := b.Client.SIM.AssignIP(ctx, sim.SIMID, &sacloud.SIMAssignIPRequest{IP: sim.IPAddress}); err != nil {
				return nil, err
			}
		}
		for _, sim := range added {
			if err := b.Client.MobileGateway.AddSIM(ctx, zone, id, &sacloud.MobileGatewayAddSIMRequest{SIMID: sim.SIMID.String()}); err != nil {
				return nil, err
			}
			if err := b.Client.SIM.AssignIP(ctx, sim.SIMID, &sacloud.SIMAssignIPRequest{IP: sim.IPAddress}); err != nil {
				return nil, err
			}
		}
		if len(b.SIMRoutes) > 0 {
			if err := b.Client.MobileGateway.SetSIMRoutes(ctx, zone, id, b.getSIMRouteSettings()); err != nil {
				return nil, err
			}
		}
	}

	if err := b.Client.MobileGateway.Config(ctx, zone, id); err != nil {
		return nil, err
	}

	if isNeedRestart {
		if err := power.BootMobileGateway(ctx, b.Client.MobileGateway, zone, id); err != nil {
			return nil, err
		}
	}

	// refresh
	mgw, err = b.Client.MobileGateway.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}
	return mgw, err
}

func (b *Builder) getInterfaceSettings() []*sacloud.MobileGatewayInterfaceSetting {
	if b.PrivateInterface == nil {
		return nil
	}
	return []*sacloud.MobileGatewayInterfaceSetting{
		{
			Index:          1,
			NetworkMaskLen: b.PrivateInterface.NetworkMaskLen,
			IPAddress:      []string{b.PrivateInterface.IPAddress},
		},
	}
}

func (b *Builder) getSIMRouteSettings() []*sacloud.MobileGatewaySIMRouteParam {
	var results []*sacloud.MobileGatewaySIMRouteParam
	for _, route := range b.SIMRoutes {
		results = append(results, &sacloud.MobileGatewaySIMRouteParam{
			ResourceID: route.SIMID.String(),
			Prefix:     route.Prefix,
		})
	}
	return results
}

func (b *Builder) collectUpdateInfo(mgw *sacloud.MobileGateway) (isNeedShutdown bool, err error) {
	// スイッチの変更/削除は再起動が必要
	isNeedShutdown = b.isPrivateInterfaceChanged(mgw)
	return
}

func (b *Builder) isPrivateInterfaceChanged(mgw *sacloud.MobileGateway) bool {
	current := b.currentPrivateInterfaceState(mgw)
	return !reflect.DeepEqual(current, b.PrivateInterface)
}

func (b *Builder) currentPrivateInterfaceState(mgw *sacloud.MobileGateway) *PrivateInterfaceSetting {
	if len(mgw.Interfaces) > 1 {
		switchID := mgw.Interfaces[1].SwitchID
		var setting *sacloud.MobileGatewayInterfaceSetting
		for _, s := range mgw.InterfaceSettings {
			if s.Index == 1 {
				setting = s
			}
		}
		if setting != nil {
			var ip string
			if len(setting.IPAddress) > 0 {
				ip = setting.IPAddress[0]
			}
			return &PrivateInterfaceSetting{
				SwitchID:       switchID,
				IPAddress:      ip,
				NetworkMaskLen: setting.NetworkMaskLen,
			}
		}
	}
	return nil
}

func (b *Builder) currentConnectedSIMs(ctx context.Context, zone string, id types.ID) ([]*SIMSetting, error) {
	var results []*SIMSetting

	sims, err := b.Client.MobileGateway.ListSIM(ctx, zone, id)
	if err != nil && !sacloud.IsNotFoundError(err) {
		return results, err
	}
	for _, sim := range sims {
		results = append(results, &SIMSetting{
			SIMID:     types.StringID(sim.ResourceID),
			IPAddress: sim.IP,
		})
	}
	return results, nil
}

func (b *Builder) currentSIMRoutes(ctx context.Context, zone string, id types.ID) ([]*sacloud.MobileGatewaySIMRoute, error) {
	return b.Client.MobileGateway.GetSIMRoutes(ctx, zone, id)
}

func (b *Builder) changedSIMs(current []*SIMSetting, desired []*SIMSetting) (added, updated, deleted []*SIMSetting) {
	for _, c := range current {
		isExists := false
		for _, d := range desired {
			if c.SIMID == d.SIMID {
				isExists = true
				if c.IPAddress != d.IPAddress {
					updated = append(updated, d)
				}
			}
		}
		if !isExists {
			deleted = append(deleted, c)
		}
	}
	for _, d := range desired {
		isExists := false
		for _, c := range current {
			if c.SIMID == d.SIMID {
				isExists = true
				continue
			}
		}
		if !isExists {
			added = append(added, d)
		}
	}
	return
}
