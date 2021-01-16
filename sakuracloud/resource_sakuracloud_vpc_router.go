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
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/helper/power"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudVPCRouter() *schema.Resource {
	resourceName := "VPCRouter"
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterCreate,
		Read:   resourceSakuraCloudVPCRouterRead,
		Update: resourceSakuraCloudVPCRouterUpdate,
		Delete: resourceSakuraCloudVPCRouterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"plan": schemaResourcePlan(resourceName, "standard", types.VPCRouterPlanStrings),
			"version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The version of the VPC Router",
				Default:     2,
				ForceNew:    true,
			},
			"public_network_interface": {
				Type:     schema.TypeList,
				Optional: true, // only required when `plan` is not `standard`
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch_id": {
							Type:         schema.TypeString,
							ForceNew:     true,
							Optional:     true,
							ValidateFunc: validateSakuracloudIDType,
							Description:  "The id of the switch to connect. This is only required when when `plan` is not `standard`",
						},
						"vip": {
							Type:        schema.TypeString,
							ForceNew:    true,
							Optional:    true,
							Description: "The virtual IP address of the VPC Router. This is only required when `plan` is not `standard`",
						},
						"ip_addresses": {
							Type:        schema.TypeList,
							ForceNew:    true,
							Optional:    true,
							MinItems:    2,
							MaxItems:    2,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The list of the IP address to assign to the VPC Router. This is required only one value when `plan` is `standard`, two values otherwise",
						},
						"vrid": {
							Type:        schema.TypeInt,
							ForceNew:    true,
							Optional:    true,
							Description: "The Virtual Router Identifier. This is only required when `plan` is not `standard`",
						},
						"aliases": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							MaxItems:    19,
							Description: "A list of ip alias to assign to the VPC Router. This can only be specified if `plan` is not `standard`",
						},
					},
				},
			},
			"syslog_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ip address of the syslog host to which the VPC Router sends logs",
			},
			"internet_connection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The flag to enable connecting to the Internet from the VPC Router",
			},
			"private_network_interface": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    7,
				Description: "A list of additional network interface setting. This doesn't include primary network interface setting",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 7),
							Description:  descf("The index of the network interface. %s", descRange(1, 7)),
						},
						"switch_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateSakuracloudIDType,
							Description:  "The id of the connected switch",
						},
						"vip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The virtual IP address to assign to the network interface. This is only required when `plan` is not `standard`",
						},
						"ip_addresses": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							MinItems:    1,
							MaxItems:    2,
							Description: "A list of ip address to assign to the network interface. This is required only one value when `plan` is `standard`, two values otherwise",
						},
						"netmask": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(16, 28),
							Description:  "The bit length of the subnet to assign to the network interface",
						},
					},
				},
			},
			"dhcp_server": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface_index": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 7),
							Description: descf(
								"The index of the network interface on which to enable the DHCP service. %s",
								descRange(1, 7),
							),
						},
						"range_start": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
							Description:  "The start value of IP address range to assign to DHCP client",
						},
						"range_stop": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
							Description:  "The end value of IP address range to assign to DHCP client",
						},
						"dns_servers": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "A list of IP address of DNS server to assign to DHCP client",
						},
					},
				},
			},
			"dhcp_static_mapping": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The static IP address to assign to DHCP client",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The source MAC address of static mapping",
						},
					},
				},
			},
			"firewall": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface_index": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 7),
							Description: descf(
								"The index of the network interface on which to enable filtering. %s",
								descRange(0, 7),
							),
						},
						"direction": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"send", "receive"}, false),
							Description: descf(
								"The direction to apply the firewall. This must be one of [%s]",
								[]string{"send", "receive"},
							),
						},
						"expression": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(types.VPCRouterFirewallProtocolStrings, false),
										Description: descf(
											"The protocol used for filtering. This must be one of [%s]",
											types.VPCRouterFirewallProtocolStrings,
										),
									},
									"source_network": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A source IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`)",
									},
									"source_port": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A source port number or port range used for filtering (e.g. `1024`, `1024-2048`). This is only used when `protocol` is `tcp` or `udp`",
									},
									"destination_network": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A destination IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`)",
									},
									"destination_port": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A destination port number or port range used for filtering (e.g. `1024`, `1024-2048`). This is only used when `protocol` is `tcp` or `udp`",
									},
									"allow": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "The flag to allow the packet through the filter",
									},
									"logging": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "The flag to enable packet logging when matching the expression",
									},
									"description": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringLenBetween(0, 512),
										Description:  descf("The description of the expression. %s", descLength(0, 512)),
									},
								},
							},
						},
					},
				},
			},
			"l2tp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pre_shared_secret": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringLenBetween(0, 40),
							Description:  "The pre shared secret for L2TP/IPsec",
						},
						"range_start": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
							Description:  "The start value of IP address range to assign to L2TP/IPsec client",
						},
						"range_stop": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
							Description:  "The end value of IP address range to assign to L2TP/IPsec client",
						},
					},
				},
			},
			"port_forwarding": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
							Description: descf(
								"The protocol used for port forwarding. This must be one of [%s]",
								[]string{"tcp", "udp"},
							),
						},
						"public_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
							Description:  "The source port number of the port forwarding. This must be a port number on a public network",
						},
						"private_ip": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
							Description:  "The destination ip address of the port forwarding",
						},
						"private_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
							Description:  "The destination port number of the port forwarding. This will be a port number on a private network",
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 512),
							Description:  descf("The description of the port forwarding. %s", descLength(0, 512)),
						},
					},
				},
			},
			"pptp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range_start": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
							Description:  "The start value of IP address range to assign to PPTP client",
						},
						"range_stop": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
							Description:  "The end value of IP address range to assign to PPTP client",
						},
					},
				},
			},
			"site_to_site_vpn": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"peer": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The IP address of the opposing appliance connected to the VPC Router",
						},
						"remote_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of the opposing appliance connected to the VPC Router. This is typically set same as value of `peer`",
						},
						"pre_shared_secret": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringLenBetween(0, 40),
							Description:  descf("The pre shared secret for the VPN. %s", descLength(0, 40)),
						},
						"routes": {
							Type:        schema.TypeSet,
							Set:         schema.HashString,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "A list of CIDR block of VPN connected networks",
						},
						"local_prefix": {
							Type:        schema.TypeSet,
							Set:         schema.HashString,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "A list of CIDR block of the network under the VPC Router",
						},
					},
				},
			},
			"static_nat": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_ip": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
							Description:  "The public IP address used for the static NAT",
						},
						"private_ip": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
							Description:  "The private IP address used for the static NAT",
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 512),
							Description:  descf("The description of the static nat. %s", descLength(0, 512)),
						},
					},
				},
			},
			"static_route": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CIDR block of destination",
						},
						"next_hop": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
							Description:  "The IP address of the next hop",
						},
					},
				},
			},
			"user": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 100,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 20),
							Description:  "The user name used to authenticate remote access",
						},
						"password": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringLenBetween(1, 20),
							Description:  "The password used to authenticate remote access",
						},
					},
				},
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"zone":        schemaResourceZone(resourceName),
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
		},
	}
}

func resourceSakuraCloudVPCRouterCreate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	builder := expandVPCRouterBuilder(d, client)
	if err := builder.Validate(ctx, zone); err != nil {
		return fmt.Errorf("validating parameter for SakuraCloud VPCRouter is failed: %s", err)
	}

	vpcRouter, err := builder.Build(ctx, zone)
	if vpcRouter != nil {
		d.SetId(vpcRouter.ID.String())
	}
	if err != nil {
		return fmt.Errorf("creating SakuraCloud VPCRouter is failed: %s", err)
	}
	return resourceSakuraCloudVPCRouterRead(d, meta)
}

func resourceSakuraCloudVPCRouterRead(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	vrOp := sacloud.NewVPCRouterOp(client)

	vpcRouter, err := vrOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud VPCRouter[%s]: %s", d.Id(), err)
	}

	return setVPCRouterResourceData(ctx, d, client, vpcRouter)
}

func resourceSakuraCloudVPCRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	vrOp := sacloud.NewVPCRouterOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	vpcRouter, err := vrOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud VPCRouter[%s]: %s", d.Id(), err)
	}

	builder := expandVPCRouterBuilder(d, client)
	if err := builder.Validate(ctx, zone); err != nil {
		return fmt.Errorf("validating parameter for SakuraCloud VPCRouter is failed: %s", err)
	}

	_, err = builder.Update(ctx, zone, vpcRouter.ID)
	if err != nil {
		return fmt.Errorf("updating SakuraCloud VPCRouter[%s] is failed: %s", d.Id(), err)
	}
	return resourceSakuraCloudVPCRouterRead(d, meta)
}

func resourceSakuraCloudVPCRouterDelete(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	vrOp := sacloud.NewVPCRouterOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	vpcRouter, err := vrOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud VPCRouter[%s]: %s", d.Id(), err)
	}

	if vpcRouter.InstanceStatus.IsUp() {
		if err := power.ShutdownVPCRouter(ctx, vrOp, zone, vpcRouter.ID, true); err != nil {
			return fmt.Errorf("stopping VPCRouter[%s] is failed: %s", d.Id(), err)
		}
	}

	if err := vrOp.Delete(ctx, zone, vpcRouter.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud VPCRouter[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setVPCRouterResourceData(_ context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.VPCRouter) error {
	if data.Availability.IsFailed() {
		d.SetId("")
		return fmt.Errorf("got unexpected state: VPCRouter[%d].Availability is failed", data.ID)
	}

	d.Set("name", data.Name)               // nolint
	d.Set("icon_id", data.IconID.String()) // nolint
	d.Set("description", data.Description) // nolint
	if err := d.Set("tags", flattenTags(data.Tags)); err != nil {
		return err
	}
	d.Set("plan", flattenVPCRouterPlan(data))                           // nolint
	d.Set("public_ip", flattenVPCRouterGlobalAddress(data))             // nolint
	d.Set("public_netmask", flattenVPCRouterGlobalNetworkMaskLen(data)) // nolint
	if err := d.Set("public_network_interface", flattenVPCRouterPublicNetworkInterface(data)); err != nil {
		return err
	}

	d.Set("syslog_host", data.Settings.SyslogHost)                               // nolint
	d.Set("internet_connection", data.Settings.InternetConnectionEnabled.Bool()) // nolint
	if err := d.Set("private_network_interface", flattenVPCRouterInterfaces(data)); err != nil {
		return err
	}
	if err := d.Set("dhcp_server", flattenVPCRouterDHCPServers(data)); err != nil {
		return err
	}
	if err := d.Set("dhcp_static_mapping", flattenVPCRouterDHCPStaticMappings(data)); err != nil {
		return err
	}
	if err := d.Set("firewall", flattenVPCRouterFirewalls(data)); err != nil {
		return err
	}
	if err := d.Set("l2tp", flattenVPCRouterL2TP(data)); err != nil {
		return err
	}
	if err := d.Set("pptp", flattenVPCRouterPPTP(data)); err != nil {
		return err
	}
	if err := d.Set("port_forwarding", flattenVPCRouterPortForwardings(data)); err != nil {
		return err
	}
	if err := d.Set("site_to_site_vpn", flattenVPCRouterSiteToSite(data)); err != nil {
		return err
	}
	if err := d.Set("static_nat", flattenVPCRouterStaticNAT(data)); err != nil {
		return err
	}
	if err := d.Set("static_route", flattenVPCRouterStaticRoutes(data)); err != nil {
		return err
	}
	if err := d.Set("user", flattenVPCRouterUsers(data)); err != nil {
		return err
	}
	d.Set("zone", getZone(d, client)) // nolint
	return nil
}
