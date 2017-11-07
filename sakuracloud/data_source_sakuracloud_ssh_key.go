package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceSakuraCloudSSHKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudSSHKeyRead,

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
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSakuraCloudSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	//filters
	if rawFilter, filterOk := d.GetOk("filter"); filterOk {
		filters := expandFilters(rawFilter)
		for key, f := range filters {
			client.SSHKey.FilterBy(key, f)
		}
	}

	res, err := client.SSHKey.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud SSHKey resource: %s", err)
	}
	if res == nil || res.Count == 0 {
		return nil
		//return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
	}
	key := res.SSHKeys[0]

	return setSSHKeyResourceData(d, client, &key)
}
