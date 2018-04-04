package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/sacloud"
)

func dataSourceSakuraCloudCDROM() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudCDROMRead,

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
			"size": {
				Type:     schema.TypeInt,
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

func dataSourceSakuraCloudCDROMRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	//filters
	if rawFilter, filterOk := d.GetOk("filter"); filterOk {
		filters := expandFilters(rawFilter)
		for key, f := range filters {
			client.CDROM.FilterBy(key, f)
		}
	}

	res, err := client.CDROM.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud CDROM resource: %s", err)
	}
	if res == nil || res.Count == 0 {
		return filterNoResultErr()
	}

	var data *sacloud.CDROM
	targets := res.CDROMs

	if rawNameSelector, ok := d.GetOk("name_selectors"); ok {
		selectors := expandStringList(rawNameSelector.([]interface{}))
		var filtered []sacloud.CDROM
		for _, a := range targets {
			if hasNames(&a, selectors) {
				filtered = append(filtered, a)
			}
		}
		targets = filtered
	}
	if rawTagSelector, ok := d.GetOk("tag_selectors"); ok {
		selectors := expandStringList(rawTagSelector.([]interface{}))
		var filtered []sacloud.CDROM
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
	return setCDROMResourceData(d, client, data)
}
