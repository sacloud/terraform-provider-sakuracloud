// Copyright 2016-2025 terraform-provider-sakuracloud authors
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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	databaseBuilder "github.com/sacloud/iaas-service-go/database/builder"
)

func expandDatabaseBuilder(d *schema.ResourceData, client *APIClient) *databaseBuilder.Builder {
	dbType := d.Get("database_type").(string)
	dbName := types.RDBMSTypeFromString(dbType)

	nic := expandDatabaseNetworkInterface(d)
	// Note: 不正な値が含まれる場合にnilが返却される
	if nic == nil {
		return nil
	}

	replicaUser := d.Get("replica_user").(string)
	replicaPassword := d.Get("replica_password").(string)

	req := &databaseBuilder.Builder{
		PlanID:         types.DatabasePlanIDMap[d.Get("plan").(string)],
		SwitchID:       nic.switchID,
		IPAddresses:    nic.ipAddresses,
		NetworkMaskLen: nic.netmask,
		DefaultRoute:   nic.gateway,
		Conf: &iaas.DatabaseRemarkDBConfCommon{
			DatabaseName:    dbName.String(),
			DatabaseVersion: d.Get("database_version").(string),
			DefaultUser:     d.Get("username").(string),
			UserPassword:    d.Get("password").(string),
		},
		CommonSetting: &iaas.DatabaseSettingCommon{
			WebUI:           expandWebUiSetting(d),
			ServicePort:     nic.port,
			SourceNetwork:   nic.sourceRanges,
			DefaultUser:     d.Get("username").(string),
			UserPassword:    d.Get("password").(string),
			ReplicaUser:     replicaUser,
			ReplicaPassword: replicaPassword,
		},
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Tags:               expandTags(d),
		IconID:             expandSakuraCloudID(d, "icon_id"),
		Client:             databaseBuilder.NewAPIClient(client),
		BackupSetting:      expandDatabaseBackupSetting(d),
		Parameters:         d.Get("parameters").(map[string]interface{}),
		ReplicationSetting: &iaas.DatabaseReplicationSetting{},
	}
	if nic.virtualIpAddress != "" {
		req.VirtualIPAddress = &nic.virtualIpAddress
	}

	if replicaUser != "" && replicaPassword != "" {
		req.ReplicationSetting = &iaas.DatabaseReplicationSetting{
			Model:    types.DatabaseReplicationModels.MasterSlave,
			User:     replicaUser,
			Password: replicaPassword,
		}
	}
	return req
}

func expandWebUiSetting(d resourceValueGettable) types.WebUI {
	if isEnabled, ok := d.GetOk("webui"); ok {
		return types.ToWebUI(isEnabled.(bool))
	}
	//Note: types.WebUI は "" もしくは "false" で無効状態を表現する
	// デフォルトで無効であり、ここでは"false"として正規化する
	return "false"
}

func expandDatabaseReadReplicaBuilder(ctx context.Context, d *schema.ResourceData, client *APIClient, zone string) (*databaseBuilder.Builder, error) {
	dbOp := iaas.NewDatabaseOp(client)

	// validate master instance
	masterID := d.Get("master_id").(string)
	masterDB, err := dbOp.Read(ctx, zone, sakuraCloudID(masterID))
	if err != nil {
		return nil, fmt.Errorf("master database instance[%s] is not found", masterID)
	}
	if masterDB.ReplicationSetting.Model != types.DatabaseReplicationModels.MasterSlave {
		return nil, fmt.Errorf("master database instance[%s] is not configured as ReplicationMaster", masterID)
	}

	nic := expandDatabaseNetworkInterface(d)
	switchID := masterDB.SwitchID.String()
	if !nic.switchID.IsEmpty() {
		switchID = nic.switchID.String()
	}
	maskLen := masterDB.NetworkMaskLen
	if nic.netmask > 0 {
		maskLen = nic.netmask
	}
	defaultRoute := masterDB.DefaultRoute
	if nic.gateway != "" {
		defaultRoute = nic.gateway
	}

	return &databaseBuilder.Builder{
		Zone:           zone,
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Tags:           expandTags(d),
		IconID:         expandSakuraCloudID(d, "icon_id"),
		PlanID:         types.ID(masterDB.PlanID.Int64() + 1),
		SwitchID:       sakuraCloudID(switchID),
		IPAddresses:    nic.ipAddresses,
		NetworkMaskLen: maskLen,
		DefaultRoute:   defaultRoute,
		Conf: &iaas.DatabaseRemarkDBConfCommon{
			DatabaseName:     masterDB.Conf.DatabaseName,
			DatabaseVersion:  masterDB.Conf.DatabaseVersion,
			DatabaseRevision: masterDB.Conf.DatabaseRevision,
		},
		CommonSetting: &iaas.DatabaseSettingCommon{
			ServicePort:   masterDB.CommonSetting.ServicePort,
			SourceNetwork: nic.sourceRanges,
		},
		ReplicationSetting: &iaas.DatabaseReplicationSetting{
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

func flattenDatabaseType(db *iaas.Database) string {
	return strings.ToLower(db.Conf.DatabaseName)
}

func flattenDatabaseTags(db *iaas.Database) *schema.Set {
	var tags types.Tags
	for _, t := range db.Tags {
		if !(strings.HasPrefix(t, "@MariaDB-") || strings.HasPrefix(t, "@postgres-")) {
			tags = append(tags, t)
		}
	}
	return flattenTags(tags)
}

func expandDatabaseBackupSetting(d resourceValueGettable) *iaas.DatabaseSettingBackup {
	d = mapFromFirstElement(d, "backup")
	if d != nil {
		backupTime := d.Get("time").(string)
		backupWeekdays := expandBackupWeekdays(d, "weekdays")
		if backupTime != "" && len(backupWeekdays) > 0 {
			return &iaas.DatabaseSettingBackup{
				Time:      backupTime,
				DayOfWeek: backupWeekdays,
			}
		}
	}
	return nil
}

func flattenDatabaseBackupSetting(db *iaas.Database) []interface{} {
	if db.BackupSetting != nil {
		setting := map[string]interface{}{
			"time":     db.BackupSetting.Time,
			"weekdays": flattenBackupWeekdays(db.BackupSetting.DayOfWeek),
		}
		return []interface{}{setting}
	}
	return nil
}

type databaseNetworkInterface struct {
	switchID         types.ID
	ipAddresses      []string //Note: 要素を2つ以上持つ場合、冗長化オプションを有効化するものと見なす。
	virtualIpAddress string   //Note: 冗長化オプションを有効化する場合に必須。
	netmask          int
	gateway          string
	port             int
	sourceRanges     []string
}

func expandDatabaseNetworkInterface(d resourceValueGettable) *databaseNetworkInterface {
	d = mapFromFirstElement(d, "network_interface")
	if d == nil {
		return nil
	}

	res := &databaseNetworkInterface{
		switchID:     expandSakuraCloudID(d, "switch_id"),
		ipAddresses:  []string{stringOrDefault(d, "ip_address")},
		netmask:      intOrDefault(d, "netmask"),
		gateway:      stringOrDefault(d, "gateway"),
		port:         intOrDefault(d, "port"),
		sourceRanges: stringListOrDefault(d, "source_ranges"),
	}

	// 冗長化設定
	vip, hasVip := d.GetOk("vip")
	secondaryIpAddress, hasSecondaryIpAddress := d.GetOk("secondary_ip_address")
	if hasVip && hasSecondaryIpAddress {
		res.virtualIpAddress = vip.(string)
		res.ipAddresses = append(res.ipAddresses, secondaryIpAddress.(string))
	} else if hasVip || hasSecondaryIpAddress {
		// Note: 冗長化設定時にはVIPとセカンダリDBのIPアドレスの双方が必須。
		// nil を返却して設定値が不正であることを通知する。
		return nil
	}
	return res
}

func flattenDatabaseNetworkInterface(db *iaas.Database) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"switch_id":     db.SwitchID.String(),
			"netmask":       db.NetworkMaskLen,
			"source_ranges": db.CommonSetting.SourceNetwork,
			"port":          db.CommonSetting.ServicePort,
			"gateway":       db.DefaultRoute,
			"ip_address":    db.IPAddresses[0],
		},
	}
}

func flattenDatabaseReadReplicaNetworkInterface(db *iaas.Database) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"switch_id":     db.SwitchID.String(),
			"netmask":       db.NetworkMaskLen,
			"source_ranges": db.CommonSetting.SourceNetwork,
			"gateway":       db.DefaultRoute,
			"ip_address":    db.IPAddresses[0],
		},
	}
}
