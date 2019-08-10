package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/nfs"
)

func dataSourceSakuraCloudNFS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudNFSRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_mask_len": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_route": {
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

func dataSourceSakuraCloudNFSRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	searcher := sacloud.NewNFSOp(client)
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
		return fmt.Errorf("could not find SakuraCloud NFS resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.NFS) == 0 {
		return filterNoResultErr()
	}

	targets := res.NFS
	d.SetId(targets[0].ID.String())
	return setNFSV2ResourceData(ctx, d, client, targets[0])
}

func setNFSV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.NFS) error {
	if data.Availability.IsFailed() {
		d.SetId("")
		return fmt.Errorf("got unexpected state: NFS[%d].Availability is failed", data.ID)
	}

	var plan string
	var size int

	planInfo, err := nfs.GetPlanInfo(ctx, sacloud.NewNoteOp(client), data.PlanID)
	if err != nil {
		return err
	}
	switch planInfo.DiskPlanID {
	case types.NFSPlans.HDD:
		plan = "hdd"
	case types.NFSPlans.SSD:
		plan = "ssd"
	}
	size = int(planInfo.Size)

	setPowerManageTimeoutValueToState(d)
	d.Set("zone", client.Zone)
	return setResourceData(d, map[string]interface{}{
		"switch_id":     data.SwitchID.String(),
		"ipaddress":     data.IPAddresses[0],
		"nw_mask_len":   data.NetworkMaskLen,
		"default_route": data.DefaultRoute,
		"plan":          plan,
		"size":          size,
		"name":          data.Name,
		"icon_id":       data.IconID.String(),
		"description":   data.Description,
		"tags":          data.Tags,
	})

}
