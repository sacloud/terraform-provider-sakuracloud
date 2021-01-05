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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/helper/builder/localrouter"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func expandLocalRouterBuilder(d *schema.ResourceData, client *APIClient) *localrouter.Builder {
	return &localrouter.Builder{
		Name:         stringOrDefault(d, "name"),
		Description:  stringOrDefault(d, "description"),
		Tags:         expandTags(d),
		IconID:       expandSakuraCloudID(d, "icon_id"),
		Switch:       expandLocalRouterSwitch(d),
		Interface:    expandLocalRouterNetworkInterface(d),
		Peers:        expandLocalRouterPeers(d),
		StaticRoutes: expandLocalStaticRoutes(d),
		Client:       localrouter.NewAPIClient(client),
	}
}

func expandLocalRouterSwitch(d resourceValueGettable) *sacloud.LocalRouterSwitch {
	if values, ok := getListFromResource(d, "switch"); ok && len(values) > 0 {
		d = mapToResourceData(values[0].(map[string]interface{}))
		return &sacloud.LocalRouterSwitch{
			Code:     stringOrDefault(d, "code"),
			Category: stringOrDefault(d, "category"),
			ZoneID:   stringOrDefault(d, "zone_id"),
		}
	}
	return nil
}

func flattenLocalRouterSwitch(data *sacloud.LocalRouter) []interface{} {
	if data.Switch != nil {
		return []interface{}{
			map[string]interface{}{
				"code":     data.Switch.Code,
				"category": data.Switch.Category,
				"zone_id":  data.Switch.ZoneID,
			},
		}
	}
	return nil
}

func expandLocalRouterNetworkInterface(d resourceValueGettable) *sacloud.LocalRouterInterface {
	if values, ok := getListFromResource(d, "network_interface"); ok && len(values) > 0 {
		d = mapToResourceData(values[0].(map[string]interface{}))
		return &sacloud.LocalRouterInterface{
			VirtualIPAddress: stringOrDefault(d, "vip"),
			IPAddress:        expandStringList(d.Get("ip_addresses").([]interface{})),
			NetworkMaskLen:   intOrDefault(d, "netmask"),
			VRID:             intOrDefault(d, "vrid"),
		}
	}
	return nil
}

func flattenLocalRouterNetworkInterface(data *sacloud.LocalRouter) []interface{} {
	if data.Interface != nil {
		return []interface{}{
			map[string]interface{}{
				"vip":          data.Interface.VirtualIPAddress,
				"ip_addresses": data.Interface.IPAddress,
				"netmask":      data.Interface.NetworkMaskLen,
				"vrid":         data.Interface.VRID,
			},
		}
	}
	return nil
}

func expandLocalRouterPeers(d resourceValueGettable) []*sacloud.LocalRouterPeer {
	var results []*sacloud.LocalRouterPeer
	if values, ok := getListFromResource(d, "peer"); ok && len(values) > 0 {
		for _, raw := range values {
			d = mapToResourceData(raw.(map[string]interface{}))
			results = append(results, &sacloud.LocalRouterPeer{
				ID:          expandSakuraCloudID(d, "peer_id"),
				SecretKey:   stringOrDefault(d, "secret_key"),
				Enabled:     boolOrDefault(d, "enabled"),
				Description: stringOrDefault(d, "description"),
			})
		}
	}
	return results
}

func flattenLocalRouterPeers(data *sacloud.LocalRouter) []interface{} {
	var results []interface{}
	for _, peer := range data.Peers {
		results = append(results, flattenLocalRouterPeer(peer))
	}
	return results
}

func flattenLocalRouterPeer(data *sacloud.LocalRouterPeer) interface{} {
	return map[string]interface{}{
		"peer_id":     data.ID.String(),
		"secret_key":  data.SecretKey,
		"enabled":     data.Enabled,
		"description": data.Description,
	}
}

func expandLocalStaticRoutes(d resourceValueGettable) []*sacloud.LocalRouterStaticRoute {
	var results []*sacloud.LocalRouterStaticRoute
	if values, ok := getListFromResource(d, "static_route"); ok && len(values) > 0 {
		for _, raw := range values {
			d = mapToResourceData(raw.(map[string]interface{}))
			results = append(results, &sacloud.LocalRouterStaticRoute{
				Prefix:  stringOrDefault(d, "prefix"),
				NextHop: stringOrDefault(d, "next_hop"),
			})
		}
	}
	return results
}

func flattenLocalRouterStaticRoutes(data *sacloud.LocalRouter) []interface{} {
	var results []interface{}
	for _, route := range data.StaticRoutes {
		results = append(results, flattenLocalRouterStaticRoute(route))
	}
	return results
}

func flattenLocalRouterStaticRoute(data *sacloud.LocalRouterStaticRoute) interface{} {
	return map[string]interface{}{
		"prefix":   data.Prefix,
		"next_hop": data.NextHop,
	}
}
