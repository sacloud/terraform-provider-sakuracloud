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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/iaas-api-go"
)

func resourceSakuraCloudAutoScale() *schema.Resource {
	resourceName := "AutoScale"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudAutoScaleCreate,
		ReadContext:   resourceSakuraCloudAutoScaleRead,
		UpdateContext: resourceSakuraCloudAutoScaleUpdate,
		DeleteContext: resourceSakuraCloudAutoScaleDelete,
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
			"zones": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: descf("List of zone names where monitored resources are located"),
			},
			"config": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validateAutoScaleConfig),
				Description:      "The configuration file for sacloud/autoscaler",
			},
			"api_key_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
				Description:      "The disk id to backed up",
				ForceNew:         true,
			},
			"cpu_threshold_scaling": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_prefix": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Server name prefix to be monitored",
						},
						"up": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Threshold for average CPU utilization to scale up/out",
						},
						"down": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Threshold for average CPU utilization to scale down/in",
						},
					},
				},
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
		},
		DeprecationMessage: "sakuracloud_auto_scale is an experimental resource. Please note that you will need to update the tfstate manually if the resource schema is changed.",
	}
}

func resourceSakuraCloudAutoScaleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	autoScaleOp := iaas.NewAutoScaleOp(client)
	autoScale, err := autoScaleOp.Create(ctx, expandAutoScaleCreateRequest(d))
	if err != nil {
		return diag.Errorf("creating SakuraCloud AutoScale is failed: %s", err)
	}

	d.SetId(autoScale.ID.String())
	return resourceSakuraCloudAutoScaleRead(ctx, d, meta)
}

func resourceSakuraCloudAutoScaleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	autoScaleOp := iaas.NewAutoScaleOp(client)
	autoScale, err := autoScaleOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not find SakuraCloud AutoScale[%s]: %s", d.Id(), err)
	}
	return setAutoScaleResourceData(d, client, autoScale)
}

func resourceSakuraCloudAutoScaleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	autoScaleOp := iaas.NewAutoScaleOp(client)
	autoScale, err := autoScaleOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud AutoScale[%s]: %s", d.Id(), err)
	}

	if _, err = autoScaleOp.Update(ctx, autoScale.ID, expandAutoScaleUpdateRequest(d, autoScale)); err != nil {
		return diag.Errorf("updating SakuraCloud AutoScale[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudAutoScaleRead(ctx, d, meta)
}

func resourceSakuraCloudAutoScaleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	autoScaleOp := iaas.NewAutoScaleOp(client)
	autoScale, err := autoScaleOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud AutoScale[%s]: %s", d.Id(), err)
	}

	if err := autoScaleOp.Delete(ctx, autoScale.ID); err != nil {
		return diag.Errorf("deleting SakuraCloud AutoScale[%s] is failed: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func setAutoScaleResourceData(d *schema.ResourceData, client *APIClient, data *iaas.AutoScale) diag.Diagnostics {
	d.Set("name", data.Name) // nolint

	if err := d.Set("zones", data.Zones); err != nil {
		return diag.FromErr(err)
	}
	d.Set("config", data.Config)       // nolint
	d.Set("api_key_id", data.APIKeyID) // nolint
	if err := d.Set("cpu_threshold_scaling", flattenAutoScaleCPUThresholdScaling(data.CPUThresholdScaling)); err != nil {
		return diag.FromErr(err)
	}

	d.Set("icon_id", data.IconID.String()) // nolint
	d.Set("description", data.Description) // nolint
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
