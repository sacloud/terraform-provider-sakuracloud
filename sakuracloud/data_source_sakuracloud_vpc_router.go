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

func dataSourceSakuraCloudVPCRouter() *schema.Resource {
	resourceName := "VPCRouter"
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudVPCRouterRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceSwitchID(resourceName),
			"plan":         schemaDataSourcePlan(resourceName, types.VPCRouterPlanStrings),
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version of the VPC Router",
			},
			"public_network_interface": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of additional network interface setting. This doesn't include primary network interface setting",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch_id": schemaDataSourceSwitchID(resourceName),
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The virtual IP address of the VPC Router. This is only used when `plan` is not `standard`",
						},
						"ip_addresses": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "The list of the IP address assigned to the VPC Router. This will be only one value when `plan` is `standard`, two values otherwise",
						},
						"vrid": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Virtual Router Identifier. This is only used when `plan` is not `standard`",
						},
						"aliases": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "A list of ip alias assigned to the VPC Router. This is only used when `plan` is not `standard`",
						},
					},
				},
			},
			"icon_id":     schemaDataSourceIconID(resourceName),
			"description": schemaDataSourceDescription(resourceName),
			"tags":        schemaDataSourceTags(resourceName),
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public ip address of the VPC Router",
			},
			"public_netmask": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The bit length of the subnet to assign to the public network interface",
			},
			"syslog_host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ip address of the syslog host to which the VPC Router sends logs",
			},
			"internet_connection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag to enable connecting to the Internet from the VPC Router",
			},
			"private_network_interface": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of additional network interface setting. This doesn't include primary network interface setting",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The index of the network interface. This will be between `1`-`7`",
						},
						"switch_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the connected switch",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The virtual IP address assigned to the network interface. This is only used when `plan` is not `standard`",
						},
						"ip_addresses": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "A list of ip address assigned to the network interface. This will be only one value when `plan` is `standard`, two values otherwise",
						},
						"netmask": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The bit length of the subnet assigned to the network interface",
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
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The index of the network interface on which to enable the DHCP service. This will be between `1`-`7`",
						},
						"range_start": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The start value of IP address range to assign to DHCP client",
						},
						"range_stop": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The end value of IP address range to assign to DHCP client",
						},
						"dns_servers": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "A list of IP address of DNS server to assign to DHCP client",
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
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The static IP address to assign to DHCP client",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The source MAC address of static mapping",
						},
					},
				},
			},
			"dns_forwarding": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface_index": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The index of the network interface on which to enable the DNS forwarding service",
						},
						"dns_servers": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "A list of IP address of DNS server to forward to",
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
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The index of the network interface on which to enable filtering. This will be between `0`-`7`",
						},
						"direction": {
							Type:     schema.TypeString,
							Computed: true,
							Description: descf(
								"The direction to apply the firewall. This will be one of [%s]",
								[]string{"send", "receive"},
							),
						},
						"expression": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
										Description: descf(
											"The protocol used for filtering. This will be one of [%s]",
											types.VPCRouterFirewallProtocolStrings,
										),
									},
									"source_network": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A source IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`)",
									},
									"source_port": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A source port number or port range used for filtering (e.g. `1024`, `1024-2048`). This is only used when `protocol` is `tcp` or `udp`",
									},
									"destination_network": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A destination IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`)",
									},
									"destination_port": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A destination port number or port range used for filtering (e.g. `1024`, `1024-2048`). This is only used when `protocol` is `tcp` or `udp`",
									},
									"allow": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "The flag to allow the packet through the filter",
									},
									"logging": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "The flag to enable packet logging when matching the expression",
									},
									"description": schemaDataSourceDescription("expression"),
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
							Type:        schema.TypeString,
							Sensitive:   true,
							Computed:    true,
							Description: "The pre shared secret for L2TP/IPsec",
						},
						"range_start": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The start value of IP address range to assign to L2TP/IPsec client",
						},
						"range_stop": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The end value of IP address range to assign to L2TP/IPsec client",
						},
					},
				},
			},
			"port_forwarding": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of `port_forwarding` blocks as defined below. This represents a `Reverse NAT`",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
							Description: descf(
								"The protocol used for port forwarding. This will be one of [%s]",
								[]string{"tcp", "udp"},
							),
						},
						"public_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The source port number of the port forwarding. This will be a port number on a public network",
						},
						"private_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The destination ip address of the port forwarding",
						},
						"private_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The destination port number of the port forwarding. This will be a port number on a private network",
						},
						"description": schemaDataSourceDescription("port forwarding"),
					},
				},
			},
			"pptp": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range_start": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The start value of IP address range to assign to PPTP client",
						},
						"range_stop": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The end value of IP address range to assign to PPTP client",
						},
					},
				},
			},
			"wire_guard": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address for WireGuard server",
						},
						"public_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the public key of the WireGuard server",
						},
						"peer": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "the of the peer",
									},
									"ip_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address for peer",
									},
									"public_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "the public key of the WireGuard client",
									},
								},
							},
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
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address of the opposing appliance connected to the VPC Router",
						},
						"remote_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the opposing appliance connected to the VPC Router. This is typically set same as value of `peer`",
						},
						"pre_shared_secret": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The pre shared secret for the VPN",
						},
						"routes": {
							Type:        schema.TypeSet,
							Set:         schema.HashString,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "A list of CIDR block of VPN connected networks",
						},
						"local_prefix": {
							Type:        schema.TypeSet,
							Set:         schema.HashString,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "A list of CIDR block of the network under the VPC Router",
						},
					},
				},
			},
			"site_to_site_vpn_parameter": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ike": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lifetime": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"dpd": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"interval": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"timeout": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"esp": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lifetime": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"encryption_algo": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hash_algo": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"static_nat": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of `static_nat` blocks as defined below. This represents a `1:1 NAT`, doing static mapping to both send/receive to/from the Internet. This is only used when `plan` is not `standard`",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public IP address used for the static NAT",
						},
						"private_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The private IP address used for the static NAT",
						},
						"description": schemaDataSourceDescription("static NAT"),
					},
				},
			},
			"static_route": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CIDR block of destination",
						},
						"next_hop": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address of the next hop",
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
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user name used to authenticate remote access",
						},
						"password": {
							Type:        schema.TypeString,
							Sensitive:   true,
							Computed:    true,
							Description: "The password used to authenticate remote access",
						},
					},
				},
			},
			"zone": schemaDataSourceZone(resourceName),
		},
	}
}

func dataSourceSakuraCloudVPCRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	searcher := iaas.NewVPCRouterOp(client)

	findCondition := &iaas.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, zone, findCondition)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud VPCRouter resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.VPCRouters) == 0 {
		return filterNoResultErr()
	}

	targets := res.VPCRouters
	d.SetId(targets[0].ID.String())
	return setVPCRouterResourceData(ctx, d, zone, client, targets[0])
}
