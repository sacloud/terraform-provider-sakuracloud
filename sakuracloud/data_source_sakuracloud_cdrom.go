package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func dataSourceSakuraCloudCDROM() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudCDROMRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
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
	searcher := sacloud.NewCDROMOp(client)
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
		return fmt.Errorf("could not find SakuraCloud CDROM resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.CDROMs) == 0 {
		return filterNoResultErr()
	}

	targets := res.CDROMs
	d.SetId(targets[0].ID.String())
	return setCDROMV2ResourceData(ctx, d, client, targets[0])
}

func setCDROMV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.CDROM) error {
	return setResourceData(d, map[string]interface{}{
		"name":        data.Name,
		"size":        data.GetSizeGB(),
		"icon_id":     data.IconID.String(),
		"description": data.Description,
		"tags":        data.Tags,
		"zone":        getV2Zone(d, client),
	})
}
