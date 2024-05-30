// Copyright 2016-2023 terraform-provider-sakuracloud authors
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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/accessor"
	"github.com/sacloud/iaas-api-go/helper/cleanup"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/setup"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudDisk() *schema.Resource {
	resourceName := "disk"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudDiskCreate,
		ReadContext:   resourceSakuraCloudDiskRead,
		UpdateContext: resourceSakuraCloudDiskUpdate,
		DeleteContext: resourceSakuraCloudDiskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(24 * time.Hour),
			Update: schema.DefaultTimeout(24 * time.Hour),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"plan": schemaResourcePlan(resourceName, types.DiskPlanNameMap[types.DiskPlans.SSD], types.DiskPlanStrings),
			"connector": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          types.DiskConnections.VirtIO,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.DiskConnectionStrings, false)),
				Description: desc.Sprintf(
					"The name of the disk connector. This must be one of [%s]",
					types.DiskConnectionStrings,
				),
			},
			"source_archive_id": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				ConflictsWith:    []string{"source_disk_id"},
				ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
				Description: desc.Sprintf(
					"The id of the source archive. %s",
					desc.Conflicts("source_disk_id"),
				),
			},
			"source_disk_id": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				ConflictsWith:    []string{"source_archive_id"},
				ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
				Description: desc.Sprintf(
					"The id of the source disk. %s",
					desc.Conflicts("source_archive_id"),
				),
			},
			"size": schemaResourceSize(resourceName, 20),
			"distant_from": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				ForceNew:    true,
				Description: "A list of disk id. The disk will be located to different storage from these disks",
			},
			"server_id":   schemaDataSourceServerID(resourceName),
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"zone":        schemaResourceZone(resourceName),
			"encryption_algorithm": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.DiskEncryptionAlgorithmStrings, false)),
				Description: desc.Sprintf(
					"The disk encryption algorithm. This must be one of [%s]",
					types.DiskEncryptionAlgorithmStrings,
				),
			},
		},
	}
}

func resourceSakuraCloudDiskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	diskOp := iaas.NewDiskOp(client)
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
		Options: &setup.Options{
			RetryCount: 3,
		},
	}

	res, err := diskBuilder.Setup(ctx, zone)
	if err != nil {
		return diag.Errorf("creating SakuraCloud Disk is failed: %s", err)
	}

	disk, ok := res.(*iaas.Disk)
	if !ok {
		return diag.Errorf("creating SakuraCloud Disk is failed: created resource is not a *iaas.Disk")
	}

	d.SetId(disk.ID.String())
	return resourceSakuraCloudDiskRead(ctx, d, meta)
}

func resourceSakuraCloudDiskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	diskOp := iaas.NewDiskOp(client)
	disk, err := diskOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud Disk[%s]: %s", d.Id(), err)
	}

	return setDiskResourceData(ctx, d, client, disk)
}

func resourceSakuraCloudDiskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	diskOp := iaas.NewDiskOp(client)
	disk, err := diskOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud Disk[%s]: %s", d.Id(), err)
	}

	_, err = diskOp.Update(ctx, zone, disk.ID, expandDiskUpdateRequest(d))
	if err != nil {
		return diag.Errorf("updating SakuraCloud Disk[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudDiskRead(ctx, d, meta)
}

func resourceSakuraCloudDiskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	diskOp := iaas.NewDiskOp(client)
	disk, err := diskOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud Disk[%s]: %s", d.Id(), err)
	}

	if err := cleanup.DeleteDisk(ctx, client, zone, disk.ID, client.checkReferencedOption()); err != nil {
		return diag.Errorf("deleting SakuraCloud Disk[%s] is failed: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func setDiskResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *iaas.Disk) diag.Diagnostics {
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
	d.Set("encryption_algorithm", data.EncryptionAlgorithm.String())
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
