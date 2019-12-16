// Copyright 2016-2019 terraform-provider-sakuracloud authors
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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/query"
)

func expandNFSDiskPlanID(ctx context.Context, client *APIClient, d resourceValueGettable) (types.ID, error) {
	var planID types.ID
	planName := d.Get("plan").(string)
	switch planName {
	case "hdd":
		planID = types.NFSPlans.HDD
	case "ssd":
		planID = types.NFSPlans.SSD
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
	return &sacloud.NFSCreateRequest{
		SwitchID:       expandSakuraCloudID(d, "switch_id"),
		PlanID:         planID,
		IPAddresses:    []string{d.Get("ip_address").(string)},
		NetworkMaskLen: d.Get("nw_mask_len").(int),
		DefaultRoute:   d.Get("gateway").(string),
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
