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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
)

func expandGSLBHealthCheckConf(d resourceValueGettable) *iaas.GSLBHealthCheck {
	healthCheckConf := d.Get("health_check").([]interface{})
	if len(healthCheckConf) == 0 {
		return nil
	}

	conf := healthCheckConf[0].(map[string]interface{})
	protocol := conf["protocol"].(string)
	switch protocol {
	case "http", "https":
		return &iaas.GSLBHealthCheck{
			Protocol:     types.EGSLBHealthCheckProtocol(protocol),
			HostHeader:   conf["host_header"].(string),
			Port:         types.StringNumber(conf["port"].(int)),
			Path:         conf["path"].(string),
			ResponseCode: types.StringNumber(forceAtoI(conf["status"].(string))),
		}
	case "tcp":
		return &iaas.GSLBHealthCheck{
			Protocol: types.EGSLBHealthCheckProtocol(protocol),
			Port:     types.StringNumber(conf["port"].(int)),
		}
	case "ping":
		return &iaas.GSLBHealthCheck{
			Protocol: types.EGSLBHealthCheckProtocol(protocol),
		}
	}
	return nil
}

func expandGSLBDelayLoop(d resourceValueGettable) int {
	healthCheckConf := d.Get("health_check").([]interface{})
	if len(healthCheckConf) == 0 {
		return 0
	}

	conf := healthCheckConf[0].(map[string]interface{})
	return conf["delay_loop"].(int)
}

func expandGSLBServers(d resourceValueGettable) []*iaas.GSLBServer {
	var servers []*iaas.GSLBServer
	for _, s := range d.Get("server").([]interface{}) {
		v := s.(map[string]interface{})
		server := expandGSLBServer(&resourceMapValue{value: v})
		servers = append(servers, server)
	}
	return servers
}

func flattenGSLBHealthCheck(data *iaas.GSLB) []interface{} {
	//health_check
	healthCheck := map[string]interface{}{}
	switch data.HealthCheck.Protocol {
	case types.GSLBHealthCheckProtocols.HTTP, types.GSLBHealthCheckProtocols.HTTPS:
		healthCheck["host_header"] = data.HealthCheck.HostHeader
		healthCheck["port"] = data.HealthCheck.Port
		healthCheck["path"] = data.HealthCheck.Path
		healthCheck["status"] = data.HealthCheck.ResponseCode.String()
	case types.GSLBHealthCheckProtocols.TCP:
		healthCheck["port"] = data.HealthCheck.Port
	}
	healthCheck["protocol"] = data.HealthCheck.Protocol
	healthCheck["delay_loop"] = data.DelayLoop

	return []interface{}{healthCheck}
}

func flattenGSLBServers(data *iaas.GSLB) []interface{} {
	var servers []interface{}
	for _, server := range data.DestinationServers {
		servers = append(servers, flattenGSLBServer(server))
	}
	return servers
}

func flattenGSLBServer(s *iaas.GSLBServer) interface{} {
	v := map[string]interface{}{}
	v["ip_address"] = s.IPAddress
	v["enabled"] = s.Enabled.Bool()
	v["weight"] = s.Weight.Int()
	return v
}

func expandGSLBServer(d resourceValueGettable) *iaas.GSLBServer {
	return &iaas.GSLBServer{
		IPAddress: d.Get("ip_address").(string),
		Enabled:   types.StringFlag(d.Get("enabled").(bool)),
		Weight:    types.StringNumber(d.Get("weight").(int)),
	}
}

func expandGSLBCreateRequest(d *schema.ResourceData) *iaas.GSLBCreateRequest {
	return &iaas.GSLBCreateRequest{
		HealthCheck:        expandGSLBHealthCheckConf(d),
		DelayLoop:          expandGSLBDelayLoop(d),
		Weighted:           types.StringFlag(d.Get("weighted").(bool)),
		SorryServer:        d.Get("sorry_server").(string),
		DestinationServers: expandGSLBServers(d),
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Tags:               expandTags(d),
		IconID:             expandSakuraCloudID(d, "icon_id"),
	}
}

func expandGSLBUpdateRequest(d *schema.ResourceData, gslb *iaas.GSLB) *iaas.GSLBUpdateRequest {
	return &iaas.GSLBUpdateRequest{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Tags:               expandTags(d),
		IconID:             expandSakuraCloudID(d, "icon_id"),
		HealthCheck:        expandGSLBHealthCheckConf(d),
		DelayLoop:          expandGSLBDelayLoop(d),
		Weighted:           types.StringFlag(d.Get("weighted").(bool)),
		SorryServer:        d.Get("sorry_server").(string),
		DestinationServers: expandGSLBServers(d),
		SettingsHash:       gslb.SettingsHash,
	}
}
