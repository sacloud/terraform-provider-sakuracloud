// Copyright 2016-2020 terraform-provider-sakuracloud authors
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
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
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
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1b", "tk1a", "is1a"}),
			},
		},
	}
}

func resourceSakuraCloudAutoBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	diskID := d.Get("disk_id").(string)
	opts := client.AutoBackup.New(d.Get("name").(string), toSakuraCloudID(diskID))
	opts.SetBackupMaximumNumberOfArchives(d.Get("max_backup_num").(int))
	rawWeekdays := d.Get("weekdays").([]interface{})
	if rawWeekdays != nil {
		weekdays, err := expandStringListWithValidateInList("weekdays", rawWeekdays, sacloud.AllowAutoBackupWeekdays())
		if err != nil {
			return err
		}
		opts.SetBackupSpanWeekdays(weekdays)
	}

	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.SetIconByID(toSakuraCloudID(iconID.(string)))
	}
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}

	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		opts.Tags = expandTags(client, rawTags)
	}

	autoBackup, err := client.AutoBackup.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud AutoBackup resource: %s", err)
	}

	d.SetId(autoBackup.GetStrID())
	return resourceSakuraCloudAutoBackupRead(d, meta)
}

func resourceSakuraCloudAutoBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	autoBackup, err := client.AutoBackup.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud AutoBackup resource: %s", err)
	}

	d.Set("name", autoBackup.Name)
	d.Set("disk_id", autoBackup.Status.DiskID)
	d.Set("max_backup_num", autoBackup.Settings.Autobackup.MaximumNumberOfArchives)
	d.Set("weekdays", autoBackup.Settings.Autobackup.BackupSpanWeekdays)
	d.Set("icon_id", autoBackup.GetIconStrID())
	d.Set("description", autoBackup.Description)
	d.Set("tags", autoBackup.Tags)
	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudAutoBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	autoBackup, err := client.AutoBackup.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud AutoBackup resource: %s", err)
	}

	autoBackup.SetBackupMaximumNumberOfArchives(d.Get("max_backup_num").(int))
	rawWeekdays := d.Get("weekdays").([]interface{})
	if rawWeekdays != nil {
		weekdays, err := expandStringListWithValidateInList("weekdays", rawWeekdays, sacloud.AllowAutoBackupWeekdays())
		if err != nil {
			return err
		}
		autoBackup.SetBackupSpanWeekdays(weekdays)
	}

	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			autoBackup.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			autoBackup.ClearIcon()
		}
	}

	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			autoBackup.Description = description.(string)
		} else {
			autoBackup.Description = ""
		}
	}
	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		autoBackup.Tags = expandTags(client, rawTags)
	} else {
		autoBackup.Tags = expandTags(client, []interface{}{})
	}

	_, err = client.AutoBackup.Update(autoBackup.ID, autoBackup)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud AutoBackup resource: %s", err)
	}

	return resourceSakuraCloudAutoBackupRead(d, meta)
}

func resourceSakuraCloudAutoBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	_, err := client.AutoBackup.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud AutoBackup resource: %s", err)
	}

	return nil
}
