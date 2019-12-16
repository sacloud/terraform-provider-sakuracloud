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
	"fmt"
	"strconv"
	"strings"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/builder"
	"github.com/sacloud/libsacloud/v2/utils/builder/vpcrouter"
)

func expandVPCRouterBuilder(d resourceValueGettable, client *APIClient) *vpcrouter.Builder {
	return &vpcrouter.Builder{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Tags:                  expandTags(d),
		IconID:                expandSakuraCloudID(d, "icon_id"),
		PlanID:                expandVPCRouterPlanID(d),
		NICSetting:            expandVPCRouterNICSetting(d),
		AdditionalNICSettings: expandVPCRouterAdditionalNICSettings(d),
		RouterSetting:         expandVPCRouterSettings(d),
		SetupOptions: &builder.RetryableSetupParameter{
			BootAfterBuild:        true,
			NICUpdateWaitDuration: builder.DefaultNICUpdateWaitDuration,
		},
		Client: sacloud.NewVPCRouterOp(client),
	}
}

func expandVPCRouterPlanID(d resourceValueGettable) types.ID {
	plan := d.Get("plan").(string)
	switch plan {
	case "standard":
		return types.VPCRouterPlans.Standard
	case "premium":
		return types.VPCRouterPlans.Premium
	case "highspec":
		return types.VPCRouterPlans.HighSpec
	case "highspec4000":
		return types.VPCRouterPlans.HighSpec4000
	default:
		return types.VPCRouterPlans.Standard
	}
}

func flattenVPCRouterPlan(vpcRouter *sacloud.VPCRouter) string {
	switch vpcRouter.PlanID {
	case types.VPCRouterPlans.Standard:
		return "standard"
	case types.VPCRouterPlans.Premium:
		return "premium"
	case types.VPCRouterPlans.HighSpec:
		return "highspec"
	case types.VPCRouterPlans.HighSpec4000:
		return "highspec4000"
	default:
		return "standard"
	}
}

func expandVPCRouterNICSetting(d resourceValueGettable) vpcrouter.NICSettingHolder {
	planID := expandVPCRouterPlanID(d)
	switch planID {
	case types.VPCRouterPlans.Standard:
		return &vpcrouter.StandardNICSetting{}
	default:
		return &vpcrouter.PremiumNICSetting{
			SwitchID:         expandSakuraCloudID(d, "switch_id"),
			IPAddress1:       d.Get("ipaddress1").(string),
			IPAddress2:       d.Get("ipaddress2").(string),
			VirtualIPAddress: d.Get("vip").(string),
			IPAliases:        expandStringList(d.Get("aliases").([]interface{})),
		}
	}
}

func expandVPCRouterAdditionalNICSettings(d resourceValueGettable) []vpcrouter.AdditionalNICSettingHolder {
	var results []vpcrouter.AdditionalNICSettingHolder
	planID := expandVPCRouterPlanID(d)
	interfaces := d.Get("network_interface").([]interface{})
	for _, iface := range interfaces {
		d = mapToResourceData(iface.(map[string]interface{}))
		var nicSetting vpcrouter.AdditionalNICSettingHolder
		ipAddresses := expandStringList(d.Get("ip_addresses").([]interface{}))

		switch planID {
		case types.VPCRouterPlans.Standard:
			nicSetting = &vpcrouter.AdditionalStandardNICSetting{
				SwitchID:       expandSakuraCloudID(d, "switch_id"),
				IPAddress:      ipAddresses[0],
				NetworkMaskLen: d.Get("nw_mask_len").(int),
				Index:          d.Get("index").(int),
			}
		default:
			nicSetting = &vpcrouter.AdditionalPremiumNICSetting{
				SwitchID:         expandSakuraCloudID(d, "switch_id"),
				NetworkMaskLen:   d.Get("nw_mask_len").(int),
				IPAddress1:       ipAddresses[0],
				IPAddress2:       ipAddresses[1],
				VirtualIPAddress: d.Get("vip").(string),
				Index:            d.Get("index").(int),
			}
		}
		results = append(results, nicSetting)
	}
	return results
}

func flattenVPCRouterInterfaces(vpcRouter *sacloud.VPCRouter) []interface{} {
	var interfaces []interface{}
	if len(vpcRouter.Interfaces) > 0 {
		for _, iface := range vpcRouter.Settings.Interfaces {
			if iface.Index == 0 {
				continue
			}
			// find nic from data.Interfaces
			var nic *sacloud.VPCRouterInterface
			for _, n := range vpcRouter.Interfaces {
				if iface.Index == n.Index {
					nic = n
					break
				}
			}

			if nic != nil {
				interfaces = append(interfaces, map[string]interface{}{
					"switch_id":    nic.SwitchID.String(),
					"vip":          iface.VirtualIPAddress,
					"ip_addresses": iface.IPAddress,
					"nw_mask_len":  iface.NetworkMaskLen,
					"index":        iface.Index,
				})
			}
		}
	}
	return interfaces
}

func flattenVPCRouterGlobalAddress(vpcRouter *sacloud.VPCRouter) string {
	if vpcRouter.PlanID == types.VPCRouterPlans.Standard {
		return vpcRouter.Interfaces[0].IPAddress
	}
	return vpcRouter.Settings.Interfaces[0].VirtualIPAddress
}

func flattenVPCRouterSwitchID(vpcRouter *sacloud.VPCRouter) string {
	if vpcRouter.PlanID != types.VPCRouterPlans.Standard {
		return vpcRouter.Interfaces[0].SwitchID.String()
	}
	return ""
}

func flattenVPCRouterVIP(vpcRouter *sacloud.VPCRouter) string {
	if vpcRouter.PlanID != types.VPCRouterPlans.Standard {
		return vpcRouter.Settings.Interfaces[0].VirtualIPAddress
	}
	return ""
}

func flattenVPCRouterIPAddress1(vpcRouter *sacloud.VPCRouter) string {
	if vpcRouter.PlanID != types.VPCRouterPlans.Standard {
		return vpcRouter.Settings.Interfaces[0].IPAddress[0]
	}
	return ""
}

func flattenVPCRouterIPAddress2(vpcRouter *sacloud.VPCRouter) string {
	if vpcRouter.PlanID != types.VPCRouterPlans.Standard {
		return vpcRouter.Settings.Interfaces[0].IPAddress[1]
	}
	return ""
}

func flattenVPCRouterIPAliases(vpcRouter *sacloud.VPCRouter) []string {
	if vpcRouter.PlanID != types.VPCRouterPlans.Standard {
		return vpcRouter.Settings.Interfaces[0].IPAliases
	}
	return []string{}
}

func flattenVPCRouterVRID(vpcRouter *sacloud.VPCRouter) int {
	if vpcRouter.PlanID != types.VPCRouterPlans.Standard {
		return vpcRouter.Settings.VRID
	}
	return 0
}

func expandVPCRouterSettings(d resourceValueGettable) *vpcrouter.RouterSetting {
	return &vpcrouter.RouterSetting{
		VRID:                      d.Get("vrid").(int),
		InternetConnectionEnabled: types.StringFlag(d.Get("internet_connection").(bool)),
		StaticNAT:                 expandVPCRouterStaticNATList(d),
		PortForwarding:            expandVPCRouterPortForwardingList(d),
		Firewall:                  expandVPCRouterFirewallList(d),
		DHCPServer:                expandVPCRouterDHCPServerList(d),
		DHCPStaticMapping:         expandVPCRouterDHCPStaticMappingList(d),
		PPTPServer:                expandVPCRouterPPTP(d),
		L2TPIPsecServer:           expandVPCRouterL2TP(d),
		RemoteAccessUsers:         expandVPCRouterUserList(d),
		SiteToSiteIPsecVPN:        expandVPCRouterSiteToSiteList(d),
		StaticRoute:               expandVPCRouterStaticRouteList(d),
		SyslogHost:                d.Get("syslog_host").(string),
	}
}

func expandVPCRouterStaticNATList(d resourceValueGettable) []*sacloud.VPCRouterStaticNAT {
	if values, ok := getListFromResource(d, "static_nat"); ok && len(values) > 0 {
		var results []*sacloud.VPCRouterStaticNAT
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterStaticNAT(v))
		}
		return results
	}
	return nil
}

func expandVPCRouterStaticNAT(d resourceValueGettable) *sacloud.VPCRouterStaticNAT {
	return &sacloud.VPCRouterStaticNAT{
		GlobalAddress:  d.Get("public_ip").(string),
		PrivateAddress: d.Get("private_ip").(string),
		Description:    d.Get("description").(string),
	}
}

func flattenVPCRouterStaticNAT(vpcRouter *sacloud.VPCRouter) []interface{} {
	var staticNATs []interface{}
	for _, s := range vpcRouter.Settings.StaticNAT {
		staticNATs = append(staticNATs, map[string]interface{}{
			"public_ip":   s.GlobalAddress,
			"private_ip":  s.PrivateAddress,
			"description": s.Description,
		})
	}
	return staticNATs
}

func expandVPCRouterDHCPServerList(d resourceValueGettable) []*sacloud.VPCRouterDHCPServer {
	if values, ok := getListFromResource(d, "dhcp_server"); ok && len(values) > 0 {
		var results []*sacloud.VPCRouterDHCPServer
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterDHCPServer(v))
		}
		return results
	}
	return nil
}

func expandVPCRouterDHCPServer(d resourceValueGettable) *sacloud.VPCRouterDHCPServer {
	return &sacloud.VPCRouterDHCPServer{
		Interface:  fmt.Sprintf("eth%d", d.Get("interface_index").(int)),
		RangeStart: d.Get("range_start").(string),
		RangeStop:  d.Get("range_stop").(string),
		DNSServers: expandStringList(d.Get("dns_servers").([]interface{})),
	}
}

func flattenVPCRouterDHCPServers(vpcRouter *sacloud.VPCRouter) []interface{} {
	var dhcpServers []interface{}
	for _, d := range vpcRouter.Settings.DHCPServer {
		dhcpServers = append(dhcpServers, map[string]interface{}{
			"range_start":     d.RangeStart,
			"range_stop":      d.RangeStop,
			"interface_index": vpcRouterInterfaceNameToIndex(d.Interface),
			"dns_servers":     d.DNSServers,
		})
	}
	return dhcpServers
}

func vpcRouterInterfaceNameToIndex(ifName string) int {
	strIndex := strings.Replace(ifName, "eth", "", -1)
	index, err := strconv.Atoi(strIndex)
	if err != nil {
		return -1
	}
	return index
}

func expandVPCRouterDHCPStaticMappingList(d resourceValueGettable) []*sacloud.VPCRouterDHCPStaticMapping {
	if values, ok := getListFromResource(d, "dhcp_static_mapping"); ok && len(values) > 0 {
		var results []*sacloud.VPCRouterDHCPStaticMapping
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterDHCPStaticMapping(v))
		}
		return results
	}
	return nil
}

func expandVPCRouterDHCPStaticMapping(d resourceValueGettable) *sacloud.VPCRouterDHCPStaticMapping {
	return &sacloud.VPCRouterDHCPStaticMapping{
		IPAddress:  d.Get("ipaddress").(string),
		MACAddress: d.Get("macaddress").(string),
	}
}

func flattenVPCRouterDHCPStaticMappings(vpcRouter *sacloud.VPCRouter) []interface{} {
	var staticMappings []interface{}
	for _, d := range vpcRouter.Settings.DHCPStaticMapping {
		staticMappings = append(staticMappings, map[string]interface{}{
			"ipaddress":  d.IPAddress,
			"macaddress": d.MACAddress,
		})
	}
	return staticMappings
}

func expandVPCRouterFirewallList(d resourceValueGettable) []*sacloud.VPCRouterFirewall {
	if values, ok := getListFromResource(d, "firewall"); ok && len(values) > 0 {
		var results []*sacloud.VPCRouterFirewall
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterFirewall(v))
		}

		// インデックスごとにSend/Receiveをまとめる
		// results: {Index: 0, Send: []Rules{...}, Receive: nil} , {Index: 0, Send: nil, Receive: []Rules{...}}
		// merged: {Index: 0, Send: []Rules{...}, Receive: []Rules{...}}
		var merged []*sacloud.VPCRouterFirewall
		for i := 0; i < 8; i++ {
			firewall := &sacloud.VPCRouterFirewall{
				Index: i,
			}
			for _, f := range results {
				if f.Index == i {
					if len(f.Send) > 0 {
						firewall.Send = f.Send
					}
					if len(f.Receive) > 0 {
						firewall.Receive = f.Receive
					}
				}
			}
			merged = append(merged, firewall)
		}
		return merged
	}
	return nil
}

func expandVPCRouterFirewall(d resourceValueGettable) *sacloud.VPCRouterFirewall {
	index := intOrDefault(d, "interface_index")
	direction := stringOrDefault(d, "direction")
	f := &sacloud.VPCRouterFirewall{
		Index: index,
	}
	if direction == "send" {
		f.Send = expandVPCRouterFirewallRuleList(d)
	}
	if direction == "receive" {
		f.Receive = expandVPCRouterFirewallRuleList(d)
	}
	return f
}

func expandVPCRouterFirewallRuleList(d resourceValueGettable) []*sacloud.VPCRouterFirewallRule {
	if values, ok := getListFromResource(d, "expression"); ok && len(values) > 0 {
		var results []*sacloud.VPCRouterFirewallRule
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterFirewallRule(v))
		}
		return results
	}
	return nil
}

func expandVPCRouterFirewallRule(d resourceValueGettable) *sacloud.VPCRouterFirewallRule {
	allow := boolOrDefault(d, "allow")
	action := types.Actions.Allow
	if !allow {
		action = types.Actions.Deny
	}

	return &sacloud.VPCRouterFirewallRule{
		Protocol:           types.Protocol(stringOrDefault(d, "protocol")),
		SourceNetwork:      types.VPCFirewallNetwork(stringOrDefault(d, "source_network")),
		SourcePort:         types.VPCFirewallPort(stringOrDefault(d, "source_port")),
		DestinationNetwork: types.VPCFirewallNetwork(stringOrDefault(d, "destination_network")),
		DestinationPort:    types.VPCFirewallPort(stringOrDefault(d, "destination_port")),
		Action:             action,
		Logging:            types.StringFlag(boolOrDefault(d, "logging")),
		Description:        stringOrDefault(d, "description"),
	}
}

func flattenVPCRouterFirewalls(vpcRouter *sacloud.VPCRouter) []interface{} {
	var firewallRules []interface{}
	for i, configs := range vpcRouter.Settings.Firewall {

		directionRules := map[string][]*sacloud.VPCRouterFirewallRule{
			"send":    configs.Send,
			"receive": configs.Receive,
		}

		for direction, rules := range directionRules {
			if len(rules) == 0 {
				continue
			}
			var expressions []interface{}
			for _, rule := range rules {
				expression := map[string]interface{}{
					"source_network":      rule.SourceNetwork,
					"source_port":         rule.SourcePort,
					"destination_network": rule.DestinationNetwork,
					"destination_port":    rule.DestinationPort,
					"allow":               rule.Action.IsAllow(),
					"protocol":            rule.Protocol,
					"logging":             rule.Logging.Bool(),
					"description":         rule.Description,
				}
				expressions = append(expressions, expression)
			}
			firewallRules = append(firewallRules, map[string]interface{}{
				"interface_index": i,
				"direction":       direction,
				"expression":      expressions,
			})
		}
	}
	return firewallRules
}

func expandVPCRouterPPTP(d resourceValueGettable) *sacloud.VPCRouterPPTPServer {
	if values, ok := getListFromResource(d, "pptp"); ok && len(values) > 0 {
		raw := values[0]
		d := mapToResourceData(raw.(map[string]interface{}))
		return &sacloud.VPCRouterPPTPServer{
			RangeStart: stringOrDefault(d, "range_start"),
			RangeStop:  stringOrDefault(d, "range_stop"),
		}
	}
	return nil
}

func flattenVPCRouterPPTP(vpcRouter *sacloud.VPCRouter) []interface{} {
	var pptp []interface{}
	if vpcRouter.Settings.PPTPServerEnabled.Bool() {
		c := vpcRouter.Settings.PPTPServer
		pptp = append(pptp, map[string]interface{}{
			"range_start": c.RangeStart,
			"range_stop":  c.RangeStop,
		})
	}
	return pptp
}

func expandVPCRouterL2TP(d resourceValueGettable) *sacloud.VPCRouterL2TPIPsecServer {
	if values, ok := getListFromResource(d, "l2tp"); ok && len(values) > 0 {
		raw := values[0]
		d := mapToResourceData(raw.(map[string]interface{}))
		return &sacloud.VPCRouterL2TPIPsecServer{
			RangeStart:      stringOrDefault(d, "range_start"),
			RangeStop:       stringOrDefault(d, "range_stop"),
			PreSharedSecret: stringOrDefault(d, "pre_shared_secret"),
		}
	}
	return nil
}

func flattenVPCRouterL2TP(vpcRouter *sacloud.VPCRouter) []interface{} {
	var l2tp []interface{}
	if vpcRouter.Settings.L2TPIPsecServerEnabled.Bool() {
		s := vpcRouter.Settings.L2TPIPsecServer
		l2tp = append(l2tp, map[string]interface{}{
			"pre_shared_secret": s.PreSharedSecret,
			"range_start":       s.RangeStart,
			"range_stop":        s.RangeStop,
		})
	}
	return l2tp
}

func expandVPCRouterPortForwardingList(d resourceValueGettable) []*sacloud.VPCRouterPortForwarding {
	if values, ok := getListFromResource(d, "port_forwarding"); ok && len(values) > 0 {
		var results []*sacloud.VPCRouterPortForwarding
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterPortForwarding(v))
		}
		return results
	}
	return nil
}

func expandVPCRouterPortForwarding(d resourceValueGettable) *sacloud.VPCRouterPortForwarding {
	return &sacloud.VPCRouterPortForwarding{
		Protocol:       types.EVPCRouterPortForwardingProtocol(d.Get("protocol").(string)),
		GlobalPort:     types.StringNumber(intOrDefault(d, "public_port")),
		PrivateAddress: stringOrDefault(d, "private_ip"),
		PrivatePort:    types.StringNumber(intOrDefault(d, "private_port")),
		Description:    stringOrDefault(d, "description"),
	}
}

func flattenVPCRouterPortForwardings(vpcRouter *sacloud.VPCRouter) []interface{} {
	var portForwardings []interface{}
	for _, p := range vpcRouter.Settings.PortForwarding {
		globalPort := p.GlobalPort.Int()
		privatePort := p.PrivatePort.Int()
		portForwardings = append(portForwardings, map[string]interface{}{
			"protocol":     string(p.Protocol),
			"public_port":  globalPort,
			"private_ip":   p.PrivateAddress,
			"private_port": privatePort,
			"description":  p.Description,
		})
	}
	return portForwardings
}

func expandVPCRouterSiteToSiteList(d resourceValueGettable) []*sacloud.VPCRouterSiteToSiteIPsecVPN {
	if values, ok := getListFromResource(d, "site_to_site_vpn"); ok && len(values) > 0 {
		var results []*sacloud.VPCRouterSiteToSiteIPsecVPN
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterSiteToSite(v))
		}
		return results
	}
	return nil
}

func expandVPCRouterSiteToSite(d resourceValueGettable) *sacloud.VPCRouterSiteToSiteIPsecVPN {
	return &sacloud.VPCRouterSiteToSiteIPsecVPN{
		Peer:            stringOrDefault(d, "peer"),
		RemoteID:        stringOrDefault(d, "remote_id"),
		PreSharedSecret: stringOrDefault(d, "pre_shared_secret"),
		Routes:          stringListOrDefault(d, "routes"),
		LocalPrefix:     stringListOrDefault(d, "local_prefix"),
	}
}

func flattenVPCRouterSiteToSite(vpcRouter *sacloud.VPCRouter) []interface{} {
	var s2sSettings []interface{}
	for _, s := range vpcRouter.Settings.SiteToSiteIPsecVPN {
		s2sSettings = append(s2sSettings, map[string]interface{}{
			"local_prefix":      s.LocalPrefix,
			"peer":              s.Peer,
			"pre_shared_secret": s.PreSharedSecret,
			"remote_id":         s.RemoteID,
			"routes":            s.Routes,
		})
	}
	return s2sSettings
}

func expandVPCRouterStaticRouteList(d resourceValueGettable) []*sacloud.VPCRouterStaticRoute {
	if values, ok := getListFromResource(d, "static_route"); ok && len(values) > 0 {
		var results []*sacloud.VPCRouterStaticRoute
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterStaticRoute(v))
		}
		return results
	}
	return nil
}

func expandVPCRouterStaticRoute(d resourceValueGettable) *sacloud.VPCRouterStaticRoute {
	return &sacloud.VPCRouterStaticRoute{
		Prefix:  stringOrDefault(d, "prefix"),
		NextHop: stringOrDefault(d, "next_hop"),
	}
}

func flattenVPCRouterStaticRoutes(vpcRouter *sacloud.VPCRouter) []interface{} {
	var staticRoutes []interface{}
	for _, s := range vpcRouter.Settings.StaticRoute {
		staticRoutes = append(staticRoutes, map[string]interface{}{
			"prefix":   s.Prefix,
			"next_hop": s.NextHop,
		})
	}
	return staticRoutes
}

func expandVPCRouterUserList(d resourceValueGettable) []*sacloud.VPCRouterRemoteAccessUser {
	if values, ok := getListFromResource(d, "user"); ok && len(values) > 0 {
		var results []*sacloud.VPCRouterRemoteAccessUser
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterUser(v))
		}
		return results
	}
	return nil
}

func expandVPCRouterUser(d resourceValueGettable) *sacloud.VPCRouterRemoteAccessUser {
	return &sacloud.VPCRouterRemoteAccessUser{
		UserName: stringOrDefault(d, "name"),
		Password: stringOrDefault(d, "password"),
	}
}

func flattenVPCRouterUsers(vpcRouter *sacloud.VPCRouter) []interface{} {
	var users []interface{}
	for _, u := range vpcRouter.Settings.RemoteAccessUsers {
		users = append(users, map[string]interface{}{
			"name":     u.UserName,
			"password": u.Password,
		})
	}
	return users
}
