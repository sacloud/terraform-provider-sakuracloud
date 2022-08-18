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
	"fmt"
	"strconv"
	"strings"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/defaults"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/setup"
	"github.com/sacloud/iaas-service-go/vpcrouter/builder"
)

func expandVPCRouterBuilder(d resourceValueGettable, client *APIClient, zone string) *builder.Builder {
	return &builder.Builder{
		Zone:                  zone,
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Tags:                  expandTags(d),
		IconID:                expandSakuraCloudID(d, "icon_id"),
		PlanID:                expandVPCRouterPlanID(d),
		Version:               d.Get("version").(int),
		NICSetting:            expandVPCRouterNICSetting(d),
		AdditionalNICSettings: expandVPCRouterAdditionalNICSettings(d),
		RouterSetting:         expandVPCRouterSettings(d),
		SetupOptions: &setup.Options{
			BootAfterBuild:        true,
			NICUpdateWaitDuration: defaults.DefaultNICUpdateWaitDuration,
		},
		Client: iaas.NewVPCRouterOp(client),
	}
}

func expandVPCRouterPlanID(d resourceValueGettable) types.ID {
	return types.VPCRouterPlanIDMap[d.Get("plan").(string)]
}

func flattenVPCRouterPlan(vpcRouter *iaas.VPCRouter) string {
	return types.VPCRouterPlanNameMap[vpcRouter.PlanID]
}

func expandVPCRouterNICSetting(d resourceValueGettable) builder.NICSettingHolder {
	planID := expandVPCRouterPlanID(d)
	switch planID {
	case types.VPCRouterPlans.Standard:
		return &builder.StandardNICSetting{}
	default:
		nic := expandVPCRouterPublicNetworkInterface(d)
		return &builder.PremiumNICSetting{
			SwitchID:         nic.switchID,
			IPAddresses:      nic.ipAddresses,
			VirtualIPAddress: nic.vip,
			IPAliases:        nic.ipAliases,
		}
	}
}

type vpcRouterPublicNetworkInterface struct {
	switchID    types.ID
	ipAddresses []string
	vip         string
	ipAliases   []string
	vrid        int
}

func expandVPCRouterPublicNetworkInterface(d resourceValueGettable) *vpcRouterPublicNetworkInterface {
	d = mapFromFirstElement(d, "public_network_interface")
	if d == nil {
		return nil
	}
	return &vpcRouterPublicNetworkInterface{
		switchID:    expandSakuraCloudID(d, "switch_id"),
		ipAddresses: stringListOrDefault(d, "ip_addresses"),
		vip:         stringOrDefault(d, "vip"),
		ipAliases:   stringListOrDefault(d, "aliases"),
		vrid:        intOrDefault(d, "vrid"),
	}
}
func flattenVPCRouterPublicNetworkInterface(vpcRouter *iaas.VPCRouter) []interface{} {
	if vpcRouter.PlanID == types.VPCRouterPlans.Standard {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"switch_id":    flattenVPCRouterSwitchID(vpcRouter),
			"vip":          flattenVPCRouterVIP(vpcRouter),
			"ip_addresses": flattenVPCRouterIPAddresses(vpcRouter),
			"aliases":      flattenVPCRouterIPAliases(vpcRouter),
			"vrid":         flattenVPCRouterVRID(vpcRouter),
		},
	}
}

func expandVPCRouterAdditionalNICSettings(d resourceValueGettable) []builder.AdditionalNICSettingHolder {
	var results []builder.AdditionalNICSettingHolder
	planID := expandVPCRouterPlanID(d)
	interfaces := d.Get("private_network_interface").([]interface{})
	for _, iface := range interfaces {
		d = mapToResourceData(iface.(map[string]interface{}))
		var nicSetting builder.AdditionalNICSettingHolder
		ipAddresses := expandStringList(d.Get("ip_addresses").([]interface{}))

		switch planID {
		case types.VPCRouterPlans.Standard:
			nicSetting = &builder.AdditionalStandardNICSetting{
				SwitchID:       expandSakuraCloudID(d, "switch_id"),
				IPAddress:      ipAddresses[0],
				NetworkMaskLen: d.Get("netmask").(int),
				Index:          d.Get("index").(int),
			}
		default:
			nicSetting = &builder.AdditionalPremiumNICSetting{
				SwitchID:         expandSakuraCloudID(d, "switch_id"),
				NetworkMaskLen:   d.Get("netmask").(int),
				IPAddresses:      ipAddresses,
				VirtualIPAddress: d.Get("vip").(string),
				Index:            d.Get("index").(int),
			}
		}
		results = append(results, nicSetting)
	}
	return results
}

func flattenVPCRouterInterfaces(vpcRouter *iaas.VPCRouter) []interface{} {
	var interfaces []interface{}
	if len(vpcRouter.Interfaces) > 0 {
		for _, iface := range vpcRouter.Settings.Interfaces {
			if iface.Index == 0 {
				continue
			}
			// find nic from data.Interfaces
			var nic *iaas.VPCRouterInterface
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
					"netmask":      iface.NetworkMaskLen,
					"index":        iface.Index,
				})
			}
		}
	}
	return interfaces
}

func flattenVPCRouterGlobalAddress(vpcRouter *iaas.VPCRouter) string {
	if vpcRouter.PlanID == types.VPCRouterPlans.Standard {
		return vpcRouter.Interfaces[0].IPAddress
	}
	return vpcRouter.Settings.Interfaces[0].VirtualIPAddress
}

func flattenVPCRouterGlobalNetworkMaskLen(vpcRouter *iaas.VPCRouter) int {
	return vpcRouter.Interfaces[0].SubnetNetworkMaskLen
}

func flattenVPCRouterSwitchID(vpcRouter *iaas.VPCRouter) string {
	if vpcRouter.PlanID != types.VPCRouterPlans.Standard {
		return vpcRouter.Interfaces[0].SwitchID.String()
	}
	return ""
}

func flattenVPCRouterVIP(vpcRouter *iaas.VPCRouter) string {
	if vpcRouter.PlanID != types.VPCRouterPlans.Standard {
		return vpcRouter.Settings.Interfaces[0].VirtualIPAddress
	}
	return ""
}

func flattenVPCRouterIPAddresses(vpcRouter *iaas.VPCRouter) []string {
	if vpcRouter.PlanID != types.VPCRouterPlans.Standard {
		return vpcRouter.Settings.Interfaces[0].IPAddress
	}
	return []string{}
}

func flattenVPCRouterIPAliases(vpcRouter *iaas.VPCRouter) []string {
	if vpcRouter.PlanID != types.VPCRouterPlans.Standard {
		return vpcRouter.Settings.Interfaces[0].IPAliases
	}
	return []string{}
}

func flattenVPCRouterVRID(vpcRouter *iaas.VPCRouter) int {
	if vpcRouter.PlanID != types.VPCRouterPlans.Standard {
		return vpcRouter.Settings.VRID
	}
	return 0
}

func expandVPCRouterSettings(d resourceValueGettable) *builder.RouterSetting {
	nic := expandVPCRouterPublicNetworkInterface(d)
	vrid := 0
	if nic != nil {
		vrid = nic.vrid
	}
	return &builder.RouterSetting{
		VRID:                      vrid,
		InternetConnectionEnabled: types.StringFlag(d.Get("internet_connection").(bool)),
		StaticNAT:                 expandVPCRouterStaticNATList(d),
		PortForwarding:            expandVPCRouterPortForwardingList(d),
		Firewall:                  expandVPCRouterFirewallList(d),
		DHCPServer:                expandVPCRouterDHCPServerList(d),
		DHCPStaticMapping:         expandVPCRouterDHCPStaticMappingList(d),
		DNSForwarding:             expandVPCRouterDNSForwarding(d),
		PPTPServer:                expandVPCRouterPPTP(d),
		L2TPIPsecServer:           expandVPCRouterL2TP(d),
		RemoteAccessUsers:         expandVPCRouterUserList(d),
		WireGuard:                 expandVPCRouterWireGuard(d),
		SiteToSiteIPsecVPN:        expandVPCRouterSiteToSite(d),
		StaticRoute:               expandVPCRouterStaticRouteList(d),
		SyslogHost:                d.Get("syslog_host").(string),
	}
}

func expandVPCRouterStaticNATList(d resourceValueGettable) []*iaas.VPCRouterStaticNAT {
	if values, ok := getListFromResource(d, "static_nat"); ok && len(values) > 0 {
		var results []*iaas.VPCRouterStaticNAT
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterStaticNAT(v))
		}
		return results
	}
	return nil
}

func expandVPCRouterStaticNAT(d resourceValueGettable) *iaas.VPCRouterStaticNAT {
	return &iaas.VPCRouterStaticNAT{
		GlobalAddress:  d.Get("public_ip").(string),
		PrivateAddress: d.Get("private_ip").(string),
		Description:    d.Get("description").(string),
	}
}

func flattenVPCRouterStaticNAT(vpcRouter *iaas.VPCRouter) []interface{} {
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

func expandVPCRouterDHCPServerList(d resourceValueGettable) []*iaas.VPCRouterDHCPServer {
	if values, ok := getListFromResource(d, "dhcp_server"); ok && len(values) > 0 {
		var results []*iaas.VPCRouterDHCPServer
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterDHCPServer(v))
		}
		return results
	}
	return nil
}

func expandVPCRouterDHCPServer(d resourceValueGettable) *iaas.VPCRouterDHCPServer {
	return &iaas.VPCRouterDHCPServer{
		Interface:  fmt.Sprintf("eth%d", d.Get("interface_index").(int)),
		RangeStart: d.Get("range_start").(string),
		RangeStop:  d.Get("range_stop").(string),
		DNSServers: expandStringList(d.Get("dns_servers").([]interface{})),
	}
}

func flattenVPCRouterDHCPServers(vpcRouter *iaas.VPCRouter) []interface{} {
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

func expandVPCRouterDHCPStaticMappingList(d resourceValueGettable) []*iaas.VPCRouterDHCPStaticMapping {
	if values, ok := getListFromResource(d, "dhcp_static_mapping"); ok && len(values) > 0 {
		var results []*iaas.VPCRouterDHCPStaticMapping
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterDHCPStaticMapping(v))
		}
		return results
	}
	return nil
}

func expandVPCRouterDHCPStaticMapping(d resourceValueGettable) *iaas.VPCRouterDHCPStaticMapping {
	return &iaas.VPCRouterDHCPStaticMapping{
		IPAddress:  d.Get("ip_address").(string),
		MACAddress: d.Get("mac_address").(string),
	}
}

func flattenVPCRouterDHCPStaticMappings(vpcRouter *iaas.VPCRouter) []interface{} {
	var staticMappings []interface{}
	for _, d := range vpcRouter.Settings.DHCPStaticMapping {
		staticMappings = append(staticMappings, map[string]interface{}{
			"ip_address":  d.IPAddress,
			"mac_address": d.MACAddress,
		})
	}
	return staticMappings
}

func expandVPCRouterDNSForwarding(d resourceValueGettable) *iaas.VPCRouterDNSForwarding {
	if values, ok := getListFromResource(d, "dns_forwarding"); ok && len(values) > 0 {
		raw := values[0]
		d := mapToResourceData(raw.(map[string]interface{}))
		return &iaas.VPCRouterDNSForwarding{
			Interface:  fmt.Sprintf("eth%d", d.Get("interface_index").(int)),
			DNSServers: expandStringList(d.Get("dns_servers").([]interface{})),
		}
	}
	return nil
}

func flattenVPCRouterDNSForwarding(vpcRouter *iaas.VPCRouter) []interface{} {
	v := vpcRouter.Settings.DNSForwarding
	if v != nil {
		return []interface{}{
			map[string]interface{}{
				"interface_index": vpcRouterInterfaceNameToIndex(v.Interface),
				"dns_servers":     v.DNSServers,
			},
		}
	}
	return nil
}

func expandVPCRouterFirewallList(d resourceValueGettable) []*iaas.VPCRouterFirewall {
	if values, ok := getListFromResource(d, "firewall"); ok && len(values) > 0 {
		var results []*iaas.VPCRouterFirewall
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterFirewall(v))
		}

		// インデックスごとにSend/Receiveをまとめる
		// results: {Index: 0, Send: []Rules{...}, Receive: nil} , {Index: 0, Send: nil, Receive: []Rules{...}}
		// merged: {Index: 0, Send: []Rules{...}, Receive: []Rules{...}}
		var merged []*iaas.VPCRouterFirewall
		for i := 0; i < 8; i++ {
			firewall := &iaas.VPCRouterFirewall{
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

func expandVPCRouterFirewall(d resourceValueGettable) *iaas.VPCRouterFirewall {
	index := intOrDefault(d, "interface_index")
	direction := stringOrDefault(d, "direction")
	f := &iaas.VPCRouterFirewall{
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

func expandVPCRouterFirewallRuleList(d resourceValueGettable) []*iaas.VPCRouterFirewallRule {
	if values, ok := getListFromResource(d, "expression"); ok && len(values) > 0 {
		var results []*iaas.VPCRouterFirewallRule
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterFirewallRule(v))
		}
		return results
	}
	return nil
}

func expandVPCRouterFirewallRule(d resourceValueGettable) *iaas.VPCRouterFirewallRule {
	allow := boolOrDefault(d, "allow")
	action := types.Actions.Allow
	if !allow {
		action = types.Actions.Deny
	}

	return &iaas.VPCRouterFirewallRule{
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

func flattenVPCRouterFirewalls(vpcRouter *iaas.VPCRouter) []interface{} {
	var firewallRules []interface{}
	for i, configs := range vpcRouter.Settings.Firewall {
		directionRules := map[string][]*iaas.VPCRouterFirewallRule{
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

func expandVPCRouterPPTP(d resourceValueGettable) *iaas.VPCRouterPPTPServer {
	if values, ok := getListFromResource(d, "pptp"); ok && len(values) > 0 {
		raw := values[0]
		d := mapToResourceData(raw.(map[string]interface{}))
		return &iaas.VPCRouterPPTPServer{
			RangeStart: stringOrDefault(d, "range_start"),
			RangeStop:  stringOrDefault(d, "range_stop"),
		}
	}
	return nil
}

func flattenVPCRouterPPTP(vpcRouter *iaas.VPCRouter) []interface{} {
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

func expandVPCRouterL2TP(d resourceValueGettable) *iaas.VPCRouterL2TPIPsecServer {
	if values, ok := getListFromResource(d, "l2tp"); ok && len(values) > 0 {
		raw := values[0]
		d := mapToResourceData(raw.(map[string]interface{}))
		return &iaas.VPCRouterL2TPIPsecServer{
			RangeStart:      stringOrDefault(d, "range_start"),
			RangeStop:       stringOrDefault(d, "range_stop"),
			PreSharedSecret: stringOrDefault(d, "pre_shared_secret"),
		}
	}
	return nil
}

func flattenVPCRouterL2TP(vpcRouter *iaas.VPCRouter) []interface{} {
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

func expandVPCRouterWireGuard(d resourceValueGettable) *iaas.VPCRouterWireGuard {
	if values, ok := getListFromResource(d, "wire_guard"); ok && len(values) > 0 {
		raw := values[0]
		d := mapToResourceData(raw.(map[string]interface{}))

		var peers []*iaas.VPCRouterWireGuardPeer
		if peerValues, ok := getListFromResource(d, "peer"); ok && len(peerValues) > 0 {
			for _, v := range peerValues {
				d := mapToResourceData(v.(map[string]interface{}))
				peers = append(peers, &iaas.VPCRouterWireGuardPeer{
					Name:      stringOrDefault(d, "name"),
					IPAddress: stringOrDefault(d, "ip_address"),
					PublicKey: stringOrDefault(d, "public_key"),
				})
			}
		}

		return &iaas.VPCRouterWireGuard{
			IPAddress: stringOrDefault(d, "ip_address"),
			Peers:     peers,
		}
	}
	return nil
}

func flattenVPCRouterWireGuard(vpcRouter *iaas.VPCRouter, publicKey string) []interface{} {
	var wireGuard []interface{}
	if vpcRouter.Settings.WireGuardEnabled.Bool() {
		s := vpcRouter.Settings.WireGuard
		var peers []interface{}
		for _, peer := range s.Peers {
			peers = append(peers, map[string]interface{}{
				"name":       peer.Name,
				"ip_address": peer.IPAddress,
				"public_key": peer.PublicKey,
			})
		}

		wireGuard = append(wireGuard, map[string]interface{}{
			"ip_address": s.IPAddress,
			"public_key": publicKey,
			"peer":       peers,
		})
	}
	return wireGuard
}

func expandVPCRouterPortForwardingList(d resourceValueGettable) []*iaas.VPCRouterPortForwarding {
	if values, ok := getListFromResource(d, "port_forwarding"); ok && len(values) > 0 {
		var results []*iaas.VPCRouterPortForwarding
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterPortForwarding(v))
		}
		return results
	}
	return nil
}

func expandVPCRouterPortForwarding(d resourceValueGettable) *iaas.VPCRouterPortForwarding {
	return &iaas.VPCRouterPortForwarding{
		Protocol:       types.EVPCRouterPortForwardingProtocol(d.Get("protocol").(string)),
		GlobalPort:     types.StringNumber(intOrDefault(d, "public_port")),
		PrivateAddress: stringOrDefault(d, "private_ip"),
		PrivatePort:    types.StringNumber(intOrDefault(d, "private_port")),
		Description:    stringOrDefault(d, "description"),
	}
}

func flattenVPCRouterPortForwardings(vpcRouter *iaas.VPCRouter) []interface{} {
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

func expandVPCRouterSiteToSite(d resourceValueGettable) *iaas.VPCRouterSiteToSiteIPsecVPN {
	siteToSiteVPN := &iaas.VPCRouterSiteToSiteIPsecVPN{}
	if values, ok := getListFromResource(d, "site_to_site_vpn"); ok && len(values) > 0 {
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			siteToSiteVPN.Config = append(siteToSiteVPN.Config, expandVPCRouterSiteToSiteConfig(v))
		}
	}
	if values, ok := getListFromResource(d, "site_to_site_vpn_parameter"); ok && len(values) > 0 {
		raw := values[0]
		d := mapToResourceData(raw.(map[string]interface{}))

		siteToSiteVPN.IKE = expandVPCRouterSiteToSiteParameterIKE(d)
		siteToSiteVPN.ESP = expandVPCRouterSiteToSiteParameterESP(d)
		siteToSiteVPN.EncryptionAlgo = stringOrDefault(d, "encryption_algo")
		siteToSiteVPN.HashAlgo = stringOrDefault(d, "hash_algo")
	}

	return siteToSiteVPN
}

func expandVPCRouterSiteToSiteConfig(d resourceValueGettable) *iaas.VPCRouterSiteToSiteIPsecVPNConfig {
	return &iaas.VPCRouterSiteToSiteIPsecVPNConfig{
		Peer:            stringOrDefault(d, "peer"),
		RemoteID:        stringOrDefault(d, "remote_id"),
		PreSharedSecret: stringOrDefault(d, "pre_shared_secret"),
		Routes:          stringSetOrDefault(d, "routes"),
		LocalPrefix:     stringSetOrDefault(d, "local_prefix"),
	}
}

func expandVPCRouterSiteToSiteParameterIKE(d resourceValueGettable) *iaas.VPCRouterSiteToSiteIPsecVPNIKE {
	if values, ok := getListFromResource(d, "ike"); ok && len(values) > 0 {
		raw := values[0]
		d := mapToResourceData(raw.(map[string]interface{}))
		return &iaas.VPCRouterSiteToSiteIPsecVPNIKE{
			Lifetime: intOrDefault(d, "lifetime"),
			DPD:      expandVPCRouterSiteToSiteParameterIKEDPD(d),
		}
	}
	return nil
}

func expandVPCRouterSiteToSiteParameterIKEDPD(d resourceValueGettable) *iaas.VPCRouterSiteToSiteIPsecVPNIKEDPD {
	if values, ok := getListFromResource(d, "dpd"); ok && len(values) > 0 {
		raw := values[0]
		d := mapToResourceData(raw.(map[string]interface{}))
		return &iaas.VPCRouterSiteToSiteIPsecVPNIKEDPD{
			Interval: intOrDefault(d, "interval"),
			Timeout:  intOrDefault(d, "timeout"),
		}
	}
	return nil
}

func expandVPCRouterSiteToSiteParameterESP(d resourceValueGettable) *iaas.VPCRouterSiteToSiteIPsecVPNESP {
	if values, ok := getListFromResource(d, "esp"); ok && len(values) > 0 {
		raw := values[0]
		d := mapToResourceData(raw.(map[string]interface{}))
		return &iaas.VPCRouterSiteToSiteIPsecVPNESP{
			Lifetime: intOrDefault(d, "lifetime"),
		}
	}
	return nil
}

func flattenVPCRouterSiteToSiteConfig(vpcRouter *iaas.VPCRouter) []interface{} {
	var s2sSettings []interface{}
	if vpcRouter.Settings.SiteToSiteIPsecVPN != nil {
		for _, s := range vpcRouter.Settings.SiteToSiteIPsecVPN.Config {
			s2sSettings = append(s2sSettings, map[string]interface{}{
				"local_prefix":      stringListToSet(s.LocalPrefix),
				"peer":              s.Peer,
				"pre_shared_secret": s.PreSharedSecret,
				"remote_id":         s.RemoteID,
				"routes":            stringListToSet(s.Routes),
			})
		}
	}
	return s2sSettings
}

func flattenVPCRouterSiteToSiteParameter(vpcRouter *iaas.VPCRouter) []interface{} {
	var s2sParameters []interface{}
	if vpcRouter.Settings.SiteToSiteIPsecVPN != nil {
		v := map[string]interface{}{
			"encryption_algo": vpcRouter.Settings.SiteToSiteIPsecVPN.EncryptionAlgo,
			"hash_algo":       vpcRouter.Settings.SiteToSiteIPsecVPN.HashAlgo,
		}
		if vpcRouter.Settings.SiteToSiteIPsecVPN.IKE != nil {
			ike := map[string]interface{}{
				"lifetime": vpcRouter.Settings.SiteToSiteIPsecVPN.IKE.Lifetime,
			}
			if vpcRouter.Settings.SiteToSiteIPsecVPN.IKE.DPD != nil {
				ike["dpd"] = []interface{}{map[string]interface{}{
					"interval": vpcRouter.Settings.SiteToSiteIPsecVPN.IKE.DPD.Interval,
					"timeout":  vpcRouter.Settings.SiteToSiteIPsecVPN.IKE.DPD.Timeout,
				}}
			}
			v["ike"] = []interface{}{ike}
		}
		if vpcRouter.Settings.SiteToSiteIPsecVPN.ESP != nil {
			v["esp"] = []interface{}{map[string]interface{}{
				"lifetime": vpcRouter.Settings.SiteToSiteIPsecVPN.ESP.Lifetime,
			}}
		}
		s2sParameters = append(s2sParameters, v)
	}
	return s2sParameters
}

func expandVPCRouterStaticRouteList(d resourceValueGettable) []*iaas.VPCRouterStaticRoute {
	if values, ok := getListFromResource(d, "static_route"); ok && len(values) > 0 {
		var results []*iaas.VPCRouterStaticRoute
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterStaticRoute(v))
		}
		return results
	}
	return nil
}

func expandVPCRouterStaticRoute(d resourceValueGettable) *iaas.VPCRouterStaticRoute {
	return &iaas.VPCRouterStaticRoute{
		Prefix:  stringOrDefault(d, "prefix"),
		NextHop: stringOrDefault(d, "next_hop"),
	}
}

func flattenVPCRouterStaticRoutes(vpcRouter *iaas.VPCRouter) []interface{} {
	var staticRoutes []interface{}
	for _, s := range vpcRouter.Settings.StaticRoute {
		staticRoutes = append(staticRoutes, map[string]interface{}{
			"prefix":   s.Prefix,
			"next_hop": s.NextHop,
		})
	}
	return staticRoutes
}

func expandVPCRouterUserList(d resourceValueGettable) []*iaas.VPCRouterRemoteAccessUser {
	if values, ok := getListFromResource(d, "user"); ok && len(values) > 0 {
		var results []*iaas.VPCRouterRemoteAccessUser
		for _, raw := range values {
			v := mapToResourceData(raw.(map[string]interface{}))
			results = append(results, expandVPCRouterUser(v))
		}
		return results
	}
	return nil
}

func expandVPCRouterUser(d resourceValueGettable) *iaas.VPCRouterRemoteAccessUser {
	return &iaas.VPCRouterRemoteAccessUser{
		UserName: stringOrDefault(d, "name"),
		Password: stringOrDefault(d, "password"),
	}
}

func flattenVPCRouterUsers(vpcRouter *iaas.VPCRouter) []interface{} {
	var users []interface{}
	for _, u := range vpcRouter.Settings.RemoteAccessUsers {
		users = append(users, map[string]interface{}{
			"name":     u.UserName,
			"password": u.Password,
		})
	}
	return users
}
