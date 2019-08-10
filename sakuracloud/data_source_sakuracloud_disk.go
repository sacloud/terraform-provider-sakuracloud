package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func dataSourceSakuraCloudDisk() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudDiskRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connector": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_archive_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_disk_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"server_id": {
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

func dataSourceSakuraCloudDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	searcher := sacloud.NewDiskOp(client)
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
		return fmt.Errorf("could not find SakuraCloud Disk resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.Disks) == 0 {
		return filterNoResultErr()
	}

	targets := res.Disks
	d.SetId(targets[0].ID.String())
	return setDiskV2ResourceData(ctx, d, client, targets[0])
}

func setDiskV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Disk) error {
	var plan string
	switch data.DiskPlanID {
	case types.DiskPlans.SSD:
		plan = "ssd"
	case types.DiskPlans.HDD:
		plan = "hdd"
	}

	var sourceDiskID, sourceArchiveID, serverID string
	if !data.SourceDiskID.IsEmpty() {
		sourceDiskID = data.SourceDiskID.String()
	}
	if !data.SourceArchiveID.IsEmpty() {
		sourceArchiveID = data.SourceArchiveID.String()
	}
	if !data.ServerID.IsEmpty() {
		serverID = data.ServerID.String()
	}

	setPowerManageTimeoutValueToState(d)

	return setResourceData(d, map[string]interface{}{
		"name":              data.Name,
		"plan":              plan,
		"source_disk_id":    sourceDiskID,
		"source_archive_id": sourceArchiveID,
		"connector":         data.Connection.String(),
		"size":              data.GetSizeGB(),
		"icon_id":           data.IconID.String(),
		"description":       data.Description,
		"tags":              data.Tags,
		"server_id":         serverID,
		"zone":              getV2Zone(d, client),
	})
}
