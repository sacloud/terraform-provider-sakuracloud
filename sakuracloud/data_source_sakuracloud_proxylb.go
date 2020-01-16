// Copyright 2016-2020 terraform-provider-sakuracloud authors
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

func dataSourceSakuraCloudProxyLB() *schema.Resource {
	resourceName := "ProxyLB"
	return &schema.Resource{
		Read: dataSourceSakuraCloudProxyLBRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"plan":         schemaDataSourceIntPlan(resourceName, types.ProxyLBPlanValues),
			"vip_failover": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag to enable VIP fail-over",
			},
			"sticky_session": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag to enable sticky session",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The timeout duration in seconds",
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
				Description: descf(
					"The name of region that the proxy LB is in. This will be one of [%s]",
					types.ProxyLBRegionStrings,
				),
			},
			"bind_port": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_mode": {
							Type:     schema.TypeString,
							Computed: true,
							Description: descf(
								"The proxy mode. This will be one of [%s]",
								types.ProxyLBProxyModeStrings,
							),
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of listening port",
						},
						"redirect_to_https": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The flag to enable redirection from http to https. This flag is used only when `proxy_mode` is `http`",
						},
						"support_http2": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The flag to enable HTTP/2. This flag is used only when `proxy_mode` is `https`",
						},
						"response_header": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"header": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: descf("The field name of HTTP header added to response by the %s", resourceName),
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: descf("The field value of HTTP header added to response by the %s", resourceName),
									},
								},
							},
						},
					},
				},
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
								types.ProxyLBProtocolStrings,
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
							Description: "The value of host header send when checking by HTTP",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The path used when checking by HTTP",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The port number used when checking by TCP",
						},
					},
				},
			},
			"sorry_server": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address of the SorryServer. This will be used when all servers are down",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The port number of the SorryServer. This will be used when all servers are down",
						},
					},
				},
			},
			"certificate": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_cert": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The certificate for a server",
						},
						"intermediate_cert": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The intermediate certificate for a server",
						},
						"private_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The private key for a server",
						},
						"additional_certificate": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_cert": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The certificate for a server",
									},
									"intermediate_cert": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The intermediate certificate for a server",
									},
									"private_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The private key for a server",
									},
								},
							},
						},
					},
				},
			},
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
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The port number of the destination server",
						},
						"group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of load balancing group. This is used when using rule-based load balancing",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The flag to enable as destination of load balancing",
						},
					},
				},
			},
			"rule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value of HTTP host header that is used as condition of rule-based balancing",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The request path that is used as condition of rule-based balancing",
						},
						"group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of load balancing group. When proxyLB received request which matched to `host` and `path`, proxyLB forwards the request to servers that having same group name",
						},
					},
				},
			},
			"icon_id":     schemaDataSourceIconID(resourceName),
			"description": schemaDataSourceDescription(resourceName),
			"tags":        schemaDataSourceTags(resourceName),
			"fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("The FQDN for accessing to the %s. This is typically used as value of CNAME record", resourceName),
			},
			"vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("The virtual IP address assigned to the %s", resourceName),
			},
			"proxy_networks": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: descf("A list of CIDR block used by the %s to access the server", resourceName),
			},
		},
	}
}

func dataSourceSakuraCloudProxyLBRead(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	searcher := sacloud.NewProxyLBOp(client)

	findCondition := &sacloud.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud ProxyLB resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.ProxyLBs) == 0 {
		return filterNoResultErr()
	}

	targets := res.ProxyLBs
	d.SetId(targets[0].ID.String())
	return setProxyLBResourceData(ctx, d, client, targets[0])
}
