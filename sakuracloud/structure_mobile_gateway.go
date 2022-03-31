// Copyright 2016-2022 terraform-provider-sakuracloud authors
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/defaults"
	mobileGatewayBuilder "github.com/sacloud/iaas-service-go/mobilegateway/builder"
	"github.com/sacloud/iaas-service-go/setup"
)

func expandMobileGatewayBuilder(d *schema.ResourceData, client *APIClient, zone string) *mobileGatewayBuilder.Builder {
	return &mobileGatewayBuilder.Builder{
		Zone: zone,

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
		SetupOptions: &setup.Options{
			BootAfterBuild:        true,
			NICUpdateWaitDuration: defaults.DefaultNICUpdateWaitDuration,
		},
		Client: mobileGatewayBuilder.NewAPIClient(client),
	}
}

func expandMobileGatewaySIMs(d resourceValueGettable) []*mobileGatewayBuilder.SIMSetting {
	var results []*mobileGatewayBuilder.SIMSetting
	if sims, ok := d.Get("sim").([]interface{}); ok && len(sims) > 0 {
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
		IPAddress: d.Get("ip_address").(string),
	}
}

func flattenMobileGatewaySIMs(sims []*iaas.MobileGatewaySIMInfo) []interface{} {
	var results []interface{}
	for _, sim := range sims {
		results = append(results, flattenMobileGatewaySIM(sim))
	}
	return results
}

func flattenMobileGatewaySIM(sim *iaas.MobileGatewaySIMInfo) interface{} {
	return map[string]interface{}{
		"sim_id":     sim.ResourceID,
		"ip_address": sim.IP,
	}
}

func expandMobileGatewayDNSSetting(d resourceValueGettable) *iaas.MobileGatewayDNSSetting {
	servers := d.Get("dns_servers").([]interface{})
	if len(servers) == 2 && servers[0].(string) != "" && servers[1].(string) != "" {
		return &iaas.MobileGatewayDNSSetting{
			DNS1: servers[0].(string),
			DNS2: servers[1].(string),
		}
	}
	return nil
}

func expandMobileGatewayPrivateNetworks(d resourceValueGettable) *mobileGatewayBuilder.PrivateInterfaceSetting {
	if raw, ok := d.Get("private_network_interface").([]interface{}); ok && len(raw) > 0 {
		d := mapToResourceData(raw[0].(map[string]interface{}))
		return &mobileGatewayBuilder.PrivateInterfaceSetting{
			SwitchID:       expandSakuraCloudID(d, "switch_id"),
			IPAddress:      d.Get("ip_address").(string),
			NetworkMaskLen: d.Get("netmask").(int),
		}
	}
	return nil
}

func flattenMobileGatewayPrivateNetworks(mgw *iaas.MobileGateway) []interface{} {
	if len(mgw.Interfaces) > 1 && !mgw.Interfaces[1].SwitchID.IsEmpty() {
		switchID := mgw.Interfaces[1].SwitchID
		var setting *iaas.MobileGatewayInterfaceSetting
		for _, s := range mgw.InterfaceSettings {
			if s.Index == 1 {
				setting = s
			}
		}
		if setting != nil {
			return []interface{}{
				map[string]interface{}{
					"switch_id":  switchID.String(),
					"ip_address": setting.IPAddress[0],
					"netmask":    setting.NetworkMaskLen,
				},
			}
		}
	}
	return nil
}

func flattenMobileGatewayPublicNetmask(mgw *iaas.MobileGateway) int {
	if len(mgw.Interfaces) > 0 {
		return mgw.Interfaces[0].SubnetNetworkMaskLen
	}
	return 0
}

func flattenMobileGatewayPublicIPAddress(mgw *iaas.MobileGateway) string {
	if len(mgw.Interfaces) > 0 {
		return mgw.Interfaces[0].IPAddress
	}
	return ""
}

func expandMobileGatewaySIMRoutes(d resourceValueGettable) []*mobileGatewayBuilder.SIMRouteSetting {
	var routes []*mobileGatewayBuilder.SIMRouteSetting
	if simRoutes, ok := d.Get("sim_route").([]interface{}); ok && len(simRoutes) > 0 {
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

func flattenMobileGatewaySIMRoutes(routes []*iaas.MobileGatewaySIMRoute) []interface{} {
	var results []interface{}
	for _, route := range routes {
		results = append(results, flattenMobileGatewaySIMRoute(route))
	}
	return results
}

func flattenMobileGatewaySIMRoute(route *iaas.MobileGatewaySIMRoute) interface{} {
	return map[string]interface{}{
		"sim_id": route.ResourceID,
		"prefix": route.Prefix,
	}
}

func expandMobileGatewayTrafficConfig(d resourceValueGettable) *iaas.MobileGatewayTrafficControl {
	values := d.Get("traffic_control").([]interface{})
	if len(values) == 0 {
		return nil
	}
	v := &resourceMapValue{value: values[0].(map[string]interface{})}
	return &iaas.MobileGatewayTrafficControl{
		TrafficQuotaInMB:       v.Get("quota").(int),
		BandWidthLimitInKbps:   v.Get("band_width_limit").(int),
		EmailNotifyEnabled:     v.Get("enable_email").(bool),
		SlackNotifyEnabled:     v.Get("enable_slack").(bool),
		SlackNotifyWebhooksURL: v.Get("slack_webhook").(string),
		AutoTrafficShaping:     v.Get("auto_traffic_shaping").(bool),
	}
}

func flattenMobileGatewayTrafficConfig(tc *iaas.MobileGatewayTrafficControl) interface{} {
	return map[string]interface{}{
		"quota":                tc.TrafficQuotaInMB,
		"band_width_limit":     tc.BandWidthLimitInKbps,
		"auto_traffic_shaping": tc.AutoTrafficShaping,
		"enable_email":         tc.EmailNotifyEnabled,
		"enable_slack":         tc.SlackNotifyEnabled,
		"slack_webhook":        tc.SlackNotifyWebhooksURL,
	}
}

func flattenMobileGatewayTrafficConfigs(tc *iaas.MobileGatewayTrafficControl) []interface{} {
	if tc == nil {
		return nil
	}
	return []interface{}{flattenMobileGatewayTrafficConfig(tc)}
}

func expandMobileGatewayStaticRoutes(d resourceValueGettable) []*iaas.MobileGatewayStaticRoute {
	var routes []*iaas.MobileGatewayStaticRoute
	if staticRoutes, ok := d.Get("static_route").([]interface{}); ok && len(staticRoutes) > 0 {
		for _, v := range staticRoutes {
			route := expandMobileGatewayStaticRoute(&resourceMapValue{v.(map[string]interface{})})
			routes = append(routes, route)
		}
	}
	return routes
}

func expandMobileGatewayStaticRoute(d resourceValueGettable) *iaas.MobileGatewayStaticRoute {
	return &iaas.MobileGatewayStaticRoute{
		Prefix:  d.Get("prefix").(string),
		NextHop: d.Get("next_hop").(string),
	}
}

func flattenMobileGatewayStaticRoutes(routes []*iaas.MobileGatewayStaticRoute) []interface{} {
	var results []interface{}
	for _, r := range routes {
		results = append(results, map[string]interface{}{
			"prefix":   r.Prefix,
			"next_hop": r.NextHop,
		})
	}
	return results
}
