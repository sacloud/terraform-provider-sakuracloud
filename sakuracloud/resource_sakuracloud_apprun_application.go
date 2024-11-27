package sakuracloud

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/sacloud/apprun-api-go"
	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

func resourceSakuraCloudApprunApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSakuraCloudApprunApplicationCreate,
		UpdateContext: resourceSakuraCloudApprunApplicationUpdate,
		ReadContext:   resourceSakuraCloudApprunApplicationRead,
		DeleteContext: resourceSakuraCloudApprunApplicationDelete,

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
							ValidateDiagFunc: *validateApprunApplicationMaxCPU(),
							Description:      "The maximum number of CPUs for a component",
						},
						"max_memory": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: *validateApprunApplicationMaxMemory(),
							Description:      "The maximum memory of component",
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
										Optional: true,
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
		},
	}
}

func resourceSakuraCloudApprunApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	d.SetId(*result.Id)
	return resourceSakuraCloudApprunApplicationRead(ctx, d, meta)
}

func resourceSakuraCloudApprunApplicationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return setApprunApplicationResourceData(d, application)
}

// NOTE: all_traffic_availableについては未対応
func resourceSakuraCloudApprunApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("deleting SakuraCloud Apprun Application[%s] is failed: %s", string(*application.Id), err)
	}
	return nil
}

func setApprunApplicationResourceData(d *schema.ResourceData, data *v1.Application) diag.Diagnostics {
	d.Set("name", data.Name)
	d.Set("timeout_seconds", data.TimeoutSeconds)
	d.Set("port", data.Port)
	d.Set("min_scale", data.MinScale)
	d.Set("max_scale", data.MaxScale)
	d.Set("components", flattenApprunApplicationComponents(d, data))

	return nil
}

func validateApprunApplicationMaxCPU() *schema.SchemaValidateDiagFunc {
	f := validation.ToDiagFunc(
		validation.StringInSlice([]string{
			(string)(v1.PostApplicationBodyComponentMaxCpuN01),
			(string)(v1.PostApplicationBodyComponentMaxCpuN02),
			(string)(v1.PostApplicationBodyComponentMaxCpuN03),
			(string)(v1.PostApplicationBodyComponentMaxCpuN04),
			(string)(v1.PostApplicationBodyComponentMaxCpuN05),
			(string)(v1.PostApplicationBodyComponentMaxCpuN06),
			(string)(v1.PostApplicationBodyComponentMaxCpuN07),
			(string)(v1.PostApplicationBodyComponentMaxCpuN08),
			(string)(v1.PostApplicationBodyComponentMaxCpuN09),
			(string)(v1.PostApplicationBodyComponentMaxCpuN1),
		}, false))
	return &f
}

func validateApprunApplicationMaxMemory() *schema.SchemaValidateDiagFunc {
	f := validation.ToDiagFunc(
		validation.StringInSlice([]string{
			(string)(v1.PostApplicationBodyComponentMaxMemoryN256Mi),
			(string)(v1.PostApplicationBodyComponentMaxMemoryN512Mi),
			(string)(v1.PostApplicationBodyComponentMaxMemoryN1Gi),
			(string)(v1.PostApplicationBodyComponentMaxMemoryN2Gi),
		}, false))
	return &f
}
