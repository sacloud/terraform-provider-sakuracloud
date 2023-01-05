// Copyright 2016-2023 terraform-provider-sakuracloud authors
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/iaas-api-go"
)

func expandAutoBackupCreateRequest(d *schema.ResourceData) *iaas.AutoBackupCreateRequest {
	return &iaas.AutoBackupCreateRequest{
		Name:                    d.Get("name").(string),
		Description:             d.Get("description").(string),
		Tags:                    expandTags(d),
		DiskID:                  expandSakuraCloudID(d, "disk_id"),
		MaximumNumberOfArchives: d.Get("max_backup_num").(int),
		BackupSpanWeekdays:      expandBackupWeekdays(d, "weekdays"),
		IconID:                  expandSakuraCloudID(d, "icon_id"),
	}
}

func expandAutoBackupUpdateRequest(d *schema.ResourceData, autoBackup *iaas.AutoBackup) *iaas.AutoBackupUpdateRequest {
	return &iaas.AutoBackupUpdateRequest{
		Name:                    d.Get("name").(string),
		Description:             d.Get("description").(string),
		Tags:                    expandTags(d),
		MaximumNumberOfArchives: d.Get("max_backup_num").(int),
		BackupSpanWeekdays:      expandBackupWeekdays(d, "weekdays"),
		IconID:                  expandSakuraCloudID(d, "icon_id"),
		SettingsHash:            autoBackup.SettingsHash,
	}
}
