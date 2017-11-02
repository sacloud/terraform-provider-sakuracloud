package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceSakuraCloudDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudDatabaseRead,

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
			"plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allow_networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"backup_rotate": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backup_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"switch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipaddress1": {
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
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateStringInWord([]string{"tk1a", "is1b"}),
			},
		},
	}
}

func dataSourceSakuraCloudDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	//filters
	if rawFilter, filterOk := d.GetOk("filter"); filterOk {
		filters := expandFilters(rawFilter)
		for key, f := range filters {
			client.Database.FilterBy(key, f)
		}
	}

	res, err := client.Database.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Database resource: %s", err)
	}
	if res == nil || res.Count == 0 {
		return nil
		//return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
	}
	data := res.Databases[0]
	return setDatabaseResourceData(d, client, &data)
}
