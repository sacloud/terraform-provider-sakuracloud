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
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/utils/setup"
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
				Default:  sacloud.DiskConnectionVirtio,
				ValidateFunc: validation.StringInSlice([]string{
					string(sacloud.DiskConnectionVirtio),
					string(sacloud.DiskConnectionIDE),
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
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			powerManageTimeoutKey: powerManageTimeoutParam,
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"hostname": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Removed:      "Use attribute in `sakuracloud_server` instead",
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(8, 64),
				Sensitive:    true,
				Removed:      "Use attribute in `sakuracloud_server` instead",
			},
			"ssh_key_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// ! Current terraform(v0.7) is not support to array validation !
				// ValidateFunc: validateSakuracloudIDArrayType,
				Removed: "Use attribute in `sakuracloud_server` instead",
			},
			"disable_pw_auth": {
				Type:     schema.TypeBool,
				Optional: true,
				Removed:  "Use attribute in `sakuracloud_server` instead",
			},
			"note_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// ! Current terraform(v0.7) is not support to array validation !
				// ValidateFunc: validateSakuracloudIDArrayType,
				Removed: "Use attribute in `sakuracloud_server` instead",
			},
		},
	}
}

func resourceSakuraCloudDiskCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	opts := client.Disk.New()

	opts.Name = d.Get("name").(string)

	plan := d.Get("plan").(string)
	switch plan {
	case "ssd":
		opts.SetDiskPlanToSSD()
	case "hdd":
		opts.SetDiskPlanToHDD()
	default:
		return fmt.Errorf("invalid disk plan [%s]", plan)
	}

	opts.Connection = sacloud.EDiskConnection(d.Get("connector").(string))

	archiveID, ok := d.GetOk("source_archive_id")
	if ok {
		opts.SetSourceArchive(toSakuraCloudID(archiveID.(string)))
	}
	diskID, ok := d.GetOk("source_disk_id")
	if ok {
		opts.SetSourceDisk(toSakuraCloudID(diskID.(string)))
	}

	opts.SizeMB = toSizeMB(d.Get("size").(int))
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

	if _, ok := d.GetOk("distant_from"); ok {
		rawDistantFrom := d.Get("distant_from").([]interface{})
		if rawDistantFrom != nil {
			strDiskIDs := expandStringList(rawDistantFrom)
			diskIDs := []int64{}
			for _, id := range strDiskIDs {
				diskIDs = append(diskIDs, int64(forceAtoI(id)))
			}

			opts.DistantFrom = diskIDs
		}
	}

	diskBuilder := &setup.RetryableSetup{
		Create: func() (sacloud.ResourceIDHolder, error) {
			return client.Disk.Create(opts)
		},
		AsyncWaitForCopy: func(id int64) (chan interface{}, chan interface{}, chan error) {
			return client.Disk.AsyncSleepWhileCopying(id, client.DefaultTimeoutDuration)
		},
		Delete: func(id int64) error {
			_, err := client.Disk.Delete(id)
			return err
		},
		ProvisionBeforeUp: func(id int64, _ interface{}) error {
			isNeedEditDisk := false

			//edit disk
			diskEditConfig := client.Disk.NewCondig()
			diskEditConfig.SetBackground(true)

			if hostName, ok := d.GetOk("hostname"); ok {
				diskEditConfig.SetHostName(hostName.(string))
				isNeedEditDisk = true
			}
			if password, ok := d.GetOk("password"); ok {
				diskEditConfig.SetPassword(password.(string))
				isNeedEditDisk = true
			}
			if sshKeyIDs, ok := d.GetOk("ssh_key_ids"); ok {
				ids := expandStringList(sshKeyIDs.([]interface{}))
				diskEditConfig.SetSSHKeys(ids)
				isNeedEditDisk = true
			}

			if disablePasswordAuth, ok := d.GetOk("disable_pw_auth"); ok {
				diskEditConfig.SetDisablePWAuth(disablePasswordAuth.(bool))
				isNeedEditDisk = true
			}

			if noteIDs, ok := d.GetOk("note_ids"); ok {
				ids := expandStringList(noteIDs.([]interface{}))
				diskEditConfig.SetNotes(ids)
				isNeedEditDisk = true
			}

			if isNeedEditDisk {
				// call disk edit API
				res, err := client.Disk.CanEditDisk(id)
				if err != nil {
					return fmt.Errorf("Failed to check CanEditDisk: %s", err)
				}
				if res {
					_, err = client.Disk.Config(id, diskEditConfig)
					if err != nil {
						return fmt.Errorf("Error editting SakuraCloud DiskConfig: %s", err)
					}
					// wait
					if err := client.Disk.SleepWhileCopying(id, client.DefaultTimeoutDuration); err != nil {
						return fmt.Errorf("Error editting SakuraCloud DiskConfig: timeout: %s", err)
					}
				} else {
					log.Printf("[WARN] Disk[%d] does not support modify disk", id)
				}
			}

			server_id, ok := d.GetOk("server_id")
			if ok {
				_, err := client.Disk.ConnectToServer(id, toSakuraCloudID(server_id.(string)))
				if err != nil {
					return fmt.Errorf("Failed to connect SakuraCloud Disk resource: %s", err)
				}
			}

			return nil
		},
		RetryCount: 3,
	}

	res, err := diskBuilder.Setup()
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Disk resource: %s", err)
	}

	disk, ok := res.(*sacloud.Disk)
	if !ok {
		return fmt.Errorf("Failed to create SakuraCloud Disk resource: created resource is not *sacloud.Disk")
	}

	d.SetId(disk.GetStrID())
	return resourceSakuraCloudDiskRead(d, meta)
}

func resourceSakuraCloudDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	disk, err := client.Disk.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud Disk resource: %s", err)
	}

	return setDiskResourceData(d, client, disk)
}

func resourceSakuraCloudDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	disk, err := client.Disk.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Disk resource: %s", err)
	}

	// has server_id and server is up,shutdown
	isRunning := disk.Server != nil && disk.Server.Instance.IsUp()
	isDiskConfigChanged := false

	if d.HasChange("hostname") || d.HasChange("password") || d.HasChange("ssh_key_ids") ||
		d.HasChange("disable_pw_auth") || d.HasChange("note_ids") {
		isDiskConfigChanged = true
	}

	if isRunning && isDiskConfigChanged {
		err := stopServer(client, disk.Server.ID, d)
		if err != nil {
			return fmt.Errorf("Error stopping SakuraCloud Server resource: %s", err)
		}
	}

	if isDiskConfigChanged {
		diskEditConfig := client.Disk.NewCondig()
		diskEditConfig.SetBackground(true)
		if d.HasChange("hostname") {
			if hostName, ok := d.GetOk("hostname"); ok {
				diskEditConfig.SetHostName(hostName.(string))
			}
		}

		if d.HasChange("password") {
			if password, ok := d.GetOk("password"); ok {
				diskEditConfig.SetPassword(password.(string))
			}
		}

		if d.HasChange("ssh_key_ids") {
			if sshKeyIDs, ok := d.GetOk("ssh_key_ids"); ok {
				ids := expandStringList(sshKeyIDs.([]interface{}))
				diskEditConfig.SetSSHKeys(ids)
			}
		}

		if d.HasChange("disable_pw_auth") {
			diskEditConfig.SetDisablePWAuth(d.Get("disable_pw_auth").(bool))
		}

		if d.HasChange("note_ids") {
			if noteIDs, ok := d.GetOk("note_ids"); ok {
				ids := expandStringList(noteIDs.([]interface{}))
				diskEditConfig.SetNotes(ids)
			}
		}

		res, err := client.Disk.CanEditDisk(disk.ID)
		if err != nil {
			return fmt.Errorf("Failed to check CanEditDisk: %s", err)
		}
		if res {
			_, err := client.Disk.Config(disk.ID, diskEditConfig)
			if err != nil {
				return fmt.Errorf("Error editting SakuraCloud DiskConfig: %s", err)
			}
			if err := client.Disk.SleepWhileCopying(disk.ID, client.DefaultTimeoutDuration); err != nil {
				return fmt.Errorf("Error editting SakuraCloud DiskConfig: timeout: %s", err)
			}
		} else {
			log.Printf("[WARN] Disk[%d] does not support modify disk", disk.ID)
		}

	}

	if d.HasChange("name") {
		disk.Name = d.Get("name").(string)
	}
	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			disk.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			disk.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			disk.Description = description.(string)
		} else {
			disk.Description = ""
		}
	}
	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			disk.Tags = expandTags(client, rawTags)
		} else {
			disk.Tags = expandTags(client, []interface{}{})
		}
	}

	disk, err = client.Disk.Update(disk.ID, disk)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Disk resource: %s", err)
	}

	if isRunning && isDiskConfigChanged {
		err := bootServer(client, disk.Server.ID)
		if err != nil {
			return fmt.Errorf("Error booting SakuraCloud Server resource: %s", err)
		}
	}

	return resourceSakuraCloudDiskRead(d, meta)
}

func resourceSakuraCloudDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	disk, err := client.Disk.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Disk resource: %s", err)
	}

	isRunning := false
	if disk.Server != nil {

		// lock during server delete operation
		lockKey := getServerDeleteAPILockKey(disk.Server.ID)
		sakuraMutexKV.Lock(lockKey)
		defer sakuraMutexKV.Unlock(lockKey)

		server, err := client.Server.Read(disk.Server.ID)
		if err == nil {
			if server.IsUp() {
				isRunning = true
				err := stopServer(client, server.ID, d)
				if err != nil {
					return fmt.Errorf("Error stopping SakuraCloud Server resource: %s", err)
				}
			}

			_, err := client.Disk.DisconnectFromServer(toSakuraCloudID(d.Id()))
			if err != nil {
				return fmt.Errorf("Error disconnecting Disk from SakuraCloud Server: %s", err)
			}
		}

	}

	_, err = client.Disk.Delete(toSakuraCloudID(d.Id()))

	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Disk resource: %s", err)
	}

	if isRunning {
		err := bootServer(client, disk.Server.ID)
		if err != nil {
			return fmt.Errorf("Error booting Server: %s", err)
		}
	}

	return nil
}

func setDiskResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.Disk) error {
	d.Set("name", data.Name)

	switch data.Plan.ID {
	case sacloud.DiskPlanSSD.ID:
		d.Set("plan", "ssd")
	case sacloud.DiskPlanHDD.ID:
		d.Set("plan", "hdd")
	}

	if data.SourceDisk != nil {
		d.Set("source_disk_id", data.SourceDisk.GetStrID())
	} else if data.SourceArchive != nil {
		d.Set("source_archive_id", data.SourceArchive.GetStrID())
	}

	d.Set("connector", string(data.Connection))
	d.Set("size", toSizeGB(data.SizeMB))
	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	d.Set("tags", data.Tags)

	if data.Server == nil {
		d.Set("server_id", "")
	} else {
		d.Set("server_id", data.Server.GetStrID())
	}

	setPowerManageTimeoutValueToState(d)

	d.Set("zone", client.Zone)
	return nil
}
