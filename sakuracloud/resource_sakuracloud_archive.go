// Copyright 2016-2022 terraform-provider-sakuracloud authors
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudArchive() *schema.Resource {
	resourceName := "archive"

	return &schema.Resource{
		CreateContext: resourceSakuraCloudArchiveCreate,
		ReadContext:   resourceSakuraCloudArchiveRead,
		UpdateContext: resourceSakuraCloudArchiveUpdate,
		DeleteContext: resourceSakuraCloudArchiveDelete,
		CustomizeDiff: customdiff.ComputedIf("hash", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
			return d.HasChange("archive_file")
		}),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(24 * time.Hour),
			Update: schema.DefaultTimeout(24 * time.Hour),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"size": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice(types.ArchiveSizes)),
				Description:      desc.Sprintf("The size of %s in GiB. This must be one of [%s]", resourceName, types.ArchiveSizes),
				ConflictsWith:    []string{"source_disk_id", "source_archive_id", "source_shared_key", "source_archive_zone"},
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
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				ConflictsWith:    []string{"source_disk_id", "size", "source_shared_key"},
				ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
				Description: desc.Sprintf(
					"The id of the source archive. %s",
					desc.Conflicts("source_disk_id"),
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
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				ConflictsWith:    []string{"source_archive_id", "size", "source_shared_key", "source_archive_zone"},
				ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
				Description: desc.Sprintf(
					"The id of the source disk. %s",
					desc.Conflicts("source_archive_id"),
				),
			},
			"source_shared_key": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				Sensitive:        true,
				ConflictsWith:    []string{"source_archive_id", "source_disk_id", "size", "source_archive_zone"},
				ValidateDiagFunc: validation.ToDiagFunc(validateSourceSharedKey),
				Description:      "The share key of source shared archive",
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"zone":        schemaResourceZone(resourceName),
		},
	}
}

func resourceSakuraCloudArchiveCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	builder, cleanup, err := expandArchiveBuilder(d, zone, client)
	if err != nil {
		return diag.FromErr(err)
	}
	if cleanup != nil {
		defer cleanup()
	}

	archive, err := builder.Build(ctx, zone)
	if archive != nil {
		d.SetId(archive.ID.String())
	}
	if err != nil {
		return diag.Errorf("creating SakuraCloud Archive is failed: %s", err)
	}

	return resourceSakuraCloudArchiveRead(ctx, d, meta)
}

func resourceSakuraCloudArchiveRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	archiveOp := iaas.NewArchiveOp(client)

	archive, err := archiveOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud Archive[%s]: %s", d.Id(), err)
	}
	return setArchiveResourceData(d, client, archive)
}

func resourceSakuraCloudArchiveUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	archiveOp := iaas.NewArchiveOp(client)

	archive, err := archiveOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud Archive[%s]: %s", d.Id(), err)
	}

	if _, err = archiveOp.Update(ctx, zone, archive.ID, expandArchiveUpdateRequest(d)); err != nil {
		return diag.Errorf("updating SakuraCloud Archive[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudArchiveRead(ctx, d, meta)
}

func resourceSakuraCloudArchiveDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	archiveOp := iaas.NewArchiveOp(client)

	archive, err := archiveOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud Archive[%s]: %s", d.Id(), err)
	}

	if err := archiveOp.Delete(ctx, zone, archive.ID); err != nil {
		return diag.Errorf("deleting SakuraCloud Archive[%s] is failed: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func setArchiveResourceData(d *schema.ResourceData, client *APIClient, data *iaas.Archive) diag.Diagnostics {
	d.Set("hash", expandArchiveHash(d))                             // nolint
	d.Set("icon_id", data.IconID.String())                          // nolint
	d.Set("name", data.Name)                                        // nolint
	d.Set("size", data.GetSizeGB())                                 // nolint
	d.Set("description", data.Description)                          // nolint
	d.Set("zone", getZone(d, client))                               // nolint
	d.Set("source_archive_id", d.Get("source_archive_id").(string)) // nolint
	d.Set("source_disk_id", d.Get("source_disk_id").(string))       // nolint
	d.Set("source_shared_key", d.Get("source_shared_key").(string)) // nolint
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
