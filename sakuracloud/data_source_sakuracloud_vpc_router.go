package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/sacloud"
)

func dataSourceSakuraCloudVPCRouter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudVPCRouterRead,

		Schema: map[string]*schema.Schema{
			"name_selectors": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tag_selectors": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
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
			"interface": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterInterfaceDataSchema(),
				},
			},
			"dhcp_server": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterDHCPServerDataSchema(),
				},
			},
			"dhcp_static_mapping": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterDHCPStaticMappingDataSchema(),
				},
			},
			"firewall": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterFirewallDataSchema(),
				},
			},
			"l2tp": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterL2TPDataSchema(),
				},
			},
			"port_forwarding": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterPortForwardingDataSchema(),
				},
			},
			"pptp": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterPPTPDataSchema(),
				},
			},
			"site_to_site_vpn": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterS2SDataSchema(),
				},
			},
			"static_nat": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterStaticNATDataSchema(),
				},
			},
			"static_route": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterStaticRouteDataSchema(),
				},
			},
			"user": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: vpcRouterUserDataSchema(),
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

	//filters
	if rawFilter, filterOk := d.GetOk("filter"); filterOk {
		filters := expandFilters(rawFilter)
		for key, f := range filters {
			client.VPCRouter.FilterBy(key, f)
		}
	}

	res, err := client.VPCRouter.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}
	if res == nil || res.Count == 0 {
		return filterNoResultErr()
	}
	var data *sacloud.VPCRouter
	targets := res.VPCRouters

	if rawNameSelector, ok := d.GetOk("name_selectors"); ok {
		selectors := expandStringList(rawNameSelector.([]interface{}))
		var filtered []sacloud.VPCRouter
		for _, a := range targets {
			if hasNames(&a, selectors) {
				filtered = append(filtered, a)
			}
		}
		targets = filtered
	}
	if rawTagSelector, ok := d.GetOk("tag_selectors"); ok {
		selectors := expandStringList(rawTagSelector.([]interface{}))
		var filtered []sacloud.VPCRouter
		for _, a := range targets {
			if hasTags(&a, selectors) {
				filtered = append(filtered, a)
			}
		}
		targets = filtered
	}

	if len(targets) == 0 {
		return filterNoResultErr()
	}
	data = &targets[0]

	d.SetId(data.GetStrID())
	return setVPCRouterResourceData(d, client, data)
}
