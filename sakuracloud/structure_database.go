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
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	databaseBuilder "github.com/sacloud/libsacloud/v2/utils/builder/database"
)

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

func expandDatabaseBuilder(d *schema.ResourceData, client *APIClient) *databaseBuilder.Builder {
	var dbVersion *types.RDBMSVersion
	dbType := d.Get("database_type").(string)
	switch dbType {
	case "postgresql":
		dbVersion = types.RDBMSVersions[types.RDBMSTypesPostgreSQL]
	case "mariadb":
		dbVersion = types.RDBMSVersions[types.RDBMSTypesMariaDB]
	}

	replicaUser := d.Get("replica_user").(string)
	replicaPassword := d.Get("replica_password").(string)

	req := &databaseBuilder.Builder{
		PlanID:         databasePlanNameToID(d.Get("plan").(string)),
		SwitchID:       expandSakuraCloudID(d, "switch_id"),
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
			SourceNetwork:   expandStringList(d.Get("source_ranges").([]interface{})),
			DefaultUser:     d.Get("user_name").(string),
			UserPassword:    d.Get("user_password").(string),
			ReplicaUser:     replicaUser,
			ReplicaPassword: replicaPassword,
		},
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),
		Client:      databaseBuilder.NewAPIClient(client),
		// 後で設定する
		BackupSetting:      &sacloud.DatabaseSettingBackup{},
		ReplicationSetting: &sacloud.DatabaseReplicationSetting{},
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
	return req
}

func expandDatabaseReadReplicaBuilder(ctx context.Context, d *schema.ResourceData, client *APIClient, zone string) (*databaseBuilder.Builder, error) {
	dbOp := sacloud.NewDatabaseOp(client)

	// validate master instance
	masterID := d.Get("master_id").(string)
	masterDB, err := dbOp.Read(ctx, zone, sakuraCloudID(masterID))
	if err != nil {
		return nil, fmt.Errorf("master database instance[%s] is not found", masterID)
	}
	if masterDB.ReplicationSetting.Model != types.DatabaseReplicationModels.MasterSlave {
		return nil, fmt.Errorf("master database instance[%s] is not configured as ReplicationMaster", masterID)
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

	return &databaseBuilder.Builder{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Tags:           expandTags(d),
		IconID:         expandSakuraCloudID(d, "icon_id"),
		PlanID:         types.ID(masterDB.PlanID.Int64() + 1),
		SwitchID:       sakuraCloudID(switchID),
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
			SourceNetwork: expandStringList(d.Get("source_ranges").([]interface{})),
		},
		ReplicationSetting: &sacloud.DatabaseReplicationSetting{
			Model:       types.DatabaseReplicationModels.AsyncReplica,
			IPAddress:   masterDB.IPAddresses[0],
			Port:        masterDB.CommonSetting.ServicePort,
			User:        masterDB.ReplicationSetting.User,
			Password:    masterDB.ReplicationSetting.Password,
			ApplianceID: masterDB.ID,
		},
		Client: databaseBuilder.NewAPIClient(client),
	}, nil
}

func flattenDatabaseType(db *sacloud.Database) string {
	var databaseType string
	switch db.Conf.DatabaseName {
	case types.RDBMSVersions[types.RDBMSTypesPostgreSQL].Name:
		databaseType = "postgresql"
	case types.RDBMSVersions[types.RDBMSTypesMariaDB].Name:
		databaseType = "mariadb"
	}
	return databaseType
}

func flattenDatabaseTags(db *sacloud.Database) []interface{} {
	var tags []interface{}
	for _, t := range db.Tags {
		if !(strings.HasPrefix(t, "@MariaDB-") || strings.HasPrefix(t, "@postgres-")) {
			tags = append(tags, t)
		}
	}
	return tags
}
