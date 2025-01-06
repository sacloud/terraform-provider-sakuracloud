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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/iaas-api-go"
)

func dataSourceSakuraCloudIcon() *schema.Resource {
	resourceName := "Icon"
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudIconRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"tags":         schemaDataSourceTags(resourceName),
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for getting the icon's raw data",
			},
		},
	}
}

func dataSourceSakuraCloudIconRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	searcher := iaas.NewIconOp(client)

	findCondition := &iaas.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud Icon resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.Icons) == 0 {
		return filterNoResultErr()
	}

	targets := res.Icons
	icon := targets[0]

	d.SetId(icon.ID.String())
	return setIconResourceData(ctx, d, client, icon)
}
