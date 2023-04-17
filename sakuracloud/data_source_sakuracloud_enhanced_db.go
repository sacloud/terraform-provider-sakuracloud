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
	"github.com/sacloud/iaas-service-go/enhanceddb/builder"
)

func dataSourceSakuraCloudEnhancedDB() *schema.Resource {
	resourceName := "EnhancedDB"
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudEnhancedDBRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"database_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of database",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region name",
			},
			"database_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of database",
			},
			"hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of database host. This will be built from `database_name` + `tidb-is1.db.sakurausercontent.com`",
			},
			"max_connections": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The value of max connections setting",
			},
			"icon_id":     schemaDataSourceIconID(resourceName),
			"description": schemaDataSourceDescription(resourceName),
			"tags":        schemaDataSourceTags(resourceName),
		},
	}
}

func dataSourceSakuraCloudEnhancedDBRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	edbOp := iaas.NewEnhancedDBOp(client)

	findCondition := &iaas.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := edbOp.Find(ctx, findCondition)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud EnhancedDB: %s", err)
	}
	if res == nil || res.Count == 0 {
		return filterNoResultErr()
	}

	targets := res.EnhancedDBs
	if len(targets) == 0 {
		return filterNoResultErr()
	}

	d.SetId(targets[0].ID.String())

	edb, err := builder.Read(ctx, edbOp, targets[0].ID)
	if err != nil {
		return diag.Errorf("could not read EnhancedDB: %s", err)
	}
	return setEnhancedDBResourceData(ctx, d, client, edb, false)
}
