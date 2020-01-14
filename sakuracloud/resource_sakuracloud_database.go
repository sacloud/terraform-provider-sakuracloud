// Copyright 2016-2019 terraform-provider-sakuracloud authors
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
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/power"
)

func resourceSakuraCloudDatabase() *schema.Resource {
	resourceName := "Database"
	return &schema.Resource{
		Create: resourceSakuraCloudDatabaseCreate,
		Read:   resourceSakuraCloudDatabaseRead,
		Update: resourceSakuraCloudDatabaseUpdate,
		Delete: resourceSakuraCloudDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"database_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(types.RDBMSTypeStrings, false),
				Default:      "postgres",
				Description: descf(
					"The type of the database. This must be one of [%s]",
					types.RDBMSTypeStrings,
				),
			},
			"plan": schemaResourcePlan(resourceName, "10g", types.DatabasePlanStrings),
			"username": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The name of default user on the database",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The password of default user on the database",
			},
			"replica_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "replica",
				Description: "The name of user that processing a replication",
			},
			"replica_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password of user that processing a replication",
			},
			"source_ranges": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: descf(
					"The range of source IP addresses that allow to access to the %s via network",
					resourceName,
				),
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5432,
				ValidateFunc: validation.IntBetween(1024, 65535),
				Description: descf(
					"The number of the listening port. %s",
					descRange(1024, 65535),
				),
			},
			"backup_weekdays": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: descf(
					"A list of weekdays to backed up. The values in the list must be in [%s]",
					types.ValidAutoBackupWeekdaysInString,
				),
			},
			"backup_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateBackupTime(),
				Description:  "The time to take backup. This must be formatted with `HH:mm`",
			},
			"switch_id": schemaResourceSwitchID(resourceName),
			"ip_addresses": {
				Type:        schema.TypeList,
				ForceNew:    true,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: descf("A list of IP address to assign to the %s", resourceName),
			},
			"netmask": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.IntBetween(8, 29),
				Description: descf(
					"The bit length of the subnet to assign to the %s. %s",
					resourceName,
					descRange(8, 29),
				),
			},
			"gateway": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: descf("The IP address of the gateway used by %s", resourceName),
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"zone":        schemaResourceZone(resourceName),
		},
	}
}

func resourceSakuraCloudDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	if err := validateDatabaseParameters(d); err != nil {
		return err
	}

	dbBuilder := expandDatabaseBuilder(d, client)
	db, err := dbBuilder.Build(ctx, zone)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Database is failed: %s", err)
	}

	// HACK データベースアプライアンスの電源投入後すぐに他の操作(Updateなど)を行うと202(Accepted)が返ってくるものの無視される。
	// この挙動はテストなどで問題となる。このためここで少しsleepすることで対応する。
	time.Sleep(client.databaseWaitAfterCreateDuration)

	d.SetId(db.ID.String())
	return resourceSakuraCloudDatabaseRead(d, meta)
}

func resourceSakuraCloudDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	dbOp := sacloud.NewDatabaseOp(client)

	data, err := dbOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not find SakuraCloud Database[%s]: %s", d.Id(), err)
	}
	return setDatabaseResourceData(ctx, d, client, data)
}

func resourceSakuraCloudDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	dbOp := sacloud.NewDatabaseOp(client)

	db, err := dbOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Database[%s]: %s", d.Id(), err)
	}

	dbBuilder := expandDatabaseBuilder(d, client)
	if _, err := dbBuilder.Update(ctx, zone, db.ID); err != nil {
		return fmt.Errorf("updating SakuraCloud Database[%s] is failed: %s", db.ID, err)
	}

	return resourceSakuraCloudDatabaseRead(d, meta)
}

func resourceSakuraCloudDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	dbOp := sacloud.NewDatabaseOp(client)

	data, err := dbOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Database[%s]: %s", d.Id(), err)
	}

	if data.InstanceStatus.IsUp() {
		if err := power.ShutdownDatabase(ctx, dbOp, zone, data.ID, true); err != nil {
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

	d.Set("database_type", flattenDatabaseType(data)) // nolint

	if data.ReplicationSetting != nil {
		d.Set("replica_user", data.CommonSetting.ReplicaUser)         // nolint
		d.Set("replica_password", data.CommonSetting.ReplicaPassword) // nolint
	}

	if data.BackupSetting != nil {
		d.Set("backup_time", data.BackupSetting.Time) // nolint
		if err := d.Set("backup_weekdays", data.BackupSetting.DayOfWeek); err != nil {
			return err
		}
	}

	if err := d.Set("tags", flattenDatabaseTags(data)); err != nil {
		return err
	}

	d.Set("name", data.Name)                              // nolint
	d.Set("username", data.CommonSetting.DefaultUser)     // nolint
	d.Set("password", data.CommonSetting.UserPassword)    // nolint
	d.Set("plan", types.DatabasePlanNameMap[data.PlanID]) // nolint
	if err := d.Set("source_ranges", data.CommonSetting.SourceNetwork); err != nil {
		return err
	}
	d.Set("port", data.CommonSetting.ServicePort) // nolint
	d.Set("switch_id", data.SwitchID.String())    // nolint
	d.Set("netmask", data.NetworkMaskLen)         // nolint
	d.Set("gateway", data.DefaultRoute)           // nolint
	if err := d.Set("ip_addresses", data.IPAddresses); err != nil {
		return err
	}
	d.Set("icon_id", data.IconID.String()) // nolint
	d.Set("description", data.Description) // nolint
	d.Set("zone", getZone(d, client))      // nolint

	return nil
}
