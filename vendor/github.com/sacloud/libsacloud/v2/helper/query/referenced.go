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
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// IsPrivateHostReferenced 指定の専有ホストが利用されている場合trueを返す
func IsPrivateHostReferenced(ctx context.Context, caller sacloud.APICaller, zone string, privateHostID types.ID) (bool, error) {
	return checkReferenced(ctx, caller, []string{zone}, privateHostID, []referenceFindFunc{
		findServerByPrivateHostID,
	})
}

// WaitWhilePrivateHostIsReferenced 指定の専有ホストが利用されている間待ち合わせる
func WaitWhilePrivateHostIsReferenced(ctx context.Context, caller sacloud.APICaller, zone string, privateHostID types.ID, option CheckReferencedOption) error {
	return waitWhileReferenced(ctx, option, func() (bool, error) {
		return IsPrivateHostReferenced(ctx, caller, zone, privateHostID)
	})
}

// IsPacketFilterReferenced 指定のパケットフィルタが利用されている場合trueを返す
func IsPacketFilterReferenced(ctx context.Context, caller sacloud.APICaller, zone string, packetFilterID types.ID) (bool, error) {
	return checkReferenced(ctx, caller, []string{zone}, packetFilterID, []referenceFindFunc{
		findServerByPacketFilterID,
	})
}

// WaitWhilePacketFilterIsReferenced 指定のパケットフィルタが利用されている間待ち合わせる
func WaitWhilePacketFilterIsReferenced(ctx context.Context, caller sacloud.APICaller, zone string, packetFilterID types.ID, option CheckReferencedOption) error {
	return waitWhileReferenced(ctx, option, func() (bool, error) {
		return IsPacketFilterReferenced(ctx, caller, zone, packetFilterID)
	})
}

// IsSIMReferenced 指定のSIMが利用されている場合trueを返す
func IsSIMReferenced(ctx context.Context, caller sacloud.APICaller, zones []string, simID types.ID) (bool, error) {
	return checkReferenced(ctx, caller, zones, simID, []referenceFindFunc{
		findMobileGatewayBySIMID,
	})
}

// WaitWhileSIMIsReferenced 指定のSIMが利用されている間待ち合わせる
func WaitWhileSIMIsReferenced(ctx context.Context, caller sacloud.APICaller, zones []string, simID types.ID, option CheckReferencedOption) error {
	return waitWhileReferenced(ctx, option, func() (bool, error) {
		return IsSIMReferenced(ctx, caller, zones, simID)
	})
}

// IsBridgeReferenced 指定のブリッジが利用されている場合trueを返す
func IsBridgeReferenced(ctx context.Context, caller sacloud.APICaller, zones []string, bridgeID types.ID) (bool, error) {
	return checkReferenced(ctx, caller, zones, bridgeID, []referenceFindFunc{
		findSwitchByBridgeID,
	})
}

// WaitWhileBridgeIsReferenced 指定のSIMが利用されている間待ち合わせる
func WaitWhileBridgeIsReferenced(ctx context.Context, caller sacloud.APICaller, zones []string, bridgeID types.ID, option CheckReferencedOption) error {
	return waitWhileReferenced(ctx, option, func() (bool, error) {
		return IsBridgeReferenced(ctx, caller, zones, bridgeID)
	})
}

// IsCDROMReferenced 指定のCD-ROM(ISOイメージ)が利用されている場合trueを返す
func IsCDROMReferenced(ctx context.Context, caller sacloud.APICaller, zone string, cdromID types.ID) (bool, error) {
	return checkReferenced(ctx, caller, []string{zone}, cdromID, []referenceFindFunc{
		findServerByCDROMID,
	})
}

// WaitWhileCDROMIsReferenced 指定のCD-ROM(ISOイメージ)が利用されている間待ち合わせる
func WaitWhileCDROMIsReferenced(ctx context.Context, caller sacloud.APICaller, zone string, cdromID types.ID, option CheckReferencedOption) error {
	return waitWhileReferenced(ctx, option, func() (bool, error) {
		return IsCDROMReferenced(ctx, caller, zone, cdromID)
	})
}

// IsDiskReferenced 指定のディスクが利用されている場合trueを返す
func IsDiskReferenced(ctx context.Context, caller sacloud.APICaller, zone string, diskID types.ID) (bool, error) {
	return checkReferenced(ctx, caller, []string{zone}, diskID, []referenceFindFunc{
		findServerByDiskID,
	})
}

// WaitWhileDiskIsReferenced 指定のディスクが利用されている間待ち合わせる
func WaitWhileDiskIsReferenced(ctx context.Context, caller sacloud.APICaller, zone string, diskID types.ID, option CheckReferencedOption) error {
	return waitWhileReferenced(ctx, option, func() (bool, error) {
		return IsDiskReferenced(ctx, caller, zone, diskID)
	})
}

// IsSwitchReferenced 指定のスイッチが利用されている場合trueを返す
//
// ハイブリッド接続情報が残っている場合にも参照されているものとみなしtrueを返す
func IsSwitchReferenced(ctx context.Context, caller sacloud.APICaller, zone string, switchID types.ID) (bool, error) {
	return checkReferenced(ctx, caller, []string{zone}, switchID, []referenceFindFunc{
		switchHasHybridConnection,
		findServerBySwitchID,
		findLoadBalancerBySwitchID,
		findVPCRouterBySwitchID,
		findDatabaseBySwitchID,
		findNFSBySwitchID,
		findMobileGatewayBySwitchID,
	})
}

// WaitWhileSwitchIsReferenced 指定のディスクが利用されている間待ち合わせる
func WaitWhileSwitchIsReferenced(ctx context.Context, caller sacloud.APICaller, zone string, switchID types.ID, option CheckReferencedOption) error {
	return waitWhileReferenced(ctx, option, func() (bool, error) {
		return IsSwitchReferenced(ctx, caller, zone, switchID)
	})
}

type referenceFindFunc func(context.Context, sacloud.APICaller, string, types.ID) (bool, error)

func checkReferenced(ctx context.Context, caller sacloud.APICaller, zones []string, id types.ID, finder []referenceFindFunc) (bool, error) {
	if len(zones) == 0 {
		zones = sacloud.SakuraCloudZones
	}

	for _, zone := range zones {
		for _, f := range finder {
			exists, err := f(ctx, caller, zone, id)
			if exists || err != nil {
				return exists, err
			}
		}
	}
	return false, nil
}

func switchHasHybridConnection(ctx context.Context, caller sacloud.APICaller, zone string, id types.ID) (bool, error) {
	swOp := sacloud.NewSwitchOp(caller)
	sw, err := swOp.Read(ctx, zone, id)
	if err != nil {
		return false, fmt.Errorf("reading switch is failed: %s", err)
	}
	return !sw.HybridConnectionID.IsEmpty(), nil
}

func findServerBySwitchID(ctx context.Context, caller sacloud.APICaller, zone string, id types.ID) (bool, error) {
	swOp := sacloud.NewSwitchOp(caller)

	searched, err := swOp.GetServers(ctx, zone, id)
	if err != nil {
		return false, fmt.Errorf("finding server is failed: %s", err)
	}
	return searched.Count != 0, nil
}

func findServerByPrivateHostID(ctx context.Context, caller sacloud.APICaller, zone string, id types.ID) (bool, error) {
	serverOp := sacloud.NewServerOp(caller)

	searched, err := serverOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return false, fmt.Errorf("finding Server is failed: %s", err)
	}

	for _, s := range searched.Servers {
		if s.PrivateHostID == id {
			return true, nil
		}
	}
	return false, nil
}

func findServerByPacketFilterID(ctx context.Context, caller sacloud.APICaller, zone string, id types.ID) (bool, error) {
	serverOp := sacloud.NewServerOp(caller)

	searched, err := serverOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return false, fmt.Errorf("finding Server is failed: %s", err)
	}

	for _, s := range searched.Servers {
		for _, iface := range s.Interfaces {
			if iface.PacketFilterID == id {
				return true, nil
			}
		}
	}
	return false, nil
}

func findServerByCDROMID(ctx context.Context, caller sacloud.APICaller, zone string, id types.ID) (bool, error) {
	serverOp := sacloud.NewServerOp(caller)

	searched, err := serverOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return false, fmt.Errorf("finding Server is failed: %s", err)
	}

	for _, server := range searched.Servers {
		if server.CDROMID == id {
			return true, nil
		}
	}
	return false, nil
}

func findServerByDiskID(ctx context.Context, caller sacloud.APICaller, zone string, id types.ID) (bool, error) {
	serverOp := sacloud.NewServerOp(caller)

	searched, err := serverOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return false, fmt.Errorf("finding Server is failed: %s", err)
	}

	for _, server := range searched.Servers {
		for _, disk := range server.Disks {
			if disk.ID == id {
				return true, nil
			}
		}
	}
	return false, nil
}

func findSwitchByBridgeID(ctx context.Context, caller sacloud.APICaller, zone string, id types.ID) (bool, error) {
	swOp := sacloud.NewSwitchOp(caller)

	searched, err := swOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return false, fmt.Errorf("finding Switch is failed: %s", err)
	}

	for _, sw := range searched.Switches {
		if sw.BridgeID == id {
			return true, nil
		}
	}
	return false, nil
}

func findMobileGatewayBySIMID(ctx context.Context, caller sacloud.APICaller, zone string, id types.ID) (bool, error) {
	mgwOp := sacloud.NewMobileGatewayOp(caller)

	searched, err := mgwOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return false, fmt.Errorf("finding MobileGateway is failed: %s", err)
	}

	for _, mgw := range searched.MobileGateways {
		sims, err := mgwOp.ListSIM(ctx, zone, mgw.ID)
		if err != nil {
			if sacloud.IsNotFoundError(err) {
				return false, nil
			}
			return false, fmt.Errorf("finding SIMs is failed: %s", err)
		}
		for _, sim := range sims {
			if sim.ResourceID == id.String() {
				return true, nil
			}
		}
	}
	return false, nil
}

func findVPCRouterBySwitchID(ctx context.Context, caller sacloud.APICaller, zone string, id types.ID) (bool, error) {
	vrOp := sacloud.NewVPCRouterOp(caller)

	searched, err := vrOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return false, fmt.Errorf("finding VPCRouter is failed: %s", err)
	}

	for _, vpcRouter := range searched.VPCRouters {
		for _, iface := range vpcRouter.Interfaces {
			if iface.SwitchScope != types.Scopes.Shared && iface.SwitchID == id {
				return true, nil
			}
		}
	}
	return false, nil
}

func findLoadBalancerBySwitchID(ctx context.Context, caller sacloud.APICaller, zone string, id types.ID) (bool, error) {
	lbOp := sacloud.NewLoadBalancerOp(caller)

	searched, err := lbOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return false, fmt.Errorf("finding LoadBalancer is failed: %s", err)
	}

	for _, lb := range searched.LoadBalancers {
		for _, iface := range lb.Interfaces {
			if iface.SwitchScope != types.Scopes.Shared && iface.SwitchID == id {
				return true, nil
			}
		}
	}
	return false, nil
}

func findDatabaseBySwitchID(ctx context.Context, caller sacloud.APICaller, zone string, id types.ID) (bool, error) {
	dbOp := sacloud.NewDatabaseOp(caller)

	searched, err := dbOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return false, fmt.Errorf("finding Database is failed: %s", err)
	}

	for _, db := range searched.Databases {
		for _, iface := range db.Interfaces {
			if iface.SwitchScope != types.Scopes.Shared && iface.SwitchID == id {
				return true, nil
			}
		}
	}
	return false, nil
}

func findNFSBySwitchID(ctx context.Context, caller sacloud.APICaller, zone string, id types.ID) (bool, error) {
	nfsOp := sacloud.NewNFSOp(caller)

	searched, err := nfsOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return false, fmt.Errorf("finding NFS is failed: %s", err)
	}

	for _, nfs := range searched.NFS {
		for _, iface := range nfs.Interfaces {
			if iface.SwitchScope != types.Scopes.Shared && iface.SwitchID == id {
				return true, nil
			}
		}
	}
	return false, nil
}

func findMobileGatewayBySwitchID(ctx context.Context, caller sacloud.APICaller, zone string, id types.ID) (bool, error) {
	mgwOp := sacloud.NewMobileGatewayOp(caller)

	searched, err := mgwOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return false, fmt.Errorf("finding MobileGateway is failed: %s", err)
	}

	for _, mgw := range searched.MobileGateways {
		for _, iface := range mgw.Interfaces {
			if iface.SwitchScope != types.Scopes.Shared && iface.SwitchID == id {
				return true, nil
			}
		}
	}
	return false, nil
}
