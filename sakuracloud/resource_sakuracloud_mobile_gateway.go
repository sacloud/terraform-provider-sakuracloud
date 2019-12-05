package sakuracloud

import (
	"context"
	"errors"
	"fmt"
	"math"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/setup"
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
			"switch_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"public_ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_nw_mask_len": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"private_ipaddress": {
				Type:         schema.TypeString,
				ValidateFunc: validateIPv4Address(),
				Optional:     true,
			},
			"private_nw_mask_len": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(8, 29),
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
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIPv4Address(),
			},
			"dns_server2": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
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
			"static_route": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sim_ids": {
				Type:     schema.TypeList,
				Computed: true,
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
	client, ctx, zone := getSacloudV2Client(d, meta)
	mgwOp := sacloud.NewMobileGatewayOp(client)

	// validate
	switchID, ip, nwMaskLen, err := expandMobileGatewayPrivateNetworks(d)
	if err != nil {
		return err
	}

	opts := &sacloud.MobileGatewayCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTagsV2(d.Get("tags").([]interface{})),
		IconID:      expandSakuraCloudID(d, "icon_id"),
		Settings: &sacloud.MobileGatewaySettingCreate{
			InternetConnectionEnabled:       expandStringFlag(d, "internet_connection"),
			InterDeviceCommunicationEnabled: expandStringFlag(d, "inter_device_communication"),
		},
	}

	builder := &setup.RetryableSetup{
		Create: func(ctx context.Context, zone string) (accessor.ID, error) {
			return mgwOp.Create(ctx, zone, opts)
		},
		Delete: func(ctx context.Context, zone string, id types.ID) error {
			return mgwOp.Delete(ctx, zone, id)
		},
		Read: func(ctx context.Context, zone string, id types.ID) (interface{}, error) {
			return mgwOp.Read(ctx, zone, id)
		},
		IsWaitForCopy: true,
		IsWaitForUp:   false,
	}

	result, err := builder.Setup(ctx, zone)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud MobileGateway is failed: %s", err)
	}
	mgw := result.(*sacloud.MobileGateway)

	// connect to switch
	if !switchID.IsEmpty() {
		patchParam := &sacloud.MobileGatewayPatchRequest{
			Settings:     mgw.Settings,
			SettingsHash: mgw.SettingsHash,
		}

		if err := mgwOp.ConnectToSwitch(ctx, zone, mgw.ID, switchID); err != nil {
			return fmt.Errorf("connecting to switch is failed: %s", err)
		}

		patchParam.Settings.Interfaces = append(patchParam.Settings.Interfaces, &sacloud.MobileGatewayInterfaceSetting{
			IPAddress:      []string{ip},
			NetworkMaskLen: nwMaskLen,
			Index:          1,
		})

		upd, err := mgwOp.Patch(ctx, zone, mgw.ID, patchParam)
		if err != nil {
			return fmt.Errorf("updating network settings is failed: %s", err)
		}
		mgw = upd
	}

	// traffic config
	if tc := expandMobileGatewayTrafficConfig(d); tc != nil {
		if err := mgwOp.SetTrafficConfig(ctx, zone, mgw.ID, tc); err != nil {
			return fmt.Errorf("updating traffic config is failed: %s", err)
		}
	}

	// dns
	dns1 := d.Get("dns_server1").(string)
	dns2 := d.Get("dns_server2").(string)
	if dns1 != "" || dns2 != "" {
		err := mgwOp.SetDNS(ctx, zone, mgw.ID, &sacloud.MobileGatewayDNSSetting{
			DNS1: dns1,
			DNS2: dns2,
		})
		if err != nil {
			return fmt.Errorf("updating dns settings is failed: %s", err)
		}
	}

	// static route
	staticRoutes := expandMobileGatewayStaticRoutes(d)
	if len(staticRoutes) > 0 {
		patchParam := &sacloud.MobileGatewayPatchRequest{
			Settings:     mgw.Settings,
			SettingsHash: mgw.SettingsHash,
		}
		patchParam.Settings.StaticRoute = staticRoutes

		upd, err := mgwOp.Patch(ctx, zone, mgw.ID, patchParam)
		if err != nil {
			return fmt.Errorf("updating static routes is failed: %s", err)
		}
		mgw = upd
	}

	// boot
	if err := mgwOp.Boot(ctx, zone, mgw.ID); err != nil {
		return fmt.Errorf("booting SakuraCloud MobileGateway is failed: %s", err)
	}
	_, err = sacloud.WaiterForUp(func() (interface{}, error) {
		return mgwOp.Read(ctx, zone, mgw.ID)
	}).WaitForState(ctx)

	if err != nil {
		return fmt.Errorf("booting SakuraCloud MobileGateway is failed: %s", err)
	}

	d.SetId(mgw.ID.String())
	return resourceSakuraCloudMobileGatewayRead(d, meta)

}

func resourceSakuraCloudMobileGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	mgwOp := sacloud.NewMobileGatewayOp(client)

	mgw, err := mgwOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud MobileGateway: %s", err)
	}

	return setMobileGatewayResourceData(ctx, d, client, mgw)
}

func resourceSakuraCloudMobileGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	mgwOp := sacloud.NewMobileGatewayOp(client)

	switchID, ip, nwMaskLen, err := expandMobileGatewayPrivateNetworks(d)
	if err != nil {
		return err
	}

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	mgw, err := mgwOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud MobileGateway: %s", err)
	}

	mgw, err = mgwOp.Update(ctx, zone, mgw.ID, &sacloud.MobileGatewayUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTagsV2(d.Get("tags").([]interface{})),
		IconID:      expandSakuraCloudID(d, "icon_id"),
		Settings: &sacloud.MobileGatewaySetting{
			Interfaces:                      mgw.Settings.Interfaces,
			StaticRoute:                     mgw.Settings.StaticRoute,
			InternetConnectionEnabled:       expandStringFlag(d, "internet_connection"),
			InterDeviceCommunicationEnabled: expandStringFlag(d, "inter_device_communication"),
		},
		SettingsHash: mgw.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("updating SakuraCloud MobileGateway is failed: %s", err)
	}

	swOptsHasChange := d.HasChange("switch_id") || d.HasChange("private_ipaddress") || d.HasChange("private_nw_mask_len")
	needRestart := false
	if mgw.InstanceStatus.IsUp() && swOptsHasChange {
		needRestart = true
	}

	// shutdown
	if needRestart {
		if err := mgwOp.Shutdown(ctx, zone, mgw.ID, nil); err != nil {
			return fmt.Errorf("updating SakuraCloud MobileGateway is failed: %s", err)
		}
		_, err = sacloud.WaiterForDown(func() (interface{}, error) {
			return mgwOp.Read(ctx, zone, mgw.ID)
		}).WaitForState(ctx)
		if err != nil {
			return fmt.Errorf("updating SakuraCloud MobileGateway is failed: %s", err)
		}
	}

	if swOptsHasChange {
		updateParam := &sacloud.MobileGatewayUpdateRequest{
			Name:         mgw.Name,
			Description:  mgw.Description,
			Tags:         mgw.Tags,
			IconID:       mgw.IconID,
			Settings:     mgw.Settings,
			SettingsHash: mgw.SettingsHash,
		}

		// disconnect from switch if already connected
		if len(mgw.Interfaces) > 1 && !mgw.Interfaces[1].SwitchID.IsEmpty() {
			if err := mgwOp.DisconnectFromSwitch(ctx, zone, mgw.ID); err != nil {
				return fmt.Errorf("disconnecting from switch is failed: %s", err)
			}
		}

		if !switchID.IsEmpty() {
			if err := mgwOp.ConnectToSwitch(ctx, zone, mgw.ID, switchID); err != nil {
				return fmt.Errorf("connecting to switch is failed: %s", err)
			}

			updateParam.Settings.Interfaces = append(updateParam.Settings.Interfaces, &sacloud.MobileGatewayInterfaceSetting{
				IPAddress:      []string{ip},
				NetworkMaskLen: nwMaskLen,
				Index:          1,
			})
		} else {
			var ifs []*sacloud.MobileGatewayInterfaceSetting
			for _, i := range updateParam.Settings.Interfaces {
				if i.Index != 1 {
					ifs = append(ifs, i)
				}
			}
			updateParam.Settings.Interfaces = ifs
		}

		upd, err := mgwOp.Update(ctx, zone, mgw.ID, updateParam)
		if err != nil {
			return fmt.Errorf("updating network settings is failed: %s", err)
		}
		mgw = upd
	}

	if d.HasChange("traffic_control") {
		if tc := expandMobileGatewayTrafficConfig(d); tc != nil {
			if err := mgwOp.SetTrafficConfig(ctx, zone, mgw.ID, tc); err != nil {
				return fmt.Errorf("updating traffic config is failed: %s", err)
			}
		} else {
			if err := mgwOp.DeleteTrafficConfig(ctx, zone, mgw.ID); err != nil {
				return fmt.Errorf("updating traffic config is failed: %s", err)
			}
		}
	}

	if d.HasChange("dns1") || d.HasChange("dns2") {
		dns1 := d.Get("dns_server1").(string)
		dns2 := d.Get("dns_server2").(string)
		if dns1 != "" || dns2 != "" {
			err := mgwOp.SetDNS(ctx, zone, mgw.ID, &sacloud.MobileGatewayDNSSetting{
				DNS1: dns1,
				DNS2: dns2,
			})
			if err != nil {
				return fmt.Errorf("updating dns settings is failed: %s", err)
			}
		}
	}

	// static route
	if d.HasChange("static_route") {
		mgw.Settings.StaticRoute = expandMobileGatewayStaticRoutes(d)
		updateParam := &sacloud.MobileGatewayUpdateRequest{
			Name:         mgw.Name,
			Description:  mgw.Description,
			Tags:         mgw.Tags,
			IconID:       mgw.IconID,
			Settings:     mgw.Settings,
			SettingsHash: mgw.SettingsHash,
		}
		upd, err := mgwOp.Update(ctx, zone, mgw.ID, updateParam)
		if err != nil {
			return fmt.Errorf("updating static routes is failed: %s", err)
		}
		mgw = upd
	}

	if needRestart {
		if err := mgwOp.Boot(ctx, zone, mgw.ID); err != nil {
			return fmt.Errorf("updating SakuraCloud MobileGateway is failed: %s", err)
		}
		_, err = sacloud.WaiterForUp(func() (interface{}, error) {
			return mgwOp.Read(ctx, zone, mgw.ID)
		}).WaitForState(ctx)

		if err != nil {
			return fmt.Errorf("updating SakuraCloud MobileGateway is failed: %s", err)
		}
	}

	return resourceSakuraCloudMobileGatewayRead(d, meta)
}

func resourceSakuraCloudMobileGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	mgwOp := sacloud.NewMobileGatewayOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	mgw, err := mgwOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud MobileGateway: %s", err)
	}

	if err := mgwOp.Shutdown(ctx, zone, mgw.ID, &sacloud.ShutdownOption{Force: true}); err != nil {
		return fmt.Errorf("stopping SakuraCloud MobileGateway is failed: %s", err)
	}
	if _, err := sacloud.WaiterForDown(func() (interface{}, error) {
		return mgwOp.Read(ctx, zone, mgw.ID)
	}).WaitForState(ctx); err != nil {
		return fmt.Errorf("stopping SakuraCloud MobileGateway is failed: %s", err)
	}

	if err := mgwOp.Delete(ctx, zone, mgw.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud MobileGateway is failed: %s", err)
	}

	return nil
}

func setMobileGatewayResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.MobileGateway) error {
	zone := getV2Zone(d, client)
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
		return fmt.Errorf("reading resolver config is failed: %s", err)
	}
	sims, err := mgwOp.ListSIM(ctx, zone, data.ID)
	if err != nil {
		return fmt.Errorf("reading sim list is failed: %s", err)
	}

	// set data
	d.Set("public_ipaddress", data.Interfaces[0].IPAddress)
	d.Set("public_nw_mask_len", data.Interfaces[0].SubnetNetworkMaskLen)
	d.Set("internet_connection", data.Settings.InternetConnectionEnabled.Bool())
	d.Set("inter_device_communication", data.Settings.InterDeviceCommunicationEnabled.Bool())

	if len(data.Interfaces) > 1 && !data.Interfaces[1].SwitchID.IsEmpty() {
		d.Set("switch_id", data.Interfaces[1].SwitchID.String())
		d.Set("private_ipaddress", data.Settings.Interfaces[0].IPAddress[0])
		d.Set("private_nw_mask_len", data.Settings.Interfaces[0].NetworkMaskLen)
	}

	if tc != nil {
		if err := d.Set("traffic_control", flattenMobileGatewayTrafficConfigs(tc)); err != nil {
			return err
		}
	}

	d.Set("dns_server1", resolver.DNS1)
	d.Set("dns_server2", resolver.DNS2)

	if err := d.Set("static_route", flattenMobileGatewayStaticRoutes(data.Settings.StaticRoute)); err != nil {
		return err
	}

	d.Set("name", data.Name)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}

	if err := d.Set("sim_ids", flattenMobileGatewaySIMList(sims)); err != nil {
		return err
	}
	d.Set("zone", zone)

	return nil
}

func expandMobileGatewayPrivateNetworks(d resourceValueGettable) (switchID types.ID, ip string, nwMaskLen int, err error) {
	if rawSwitchID, ok := d.GetOk("switch_id"); ok {
		switchID = types.StringID(rawSwitchID.(string))
		if !switchID.IsEmpty() {
			ip = d.Get("private_ipaddress").(string)
			nwMaskLen = d.Get("private_nw_mask_len").(int)

			if ip == "" || nwMaskLen == 0 {
				err = errors.New("private_ipaddress and private_nw_mask_len is required when switch_id is specified")
			}
		}
	}
	return
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
	return []interface{}{flattenMobileGatewayTrafficConfig(tc)}
}

func expandMobileGatewayStaticRoutes(d resourceValueGettable) []*sacloud.MobileGatewayStaticRoute {
	var routes []*sacloud.MobileGatewayStaticRoute
	if staticRoutes, ok := d.Get("static_route").([]interface{}); ok && len(staticRoutes) > 0 {
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

func flattenMobileGatewaySIMList(sims []*sacloud.MobileGatewaySIMInfo) []interface{} {
	var result []interface{}
	for _, s := range sims {
		result = append(result, s.ResourceID)
	}
	return result
}
