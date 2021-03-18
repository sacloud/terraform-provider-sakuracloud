// Copyright 2016-2021 terraform-provider-sakuracloud authors
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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func dataSourceSakuraCloudInternet() *schema.Resource {
	resourceName := "Switch+Router"
	return &schema.Resource{
		Read: dataSourceSakuraCloudInternetRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"icon_id":      schemaDataSourceIconID(resourceName),
			"description":  schemaDataSourceDescription(resourceName),
			"tags":         schemaDataSourceTags(resourceName),
			"netmask":      schemaDataSourceNetMask(resourceName),
			"band_width": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The bandwidth of the network connected to the Internet in Mbps",
			},
			"switch_id": schemaDataSourceSwitchID(resourceName),
			"server_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: descf("A list of the ID of Servers connected to the %s", resourceName),
			},
			"network_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("The IPv4 network address assigned to the %s", resourceName),
			},
			"gateway": schemaDataSourceGateway(resourceName),
			"min_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("Minimum IP address in assigned global addresses to the %s", resourceName),
			},
			"max_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("Maximum IP address in assigned global addresses to the %s", resourceName),
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: descf("A list of assigned global address to the %s", resourceName),
			},
			"enable_ipv6": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag to enable IPv6",
			},
			"ipv6_prefix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("The network prefix of assigned IPv6 addresses to the %s", resourceName),
			},
			"ipv6_prefix_len": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The bit length of IPv6 network prefix",
			},
			"ipv6_network_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("The IPv6 network address assigned to the %s", resourceName),
			},
			"zone": schemaDataSourceZone(resourceName),
		},
	}
}

func dataSourceSakuraCloudInternetRead(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	searcher := sacloud.NewInternetOp(client)

	findCondition := &sacloud.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, zone, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud Internet resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.Internet) == 0 {
		return filterNoResultErr()
	}

	targets := res.Internet
	d.SetId(targets[0].ID.String())
	return setInternetResourceData(ctx, d, client, targets[0])
}
