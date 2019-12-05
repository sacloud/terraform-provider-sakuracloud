package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func dataSourceSakuraCloudInternet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudInternetRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name": {
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
			"nw_mask_len": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"band_width": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"switch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nw_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipaddresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ipv6_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_prefix_len": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ipv6_nw_address": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceSakuraCloudInternetRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	searcher := sacloud.NewInternetOp(client)

	findCondition := &sacloud.FindCondition{
		Count: defaultSearchLimit,
	}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, zone, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud Internet resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.Internet) == 0 {
		return filterNoResultErr()
	}

	targets := res.Internet
	d.SetId(targets[0].ID.String())
	return setInternetResourceData(ctx, d, client, targets[0])
}
