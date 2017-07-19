package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
)

func dataSourceSakuraCloudServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudServerRead,

		Schema: map[string]*schema.Schema{
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
			"core": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"base_interface": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Use field 'nic' instead",
			},
			"nic": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cdrom_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"additional_interfaces": {
				Type:       schema.TypeList,
				Computed:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: "Use field 'additional_nics' instead",
			},
			"additional_nics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"packet_filter_ids": {
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
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateStringInWord([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"macaddresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"base_nw_ipaddress": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Use field 'ipaddress' instead",
			},
			"ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"base_nw_dns_servers": {
				Type:       schema.TypeList,
				Computed:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: "Use field 'dns_servers' instead",
			},
			"dns_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"base_nw_gateway": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Use field 'gateway' instead",
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"base_nw_address": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Use field 'nw_address' instead",
			},
			"nw_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"base_nw_mask_len": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Use field 'nw_mask_len' instead",
			},
			"nw_mask_len": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSakuraCloudServerRead(d *schema.ResourceData, meta interface{}) error {
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
			client.Server.FilterBy(key, f)
		}
	}

	res, err := client.Server.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Server resource: %s", err)
	}
	if res == nil || res.Count == 0 {
		return nil
		//return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
	}
	server := res.Servers[0]

	return setServerResourceData(d, client, &server)
}
