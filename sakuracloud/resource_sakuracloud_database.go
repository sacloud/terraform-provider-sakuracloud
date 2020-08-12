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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/utils/setup"
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
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"postgresql", "mariadb"}, false),
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
				ValidateFunc: validation.StringInSlice([]string{"10g", "30g", "90g", "240g", "500g", "1t"}, false),
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
			"replica_user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replica_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
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
				ValidateFunc: validation.IntBetween(1024, 65535),
			},
			"backup_weekdays": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"backup_time": {
				Type:         schema.TypeString,
				Optional:     true,
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
				ValidateFunc: validation.IntBetween(8, 29),
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
				ValidateFunc: validateZone([]string{"tk1a", "tk1b", "is1b", "is1a"}),
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
	case "mariadb":
		opts = sacloud.NewCreateMariaDBDatabaseValue()
	default:
		return fmt.Errorf("Unknown database_type [%s]", dbType)
	}

	opts.Name = d.Get("name").(string)
	opts.DefaultUser = d.Get("user_name").(string)
	opts.UserPassword = d.Get("user_password").(string)

	replicaPassword := d.Get("replica_password").(string)
	if replicaPassword != "" {
		opts.ReplicaPassword = replicaPassword
	}

	if rawNetworks, ok := d.GetOk("allow_networks"); ok {
		if rawNetworks != nil {
			opts.SourceNetwork = expandStringList(rawNetworks.([]interface{}))
		}
	}
	opts.ServicePort = d.Get("port").(int)

	opts.BackupTime = d.Get("backup_time").(string)
	rawBackupWeekdays := d.Get("backup_weekdays").([]interface{})
	backupWeekdays, err := expandStringListWithValidateInList("backup_weekdays", rawBackupWeekdays, sacloud.AllowDatabaseBackupWeekdays())
	if err != nil {
		return err
	}
	opts.BackupDayOfWeek = backupWeekdays
	if opts.BackupTime != "" && len(backupWeekdays) > 0 {
		opts.EnableBackup = true
	}

	opts.SwitchID = sacloud.StringID(d.Get("switch_id").(string))
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
	case "500g":
		opts.Plan = sacloud.DatabasePlan500G
	case "1t":
		opts.Plan = sacloud.DatabasePlan1T
	}

	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.Icon = sacloud.NewResource(toSakuraCloudID(iconID.(string)))
	}

	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	if rawTags, ok := d.GetOk("tags"); ok {
		if rawTags != nil {
			opts.Tags = expandTags(client, rawTags.([]interface{}))
		}
	}

	createDB := sacloud.CreateNewDatabase(opts)

	dbBuilder := &setup.RetryableSetup{
		Create: func() (sacloud.ResourceIDHolder, error) {
			return client.Database.Create(createDB)
		},
		AsyncWaitForCopy: func(id sacloud.ID) (chan interface{}, chan interface{}, chan error) {
			return client.Database.AsyncSleepWhileCopying(id, client.DefaultTimeoutDuration, 20)
		},
		Delete: func(id sacloud.ID) error {
			_, err := client.Database.Delete(id)
			return err
		},
		WaitForUp: func(id sacloud.ID) error {
			return client.Database.SleepUntilUp(id, client.DefaultTimeoutDuration)
		},
		RetryCount: 3,
	}

	res, err := dbBuilder.Setup()
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Database resource: %s", err)
	}

	database, ok := res.(*sacloud.Database)
	if !ok {
		return fmt.Errorf("Failed to create SakuraCloud Database resource: created resource is not *sacloud.Database")
	}

	err = client.Database.SleepUntilDatabaseRunning(database.ID, client.DefaultTimeoutDuration, 5)
	if err != nil {
		return fmt.Errorf("Failed to wait SakuraCloud Database start: %s", err)
	}

	d.SetId(database.GetStrID())
	return resourceSakuraCloudDatabaseRead(d, meta)
}

func resourceSakuraCloudDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	data, err := client.Database.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
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

	needRestart := false
	if d.HasChange("replica_password") {
		if database.IsUp() {
			needRestart = true
			err = handleShutdown(client.Database, database.ID, d, client.DefaultTimeoutDuration)
			if err != nil {
				return fmt.Errorf("Error updating SakuraCloud Database resource: %s", err)
			}
			err = client.Database.SleepUntilDown(database.ID, client.DefaultTimeoutDuration)
			if err != nil {
				return fmt.Errorf("Error updating SakuraCloud Database resource: %s", err)
			}
		}
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
			database.Tags = expandTags(client, rawTags)
		} else {
			database.Tags = expandTags(client, []interface{}{})
		}
	}

	database, err = client.Database.Update(database.ID, database)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Database resource: %s", err)
	}

	if d.HasChange("user_password") {
		database.Settings.DBConf.Common.UserPassword = d.Get("user_password").(string)
	}

	if d.HasChange("replica_password") {
		replicaPassword := d.Get("replica_password").(string)
		if replicaPassword == "" {
			database.Settings.DBConf.Common.ReplicaPassword = ""
			database.Settings.DBConf.Replication = nil
		} else {
			database.Settings.DBConf.Common.ReplicaPassword = replicaPassword
			database.Settings.DBConf.Replication = &sacloud.DatabaseReplicationSetting{
				Model: sacloud.DatabaseReplicationModelMasterSlave,
			}
		}
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
		rawPort := d.Get("port").(int)
		database.Settings.DBConf.Common.ServicePort = json.Number(fmt.Sprintf("%d", rawPort))
	}
	if d.HasChange("backup_weekdays") {
		rawBackupWeekdays := d.Get("backup_weekdays").([]interface{})
		backupWeekdays, err := expandStringListWithValidateInList("backup_weekdays", rawBackupWeekdays, sacloud.AllowDatabaseBackupWeekdays())
		if err != nil {
			return err
		}

		if database.Settings.DBConf.Backup == nil {
			database.Settings.DBConf.Backup = &sacloud.DatabaseBackupSetting{}
		}
		database.Settings.DBConf.Backup.DayOfWeek = backupWeekdays
	}

	if d.HasChange("backup_time") {
		backupTime := d.Get("backup_time").(string)
		if database.Settings.DBConf.Backup == nil {
			database.Settings.DBConf.Backup = &sacloud.DatabaseBackupSetting{}
		}
		database.Settings.DBConf.Backup.Time = backupTime
	}
	if database.Settings.DBConf.Backup != nil &&
		(len(database.Settings.DBConf.Backup.DayOfWeek) == 0 || database.Settings.DBConf.Backup.Time == "") {
		database.Settings.DBConf.Backup = nil
	}

	database, err = client.Database.UpdateSetting(database.ID, database)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Database resource: %s", err)
	}
	_, err = client.Database.Config(database.ID)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Database resource: %s", err)
	}

	if needRestart {
		_, err = client.Database.Boot(database.ID)
		if err != nil {
			return fmt.Errorf("Error updating SakuraCloud Database resource: %s", err)
		}
		err = client.Database.SleepUntilDatabaseRunning(database.ID, client.DefaultTimeoutDuration, 5)
		if err != nil {
			return fmt.Errorf("Error updating SakuraCloud Database resource: %s", err)
		}
	}

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

func setDatabaseResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.Database) error {
	if data.IsFailed() {
		d.SetId("")
		return fmt.Errorf("Database[%d] state is failed", data.ID)
	}

	switch data.Remark.DBConf.Common.DatabaseName {
	case "postgres":
		d.Set("database_type", "postgresql")
	case "MariaDB":
		d.Set("database_type", "mariadb")
	}

	d.Set("name", data.Name)
	d.Set("user_name", data.Settings.DBConf.Common.DefaultUser)
	d.Set("user_password", data.Settings.DBConf.Common.UserPassword)
	d.Set("replica_user", data.Settings.DBConf.Common.ReplicaUser)
	d.Set("replica_password", data.Settings.DBConf.Common.ReplicaPassword)

	//plan
	switch data.Plan.ID.Int64() {
	case int64(sacloud.DatabasePlan10G):
		d.Set("plan", "10g")
	case int64(sacloud.DatabasePlan30G):
		d.Set("plan", "30g")
	case int64(sacloud.DatabasePlan90G):
		d.Set("plan", "90g")
	case int64(sacloud.DatabasePlan240G):
		d.Set("plan", "240g")
	case int64(sacloud.DatabasePlan500G):
		d.Set("plan", "500g")
	case int64(sacloud.DatabasePlan1T):
		d.Set("plan", "1t")
	}

	d.Set("allow_networks", data.Settings.DBConf.Common.SourceNetwork)
	port, _ := data.Settings.DBConf.Common.ServicePort.Int64()
	d.Set("port", port)

	if data.Settings.DBConf.Backup != nil {
		d.Set("backup_time", data.Settings.DBConf.Backup.Time)
		d.Set("backup_weekdays", data.Settings.DBConf.Backup.DayOfWeek)
	}

	d.Set("switch_id", "")
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
	setPowerManageTimeoutValueToState(d)

	d.Set("zone", client.Zone)
	return nil
}
