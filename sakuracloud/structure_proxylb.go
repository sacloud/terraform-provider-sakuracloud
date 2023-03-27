// Copyright 2016-2023 terraform-provider-sakuracloud authors
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

func expandProxyLBCreateRequest(d *schema.ResourceData) *iaas.ProxyLBCreateRequest {
	return &iaas.ProxyLBCreateRequest{
		Plan:                 types.EProxyLBPlan(d.Get("plan").(int)),
		HealthCheck:          expandProxyLBHealthCheck(d),
		SorryServer:          expandProxyLBSorryServer(d),
		BindPorts:            expandProxyLBBindPorts(d),
		Servers:              expandProxyLBServers(d),
		Rules:                expandProxyLBRules(d),
		StickySession:        expandProxyLBStickySession(d),
		Gzip:                 expandProxyLBGzip(d),
		BackendHttpKeepAlive: expandProxyLBBackendHttpKeepAlive(d),
		ProxyProtocol:        expandProxyLBProxyProtocol(d),
		Syslog:               expandProxyLBSyslog(d),
		Timeout:              expandProxyLBTimeout(d),
		UseVIPFailover:       d.Get("vip_failover").(bool),
		Region:               types.EProxyLBRegion(d.Get("region").(string)),
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		Tags:                 expandTags(d),
		IconID:               expandSakuraCloudID(d, "icon_id"),
	}
}

func expandProxyLBUpdateRequest(d *schema.ResourceData) *iaas.ProxyLBUpdateRequest {
	return &iaas.ProxyLBUpdateRequest{
		HealthCheck:          expandProxyLBHealthCheck(d),
		SorryServer:          expandProxyLBSorryServer(d),
		BindPorts:            expandProxyLBBindPorts(d),
		Servers:              expandProxyLBServers(d),
		Rules:                expandProxyLBRules(d),
		StickySession:        expandProxyLBStickySession(d),
		Gzip:                 expandProxyLBGzip(d),
		BackendHttpKeepAlive: expandProxyLBBackendHttpKeepAlive(d),
		ProxyProtocol:        expandProxyLBProxyProtocol(d),
		Syslog:               expandProxyLBSyslog(d),
		Timeout:              expandProxyLBTimeout(d),
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		Tags:                 expandTags(d),
		IconID:               expandSakuraCloudID(d, "icon_id"),
	}
}

func flattenProxyLBSyslog(proxyLB *iaas.ProxyLB) []interface{} {
	syslog := proxyLB.Syslog
	if syslog != nil {
		return []interface{}{
			map[string]interface{}{
				"server": syslog.Server,
				"port":   syslog.Port,
			},
		}
	}
	return nil
}

func flattenProxyLBBindPorts(proxyLB *iaas.ProxyLB) []interface{} {
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
			"ssl_policy":        bindPort.SSLPolicy,
			"response_header":   headers,
		})
	}
	return bindPorts
}

func flattenProxyLBHealthCheck(proxyLB *iaas.ProxyLB) []interface{} {
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

func flattenProxyLBSorryServer(proxyLB *iaas.ProxyLB) []interface{} {
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

func flattenProxyLBServers(proxyLB *iaas.ProxyLB) []interface{} {
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

func flattenProxyLBRules(proxyLB *iaas.ProxyLB) []interface{} {
	var results []interface{}
	for _, rule := range proxyLB.Rules {
		results = append(results, map[string]interface{}{
			"host":                 rule.Host,
			"path":                 rule.Path,
			"source_ips":           rule.SourceIPs,
			"group":                rule.ServerGroup,
			"action":               rule.Action.String(),
			"redirect_location":    rule.RedirectLocation,
			"redirect_status_code": rule.RedirectStatusCode.String(),
			"fixed_status_code":    rule.FixedStatusCode.String(),
			"fixed_content_type":   rule.FixedContentType.String(),
			"fixed_message_body":   rule.FixedMessageBody,
		})
	}
	return results
}

func flattenProxyLBCerts(certs *iaas.ProxyLBCertificates) []interface{} {
	if certs == nil {
		return nil
	}
	proxylbCert := make(map[string]interface{})
	if certs.PrimaryCert != nil {
		proxylbCert["server_cert"] = certs.PrimaryCert.ServerCertificate
		proxylbCert["intermediate_cert"] = certs.PrimaryCert.IntermediateCertificate
		proxylbCert["private_key"] = certs.PrimaryCert.PrivateKey
		proxylbCert["common_name"] = certs.PrimaryCert.CertificateCommonName
		proxylbCert["subject_alt_names"] = certs.PrimaryCert.CertificateAltNames
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

func flattenProxyLBStickySession(proxyLB *iaas.ProxyLB) bool {
	if proxyLB.StickySession != nil {
		return proxyLB.StickySession.Enabled
	}
	return false
}

func flattenProxyLBGzip(proxyLB *iaas.ProxyLB) bool {
	if proxyLB.Gzip != nil {
		return proxyLB.Gzip.Enabled
	}
	return false
}

func flattenProxyLBBackendHttpKeepAlive(proxyLB *iaas.ProxyLB) string {
	if proxyLB.BackendHttpKeepAlive != nil {
		return proxyLB.BackendHttpKeepAlive.Mode.String()
	}
	return ""
}

func flattenProxyLBProxyProtocol(proxyLB *iaas.ProxyLB) bool {
	if proxyLB.ProxyProtocol != nil {
		return proxyLB.ProxyProtocol.Enabled
	}
	return false
}

func flattenProxyLBTimeout(proxyLB *iaas.ProxyLB) int {
	if proxyLB.Timeout != nil {
		return proxyLB.Timeout.InactiveSec
	}
	return 0
}

func expandProxyLBStickySession(d resourceValueGettable) *iaas.ProxyLBStickySession {
	stickySession := d.Get("sticky_session").(bool)
	if stickySession {
		return &iaas.ProxyLBStickySession{
			Enabled: true,
			Method:  "cookie",
		}
	}
	return nil
}

func expandProxyLBGzip(d resourceValueGettable) *iaas.ProxyLBGzip {
	gzip := d.Get("gzip").(bool)
	if gzip {
		return &iaas.ProxyLBGzip{
			Enabled: true,
		}
	}
	return nil
}

func expandProxyLBBackendHttpKeepAlive(d resourceValueGettable) *iaas.ProxyLBBackendHttpKeepAlive {
	s := d.Get("backend_http_keep_alive").(string)
	if s != "" {
		return &iaas.ProxyLBBackendHttpKeepAlive{
			Mode: types.EProxyLBBackendHttpKeepAlive(s),
		}
	}
	return nil
}

func expandProxyLBProxyProtocol(d resourceValueGettable) *iaas.ProxyLBProxyProtocol {
	v := d.Get("proxy_protocol").(bool)
	if v {
		return &iaas.ProxyLBProxyProtocol{
			Enabled: true,
		}
	}
	return nil
}

func expandProxyLBSyslog(d resourceValueGettable) *iaas.ProxyLBSyslog {
	if syslog, ok := getListFromResource(d, "syslog"); ok && len(syslog) == 1 {
		values := mapToResourceData(syslog[0].(map[string]interface{}))
		return &iaas.ProxyLBSyslog{
			Server: values.Get("server").(string),
			Port:   values.Get("port").(int),
		}
	}
	return &iaas.ProxyLBSyslog{Port: 514}
}

func expandProxyLBBindPorts(d resourceValueGettable) []*iaas.ProxyLBBindPort {
	var results []*iaas.ProxyLBBindPort
	if bindPorts, ok := getListFromResource(d, "bind_port"); ok {
		for _, bindPort := range bindPorts {
			values := mapToResourceData(bindPort.(map[string]interface{}))
			var headers []*iaas.ProxyLBResponseHeader
			if rawHeaders, ok := values.GetOk("response_header"); ok {
				for _, rawHeader := range rawHeaders.([]interface{}) {
					if rawHeader == nil {
						continue
					}
					v := rawHeader.(map[string]interface{})
					headers = append(headers, &iaas.ProxyLBResponseHeader{
						Header: v["header"].(string),
						Value:  v["value"].(string),
					})
				}
			}

			results = append(results, &iaas.ProxyLBBindPort{
				ProxyMode:         types.EProxyLBProxyMode(values.Get("proxy_mode").(string)),
				Port:              values.Get("port").(int),
				RedirectToHTTPS:   values.Get("redirect_to_https").(bool),
				SupportHTTP2:      values.Get("support_http2").(bool),
				SSLPolicy:         values.Get("ssl_policy").(string),
				AddResponseHeader: headers,
			})
		}
	}
	return results
}

func expandProxyLBHealthCheck(d resourceValueGettable) *iaas.ProxyLBHealthCheck {
	if healthChecks, ok := getListFromResource(d, "health_check"); ok {
		v := mapToResourceData(healthChecks[0].(map[string]interface{}))
		protocol := v.Get("protocol").(string)
		switch protocol {
		case "http":
			return &iaas.ProxyLBHealthCheck{
				Protocol:  types.ProxyLBProtocols.HTTP,
				Path:      v.Get("path").(string),
				Host:      v.Get("host_header").(string),
				DelayLoop: v.Get("delay_loop").(int),
			}
		case "tcp":
			return &iaas.ProxyLBHealthCheck{
				Protocol:  types.ProxyLBProtocols.TCP,
				DelayLoop: v.Get("delay_loop").(int),
			}
		}
	}
	return nil
}

func expandProxyLBSorryServer(d resourceValueGettable) *iaas.ProxyLBSorryServer {
	if sorryServers, ok := getListFromResource(d, "sorry_server"); ok && len(sorryServers) > 0 {
		v := mapToResourceData(sorryServers[0].(map[string]interface{}))
		return &iaas.ProxyLBSorryServer{
			IPAddress: v.Get("ip_address").(string),
			Port:      v.Get("port").(int),
		}
	}
	return nil
}

func expandProxyLBServers(d resourceValueGettable) []*iaas.ProxyLBServer {
	var results []*iaas.ProxyLBServer
	if servers, ok := getListFromResource(d, "server"); ok && len(servers) > 0 {
		for _, server := range servers {
			v := mapToResourceData(server.(map[string]interface{}))
			results = append(results, &iaas.ProxyLBServer{
				IPAddress:   v.Get("ip_address").(string),
				Port:        v.Get("port").(int),
				Enabled:     v.Get("enabled").(bool),
				ServerGroup: v.Get("group").(string),
			})
		}
	}
	return results
}

func expandProxyLBRules(d resourceValueGettable) []*iaas.ProxyLBRule {
	var results []*iaas.ProxyLBRule
	if rules, ok := getListFromResource(d, "rule"); ok && len(rules) > 0 {
		for _, rule := range rules {
			v := mapToResourceData(rule.(map[string]interface{}))
			results = append(results, &iaas.ProxyLBRule{
				Host:               v.Get("host").(string),
				Path:               v.Get("path").(string),
				SourceIPs:          v.Get("source_ips").(string),
				ServerGroup:        v.Get("group").(string),
				Action:             types.EProxyLBRuleAction(v.Get("action").(string)),
				RedirectLocation:   v.Get("redirect_location").(string),
				RedirectStatusCode: types.EProxyLBRedirectStatusCode(forceAtoI(v.Get("redirect_status_code").(string))),
				FixedStatusCode:    types.EProxyLBFixedStatusCode(forceAtoI(v.Get("fixed_status_code").(string))),
				FixedContentType:   types.EProxyLBFixedContentType(v.Get("fixed_content_type").(string)),
				FixedMessageBody:   v.Get("fixed_message_body").(string),
			})
		}
	}
	return results
}

func expandProxyLBTimeout(d resourceValueGettable) *iaas.ProxyLBTimeout {
	return &iaas.ProxyLBTimeout{InactiveSec: d.Get("timeout").(int)}
}

func expandProxyLBCerts(d resourceValueGettable) *iaas.ProxyLBCertificates {
	// set cert
	if certs, ok := getListFromResource(d, "certificate"); ok && len(certs) > 0 {
		values := mapToResourceData(certs[0].(map[string]interface{}))
		cert := &iaas.ProxyLBCertificates{
			PrimaryCert: &iaas.ProxyLBPrimaryCert{
				ServerCertificate:       values.Get("server_cert").(string),
				IntermediateCertificate: values.Get("intermediate_cert").(string),
				PrivateKey:              values.Get("private_key").(string),
			},
		}

		if rawAdditionalCerts, ok := getListFromResource(values, "additional_certificate"); ok && len(rawAdditionalCerts) > 0 {
			for _, rawCert := range rawAdditionalCerts {
				values := mapToResourceData(rawCert.(map[string]interface{}))
				cert.AdditionalCerts = append(cert.AdditionalCerts, &iaas.ProxyLBAdditionalCert{
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
