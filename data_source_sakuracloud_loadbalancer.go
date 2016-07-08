package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
)

func dataSourceSakuraCloudLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudLoadBalancerRead,

		Schema: map[string]*schema.Schema{
			"filter": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"values": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"VRID": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_double": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"plan": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipaddress1": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipaddress2": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_mask_len": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_route": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zone": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateStringInWord([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		},
	}
}

func dataSourceSakuraCloudLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	//filters
	if rawFilter, filterOk := d.GetOk("filter"); filterOk {
		filters := expandFilters(rawFilter)
		for key, f := range filters {
			client.LoadBalancer.FilterBy(key, f)
		}
	}

	res, err := client.LoadBalancer.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}
	if res == nil || res.Count == 0 {
		return nil
		//return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
	}
	loadBalancer := res.LoadBalancers[0]

	d.SetId(loadBalancer.ID)
	d.Set("switch_id", loadBalancer.Switch.ID)
	d.Set("VRID", loadBalancer.Remark.VRRP.VRID)
	if len(loadBalancer.Remark.Servers) > 1 {
		d.Set("is_double", true)
		d.Set("ipaddress1", loadBalancer.Remark.Servers[0].(map[string]interface{})["IPAddress"])
		d.Set("ipaddress2", loadBalancer.Remark.Servers[1].(map[string]interface{})["IPAddress"])
	} else {
		d.Set("is_double", false)
		d.Set("ipaddress1", loadBalancer.Remark.Servers[0].(map[string]interface{})["IPAddress"])
	}
	d.Set("nw_mask_len", loadBalancer.Remark.Network.NetworkMaskLen)
	d.Set("default_route", loadBalancer.Remark.Network.DefaultRoute)

	d.Set("name", loadBalancer.Name)
	d.Set("description", loadBalancer.Description)
	d.Set("tags", loadBalancer.Tags)
	d.Set("zone", client.Zone)

	return nil
}
