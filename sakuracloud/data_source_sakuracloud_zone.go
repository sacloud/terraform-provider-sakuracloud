package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func dataSourceSakuraCloudZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudZoneRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateZone(allowZones),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceSakuraCloudZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	zoneOp := sacloud.NewZoneOp(client)
	ctx := context.Background()

	zoneSlug := client.Zone
	if v, ok := d.GetOk("name"); ok {
		zoneSlug = v.(string)
	}

	res, err := zoneOp.Find(ctx, &sacloud.FindCondition{
		Count: defaultSearchLimit,
	})
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud Zone resource: %s", err)
	}
	if res == nil || len(res.Zones) == 0 {
		return filterNoResultErr()
	}
	var data *sacloud.Zone

	for _, z := range res.Zones {
		if z.Name == zoneSlug {
			data = z
			break
		}
	}
	if data == nil {
		return filterNoResultErr()
	}

	d.SetId(data.ID.String())
	return setResourceData(d, map[string]interface{}{
		"name":        data.Name,
		"zone_id":     data.ID.String(),
		"description": data.Description,
		"region_id":   data.Region.ID.String(),
		"region_name": data.Region.Name,
		"dns_servers": data.Region.NameServers,
	})
}
