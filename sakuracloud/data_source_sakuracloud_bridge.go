package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceSakuraCloudBridge() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudBridgeRead,

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
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_ids": {
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
		},
	}
}

func dataSourceSakuraCloudBridgeRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	//filters
	if rawFilter, filterOk := d.GetOk("filter"); filterOk {
		filters := expandFilters(rawFilter)
		for key, f := range filters {
			client.Bridge.FilterBy(key, f)
		}
	}

	res, err := client.Bridge.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Bridge resource: %s", err)
	}
	if res == nil || res.Count == 0 {
		return nil
		//return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
	}
	data := res.Bridges[0]
	return setBridgeResourceData(d, client, &data)
}
