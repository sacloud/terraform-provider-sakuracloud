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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/utils/vpcrouter"
)

const vpcRouterPowerAPILockKey = "sakuracloud_vpc_router.power.%d.lock"

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
			"ipaddress1": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"ipaddress2": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
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
			"interfaces": {
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
						"ipaddresses": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							MinItems: 1,
							MaxItems: 2,
						},
						"nw_mask_len": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(16, 28),
						},
					},
				},
			},
			"dhcp_servers": {
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
			"dhcp_static_mappings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipaddress": {
							Type:     schema.TypeString,
							Required: true,
						},
						"macaddress": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"firewalls": {
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
						"expressions": {
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
			"port_forwardings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
						},
						"global_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"private_address": {
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
						"global_address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPv4Address(),
						},
						"private_address": {
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
			"static_routes": {
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
			"users": {
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
			"global_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudVPCRouterCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	vrOp := sacloud.NewVPCRouterOp(client)

	builder := vpcrouter.Builder{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Tags:                  expandTags(d),
		IconID:                expandSakuraCloudID(d, "icon_id"),
		PlanID:                expandVPCRouterPlanID(d),
		NICSetting:            expandVPCRouterNICSetting(d),
		AdditionalNICSettings: expandVPCRouterAdditionalNICSettings(d),
		RouterSetting:         expandVPCRouterSettings(d),
	}

	if err := builder.Validate(ctx, vrOp, zone); err != nil {
		return fmt.Errorf("validating parameter for SakuraCloud VPCRouter is failed: %s", err)
	}

	vpcRouter, err := builder.Build(ctx, vrOp, zone)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud VPCRouter is failed: %s", err)
	}
	d.SetId(vpcRouter.ID.String())
	return resourceSakuraCloudVPCRouterRead(d, meta)
}

func resourceSakuraCloudVPCRouterRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	vrOp := sacloud.NewVPCRouterOp(client)

	vpcRouter, err := vrOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud VPCRouter: %s", err)
	}

	return setVPCRouterResourceData(ctx, d, client, vpcRouter)
}

func resourceSakuraCloudVPCRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	vrOp := sacloud.NewVPCRouterOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	vpcRouter, err := vrOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud VPCRouter: %s", err)
	}

	builder := vpcrouter.Builder{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Tags:                  expandTags(d),
		IconID:                expandSakuraCloudID(d, "icon_id"),
		PlanID:                expandVPCRouterPlanID(d),
		NICSetting:            expandVPCRouterNICSetting(d),
		AdditionalNICSettings: expandVPCRouterAdditionalNICSettings(d),
		RouterSetting:         expandVPCRouterSettings(d),
	}

	if err := builder.Validate(ctx, vrOp, zone); err != nil {
		return fmt.Errorf("validating parameter for SakuraCloud VPCRouter is failed: %s", err)
	}

	vpcRouter, err = builder.Update(ctx, vrOp, zone, vpcRouter.ID)
	if err != nil {
		return fmt.Errorf("updating SakuraCloud VPCRouter is failed: %s", err)
	}
	return resourceSakuraCloudVPCRouterRead(d, meta)
}

func resourceSakuraCloudVPCRouterDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	vrOp := sacloud.NewVPCRouterOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	vpcRouter, err := vrOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud VPCRouter: %s", err)
	}

	if vpcRouter.InstanceStatus.IsUp() {
		if err := shutdownVPCRouterSync(ctx, client, zone, vpcRouter.ID); err != nil {
			return fmt.Errorf("stopping VPCRouter is failed: %s", err)
		}
	}

	if err := vrOp.Delete(ctx, zone, vpcRouter.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud VPCRouter is failed: %s", err)
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
	d.Set("global_address", flattenVPCRouterGlobalAddress(data))
	d.Set("vip", flattenVPCRouterVIP(data))
	d.Set("ipaddress1", flattenVPCRouterIPAddress1(data))
	d.Set("ipaddress2", flattenVPCRouterIPAddress2(data))
	if err := d.Set("aliases", flattenVPCRouterIPAliases(data)); err != nil {
		return err
	}
	d.Set("vrid", flattenVPCRouterVRID(data))
	d.Set("syslog_host", data.Settings.SyslogHost)
	d.Set("internet_connection", data.Settings.InternetConnectionEnabled.Bool())
	if err := d.Set("interfaces", flattenVPCRouterInterfaces(data)); err != nil {
		return err
	}
	if err := d.Set("dhcp_servers", flattenVPCRouterDHCPServers(data)); err != nil {
		return err
	}
	if err := d.Set("dhcp_static_mappings", flattenVPCRouterDHCPStaticMappings(data)); err != nil {
		return err
	}
	if err := d.Set("firewalls", flattenVPCRouterFirewalls(data)); err != nil {
		return err
	}
	if err := d.Set("l2tp", flattenVPCRouterL2TP(data)); err != nil {
		return err
	}
	if err := d.Set("pptp", flattenVPCRouterPPTP(data)); err != nil {
		return err
	}
	if err := d.Set("port_forwardings", flattenVPCRouterPortForwardings(data)); err != nil {
		return err
	}
	if err := d.Set("site_to_site_vpn", flattenVPCRouterSiteToSite(data)); err != nil {
		return err
	}
	if err := d.Set("static_nat", flattenVPCRouterStaticNAT(data)); err != nil {
		return err
	}
	if err := d.Set("static_routes", flattenVPCRouterStaticRoutes(data)); err != nil {
		return err
	}
	if err := d.Set("users", flattenVPCRouterUsers(data)); err != nil {
		return err
	}
	d.Set("zone", getZone(d, client))
	return nil
}
