package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
)

func dataSourceSakuraCloudInternet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudInternetRead,

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
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nw_mask_len": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"band_width": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"switch_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nw_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_gateway": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_min_ipaddress": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_max_ipaddress": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_ipaddresses": &schema.Schema{
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

func dataSourceSakuraCloudInternetRead(d *schema.ResourceData, meta interface{}) error {
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
			client.Internet.FilterBy(key, f)
		}
	}

	res, err := client.Internet.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Internet resource: %s", err)
	}
	if res == nil || res.Count == 0 {
		return nil
		//return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
	}
	internet := res.Internet[0]

	d.SetId(internet.ID)
	d.Set("name", internet.Name)
	d.Set("description", internet.Description)
	d.Set("tags", internet.Tags)
	d.Set("zone", client.Zone)

	d.Set("nw_mask_len", internet.NetworkMaskLen)
	d.Set("band_width", internet.BandWidthMbps)

	sw, err := client.Switch.Read(internet.Switch.ID)
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Switch resource: %s", err)
	}

	d.Set("switch_id", sw.ID)
	d.Set("nw_address", sw.Subnets[0].NetworkAddress)
	d.Set("nw_gateway", sw.Subnets[0].DefaultRoute)
	d.Set("nw_min_ipaddress", sw.Subnets[0].IPAddresses.Min)
	d.Set("nw_max_ipaddress", sw.Subnets[0].IPAddresses.Max)

	ipList, err := sw.GetIPAddressList()
	if err != nil {
		return fmt.Errorf("Error reading Switch resource(IPAddresses): %s", err)
	}
	d.Set("nw_ipaddresses", ipList)

	if sw.ServerCount > 0 {
		servers, err := client.Switch.GetServers(sw.ID)
		if err != nil {
			return fmt.Errorf("Couldn't find SakuraCloud Servers( is connected Switch): %s", err)
		}
		d.Set("server_ids", flattenServers(servers))
	} else {
		d.Set("server_ids", []string{})
	}

	return nil
}
