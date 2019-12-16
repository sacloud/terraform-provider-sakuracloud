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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func expandLoadBalancerVIPs(d resourceValueGettable) []*sacloud.LoadBalancerVirtualIPAddress {
	var vips []*sacloud.LoadBalancerVirtualIPAddress
	vipsConf := d.Get("vip").([]interface{})
	for _, vip := range vipsConf {
		v := &resourceMapValue{vip.(map[string]interface{})}
		vips = append(vips, expandLoadBalancerVIP(v))
	}
	return vips
}

func expandLoadBalancerVIP(d resourceValueGettable) *sacloud.LoadBalancerVirtualIPAddress {
	servers := expandLoadBalancerServers(d, d.Get("port").(int))
	return &sacloud.LoadBalancerVirtualIPAddress{
		VirtualIPAddress: d.Get("vip").(string),
		Port:             types.StringNumber(d.Get("port").(int)),
		DelayLoop:        types.StringNumber(d.Get("delay_loop").(int)),
		SorryServer:      d.Get("sorry_server").(string),
		Description:      d.Get("description").(string),
		Servers:          servers,
	}
}

func flattenLoadBalancerVIPs(lb *sacloud.LoadBalancer) []interface{} {
	var vips []interface{}
	for _, v := range lb.VirtualIPAddresses {
		vips = append(vips, flattenLoadBalancerVIP(v))
	}
	return vips
}

func flattenLoadBalancerVIP(vip *sacloud.LoadBalancerVirtualIPAddress) interface{} {
	return map[string]interface{}{
		"vip":          vip.VirtualIPAddress,
		"port":         vip.Port.Int(),
		"delay_loop":   vip.DelayLoop.Int(),
		"sorry_server": vip.SorryServer,
		"server":       flattenLoadBalancerServers(vip.Servers),
	}
}

func expandLoadBalancerServers(d resourceValueGettable, vipPort int) []*sacloud.LoadBalancerServer {
	var servers []*sacloud.LoadBalancerServer
	for _, v := range d.Get("server").([]interface{}) {
		data := &resourceMapValue{v.(map[string]interface{})}
		server := expandLoadBalancerServer(data, vipPort)
		servers = append(servers, server)
	}
	return servers
}

func expandLoadBalancerServer(d resourceValueGettable, vipPort int) *sacloud.LoadBalancerServer {
	return &sacloud.LoadBalancerServer{
		IPAddress: d.Get("ipaddress").(string),
		Port:      types.StringNumber(vipPort),
		Enabled:   expandStringFlag(d, "enabled"),
		HealthCheck: &sacloud.LoadBalancerServerHealthCheck{
			Protocol:     types.ELoadBalancerHealthCheckProtocol(d.Get("check_protocol").(string)),
			Path:         d.Get("check_path").(string),
			ResponseCode: expandStringNumber(d, "check_status"),
		},
	}
}

func flattenLoadBalancerServers(servers []*sacloud.LoadBalancerServer) []interface{} {
	var results []interface{}
	for _, s := range servers {
		results = append(results, flattenLoadBalancerServer(s))
	}
	return results
}

func flattenLoadBalancerServer(server *sacloud.LoadBalancerServer) interface{} {
	return map[string]interface{}{
		"ipaddress":      server.IPAddress,
		"check_protocol": server.HealthCheck.Protocol,
		"check_path":     server.HealthCheck.Path,
		"check_status":   server.HealthCheck.ResponseCode.String(),
		"enabled":        server.Enabled.Bool(),
	}
}

func expandLoadBalancerPlanID(d resourceValueGettable) types.ID {
	plan := d.Get("plan").(string)
	if plan == "standard" {
		return types.LoadBalancerPlans.Standard
	}

	return types.LoadBalancerPlans.Premium
}

func flattenLoadBalancerPlanID(lb *sacloud.LoadBalancer) string {
	var plan string
	switch lb.PlanID {
	case types.LoadBalancerPlans.Standard:
		plan = "standard"
	case types.LoadBalancerPlans.Premium:
		plan = "highspec"
	}
	return plan
}

func expandLoadBalancerIPAddresses(d resourceValueGettable) []string {
	ipAddresses := []string{d.Get("ipaddress1").(string)}
	if ip2, ok := d.GetOk("ipaddress2"); ok {
		ipAddresses = append(ipAddresses, ip2.(string))
	}
	return ipAddresses
}

func flattenLoadBalancerIPAddresses(lb *sacloud.LoadBalancer) (ha bool, ipaddress1, ipaddress2 string) {
	ipaddress1 = lb.IPAddresses[0]
	if len(lb.IPAddresses) > 1 {
		ha = true
		ipaddress2 = lb.IPAddresses[1]
	}
	return
}

func expandLoadBalancerCreateRequest(d *schema.ResourceData) *sacloud.LoadBalancerCreateRequest {
	return &sacloud.LoadBalancerCreateRequest{
		SwitchID:           expandSakuraCloudID(d, "switch_id"),
		PlanID:             expandLoadBalancerPlanID(d),
		VRID:               d.Get("vrid").(int),
		IPAddresses:        expandLoadBalancerIPAddresses(d),
		NetworkMaskLen:     d.Get("nw_mask_len").(int),
		DefaultRoute:       d.Get("gateway").(string),
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Tags:               expandTags(d),
		IconID:             expandSakuraCloudID(d, "icon_id"),
		VirtualIPAddresses: expandLoadBalancerVIPs(d),
	}
}
func expandLoadBalancerUpdateRequest(d *schema.ResourceData, lb *sacloud.LoadBalancer) *sacloud.LoadBalancerUpdateRequest {
	return &sacloud.LoadBalancerUpdateRequest{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Tags:               expandTags(d),
		IconID:             expandSakuraCloudID(d, "icon_id"),
		VirtualIPAddresses: expandLoadBalancerVIPs(d),
		SettingsHash:       lb.SettingsHash,
	}
}
