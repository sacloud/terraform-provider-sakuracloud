// Copyright 2016-2022 terraform-provider-sakuracloud authors
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

package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/libsacloud/v2/helper/query"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func expandNFSDiskPlanID(ctx context.Context, client *APIClient, d resourceValueGettable) (types.ID, error) {
	var planID types.ID
	planName := d.Get("plan").(string)
	planID, ok := types.NFSPlanIDMap[planName]
	if !ok {
		return types.ID(0), fmt.Errorf("plan is not found: %s", planName)
	}
	size := d.Get("size").(int)

	return query.FindNFSPlanID(ctx, sacloud.NewNoteOp(client), planID, types.ENFSSize(size))
}

func flattenNFSDiskPlan(ctx context.Context, client *APIClient, planID types.ID) (string, int, error) {
	planInfo, err := query.GetNFSPlanInfo(ctx, sacloud.NewNoteOp(client), planID)
	if err != nil {
		return "", 0, err
	}
	var planName string
	size := int(planInfo.Size)

	switch planInfo.DiskPlanID {
	case types.NFSPlans.HDD:
		planName = "hdd"
	case types.NFSPlans.SSD:
		planName = "ssd"
	}

	return planName, size, nil
}

func expandNFSCreateRequest(d *schema.ResourceData, planID types.ID) *sacloud.NFSCreateRequest {
	nic := expandNFSNetworkInterface(d)
	return &sacloud.NFSCreateRequest{
		SwitchID:       nic.switchID,
		PlanID:         planID,
		IPAddresses:    []string{nic.ipAddress},
		NetworkMaskLen: nic.netmask,
		DefaultRoute:   nic.gateway,
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Tags:           expandTags(d),
		IconID:         expandSakuraCloudID(d, "icon_id"),
	}
}

func expandNFSUpdateRequest(d *schema.ResourceData) *sacloud.NFSUpdateRequest {
	return &sacloud.NFSUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),
	}
}

type nfsNetworkInterface struct {
	switchID  types.ID
	ipAddress string
	netmask   int
	gateway   string
}

func expandNFSNetworkInterface(d resourceValueGettable) *nfsNetworkInterface {
	d = mapFromFirstElement(d, "network_interface")
	if d == nil {
		return nil
	}
	return &nfsNetworkInterface{
		switchID:  expandSakuraCloudID(d, "switch_id"),
		ipAddress: stringOrDefault(d, "ip_address"),
		netmask:   intOrDefault(d, "netmask"),
		gateway:   stringOrDefault(d, "gateway"),
	}
}

func flattenNFSNetworkInterface(nfs *sacloud.NFS) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"switch_id":  nfs.SwitchID.String(),
			"ip_address": nfs.IPAddresses[0],
			"netmask":    nfs.NetworkMaskLen,
			"gateway":    nfs.DefaultRoute,
		},
	}
}
