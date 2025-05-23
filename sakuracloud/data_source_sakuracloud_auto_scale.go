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
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func dataSourceSakuraCloudAutoScale() *schema.Resource {
	resourceName := "AutoScale"
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudAutoScaleRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"icon_id":      schemaDataSourceIconID(resourceName),
			"description":  schemaDataSourceDescription(resourceName),
			"tags":         schemaDataSourceTags(resourceName),
			"zones": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: desc.Sprintf("List of zone names where monitored resources are located"),
			},
			"config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration file for sacloud/autoscaler",
			},
			"api_key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the API key",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag to stop trigger",
			},
			"trigger_type": {
				Type:     schema.TypeString,
				Computed: true,
				Description: desc.Sprintf(
					"This must be one of [%s]",
					[]string{"cpu", "router", "scheudle"},
				),
			},
			"router_threshold_scaling": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"router_prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Router name prefix to be monitored",
						},
						"direction": {
							Type:     schema.TypeString,
							Computed: true,
							Description: desc.Sprintf(
								"This must be one of [%s]",
								[]string{"in", "out"},
							),
						},
						"mbps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Mbps",
						},
					},
				},
			},
			"cpu_threshold_scaling": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Server name prefix to be monitored",
						},
						"up": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Threshold for average CPU utilization to scale up/out",
						},
						"down": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Threshold for average CPU utilization to scale down/in",
						},
					},
				},
			},
			"schedule_scaling": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Computed: true,
							Description: desc.Sprintf(
								"This must be one of [%s]",
								[]string{"up", "down"},
							),
						},
						"hour": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Hour to be triggered",
						},
						"minute": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: desc.Sprintf(
								"Minute to be triggered. This must be one of [%s]",
								[]string{"0", "15", "30", "45"},
							),
						},
						"days_of_week": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
							Description: desc.Sprintf(
								"A list of weekdays to backed up. The values in the list must be in [%s]",
								types.DaysOfTheWeekStrings,
							),
						},
					},
				},
			},
		},
	}
}

func dataSourceSakuraCloudAutoScaleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	searcher := iaas.NewAutoScaleOp(client)

	findCondition := &iaas.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud AutoScale: %s", err)
	}
	if res == nil || res.Count == 0 {
		return filterNoResultErr()
	}

	targets := res.AutoScale
	if len(targets) == 0 {
		return filterNoResultErr()
	}

	d.SetId(targets[0].ID.String())
	return setAutoScaleResourceData(d, client, targets[0])
}
