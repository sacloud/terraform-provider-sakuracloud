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

func dataSourceSakuraCloudLoadBalancer() *schema.Resource {
	resourceName := "LoadBalancer"
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudLoadBalancerRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"plan":         schemaDataSourcePlan(resourceName, []string{"standard", "highspec"}),
			"network_interface": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch_id":    schemaDataSourceSwitchID(resourceName),
						"ip_addresses": schemaDataSourceIPAddresses(resourceName),
						"netmask":      schemaDataSourceNetMask(resourceName),
						"gateway":      schemaDataSourceGateway(resourceName),
						"vrid": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Virtual Router Identifier",
						},
					},
				},
			},
			"vip": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The virtual IP address",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The target port number for load-balancing",
						},
						"delay_loop": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The interval in seconds between checks",
						},
						"sorry_server": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address of the SorryServer. This will be used when all servers under this VIP are down",
						},
						"description": schemaDataSourceDescription("VIP"),
						"server": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address of the destination server",
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
										Description: desc.Sprintf(
											"The protocol used for health checks. This will be one of [%s]",
											types.LoadBalancerHealthCheckProtocolStrings,
										),
									},
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The path used when checking by HTTP/HTTPS",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The response code to expect when checking by HTTP/HTTPS",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "The flag to enable as destination of load balancing",
									},
								},
							},
						},
					},
				},
			},
			"icon_id":     schemaDataSourceIconID(resourceName),
			"description": schemaDataSourceDescription(resourceName),
			"tags":        schemaDataSourceTags(resourceName),
			"zone":        schemaDataSourceZone(resourceName),
		},
	}
}

func dataSourceSakuraCloudLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	searcher := iaas.NewLoadBalancerOp(client)

	findCondition := &iaas.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, zone, findCondition)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud LoadBalancer resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.LoadBalancers) == 0 {
		return filterNoResultErr()
	}

	targets := res.LoadBalancers
	d.SetId(targets[0].ID.String())
	return setLoadBalancerResourceData(ctx, d, client, targets[0])
}
