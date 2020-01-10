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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func dataSourceSakuraCloudSubnet() *schema.Resource {
	resourceName := "subnet"

	return &schema.Resource{
		Read: dataSourceSakuraCloudSubnetRead,

		Schema: map[string]*schema.Schema{
			"internet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
				Description:  "The id of the switch+router resource that the subnet belongs",
			},
			"index": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "The index of the subnet in assigned to the switch+router",
			},

			"netmask": schemaDataSourceNetMask(resourceName),
			"next_hop": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ip address of the next-hop at the subnet",
			},
			"switch_id": schemaDataSourceSwitchID(resourceName),
			"network_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPv4 network address assigned to the subnet",
			},
			"min_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Minimum IP address in assigned global addresses to the subnet",
			},
			"max_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Maximum IP address in assigned global addresses to the subnet",
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of assigned global address to the subnet",
			},
			"zone": schemaDataSourceZone(resourceName),
		},
	}
}

func dataSourceSakuraCloudSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	internetOp := sacloud.NewInternetOp(client)
	subnetOp := sacloud.NewSubnetOp(client)

	internetID := expandSakuraCloudID(d, "internet_id")
	subnetIndex := d.Get("index").(int)

	res, err := internetOp.Read(ctx, zone, internetID)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud Internet[%d]: %s", internetID, err)
	}
	if subnetIndex >= len(res.Switch.Subnets) {
		return fmt.Errorf("could not find SakuraCloud Subnet: invalid subneet index: %d", subnetIndex)
	}

	subnetID := res.Switch.Subnets[subnetIndex].ID
	subnet, err := subnetOp.Read(ctx, zone, subnetID)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud Subnet[%d]: %s", subnetID, err)
	}

	d.SetId(subnetID.String())
	return setSubnetResourceData(ctx, d, client, subnet)
}
