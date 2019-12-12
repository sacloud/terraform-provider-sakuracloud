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

	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/setup"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
	client, ctx, zone := getSacloudClient(d, meta)
	dbOp := sacloud.NewDatabaseOp(client)

	// validate master instance
	req, err := expandDatabaseReadReplicaCreateRequest(ctx, d, client, zone)
	if err != nil {
		return nil
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
	time.Sleep(client.databaseWaitAfterCreateDuration)

	d.SetId(database.ID.String())
	return setDatabaseReadReplicaResourceData(ctx, d, client, database)
}

func resourceSakuraCloudDatabaseReadReplicaRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	dbOp := sacloud.NewDatabaseOp(client)

	data, err := dbOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not find SakuraCloud Database ReadReplica[%s] : %s", d.Id(), err)
	}
	return setDatabaseReadReplicaResourceData(ctx, d, client, data)
}

func resourceSakuraCloudDatabaseReadReplicaUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	dbOp := sacloud.NewDatabaseOp(client)

	db, err := dbOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Database[%s]: %s", d.Id(), err)
	}

	db, err = dbOp.Update(ctx, zone, db.ID, expandDatabaseReadReplicaUpdateRequest(d, db))
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Database[%s] is failed: %s", d.Id(), err)
	}

	return setDatabaseReadReplicaResourceData(ctx, d, client, db)
}

func resourceSakuraCloudDatabaseReadReplicaDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	dbOp := sacloud.NewDatabaseOp(client)

	data, err := dbOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Database[%s]: %s", d.Id(), err)
	}

	// shutdown(force) if running
	if data.InstanceStatus.IsUp() {
		if err := shutdownDatabaseSync(ctx, dbOp, zone, data.ID, true); err != nil {
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

	if err := d.Set("tags", flattenDatabaseTags(data)); err != nil {
		return err
	}

	d.Set("master_id", data.ReplicationSetting.ApplianceID.String())
	d.Set("name", data.Name)
	d.Set("switch_id", data.SwitchID.String())
	d.Set("nw_mask_len", data.NetworkMaskLen)
	d.Set("default_route", data.DefaultRoute)
	d.Set("ipaddress1", data.IPAddresses[0])
	if err := d.Set("allow_networks", data.CommonSetting.SourceNetwork); err != nil {
		return err
	}
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	d.Set("zone", getZone(d, client))
	return nil
}
