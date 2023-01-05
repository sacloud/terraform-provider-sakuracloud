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

func dataSourceSakuraCloudBridge() *schema.Resource {
	resourceName := "Bridge"

	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudBridgeRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{excludeTags: true}),
			"name":         schemaDataSourceName(resourceName),
			"description":  schemaDataSourceDescription(resourceName),
			"zone":         schemaDataSourceZone(resourceName),
		},
	}
}

func dataSourceSakuraCloudBridgeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	searcher := iaas.NewBridgeOp(client)

	findCondition := &iaas.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, zone, findCondition)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud Bridge: %s", err)
	}
	if res == nil || res.Count == 0 {
		return filterNoResultErr()
	}

	targets := res.Bridges
	if len(targets) == 0 {
		return filterNoResultErr()
	}

	d.SetId(targets[0].ID.String())
	return setBridgeResourceData(ctx, d, client, targets[0])
}
