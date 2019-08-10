package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func dataSourceSakuraCloudLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudLoadBalancerRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_double": {
				Type:     schema.TypeBool,
				Computed: true,
				Removed:  "Use field 'high_availability' instead",
			},
			"high_availability": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"plan": {
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
			"nw_mask_len": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_route": {
				Type:     schema.TypeString,
				Computed: true,
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
			"vips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"delay_loop": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sorry_server": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipaddress": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"check_protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"check_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"check_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
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

func dataSourceSakuraCloudLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	searcher := sacloud.NewLoadBalancerOp(client)
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
		return fmt.Errorf("could not find SakuraCloud LoadBalancer resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.LoadBalancers) == 0 {
		return filterNoResultErr()
	}

	targets := res.LoadBalancers
	d.SetId(targets[0].ID.String())
	return setLoadBalancerV2ResourceData(ctx, d, client, targets[0])
}

func setLoadBalancerV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.LoadBalancer) error {
	if data.Availability.IsFailed() {
		d.SetId("")
		return fmt.Errorf("got unexpected state: LoadBalancer[%d].Availability is failed", data.ID)
	}

	var ha bool
	var ipaddress1, ipaddress2 string
	ipaddress1 = data.IPAddresses[0]
	if len(data.IPAddresses) > 1 {
		ha = true
		ipaddress2 = data.IPAddresses[1]
	}

	var plan string
	switch data.PlanID {
	case types.LoadBalancerPlans.Standard:
		plan = "standard"
	case types.LoadBalancerPlans.Premium:
		plan = "highspec"
	}

	var vips []interface{}
	for _, v := range data.VirtualIPAddresses {
		vip := map[string]interface{}{
			"vip":          v.VirtualIPAddress,
			"port":         v.Port.Int(),
			"delay_loop":   v.DelayLoop.Int(),
			"sorry_server": v.SorryServer,
		}
		var servers []interface{}
		for _, server := range v.Servers {
			s := map[string]interface{}{}
			s["ipaddress"] = server.IPAddress
			s["check_protocol"] = server.HealthCheck.Protocol
			s["check_path"] = server.HealthCheck.Path
			s["check_status"] = server.HealthCheck.ResponseCode.String()
			s["enabled"] = server.Enabled.Bool()
			servers = append(servers, s)
		}
		vip["servers"] = servers
		vips = append(vips, vip)
	}

	setPowerManageTimeoutValueToState(d)
	return setResourceData(d, map[string]interface{}{
		"switch_id":         data.SwitchID.String(),
		"vrid":              data.VRID,
		"plan":              plan,
		"high_availability": ha,
		"ipaddress1":        ipaddress1,
		"ipaddress2":        ipaddress2,
		"nw_mask_len":       data.NetworkMaskLen,
		"default_route":     data.DefaultRoute,
		"name":              data.Name,
		"icon_id":           data.IconID.String(),
		"description":       data.Description,
		"tags":              data.Tags,
		"vips":              vips,
		"zone":              getV2Zone(d, client),
	})
}
