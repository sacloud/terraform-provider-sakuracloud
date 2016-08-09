package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
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

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backup_hour": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validateIntInList(sacloud.AllowAutoBackupHour()),
			},
			"weekdays": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"max_backup_num": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validateIntegerInRange(1, 10),
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
				ValidateFunc: validateStringInWord([]string{"is1b", "tk1a"}),
			},
		},
	}
}

func resourceSakuraCloudAutoBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	opts := client.AutoBackup.New(d.Get("name").(string), d.Get("disk_id").(string))
	opts.SetBackupHour(d.Get("backup_hour").(int))
	opts.SetBackupMaximumNumberOfArchives(d.Get("max_backup_num").(int))
	rawWeekdays := d.Get("weekdays").([]interface{})
	if rawWeekdays != nil {
		weekdays, err := expandStringListWithValidateInList("weekdays", rawWeekdays, sacloud.AllowAutoBackupWeekdays())
		if err != nil {
			return err
		}
		opts.SetBackupSpanWeekdays(weekdays)
	}

	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}

	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		opts.Tags = expandStringList(rawTags)
	}

	autoBackup, err := client.AutoBackup.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud AutoBackup resource: %s", err)
	}

	d.SetId(autoBackup.ID)
	return resourceSakuraCloudAutoBackupRead(d, meta)
}

func resourceSakuraCloudAutoBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	autoBackup, err := client.AutoBackup.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud AutoBackup resource: %s", err)
	}

	d.Set("name", autoBackup.Name)
	d.Set("disk_id", autoBackup.Status.DiskID)
	d.Set("backup_hour", autoBackup.Settings.Autobackup.BackupHour)
	d.Set("max_backup_num", autoBackup.Settings.Autobackup.MaximumNumberOfArchives)
	d.Set("weekdays", autoBackup.Settings.Autobackup.BackupSpanWeekdays)

	d.Set("description", autoBackup.Description)
	d.Set("tags", autoBackup.Tags)
	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudAutoBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	autoBackup, err := client.AutoBackup.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud AutoBackup resource: %s", err)
	}

	autoBackup.SetBackupHour(d.Get("backup_hour").(int))
	autoBackup.SetBackupMaximumNumberOfArchives(d.Get("max_backup_num").(int))
	rawWeekdays := d.Get("weekdays").([]interface{})
	if rawWeekdays != nil {
		weekdays, err := expandStringListWithValidateInList("weekdays", rawWeekdays, sacloud.AllowAutoBackupWeekdays())
		if err != nil {
			return err
		}
		autoBackup.SetBackupSpanWeekdays(weekdays)
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
		autoBackup.Tags = expandStringList(rawTags)
	}

	autoBackup, err = client.AutoBackup.Update(autoBackup.ID, autoBackup)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud AutoBackup resource: %s", err)
	}

	d.SetId(autoBackup.ID)
	return resourceSakuraCloudAutoBackupRead(d, meta)

}

func resourceSakuraCloudAutoBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	_, err := client.AutoBackup.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud AutoBackup resource: %s", err)
	}

	return nil
}
