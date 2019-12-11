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
	"github.com/sacloud/libsacloud/v2/utils/builder"
	"math"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	mobileGatewayBuilder "github.com/sacloud/libsacloud/v2/utils/builder/mobilegateway"
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

func expandMobileGatewayBuilder(d *schema.ResourceData, client *APIClient) *mobileGatewayBuilder.Builder {
	return &mobileGatewayBuilder.Builder{
		Name:                            d.Get("name").(string),
		Description:                     d.Get("description").(string),
		Tags:                            expandTags(d),
		IconID:                          expandSakuraCloudID(d, "icon_id"),
		PrivateInterface:                expandMobileGatewayPrivateNetworks(d),
		StaticRoutes:                    expandMobileGatewayStaticRoutes(d),
		SIMRoutes:                       expandMobileGatewaySIMRoutes(d),
		InternetConnectionEnabled:       d.Get("internet_connection").(bool),
		InterDeviceCommunicationEnabled: d.Get("inter_device_communication").(bool),
		DNS:                             expandMobileGatewayDNSSetting(d),
		SIMs:                            expandMobileGatewaySIMs(d),
		TrafficConfig:                   expandMobileGatewayTrafficConfig(d),
		SetupOptions: &builder.RetryableSetupParameter{
			BootAfterBuild:        true,
			NICUpdateWaitDuration: builder.DefaultNICUpdateWaitDuration,
		},
		Client: mobileGatewayBuilder.NewAPIClient(client),
	}
}

func expandMobileGatewaySIMs(d resourceValueGettable) []*mobileGatewayBuilder.SIMSetting {
	var results []*mobileGatewayBuilder.SIMSetting
	if sims, ok := d.Get("sims").([]interface{}); ok && len(sims) > 0 {
		for _, v := range sims {
			sim := expandMobileGatewaySIM(mapToResourceData(v.(map[string]interface{})))
			results = append(results, sim)
		}
	}
	return results
}

func expandMobileGatewaySIM(d resourceValueGettable) *mobileGatewayBuilder.SIMSetting {
	return &mobileGatewayBuilder.SIMSetting{
		SIMID:     expandSakuraCloudID(d, "sim_id"),
		IPAddress: d.Get("ipaddress").(string),
	}
}

func flattenMobileGatewaySIMs(sims []*sacloud.MobileGatewaySIMInfo) []interface{} {
	var results []interface{}
	for _, sim := range sims {
		results = append(results, flattenMobileGatewaySIM(sim))
	}
	return results
}

func flattenMobileGatewaySIM(sim *sacloud.MobileGatewaySIMInfo) interface{} {
	return map[string]interface{}{
		"sim_id":    sim.ResourceID,
		"ipaddress": sim.IP,
	}
}

func expandMobileGatewayDNSSetting(d resourceValueGettable) *sacloud.MobileGatewayDNSSetting {
	dns1 := d.Get("dns_server1").(string)
	dns2 := d.Get("dns_server2").(string)
	if dns1 != "" || dns2 != "" {
		return &sacloud.MobileGatewayDNSSetting{
			DNS1: dns1,
			DNS2: dns2,
		}
	}
	return nil
}

func expandMobileGatewayPrivateNetworks(d resourceValueGettable) *mobileGatewayBuilder.PrivateInterfaceSetting {
	if raw, ok := d.Get("private_interface").([]interface{}); ok && len(raw) > 0 {
		d := mapToResourceData(raw[0].(map[string]interface{}))
		return &mobileGatewayBuilder.PrivateInterfaceSetting{
			SwitchID:       expandSakuraCloudID(d, "switch_id"),
			IPAddress:      d.Get("ipaddress").(string),
			NetworkMaskLen: d.Get("nw_mask_len").(int),
		}
	}
	return nil
}

func flattenMobileGatewayPrivateNetworks(mgw *sacloud.MobileGateway) []interface{} {
	if len(mgw.Interfaces) > 1 && !mgw.Interfaces[1].SwitchID.IsEmpty() {
		switchID := mgw.Interfaces[1].SwitchID
		var setting *sacloud.MobileGatewayInterfaceSetting
		for _, s := range mgw.Settings.Interfaces {
			if s.Index == 1 {
				setting = s
			}
		}
		if setting != nil {
			return []interface{}{
				map[string]interface{}{
					"switch_id":   switchID.String(),
					"ipaddress":   setting.IPAddress[0],
					"nw_mask_len": setting.NetworkMaskLen,
				},
			}
		}
	}
	return nil
}

func flattenMobileGatewayPublicNetworks(mgw *sacloud.MobileGateway) []interface{} {
	var results []interface{}
	if len(mgw.Interfaces) > 0 {
		results = append(results, map[string]interface{}{
			"ipaddress":   mgw.Interfaces[0].IPAddress,
			"nw_mask_len": mgw.Interfaces[0].SubnetNetworkMaskLen,
		})
	}
	return results
}

func expandMobileGatewaySIMRoutes(d resourceValueGettable) []*mobileGatewayBuilder.SIMRouteSetting {
	var routes []*mobileGatewayBuilder.SIMRouteSetting
	if simRoutes, ok := d.Get("sim_routes").([]interface{}); ok && len(simRoutes) > 0 {
		for _, v := range simRoutes {
			route := expandMobileGatewaySIMRoute(mapToResourceData(v.(map[string]interface{})))
			routes = append(routes, route)
		}
	}
	return routes
}

func expandMobileGatewaySIMRoute(d resourceValueGettable) *mobileGatewayBuilder.SIMRouteSetting {
	var simRoute = &mobileGatewayBuilder.SIMRouteSetting{
		Prefix: d.Get("prefix").(string),
		SIMID:  expandSakuraCloudID(d, "sim_id"),
	}
	return simRoute
}

func flattenMobileGatewaySIMRoutes(routes []*sacloud.MobileGatewaySIMRoute) []interface{} {
	var results []interface{}
	for _, route := range routes {
		results = append(results, flattenMobileGatewaySIMRoute(route))
	}
	return results
}

func flattenMobileGatewaySIMRoute(route *sacloud.MobileGatewaySIMRoute) interface{} {
	return map[string]interface{}{
		"sim_id": route.ResourceID,
		"prefix": route.Prefix,
	}
}

func expandMobileGatewayTrafficConfig(d resourceValueGettable) *sacloud.MobileGatewayTrafficControl {
	values := d.Get("traffic_control").([]interface{})
	if len(values) == 0 {
		return nil
	}
	v := &resourceMapValue{value: values[0].(map[string]interface{})}
	return &sacloud.MobileGatewayTrafficControl{
		TrafficQuotaInMB:       v.Get("quota").(int),
		BandWidthLimitInKbps:   v.Get("band_width_limit").(int),
		EmailNotifyEnabled:     v.Get("enable_email").(bool),
		SlackNotifyEnabled:     v.Get("enable_slack").(bool),
		SlackNotifyWebhooksURL: v.Get("slack_webhook").(string),
		AutoTrafficShaping:     v.Get("auto_traffic_shaping").(bool),
	}
}

func flattenMobileGatewayTrafficConfig(tc *sacloud.MobileGatewayTrafficControl) interface{} {
	return map[string]interface{}{
		"quota":                tc.TrafficQuotaInMB,
		"band_width_limit":     tc.BandWidthLimitInKbps,
		"auto_traffic_shaping": tc.AutoTrafficShaping,
		"enable_email":         tc.EmailNotifyEnabled,
		"enable_slack":         tc.SlackNotifyEnabled,
		"slack_webhook":        tc.SlackNotifyWebhooksURL,
	}
}

func flattenMobileGatewayTrafficConfigs(tc *sacloud.MobileGatewayTrafficControl) []interface{} {
	if tc == nil {
		return nil
	}
	return []interface{}{flattenMobileGatewayTrafficConfig(tc)}
}

func expandMobileGatewayStaticRoutes(d resourceValueGettable) []*sacloud.MobileGatewayStaticRoute {
	var routes []*sacloud.MobileGatewayStaticRoute
	if staticRoutes, ok := d.Get("static_routes").([]interface{}); ok && len(staticRoutes) > 0 {
		for _, v := range staticRoutes {
			route := expandMobileGatewayStaticRoute(&resourceMapValue{v.(map[string]interface{})})
			routes = append(routes, route)
		}
	}
	return routes
}

func expandMobileGatewayStaticRoute(d resourceValueGettable) *sacloud.MobileGatewayStaticRoute {
	return &sacloud.MobileGatewayStaticRoute{
		Prefix:  d.Get("prefix").(string),
		NextHop: d.Get("next_hop").(string),
	}
}

func flattenMobileGatewayStaticRoutes(routes []*sacloud.MobileGatewayStaticRoute) []interface{} {
	var results []interface{}
	for _, r := range routes {
		results = append(results, map[string]interface{}{
			"prefix":   r.Prefix,
			"next_hop": r.NextHop,
		})
	}
	return results
}
