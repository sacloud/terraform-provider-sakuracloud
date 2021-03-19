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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudArchiveShare() *schema.Resource {
	resourceName := "ArchiveShare"

	return &schema.Resource{
		CreateContext: resourceSakuraCloudArchiveShareCreate,
		ReadContext:   resourceSakuraCloudArchiveShareRead,
		DeleteContext: resourceSakuraCloudArchiveShareDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"archive_id": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
				Description:      "The id of the archive",
			},
			"share_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The key to use sharing the Archive",
			},
			"zone": schemaResourceZone(resourceName),
		},
	}
}

func resourceSakuraCloudArchiveShareCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	archiveOp := sacloud.NewArchiveOp(client)
	archiveID := expandSakuraCloudID(d, "archive_id")

	archive, err := archiveOp.Read(ctx, zone, archiveID)
	if err != nil {
		return diag.Errorf("sharing SakuraCloud Archive is failed: %s", err)
	}

	// share
	shareInfo, err := archiveOp.Share(ctx, zone, archiveID)
	if err != nil {
		return diag.Errorf("sharing SakuraCloud Archive is failed: %s", err)
	}

	d.SetId(archive.ID.String())
	d.Set("share_key", shareInfo.SharedKey)
	d.Set("zone", zone)
	return nil
}

func resourceSakuraCloudArchiveShareRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	archiveOp := sacloud.NewArchiveOp(client)

	archive, err := archiveOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud Archive[%s]: %s", d.Id(), err)
	}

	if !archive.Availability.IsUploading() {
		d.SetId("")
		return nil
	}

	d.SetId(archive.ID.String())
	d.Set("share_key", d.Get("share_key").(string))
	d.Set("zone", zone)
	return nil
}

func resourceSakuraCloudArchiveShareDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	archiveOp := sacloud.NewArchiveOp(client)

	archive, err := archiveOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud Archive[%s]: %s", d.Id(), err)
	}

	if !archive.Availability.IsUploading() {
		d.SetId("")
		return nil
	}

	if err := archiveOp.CloseFTP(ctx, zone, archive.ID); err != nil {
		return diag.Errorf("deleting SakuraCloud Archive Share[%s] is failed: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
