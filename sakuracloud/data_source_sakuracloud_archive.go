package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/sacloud/ostype"
)

func dataSourceSakuraCloudArchive() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudArchiveRead,

		Schema: map[string]*schema.Schema{
			"os_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validateStringInWord(ostype.OSTypeShortNames),
				ConflictsWith: []string{"filter"},
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
				ConflictsWith: []string{"os_type"},
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
				ValidateFunc: validateStringInWord([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		},
	}
}

func dataSourceSakuraCloudArchiveRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	var archive *sacloud.Archive

	if osType, ok := d.GetOk("os_type"); ok {
		strOSType := osType.(string)
		if strOSType != "" {

			res, err := client.Archive.FindByOSType(strToOSType(strOSType))
			if err != nil {
				return fmt.Errorf("Couldn't find SakuraCloud Archive resource: %s", err)
			}
			archive = res
		}
	} else {

		//filters
		if rawFilter, filterOk := d.GetOk("filter"); filterOk {
			filters := expandFilters(rawFilter)
			for key, f := range filters {
				client.Archive.FilterBy(key, f)
			}
		}

		res, err := client.Archive.Find()
		if err != nil {
			return fmt.Errorf("Couldn't find SakuraCloud Archive resource: %s", err)
		}
		if res == nil || res.Count == 0 {
			return nil
			//return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
		}
		archive = &res.Archives[0]
	}

	if archive != nil {

		d.SetId(archive.GetStrID())
		d.Set("name", archive.Name)
		d.Set("size", toSizeGB(archive.SizeMB))
		d.Set("icon_id", archive.GetIconStrID())
		d.Set("description", archive.Description)
		d.Set("tags", archive.Tags)

		d.Set("zone", client.Zone)
	}

	return nil
}

func strToOSType(strType string) ostype.ArchiveOSTypes {
	return ostype.StrToOSType(strType)
}
