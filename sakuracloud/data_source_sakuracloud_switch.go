package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func dataSourceSakuraCloudSwitch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudSwitchRead,

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
			"bridge_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_ids": {
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

func dataSourceSakuraCloudSwitchRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	searcher := sacloud.NewSwitchOp(client)
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
		return fmt.Errorf("could not find SakuraCloud Switch resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.Switches) == 0 {
		return filterNoResultErr()
	}

	targets := res.Switches
	d.SetId(targets[0].ID.String())
	return setSwitchV2ResourceData(ctx, d, client, targets[0])
}

func setSwitchV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Switch) error {
	zone := getV2Zone(d, client)

	var serverIDs []string
	if data.ServerCount > 0 {
		swOp := sacloud.NewSwitchOp(client)
		searched, err := swOp.GetServers(ctx, zone, data.ID)
		if err != nil {
			return fmt.Errorf("could not find SakuraCloud Servers: switch[%s]", err)
		}
		for _, s := range searched.Servers {
			serverIDs = append(serverIDs, s.ID.String())
		}
	}

	setPowerManageTimeoutValueToState(d)
	return setResourceData(d, map[string]interface{}{
		"name":        data.Name,
		"icon_id":     data.IconID.String(),
		"description": data.Description,
		"tags":        data.Tags,
		"bridge_id":   data.BridgeID.String(),
		"server_ids":  serverIDs,
		"zone":        zone,
	})

}
