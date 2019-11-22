package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func dataSourceSakuraCloudIcon() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudIconRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// TODO 廃止
			//"body": {
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceSakuraCloudIconRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	searcher := sacloud.NewIconOp(client)
	ctx := context.Background()

	findCondition := &sacloud.FindCondition{
		Count: defaultSearchLimit,
	}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud Icon resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.Icons) == 0 {
		return filterNoResultErr()
	}

	targets := res.Icons
	icon := targets[0]

	d.SetId(icon.ID.String())
	return setIconResourceData(ctx, d, client, icon)
}
