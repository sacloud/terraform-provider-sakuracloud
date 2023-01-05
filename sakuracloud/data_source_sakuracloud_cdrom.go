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
)

func dataSourceSakuraCloudCDROM() *schema.Resource {
	resourceName := "CD-ROM"
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudCDROMRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"size":         schemaDataSourceSize(resourceName),
			"icon_id":      schemaDataSourceIconID(resourceName),
			"description":  schemaDataSourceDescription(resourceName),
			"tags":         schemaDataSourceTags(resourceName),
			"zone":         schemaDataSourceZone(resourceName),
		},
	}
}

func dataSourceSakuraCloudCDROMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	searcher := iaas.NewCDROMOp(client)

	findCondition := &iaas.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, zone, findCondition)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud CDROM: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.CDROMs) == 0 {
		return filterNoResultErr()
	}

	targets := res.CDROMs
	target := targets[0]

	d.SetId(target.ID.String())
	d.Set("name", target.Name)               // nolint
	d.Set("size", target.GetSizeGB())        // nolint
	d.Set("icon_id", target.IconID.String()) // nolint
	d.Set("description", target.Description) // nolint
	d.Set("zone", getZone(d, client))        // nolint
	return diag.FromErr(d.Set("tags", flattenTags(target.Tags)))
}
