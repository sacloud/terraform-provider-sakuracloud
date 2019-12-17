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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/utils/cleanup"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_network_interface": {
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
						"ip_address": {
							Type:         schema.TypeString,
							ValidateFunc: validateIPv4Address(),
							Required:     true,
						},
						"netmask": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(8, 29),
						},
					},
				},
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_netmask": {
				Type:     schema.TypeInt,
				Computed: true,
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
			"dns_servers": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 2,
				MinItems: 2,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"sim": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sim_id": {
							Type:         schema.TypeString,
							ValidateFunc: validateSakuracloudIDType,
							Required:     true,
						},
						"ip_address": {
							Type:         schema.TypeString,
							ValidateFunc: validateIPv4Address(),
							Required:     true,
						},
					},
				},
			},
			"sim_route": {
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
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

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
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

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
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

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
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

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

	if err := cleanup.DeleteMobileGateway(ctx, mgwOp, simOp, zone, mgw.ID); err != nil {
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
	if err := d.Set("private_network_interface", flattenMobileGatewayPrivateNetworks(data)); err != nil {
		return err
	}
	d.Set("public_ip", flattenMobileGatewayPublicIPAddress(data))
	d.Set("public_netmask", flattenMobileGatewayPublicNetmask(data))
	d.Set("internet_connection", data.Settings.InternetConnectionEnabled.Bool())
	d.Set("inter_device_communication", data.Settings.InterDeviceCommunicationEnabled.Bool())

	if err := d.Set("traffic_control", flattenMobileGatewayTrafficConfigs(tc)); err != nil {
		return err
	}
	d.Set("dns_servers", []string{resolver.DNS1, resolver.DNS2})
	if err := d.Set("static_route", flattenMobileGatewayStaticRoutes(data.Settings.StaticRoute)); err != nil {
		return err
	}
	d.Set("name", data.Name)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	if err := d.Set("sim", flattenMobileGatewaySIMs(sims)); err != nil {
		return err
	}
	if err := d.Set("sim_route", flattenMobileGatewaySIMRoutes(simRoutes)); err != nil {
		return err
	}
	d.Set("zone", zone)

	return nil
}
