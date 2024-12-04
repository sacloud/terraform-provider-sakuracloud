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
							Type:        schema.TypeList,
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
													Type:     schema.TypeString,
													Required: true,
												},
												"port": {
													Type:     schema.TypeInt,
													Required: true,
												},
												"headers": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"value": {
																Type:     schema.TypeString,
																Optional: true,
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
		Components:     *expandApprunApplicationComponents(d),
	}
	result, err := appOp.Create(ctx, &params)
	if err != nil {
		return diag.FromErr(err)
	}

	// 内部的にVersions/Traffics APIを利用してトラフィック分散の状態も変更する
	versions, err := getVersions(ctx, d, meta, *result.Id)
	if err != nil {
		return diag.FromErr(err)
	}

	trafficOp := apprun.NewTrafficOp(client.apprunClient)
	traffics, err := expandApprunApplicationTraffics(d, versions)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = trafficOp.Update(ctx, *result.Id, traffics)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*result.Id)
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
			if e.Detail.Code != nil && *e.Detail.Code == 404 {
				d.SetId("")
				return nil
			}
		}
		return diag.Errorf("could not read SakuraCloud Apprun Application[%s]: %s", d.Id(), err)
	}
	d.SetId(*application.Id)

	versions, err := getVersions(ctx, d, meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	trafficOp := apprun.NewTrafficOp(client.apprunClient)
	traffics, err := trafficOp.List(ctx, d.Id())
	if err != nil {
		return diag.Errorf("could not read SakuraCloud Apprun Application Traffics[%s]: %s", d.Id(), err)
	}

	return setApprunApplicationResourceData(d, application, traffics.Data, versions)
}

// NOTE: all_traffic_availableについては未対応
func resourceSakuraCloudApprunApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// Applicationの状態を変更
	appOp := apprun.NewApplicationOp(client.apprunClient)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	application, err := appOp.Read(ctx, d.Id())
	if err != nil {
		return diag.Errorf("could not read SakuraCloud Apprun Application[%s]: %s", d.Id(), err)
	}

	patchedName := d.Get("name").(string)
	patchedTimeoutSeconds := d.Get("timeout_seconds").(int)
	patchedPort := d.Get("port").(int)
	patchedMinScale := d.Get("min_scale").(int)
	patchedMaxScale := d.Get("max_scale").(int)
	params := v1.PatchApplicationBody{
		Name:           &patchedName,
		TimeoutSeconds: &patchedTimeoutSeconds,
		Port:           &patchedPort,
		MinScale:       &patchedMinScale,
		MaxScale:       &patchedMaxScale,
		Components:     expandApprunApplicationComponentsForUpdate(d),
	}
	result, err := appOp.Update(ctx, *application.Id, &params)
	if err != nil {
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

	d.SetId(*result.Id)
	return resourceSakuraCloudApprunApplicationRead(ctx, d, meta)
}

func resourceSakuraCloudApprunApplicationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	appOp := apprun.NewApplicationOp(client.apprunClient)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	application, err := appOp.Read(ctx, d.Id())
	if err != nil {
		return diag.Errorf("could not read SakuraCloud Apprun Application[%s]: %s", d.Id(), err)
	}

	if err := appOp.Delete(ctx, *application.Id); err != nil {
		return diag.Errorf("deleting SakuraCloud Apprun Application[%s] is failed: %s", *application.Id, err)
	}
	return nil
}

func setApprunApplicationResourceData(d *schema.ResourceData, application *v1.Application, traffics *[]v1.Traffic, versions *[]v1.Version) diag.Diagnostics {
	d.Set("name", application.Name)
	d.Set("timeout_seconds", application.TimeoutSeconds)
	d.Set("port", application.Port)
	d.Set("min_scale", application.MinScale)
	d.Set("max_scale", application.MaxScale)
	d.Set("components", flattenApprunApplicationComponents(d, application, true))
	d.Set("traffics", flattenApprunApplicationTraffics(traffics, versions))

	return nil
}

func validateApprunApplicationMaxCPU() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice(apprun.ApplicationMaxCPUs, false))
}

func validateApprunApplicationMaxMemory() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice(apprun.ApplicationMaxMemories, false))
}

func getVersions(ctx context.Context, d *schema.ResourceData, meta interface{}, applicationId string) (*[]v1.Version, error) {
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
		if len(*vs.Data) == 0 {
			break
		}

		versions = append(versions, *vs.Data...)
		pageNum++
	}

	return &versions, nil
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
