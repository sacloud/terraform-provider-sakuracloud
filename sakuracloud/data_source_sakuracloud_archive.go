package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/ostype"
	"github.com/sacloud/libsacloud/v2/utils/archive"
)

func dataSourceSakuraCloudArchive() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudArchiveRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"os_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringInSlice(ostype.OSTypeShortNames, false),
				ConflictsWith: []string{"filters"},
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

func dataSourceSakuraCloudArchiveRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	searcher := sacloud.NewArchiveOp(client)

	ctx := context.Background()
	zone := getV2Zone(d, client)

	var data *sacloud.Archive
	if osType, ok := d.GetOk("os_type"); ok {
		strOSType := osType.(string)
		if strOSType != "" {
			res, err := archive.FindByOSType(ctx, searcher, zone, ostype.StrToOSType(strOSType))
			if err != nil {
				return err
			}
			data = res
		}
	} else {
		findCondition := &sacloud.FindCondition{
			Count: defaultSearchLimit,
		}
		if rawFilter, ok := d.GetOk(filterAttrName); ok {
			findCondition.Filter = expandSearchFilter(rawFilter)
		}

		res, err := searcher.Find(ctx, zone, findCondition)
		if err != nil {
			return fmt.Errorf("could not find SakuraCloud Archive resource: %s", err)
		}
		if res == nil || res.Count == 0 {
			return filterNoResultErr()
		}

		targets := res.Archives
		if len(targets) == 0 {
			return filterNoResultErr()
		}
		data = targets[0]
	}

	if data != nil {
		d.SetId(data.ID.String())
		return setResourceData(d, map[string]interface{}{
			"name":        data.Name,
			"size":        data.GetSizeGB(),
			"icon_id":     data.IconID.String(),
			"description": data.Description,
			"tags":        data.Tags,
			"zone":        getV2Zone(d, client),
		})
	}
	return nil
}
