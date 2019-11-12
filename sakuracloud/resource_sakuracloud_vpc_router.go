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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/utils/setup"
)

const vpcRouterPowerAPILockKey = "sakuracloud_vpc_router.power.%d.lock"

func resourceSakuraCloudVPCRouter() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterCreate,
		Read:   resourceSakuraCloudVPCRouterRead,
		Update: resourceSakuraCloudVPCRouterUpdate,
		Delete: resourceSakuraCloudVPCRouterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"plan": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "standard",
				ValidateFunc: validation.StringInSlice([]string{"standard", "premium", "highspec"}, false),
			},
			"switch_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"vip": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"ipaddress1": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"ipaddress2": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"vrid": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
			},
			"aliases": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MaxItems: 19,
			},
			"syslog_host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internet_connection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"interface": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 7,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterInterfaceEmbeddedSchema(),
				},
			},
			"dhcp_server": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterDHCPServerEmbeddedSchema(),
				},
			},
			"dhcp_static_mapping": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterDHCPStaticMappingEmbeddedSchema(),
				},
			},
			"firewall": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterFirewallEmbeddedSchema(),
				},
			},
			"l2tp": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: vpcRouterL2TPEmbeddedSchema(),
				},
			},
			"port_forwarding": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterPortForwardingEmbeddedSchema(),
				},
			},
			"pptp": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: vpcRouterPPTPEmbeddedSchema(),
				},
			},
			"site_to_site_vpn": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterS2SEmbeddedSchema(),
				},
			},
			"static_nat": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterStaticNATEmbeddedSchema(),
				},
			},
			"static_route": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterStaticRouteEmbeddedSchema(),
				},
			},
			"user": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 100,
				Elem: &schema.Resource{
					Schema: vpcRouterUserEmbeddedSchema(),
				},
			},
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 512),
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			powerManageTimeoutKey: powerManageTimeoutParam,
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"global_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudVPCRouterCreate(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	opts := client.VPCRouter.New()

	plan := d.Get("plan").(string)
	switch plan {
	case "standard":
		opts.SetStandardPlan()
	case "premium", "highspec":
		switchID := ""
		vip := ""
		ipaddress1 := ""
		ipaddress2 := ""
		vrid := -1
		aliases := []string{}

		//validate
		errFormat := "Failed to create SakuraCloud VPCRouter resource : %s is Required when plan is 'premium' or 'highspec'"
		if s, ok := d.GetOk("switch_id"); ok {
			switchID = s.(string)
		} else {
			return fmt.Errorf(errFormat, "switch_id")
		}
		if s, ok := d.GetOk("vip"); ok {
			vip = s.(string)
		} else {
			return fmt.Errorf(errFormat, "vip")
		}

		if s, ok := d.GetOk("ipaddress1"); ok {
			ipaddress1 = s.(string)
		} else {
			return fmt.Errorf(errFormat, "ipaddress1")
		}
		if s, ok := d.GetOk("ipaddress2"); ok {
			ipaddress2 = s.(string)
		} else {
			return fmt.Errorf(errFormat, "ipaddress2")
		}

		if s, ok := d.GetOk("vrid"); ok {
			vrid = s.(int)
		} else {
			return fmt.Errorf(errFormat, "vrid")
		}

		if list, ok := d.GetOk("aliases"); ok {
			rawAliases := list.([]interface{})
			for _, a := range rawAliases {
				aliases = append(aliases, a.(string))
			}
		}

		if plan == "premium" {
			opts.SetPremiumPlan(switchID, vip, ipaddress1, ipaddress2, vrid, aliases)
		} else {
			opts.SetHighSpecPlan(switchID, vip, ipaddress1, ipaddress2, vrid, aliases)
		}
	}

	opts.Name = d.Get("name").(string)
	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.SetIconByID(toSakuraCloudID(iconID.(string)))
	}
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		opts.Tags = expandTags(client, rawTags)
	}

	opts.InitVPCRouterSetting()
	if syslogHost, ok := d.GetOk("syslog_host"); ok {
		opts.Settings.Router.SyslogHost = syslogHost.(string)
	}

	if d.Get("internet_connection").(bool) {
		opts.Settings.Router.InternetConnection = &sacloud.VPCRouterInternetConnection{
			Enabled: "True",
		}
	}

	vpcRouterBuilder := &setup.RetryableSetup{
		Create: func() (sacloud.ResourceIDHolder, error) {
			return client.VPCRouter.Create(opts)
		},
		AsyncWaitForCopy: func(id int64) (chan interface{}, chan interface{}, chan error) {
			return client.VPCRouter.AsyncSleepWhileCopying(id, client.DefaultTimeoutDuration, 20)
		},
		Delete: func(id int64) error {
			_, err := client.VPCRouter.Delete(id)
			return err
		},
		ProvisionBeforeUp: func(id int64, created interface{}) error {
			vpcRouter := created.(*sacloud.VPCRouter)

			if interfaces, ok := getListFromResource(d, "interface"); ok && len(interfaces) > 0 {
				for i, iface := range interfaces {
					if iface == nil {
						continue
					}
					values := mapToResourceData(iface.(map[string]interface{}))

					index := i + 1
					switchID := values.Get("switch_id").(string)
					var vip string
					if v, ok := values.GetOk("vip"); ok {
						vip = v.(string)
					}
					nwMaskLen := values.Get("nw_mask_len").(int)
					var ipaddresses []string
					if ipList, ok := getListFromResource(values, "ipaddress"); ok && len(ipList) > 0 {
						for _, ip := range ipList {
							ipaddresses = append(ipaddresses, ip.(string))
						}
					}

					if len(ipaddresses) == 0 {
						return fmt.Errorf("SakuraCloud VPCRouter: ipaddresses is required on interface.%d", i)
					}

					if vpcRouter.IsStandardPlan() {
						v, err := client.VPCRouter.AddStandardInterfaceAt(vpcRouter.ID, toSakuraCloudID(switchID), ipaddresses[0], nwMaskLen, index)
						if err != nil {
							return err
						}
						vpcRouter = v
					} else {
						v, err := client.VPCRouter.AddPremiumInterfaceAt(vpcRouter.ID, toSakuraCloudID(switchID), ipaddresses, nwMaskLen, vip, index)
						if err != nil {
							return err
						}
						vpcRouter = v
					}
				}
			}

			// DHCP Server
			if dhcpServers, ok := getListFromResource(d, "dhcp_server"); ok && len(dhcpServers) > 0 {
				for _, rawDHCPServer := range dhcpServers {
					values := mapToResourceData(rawDHCPServer.(map[string]interface{}))

					dhcpServer := expandVPCRouterDHCPServer(values)
					vpcRouter.Settings.Router.AddDHCPServer(values.Get("vpc_router_interface_index").(int),
						dhcpServer.RangeStart, dhcpServer.RangeStop,
						dhcpServer.DNSServers...)
				}
			}

			// DHCP static mapping
			if staticMappings, ok := getListFromResource(d, "dhcp_static_mapping"); ok && len(staticMappings) > 0 {
				for _, rawMapping := range staticMappings {
					values := mapToResourceData(rawMapping.(map[string]interface{}))

					mapping := expandVPCRouterDHCPStaticMapping(values)
					vpcRouter.Settings.Router.AddDHCPStaticMapping(mapping.IPAddress, mapping.MACAddress)
				}
			}

			// Firewall rules
			if firewallRules, ok := getListFromResource(d, "firewall"); ok && len(firewallRules) > 0 {
				for _, rawRules := range firewallRules {
					values := mapToResourceData(rawRules.(map[string]interface{}))

					ifIndex := values.Get("vpc_router_interface_index").(int)
					direction := values.Get("direction").(string)

					// clear rules
					if vpcRouter.HasFirewall() && len(vpcRouter.Settings.Router.Firewall.Config) > ifIndex {
						switch direction {
						case "send":
							vpcRouter.Settings.Router.Firewall.Config[ifIndex].Send = nil
						case "receive":
							vpcRouter.Settings.Router.Firewall.Config[ifIndex].Receive = nil
						}
					}

					if rawExpressions, ok := values.GetOk("expressions"); ok {
						expressions := rawExpressions.([]interface{})
						for _, e := range expressions {
							exp := e.(map[string]interface{})

							allow := exp["allow"].(bool)
							protocol := exp["protocol"].(string)
							sourceNW := exp["source_nw"].(string)
							sourcePort := exp["source_port"].(string)
							destNW := exp["dest_nw"].(string)
							destPort := exp["dest_port"].(string)
							logging := exp["logging"].(bool)
							desc := ""
							if de, ok := exp["description"]; ok {
								desc = de.(string)
							}

							switch direction {
							case "send":
								vpcRouter.Settings.Router.AddFirewallRuleSend(ifIndex, allow, protocol, sourceNW, sourcePort, destNW, destPort, logging, desc)
							case "receive":
								vpcRouter.Settings.Router.AddFirewallRuleReceive(ifIndex, allow, protocol, sourceNW, sourcePort, destNW, destPort, logging, desc)
							}
						}
					}
				}
			}

			// L2TP
			if l2tpSettings, ok := getListFromResource(d, "l2tp"); ok && len(l2tpSettings) > 0 {
				if l2tpSettings[0] != nil {
					values := mapToResourceData(l2tpSettings[0].(map[string]interface{}))
					l2tp := expandVPCRouterL2TP(values)
					vpcRouter.Settings.Router.EnableL2TPIPsecServer(l2tp.PreSharedSecret, l2tp.RangeStart, l2tp.RangeStop)
				}
			}

			// PortForwarding
			if portForwardings, ok := getListFromResource(d, "port_forwarding"); ok && len(portForwardings) > 0 {
				for _, rawPortForwarding := range portForwardings {
					values := mapToResourceData(rawPortForwarding.(map[string]interface{}))
					pf := expandVPCRouterPortForwarding(values)
					vpcRouter.Settings.Router.AddPortForwarding(pf.Protocol, pf.GlobalPort, pf.PrivateAddress, pf.PrivatePort, pf.Description)
				}
			}

			// PPTP
			if pptpSettings, ok := getListFromResource(d, "pptp"); ok && len(pptpSettings) > 0 {
				if pptpSettings[0] != nil {
					values := mapToResourceData(pptpSettings[0].(map[string]interface{}))
					pptp := expandVPCRouterPPTP(values)
					vpcRouter.Settings.Router.EnablePPTPServer(pptp.RangeStart, pptp.RangeStop)
				}
			}

			// SiteToSite VPN
			if s2sSettings, ok := getListFromResource(d, "site_to_site_vpn"); ok && len(s2sSettings) > 0 {
				for _, rawS2s := range s2sSettings {
					values := mapToResourceData(rawS2s.(map[string]interface{}))
					s2s := expandVPCRouterSiteToSiteIPsecVPN(values)
					vpcRouter.Settings.Router.AddSiteToSiteIPsecVPN(s2s.LocalPrefix, s2s.Peer, s2s.PreSharedSecret, s2s.RemoteID, s2s.Routes)
				}
			}

			// Static NAT
			if staticNATSettings, ok := getListFromResource(d, "static_nat"); ok && len(staticNATSettings) > 0 {
				for _, rawStaticNAT := range staticNATSettings {
					values := mapToResourceData(rawStaticNAT.(map[string]interface{}))
					staticNAT := expandVPCRouterStaticNAT(values)
					vpcRouter.Settings.Router.AddStaticNAT(staticNAT.GlobalAddress, staticNAT.PrivateAddress, staticNAT.Description)
				}
			}

			// Static Routes
			if staticRoutes, ok := getListFromResource(d, "static_route"); ok && len(staticRoutes) > 0 {
				for _, rawStaticRoute := range staticRoutes {
					values := mapToResourceData(rawStaticRoute.(map[string]interface{}))
					staticRoute := expandVPCRouterStaticRoute(values)
					vpcRouter.Settings.Router.AddStaticRoute(staticRoute.Prefix, staticRoute.NextHop)
				}
			}

			// Users
			if users, ok := getListFromResource(d, "user"); ok && len(users) > 0 {
				for _, rawUser := range users {
					values := mapToResourceData(rawUser.(map[string]interface{}))
					user := expandVPCRouterRemoteAccessUser(values)
					vpcRouter.Settings.Router.AddRemoteAccessUser(user.UserName, user.Password)
				}
			}

			var err error
			vpcRouter, err = client.VPCRouter.UpdateSetting(vpcRouter.ID, vpcRouter)
			if err != nil {
				return fmt.Errorf("Error creating SakuraCloud VPCRouter resource: %s", err)
			}
			if _, err = client.VPCRouter.Config(vpcRouter.ID); err != nil {
				return fmt.Errorf("Error creating SakuraCloud VPCRouter settings: %s", err)
			}
			if _, err := client.VPCRouter.Boot(id); err != nil {
				return fmt.Errorf("Failed to boot SakuraCloud VPCRouter resource: %s", err)
			}
			return nil
		},
		WaitForUp: func(id int64) error {
			return client.VPCRouter.SleepUntilUp(id, client.DefaultTimeoutDuration)
		},
		RetryCount:             3,
		ProvisioningRetryCount: 1,
	}

	res, err := vpcRouterBuilder.Setup()
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud VPCRouter resource: %s", err)
	}

	vpcRouter, ok := res.(*sacloud.VPCRouter)
	if !ok {
		return fmt.Errorf("Failed to create SakuraCloud VPCRouter resource: created resource is not *sacloud.VPCRouter")
	}

	d.SetId(vpcRouter.GetStrID())
	return resourceSakuraCloudVPCRouterRead(d, meta)
}

func resourceSakuraCloudVPCRouterRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	return setVPCRouterResourceData(d, client, vpcRouter)
}

func setVPCRouterResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.VPCRouter) error {

	if data.IsFailed() {
		d.SetId("")
		return fmt.Errorf("VPCRouter[%d] state is failed", data.ID)
	}

	d.Set("name", data.Name)
	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	if data.Settings != nil && data.Settings.Router != nil {
		d.Set("syslog_host", data.Settings.Router.SyslogHost)

		in := data.Settings.Router.InternetConnection
		if in != nil && in.Enabled == "True" {
			d.Set("internet_connection", true)
		} else {
			d.Set("internet_connection", false)
		}

	} else {
		d.Set("syslog_host", "")
		d.Set("internet_connection", false)
	}
	d.Set("tags", data.Tags)

	//plan
	planID := data.Plan.ID
	switch planID {
	case 1:
		d.Set("plan", "standard")
	case 2:
		d.Set("plan", "premium")
	case 3:
		d.Set("plan", "highspec")
	}
	if planID == 1 {
		d.Set("global_address", data.Interfaces[0].IPAddress)
	} else {
		d.Set("switch_id", data.Switch.GetStrID())
		d.Set("vip", data.Settings.Router.Interfaces[0].VirtualIPAddress)
		d.Set("ipaddress1", data.Settings.Router.Interfaces[0].IPAddress[0])
		d.Set("ipaddress2", data.Settings.Router.Interfaces[0].IPAddress[1])
		d.Set("aliases", data.Settings.Router.Interfaces[0].IPAliases)
		d.Set("vrid", data.Settings.Router.VRID)

		d.Set("global_address", data.Settings.Router.Interfaces[0].VirtualIPAddress)
	}

	setPowerManageTimeoutValueToState(d)

	// interface
	var interfaces []map[string]interface{}
	if data.HasInterfaces() {
		for i, iface := range data.Settings.Router.Interfaces {
			if i == 0 {
				continue
			}
			interfaces = append(interfaces, map[string]interface{}{
				"switch_id":   data.Interfaces[i].Switch.GetStrID(),
				"vip":         iface.VirtualIPAddress,
				"ipaddress":   iface.IPAddress,
				"nw_mask_len": iface.NetworkMaskLen,
			})
		}
	}
	d.Set("interface", interfaces)

	var dhcpServers []map[string]interface{}
	if data.HasDHCPServer() {
		for _, c := range data.Settings.Router.DHCPServer.Config {
			dhcpServers = append(dhcpServers, map[string]interface{}{
				"range_start":                c.RangeStart,
				"range_stop":                 c.RangeStop,
				"vpc_router_interface_index": c.InterfaceIndex(),
				"dns_servers":                c.DNSServers,
			})
		}
	}
	d.Set("dhcp_server", dhcpServers)

	var staticMappings []map[string]interface{}
	if data.HasDHCPStaticMapping() {
		for _, c := range data.Settings.Router.DHCPStaticMapping.Config {
			staticMappings = append(staticMappings, map[string]interface{}{
				"ipaddress":  c.IPAddress,
				"macaddress": c.MACAddress,
			})
		}
	}
	d.Set("dhcp_static_mapping", staticMappings)

	var firewallRules []map[string]interface{}
	if data.HasFirewall() {
		for i, configs := range data.Settings.Router.Firewall.Config {

			directionRules := map[string][]*sacloud.VPCRouterFirewallRule{
				"send":    configs.Send,
				"receive": configs.Receive,
			}

			for direction, rules := range directionRules {
				if len(rules) == 0 {
					continue
				}
				expressions := []interface{}{}
				for _, rule := range rules {
					expression := map[string]interface{}{
						"source_nw":   rule.SourceNetwork,
						"source_port": rule.SourcePort,
						"dest_nw":     rule.DestinationNetwork,
						"dest_port":   rule.DestinationPort,
						"allow":       rule.Action == "allow",
						"protocol":    rule.Protocol,
						"logging":     strings.ToLower(rule.Logging) == "true",
						"description": rule.Description,
					}
					expressions = append(expressions, expression)
				}
				firewallRules = append(firewallRules, map[string]interface{}{
					"vpc_router_interface_index": i,
					"direction":                  direction,
					"expressions":                expressions,
				})
			}
		}
	}
	d.Set("firewall", firewallRules)

	var l2tp []map[string]interface{}
	if data.HasL2TPIPsecServer() {
		c := data.Settings.Router.L2TPIPsecServer.Config
		l2tp = append(l2tp, map[string]interface{}{
			"pre_shared_secret": c.PreSharedSecret,
			"range_start":       c.RangeStart,
			"range_stop":        c.RangeStop,
		})
	}
	d.Set("l2tp", l2tp)

	var portForwardings []map[string]interface{}
	if data.HasPortForwarding() {
		for _, c := range data.Settings.Router.PortForwarding.Config {
			globalPort, _ := strconv.Atoi(c.GlobalPort)
			privatePort, _ := strconv.Atoi(c.PrivatePort)
			portForwardings = append(portForwardings, map[string]interface{}{
				"protocol":        c.Protocol,
				"global_port":     globalPort,
				"private_address": c.PrivateAddress,
				"private_port":    privatePort,
				"description":     c.Description,
			})
		}
	}
	d.Set("port_forwarding", portForwardings)

	var pptp []map[string]interface{}
	if data.HasPPTPServer() {
		c := data.Settings.Router.PPTPServer.Config
		pptp = append(pptp, map[string]interface{}{
			"range_start": c.RangeStart,
			"range_stop":  c.RangeStop,
		})
	}
	d.Set("pptp", pptp)

	var s2sSettings []map[string]interface{}
	if data.HasSiteToSiteIPsecVPN() {
		// SiteToSiteConnectionDetail
		connInfo, err := client.VPCRouter.SiteToSiteConnectionDetails(data.ID)
		if err != nil {
			return fmt.Errorf("Reading VPCRouter SiteToSiteConnectionDetail is failed: %s", err)
		}

		for i, c := range data.Settings.Router.SiteToSiteIPsecVPN.Config {
			detail := connInfo.Details.Config[i]
			s2sSettings = append(s2sSettings, map[string]interface{}{
				"local_prefix":                 c.LocalPrefix,
				"peer":                         c.Peer,
				"pre_shared_secret":            c.PreSharedSecret,
				"remote_id":                    c.RemoteID,
				"routes":                       c.Routes,
				"esp_authentication_protocol":  detail.ESP.AuthenticationProtocol,
				"esp_dh_group":                 detail.ESP.DHGroup,
				"esp_encryption_protocol":      detail.ESP.EncryptionProtocol,
				"esp_lifetime":                 detail.ESP.Lifetime,
				"esp_mode":                     detail.ESP.Mode,
				"esp_perfect_forward_secrecy":  detail.ESP.PerfectForwardSecrecy,
				"ike_authentication_protocol":  detail.IKE.AuthenticationProtocol,
				"ike_encryption_protocol":      detail.IKE.EncryptionProtocol,
				"ike_lifetime":                 detail.IKE.Lifetime,
				"ike_mode":                     detail.IKE.Mode,
				"ike_perfect_forward_secrecy":  detail.IKE.PerfectForwardSecrecy,
				"ike_pre_shared_secret":        detail.IKE.PreSharedSecret,
				"peer_id":                      detail.Peer.ID,
				"peer_inside_networks":         detail.Peer.InsideNetworks,
				"peer_outside_ipaddress":       detail.Peer.OutsideIPAddress,
				"vpc_router_inside_networks":   detail.VPCRouter.InsideNetworks,
				"vpc_router_outside_ipaddress": detail.VPCRouter.OutsideIPAddress,
			})
		}
	}
	d.Set("site_to_site_vpn", s2sSettings)

	var staticNATs []map[string]interface{}
	if data.HasStaticNAT() {
		for _, c := range data.Settings.Router.StaticNAT.Config {
			staticNATs = append(staticNATs, map[string]interface{}{
				"global_address":  c.GlobalAddress,
				"private_address": c.PrivateAddress,
				"description":     c.Description,
			})
		}
	}
	d.Set("static_nat", staticNATs)

	var staticRoutes []map[string]interface{}
	if data.HasStaticRoutes() {
		for _, c := range data.Settings.Router.StaticRoutes.Config {
			staticRoutes = append(staticRoutes, map[string]interface{}{
				"prefix":   c.Prefix,
				"next_hop": c.NextHop,
			})
		}
	}
	d.Set("static_route", staticRoutes)

	var users []map[string]interface{}
	if data.HasRemoteAccessUsers() {
		for _, c := range data.Settings.Router.RemoteAccessUsers.Config {
			users = append(users, map[string]interface{}{
				"name":     c.UserName,
				"password": c.Password,
			})
		}
	}
	d.Set("user", users)

	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudVPCRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	isNeedRestart := false
	if vpcRouter.IsUp() && d.HasChange("interface") {
		isNeedRestart = true
	}

	if isNeedRestart {
		// power API lock
		lockKey := getVPCRouterPowerAPILockKey(vpcRouter.ID)
		sakuraMutexKV.Lock(lockKey)
		defer sakuraMutexKV.Unlock(lockKey)

		err = nil
		for i := 0; i < 10; i++ {
			vpcRouter, err := client.VPCRouter.Read(vpcRouter.ID)
			if err != nil {
				return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
			}
			if vpcRouter.Instance.IsDown() {
				err = nil
				break
			}
			err = handleShutdown(client.VPCRouter, vpcRouter.ID, d, 60*time.Second)
		}
		if err != nil {
			return fmt.Errorf("Error stopping SakuraCloud VPCRouter resource: %s", err)
		}
	}

	if d.HasChange("name") {
		vpcRouter.Name = d.Get("name").(string)
	}
	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			vpcRouter.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			vpcRouter.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			vpcRouter.Description = description.(string)
		} else {
			vpcRouter.Description = ""
		}
	}
	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			vpcRouter.Tags = expandTags(client, rawTags)
		} else {
			vpcRouter.Tags = expandTags(client, []interface{}{})
		}
	}
	if d.HasChange("syslog_host") {

		if vpcRouter.Settings == nil || vpcRouter.Settings.Router == nil {
			vpcRouter.InitVPCRouterSetting()
		}

		if syslogHost, ok := d.GetOk("syslog_host"); ok {
			vpcRouter.Settings.Router.SyslogHost = syslogHost.(string)
		} else {
			vpcRouter.Settings.Router.SyslogHost = ""
		}
	}
	if d.HasChange("internet_connection") {
		vpcRouter.Settings.Router.InternetConnection = &sacloud.VPCRouterInternetConnection{
			Enabled: "False",
		}
		if d.Get("internet_connection").(bool) {
			vpcRouter.Settings.Router.InternetConnection.Enabled = "True"
		}
	}

	if d.HasChange("interface") {
		if vpcRouter.HasInterfaces() {
			for i := range vpcRouter.Settings.Router.Interfaces {
				if i == 0 || vpcRouter.Settings.Router.Interfaces[i] == nil {
					continue
				}
				if _, err := client.VPCRouter.DisconnectFromSwitch(vpcRouter.ID, i); err != nil {
					return fmt.Errorf("Error updating SakuraCloud VPCRouter interface: %s", err)
				}
			}
		}
		if interfaces, ok := getListFromResource(d, "interface"); ok && len(interfaces) > 0 {
			for i, iface := range interfaces {
				if iface == nil {
					continue
				}
				values := mapToResourceData(iface.(map[string]interface{}))

				index := i + 1
				switchID := values.Get("switch_id").(string)
				var vip string
				if v, ok := values.GetOk("vip"); ok {
					vip = v.(string)
				}
				nwMaskLen := values.Get("nw_mask_len").(int)
				var ipaddresses []string
				if ipList, ok := getListFromResource(values, "ipaddress"); ok && len(ipList) > 0 {
					for _, ip := range ipList {
						ipaddresses = append(ipaddresses, ip.(string))
					}
				}

				if len(ipaddresses) == 0 {
					return fmt.Errorf("SakuraCloud VPCRouter: ipaddresses is required on interface.%d", i)
				}

				if vpcRouter.IsStandardPlan() {
					_, err := client.VPCRouter.AddStandardInterfaceAt(vpcRouter.ID, toSakuraCloudID(switchID), ipaddresses[0], nwMaskLen, index)
					if err != nil {
						return err
					}
				} else {
					_, err := client.VPCRouter.AddPremiumInterfaceAt(vpcRouter.ID, toSakuraCloudID(switchID), ipaddresses, nwMaskLen, vip, index)
					if err != nil {
						return err
					}
				}
			}
		}

		refreshedRouter, err := client.VPCRouter.Read(vpcRouter.ID)
		if err != nil {
			return fmt.Errorf("Error updating SakuraCloud VPCRouter resource: can't read VPCRouter %d: %s", vpcRouter.ID, err)
		}
		vpcRouter = refreshedRouter
	}

	if d.HasChange("dhcp_server") {
		vpcRouter.Settings.Router.DHCPServer = nil
		// DHCP Server
		if dhcpServers, ok := getListFromResource(d, "dhcp_server"); ok && len(dhcpServers) > 0 {
			for _, rawDHCPServer := range dhcpServers {
				values := mapToResourceData(rawDHCPServer.(map[string]interface{}))

				dhcpServer := expandVPCRouterDHCPServer(values)
				vpcRouter.Settings.Router.AddDHCPServer(values.Get("vpc_router_interface_index").(int),
					dhcpServer.RangeStart, dhcpServer.RangeStop,
					dhcpServer.DNSServers...)
			}
		}
	}

	if d.HasChange("dhcp_static_mapping") {
		vpcRouter.Settings.Router.DHCPStaticMapping = nil
		if staticMappings, ok := getListFromResource(d, "dhcp_static_mapping"); ok && len(staticMappings) > 0 {
			for _, rawMapping := range staticMappings {
				values := mapToResourceData(rawMapping.(map[string]interface{}))

				mapping := expandVPCRouterDHCPStaticMapping(values)
				vpcRouter.Settings.Router.AddDHCPStaticMapping(mapping.IPAddress, mapping.MACAddress)
			}
		}
	}

	if d.HasChange("firewall") {
		// Firewall rules
		if firewallRules, ok := getListFromResource(d, "firewall"); ok && len(firewallRules) > 0 {
			for _, rawRules := range firewallRules {
				values := mapToResourceData(rawRules.(map[string]interface{}))

				ifIndex := values.Get("vpc_router_interface_index").(int)
				direction := values.Get("direction").(string)

				// clear rules
				if vpcRouter.HasFirewall() && len(vpcRouter.Settings.Router.Firewall.Config) > ifIndex {
					switch direction {
					case "send":
						vpcRouter.Settings.Router.Firewall.Config[ifIndex].Send = nil
					case "receive":
						vpcRouter.Settings.Router.Firewall.Config[ifIndex].Receive = nil
					}
				}

				if rawExpressions, ok := values.GetOk("expressions"); ok {
					expressions := rawExpressions.([]interface{})
					for _, e := range expressions {
						exp := e.(map[string]interface{})

						allow := exp["allow"].(bool)
						protocol := exp["protocol"].(string)
						sourceNW := exp["source_nw"].(string)
						sourcePort := exp["source_port"].(string)
						destNW := exp["dest_nw"].(string)
						destPort := exp["dest_port"].(string)
						logging := exp["logging"].(bool)
						desc := ""
						if de, ok := exp["description"]; ok {
							desc = de.(string)
						}

						switch direction {
						case "send":
							vpcRouter.Settings.Router.AddFirewallRuleSend(ifIndex, allow, protocol, sourceNW, sourcePort, destNW, destPort, logging, desc)
						case "receive":
							vpcRouter.Settings.Router.AddFirewallRuleReceive(ifIndex, allow, protocol, sourceNW, sourcePort, destNW, destPort, logging, desc)
						}
					}
				}
			}
		}
	}

	if d.HasChange("l2tp") {
		// L2TP
		if l2tpSettings, ok := getListFromResource(d, "l2tp"); ok && len(l2tpSettings) > 0 {
			if l2tpSettings[0] != nil {
				values := mapToResourceData(l2tpSettings[0].(map[string]interface{}))
				l2tp := expandVPCRouterL2TP(values)
				vpcRouter.Settings.Router.EnableL2TPIPsecServer(l2tp.PreSharedSecret, l2tp.RangeStart, l2tp.RangeStop)
			}
		}
	}

	if d.HasChange("port_forwarding") {
		vpcRouter.Settings.Router.PortForwarding = nil
		if portForwardings, ok := getListFromResource(d, "port_forwarding"); ok && len(portForwardings) > 0 {
			for _, rawPortForwarding := range portForwardings {
				values := mapToResourceData(rawPortForwarding.(map[string]interface{}))
				pf := expandVPCRouterPortForwarding(values)
				vpcRouter.Settings.Router.AddPortForwarding(pf.Protocol, pf.GlobalPort, pf.PrivateAddress, pf.PrivatePort, pf.Description)
			}
		}
	}

	if d.HasChange("pptp") {
		if pptpSettings, ok := getListFromResource(d, "pptp"); ok && len(pptpSettings) > 0 {
			if pptpSettings[0] != nil {
				values := mapToResourceData(pptpSettings[0].(map[string]interface{}))
				pptp := expandVPCRouterPPTP(values)
				vpcRouter.Settings.Router.EnablePPTPServer(pptp.RangeStart, pptp.RangeStop)
			}
		}
	}

	if d.HasChange("site_to_site_vpn") {
		vpcRouter.Settings.Router.SiteToSiteIPsecVPN = nil
		if s2sSettings, ok := getListFromResource(d, "site_to_site_vpn"); ok && len(s2sSettings) > 0 {
			for _, rawS2s := range s2sSettings {
				values := mapToResourceData(rawS2s.(map[string]interface{}))
				s2s := expandVPCRouterSiteToSiteIPsecVPN(values)
				vpcRouter.Settings.Router.AddSiteToSiteIPsecVPN(s2s.LocalPrefix, s2s.Peer, s2s.PreSharedSecret, s2s.RemoteID, s2s.Routes)
			}
		}
	}

	if d.HasChange("static_nat") {
		vpcRouter.Settings.Router.StaticNAT = nil
		if staticNATSettings, ok := getListFromResource(d, "static_nat"); ok && len(staticNATSettings) > 0 {
			for _, rawStaticNAT := range staticNATSettings {
				values := mapToResourceData(rawStaticNAT.(map[string]interface{}))
				staticNAT := expandVPCRouterStaticNAT(values)
				vpcRouter.Settings.Router.AddStaticNAT(staticNAT.GlobalAddress, staticNAT.PrivateAddress, staticNAT.Description)
			}
		}
	}

	if d.HasChange("static_route") {
		vpcRouter.Settings.Router.StaticRoutes = nil
		if staticRoutes, ok := getListFromResource(d, "static_route"); ok && len(staticRoutes) > 0 {
			for _, rawStaticRoute := range staticRoutes {
				values := mapToResourceData(rawStaticRoute.(map[string]interface{}))
				staticRoute := expandVPCRouterStaticRoute(values)
				vpcRouter.Settings.Router.AddStaticRoute(staticRoute.Prefix, staticRoute.NextHop)
			}
		}
	}

	if d.HasChange("user") {
		vpcRouter.Settings.Router.RemoteAccessUsers = nil
		if users, ok := getListFromResource(d, "user"); ok && len(users) > 0 {
			for _, rawUser := range users {
				values := mapToResourceData(rawUser.(map[string]interface{}))
				user := expandVPCRouterRemoteAccessUser(values)
				vpcRouter.Settings.Router.AddRemoteAccessUser(user.UserName, user.Password)
			}
		}
	}

	vpcRouter, err = client.VPCRouter.Update(vpcRouter.ID, vpcRouter)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud VPCRouter resource: %s", err)
	}
	if _, err := client.VPCRouter.Config(vpcRouter.ID); err != nil {
		return fmt.Errorf("Error updating SakuraCloud VPCRouter settings: %s", err)
	}

	if isNeedRestart {
		_, err = client.VPCRouter.Boot(vpcRouter.ID)
		if err != nil {
			return fmt.Errorf("Failed to boot SakuraCloud VPCRouter resource: %s", err)
		}

		err = client.VPCRouter.SleepUntilUp(vpcRouter.ID, client.DefaultTimeoutDuration)
		if err != nil {
			return fmt.Errorf("Failed to boot SakuraCloud VPCRouter resource: %s", err)
		}
	}

	return resourceSakuraCloudVPCRouterRead(d, meta)
}

func resourceSakuraCloudVPCRouterDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Servers: %s", err)
	}

	if vpcRouter.Instance.IsUp() {
		// power API lock
		lockKey := getVPCRouterPowerAPILockKey(vpcRouter.ID)
		sakuraMutexKV.Lock(lockKey)
		defer sakuraMutexKV.Unlock(lockKey)

		err = nil
		for i := 0; i < 10; i++ {
			vpcRouter, err = client.VPCRouter.Read(vpcRouter.ID)
			if err != nil {
				return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
			}
			if vpcRouter.Instance.IsDown() {
				err = nil
				break
			}
			err = handleShutdown(client.VPCRouter, vpcRouter.ID, d, 60*time.Second)
		}
		if err != nil {
			return fmt.Errorf("Error stopping SakuraCloud VPCRouter resource: %s", err)
		}
	}

	_, err = client.VPCRouter.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud VPCRouter resource: %s", err)
	}

	return nil
}

func getVPCRouterPowerAPILockKey(id int64) string {
	return fmt.Sprintf(vpcRouterPowerAPILockKey, id)
}
