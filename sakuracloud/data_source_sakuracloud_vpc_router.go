package sakuracloud

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func dataSourceSakuraCloudVPCRouter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudVPCRouterRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipaddress1": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipaddress2": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"aliases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"icon_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"global_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"syslog_host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_connection": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"switch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipaddresses": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"nw_mask_len": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"dhcp_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface_index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"range_start": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"range_stop": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_servers": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
					},
				},
			},
			"dhcp_static_mappings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipaddress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"macaddress": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"firewall": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface_index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expressions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_nw": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dest_nw": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dest_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"allow": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"logging": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"l2tp": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pre_shared_secret": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
						"range_start": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"range_stop": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"port_forwardings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"global_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"private_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"pptp": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range_start": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"range_stop": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"site_to_site_vpn": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"peer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pre_shared_secret": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"routes": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"local_prefix": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
					},
				},
			},
			"static_nat": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"global_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"static_routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
					},
				},
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		},
	}
}

func dataSourceSakuraCloudVPCRouterRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	searcher := sacloud.NewVPCRouterOp(client)
	ctx := context.Background()
	zone := getV2Zone(d, client)

	findCondition := &sacloud.FindCondition{
		Count: defaultSearchLimit,
	}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, zone, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud VPCRouter resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.VPCRouters) == 0 {
		return filterNoResultErr()
	}

	targets := res.VPCRouters
	d.SetId(targets[0].ID.String())
	return setVPCRouterV2ResourceData(ctx, d, client, targets[0])
}

func setVPCRouterV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.VPCRouter) error {
	if data.Availability.IsFailed() {
		d.SetId("")
		return fmt.Errorf("got unexpected state: VPCRouter[%d].Availability is failed", data.ID)
	}

	//plan1
	var planName string
	switch data.PlanID {
	case types.VPCRouterPlans.Standard:
		planName = "standard"
	case types.VPCRouterPlans.Premium:
		planName = "premium"
	case types.VPCRouterPlans.HighSpec:
		planName = "highspec"
	}

	var globalAddress, switchID, vip, ipaddress1, ipaddress2 string
	var aliases []string
	var vrid int
	if data.PlanID == types.VPCRouterPlans.Standard {
		globalAddress = data.Interfaces[0].IPAddress
	} else {
		switchID = data.Interfaces[0].SwitchID.String()
		vip = data.Settings.Interfaces[0].VirtualIPAddress
		ipaddress1 = data.Settings.Interfaces[0].IPAddress[0]
		ipaddress2 = data.Settings.Interfaces[0].IPAddress[1]
		aliases = data.Settings.Interfaces[0].IPAliases
		vrid = data.Settings.VRID
		globalAddress = data.Settings.Interfaces[0].VirtualIPAddress
	}

	setPowerManageTimeoutValueToState(d)

	// interface
	var interfaces []map[string]interface{}
	if len(data.Interfaces) > 0 {
		for _, iface := range data.Settings.Interfaces {
			// find nic from data.Interfaces
			var nic *sacloud.VPCRouterInterface
			for _, n := range data.Interfaces {
				if iface.Index == n.Index {
					nic = n
					break
				}
			}

			if nic != nil {
				interfaces = append(interfaces, map[string]interface{}{
					"switch_id":   nic.SwitchID.String(),
					"vip":         iface.VirtualIPAddress,
					"ipaddresses": iface.IPAddress,
					"nw_mask_len": iface.NetworkMaskLen,
					"index":       iface.Index,
				})
			}
		}
	}

	var dhcpServers []map[string]interface{}
	for _, d := range data.Settings.DHCPServer {
		dhcpServers = append(dhcpServers, map[string]interface{}{
			"range_start":     d.RangeStart,
			"range_stop":      d.RangeStop,
			"interface_index": vpcRouterInterfaceNameToIndex(d.Interface),
			"dns_servers":     d.DNSServers,
		})
	}

	var staticMappings []map[string]interface{}
	for _, d := range data.Settings.DHCPStaticMapping {
		staticMappings = append(staticMappings, map[string]interface{}{
			"ipaddress":  d.IPAddress,
			"macaddress": d.MACAddress,
		})
	}

	var firewallRules []map[string]interface{}
	for i, configs := range data.Settings.Firewall {

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
					"source_nw":   rule.SourceNetwork,
					"source_port": rule.SourcePort,
					"dest_nw":     rule.DestinationNetwork,
					"dest_port":   rule.DestinationPort,
					"allow":       rule.Action.IsAllow(),
					"protocol":    rule.Protocol,
					"logging":     rule.Logging.Bool(),
					"description": rule.Description,
				}
				expressions = append(expressions, expression)
			}
			firewallRules = append(firewallRules, map[string]interface{}{
				"interface_index": i,
				"direction":       direction,
				"expressions":     expressions,
			})
		}
	}

	var l2tp []map[string]interface{}
	if data.Settings.L2TPIPsecServerEnabled.Bool() {
		s := data.Settings.L2TPIPsecServer
		l2tp = append(l2tp, map[string]interface{}{
			"pre_shared_secret": s.PreSharedSecret,
			"range_start":       s.RangeStart,
			"range_stop":        s.RangeStop,
		})
	}

	var portForwardings []map[string]interface{}
	for _, p := range data.Settings.PortForwarding {
		globalPort := p.GlobalPort.Int()
		privatePort := p.PrivatePort.Int()
		portForwardings = append(portForwardings, map[string]interface{}{
			"protocol":        string(p.Protocol),
			"global_port":     globalPort,
			"private_address": p.PrivateAddress,
			"private_port":    privatePort,
			"description":     p.Description,
		})
	}

	var pptp []map[string]interface{}
	if data.Settings.PPTPServerEnabled.Bool() {
		c := data.Settings.PPTPServer
		pptp = append(pptp, map[string]interface{}{
			"range_start": c.RangeStart,
			"range_stop":  c.RangeStop,
		})
	}

	var s2sSettings []map[string]interface{}
	for _, s := range data.Settings.SiteToSiteIPsecVPN {
		s2sSettings = append(s2sSettings, map[string]interface{}{
			"local_prefix":      s.LocalPrefix,
			"peer":              s.Peer,
			"pre_shared_secret": s.PreSharedSecret,
			"remote_id":         s.RemoteID,
			"routes":            s.Routes,
		})
	}

	var staticNATs []map[string]interface{}
	for _, s := range data.Settings.StaticNAT {
		staticNATs = append(staticNATs, map[string]interface{}{
			"global_address":  s.GlobalAddress,
			"private_address": s.PrivateAddress,
			"description":     s.Description,
		})
	}

	var staticRoutes []map[string]interface{}
	for _, s := range data.Settings.StaticRoute {
		staticRoutes = append(staticRoutes, map[string]interface{}{
			"prefix":   s.Prefix,
			"next_hop": s.NextHop,
		})
	}

	var users []map[string]interface{}
	for _, u := range data.Settings.RemoteAccessUsers {
		users = append(users, map[string]interface{}{
			"name":     u.UserName,
			"password": u.Password,
		})
	}

	return setResourceData(d, map[string]interface{}{
		"name":                 data.Name,
		"icon_id":              data.IconID.String(),
		"description":          data.Description,
		"tags":                 data.Tags,
		"plan":                 planName,
		"switch_id":            switchID,
		"global_address":       globalAddress,
		"vip":                  vip,
		"ipaddress1":           ipaddress1,
		"ipaddress2":           ipaddress2,
		"aliases":              aliases,
		"vrid":                 vrid,
		"syslog_host":          data.Settings.SyslogHost,
		"internet_connection":  data.Settings.InternetConnectionEnabled.Bool(),
		"interfaces":           interfaces,
		"dhcp_servers":         dhcpServers,
		"dhcp_static_mappings": staticMappings,
		"firewall":             firewallRules,
		"l2tp":                 l2tp,
		"pptp":                 pptp,
		"port_forwardings":     portForwardings,
		"site_to_site_vpn":     s2sSettings,
		"static_nat":           staticNATs,
		"static_routes":        staticRoutes,
		"users":                users,
		"zone":                 getV2Zone(d, client),
	})
}

func vpcRouterInterfaceNameToIndex(ifName string) int {
	strIndex := strings.Replace(ifName, "eth", "", -1)
	index, err := strconv.Atoi(strIndex)
	if err != nil {
		return -1
	}
	return index
}
