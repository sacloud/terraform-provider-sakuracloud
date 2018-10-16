package sakuracloud

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/utils/setup"
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
				ValidateFunc: validateZone([]string{"tk1a", "is1b", "is1a"}),
			},
		},
	}
}

func resourceSakuraCloudDatabaseReadReplicaCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	// validate master instance
	masterID := d.Get("master_id").(string)
	masterDB, err := client.Database.Read(toSakuraCloudID(masterID))
	if err != nil {
		return fmt.Errorf("master database instance[%s] is not found", masterID)
	}
	if !masterDB.IsReplicationMaster() {
		return fmt.Errorf("master database instance[%s] is not configured as ReplicationMaster", masterID)
	}

	servicePort := masterDB.Settings.DBConf.Common.ServicePort
	port, err := servicePort.Int64()
	if servicePort.String() != "" && err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Database ReadReplica resource: %s", err)
	}

	switchID := masterDB.Switch.GetStrID()
	if v, ok := d.GetOk("switch_id"); ok {
		switchID = v.(string)
	}
	maskLen := masterDB.Remark.Network.NetworkMaskLen
	if v, ok := d.GetOk("nw_mask_len"); ok {
		maskLen = v.(int)
	}
	defaultRoute := masterDB.Remark.Network.DefaultRoute
	if v, ok := d.GetOk("default_route"); ok {
		defaultRoute = v.(string)
	}

	var opts = &sacloud.SlaveDatabaseValue{
		Plan:              sacloud.DatabasePlan(masterDB.Plan.ID),
		DefaultUser:       masterDB.Settings.DBConf.Common.DefaultUser,
		UserPassword:      masterDB.Settings.DBConf.Common.UserPassword,
		SwitchID:          switchID,
		IPAddress1:        d.Get("ipaddress1").(string),
		MaskLen:           maskLen,
		DefaultRoute:      defaultRoute,
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		Tags:              expandTags(client, d.Get("tags").([]interface{})),
		Icon:              sacloud.NewResource(toSakuraCloudID(d.Get("icon_id").(string))),
		DatabaseName:      masterDB.Remark.DBConf.Common.DatabaseName,
		DatabaseVersion:   masterDB.Remark.DBConf.Common.DatabaseVersion,
		ReplicaPassword:   masterDB.Settings.DBConf.Common.ReplicaPassword,
		MasterApplianceID: masterDB.ID,
		MasterIPAddress:   masterDB.Remark.Servers[0].(map[string]interface{})["IPAddress"].(string),
		MasterPort:        int(port),
	}

	createDB := sacloud.NewSlaveDatabaseValue(opts)
	dbBuilder := &setup.RetryableSetup{
		Create: func() (sacloud.ResourceIDHolder, error) {
			return client.Database.Create(createDB)
		},
		AsyncWaitForCopy: func(id int64) (chan interface{}, chan interface{}, chan error) {
			return client.Database.AsyncSleepWhileCopying(id, client.DefaultTimeoutDuration, 5)
		},
		Delete: func(id int64) error {
			_, err := client.Database.Delete(id)
			return err
		},
		WaitForUp: func(id int64) error {
			return client.Database.SleepUntilUp(id, client.DefaultTimeoutDuration)
		},
		RetryCount: 3,
	}

	res, err := dbBuilder.Setup()
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Database ReadReplica resource: %s", err)
	}

	database, ok := res.(*sacloud.Database)
	if !ok {
		return fmt.Errorf("Failed to create SakuraCloud Database ReadReplica resource: created resource is not *sacloud.Database")
	}

	err = client.Database.SleepUntilDatabaseRunning(database.ID, client.DefaultTimeoutDuration, 5)
	if err != nil {
		return fmt.Errorf("Failed to wait SakuraCloud Database ReadReplica start: %s", err)
	}

	d.SetId(database.GetStrID())
	return resourceSakuraCloudDatabaseReadReplicaRead(d, meta)
}

func resourceSakuraCloudDatabaseReadReplicaRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	data, err := client.Database.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud Database ReadReplica resource: %s", err)
	}

	// validate is replication slave?
	if data.Settings.DBConf.Replication == nil || data.Settings.DBConf.Replication.Appliance == nil {
		return fmt.Errorf("database instance[%s] is not configured as ReplicationSlave", d.Id())
	}

	return setDatabaseReadReplicaResourceData(d, client, data)
}

func resourceSakuraCloudDatabaseReadReplicaUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	database, err := client.Database.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Database ReadReplica resource: %s", err)
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
		return fmt.Errorf("Error updating SakuraCloud Database ReadReplica resource: %s", err)
	}

	return resourceSakuraCloudDatabaseReadReplicaRead(d, meta)
}

func resourceSakuraCloudDatabaseReadReplicaDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	err := handleShutdown(client.Database, toSakuraCloudID(d.Id()), d, client.DefaultTimeoutDuration)
	if err != nil {
		return fmt.Errorf("Error stopping SakuraCloud Database ReadReplica resource: %s", err)
	}

	_, err = client.Database.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Database ReadReplica resource: %s", err)
	}

	return nil
}

func setDatabaseReadReplicaResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.Database) error {

	if data.IsFailed() {
		d.SetId("")
		return fmt.Errorf("Database[%d] state is failed", data.ID)
	}

	d.Set("name", data.Name)
	d.Set("switch_id", data.Interfaces[0].Switch.GetStrID())
	d.Set("nw_mask_len", data.Remark.Network.NetworkMaskLen)
	d.Set("default_route", data.Remark.Network.DefaultRoute)
	d.Set("ipaddress1", data.Remark.Servers[0].(map[string]interface{})["IPAddress"])
	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	d.Set("master_id", data.Settings.DBConf.Replication.Appliance.GetStrID())
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
