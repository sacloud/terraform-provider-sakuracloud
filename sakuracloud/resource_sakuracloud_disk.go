package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/utils/setup"

	"log"

	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
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
					fmt.Sprintf("%s", sacloud.DiskConnectionVirtio),
					fmt.Sprintf("%s", sacloud.DiskConnectionIDE),
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
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(8, 64),
				Sensitive:    true,
			},
			"ssh_key_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// ! Current terraform(v0.7) is not support to array validation !
				// ValidateFunc: validateSakuracloudIDArrayType,
			},
			"disable_pw_auth": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"note_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// ! Current terraform(v0.7) is not support to array validation !
				// ValidateFunc: validateSakuracloudIDArrayType,
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
		break
	case "hdd":
		opts.SetDiskPlanToHDD()
		break
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
			//edit disk
			diskEditConfig := client.Disk.NewCondig()
			if hostName, ok := d.GetOk("hostname"); ok {
				diskEditConfig.SetHostName(hostName.(string))
			}
			if password, ok := d.GetOk("password"); ok {
				diskEditConfig.SetPassword(password.(string))
			}
			if sshKeyIDs, ok := d.GetOk("ssh_key_ids"); ok {
				ids := expandStringList(sshKeyIDs.([]interface{}))
				diskEditConfig.SetSSHKeys(ids)
			}

			if disablePasswordAuth, ok := d.GetOk("disable_pw_auth"); ok {
				diskEditConfig.SetDisablePWAuth(disablePasswordAuth.(bool))
			}

			if noteIDs, ok := d.GetOk("note_ids"); ok {
				ids := expandStringList(noteIDs.([]interface{}))
				diskEditConfig.SetNotes(ids)
			}

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
			} else {
				log.Printf("[WARN] Disk[%d] does not support modify disk", id)
			}

			server_id, ok := d.GetOk("server_id")
			if ok {
				_, err = client.Disk.ConnectToServer(id, toSakuraCloudID(server_id.(string)))

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

	if d.HasChange("hostname") || d.HasChange("password") || d.HasChange("ssh_key_ids") || d.HasChange("disable_pw_auth") || d.HasChange("note_ids") {
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
		if d.HasChange("hostname") {
			if hostName, ok := d.GetOk("hostname"); ok {
				diskEditConfig.SetHostName(hostName.(string))
			} else {
				diskEditConfig.HostName = nil
			}
		}

		if d.HasChange("password") {
			if password, ok := d.GetOk("password"); ok {
				diskEditConfig.SetPassword(password.(string))
			} else {
				diskEditConfig.SetPassword("")
			}
		}

		if d.HasChange("ssh_key_ids") {
			if sshKeyIDs, ok := d.GetOk("ssh_key_ids"); ok {
				ids := expandStringList(sshKeyIDs.([]interface{}))
				diskEditConfig.SetSSHKeys(ids)
			} else {
				diskEditConfig.SSHKeys = nil
			}
		}

		if d.HasChange("disable_pw_auth") {
			diskEditConfig.SetDisablePWAuth(d.Get("disable_pw_auth").(bool))
		}

		if d.HasChange("note_ids") {
			if noteIDs, ok := d.GetOk("note_ids"); ok {
				ids := expandStringList(noteIDs.([]interface{}))
				diskEditConfig.SetNotes(ids)
			} else {
				diskEditConfig.Notes = nil
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
		break
	case sacloud.DiskPlanHDD.ID:
		d.Set("plan", "hdd")
		break

	}

	if data.SourceDisk != nil {
		d.Set("source_disk_id", data.SourceDisk.GetStrID())
	} else if data.SourceArchive != nil {
		d.Set("source_archive_id", data.SourceArchive.GetStrID())
	}

	d.Set("connector", fmt.Sprintf("%s", data.Connection))
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
