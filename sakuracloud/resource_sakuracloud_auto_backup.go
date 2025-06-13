// Copyright 2016-2025 terraform-provider-sakuracloud authors
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
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudAutoBackup() *schema.Resource {
	resourceName := "AutoBackup"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudAutoBackupCreate,
		ReadContext:   resourceSakuraCloudAutoBackupRead,
		UpdateContext: resourceSakuraCloudAutoBackupUpdate,
		DeleteContext: resourceSakuraCloudAutoBackupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"disk_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
				Description:      "The disk id to backed up",
			},
			"weekdays": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Description: desc.Sprintf(
					"A list of weekdays to backed up. The values in the list must be in [%s]",
					types.DaysOfTheWeekStrings,
				),
			},
			"max_backup_num": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 10)),
				Description:      desc.Sprintf("The number backup files to keep. %s", desc.Range(1, 10)),
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"zone":        schemaResourceZone(resourceName),
		},
	}
}

func resourceSakuraCloudAutoBackupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	autoBackupOp := iaas.NewAutoBackupOp(client)

	if err := validateBackupWeekdays(d, "weekdays"); err != nil {
		return diag.FromErr(err)
	}

	autoBackup, err := autoBackupOp.Create(ctx, zone, expandAutoBackupCreateRequest(d))
	if err != nil {
		return diag.Errorf("creating SakuraCloud AutoBackup is failed: %s", err)
	}

	d.SetId(autoBackup.ID.String())
	return resourceSakuraCloudAutoBackupRead(ctx, d, meta)
}

func resourceSakuraCloudAutoBackupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	autoBackupOp := iaas.NewAutoBackupOp(client)
	autoBackup, err := autoBackupOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not find SakuraCloud AutoBackup[%s]: %s", d.Id(), err)
	}
	return setAutoBackupResourceData(d, client, autoBackup)
}

func resourceSakuraCloudAutoBackupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	autoBackupOp := iaas.NewAutoBackupOp(client)

	autoBackup, err := autoBackupOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud AutoBackup[%s]: %s", d.Id(), err)
	}

	if err := validateBackupWeekdays(d, "weekdays"); err != nil {
		return diag.FromErr(err)
	}

	if _, err = autoBackupOp.Update(ctx, zone, autoBackup.ID, expandAutoBackupUpdateRequest(d, autoBackup)); err != nil {
		return diag.Errorf("updating SakuraCloud AutoBackup[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudAutoBackupRead(ctx, d, meta)
}

func resourceSakuraCloudAutoBackupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	autoBackupOp := iaas.NewAutoBackupOp(client)
	autoBackup, err := autoBackupOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud AutoBackup[%s]: %s", d.Id(), err)
	}

	if err := autoBackupOp.Delete(ctx, zone, autoBackup.ID); err != nil {
		return diag.Errorf("deleting SakuraCloud AutoBackup[%s] is failed: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func setAutoBackupResourceData(d *schema.ResourceData, client *APIClient, data *iaas.AutoBackup) diag.Diagnostics {
	d.Set("name", data.Name)                              //nolint
	d.Set("disk_id", data.DiskID.String())                //nolint
	d.Set("max_backup_num", data.MaximumNumberOfArchives) //nolint
	d.Set("icon_id", data.IconID.String())                //nolint
	d.Set("description", data.Description)                //nolint
	d.Set("zone", getZone(d, client))                     //nolint
	if err := d.Set("weekdays", flattenBackupWeekdays(data.BackupSpanWeekdays)); err != nil {
		return diag.FromErr(err)
	}
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
