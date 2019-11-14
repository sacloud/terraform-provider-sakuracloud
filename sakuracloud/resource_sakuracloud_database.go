package sakuracloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/setup"
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
				ValidateFunc: validation.StringInSlice([]string{"postgresql", "mariadb"}, false),
				Default:      "postgresql",
			},
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
				Optional: true,
				Default:  "replica",
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
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"tk1a", "is1b", "is1a"}),
			},
		},
	}
}

func resourceSakuraCloudDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	dbOp := sacloud.NewDatabaseOp(client)

	if err := validateBackupWeekdays(d, "backup_weekdays"); err != nil {
		return err
	}

	var dbVersion *types.RDBMSVersion
	dbType := d.Get("database_type").(string)
	switch dbType {
	case "postgresql":
		dbVersion = types.RDBMSVersions[types.RDBMSTypesPostgreSQL]
	case "mariadb":
		dbVersion = types.RDBMSVersions[types.RDBMSTypesMariaDB]
	default:
		return fmt.Errorf("unknown database_type[%s]", dbType)
	}

	replicaUser := d.Get("replica_user").(string)
	replicaPassword := d.Get("replica_password").(string)

	req := &sacloud.DatabaseCreateRequest{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Tags:           expandTagsV2(d.Get("tags").([]interface{})),
		IconID:         types.StringID(d.Get("icon_id").(string)),
		PlanID:         databasePlanNameToID(d.Get("plan").(string)),
		SwitchID:       types.StringID(d.Get("switch_id").(string)),
		IPAddresses:    []string{d.Get("ipaddress1").(string)},
		NetworkMaskLen: d.Get("nw_mask_len").(int),
		DefaultRoute:   d.Get("default_route").(string),
		Conf: &sacloud.DatabaseRemarkDBConfCommon{
			DatabaseName:     dbVersion.Name,
			DatabaseVersion:  dbVersion.Version,
			DatabaseRevision: dbVersion.Revision,
			DefaultUser:      d.Get("user_name").(string),
			UserPassword:     d.Get("user_password").(string),
		},
		CommonSetting: &sacloud.DatabaseSettingCommon{
			ServicePort:     d.Get("port").(int),
			SourceNetwork:   expandStringList(d.Get("allow_networks").([]interface{})),
			DefaultUser:     d.Get("user_name").(string),
			UserPassword:    d.Get("user_password").(string),
			ReplicaUser:     replicaUser,
			ReplicaPassword: replicaPassword,
		},
	}

	backupTime := d.Get("backup_time").(string)
	backupWeekdays := expandBackupWeekdays(d.Get("backup_weekdays").([]interface{}))
	if backupTime != "" && len(backupWeekdays) > 0 {
		req.BackupSetting = &sacloud.DatabaseSettingBackup{
			Time:      backupTime,
			DayOfWeek: backupWeekdays,
		}
	}

	if replicaUser != "" && replicaPassword != "" {
		req.ReplicationSetting = &sacloud.DatabaseReplicationSetting{
			Model:    types.DatabaseReplicationModels.MasterSlave,
			User:     replicaUser,
			Password: replicaPassword,
		}
	}

	dbBuilder := &setup.RetryableSetup{
		IsWaitForCopy: true,
		IsWaitForUp:   true,
		Create: func(ctx context.Context, zone string) (accessor.ID, error) {
			return dbOp.Create(ctx, zone, req)
		},
		Delete: func(ctx context.Context, zone string, id types.ID) error {
			return dbOp.Delete(ctx, zone, id)
		},
		Read: func(ctx context.Context, zone string, id types.ID) (interface{}, error) {
			return dbOp.Read(ctx, zone, id)
		},
		RetryCount: 3,
	}

	res, err := dbBuilder.Setup(ctx, zone)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Database is failed: %s", err)
	}
	db := res.(*sacloud.Database)

	// HACK データベースアプライアンスの電源投入後すぐに他の操作(Updateなど)を行うと202(Accepted)が返ってくるものの無視される。
	// この挙動はテストなどで問題となる。このためここで少しsleepすることで対応する。
	time.Sleep(1 * time.Minute)

	d.SetId(db.ID.String())
	return resourceSakuraCloudDatabaseRead(d, meta)
}

func resourceSakuraCloudDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	dbOp := sacloud.NewDatabaseOp(client)

	data, err := dbOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not find SakuraCloud Database resource: %s", err)
	}
	return setDatabaseResourceData(ctx, d, client, data)
}

func resourceSakuraCloudDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	dbOp := sacloud.NewDatabaseOp(client)

	db, err := dbOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Database[%s]: %s", d.Id(), err)
	}

	needRestart := false
	if d.HasChange("replica_password") && db.InstanceStatus.IsUp() {
		if err := shutdownDatabase(ctx, dbOp, zone, db.ID, false); err != nil {
			return err
		}
		needRestart = true
	}

	replicaUser := d.Get("replica_user").(string)
	replicaPassword := d.Get("replica_password").(string)

	req := &sacloud.DatabaseUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTagsV2(d.Get("tags").([]interface{})),
		IconID:      types.StringID(d.Get("icon_id").(string)),
		CommonSetting: &sacloud.DatabaseSettingCommon{
			ServicePort:     d.Get("port").(int),
			SourceNetwork:   expandStringList(d.Get("allow_networks").([]interface{})),
			DefaultUser:     d.Get("user_name").(string),
			UserPassword:    d.Get("user_password").(string),
			ReplicaUser:     replicaUser,
			ReplicaPassword: replicaPassword,
		},
		BackupSetting:      &sacloud.DatabaseSettingBackup{},
		ReplicationSetting: &sacloud.DatabaseReplicationSetting{},
		SettingsHash:       db.SettingsHash,
	}
	backupTime := d.Get("backup_time").(string)
	backupWeekdays := expandBackupWeekdays(d.Get("backup_weekdays").([]interface{}))
	if backupTime != "" && len(backupWeekdays) > 0 {
		req.BackupSetting = &sacloud.DatabaseSettingBackup{
			Time:      backupTime,
			DayOfWeek: backupWeekdays,
		}
	}

	if replicaUser != "" && replicaPassword != "" {
		req.ReplicationSetting = &sacloud.DatabaseReplicationSetting{
			Model: types.DatabaseReplicationModels.MasterSlave,
		}
	}

	db, err = dbOp.Update(ctx, zone, db.ID, req)
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Database[%s] is failed: %s", d.Id(), err)
	}

	if needRestart && !db.InstanceStatus.IsUp() {
		if err := bootDatabase(ctx, dbOp, zone, db.ID); err != nil {
			return err
		}
	}
	return resourceSakuraCloudDatabaseRead(d, meta)
}

func resourceSakuraCloudDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	dbOp := sacloud.NewDatabaseOp(client)

	data, err := dbOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Database resource: %s", err)
	}

	// shutdown(force) if running
	if data.InstanceStatus.IsUp() {
		if err := shutdownDatabase(ctx, dbOp, zone, data.ID, true); err != nil {
			return err
		}
	}

	// delete
	if err = dbOp.Delete(ctx, zone, data.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud Database[%s] is failed: %s", data.ID, err)
	}

	d.SetId("")
	return nil
}

func setDatabaseResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Database) error {
	if data.Availability.IsFailed() {
		d.SetId("")
		return fmt.Errorf("got unexpected state: Database[%d].Availability is failed", data.ID)
	}

	var databaseType string
	switch data.Conf.DatabaseName {
	case types.RDBMSVersions[types.RDBMSTypesPostgreSQL].Name:
		databaseType = "postgresql"
	case types.RDBMSVersions[types.RDBMSTypesMariaDB].Name:
		databaseType = "mariadb"
	}
	d.Set("database_type", databaseType)

	if data.ReplicationSetting != nil {
		d.Set("replica_user", data.CommonSetting.ReplicaUser)
		d.Set("replica_password", data.CommonSetting.ReplicaPassword)
	}

	if data.BackupSetting != nil {
		d.Set("backup_time", data.BackupSetting.Time)
		if err := d.Set("backup_weekdays", data.BackupSetting.DayOfWeek); err != nil {
			return fmt.Errorf("error setting backup_weekdays: %v", data.BackupSetting.DayOfWeek)
		}
	}

	var tags []string
	for _, t := range data.Tags {
		if !(strings.HasPrefix(t, "@MariaDB-") || strings.HasPrefix(t, "@postgres-")) {
			tags = append(tags, t)
		}
	}
	if err := d.Set("tags", tags); err != nil {
		return fmt.Errorf("error setting tags: %v", tags)
	}

	d.Set("name", data.Name)
	d.Set("user_name", data.CommonSetting.DefaultUser)
	d.Set("user_password", data.CommonSetting.UserPassword)
	d.Set("plan", databasePlanIDToName(data.PlanID))
	if err := d.Set("allow_networks", data.CommonSetting.SourceNetwork); err != nil {
		return fmt.Errorf("error setting allow_networks: %v", data.CommonSetting.SourceNetwork)
	}
	d.Set("port", data.CommonSetting.ServicePort)
	d.Set("switch_id", data.SwitchID.String())
	d.Set("nw_mask_len", data.NetworkMaskLen)
	d.Set("default_route", data.DefaultRoute)
	d.Set("ipaddress1", data.IPAddresses[0])
	if !data.IconID.IsEmpty() {
		d.Set("icon_id", data.IconID.String())
	}
	d.Set("description", data.Description)
	d.Set("zone", getV2Zone(d, client))

	return nil
}

func bootDatabase(ctx context.Context, dbOp sacloud.DatabaseAPI, zone string, id types.ID) error {
	// boot
	if err := dbOp.Boot(ctx, zone, id); err != nil {
		return fmt.Errorf("booting Database[%s] is failed: %s", id, err)
	}

	// wait
	waiter := sacloud.WaiterForUp(func() (interface{}, error) {
		return dbOp.Read(ctx, zone, id)
	})
	if _, err := waiter.WaitForState(ctx); err != nil {
		return fmt.Errorf("waiting for Database[%s] up is failed: %s", id, err)
	}
	return nil
}

func shutdownDatabase(ctx context.Context, dbOp sacloud.DatabaseAPI, zone string, id types.ID, forceShutdown bool) error {
	// shutdown
	if err := dbOp.Shutdown(ctx, zone, id, &sacloud.ShutdownOption{Force: forceShutdown}); err != nil {
		return fmt.Errorf("stopping Database[%s] is failed: %s", id, err)
	}

	// wait
	waiter := sacloud.WaiterForDown(func() (interface{}, error) {
		return dbOp.Read(ctx, zone, id)
	})
	if _, err := waiter.WaitForState(ctx); err != nil {
		return fmt.Errorf("waiting for Database[%s] down is failed: %s", id, err)
	}
	return nil
}

func databasePlanIDToName(planID types.ID) string {
	switch planID {
	case types.DatabasePlans.DB10GB:
		return "10g"
	case types.DatabasePlans.DB30GB:
		return "30g"
	case types.DatabasePlans.DB90GB:
		return "90g"
	case types.DatabasePlans.DB240GB:
		return "240g"
	case types.DatabasePlans.DB500GB:
		return "500g"
	case types.DatabasePlans.DB1TB:
		return "1t"
	}
	return ""
}

func databasePlanNameToID(planName string) types.ID {
	switch planName {
	case "10g":
		return types.DatabasePlans.DB10GB
	case "30g":
		return types.DatabasePlans.DB30GB
	case "90g":
		return types.DatabasePlans.DB90GB
	case "240g":
		return types.DatabasePlans.DB240GB
	case "500g":
		return types.DatabasePlans.DB500GB
	case "1t":
		return types.DatabasePlans.DB1TB
	}
	return types.ID(0)
}
