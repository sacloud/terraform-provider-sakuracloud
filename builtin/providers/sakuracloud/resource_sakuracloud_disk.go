package sakuracloud

import (
	"fmt"
	"github.com/docker/go-units"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"log"
	"time"
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
				ValidateFunc: validateStringInWord([]string{"ssd", "hdd"}),
			},
			"connector": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  sacloud.DiskConnectionVirtio,
				ValidateFunc: validateStringInWord([]string{
					fmt.Sprintf("%s", sacloud.DiskConnectionVirtio),
					fmt.Sprintf("%s", sacloud.DiskConnectionIDE),
				}),
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
				ValidateFunc: validateStringInWord([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"hostname": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateMaxLength(1, 64),
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateMaxLength(8, 64),
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
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

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

	opts.SizeMB = d.Get("size").(int) * units.GiB / units.MiB
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		opts.Tags = expandStringList(rawTags)
	}

	disk, err := client.Disk.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Disk resource: %s", err)
	}

	err = client.Disk.SleepWhileCopying(disk.ID, client.DefaultTimeoutDuration)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Disk resource: %s", err)
	}

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
	res, err := client.Disk.CanEditDisk(disk.ID)
	if err != nil {
		return fmt.Errorf("Failed to check CanEditDisk: %s", err)
	}
	if res {
		_, err = client.Disk.Config(disk.ID, diskEditConfig)
		if err != nil {
			return fmt.Errorf("Error editting SakuraCloud DiskConfig: %s", err)
		}
	} else {
		log.Printf("[WARN] Disk[%d] does not support modify disk", disk.ID)
	}

	server_id, ok := d.GetOk("server_id")
	if ok {
		_, err = client.Disk.ConnectToServer(disk.ID, toSakuraCloudID(server_id.(string)))

		if err != nil {
			return fmt.Errorf("Failed to connect SakuraCloud Disk resource: %s", err)
		}
	}

	d.SetId(disk.GetStrID())
	return resourceSakuraCloudDiskRead(d, meta)
}

func resourceSakuraCloudDiskRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	disk, err := client.Disk.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Disk resource: %s", err)
	}

	return setDiskResourceData(d, client, disk)
}

func resourceSakuraCloudDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

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
		_, err := client.Server.Shutdown(disk.Server.ID)
		if err != nil {
			return fmt.Errorf("Error stopping SakuraCloud Server resource: %s", err)
		}

		err = client.Server.SleepUntilDown(disk.Server.ID, client.DefaultTimeoutDuration)
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
			disk.Tags = expandStringList(rawTags)
		}
	}

	disk, err = client.Disk.Update(disk.ID, disk)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Disk resource: %s", err)
	}

	d.SetId(disk.GetStrID())

	if isRunning && isDiskConfigChanged {
		_, err := client.Server.Boot(disk.Server.ID)
		if err != nil {
			return fmt.Errorf("Error booting SakuraCloud Server resource: %s", err)
		}

		err = client.Server.SleepUntilUp(disk.Server.ID, client.DefaultTimeoutDuration)
		if err != nil {
			return fmt.Errorf("Error booting SakuraCloud Server resource: %s", err)
		}
	}

	return resourceSakuraCloudDiskRead(d, meta)
}

func resourceSakuraCloudDiskDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	disk, err := client.Disk.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Disk resource: %s", err)
	}

	isRunning := false
	if disk.Server != nil {

		if disk.Server.Instance.IsUp() {
			isRunning = true
			time.Sleep(2 * time.Second)
			_, err := client.Server.Stop(disk.Server.ID)
			if err != nil {
				return fmt.Errorf("Error stopping Server: %s", err)
			}
			err = client.Server.SleepUntilDown(disk.Server.ID, client.DefaultTimeoutDuration)
			if err != nil {
				return fmt.Errorf("Error stopping Server: %s", err)
			}
		}

		_, err := client.Disk.DisconnectFromServer(toSakuraCloudID(d.Id()))
		if err != nil {
			return fmt.Errorf("Error disconnecting Disk from SakuraCloud Server: %s", err)
		}

	}

	_, err = client.Disk.Delete(toSakuraCloudID(d.Id()))

	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Disk resource: %s", err)
	}

	if isRunning {
		_, err := client.Server.Boot(disk.Server.ID)
		if err != nil {
			return fmt.Errorf("Error booting Server: %s", err)
		}
		err = client.Server.SleepUntilUp(disk.Server.ID, client.DefaultTimeoutDuration)
		if err != nil {
			return fmt.Errorf("Error booting Server: %s", err)
		}

	}

	return nil
}

func setDiskResourceData(d *schema.ResourceData, client *api.Client, data *sacloud.Disk) error {

	d.Set("name", data.Name)

	switch data.Plan.ID {
	case sacloud.DiskPlanSSD.ID:
		d.Set("plan", "ssd")
		break
	case sacloud.DiskPlanHDD.ID:
		d.Set("plan", "hdd")
		break

	}

	d.Set("connector", fmt.Sprintf("%s", data.Connection))
	d.Set("size", data.SizeMB*units.MiB/units.GiB)
	d.Set("description", data.Description)
	d.Set("tags", data.Tags)

	if data.Server != nil {
		d.Set("server_id", data.Server.GetStrID())
	}

	d.Set("zone", client.Zone)
	d.SetId(data.GetStrID())
	return nil
}
