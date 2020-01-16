// Copyright 2016-2020 terraform-provider-sakuracloud authors
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
)

func dataSourceSakuraCloudSubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudSubnetRead,

		Schema: map[string]*schema.Schema{
			"internet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"index": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},

			"nw_mask_len": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"next_hop": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipaddresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		},
	}
}

func dataSourceSakuraCloudSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	internetID := toSakuraCloudID(d.Get("internet_id").(string))
	subnetIndex := d.Get("index").(int)

	res, err := client.Internet.Read(internetID)
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Internet resource(id:%d): %s", internetID, err)
	}
	if subnetIndex >= len(res.Switch.Subnets) {
		return fmt.Errorf("Couldn't find SakuraCloud Subnet: invalid subneet index: %d", subnetIndex)
	}

	subnetID := res.Switch.Subnets[subnetIndex].ID
	subnet, err := client.Subnet.Read(subnetID)
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Subnet(id:%d) resource: %s", subnetID, err)
	}

	d.SetId(subnet.GetStrID())
	return setSubnetResourceData(d, client, subnet)
}
