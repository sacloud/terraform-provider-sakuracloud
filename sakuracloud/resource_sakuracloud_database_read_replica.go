package sakuracloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/setup"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudDatabaseReadReplica() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudDatabaseReadReplicaCreate,
		Read:   resourceSakuraCloudDatabaseReadReplicaRead,
		Update: resourceSakuraCloudDatabaseReadReplicaUpdate,
		Delete: resourceSakuraCloudDatabaseReadReplicaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"master_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"switch_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
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
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(8, 29),
			},
			"default_route": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"allow_networks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func resourceSakuraCloudDatabaseReadReplicaCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	dbOp := sacloud.NewDatabaseOp(client)

	// validate master instance
	masterID := d.Get("master_id").(string)
	masterDB, err := dbOp.Read(ctx, zone, types.StringID(masterID))
	if err != nil {
		return fmt.Errorf("master database instance[%s] is not found", masterID)
	}
	if masterDB.ReplicationSetting.Model != types.DatabaseReplicationModels.MasterSlave {
		return fmt.Errorf("master database instance[%s] is not configured as ReplicationMaster", masterID)
	}

	switchID := masterDB.SwitchID.String()
	if v, ok := d.GetOk("switch_id"); ok {
		switchID = v.(string)
	}
	maskLen := masterDB.NetworkMaskLen
	if v, ok := d.GetOk("nw_mask_len"); ok {
		maskLen = v.(int)
	}
	defaultRoute := masterDB.DefaultRoute
	if v, ok := d.GetOk("default_route"); ok {
		defaultRoute = v.(string)
	}

	req := &sacloud.DatabaseCreateRequest{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Tags:           expandTagsV2(d.Get("tags").([]interface{})),
		IconID:         types.StringID(d.Get("icon_id").(string)),
		PlanID:         types.ID(masterDB.PlanID.Int64() + 1),
		SwitchID:       types.StringID(switchID),
		IPAddresses:    []string{d.Get("ipaddress1").(string)},
		NetworkMaskLen: maskLen,
		DefaultRoute:   defaultRoute,
		Conf: &sacloud.DatabaseRemarkDBConfCommon{
			DatabaseName:     masterDB.Conf.DatabaseName,
			DatabaseVersion:  masterDB.Conf.DatabaseVersion,
			DatabaseRevision: masterDB.Conf.DatabaseRevision,
		},
		CommonSetting: &sacloud.DatabaseSettingCommon{
			ServicePort:   masterDB.CommonSetting.ServicePort,
			SourceNetwork: expandStringList(d.Get("allow_networks").([]interface{})),
		},
		ReplicationSetting: &sacloud.DatabaseReplicationSetting{
			Model:       types.DatabaseReplicationModels.AsyncReplica,
			IPAddress:   masterDB.IPAddresses[0],
			Port:        masterDB.CommonSetting.ServicePort,
			User:        masterDB.ReplicationSetting.User,
			Password:    masterDB.ReplicationSetting.Password,
			ApplianceID: masterDB.ID,
		},
	}

	dbBuilder := &setup.RetryableSetup{
		IsWaitForCopy: true,
		IsWaitForUp:   true,
		Create: func(ctx context.Context, zone string) (accessor.ID, error) {
			return dbOp.Create(ctx, zone, req)
		},
		Read: func(ctx context.Context, zone string, id types.ID) (interface{}, error) {
			return dbOp.Read(ctx, zone, id)
		},
		Delete: func(ctx context.Context, zone string, id types.ID) error {
			return dbOp.Delete(ctx, zone, id)
		},
		RetryCount: 3,
	}

	res, err := dbBuilder.Setup(ctx, zone)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Database ReadReplica is failed: %s", err)
	}

	database, ok := res.(*sacloud.Database)
	if !ok {
		return fmt.Errorf("creating SakuraCloud Database ReadReplica is failed: resource is not *sacloud.Database")
	}

	// HACK データベースアプライアンスの電源投入後すぐに他の操作(Updateなど)を行うと202(Accepted)が返ってくるものの無視される。
	// この挙動はテストなどで問題となる。このためここで少しsleepすることで対応する。
	time.Sleep(1 * time.Minute)

	d.SetId(database.ID.String())
	return setDatabaseReadReplicaResourceData(ctx, d, client, database)
}

func resourceSakuraCloudDatabaseReadReplicaRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	dbOp := sacloud.NewDatabaseOp(client)

	data, err := dbOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not find SakuraCloud Database ReadReplica resource: %s", err)
	}
	return setDatabaseReadReplicaResourceData(ctx, d, client, data)
}

func resourceSakuraCloudDatabaseReadReplicaUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	dbOp := sacloud.NewDatabaseOp(client)

	db, err := dbOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Database[%s]: %s", d.Id(), err)
	}

	req := &sacloud.DatabasePatchRequest{
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		Tags:         expandTagsV2(d.Get("tags").([]interface{})),
		IconID:       types.StringID(d.Get("icon_id").(string)),
		SettingsHash: db.SettingsHash,
	}
	db, err = dbOp.Patch(ctx, zone, db.ID, req)
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Database[%s] is failed: %s", d.Id(), err)
	}

	return setDatabaseReadReplicaResourceData(ctx, d, client, db)
}

func resourceSakuraCloudDatabaseReadReplicaDelete(d *schema.ResourceData, meta interface{}) error {
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

func setDatabaseReadReplicaResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Database) error {

	if data.Availability.IsFailed() {
		d.SetId("")
		return fmt.Errorf("got unexpected state: Database[%d].Availability is failed", data.ID)
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

	d.Set("master_id", data.ReplicationSetting.ApplianceID.String())
	d.Set("name", data.Name)
	d.Set("switch_id", data.SwitchID.String())
	d.Set("nw_mask_len", data.NetworkMaskLen)
	d.Set("default_route", data.DefaultRoute)
	d.Set("ipaddress1", data.IPAddresses[0])
	if err := d.Set("allow_networks", data.CommonSetting.SourceNetwork); err != nil {
		return fmt.Errorf("error setting allow_networks: %v", data.CommonSetting.SourceNetwork)
	}
	if !data.IconID.IsEmpty() {
		d.Set("icon_id", data.IconID.String())
	}
	d.Set("description", data.Description)
	d.Set("zone", getV2Zone(d, client))

	return nil
}
