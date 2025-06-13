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
	"math"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/cleanup"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudMobileGateway() *schema.Resource {
	resourceName := "MobileGateway"

	return &schema.Resource{
		CreateContext: resourceSakuraCloudMobileGatewayCreate,
		ReadContext:   resourceSakuraCloudMobileGatewayRead,
		UpdateContext: resourceSakuraCloudMobileGatewayUpdate,
		DeleteContext: resourceSakuraCloudMobileGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
							Required:         true,
							Description:      desc.Sprintf("The id of the switch to which the %s connects", resourceName),
						},
						"ip_address": {
							Type:             schema.TypeString,
							ValidateDiagFunc: validateIPv4Address(),
							Required:         true,
							Description:      desc.Sprintf("The IP address to assign to the %s", resourceName),
						},
						"netmask": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(8, 29)),
							Description: desc.Sprintf(
								"The bit length of the subnet to assign to the %s. %s",
								resourceName,
								desc.Range(8, 29),
							),
						},
					},
				},
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: desc.Sprintf("The public IP address assigned to the %s", resourceName),
			},
			"public_netmask": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: desc.Sprintf("The bit length of the subnet assigned to the %s", resourceName),
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
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, math.MaxInt32)),
							Description:      "The threshold of monthly traffic usage to enable to the traffic shaping",
						},
						"band_width_limit": {
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, math.MaxInt32)),
							Description:      "The bandwidth allowed when the traffic shaping is enabled",
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
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(`^https://hooks.slack.com/services/\w+/\w+/\w+$`), "slack_webhook")),
							Description:      "The webhook URL used when sends notification. It will only used when `enable_slack` is set `true`",
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
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
							Description:      "The IP address of next hop",
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
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
							Required:         true,
							Description:      desc.Sprintf("The id of the Switch connected to the %s", resourceName),
						},
						"ip_address": {
							Type:             schema.TypeString,
							ValidateDiagFunc: validateIPv4Address(),
							Required:         true,
							Description:      "The IP address to assign to the SIM",
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
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
							Required:         true,
							Description:      "The id of the routing destination SIM",
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

func resourceSakuraCloudMobileGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	builder := expandMobileGatewayBuilder(d, client, zone)
	if err := builder.Validate(ctx, zone); err != nil {
		return diag.Errorf("validating SakuraCloud MobileGateway is failed: %s", err)
	}

	mgw, err := builder.Build(ctx)
	if mgw != nil {
		d.SetId(mgw.ID.String())
	}
	if err != nil {
		return diag.Errorf("creating SakuraCloud MobileGateway is failed: %s", err)
	}

	return resourceSakuraCloudMobileGatewayRead(ctx, d, meta)
}

func resourceSakuraCloudMobileGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	mgwOp := iaas.NewMobileGatewayOp(client)

	mgw, err := mgwOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud MobileGateway[%s]: %s", d.Id(), err)
	}

	return setMobileGatewayResourceData(ctx, d, client, mgw)
}

func resourceSakuraCloudMobileGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	mgwOp := iaas.NewMobileGatewayOp(client)

	mgw, err := mgwOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud MobileGateway[%s]: %s", d.Id(), err)
	}

	builder := expandMobileGatewayBuilder(d, client, zone)
	if err := builder.Validate(ctx, zone); err != nil {
		return diag.Errorf("validating SakuraCloud MobileGateway is failed: %s", err)
	}
	builder.ID = mgw.ID

	if _, err = builder.Build(ctx); err != nil {
		return diag.Errorf("updating SakuraCloud MobileGateway[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudMobileGatewayRead(ctx, d, meta)
}

func resourceSakuraCloudMobileGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	mgwOp := iaas.NewMobileGatewayOp(client)
	simOp := iaas.NewSIMOp(client)

	mgw, err := mgwOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud MobileGateway: %s", err)
	}

	if err := cleanup.DeleteMobileGateway(ctx, mgwOp, simOp, zone, mgw.ID); err != nil {
		return diag.Errorf("deleting SakuraCloud MobileGateway[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setMobileGatewayResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *iaas.MobileGateway) diag.Diagnostics {
	zone := getZone(d, client)
	mgwOp := iaas.NewMobileGatewayOp(client)

	if data.Availability.IsFailed() {
		d.SetId("")
		return diag.Errorf("got unexpected state: MobileGateway[%d].Availability is failed", data.ID)
	}

	// fetch configs
	tc, err := mgwOp.GetTrafficConfig(ctx, zone, data.ID)
	if err != nil && !iaas.IsNotFoundError(err) {
		return diag.Errorf("reading TrafficConfig is failed: %s", err)
	}
	resolver, err := mgwOp.GetDNS(ctx, zone, data.ID)
	if err != nil {
		return diag.Errorf("reading ResolverConfig is failed: %s", err)
	}
	sims, err := mgwOp.ListSIM(ctx, zone, data.ID)
	if err != nil && !iaas.IsNotFoundError(err) {
		return diag.Errorf("reading SIMs is failed: %s", err)
	}
	simRoutes, err := mgwOp.GetSIMRoutes(ctx, zone, data.ID)
	if err != nil {
		return diag.Errorf("reading SIM Routes is failed: %s", err)
	}

	// set data
	if err := d.Set("private_network_interface", flattenMobileGatewayPrivateNetworks(data)); err != nil {
		return diag.FromErr(err)
	}
	d.Set("public_ip", flattenMobileGatewayPublicIPAddress(data))                    //nolint
	d.Set("public_netmask", flattenMobileGatewayPublicNetmask(data))                 //nolint
	d.Set("internet_connection", data.InternetConnectionEnabled.Bool())              //nolint
	d.Set("inter_device_communication", data.InterDeviceCommunicationEnabled.Bool()) //nolint

	if err := d.Set("traffic_control", flattenMobileGatewayTrafficConfigs(tc)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dns_servers", []string{resolver.DNS1, resolver.DNS2}); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("static_route", flattenMobileGatewayStaticRoutes(data.StaticRoutes)); err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", data.Name)               //nolint
	d.Set("icon_id", data.IconID.String()) //nolint
	d.Set("description", data.Description) //nolint
	if err := d.Set("tags", flattenTags(data.Tags)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("sim", flattenMobileGatewaySIMs(sims)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("sim_route", flattenMobileGatewaySIMRoutes(simRoutes)); err != nil {
		return diag.FromErr(err)
	}
	d.Set("zone", zone) //nolint

	return nil
}
