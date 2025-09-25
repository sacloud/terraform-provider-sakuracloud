// Copyright 2016-2025 The sacloud/terraform-provider-sakuracloud Authors
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
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/apprun-api-go"
	v1 "github.com/sacloud/apprun-api-go/apis/v1"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudApprunApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSakuraCloudApprunApplicationCreate,
		UpdateContext: resourceSakuraCloudApprunApplicationUpdate,
		ReadContext:   resourceSakuraCloudApprunApplicationRead,
		DeleteContext: resourceSakuraCloudApprunApplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of application",
			},
			"timeout_seconds": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The time limit between accessing the application's public URL, starting the instance, and receiving a response",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The port number where the application listens for requests",
			},
			"min_scale": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The minimum number of scales for the entire application",
			},
			"max_scale": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The maximum number of scales for the entire application",
			},
			"components": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The application component information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The component name",
							ForceNew:    true,
						},
						"max_cpu": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateApprunApplicationMaxCPU(),
							Description: desc.Sprintf(
								"The maximum number of CPUs for a component. The values in the list must be in [%s]",
								apprun.ApplicationMaxCPUs,
							),
						},
						"max_memory": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateApprunApplicationMaxMemory(),
							Description: desc.Sprintf(
								"The maximum memory of component. The values in the list must be in [%s]",
								apprun.ApplicationMaxMemories,
							),
						},
						"deploy_source": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "The sources that make up the component",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"container_registry": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"image": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The container image name",
												},
												"server": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The container registry server name",
												},
												"username": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The container registry credentials",
												},
												"password": {
													Type:        schema.TypeString,
													Optional:    true,
													Sensitive:   true,
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
							Optional:    true,
							Description: "The environment variables passed to components",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The environment variable name",
									},
									"value": {
										Type:        schema.TypeString,
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
									"value": {
										Type: schema.TypeString,
									},
								},
							}),
						},
						"probe": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The component probe settings",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"http_get": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The path to access HTTP server to check probes",
												},
												"port": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "The port number for accessing HTTP server and checking probes",
												},
												"headers": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The header field name",
															},
															"value": {
																Type:        schema.TypeString,
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
			"traffics": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The application traffic",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_index": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The application version index",
						},
						"percent": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The percentage of traffic dispersion",
						},
					},
				},
			},
			"packet_filter": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "The packet filter for the application",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"settings": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    5,
							Description: "The list of packet filter rule",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"from_ip": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The source IP address of the rule",
									},
									"from_ip_prefix_length": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The prefix length (CIDR notation) of the from_ip address, indicating the network size",
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

func resourceSakuraCloudApprunApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	createUserIfNotExist(ctx, d, meta)

	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	appOp := apprun.NewApplicationOp(client.apprunClient)
	params := v1.PostApplicationBody{
		Name:           d.Get("name").(string),
		TimeoutSeconds: d.Get("timeout_seconds").(int),
		Port:           d.Get("port").(int),
		MinScale:       d.Get("min_scale").(int),
		MaxScale:       d.Get("max_scale").(int),
		Components:     expandApprunApplicationComponents(d),
	}
	result, err := appOp.Create(ctx, &params)
	if err != nil {
		return diag.FromErr(err)
	}

	pfOp := apprun.NewPacketFilterOp(client.apprunClient)
	if _, err := pfOp.Update(ctx, result.Id, expandApprunPacketFilter(d)); err != nil {
		return diag.FromErr(err)
	}

	// 内部的にVersions/Traffics APIを利用してトラフィック分散の状態も変更する
	versions, err := getVersions(ctx, d, meta, result.Id)
	if err != nil {
		return diag.FromErr(err)
	}

	trafficOp := apprun.NewTrafficOp(client.apprunClient)
	traffics, err := expandApprunApplicationTraffics(d, versions)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = trafficOp.Update(ctx, result.Id, traffics)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.Id)
	return resourceSakuraCloudApprunApplicationRead(ctx, d, meta)
}

func resourceSakuraCloudApprunApplicationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	createUserIfNotExist(ctx, d, meta)

	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	appOp := apprun.NewApplicationOp(client.apprunClient)

	application, err := appOp.Read(ctx, d.Id())
	if err != nil {
		if e, ok := err.(*v1.ModelDefaultError); ok {
			if e.Detail.Code == 404 {
				d.SetId("")
				return nil
			}
		}
		return diag.Errorf("could not read SakuraCloud Apprun Application[%s]: %s", d.Id(), err)
	}
	d.SetId(application.Id)

	versions, err := getVersions(ctx, d, meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	trafficOp := apprun.NewTrafficOp(client.apprunClient)
	traffics, err := trafficOp.List(ctx, d.Id())
	if err != nil {
		return diag.Errorf("could not read SakuraCloud Apprun Application Traffics[%s]: %s", d.Id(), err)
	}

	pfOp := apprun.NewPacketFilterOp(client.apprunClient)
	pf, err := pfOp.Read(ctx, application.Id)
	if err != nil {
		return diag.FromErr(err)
	}

	return setApprunApplicationResourceData(d, application, traffics.Data, versions, pf)
}

// NOTE: all_traffic_availableについては未対応
func resourceSakuraCloudApprunApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// Applicationの状態を変更
	appOp := apprun.NewApplicationOp(client.apprunClient)
	application, err := appOp.Read(ctx, d.Id())
	if err != nil {
		return diag.Errorf("could not read SakuraCloud Apprun Application[%s]: %s", d.Id(), err)
	}

	patchedTimeoutSeconds := d.Get("timeout_seconds").(int)
	patchedPort := d.Get("port").(int)
	patchedMinScale := d.Get("min_scale").(int)
	patchedMaxScale := d.Get("max_scale").(int)
	params := v1.PatchApplicationBody{
		TimeoutSeconds: &patchedTimeoutSeconds,
		Port:           &patchedPort,
		MinScale:       &patchedMinScale,
		MaxScale:       &patchedMaxScale,
		Components:     expandApprunApplicationComponentsForUpdate(d),
	}
	result, err := appOp.Update(ctx, application.Id, &params)
	if err != nil {
		return diag.FromErr(err)
	}

	pfOp := apprun.NewPacketFilterOp(client.apprunClient)
	if _, err := pfOp.Update(ctx, result.Id, expandApprunPacketFilter(d)); err != nil {
		return diag.FromErr(err)
	}

	// 内部的にVersions/Traffics APIを利用してトラフィック分散の状態も変更する
	versions, err := getVersions(ctx, d, meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	trafficOp := apprun.NewTrafficOp(client.apprunClient)
	traffics, err := expandApprunApplicationTraffics(d, versions)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = trafficOp.Update(ctx, d.Id(), traffics)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.Id)
	return resourceSakuraCloudApprunApplicationRead(ctx, d, meta)
}

func resourceSakuraCloudApprunApplicationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	appOp := apprun.NewApplicationOp(client.apprunClient)
	application, err := appOp.Read(ctx, d.Id())
	if err != nil {
		return diag.Errorf("could not read SakuraCloud Apprun Application[%s]: %s", d.Id(), err)
	}

	if err := appOp.Delete(ctx, application.Id); err != nil {
		return diag.Errorf("deleting SakuraCloud Apprun Application[%s] is failed: %s", application.Id, err)
	}
	return nil
}

func setApprunApplicationResourceData(d *schema.ResourceData, application *v1.Application, traffics []v1.Traffic, versions []v1.Version, pf *v1.HandlerGetPacketFilter) diag.Diagnostics {
	d.Set("name", application.Name)                                               //nolint:errcheck,gosec
	d.Set("timeout_seconds", application.TimeoutSeconds)                          //nolint:errcheck,gosec
	d.Set("port", application.Port)                                               //nolint:errcheck,gosec
	d.Set("min_scale", application.MinScale)                                      //nolint:errcheck,gosec
	d.Set("max_scale", application.MaxScale)                                      //nolint:errcheck,gosec
	d.Set("components", flattenApprunApplicationComponents(d, application, true)) //nolint:errcheck,gosec
	d.Set("traffics", flattenApprunApplicationTraffics(traffics, versions))       //nolint:errcheck,gosec
	d.Set("packet_filter", flattenApprunPacketFilter(pf))                         //nolint:errcheck,gosec
	d.Set("status", application.Status)                                           //nolint:errcheck,gosec
	d.Set("public_url", application.PublicUrl)                                    //nolint:errcheck,gosec

	return nil
}

func validateApprunApplicationMaxCPU() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice(apprun.ApplicationMaxCPUs, false))
}

func validateApprunApplicationMaxMemory() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice(apprun.ApplicationMaxMemories, false))
}

func getVersions(ctx context.Context, d *schema.ResourceData, meta interface{}, applicationId string) ([]v1.Version, error) {
	var versions []v1.Version

	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return nil, err
	}

	versionOp := apprun.NewVersionOp(client.apprunClient)

	pageNum := 1
	pageSize := 100
	for {
		vs, err := versionOp.List(ctx, applicationId, &v1.ListApplicationVersionsParams{
			PageNum:  &pageNum,
			PageSize: &pageSize,
		})
		if err != nil {
			return nil, err
		}
		if len(vs.Data) == 0 {
			break
		}

		versions = append(versions, vs.Data...)
		pageNum++
	}

	return versions, nil
}

// NOTE: AppRunは初回利用時に一度のみユーザーの作成を必要とする。
// SakuraCloud Providerでは明示的にユーザーの作成を行わず、CURD操作の開始時に暗黙的にユーザーの存在確認と作成を行う。
// ref. https://manual.sakura.ad.jp/sakura-apprun-api/spec.html#tag/%E3%83%A6%E3%83%BC%E3%82%B6%E3%83%BC
func createUserIfNotExist(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	userOp := apprun.NewUserOp(client.apprunClient)
	res, err := userOp.Read(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if res.StatusCode == http.StatusNotFound {
		_, err := userOp.Create(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}
