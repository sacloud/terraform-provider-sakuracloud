package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"strings"
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateStringInWord([]string{"postgresql", "mariadb"}),
				Default:      "postgresql",
			},
			//"is_double": {
			//	Type:     schema.TypeBool,
			//	ForceNew: true,
			//	Optional: true,
			//	Default:  false,
			//},
			"plan": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "10g",
				ValidateFunc: validateStringInWord([]string{"10g", "30g", "90g", "240g"}),
			},
			"user_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"user_password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"allow_networks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5432,
				ValidateFunc: validateIntegerInRange(1024, 65535),
			},
			"backup_time": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateBackupTime(),
			},

			"switch_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"ipaddress1": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"nw_mask_len": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateIntegerInRange(8, 29),
			},
			"default_route": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
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
			powerManageTimeoutKey: powerManageTimeoutParam,
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateStringInWord([]string{"tk1a", "is1b"}),
			},
		},
	}
}

func resourceSakuraCloudDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	var opts *sacloud.CreateDatabaseValue
	dbType := d.Get("database_type").(string)
	switch dbType {
	case "postgresql":
		opts = sacloud.NewCreatePostgreSQLDatabaseValue()
		break
	case "mariadb":
		opts = sacloud.NewCreateMariaDBDatabaseValue()
		break
	default:
		return fmt.Errorf("Unknown database_type [%s]", dbType)
	}

	opts.Name = d.Get("name").(string)
	opts.DefaultUser = d.Get("user_name").(string)
	opts.UserPassword = d.Get("user_password").(string)
	if rawNetworks, ok := d.GetOk("allow_networks"); ok {
		if rawNetworks != nil {
			opts.SourceNetwork = expandStringList(rawNetworks.([]interface{}))
		}
	}
	opts.ServicePort = fmt.Sprintf("%d", d.Get("port").(int))
	//opts.BackupRotate = d.Get("backup_rotate").(int)
	opts.BackupTime = d.Get("backup_time").(string)

	opts.SwitchID = d.Get("switch_id").(string)
	ipAddress1 := d.Get("ipaddress1").(string)
	nwMaskLen := d.Get("nw_mask_len").(int)
	defaultRoute := ""
	if df, ok := d.GetOk("default_route"); ok {
		defaultRoute = df.(string)
	}

	opts.IPAddress1 = ipAddress1
	opts.MaskLen = nwMaskLen
	opts.DefaultRoute = defaultRoute

	//
	strPlan := d.Get("plan").(string)
	switch strPlan {
	case "10g":
		opts.Plan = sacloud.DatabasePlan10G
	case "30g":
		opts.Plan = sacloud.DatabasePlan30G
	case "90g":
		opts.Plan = sacloud.DatabasePlan90G
	case "240g":
		opts.Plan = sacloud.DatabasePlan240G
	}

	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.Icon = sacloud.NewResource(toSakuraCloudID(iconID.(string)))
	}

	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	if rawTags, ok := d.GetOk("tags"); ok {
		if rawTags != nil {
			opts.Tags = expandStringList(rawTags.([]interface{}))
		}
	}

	createDB := sacloud.CreateNewDatabase(opts)
	database, err := client.Database.Create(createDB)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Database resource: %s", err)
	}

	d.SetId(database.GetStrID())

	//wait
	err = client.Database.SleepWhileCopying(database.ID, client.DefaultTimeoutDuration, 5)
	if err != nil {
		return fmt.Errorf("Failed to wait SakuraCloud Database copy: %s", err)
	}
	err = client.Database.SleepUntilUp(database.ID, client.DefaultTimeoutDuration)
	if err != nil {
		return fmt.Errorf("Failed to wait SakuraCloud Database boot: %s", err)
	}
	err = client.Database.SleepUntilDatabaseRunning(database.ID, client.DefaultTimeoutDuration, 5)
	if err != nil {
		return fmt.Errorf("Failed to wait SakuraCloud Database start: %s", err)
	}

	return resourceSakuraCloudDatabaseRead(d, meta)
}

func resourceSakuraCloudDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	data, err := client.Database.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Database resource: %s", err)
	}

	return setDatabaseResourceData(d, client, data)
}

func resourceSakuraCloudDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	database, err := client.Database.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Database resource: %s", err)
	}

	if d.HasChange("name") {
		database.Name = d.Get("name").(string)
	}
	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			database.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			database.ClearIcon()
		}
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
		database.Settings.DBConf.Backup.Rotate = d.Get("backup_rotate").(int)
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

	d.SetId(database.GetStrID())

	return resourceSakuraCloudDatabaseRead(d, meta)
}

func resourceSakuraCloudDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	err := handleShutdown(client.Database, toSakuraCloudID(d.Id()), d, client.DefaultTimeoutDuration)
	if err != nil {
		return fmt.Errorf("Error stopping SakuraCloud Database resource: %s", err)
	}

	_, err = client.Database.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Database resource: %s", err)
	}

	return nil
}

func setDatabaseResourceData(d *schema.ResourceData, client *api.Client, data *sacloud.Database) error {

	switch data.Remark.DBConf.Common.DatabaseName {
	case "postgres":
		d.Set("database_type", "postgresql")
		break
	case "MariaDB":
		d.Set("database_type", "mariadb")
		break
	}

	d.Set("name", data.Name)
	d.Set("user_name", data.Settings.DBConf.Common.DefaultUser)
	d.Set("user_password", data.Settings.DBConf.Common.UserPassword)

	//plan
	switch data.Plan.ID {
	case int64(sacloud.DatabasePlan10G):
		d.Set("plan", "10g")
	case int64(sacloud.DatabasePlan30G):
		d.Set("plan", "30g")
	case int64(sacloud.DatabasePlan90G):
		d.Set("plan", "90g")
	case int64(sacloud.DatabasePlan240G):
		d.Set("plan", "240g")

	}

	d.Set("allow_networks", data.Settings.DBConf.Common.SourceNetwork)
	d.Set("port", data.Settings.DBConf.Common.ServicePort)

	d.Set("backup_rotate", data.Settings.DBConf.Backup.Rotate)
	d.Set("backup_time", data.Settings.DBConf.Backup.Time)

	d.Set("switch_id", data.Interfaces[0].Switch.GetStrID())
	d.Set("nw_mask_len", data.Remark.Network.NetworkMaskLen)
	d.Set("default_route", data.Remark.Network.DefaultRoute)
	d.Set("ipaddress1", data.Remark.Servers[0].(map[string]interface{})["IPAddress"])

	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	tags := []string{}
	for _, t := range data.Tags {
		if !(strings.HasPrefix(t, "@MariaDB-") || strings.HasPrefix(t, "@postgres-")) {
			tags = append(tags, t)
		}
	}
	d.Set("tags", tags)

	d.Set("zone", client.Zone)
	d.SetId(data.GetStrID())
	return nil
}
