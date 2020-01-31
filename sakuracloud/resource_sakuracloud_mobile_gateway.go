// Copyright 2016-2020 terraform-provider-sakuracloud authors
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
	resourceName := "MobileGateway"

	return &schema.Resource{
		Create: resourceSakuraCloudMobileGatewayCreate,
		Read:   resourceSakuraCloudMobileGatewayRead,
		Update: resourceSakuraCloudMobileGatewayUpdate,
		Delete: resourceSakuraCloudMobileGatewayDelete,
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
							Description:  descf("The id of the switch to which the %s connects", resourceName),
						},
						"ip_address": {
							Type:         schema.TypeString,
							ValidateFunc: validateIPv4Address(),
							Required:     true,
							Description:  descf("The IP address to assign to the %s", resourceName),
						},
						"netmask": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(8, 29),
							Description: descf(
								"The bit length of the subnet to assign to the %s. %s",
								resourceName,
								descRange(8, 29),
							),
						},
					},
				},
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("The public IP address assigned to the %s", resourceName),
			},
			"public_netmask": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: descf("The bit length of the subnet assigned to the %s", resourceName),
			},
			"internet_connection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "The flag to enable connect to the Internet",
			},
			"inter_device_communication": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "The flag to allow communication between each connected devices",
			},
			"dns_servers": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    2,
				MinItems:    2,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of IP address used by each connected devices",
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
							Description:  "The threshold of monthly traffic usage to enable to the traffic shaping",
						},
						"band_width_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, math.MaxInt32),
							Description:  "The bandwidth allowed when the traffic shaping is enabled",
						},
						"enable_email": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "The flag to enable email notification when the traffic shaping is enabled",
						},
						"enable_slack": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "The flag to enable slack notification when the traffic shaping is enabled",
						},
						"slack_webhook": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^https://hooks.slack.com/services/\w+/\w+/\w+$`), ""),
							Description:  "The webhook URL used when sends notification. It will only used when `enable_slack` is set `true`",
						},
						"auto_traffic_shaping": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "The flag to enable the traffic shaping",
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
							Description: "The destination network prefix used by static routing. This must be specified by CIDR block formatted string",
						},
						"next_hop": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.SingleIP(),
							Description:  "The IP address of next hop",
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
							Description:  descf("The id of the Switch connected to the %s", resourceName),
						},
						"ip_address": {
							Type:         schema.TypeString,
							ValidateFunc: validateIPv4Address(),
							Required:     true,
							Description:  "The IP address to assign to the SIM",
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
							Description:  "The id of the routing destination SIM",
						},
						"prefix": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The destination network prefix used by the sim routing. This must be specified by CIDR block formatted string",
						},
					},
				},
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"zone":        schemaResourceZone(resourceName),
		},
	}
}

func resourceSakuraCloudMobileGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
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
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
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
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
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
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
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
	d.Set("public_ip", flattenMobileGatewayPublicIPAddress(data))                             // nolint
	d.Set("public_netmask", flattenMobileGatewayPublicNetmask(data))                          // nolint
	d.Set("internet_connection", data.Settings.InternetConnectionEnabled.Bool())              // nolint
	d.Set("inter_device_communication", data.Settings.InterDeviceCommunicationEnabled.Bool()) // nolint

	if err := d.Set("traffic_control", flattenMobileGatewayTrafficConfigs(tc)); err != nil {
		return err
	}
	if err := d.Set("dns_servers", []string{resolver.DNS1, resolver.DNS2}); err != nil {
		return err
	}
	if err := d.Set("static_route", flattenMobileGatewayStaticRoutes(data.Settings.StaticRoute)); err != nil {
		return err
	}
	d.Set("name", data.Name)               // nolint
	d.Set("icon_id", data.IconID.String()) // nolint
	d.Set("description", data.Description) // nolint
	if err := d.Set("tags", flattenTags(data.Tags)); err != nil {
		return err
	}
	if err := d.Set("sim", flattenMobileGatewaySIMs(sims)); err != nil {
		return err
	}
	if err := d.Set("sim_route", flattenMobileGatewaySIMRoutes(simRoutes)); err != nil {
		return err
	}
	d.Set("zone", zone) // nolint

	return nil
}
