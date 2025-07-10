// Copyright 2016-2025 terraform-provider-sakuracloud authors
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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/power"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudVPCRouter() *schema.Resource {
	resourceName := "VPCRouter"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudVPCRouterCreate,
		ReadContext:   resourceSakuraCloudVPCRouterRead,
		UpdateContext: resourceSakuraCloudVPCRouterUpdate,
		DeleteContext: resourceSakuraCloudVPCRouterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": func() *schema.Schema { s := schemaResourceName(resourceName); s.Sensitive = true; return s }(),
			"plan": func() *schema.Schema {
				s := schemaResourcePlan(resourceName, "standard", types.VPCRouterPlanStrings)
				s.Sensitive = true
				return s
			}(),
			"version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The version of the VPC Router",
				Default:     2,
				ForceNew:    true,
				Sensitive:   true,
			},
			"public_network_interface": {
				Type:      schema.TypeList,
				Optional:  true, // only required when `plan` is not `standard`
				MinItems:  1,
				MaxItems:  1,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch_id": {
							Type:             schema.TypeString,
							ForceNew:         true,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
							Description:      "The id of the switch to connect. This is only required when when `plan` is not `standard`",
							Sensitive:        true,
						},
						"vip": {
							Type:        schema.TypeString,
							ForceNew:    true,
							Optional:    true,
							Description: "The virtual IP address of the VPC Router. This is only required when `plan` is not `standard`",
							Sensitive:   true,
						},
						"ip_addresses": {
							Type:        schema.TypeList,
							ForceNew:    true,
							Optional:    true,
							MinItems:    2,
							MaxItems:    2,
							Elem:        &schema.Schema{Type: schema.TypeString, Sensitive: true},
							Description: "The list of the IP address to assign to the VPC Router. This is required only one value when `plan` is `standard`, two values otherwise",
							Sensitive:   true,
						},
						"vrid": {
							Type:        schema.TypeInt,
							ForceNew:    true,
							Optional:    true,
							Description: "The Virtual Router Identifier. This is only required when `plan` is not `standard`",
							Sensitive:   true,
						},
						"aliases": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString, Sensitive: true},
							MaxItems:    19,
							Description: "A list of ip alias to assign to the VPC Router. This can only be specified if `plan` is not `standard`",
							Sensitive:   true,
						},
					},
				},
			},
			"syslog_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ip address of the syslog host to which the VPC Router sends logs",
				Sensitive:   true,
			},
			"internet_connection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The flag to enable connecting to the Internet from the VPC Router",
				Sensitive:   true,
			},
			"private_network_interface": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    7,
				Description: "A list of additional network interface setting. This doesn't include primary network interface setting",
				Sensitive:   true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 7)),
							Description:      desc.Sprintf("The index of the network interface. %s", desc.Range(1, 7)),
							Sensitive:        true,
						},
						"switch_id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
							Description:      "The id of the connected switch",
							Sensitive:        true,
						},
						"vip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The virtual IP address to assign to the network interface. This is only required when `plan` is not `standard`",
							Sensitive:   true,
						},
						"ip_addresses": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString, Sensitive: true},
							MinItems:    1,
							MaxItems:    2,
							Description: "A list of ip address to assign to the network interface. This is required only one value when `plan` is `standard`, two values otherwise",
							Sensitive:   true,
						},
						"netmask": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(16, 29)),
							Description:      "The bit length of the subnet to assign to the network interface",
							Sensitive:        true,
						},
					},
				},
			},
			"dhcp_server": {
				Type:      schema.TypeList,
				Optional:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface_index": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 7)),
							Description: desc.Sprintf(
								"The index of the network interface on which to enable the DHCP service. %s",
								desc.Range(1, 7),
							),
							Sensitive: true,
						},
						"range_start": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateIPv4Address(),
							Description:      "The start value of IP address range to assign to DHCP client",
							Sensitive:        true,
						},
						"range_stop": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateIPv4Address(),
							Description:      "The end value of IP address range to assign to DHCP client",
							Sensitive:        true,
						},
						"dns_servers": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString, Sensitive: true},
							Description: "A list of IP address of DNS server to assign to DHCP client",
							Sensitive:   true,
						},
					},
				},
			},
			"dhcp_static_mapping": {
				Type:      schema.TypeList,
				Optional:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The static IP address to assign to DHCP client",
							Sensitive:   true,
						},
						"mac_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The source MAC address of static mapping",
							Sensitive:   true,
						},
					},
				},
			},
			"dns_forwarding": {
				Type:      schema.TypeList,
				Optional:  true,
				MaxItems:  1,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface_index": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 7)),
							Description: desc.Sprintf(
								"The index of the network interface on which to enable the DNS forwarding service. %s",
								desc.Range(1, 7),
							),
							Sensitive: true,
						},
						"dns_servers": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString, Sensitive: true},
							Description: "A list of IP address of DNS server to forward to",
							Sensitive:   true,
						},
					},
				},
			},
			"firewall": {
				Type:      schema.TypeList,
				Optional:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface_index": {
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 7)),
							Description: desc.Sprintf(
								"The index of the network interface on which to enable filtering. %s",
								desc.Range(0, 7),
							),
							Sensitive: true,
						},
						"direction": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"send", "receive"}, false)),
							Description: desc.Sprintf(
								"The direction to apply the firewall. This must be one of [%s]",
								[]string{"send", "receive"},
							),
							Sensitive: true,
						},
						"expression": {
							Type:      schema.TypeList,
							Required:  true,
							Sensitive: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.VPCRouterFirewallProtocolStrings, false)),
										Description: desc.Sprintf(
											"The protocol used for filtering. This must be one of [%s]",
											types.VPCRouterFirewallProtocolStrings,
										),
										Sensitive: true,
									},
									"source_network": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A source IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`)",
										Sensitive:   true,
									},
									"source_port": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A source port number or port range used for filtering (e.g. `1024`, `1024-2048`). This is only used when `protocol` is `tcp` or `udp`",
										Sensitive:   true,
									},
									"destination_network": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A destination IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`)",
										Sensitive:   true,
									},
									"destination_port": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A destination port number or port range used for filtering (e.g. `1024`, `1024-2048`). This is only used when `protocol` is `tcp` or `udp`",
										Sensitive:   true,
									},
									"allow": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "The flag to allow the packet through the filter",
										Sensitive:   true,
									},
									"logging": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "The flag to enable packet logging when matching the expression",
										Sensitive:   true,
									},
									"description": {
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: isValidLengthBetween(0, 512),
										Description:      desc.Sprintf("The description of the expression. %s", desc.Length(0, 512)),
										Sensitive:        true,
									},
								},
							},
						},
					},
				},
			},
			"l2tp": {
				Type:      schema.TypeList,
				Optional:  true,
				MaxItems:  1,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pre_shared_secret": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: isValidLengthBetween(0, 40),
							Description:      "The pre shared secret for L2TP/IPsec",
						},
						"range_start": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateIPv4Address(),
							Description:      "The start value of IP address range to assign to L2TP/IPsec client",
							Sensitive:        true,
						},
						"range_stop": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateIPv4Address(),
							Description:      "The end value of IP address range to assign to L2TP/IPsec client",
							Sensitive:        true,
						},
					},
				},
			},
			"port_forwarding": {
				Type:      schema.TypeList,
				Optional:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"tcp", "udp"}, false)),
							Description: desc.Sprintf(
								"The protocol used for port forwarding. This must be one of [%s]",
								[]string{"tcp", "udp"},
							),
							Sensitive: true,
						},
						"public_port": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 65535)),
							Description:      "The source port number of the port forwarding. This must be a port number on a public network",
							Sensitive:        true,
						},
						"private_ip": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateIPv4Address(),
							Description:      "The destination ip address of the port forwarding",
							Sensitive:        true,
						},
						"private_port": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 65535)),
							Description:      "The destination port number of the port forwarding. This will be a port number on a private network",
							Sensitive:        true,
						},
						"description": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: isValidLengthBetween(0, 512),
							Description:      desc.Sprintf("The description of the port forwarding. %s", desc.Length(0, 512)),
							Sensitive:        true,
						},
					},
				},
			},
			"pptp": {
				Type:      schema.TypeList,
				Optional:  true,
				MaxItems:  1,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range_start": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateIPv4Address(),
							Description:      "The start value of IP address range to assign to PPTP client",
							Sensitive:        true,
						},
						"range_stop": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateIPv4Address(),
							Description:      "The end value of IP address range to assign to PPTP client",
							Sensitive:        true,
						},
					},
				},
			},
			"wire_guard": {
				Type:      schema.TypeList,
				Optional:  true,
				MaxItems:  1,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The IP address for WireGuard server. This must be formatted with xxx.xxx.xxx.xxx/nn",
							Sensitive:   true,
						},
						"public_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the public key of the WireGuard server",
							Sensitive:   true,
						},
						"peer": {
							Type:      schema.TypeList,
							Optional:  true,
							Sensitive: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "the of the peer",
										Sensitive:   true,
									},
									"ip_address": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The IP address for peer",
										Sensitive:   true,
									},
									"public_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "the public key of the WireGuard client",
										Sensitive:   true,
									},
								},
							},
						},
					},
				},
			},
			"site_to_site_vpn": {
				Type:      schema.TypeList,
				Optional:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"peer": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The IP address of the opposing appliance connected to the VPC Router",
							Sensitive:   true,
						},
						"remote_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of the opposing appliance connected to the VPC Router. This is typically set same as value of `peer`",
							Sensitive:   true,
						},
						"pre_shared_secret": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: isValidLengthBetween(0, 40),
							Description:      desc.Sprintf("The pre shared secret for the VPN. %s", desc.Length(0, 40)),
						},
						"routes": {
							Type:        schema.TypeSet,
							Set:         schema.HashString,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString, Sensitive: true},
							Description: "A list of CIDR block of VPN connected networks",
							Sensitive:   true,
						},
						"local_prefix": {
							Type:        schema.TypeSet,
							Set:         schema.HashString,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString, Sensitive: true},
							Description: "A list of CIDR block of the network under the VPC Router",
							Sensitive:   true,
						},
					},
				},
			},
			"site_to_site_vpn_parameter": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ike": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lifetime": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Lifetime of IKE SA. Default: 28800",
									},
									"dpd": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"interval": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Default: 15",
												},
												"timeout": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Default: 30",
												},
											},
										},
									},
								},
							},
						},
						"esp": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lifetime": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Default: 1800",
									},
								},
							},
						},
						"encryption_algo": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.VPCRouterSiteToSiteVPNEncryptionAlgos, false)),
							Description: desc.Sprintf(
								"This must be one of [%s]",
								types.VPCRouterSiteToSiteVPNEncryptionAlgos,
							),
						},
						"hash_algo": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.VPCRouterSiteToSiteVPNHashAlgos, false)),
							Description: desc.Sprintf(
								"This must be one of [%s]",
								types.VPCRouterSiteToSiteVPNHashAlgos,
							),
						},
						"dh_group": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.VPCRouterSiteToSiteVPNDHGroups, false)),
							Description: desc.Sprintf(
								"This must be one of [%s]",
								types.VPCRouterSiteToSiteVPNDHGroups,
							),
						},
					},
				},
			},
			"static_nat": {
				Type:      schema.TypeList,
				Optional:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_ip": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateIPv4Address(),
							Description:      "The public IP address used for the static NAT",
							Sensitive:        true,
						},
						"private_ip": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateIPv4Address(),
							Description:      "The private IP address used for the static NAT",
							Sensitive:        true,
						},
						"description": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: isValidLengthBetween(0, 512),
							Description:      desc.Sprintf("The description of the static nat. %s", desc.Length(0, 512)),
							Sensitive:        true,
						},
					},
				},
			},
			"static_route": {
				Type:      schema.TypeList,
				Optional:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CIDR block of destination",
							Sensitive:   true,
						},
						"next_hop": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateIPv4Address(),
							Description:      "The IP address of the next hop",
							Sensitive:        true,
						},
					},
				},
			},
			"scheduled_maintenance": {
				Type:      schema.TypeList,
				Optional:  true,
				Computed:  true,
				MaxItems:  1,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day_of_week": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  types.DaysOfTheWeek.Monday.String(),
							Description: desc.Sprintf(
								"The value must be in [%s]",
								types.DaysOfTheWeekStrings,
							),
							Sensitive: true,
						},
						"hour": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          3,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 23)),
							Description:      "The time to start maintenance",
							Sensitive:        true,
						},
					},
				},
			},
			"user": {
				Type:      schema.TypeList,
				Optional:  true,
				MaxItems:  100,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: isValidLengthBetween(1, 20),
							Description:      "The user name used to authenticate remote access",
							Sensitive:        true,
						},
						"password": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: isValidLengthBetween(1, 20),
							Description:      "The password used to authenticate remote access",
						},
					},
				},
			},
			"icon_id":     func() *schema.Schema { s := schemaResourceIconID(resourceName); s.Sensitive = true; return s }(),
			"description": func() *schema.Schema { s := schemaResourceDescription(resourceName); s.Sensitive = true; return s }(),
			"tags":        func() *schema.Schema { s := schemaResourceTags(resourceName); s.Sensitive = true; return s }(),
			"zone":        func() *schema.Schema { s := schemaResourceZone(resourceName); s.Sensitive = true; return s }(),
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public ip address of the VPC Router",
				Sensitive:   true,
			},
			"public_netmask": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The bit length of the subnet to assign to the public network interface",
				Sensitive:   true,
			},
		},
	}
}

func resourceSakuraCloudVPCRouterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	builder := expandVPCRouterBuilder(d, client, zone)
	if err := builder.Validate(ctx, zone); err != nil {
		return diag.Errorf("validating parameter for SakuraCloud VPCRouter is failed: %s", err)
	}

	vpcRouter, err := builder.Build(ctx)
	if vpcRouter != nil {
		d.SetId(vpcRouter.ID.String())
	}
	if err != nil {
		return diag.Errorf("creating SakuraCloud VPCRouter is failed: %s", err)
	}

	// Note: 起動してからしばらくは/:id/Statusが空となるため、数秒待つようにする。
	time.Sleep(client.vpcRouterWaitAfterCreateDuration)

	return resourceSakuraCloudVPCRouterRead(ctx, d, meta)
}

func resourceSakuraCloudVPCRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	vrOp := iaas.NewVPCRouterOp(client)

	vpcRouter, err := vrOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud VPCRouter[%s]: %s", d.Id(), err)
	}

	return setVPCRouterResourceData(ctx, d, zone, client, vpcRouter)
}

func resourceSakuraCloudVPCRouterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	vrOp := iaas.NewVPCRouterOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	vpcRouter, err := vrOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud VPCRouter[%s]: %s", d.Id(), err)
	}

	builder := expandVPCRouterBuilder(d, client, zone)
	if err := builder.Validate(ctx, zone); err != nil {
		return diag.Errorf("validating parameter for SakuraCloud VPCRouter is failed: %s", err)
	}
	builder.ID = vpcRouter.ID

	_, err = builder.Build(ctx)
	if err != nil {
		return diag.Errorf("updating SakuraCloud VPCRouter[%s] is failed: %s", d.Id(), err)
	}

	// Note: 起動してからしばらくは/:id/Statusが空となるため、数秒待つようにする。
	time.Sleep(client.vpcRouterWaitAfterCreateDuration)

	return resourceSakuraCloudVPCRouterRead(ctx, d, meta)
}

func resourceSakuraCloudVPCRouterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	vrOp := iaas.NewVPCRouterOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	vpcRouter, err := vrOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud VPCRouter[%s]: %s", d.Id(), err)
	}

	if vpcRouter.InstanceStatus.IsUp() {
		if err := power.ShutdownVPCRouter(ctx, vrOp, zone, vpcRouter.ID, true); err != nil {
			return diag.Errorf("stopping VPCRouter[%s] is failed: %s", d.Id(), err)
		}
	}

	if err := vrOp.Delete(ctx, zone, vpcRouter.ID); err != nil {
		return diag.Errorf("deleting SakuraCloud VPCRouter[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setVPCRouterResourceData(ctx context.Context, d *schema.ResourceData, zone string, client *APIClient, data *iaas.VPCRouter) diag.Diagnostics {
	if data.Availability.IsFailed() {
		d.SetId("")
		return diag.Errorf("got unexpected state: VPCRouter[%d].Availability is failed", data.ID)
	}

	d.Set("name", data.Name)               //nolint
	d.Set("icon_id", data.IconID.String()) //nolint
	d.Set("description", data.Description) //nolint
	if err := d.Set("tags", flattenTags(data.Tags)); err != nil {
		return diag.FromErr(err)
	}
	d.Set("plan", flattenVPCRouterPlan(data))                           //nolint
	d.Set("public_ip", flattenVPCRouterGlobalAddress(data))             //nolint
	d.Set("public_netmask", flattenVPCRouterGlobalNetworkMaskLen(data)) //nolint
	if err := d.Set("public_network_interface", flattenVPCRouterPublicNetworkInterface(data)); err != nil {
		return diag.FromErr(err)
	}

	d.Set("syslog_host", data.Settings.SyslogHost)                               //nolint
	d.Set("internet_connection", data.Settings.InternetConnectionEnabled.Bool()) //nolint
	if err := d.Set("private_network_interface", flattenVPCRouterInterfaces(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dhcp_server", flattenVPCRouterDHCPServers(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dhcp_static_mapping", flattenVPCRouterDHCPStaticMappings(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dns_forwarding", flattenVPCRouterDNSForwarding(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("firewall", flattenVPCRouterFirewalls(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("l2tp", flattenVPCRouterL2TP(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("pptp", flattenVPCRouterPPTP(data)); err != nil {
		return diag.FromErr(err)
	}
	// get public key from /:id/Status API
	status, err := iaas.NewVPCRouterOp(client).Status(ctx, zone, data.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	wireGuardPublicKey := ""
	if status != nil && status.WireGuard != nil {
		wireGuardPublicKey = status.WireGuard.PublicKey
	}
	if err := d.Set("wire_guard", flattenVPCRouterWireGuard(data, wireGuardPublicKey)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("port_forwarding", flattenVPCRouterPortForwardings(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("site_to_site_vpn", flattenVPCRouterSiteToSiteConfig(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("site_to_site_vpn_parameter", flattenVPCRouterSiteToSiteParameter(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("static_nat", flattenVPCRouterStaticNAT(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("static_route", flattenVPCRouterStaticRoutes(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("user", flattenVPCRouterUsers(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("scheduled_maintenance", flattenVPCRouterScheduledMaintenance(data)); err != nil {
		return diag.FromErr(err)
	}
	d.Set("zone", getZone(d, client)) //nolint
	return nil
}
