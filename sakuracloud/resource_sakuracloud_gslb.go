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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudGSLB() *schema.Resource {
	resourceName := "GSLB"

	return &schema.Resource{
		CreateContext: resourceSakuraCloudGSLBCreate,
		ReadContext:   resourceSakuraCloudGSLBRead,
		UpdateContext: resourceSakuraCloudGSLBUpdate,
		DeleteContext: resourceSakuraCloudGSLBDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The FQDN for accessing to the GSLB. This is typically used as value of CNAME record",
			},
			"health_check": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.GSLBHealthCheckProtocolStrings, false)),
							Description: desc.Sprintf(
								"The protocol used for health checks. This must be one of [%s]",
								types.GSLBHealthCheckProtocolStrings,
							),
						},
						"delay_loop": {
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(10, 60)),
							Default:          10,
							Description:      desc.Sprintf("The interval in seconds between checks. %s", desc.Range(10, 60)),
						},
						"host_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value of host header send when checking by HTTP/HTTPS",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The path used when checking by HTTP/HTTPS",
						},
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The response-code to expect when checking by HTTP/HTTPS",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The port number used when checking by TCP/HTTP/HTTPS",
						},
					},
				},
			},
			"weighted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The flag to enable weighted load-balancing",
			},
			"sorry_server": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP address of the SorryServer. This will be used when all servers are down",
			},
			"server": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 12,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "The IP address of the server",
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "The flag to enable as destination of load balancing",
						},
						"weight": {
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 10000)),
							Default:          1,
							Description: desc.Sprintf(
								"The weight used when weighted load balancing is enabled. %s",
								desc.Range(1, 10000),
							),
						},
					},
				},
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
		},
	}
}

func resourceSakuraCloudGSLBCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	gslbOp := iaas.NewGSLBOp(client)
	gslb, err := gslbOp.Create(ctx, expandGSLBCreateRequest(d))
	if err != nil {
		return diag.Errorf("creating SakuraCloud GSLB is failed: %s", err)
	}

	d.SetId(gslb.ID.String())
	return resourceSakuraCloudGSLBRead(ctx, d, meta)
}

func resourceSakuraCloudGSLBRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	gslbOp := iaas.NewGSLBOp(client)
	gslb, err := gslbOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud GSLB[%s]: %s", d.Id(), err)
	}

	return setGSLBResourceData(ctx, d, client, gslb)
}

func resourceSakuraCloudGSLBUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	gslbOp := iaas.NewGSLBOp(client)
	gslb, err := gslbOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud GSLB[%s]: %s", d.Id(), err)
	}

	_, err = gslbOp.Update(ctx, sakuraCloudID(d.Id()), expandGSLBUpdateRequest(d, gslb))
	if err != nil {
		return diag.Errorf("updating SakuraCloud GSLB[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudGSLBRead(ctx, d, meta)
}

func resourceSakuraCloudGSLBDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	gslbOp := iaas.NewGSLBOp(client)
	gslb, err := gslbOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud GSLB[%s]: %s", d.Id(), err)
	}
	if err := gslbOp.Delete(ctx, gslb.ID); err != nil {
		return diag.Errorf("deleting SakuraCloud GSLB[%s] is failed: %s", d.Id(), err)
	}

	return nil
}

func setGSLBResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *iaas.GSLB) diag.Diagnostics {
	d.Set("name", data.Name)                //nolint
	d.Set("fqdn", data.FQDN)                //nolint
	d.Set("sorry_server", data.SorryServer) //nolint
	d.Set("icon_id", data.IconID.String())  //nolint
	d.Set("description", data.Description)  //nolint
	d.Set("weighted", data.Weighted.Bool()) //nolint
	if err := d.Set("health_check", flattenGSLBHealthCheck(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("server", flattenGSLBServers(data)); err != nil {
		return diag.FromErr(err)
	}
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
