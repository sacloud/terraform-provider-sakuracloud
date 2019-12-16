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
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/utils/power"
)

func resourceSakuraCloudVPCRouter() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterCreate,
		Read:   resourceSakuraCloudVPCRouterRead,
		Update: resourceSakuraCloudVPCRouterUpdate,
		Delete: resourceSakuraCloudVPCRouterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"plan": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "standard",
				ValidateFunc: validation.StringInSlice([]string{"standard", "premium", "highspec", "highspec4000"}, false),
			},
			"switch_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"vip": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				MinItems: 2,
				MaxItems: 2,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vrid": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
			},
			"aliases": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MaxItems: 19,
			},
			"syslog_host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internet_connection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"network_interface": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 7,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 7),
						},
						"switch_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateSakuracloudIDType,
						},
						"vip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip_addresses": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							MinItems: 1,
							MaxItems: 2,
						},
						"netmask": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(16, 28),
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
						},
						"range_start": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
						},
						"range_stop": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
						},
						"dns_servers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
							Type:     schema.TypeString,
							Required: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Required: true,
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
						},
						"direction": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"send", "receive"}, false),
						},
						"expression": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "icmp", "ip"}, false),
									},
									"source_network": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"source_port": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"destination_network": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"destination_port": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allow": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"logging": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"description": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringLenBetween(0, 512),
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
						},
						"range_start": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
						},
						"range_stop": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
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
						},
						"public_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"private_ip": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
						},
						"private_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 512),
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
						},
						"range_stop": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
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
							Type:     schema.TypeString,
							Required: true,
						},
						"remote_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"pre_shared_secret": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringLenBetween(0, 40),
						},
						"routes": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"local_prefix": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
						},
						"private_ip": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 512),
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
							Type:     schema.TypeString,
							Required: true,
						},
						"next_hop": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
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
						},
						"password": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringLenBetween(1, 20),
						},
					},
				},
			},
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 512),
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
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
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudVPCRouterCreate(d *schema.ResourceData, meta interface{}) error {
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	builder := expandVPCRouterBuilder(d, client)
	if err := builder.Validate(ctx, zone); err != nil {
		return fmt.Errorf("validating parameter for SakuraCloud VPCRouter is failed: %s", err)
	}

	vpcRouter, err := builder.Build(ctx, zone)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud VPCRouter is failed: %s", err)
	}
	d.SetId(vpcRouter.ID.String())
	return resourceSakuraCloudVPCRouterRead(d, meta)
}

func resourceSakuraCloudVPCRouterRead(d *schema.ResourceData, meta interface{}) error {
	client, zone := getSacloudClient(d, meta)
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
	client, zone := getSacloudClient(d, meta)
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

	vpcRouter, err = builder.Update(ctx, zone, vpcRouter.ID)
	if err != nil {
		return fmt.Errorf("updating SakuraCloud VPCRouter[%s] is failed: %s", vpcRouter.ID, err)
	}
	return resourceSakuraCloudVPCRouterRead(d, meta)
}

func resourceSakuraCloudVPCRouterDelete(d *schema.ResourceData, meta interface{}) error {
	client, zone := getSacloudClient(d, meta)
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
			return fmt.Errorf("stopping VPCRouter[%s] is failed: %s", vpcRouter.ID, err)
		}
	}

	if err := vrOp.Delete(ctx, zone, vpcRouter.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud VPCRouter[%s] is failed: %s", vpcRouter.ID, err)
	}
	return nil
}

func setVPCRouterResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.VPCRouter) error {
	if data.Availability.IsFailed() {
		d.SetId("")
		return fmt.Errorf("got unexpected state: VPCRouter[%d].Availability is failed", data.ID)
	}

	d.Set("name", data.Name)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	d.Set("plan", flattenVPCRouterPlan(data))
	d.Set("switch_id", flattenVPCRouterSwitchID(data))
	d.Set("public_ip", flattenVPCRouterGlobalAddress(data))
	d.Set("vip", flattenVPCRouterVIP(data))
	if err := d.Set("ip_addresses", flattenVPCRouterIPAddresses(data)); err != nil {
		return err
	}
	if err := d.Set("aliases", flattenVPCRouterIPAliases(data)); err != nil {
		return err
	}
	d.Set("vrid", flattenVPCRouterVRID(data))
	d.Set("syslog_host", data.Settings.SyslogHost)
	d.Set("internet_connection", data.Settings.InternetConnectionEnabled.Bool())
	if err := d.Set("network_interface", flattenVPCRouterInterfaces(data)); err != nil {
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
	d.Set("zone", getZone(d, client))
	return nil
}
