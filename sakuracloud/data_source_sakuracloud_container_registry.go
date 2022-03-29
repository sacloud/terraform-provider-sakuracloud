// Copyright 2016-2022 terraform-provider-sakuracloud authors
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
)

func dataSourceSakuraCloudContainerRegistry() *schema.Resource {
	resourceName := "ContainerRegistry"
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudContainerRegistryRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"access_level": {
				Type:     schema.TypeString,
				Computed: true,
				Description: descf(
					"The level of access that allow to users. This will be one of [%s]",
					types.ContainerRegistryAccessLevelStrings,
				),
			},
			"virtual_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The alias for accessing the container registry",
			},
			"subdomain_label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The label at the lowest of the FQDN used when be accessed from users",
			},
			"fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The FQDN for accessing the container registry. FQDN is built from `subdomain_label` + `.sakuracr.jp`",
			},
			"user": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user name used to authenticate remote access",
						},
						"permission": {
							Type:     schema.TypeString,
							Computed: true,
							Description: descf(
								"The level of access that allow to the user. This will be one of [%s]",
								types.ContainerRegistryPermissionStrings,
							),
						},
					},
				},
			},
			"icon_id":     schemaDataSourceIconID(resourceName),
			"description": schemaDataSourceDescription(resourceName),
			"tags":        schemaDataSourceTags(resourceName),
		},
	}
}

func dataSourceSakuraCloudContainerRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	searcher := iaas.NewContainerRegistryOp(client)

	findCondition := &iaas.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud ContainerRegistry: %s", err)
	}
	if res == nil || res.Count == 0 {
		return filterNoResultErr()
	}

	targets := res.ContainerRegistries
	if len(targets) == 0 {
		return filterNoResultErr()
	}

	d.SetId(targets[0].ID.String())
	return setContainerRegistryResourceData(ctx, d, client, targets[0], false)
}
