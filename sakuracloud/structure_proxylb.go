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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func expandProxyLBCreateRequest(d *schema.ResourceData) *sacloud.ProxyLBCreateRequest {
	return &sacloud.ProxyLBCreateRequest{
		Plan:           types.EProxyLBPlan(d.Get("plan").(int)),
		HealthCheck:    expandProxyLBHealthCheck(d),
		SorryServer:    expandProxyLBSorryServer(d),
		BindPorts:      expandProxyLBBindPorts(d),
		Servers:        expandProxyLBServers(d),
		Rules:          expandProxyLBRules(d),
		StickySession:  expandProxyLBStickySession(d),
		Gzip:           expandProxyLBGzip(d),
		Timeout:        expandProxyLBTimeout(d),
		UseVIPFailover: d.Get("vip_failover").(bool),
		Region:         types.EProxyLBRegion(d.Get("region").(string)),
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Tags:           expandTags(d),
		IconID:         expandSakuraCloudID(d, "icon_id"),
	}
}

func expandProxyLBUpdateRequest(d *schema.ResourceData) *sacloud.ProxyLBUpdateRequest {
	return &sacloud.ProxyLBUpdateRequest{
		HealthCheck:   expandProxyLBHealthCheck(d),
		SorryServer:   expandProxyLBSorryServer(d),
		BindPorts:     expandProxyLBBindPorts(d),
		Servers:       expandProxyLBServers(d),
		Rules:         expandProxyLBRules(d),
		StickySession: expandProxyLBStickySession(d),
		Gzip:          expandProxyLBGzip(d),
		Timeout:       expandProxyLBTimeout(d),
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Tags:          expandTags(d),
		IconID:        expandSakuraCloudID(d, "icon_id"),
	}
}

func flattenProxyLBBindPorts(proxyLB *sacloud.ProxyLB) []interface{} {
	var bindPorts []interface{}
	for _, bindPort := range proxyLB.BindPorts {
		var headers []interface{}
		for _, header := range bindPort.AddResponseHeader {
			headers = append(headers, map[string]interface{}{
				"header": header.Header,
				"value":  header.Value,
			})
		}

		bindPorts = append(bindPorts, map[string]interface{}{
			"proxy_mode":        bindPort.ProxyMode,
			"port":              bindPort.Port,
			"redirect_to_https": bindPort.RedirectToHTTPS,
			"support_http2":     bindPort.SupportHTTP2,
			"response_header":   headers,
		})
	}
	return bindPorts
}

func flattenProxyLBHealthCheck(proxyLB *sacloud.ProxyLB) []interface{} {
	var results []interface{}
	if proxyLB.HealthCheck != nil {
		results = []interface{}{
			map[string]interface{}{
				"protocol":    proxyLB.HealthCheck.Protocol,
				"delay_loop":  proxyLB.HealthCheck.DelayLoop,
				"host_header": proxyLB.HealthCheck.Host,
				"path":        proxyLB.HealthCheck.Path,
			},
		}
	}
	return results
}

func flattenProxyLBSorryServer(proxyLB *sacloud.ProxyLB) []interface{} {
	var results []interface{}
	if proxyLB.SorryServer != nil && proxyLB.SorryServer.IPAddress != "" {
		results = []interface{}{
			map[string]interface{}{
				"ip_address": proxyLB.SorryServer.IPAddress,
				"port":       proxyLB.SorryServer.Port,
			},
		}
	}
	return results
}

func flattenProxyLBServers(proxyLB *sacloud.ProxyLB) []interface{} {
	var results []interface{}
	for _, server := range proxyLB.Servers {
		results = append(results, map[string]interface{}{
			"ip_address": server.IPAddress,
			"port":       server.Port,
			"enabled":    server.Enabled,
			"group":      server.ServerGroup,
		})
	}
	return results
}

func flattenProxyLBRules(proxyLB *sacloud.ProxyLB) []interface{} {
	var results []interface{}
	for _, rule := range proxyLB.Rules {
		results = append(results, map[string]interface{}{
			"host":  rule.Host,
			"path":  rule.Path,
			"group": rule.ServerGroup,
		})
	}
	return results
}

func flattenProxyLBCerts(certs *sacloud.ProxyLBCertificates) []interface{} {
	if certs == nil {
		return nil
	}
	proxylbCert := make(map[string]interface{})
	if certs.PrimaryCert != nil {
		proxylbCert["server_cert"] = certs.PrimaryCert.ServerCertificate
		proxylbCert["intermediate_cert"] = certs.PrimaryCert.IntermediateCertificate
		proxylbCert["private_key"] = certs.PrimaryCert.PrivateKey
	}
	if len(certs.AdditionalCerts) > 0 {
		var additionalCerts []interface{}
		for _, cert := range certs.AdditionalCerts {
			additionalCerts = append(additionalCerts, map[string]interface{}{
				"server_cert":       cert.ServerCertificate,
				"intermediate_cert": cert.IntermediateCertificate,
				"private_key":       cert.PrivateKey,
			})
		}
		proxylbCert["additional_certificate"] = additionalCerts
	}
	return []interface{}{proxylbCert}
}

func flattenProxyLBStickySession(proxyLB *sacloud.ProxyLB) bool {
	if proxyLB.StickySession != nil {
		return proxyLB.StickySession.Enabled
	}
	return false
}

func flattenProxyLBGzip(proxyLB *sacloud.ProxyLB) bool {
	if proxyLB.Gzip != nil {
		return proxyLB.Gzip.Enabled
	}
	return false
}

func flattenProxyLBTimeout(proxyLB *sacloud.ProxyLB) int {
	if proxyLB.Timeout != nil {
		return proxyLB.Timeout.InactiveSec
	}
	return 0
}

func expandProxyLBStickySession(d resourceValueGettable) *sacloud.ProxyLBStickySession {
	stickySession := d.Get("sticky_session").(bool)
	if stickySession {
		return &sacloud.ProxyLBStickySession{
			Enabled: true,
			Method:  "cookie",
		}
	}
	return nil
}

func expandProxyLBGzip(d resourceValueGettable) *sacloud.ProxyLBGzip {
	gzip := d.Get("gzip").(bool)
	if gzip {
		return &sacloud.ProxyLBGzip{
			Enabled: true,
		}
	}
	return nil
}

func expandProxyLBBindPorts(d resourceValueGettable) []*sacloud.ProxyLBBindPort {
	var results []*sacloud.ProxyLBBindPort
	if bindPorts, ok := getListFromResource(d, "bind_port"); ok {
		for _, bindPort := range bindPorts {
			values := mapToResourceData(bindPort.(map[string]interface{}))
			var headers []*sacloud.ProxyLBResponseHeader
			if rawHeaders, ok := values.GetOk("response_header"); ok {
				for _, rawHeader := range rawHeaders.([]interface{}) {
					if rawHeader == nil {
						continue
					}
					v := rawHeader.(map[string]interface{})
					headers = append(headers, &sacloud.ProxyLBResponseHeader{
						Header: v["header"].(string),
						Value:  v["value"].(string),
					})
				}
			}

			results = append(results, &sacloud.ProxyLBBindPort{
				ProxyMode:         types.EProxyLBProxyMode(values.Get("proxy_mode").(string)),
				Port:              values.Get("port").(int),
				RedirectToHTTPS:   values.Get("redirect_to_https").(bool),
				SupportHTTP2:      values.Get("support_http2").(bool),
				AddResponseHeader: headers,
			})
		}
	}
	return results
}

func expandProxyLBHealthCheck(d resourceValueGettable) *sacloud.ProxyLBHealthCheck {
	if healthChecks, ok := getListFromResource(d, "health_check"); ok {
		v := mapToResourceData(healthChecks[0].(map[string]interface{}))
		protocol := v.Get("protocol").(string)
		switch protocol {
		case "http":
			return &sacloud.ProxyLBHealthCheck{
				Protocol:  types.ProxyLBProtocols.HTTP,
				Path:      v.Get("path").(string),
				Host:      v.Get("host_header").(string),
				DelayLoop: v.Get("delay_loop").(int),
			}
		case "tcp":
			return &sacloud.ProxyLBHealthCheck{
				Protocol:  types.ProxyLBProtocols.TCP,
				DelayLoop: v.Get("delay_loop").(int),
			}
		}
	}
	return nil
}

func expandProxyLBSorryServer(d resourceValueGettable) *sacloud.ProxyLBSorryServer {
	if sorryServers, ok := getListFromResource(d, "sorry_server"); ok && len(sorryServers) > 0 {
		v := mapToResourceData(sorryServers[0].(map[string]interface{}))
		return &sacloud.ProxyLBSorryServer{
			IPAddress: v.Get("ip_address").(string),
			Port:      v.Get("port").(int),
		}
	}
	return nil
}

func expandProxyLBServers(d resourceValueGettable) []*sacloud.ProxyLBServer {
	var results []*sacloud.ProxyLBServer
	if servers, ok := getListFromResource(d, "server"); ok && len(servers) > 0 {
		for _, server := range servers {
			v := mapToResourceData(server.(map[string]interface{}))
			results = append(results, &sacloud.ProxyLBServer{
				IPAddress:   v.Get("ip_address").(string),
				Port:        v.Get("port").(int),
				Enabled:     v.Get("enabled").(bool),
				ServerGroup: v.Get("group").(string),
			})
		}
	}
	return results
}

func expandProxyLBRules(d resourceValueGettable) []*sacloud.ProxyLBRule {
	var results []*sacloud.ProxyLBRule
	if rules, ok := getListFromResource(d, "rule"); ok && len(rules) > 0 {
		for _, rule := range rules {
			v := mapToResourceData(rule.(map[string]interface{}))
			results = append(results, &sacloud.ProxyLBRule{
				Host:        v.Get("host").(string),
				Path:        v.Get("path").(string),
				ServerGroup: v.Get("group").(string),
			})
		}
	}
	return results
}

func expandProxyLBTimeout(d resourceValueGettable) *sacloud.ProxyLBTimeout {
	return &sacloud.ProxyLBTimeout{InactiveSec: d.Get("timeout").(int)}
}

func expandProxyLBCerts(d resourceValueGettable) *sacloud.ProxyLBCertificates {
	// set cert
	if certs, ok := getListFromResource(d, "certificate"); ok && len(certs) > 0 {
		values := mapToResourceData(certs[0].(map[string]interface{}))
		cert := &sacloud.ProxyLBCertificates{
			PrimaryCert: &sacloud.ProxyLBPrimaryCert{
				ServerCertificate:       values.Get("server_cert").(string),
				IntermediateCertificate: values.Get("intermediate_cert").(string),
				PrivateKey:              values.Get("private_key").(string),
			},
		}

		if rawAdditionalCerts, ok := getListFromResource(values, "additional_certificate"); ok && len(rawAdditionalCerts) > 0 {
			for _, rawCert := range rawAdditionalCerts {
				values := mapToResourceData(rawCert.(map[string]interface{}))
				cert.AdditionalCerts = append(cert.AdditionalCerts, &sacloud.ProxyLBAdditionalCert{
					ServerCertificate:       values.Get("server_cert").(string),
					IntermediateCertificate: values.Get("intermediate_cert").(string),
					PrivateKey:              values.Get("private_key").(string),
				})
			}
		}

		return cert
	}
	return nil
}
