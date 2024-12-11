// Copyright 2016-2023 The sacloud/terraform-provider-sakuracloud Authors
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
	"github.com/sacloud/apprun-api-go"
	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

func dataSourceSakuraCloudApprunApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudApprunApplicationRead,

		Schema: map[string]*schema.Schema{
			// input/condition
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of application",
			},

			// computed fields
			"timeout_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The time limit between accessing the application's public URL, starting the instance, and receiving a response",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port number where the application listens for requests",
			},
			"min_scale": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The minimum number of scales for the entire application",
			},
			"max_scale": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of scales for the entire application",
			},
			"components": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The application component information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The component name",
						},
						"max_cpu": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The maximum number of CPUs for a component",
						},
						"max_memory": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The maximum memory of component",
						},
						"deploy_source": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The sources that make up the component",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"container_registry": {
										Type:     schema.TypeList,
										Computed: true,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"image": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The container image name",
												},
												"server": {
													Type:        schema.TypeString,
													Computed:    true,
													Optional:    true,
													Description: "The container registry server name",
												},
												"username": {
													Type:        schema.TypeString,
													Computed:    true,
													Optional:    true,
													Description: "The container registry credentials",
												},
											},
										},
									},
								},
							},
						},
						"env": {
							Type:        schema.TypeSet,
							Computed:    true,
							Optional:    true,
							Description: "The environment variables passed to components",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Optional:    true,
										Description: "The environment variable name",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Optional:    true,
										Sensitive:   true,
										Description: "environment variable value",
									},
								},
							},
							Set: schema.HashResource(&schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type: schema.TypeString,
									},
								},
							}),
						},
						"probe": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "The component probe settings",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"http_get": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The path to access HTTP server to check probes",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The port number for accessing HTTP server and checking probes",
												},
												"headers": {
													Type:     schema.TypeList,
													Computed: true,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Optional:    true,
																Description: "The header field name",
															},
															"value": {
																Type:        schema.TypeString,
																Computed:    true,
																Optional:    true,
																Description: "The header field value",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The application status",
			},
			"public_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public URL",
			},
		},
	}
}

func dataSourceSakuraCloudApprunApplicationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	if name == "" {
		return diag.Errorf("name required")
	}

	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	appOp := apprun.NewApplicationOp(client.apprunClient)

	apps, err := appOp.List(ctx, &v1.ListApplicationsParams{})
	if err != nil {
		return diag.Errorf("could not find SakuraCloud AppRun resource: %s", err)
	}
	if apps == nil || len(*apps.Data) == 0 {
		return filterNoResultErr()
	}

	var data *v1.Application
	for _, d := range *apps.Data {
		if *d.Name == name {
			a, err := appOp.Read(ctx, *d.Id)
			if err != nil {
				return diag.FromErr(err)
			}
			data = a
			break
		}
	}
	if data == nil {
		return filterNoResultErr()
	}

	d.SetId(*data.Id)
	d.Set("name", *data.Name)
	d.Set("timeout_seconds", *data.TimeoutSeconds)
	d.Set("port", *data.Port)
	d.Set("min_scale", *data.MinScale)
	d.Set("max_scale", *data.MaxScale)
	d.Set("components", flattenApprunApplicationComponents(d, data, false))
	d.Set("status", *data.Status)
	d.Set("public_url", *data.PublicUrl)

	return nil
}
