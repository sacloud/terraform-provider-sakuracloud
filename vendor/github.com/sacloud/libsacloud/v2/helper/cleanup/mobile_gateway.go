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

package cleanup

import (
	"context"

	"github.com/sacloud/libsacloud/v2/helper/power"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// DeleteMobileGateway 削除
func DeleteMobileGateway(ctx context.Context, mgwAPI sacloud.MobileGatewayAPI, simAPI sacloud.SIMAPI, zone string, id types.ID) error {
	// check MobileGateway is exists
	mgw, err := mgwAPI.Read(ctx, zone, id)
	if err != nil {
		return err
	}

	if mgw.InstanceStatus.IsUp() {
		if err := power.ShutdownMobileGateway(ctx, mgwAPI, zone, id, true); err != nil {
			return err
		}
	}

	// sim route
	simRoutes, err := mgwAPI.GetSIMRoutes(ctx, zone, id)
	if err != nil {
		return err
	}
	if len(simRoutes) > 0 {
		if err := mgwAPI.SetSIMRoutes(ctx, zone, id, []*sacloud.MobileGatewaySIMRouteParam{}); err != nil {
			return err
		}
	}

	// sim
	sims, err := mgwAPI.ListSIM(ctx, zone, id)
	if err != nil {
		return err
	}
	for _, sim := range sims {
		if err := simAPI.ClearIP(ctx, types.StringID(sim.ResourceID)); err != nil {
			return err
		}
		if err := mgwAPI.DeleteSIM(ctx, zone, id, types.StringID(sim.ResourceID)); err != nil {
			return err
		}
	}

	return mgwAPI.Delete(ctx, zone, id)
}
