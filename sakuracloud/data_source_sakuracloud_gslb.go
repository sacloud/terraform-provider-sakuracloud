// Copyright 2016-2021 terraform-provider-sakuracloud authors
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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func dataSourceSakuraCloudGSLB() *schema.Resource {
	resourceName := "GSLB"
	return &schema.Resource{
		Read: dataSourceSakuraCloudGSLBRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The FQDN for accessing to the GSLB. This is typically used as value of CNAME record",
			},
			"health_check": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
							Description: descf(
								"The protocol used for health checks. This will be one of [%s]",
								types.GSLBHealthCheckProtocolStrings,
							),
						},
						"delay_loop": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The interval in seconds between checks",
						},
						"host_header": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value of host header send when checking by HTTP/HTTPS",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The path used when checking by HTTP/HTTPS",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The response-code to expect when checking by HTTP/HTTPS",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The port number used when checking by TCP",
						},
					},
				},
			},
			"weighted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag to enable weighted load-balancing",
			},
			"sorry_server": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of the SorryServer. This will be used when all servers are down",
			},
			"icon_id":     schemaDataSourceIconID(resourceName),
			"description": schemaDataSourceDescription(resourceName),
			"tags":        schemaDataSourceTags(resourceName),
			"server": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address of the server",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The flag to enable as destination of load balancing",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The weight used when weighted load balancing is enabled",
						},
					},
				},
			},
		},
	}
}

func dataSourceSakuraCloudGSLBRead(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	searcher := sacloud.NewGSLBOp(client)

	findCondition := &sacloud.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud GSLB resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.GSLBs) == 0 {
		return filterNoResultErr()
	}

	targets := res.GSLBs
	d.SetId(targets[0].ID.String())
	return setGSLBResourceData(ctx, d, client, targets[0])
}
