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
	"math"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	mobileGatewayUtil "github.com/sacloud/libsacloud/v2/utils/mobilegateway"
)

func resourceSakuraCloudMobileGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudMobileGatewayCreate,
		Read:   resourceSakuraCloudMobileGatewayRead,
		Update: resourceSakuraCloudMobileGatewayUpdate,
		Delete: resourceSakuraCloudMobileGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_interface": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch_id": {
							Type:         schema.TypeString,
							ValidateFunc: validateSakuracloudIDType,
							Required:     true,
						},
						"ipaddress": {
							Type:         schema.TypeString,
							ValidateFunc: validateIPv4Address(),
							Required:     true,
						},
						"nw_mask_len": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(8, 29),
						},
					},
				},
			},
			"public_interface": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipaddress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nw_mask_len": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"internet_connection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"inter_device_communication": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"dns_server1": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIPv4Address(),
			},
			"dns_server2": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIPv4Address(),
			},
			"traffic_control": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quota": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, math.MaxInt32),
						},
						"band_width_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, math.MaxInt32),
						},
						"enable_email": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"enable_slack": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"slack_webhook": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^https://hooks.slack.com/services/\w+/\w+/\w+$`), ""),
						},
						"auto_traffic_shaping": {
							Type:     schema.TypeBool,
							Optional: true,
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
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"sims": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sim_id": {
							Type:         schema.TypeString,
							ValidateFunc: validateSakuracloudIDType,
							Required:     true,
						},
						"ipaddress": {
							Type:         schema.TypeString,
							ValidateFunc: validateIPv4Address(),
							Required:     true,
						},
					},
				},
			},
			"sim_routes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sim_id": {
							Type:         schema.TypeString,
							ValidateFunc: validateSakuracloudIDType,
							Required:     true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Required: true,
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
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}

func resourceSakuraCloudMobileGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)

	builder := expandMobileGatewayBuilder(d, client)
	if err := builder.Validate(ctx, zone); err != nil {
		return fmt.Errorf("validating SakuraCloud MobileGateway is failed: %s", err)
	}

	mgw, err := builder.Build(ctx, zone)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud MobileGateway is failed: %s", err)
	}

	d.SetId(mgw.ID.String())
	return resourceSakuraCloudMobileGatewayRead(d, meta)

}

func resourceSakuraCloudMobileGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	mgwOp := sacloud.NewMobileGatewayOp(client)

	mgw, err := mgwOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud MobileGateway[%s]: %s", d.Id(), err)
	}

	return setMobileGatewayResourceData(ctx, d, client, mgw)
}

func resourceSakuraCloudMobileGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	mgwOp := sacloud.NewMobileGatewayOp(client)

	mgw, err := mgwOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud MobileGateway[%s]: %s", d.Id(), err)
	}

	builder := expandMobileGatewayBuilder(d, client)
	if err := builder.Validate(ctx, zone); err != nil {
		return fmt.Errorf("validating SakuraCloud MobileGateway is failed: %s", err)
	}

	mgw, err = builder.Update(ctx, zone, mgw.ID)
	if err != nil {
		return fmt.Errorf("updating SakuraCloud MobileGateway[%s] is failed: %s", mgw.ID, err)
	}

	return resourceSakuraCloudMobileGatewayRead(d, meta)
}

func resourceSakuraCloudMobileGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	mgwOp := sacloud.NewMobileGatewayOp(client)
	simOp := sacloud.NewSIMOp(client)

	mgw, err := mgwOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud MobileGateway: %s", err)
	}

	if err := mobileGatewayUtil.Delete(ctx, mgwOp, simOp, zone, mgw.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud MobileGateway[%s] is failed: %s", mgw.ID, err)
	}
	return nil
}

func setMobileGatewayResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.MobileGateway) error {
	zone := getZone(d, client)
	mgwOp := sacloud.NewMobileGatewayOp(client)

	if data.Availability.IsFailed() {
		d.SetId("")
		return fmt.Errorf("got unexpected state: MobileGateway[%d].Availability is failed", data.ID)
	}

	// fetch configs
	tc, err := mgwOp.GetTrafficConfig(ctx, zone, data.ID)
	if err != nil && !sacloud.IsNotFoundError(err) {
		return fmt.Errorf("reading TrafficConfig is failed: %s", err)
	}
	resolver, err := mgwOp.GetDNS(ctx, zone, data.ID)
	if err != nil {
		return fmt.Errorf("reading ResolverConfig is failed: %s", err)
	}
	sims, err := mgwOp.ListSIM(ctx, zone, data.ID)
	if err != nil && !sacloud.IsNotFoundError(err) {
		return fmt.Errorf("reading SIMs is failed: %s", err)
	}
	simRoutes, err := mgwOp.GetSIMRoutes(ctx, zone, data.ID)
	if err != nil {
		return fmt.Errorf("reading SIM Routes is failed: %s", err)
	}

	// set data
	if err := d.Set("private_interface", flattenMobileGatewayPrivateNetworks(data)); err != nil {
		return err
	}
	if err := d.Set("public_interface", flattenMobileGatewayPublicNetworks(data)); err != nil {
		return err
	}
	d.Set("internet_connection", data.Settings.InternetConnectionEnabled.Bool())
	d.Set("inter_device_communication", data.Settings.InterDeviceCommunicationEnabled.Bool())

	if err := d.Set("traffic_control", flattenMobileGatewayTrafficConfigs(tc)); err != nil {
		return err
	}
	d.Set("dns_server1", resolver.DNS1)
	d.Set("dns_server2", resolver.DNS2)
	if err := d.Set("static_routes", flattenMobileGatewayStaticRoutes(data.Settings.StaticRoute)); err != nil {
		return err
	}
	d.Set("name", data.Name)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	if err := d.Set("sims", flattenMobileGatewaySIMs(sims)); err != nil {
		return err
	}
	if err := d.Set("sim_routes", flattenMobileGatewaySIMRoutes(simRoutes)); err != nil {
		return err
	}
	d.Set("zone", zone)

	return nil
}
