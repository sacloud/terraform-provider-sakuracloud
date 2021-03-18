// Copyright 2016-2021 terraform-provider-sakuracloud authors
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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudArchive() *schema.Resource {
	resourceName := "archive"

	return &schema.Resource{
		Create: resourceSakuraCloudArchiveCreate,
		Read:   resourceSakuraCloudArchiveRead,
		Update: resourceSakuraCloudArchiveUpdate,
		Delete: resourceSakuraCloudArchiveDelete,
		CustomizeDiff: customdiff.ComputedIf("hash", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
			return d.HasChange("archive_file")
		}),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(24 * time.Hour),
			Update: schema.DefaultTimeout(24 * time.Hour),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"size": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				Default:       20,
				ValidateFunc:  validation.IntInSlice(types.ArchiveSizes),
				Description:   descf("The size of %s in GiB. This must be one of [%s]", resourceName, types.ArchiveSizes),
				ConflictsWith: []string{"source_disk_id", "source_archive_id", "source_shared_key", "source_archive_zone"},
			},
			"archive_file": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The file path to upload to the SakuraCloud",
			},
			"hash": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The md5 checksum calculated from the base64 encoded file body",
			},
			"source_archive_id": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"source_disk_id", "size", "source_shared_key"},
				ValidateFunc:  validateSakuracloudIDType,
				Description: descf(
					"The id of the source archive. %s",
					descConflicts("source_disk_id"),
				),
			},
			"source_archive_zone": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"source_shared_key", "source_disk_id", "size"},
				Description:   "The share key of source shared archive",
			},
			"source_disk_id": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"source_archive_id", "size", "source_shared_key", "source_archive_zone"},
				ValidateFunc:  validateSakuracloudIDType,
				Description: descf(
					"The id of the source disk. %s",
					descConflicts("source_archive_id"),
				),
			},
			"source_shared_key": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"source_archive_id", "source_disk_id", "size", "source_archive_zone"},
				ValidateFunc:  validateSourceSharedKey,
				Description:   "The share key of source shared archive",
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"zone":        schemaResourceZone(resourceName),
		},
	}
}

func resourceSakuraCloudArchiveCreate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	builder, cleanup, err := expandArchiveBuilder(d, zone, client)
	if err != nil {
		return err
	}
	if cleanup != nil {
		defer cleanup()
	}

	archive, err := builder.Build(ctx, zone)
	if archive != nil {
		d.SetId(archive.ID.String())
	}
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Archive is failed: %s", err)
	}

	return resourceSakuraCloudArchiveRead(d, meta)
}

func resourceSakuraCloudArchiveRead(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	archiveOp := sacloud.NewArchiveOp(client)

	archive, err := archiveOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Archive[%s]: %s", d.Id(), err)
	}
	return setArchiveResourceData(d, client, archive)
}

func resourceSakuraCloudArchiveUpdate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	archiveOp := sacloud.NewArchiveOp(client)

	archive, err := archiveOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Archive[%s]: %s", d.Id(), err)
	}

	if _, err = archiveOp.Update(ctx, zone, archive.ID, expandArchiveUpdateRequest(d)); err != nil {
		return fmt.Errorf("updating SakuraCloud Archive[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudArchiveRead(d, meta)
}

func resourceSakuraCloudArchiveDelete(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	archiveOp := sacloud.NewArchiveOp(client)

	archive, err := archiveOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Archive[%s]: %s", d.Id(), err)
	}

	if err := archiveOp.Delete(ctx, zone, archive.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud Archive[%s] is failed: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func setArchiveResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.Archive) error {
	d.Set("hash", expandArchiveHash(d))                             // nolint
	d.Set("icon_id", data.IconID.String())                          // nolint
	d.Set("name", data.Name)                                        // nolint
	d.Set("size", data.GetSizeGB())                                 // nolint
	d.Set("description", data.Description)                          // nolint
	d.Set("zone", getZone(d, client))                               // nolint
	d.Set("source_archive_id", d.Get("source_archive_id").(string)) // nolint
	d.Set("source_disk_id", d.Get("source_disk_id").(string))       // nolint
	d.Set("source_shared_key", d.Get("source_shared_key").(string)) // nolint
	return d.Set("tags", flattenTags(data.Tags))
}
