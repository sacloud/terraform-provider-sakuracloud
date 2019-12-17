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
	"time"

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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(24 * time.Hour),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(24 * time.Hour),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "target SakuraCloud zone",
			},
		},
	}
}

func resourceSakuraCloudDiskCreate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	diskOp := sacloud.NewDiskOp(client)

	diskBuilder := &setup.RetryableSetup{
		IsWaitForCopy: true,
		Create: func(ctx context.Context, zone string) (accessor.ID, error) {
			return diskOp.Create(ctx, zone, expandDiskCreateRequest(d), expandSakuraCloudIDs(d, "distant_from"))
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
		return fmt.Errorf("creating SakuraCloud Disk is failed: %s", err)
	}

	disk, ok := res.(*sacloud.Disk)
	if !ok {
		return fmt.Errorf("creating SakuraCloud Disk is failed: created resource is not a *sacloud.Disk")
	}

	d.SetId(disk.ID.String())
	return resourceSakuraCloudDiskRead(d, meta)
}

func resourceSakuraCloudDiskRead(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	diskOp := sacloud.NewDiskOp(client)

	disk, err := diskOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Disk[%s]: %s", d.Id(), err)
	}

	return setDiskResourceData(ctx, d, client, disk)
}

func resourceSakuraCloudDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	diskOp := sacloud.NewDiskOp(client)

	disk, err := diskOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Disk[%s]: %s", d.Id(), err)
	}

	_, err = diskOp.Update(ctx, zone, disk.ID, expandDiskUpdateRequest(d))
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Disk[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudDiskRead(d, meta)
}

func resourceSakuraCloudDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	diskOp := sacloud.NewDiskOp(client)

	disk, err := diskOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Disk[%s]: %s", d.Id(), err)
	}

	if err := waitForDeletionByDiskID(ctx, client, zone, disk.ID); err != nil {
		return fmt.Errorf("waiting deletion is failed: Disk[%s] still used by Servers: %s", disk.ID, err)
	}

	if err := diskOp.Delete(ctx, zone, disk.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud Disk[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setDiskResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Disk) error {
	d.Set("name", data.Name)                                  // nolint
	d.Set("plan", flattenDiskPlan(data))                      // nolint
	d.Set("source_disk_id", data.SourceDiskID.String())       // nolint
	d.Set("source_archive_id", data.SourceArchiveID.String()) // nolint
	d.Set("connector", data.Connection.String())              // nolint
	d.Set("size", data.GetSizeGB())                           // nolint
	d.Set("icon_id", data.IconID.String())                    // nolint
	d.Set("description", data.Description)                    // nolint
	d.Set("server_id", data.ServerID.String())                // nolint
	d.Set("zone", getZone(d, client))                         // nolint
	return d.Set("tags", data.Tags)
}
