package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func dataSourceSakuraCloudSSHKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudSSHKeyRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{excludeTags: true}),
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
	searcher := sacloud.NewSSHKeyOp(client)
	ctx := context.Background()

	findCondition := &sacloud.FindCondition{
		Count: defaultSearchLimit,
	}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud SSHKey resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.SSHKeys) == 0 {
		return filterNoResultErr()
	}

	targets := res.SSHKeys
	d.SetId(targets[0].ID.String())
	return setSSHKeyV2ResourceData(ctx, d, client, targets[0])
}

func setSSHKeyV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.SSHKey) error {
	return setResourceData(d, map[string]interface{}{
		"name":        data.Name,
		"public_key":  data.PublicKey,
		"fingerprint": data.Fingerprint,
		"description": data.Description,
	})
}
