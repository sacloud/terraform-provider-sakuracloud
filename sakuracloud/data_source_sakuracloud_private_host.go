package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func dataSourceSakuraCloudPrivateHost() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudPrivateHostRead,

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
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assigned_core": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"assigned_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"tk1a"}),
			},
		},
	}
}

func dataSourceSakuraCloudPrivateHostRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	searcher := sacloud.NewPrivateHostOp(client)
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
		return fmt.Errorf("could not find SakuraCloud PrivateHost resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.PrivateHosts) == 0 {
		return filterNoResultErr()
	}

	targets := res.PrivateHosts
	d.SetId(targets[0].ID.String())
	return setPrivateHostV2ResourceData(ctx, d, client, targets[0])
}

func setPrivateHostV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.PrivateHost) error {
	setPowerManageTimeoutValueToState(d)

	return setResourceData(d, map[string]interface{}{
		"name":            data.Name,
		"icon_id":         data.IconID.String(),
		"description":     data.Description,
		"tags":            data.Tags,
		"hostname":        data.HostName,
		"assigned_core":   data.AssignedCPU,
		"assigned_memory": data.GetAssignedMemoryGB(),
		"zone":            client.Zone,
	})
}
