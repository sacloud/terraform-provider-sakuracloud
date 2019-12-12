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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudAutoBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudAutoBackupCreate,
		Read:   resourceSakuraCloudAutoBackupRead,
		Update: resourceSakuraCloudAutoBackupUpdate,
		Delete: resourceSakuraCloudAutoBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"weekdays": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"max_backup_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 10),
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

func resourceSakuraCloudAutoBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	autoBackupOp := sacloud.NewAutoBackupOp(client)

	if err := validateBackupWeekdays(d, "weekdays"); err != nil {
		return err
	}

	autoBackup, err := autoBackupOp.Create(ctx, zone, expandAutoBackupCreateRequest(d))
	if err != nil {
		return fmt.Errorf("creating SakuraCloud AutoBackup is failed: %s", err)
	}

	d.SetId(autoBackup.ID.String())
	return resourceSakuraCloudAutoBackupRead(d, meta)
}

func resourceSakuraCloudAutoBackupRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	autoBackupOp := sacloud.NewAutoBackupOp(client)

	autoBackup, err := autoBackupOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not find SakuraCloud AutoBackup[%s]: %s", d.Id(), err)
	}
	return setAutoBackupResourceData(d, client, autoBackup)
}

func resourceSakuraCloudAutoBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	autoBackupOp := sacloud.NewAutoBackupOp(client)

	autoBackup, err := autoBackupOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud AutoBackup[%s]: %s", d.Id(), err)
	}

	if err := validateBackupWeekdays(d, "weekdays"); err != nil {
		return err
	}

	autoBackup, err = autoBackupOp.Update(ctx, zone, autoBackup.ID, expandAutoBackupUpdateRequest(d, autoBackup))
	if err != nil {
		return fmt.Errorf("updating SakuraCloud AutoBackup[%s] is failed: %s", autoBackup.ID, err)
	}

	return resourceSakuraCloudAutoBackupRead(d, meta)
}

func resourceSakuraCloudAutoBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	autoBackupOp := sacloud.NewAutoBackupOp(client)

	autoBackup, err := autoBackupOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud AutoBackup[%s]: %s", d.Id(), err)
	}

	if err := autoBackupOp.Delete(ctx, zone, autoBackup.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud AutoBackup[%s] is failed: %s", autoBackup.ID, err)
	}

	d.SetId("")
	return nil
}

func setAutoBackupResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.AutoBackup) error {
	d.Set("name", data.Name)
	d.Set("disk_id", data.DiskID.String())
	if err := d.Set("weekdays", flattenBackupWeekdays(data.BackupSpanWeekdays)); err != nil {
		return err
	}
	d.Set("max_backup_num", data.MaximumNumberOfArchives)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	d.Set("zone", getZone(d, client))
	return nil
}
