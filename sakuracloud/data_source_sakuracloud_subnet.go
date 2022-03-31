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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/iaas-api-go"
)

func dataSourceSakuraCloudSubnet() *schema.Resource {
	resourceName := "Subnet"

	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudSubnetRead,

		Schema: map[string]*schema.Schema{
			"internet_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
				Description:      "The id of the switch+router resource that the Subnet belongs",
			},
			"index": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The index of the subnet in assigned to the Switch+Router",
			},

			"netmask": schemaDataSourceNetMask(resourceName),
			"next_hop": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ip address of the next-hop at the Subnet",
			},
			"switch_id": schemaDataSourceSwitchID(resourceName),
			"network_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPv4 network address assigned to the Subnet",
			},
			"min_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Minimum IP address in assigned global addresses to the Subnet",
			},
			"max_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Maximum IP address in assigned global addresses to the Subnet",
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of assigned global address to the Subnet",
			},
			"zone": schemaDataSourceZone(resourceName),
		},
	}
}

func dataSourceSakuraCloudSubnetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	internetOp := iaas.NewInternetOp(client)
	subnetOp := iaas.NewSubnetOp(client)

	internetID := expandSakuraCloudID(d, "internet_id")
	subnetIndex := d.Get("index").(int)

	res, err := internetOp.Read(ctx, zone, internetID)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud Internet[%d]: %s", internetID, err)
	}
	if subnetIndex >= len(res.Switch.Subnets) {
		return diag.Errorf("could not find SakuraCloud Subnet: invalid subneet index: %d", subnetIndex)
	}

	subnetID := res.Switch.Subnets[subnetIndex].ID
	subnet, err := subnetOp.Read(ctx, zone, subnetID)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud Subnet[%d]: %s", subnetID, err)
	}

	d.SetId(subnetID.String())
	return setSubnetResourceData(ctx, d, client, subnet)
}
