package sakuracloud

import (
	"fmt"
	"github.com/docker/go-units"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"plan": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  sacloud.DiskPlanSSD.ID.String(),
				ValidateFunc: validateStringInWord([]string{
					sacloud.DiskPlanSSD.ID.String(),
					sacloud.DiskPlanHDD.ID.String(),
				}),
			},
			"connection": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  sacloud.DiskConnectionVirtio,
				ValidateFunc: validateStringInWord([]string{
					fmt.Sprintf("%s", sacloud.DiskConnectionVirtio),
					fmt.Sprintf("%s", sacloud.DiskConnectionIDE),
				}),
			},
			"source_archive_id": &schema.Schema{
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"source_disk_id"},
			},
			"source_disk_id": &schema.Schema{
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"source_archive_id"},
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  20,
			},
			"server_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true, //ReadOnly
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zone": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateStringInWord([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			//"hostname": &schema.Schema{
			//	Type:         schema.TypeString,
			//	Optional:     true,
			//	ValidateFunc: validateMaxLength(1, 64),
			//},
			"password": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateMaxLength(8, 64),
			},
			"ssh_key_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"disable_pw_auth": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"user_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_route": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_mask_len": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"note_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
	opts.Plan.SetIDByString(d.Get("plan").(string))
	opts.Connection = sacloud.EDiskConnection(d.Get("connection").(string))

	archiveID, ok := d.GetOk("source_archive_id")
	if ok {
		opts.SetSourceArchive(archiveID.(string))
	}
	diskID, ok := d.GetOk("source_disk_id")
	if ok {
		opts.SetSourceDisk(diskID.(string))
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

	err = client.Disk.SleepWhileCopying(disk.ID, 10*time.Minute)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Disk resource: %s", err)
	}

	//edit disk
	diskEditCondig := client.Disk.NewCondig()
	//if hostName, ok := d.GetOk("hostname"); ok {
	//	diskEditCondig.SetHostName(hostName.(string))
	//}
	if password, ok := d.GetOk("password"); ok {
		diskEditCondig.SetPassword(password.(string))
	}
	if sshKeyIDs, ok := d.GetOk("ssh_key_ids"); ok {
		ids := expandStringList(sshKeyIDs.([]interface{}))
		diskEditCondig.SetSSHKeys(ids)
	}

	if disablePasswordAuth, ok := d.GetOk("disable_pw_auth"); ok {
		diskEditCondig.SetDisablePWAuth(disablePasswordAuth.(bool))
	}

	if noteIDs, ok := d.GetOk("note_ids"); ok {
		ids := expandStringList(noteIDs.([]interface{}))
		diskEditCondig.SetNotes(ids)
	}

	if userIPAddress, ok := d.GetOk("user_ip_address"); ok {
		diskEditCondig.SetUserIPAddress(userIPAddress.(string))
	}

	if defaultRoute, ok := d.GetOk("default_route"); ok {
		diskEditCondig.SetDefaultRoute(defaultRoute.(string))
	}

	if networkMaskLen, ok := d.GetOk("network_mask_len"); ok {
		diskEditCondig.SetNetworkMaskLen(fmt.Sprint(networkMaskLen.(int)))
	}

	// call disk edit API
	_, err = client.Disk.Config(disk.ID, diskEditCondig)
	if err != nil {
		return fmt.Errorf("Error editting SakuraCloud DiskConfig: %s", err)
	}

	server_id, ok := d.GetOk("server_id")
	if ok {
		_, err = client.Disk.ConnectToServer(disk.ID, server_id.(string))

		if err != nil {
			return fmt.Errorf("Failed to connect SakuraCloud Disk resource: %s", err)
		}
	}

	d.SetId(disk.ID)
	return resourceSakuraCloudDiskRead(d, meta)
}

func resourceSakuraCloudDiskRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	disk, err := client.Disk.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Disk resource: %s", err)
	}

	d.Set("name", disk.Name)
	d.Set("plan", disk.Plan.ID.String())
	d.Set("connection", fmt.Sprintf("%s", disk.Connection))
	d.Set("size", disk.SizeMB*units.MiB/units.GiB)
	d.Set("description", disk.Description)
	d.Set("tags", disk.Tags)

	if disk.Server != nil {
		d.Set("server_id", disk.Server.ID)
	}

	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	disk, err := client.Disk.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Disk resource: %s", err)
	}

	// has server_id and server is up,shutdown
	isRunning := disk.Server != nil && disk.Server.Instance.IsUp()
	isDiskConfigChanged := false

	if d.HasChange("passowrd") || d.HasChange("ssh_key_ids") || d.HasChange("disable_pw_auth") || d.HasChange("note_ids") || d.HasChange("user_ip_address") || d.HasChange("default_route") || d.HasChange("network_mask_len") {
		isDiskConfigChanged = true
	}

	if isRunning && isDiskConfigChanged {
		_, err := client.Server.Shutdown(disk.Server.ID)
		if err != nil {
			return fmt.Errorf("Error stopping SakuraCloud Server resource: %s", err)
		}

		err = client.Server.SleepUntilDown(disk.Server.ID, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("Error stopping SakuraCloud Server resource: %s", err)
		}
	}

	if isDiskConfigChanged {
		diskEditCondig := client.Disk.NewCondig()
		//if d.HasChange("hostname") {
		//	if hostName, ok := d.GetOk("hostname"); ok {
		//		diskEditCondig.SetHostName(hostName.(string))
		//	} else {
		//		diskEditCondig.SetHostName("")
		//	}
		//}

		if d.HasChange("password") {
			if password, ok := d.GetOk("password"); ok {
				diskEditCondig.SetPassword(password.(string))
			} else {
				diskEditCondig.SetPassword("")
			}
		}

		if d.HasChange("ssh_key_ids") {
			if sshKeyIDs, ok := d.GetOk("ssh_key_ids"); ok {
				ids := expandStringList(sshKeyIDs.([]interface{}))
				diskEditCondig.SetSSHKeys(ids)
			} else {
				diskEditCondig.SSHKeys = nil
			}
		}

		if d.HasChange("disable_pw_auth") {
			diskEditCondig.SetDisablePWAuth(d.Get("disable_pw_auth").(bool))
		}

		if d.HasChange("user_ip_address") || d.HasChange("default_route") || d.HasChange("network_mask_len") {
			diskEditCondig.SetUserIPAddress(d.Get("user_ip_address").(string))
			diskEditCondig.SetDefaultRoute(d.Get("default_route").(string))
			diskEditCondig.SetNetworkMaskLen(fmt.Sprint(d.Get("network_mask_len").(int)))
		}

		if d.HasChange("note_ids") {
			if noteIDs, ok := d.GetOk("note_ids"); ok {
				ids := expandStringList(noteIDs.([]interface{}))
				diskEditCondig.SetNotes(ids)
			} else {
				diskEditCondig.Notes = nil
			}
		}

		_, err := client.Disk.Config(disk.ID, diskEditCondig)
		if err != nil {
			return fmt.Errorf("Error editting SakuraCloud DiskConfig: %s", err)
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

	d.SetId(disk.ID)

	if isRunning && isDiskConfigChanged {
		_, err := client.Server.Boot(disk.Server.ID)
		if err != nil {
			return fmt.Errorf("Error booting SakuraCloud Server resource: %s", err)
		}

		err = client.Server.SleepUntilUp(disk.Server.ID, 10*time.Minute)
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

	disk, err := client.Disk.Read(d.Id())
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
			err = client.Server.SleepUntilDown(disk.Server.ID, 10*time.Minute)
			if err != nil {
				return fmt.Errorf("Error stopping Server: %s", err)
			}
		}

		_, err := client.Disk.DisconnectFromServer(d.Id())
		if err != nil {
			return fmt.Errorf("Error disconnecting Disk from SakuraCloud Server: %s", err)
		}

	}

	_, err = client.Disk.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Disk resource: %s", err)
	}

	if isRunning {
		_, err := client.Server.Boot(disk.Server.ID)
		if err != nil {
			return fmt.Errorf("Error booting Server: %s", err)
		}
		err = client.Server.SleepUntilUp(disk.Server.ID, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("Error booting Server: %s", err)
		}

	}

	return nil
}
