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
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func dataSourceSakuraCloudDisk() *schema.Resource {
	resourceName := "Disk"

	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudDiskRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"plan":         schemaDataSourcePlan(resourceName, types.DiskPlanStrings),
			"connector": {
				Type:     schema.TypeString,
				Computed: true,
				Description: desc.Sprintf(
					"The name of the disk connector. This will be one of [%s]",
					types.DiskConnectionStrings,
				),
			},
			"source_archive_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the source archive",
			},
			"source_disk_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the source disk",
			},
			"size":        schemaDataSourceSize(resourceName),
			"server_id":   schemaDataSourceServerID(resourceName),
			"icon_id":     schemaDataSourceIconID(resourceName),
			"description": schemaDataSourceDescription(resourceName),
			"tags":        schemaDataSourceTags(resourceName),
			"zone":        schemaDataSourceZone(resourceName),
		},
	}
}

func dataSourceSakuraCloudDiskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	searcher := iaas.NewDiskOp(client)

	findCondition := &iaas.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, zone, findCondition)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud Disk resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.Disks) == 0 {
		return filterNoResultErr()
	}

	targets := res.Disks
	d.SetId(targets[0].ID.String())
	return setDiskResourceData(ctx, d, client, targets[0])
}
