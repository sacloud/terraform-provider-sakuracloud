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
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func expandSimpleMonitorCreateRequest(d *schema.ResourceData) *sacloud.SimpleMonitorCreateRequest {
	return &sacloud.SimpleMonitorCreateRequest{
		Target:             d.Get("target").(string),
		Enabled:            types.StringFlag(d.Get("enabled").(bool)),
		HealthCheck:        expandSimpleMonitorHealthCheck(d),
		DelayLoop:          d.Get("delay_loop").(int),
		MaxCheckAttempts:   d.Get("max_check_attempts").(int),
		RetryInterval:      d.Get("retry_interval").(int),
		Timeout:            d.Get("timeout").(int),
		NotifyEmailEnabled: types.StringFlag(d.Get("notify_email_enabled").(bool)),
		NotifyEmailHTML:    types.StringFlag(d.Get("notify_email_html").(bool)),
		NotifySlackEnabled: types.StringFlag(d.Get("notify_slack_enabled").(bool)),
		SlackWebhooksURL:   d.Get("notify_slack_webhook").(string),
		NotifyInterval:     expandSimpleMonitorNotifyInterval(d),
		Description:        d.Get("description").(string),
		Tags:               expandTags(d),
		IconID:             expandSakuraCloudID(d, "icon_id"),
	}
}

func expandSimpleMonitorUpdateRequest(d *schema.ResourceData) *sacloud.SimpleMonitorUpdateRequest {
	return &sacloud.SimpleMonitorUpdateRequest{
		Enabled:            types.StringFlag(d.Get("enabled").(bool)),
		HealthCheck:        expandSimpleMonitorHealthCheck(d),
		DelayLoop:          d.Get("delay_loop").(int),
		MaxCheckAttempts:   d.Get("max_check_attempts").(int),
		RetryInterval:      d.Get("retry_interval").(int),
		Timeout:            d.Get("timeout").(int),
		NotifyEmailEnabled: types.StringFlag(d.Get("notify_email_enabled").(bool)),
		NotifyEmailHTML:    types.StringFlag(d.Get("notify_email_html").(bool)),
		NotifySlackEnabled: types.StringFlag(d.Get("notify_slack_enabled").(bool)),
		SlackWebhooksURL:   d.Get("notify_slack_webhook").(string),
		NotifyInterval:     expandSimpleMonitorNotifyInterval(d),
		Description:        d.Get("description").(string),
		Tags:               expandTags(d),
		IconID:             expandSakuraCloudID(d, "icon_id"),
	}
}

func expandSimpleMonitorNotifyInterval(d *schema.ResourceData) int {
	return d.Get("notify_interval").(int) * 60 * 60 // hours => seconds
}

func flattenSimpleMonitorNotifyInterval(simpleMonitor *sacloud.SimpleMonitor) int {
	interval := simpleMonitor.NotifyInterval
	if interval == 0 {
		return 0
	}
	// seconds => hours
	return interval / 60 / 60
}

func flattenSimpleMonitorHealthCheck(simpleMonitor *sacloud.SimpleMonitor) []interface{} {
	healthCheck := map[string]interface{}{}
	hc := simpleMonitor.HealthCheck
	switch hc.Protocol {
	case types.SimpleMonitorProtocols.HTTP:
		healthCheck["path"] = hc.Path
		healthCheck["status"] = hc.Status.Int()
		healthCheck["contains_string"] = hc.ContainsString
		healthCheck["host_header"] = hc.Host
		healthCheck["port"] = hc.Port.Int()
		healthCheck["username"] = hc.BasicAuthUsername
		healthCheck["password"] = hc.BasicAuthPassword
	case types.SimpleMonitorProtocols.HTTPS:
		healthCheck["path"] = hc.Path
		healthCheck["status"] = hc.Status.Int()
		healthCheck["contains_string"] = hc.ContainsString
		healthCheck["host_header"] = hc.Host
		healthCheck["port"] = hc.Port.Int()
		healthCheck["sni"] = hc.SNI.Bool()
		healthCheck["username"] = hc.BasicAuthUsername
		healthCheck["password"] = hc.BasicAuthPassword
		healthCheck["http2"] = hc.HTTP2
		healthCheck["verify_sni"] = hc.VerifySNI
	case types.SimpleMonitorProtocols.TCP, types.SimpleMonitorProtocols.SSH, types.SimpleMonitorProtocols.SMTP, types.SimpleMonitorProtocols.POP3:
		healthCheck["port"] = hc.Port.Int()
	case types.SimpleMonitorProtocols.SNMP:
		healthCheck["community"] = hc.Community
		healthCheck["snmp_version"] = hc.SNMPVersion
		healthCheck["oid"] = hc.OID
		healthCheck["expected_data"] = hc.ExpectedData
	case types.SimpleMonitorProtocols.DNS:
		healthCheck["qname"] = hc.QName
		healthCheck["expected_data"] = hc.ExpectedData
	case types.SimpleMonitorProtocols.FTP:
		healthCheck["ftps"] = hc.FTPS.String()
	case types.SimpleMonitorProtocols.SSLCertificate:
	}
	days := hc.RemainingDays
	if days == 0 {
		days = 30
	}
	healthCheck["remaining_days"] = days
	healthCheck["protocol"] = hc.Protocol
	return []interface{}{healthCheck}
}

func expandSimpleMonitorHealthCheck(d resourceValueGettable) *sacloud.SimpleMonitorHealthCheck {
	healthCheckConf := d.Get("health_check").([]interface{})
	conf := healthCheckConf[0].(map[string]interface{})
	protocol := conf["protocol"].(string)
	port := conf["port"].(int)

	switch protocol {
	case "http":
		if port == 0 {
			port = 80
		}
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol:          types.SimpleMonitorProtocols.HTTP,
			Port:              types.StringNumber(port),
			Path:              forceString(conf["path"]),
			Status:            types.StringNumber(conf["status"].(int)),
			ContainsString:    forceString(conf["contains_string"]),
			Host:              forceString(conf["host_header"]),
			BasicAuthUsername: forceString(conf["username"]),
			BasicAuthPassword: forceString(conf["password"]),
		}
	case "https":
		if port == 0 {
			port = 443
		}
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol:          types.SimpleMonitorProtocols.HTTPS,
			Port:              types.StringNumber(port),
			Path:              forceString(conf["path"]),
			Status:            types.StringNumber(conf["status"].(int)),
			ContainsString:    forceString(conf["contains_string"]),
			SNI:               types.StringFlag(forceBool(conf["sni"])),
			Host:              forceString(conf["host_header"]),
			BasicAuthUsername: forceString(conf["username"]),
			BasicAuthPassword: forceString(conf["password"]),
			HTTP2:             types.StringFlag(forceBool(conf["http2"])),
			VerifySNI:         types.StringFlag(forceBool(conf["verify_sni"])),
		}

	case "dns":
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol:     types.SimpleMonitorProtocols.DNS,
			QName:        forceString(conf["qname"]),
			ExpectedData: forceString(conf["expected_data"]),
		}
	case "snmp":
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol:     types.SimpleMonitorProtocols.SNMP,
			Community:    forceString(conf["community"]),
			SNMPVersion:  forceString(conf["snmp_version"]),
			OID:          forceString(conf["oid"]),
			ExpectedData: forceString(conf["expected_data"]),
		}
	case "tcp":
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol: types.SimpleMonitorProtocols.TCP,
			Port:     types.StringNumber(port),
		}
	case "ssh":
		if port == 0 {
			port = 22
		}
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol: types.SimpleMonitorProtocols.SSH,
			Port:     types.StringNumber(port),
		}
	case "smtp":
		if port == 0 {
			port = 25
		}
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol: types.SimpleMonitorProtocols.SMTP,
			Port:     types.StringNumber(port),
		}
	case "pop3":
		if port == 0 {
			port = 110
		}
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol: types.SimpleMonitorProtocols.POP3,
			Port:     types.StringNumber(port),
		}
	case "ping":
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol: types.SimpleMonitorProtocols.Ping,
		}
	case "sslcertificate":
		days := 0
		if v, ok := conf["remaining_days"]; ok {
			days = v.(int)
		}
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol:      types.SimpleMonitorProtocols.SSLCertificate,
			RemainingDays: days,
		}
	case "ftp":
		if port == 0 {
			port = 21
		}
		ftps := ""
		if v, ok := conf["ftps"]; ok {
			ftps = v.(string)
		}
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol: types.SimpleMonitorProtocols.FTP,
			Port:     types.StringNumber(port),
			FTPS:     types.ESimpleMonitorFTPS(ftps),
		}
	}
	return nil
}
