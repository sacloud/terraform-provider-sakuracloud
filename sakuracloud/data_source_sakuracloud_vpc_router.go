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

func dataSourceSakuraCloudVPCRouter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudVPCRouterRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vrid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"aliases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"icon_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"syslog_host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_connection": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"network_interface": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"switch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_addresses": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"nw_mask_len": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"dhcp_server": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface_index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"range_start": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"range_stop": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_servers": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
					},
				},
			},
			"dhcp_static_mapping": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"macaddress": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"firewall": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface_index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expression": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_nw": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dest_nw": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dest_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"allow": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"logging": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"l2tp": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pre_shared_secret": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
						"range_start": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"range_stop": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"port_forwarding": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"pptp": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range_start": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"range_stop": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"site_to_site_vpn": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"peer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pre_shared_secret": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"routes": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"local_prefix": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
					},
				},
			},
			"static_nat": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"static_route": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"user": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
					},
				},
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

func dataSourceSakuraCloudVPCRouterRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	searcher := sacloud.NewVPCRouterOp(client)

	findCondition := &sacloud.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, zone, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud VPCRouter resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.VPCRouters) == 0 {
		return filterNoResultErr()
	}

	targets := res.VPCRouters
	d.SetId(targets[0].ID.String())
	return setVPCRouterResourceData(ctx, d, client, targets[0])
}
