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
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/packages-go/size"
)

func flattenDiskPlan(data *iaas.Disk) string {
	plan, ok := types.DiskPlanNameMap[data.DiskPlanID]
	if !ok {
		return ""
	}
	return plan
}

func expandDiskPlan(d *schema.ResourceData) types.ID {
	return types.DiskPlanIDMap[d.Get("plan").(string)]
}

func expandDiskCreateRequest(d *schema.ResourceData) *iaas.DiskCreateRequest {
	return &iaas.DiskCreateRequest{
		DiskPlanID:      expandDiskPlan(d),
		Connection:      types.EDiskConnection(d.Get("connector").(string)),
		SourceDiskID:    expandSakuraCloudID(d, "source_disk_id"),
		SourceArchiveID: expandSakuraCloudID(d, "source_archive_id"),
		SizeMB:          d.Get("size").(int) * size.GiB,
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		Tags:            expandTags(d),
		IconID:          expandSakuraCloudID(d, "icon_id"),
	}
}

func expandDiskUpdateRequest(d *schema.ResourceData) *iaas.DiskUpdateRequest {
	return &iaas.DiskUpdateRequest{
		Connection:  types.EDiskConnection(d.Get("connector").(string)),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),
	}
}
