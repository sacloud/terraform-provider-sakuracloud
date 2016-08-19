package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"time"
)

func resourceSakuraCloudDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudDatabaseCreate,
		Read:   resourceSakuraCloudDatabaseRead,
		Update: resourceSakuraCloudDatabaseUpdate,
		Delete: resourceSakuraCloudDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			//"is_double": &schema.Schema{
			//	Type:     schema.TypeBool,
			//	ForceNew: true,
			//	Optional: true,
			//	Default:  false,
			//},
			//"plan": &schema.Schema{
			//	Type:         schema.TypeString,
			//	ForceNew:     true,
			//	Optional:     true,
			//	Default:      "standard",
			//	ValidateFunc: validateStringInWord([]string{"mini"}),
			//},
			"admin_password": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"user_name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"user_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"allow_networks": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"port": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5432,
				ValidateFunc: validateIntegerInRange(1024, 65535),
			},

			"backup_rotate": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      8,
				ValidateFunc: validateIntegerInRange(1, 8),
			},
			"backup_time": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateBackupTime(),
			},

			"switch_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "shared",
			},
			"ipaddress1": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
				Optional: true,
			},
			"nw_mask_len": &schema.Schema{
				Type:         schema.TypeInt,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(8, 29),
			},
			"default_route": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
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
				ValidateFunc: validateStringInWord([]string{"tk1a"}),
			},
		},
	}
}

func resourceSakuraCloudDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	opts := sacloud.NewCreatePostgreSQLDatabaseValue()

	opts.Name = d.Get("name").(string)
	opts.AdminPassword = d.Get("admin_password").(string)
	opts.DefaultUser = d.Get("user_name").(string)
	opts.UserPassword = d.Get("user_password").(string)
	if rawNetworks, ok := d.GetOk("allow_networks"); ok {
		if rawNetworks != nil {
			opts.SourceNetwork = expandStringList(rawNetworks.([]interface{}))
		}
	}
	opts.ServicePort = fmt.Sprintf("%d", d.Get("port").(int))
	opts.BackupRotate = d.Get("backup_rotate").(int)
	opts.BackupTime = d.Get("backup_time").(string)

	opts.SwitchID = d.Get("switch_id").(string)

	if opts.SwitchID != "shared" {
		_, ok1 := d.GetOk("ipaddress1")
		_, ok2 := d.GetOk("nw_mask_len")
		if !ok1 || !ok2 {
			msg := "ipaddress1 and nw_mask_len is required when SwitchID is exists"
			return fmt.Errorf("Faild to create SakuraCloud Database resource %s", msg)
		}

		ipAddress1 := d.Get("ipaddress1").(string)
		nwMaskLen := d.Get("nw_mask_len").(int)
		defaultRoute := ""
		if df, ok := d.GetOk("default_route"); ok {
			defaultRoute = df.(string)
		}

		opts.IPAddress1 = ipAddress1
		opts.MaskLen = nwMaskLen
		opts.DefaultRoute = defaultRoute
	}

	opts.Plan = sacloud.DatabasePlanMini

	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	if rawTags, ok := d.GetOk("tags"); ok {
		if rawTags != nil {
			opts.Tags = expandStringList(rawTags.([]interface{}))
		}
	}

	createDB := sacloud.CreateNewPostgreSQLDatabase(opts)
	database, err := client.Database.Create(createDB)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Database resource: %s", err)
	}

	d.SetId(database.ID)

	//wait
	err = client.Database.SleepWhileCopying(database.ID, 20*time.Minute, 5)
	if err != nil {
		return fmt.Errorf("Failed to wait SakuraCloud Database copy: %s", err)
	}
	err = client.Database.SleepUntilUp(database.ID, 20*time.Minute)
	if err != nil {
		return fmt.Errorf("Failed to wait SakuraCloud Database boot: %s", err)
	}

	return resourceSakuraCloudDatabaseRead(d, meta)
}

func resourceSakuraCloudDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	database, err := client.Database.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Database resource: %s", err)
	}

	d.Set("name", database.Name)
	d.Set("admin_password", database.Settings.DBConf.Common.AdminPassword)
	d.Set("user_password", database.Settings.DBConf.Common.DefaultUser)
	d.Set("user_password", database.Settings.DBConf.Common.UserPassword)

	d.Set("allow_networks", database.Settings.DBConf.Common.SourceNetwork)
	d.Set("port", database.Settings.DBConf.Common.ServicePort)

	d.Set("backup_rotate", database.Settings.DBConf.Backup.Rotate)
	d.Set("backup_time", database.Settings.DBConf.Backup.Time)

	if database.Interfaces[0].Switch.Scope == sacloud.ESCopeShared {
		d.Set("switch_id", "shared")
		d.Set("nw_mask_len", nil)
		d.Set("default_route", nil)
		d.Set("ipaddress1", database.Interfaces[0].IPAddress)
	} else {
		d.Set("switch_id", database.Interfaces[0].Switch.ID)
		d.Set("nw_mask_len", database.Remark.Network.NetworkMaskLen)
		d.Set("default_route", database.Remark.Network.DefaultRoute)
		d.Set("ipaddress1", database.Remark.Servers[0].(map[string]interface{})["IPAddress"])
	}

	d.Set("description", database.Description)
	d.Set("tags", database.Tags)

	d.Set("zone", client.Zone)
	return nil
}

func resourceSakuraCloudDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	database, err := client.Database.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Database resource: %s", err)
	}

	if d.HasChange("name") {
		database.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			database.Description = description.(string)
		} else {
			database.Description = ""
		}
	}

	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			database.Tags = expandStringList(rawTags)
		} else {
			database.Tags = []string{}
		}
	}

	database, err = client.Database.Update(database.ID, database)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Database resource: %s", err)
	}

	if d.HasChange("user_password") {
		database.Settings.DBConf.Common.UserPassword = d.Get("user_password").(string)
	}
	if d.HasChange("allow_networks") {
		if rawNetworks, ok := d.GetOk("allow_networks"); ok {
			if rawNetworks != nil {
				database.Settings.DBConf.Common.SourceNetwork = expandStringList(rawNetworks.([]interface{}))
			} else {
				database.Settings.DBConf.Common.SourceNetwork = nil
			}
		}

	}
	if d.HasChange("port") {
		database.Settings.DBConf.Common.ServicePort = fmt.Sprintf("%d", d.Get("port").(int))
	}
	if d.HasChange("backup_rotate") {
		database.Settings.DBConf.Backup.Rotate = fmt.Sprintf("%d", d.Get("backup_rotate").(int))
	}
	if d.HasChange("backup_time") {
		database.Settings.DBConf.Backup.Time = d.Get("backup_time").(string)
	}

	database, err = client.Database.UpdateSetting(database.ID, database)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Database resource: %s", err)
	}
	_, err = client.Database.Config(database.ID)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Database resource: %s", err)
	}

	d.SetId(database.ID)

	return resourceSakuraCloudDatabaseRead(d, meta)
}

func resourceSakuraCloudDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	time.Sleep(2 * time.Second)
	_, err := client.Database.Stop(d.Id())
	if err != nil {
		return fmt.Errorf("Error stopping SakuraCloud Database resource: %s", err)
	}

	err = client.Database.SleepUntilDown(d.Id(), 20*time.Minute)
	if err != nil {
		return fmt.Errorf("Error stopping SakuraCloud Database resource: %s", err)
	}

	_, err = client.Database.Delete(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Database resource: %s", err)
	}

	return nil
}
