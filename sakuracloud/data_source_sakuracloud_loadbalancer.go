package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
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
	return setLoadBalancerResourceData(ctx, d, client, targets[0])
}
