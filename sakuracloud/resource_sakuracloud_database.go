// Copyright 2016-2021 terraform-provider-sakuracloud authors
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
	"github.com/sacloud/libsacloud/v2/helper/power"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
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
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(3, 20),
				Description:  descf("The name of default user on the database. %s", descLength(3, 20)),
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
			"network_interface": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch_id": schemaResourceSwitchID(resourceName),
						"ip_address": {
							Type:        schema.TypeString,
							ForceNew:    true,
							Required:    true,
							Description: descf("The IP address to assign to the %s", resourceName),
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
						"source_ranges": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Description: descf(
								"The range of source IP addresses that allow to access to the %s via network",
								resourceName,
							),
						},
					},
				},
			},
			"backup": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"weekdays": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
							Description: descf(
								"A list of weekdays to backed up. The values in the list must be in [%s]",
								types.BackupWeekdayStrings,
							),
						},
						"time": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateBackupTime(),
							Description:  "The time to take backup. This must be formatted with `HH:mm`",
						},
					},
				},
			},
			"parameters": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The map for setting RDMBS-specific parameters. Valid keys can be found with the `usacloud database list-parameters` command",
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
	if db != nil {
		d.SetId(db.ID.String())
	}
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Database is failed: %s", err)
	}

	// HACK データベースアプライアンスの電源投入後すぐに他の操作(Updateなど)を行うと202(Accepted)が返ってくるものの無視される。
	// この挙動はテストなどで問題となる。このためここで少しsleepすることで対応する。
	time.Sleep(client.databaseWaitAfterCreateDuration)

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
	return setDatabaseResourceData(ctx, d, client, data, zone)
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
		return fmt.Errorf("updating SakuraCloud Database[%s] is failed: %s", d.Id(), err)
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
		return fmt.Errorf("deleting SakuraCloud Database[%s] is failed: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func setDatabaseResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Database, zone string) error {
	if data.Availability.IsFailed() {
		d.SetId("")
		return fmt.Errorf("got unexpected state: Database[%d].Availability is failed", data.ID)
	}

	d.Set("database_type", flattenDatabaseType(data)) // nolint
	if data.ReplicationSetting != nil {
		d.Set("replica_user", data.CommonSetting.ReplicaUser)         // nolint
		d.Set("replica_password", data.CommonSetting.ReplicaPassword) // nolint
	}
	if err := d.Set("backup", flattenDatabaseBackupSetting(data)); err != nil {
		return err
	}
	if err := d.Set("tags", flattenDatabaseTags(data)); err != nil {
		return err
	}
	d.Set("name", data.Name)                              // nolint
	d.Set("username", data.CommonSetting.DefaultUser)     // nolint
	d.Set("password", data.CommonSetting.UserPassword)    // nolint
	d.Set("plan", types.DatabasePlanNameMap[data.PlanID]) // nolint
	if err := d.Set("network_interface", flattenDatabaseNetworkInterface(data)); err != nil {
		return err
	}
	d.Set("icon_id", data.IconID.String()) // nolint
	d.Set("description", data.Description) // nolint
	d.Set("zone", getZone(d, client))      // nolint

	parameters, err := sacloud.NewDatabaseOp(client).GetParameter(ctx, zone, data.ID)
	if err != nil {
		return err
	}
	parameterSettings := convertDatabaseParametersToStringKeyValues(parameters)
	if err := d.Set("parameters", parameterSettings); err != nil {
		return err
	}

	return nil
}

func convertDatabaseParametersToStringKeyValues(parameter *sacloud.DatabaseParameter) map[string]interface{} {
	stringMap := make(map[string]interface{})
	// convert to string
	for k, v := range parameter.Settings {
		switch v := v.(type) {
		case fmt.Stringer:
			stringMap[k] = v.String()
		case string:
			stringMap[k] = v
		default:
			stringMap[k] = fmt.Sprintf("%v", v)
		}
	}

	dest := make(map[string]interface{})
	for k, v := range stringMap {
		for _, meta := range parameter.MetaInfo {
			if k == meta.Name {
				dest[meta.Label] = v
			}
		}
	}

	return dest
}
