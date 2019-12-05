// Copyright 2016-2019 terraform-provider-sakuracloud authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/setup"
)

func resourceSakuraCloudDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudDiskCreate,
		Read:   resourceSakuraCloudDiskRead,
		Update: resourceSakuraCloudDiskUpdate,
		Delete: resourceSakuraCloudDiskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plan": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "ssd",
				ValidateFunc: validation.StringInSlice([]string{"ssd", "hdd"}, false),
			},
			"connector": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  types.DiskConnections.VirtIO,
				ValidateFunc: validation.StringInSlice([]string{
					types.DiskConnections.VirtIO.String(),
					types.DiskConnections.IDE.String(),
				}, false),
			},
			"source_archive_id": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"source_disk_id"},
				ValidateFunc:  validateSakuracloudIDType,
			},
			"source_disk_id": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"source_archive_id"},
				ValidateFunc:  validateSakuracloudIDType,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  20,
			},
			"distant_from": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Computed: true, //ReadOnly
			},
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
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

func resourceSakuraCloudDiskCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	diskOp := sacloud.NewDiskOp(client)

	var planID types.ID
	plan := d.Get("plan").(string)
	switch d.Get("plan").(string) {
	case "ssd":
		planID = types.DiskPlans.SSD
	case "hdd":
		planID = types.DiskPlans.HDD
	default:
		return fmt.Errorf("invalid disk plan [%s]", plan)
	}
	distantFrom := expandSakuraCloudIDs(d, "distant_from")

	ops := &sacloud.DiskCreateRequest{
		DiskPlanID:      planID,
		Connection:      types.EDiskConnection(d.Get("connector").(string)),
		SourceDiskID:    expandSakuraCloudID(d, "source_disk_id"),
		SourceArchiveID: expandSakuraCloudID(d, "source_archive_id"),
		SizeMB:          toSizeMB(d.Get("size").(int)),
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		Tags:            expandTags(d),
		IconID:          expandSakuraCloudID(d, "icon_id"),
	}

	diskBuilder := &setup.RetryableSetup{
		IsWaitForCopy: true,
		Create: func(ctx context.Context, zone string) (accessor.ID, error) {
			return diskOp.Create(ctx, zone, ops, distantFrom)
		},
		Read: func(ctx context.Context, zone string, id types.ID) (interface{}, error) {
			return diskOp.Read(ctx, zone, id)
		},
		Delete: func(ctx context.Context, zone string, id types.ID) error {
			return diskOp.Delete(ctx, zone, id)
		},
		RetryCount: 3,
	}

	res, err := diskBuilder.Setup(ctx, zone)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Disk resource is failed: %s", err)
	}

	disk, ok := res.(*sacloud.Disk)
	if !ok {
		return fmt.Errorf("creating SakuraCloud Disk resource is failed: created resource is not a *sacloud.Disk")
	}

	d.SetId(disk.ID.String())
	return resourceSakuraCloudDiskRead(d, meta)
}

func resourceSakuraCloudDiskRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	diskOp := sacloud.NewDiskOp(client)

	disk, err := diskOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not find SakuraCloud Disk resource: %s", err)
	}

	return setDiskResourceData(ctx, d, client, disk)
}

func resourceSakuraCloudDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	diskOp := sacloud.NewDiskOp(client)

	disk, err := diskOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Disk resource: %s", err)
	}

	ops := &sacloud.DiskUpdateRequest{
		Connection:  types.EDiskConnection(d.Get("connector").(string)),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),
	}

	disk, err = diskOp.Update(ctx, zone, disk.ID, ops)
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Disk resource is failed: %s", err)
	}

	return resourceSakuraCloudDiskRead(d, meta)
}

func resourceSakuraCloudDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	diskOp := sacloud.NewDiskOp(client)

	disk, err := diskOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Disk resource: %s", err)
	}

	if err := diskOp.Delete(ctx, zone, disk.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud Disk resource is failed: %s", err)
	}
	return nil
}

func setDiskResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Disk) error {
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

	d.Set("name", data.Name)
	d.Set("plan", plan)
	d.Set("source_disk_id", sourceDiskID)
	d.Set("source_archive_id", sourceArchiveID)
	d.Set("connector", data.Connection.String())
	d.Set("size", data.GetSizeGB())
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	d.Set("server_id", serverID)
	d.Set("zone", getZone(d, client))
	return nil
}
