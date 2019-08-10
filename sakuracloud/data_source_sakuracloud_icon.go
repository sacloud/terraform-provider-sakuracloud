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
	d.SetId(targets[0].ID.String())
	return setIconV2ResourceData(ctx, d, client, targets[0])
}

func setIconV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Icon) error {
	return setResourceData(d, map[string]interface{}{
		"name": data.Name,
		"tags": data.Tags,
		"url":  data.URL,
	})
}
