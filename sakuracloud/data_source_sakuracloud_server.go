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
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
)

func dataSourceSakuraCloudServer() *schema.Resource {
	resourceName := "Server"

	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudServerRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"core": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of virtual CPUs",
			},
			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of memory in GiB",
			},
			"gpu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of GPUs",
			},
			"commitment": {
				Type:     schema.TypeString,
				Computed: true,
				Description: descf(
					"The policy of how to allocate virtual CPUs to the server. This will be one of [%s]",
					types.CommitmentStrings,
				),
			},
			"disks": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of disk id connected to the server",
			},
			"interface_driver": {
				Type:     schema.TypeString,
				Computed: true,
				Description: descf(
					"The driver name of network interface. This will be one of [%s]",
					types.InterfaceDriverStrings,
				),
			},
			"network_interface": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"upstream": {
							Type:     schema.TypeString,
							Computed: true,
							Description: descf(
								"The upstream type or upstream switch id. This will be one of [%s]",
								[]string{"shared", "disconnect", "<switch id>"},
							),
						},
						"user_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address for only display. This value doesn't affect actual NIC settings",
						},
						"packet_filter_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the packet filter attached to the network interface",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MAC address",
						},
					},
				},
			},
			"cdrom_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the CD-ROM attached to the server",
			},
			"private_host_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the private host which the server is assigned",
			},
			"private_host_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the private host which the server is assigned",
			},
			"icon_id":     schemaDataSourceIconID(resourceName),
			"description": schemaDataSourceDescription(resourceName),
			"tags":        schemaDataSourceTags(resourceName),
			"zone":        schemaDataSourceZone(resourceName),
			"ip_address":  schemaDataSourceIPAddress(resourceName),
			"gateway":     schemaDataSourceGateway(resourceName),
			"netmask":     schemaDataSourceNetMask(resourceName),
			"network_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network address which the `ip_address` belongs",
			},
			"hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The hostname of the Server",
			},
			"dns_servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of IP address of DNS server in the zone",
			},
		},
	}
}

func dataSourceSakuraCloudServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	searcher := iaas.NewServerOp(client)

	findCondition := &iaas.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, zone, findCondition)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud Server resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.Servers) == 0 {
		return filterNoResultErr()
	}

	targets := res.Servers
	d.SetId(targets[0].ID.String())
	return setServerResourceData(ctx, d, client, targets[0])
}
