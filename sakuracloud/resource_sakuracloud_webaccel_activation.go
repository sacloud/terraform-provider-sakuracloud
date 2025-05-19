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
	"github.com/sacloud/webaccel-api-go"
)

func resourceSakuraCloudWebAccelActivation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSakuraCloudWebAccelActivationCreate,
		ReadContext:   resourceSakuraCloudWebAccelActivationRead,
		UpdateContext: resourceSakuraCloudWebAccelActivationUpdate,
		DeleteContext: resourceSakuraCloudWebAccelActivationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"site_id": {
				Type:        schema.TypeString,
				Description: "target site id",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "whether the site activation is enabled or not",
				Required:    true,
			},
		},
	}
}

func resourceSakuraCloudWebAccelActivationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	siteID := d.Get("site_id").(string)

	statusString := "disabled"
	if d.Get("enabled").(bool) {
		statusString = "enabled"
	}

	op := webaccel.NewOp(client.webaccelClient)
	site, err := op.Read(ctx, siteID)
	if err != nil {
		return diag.FromErr(err)
	}

	//for avoiding status update confliction
	if statusString == site.Status {
		return resourceSakuraCloudWebAccelActivationRead(ctx, d, meta)
	}

	_, err = op.UpdateStatus(ctx, siteID, &webaccel.UpdateSiteStatusRequest{Status: statusString})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(site.ID)
	return resourceSakuraCloudWebAccelActivationRead(ctx, d, meta)
}

func resourceSakuraCloudWebAccelActivationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	siteID := d.Get("site_id").(string)

	site, err := webaccel.NewOp(client.webaccelClient).Read(ctx, siteID)
	if err != nil {
		return diag.Errorf("could not read SakuraCloud WebAccel Activation[%s]: %s", d.Id(), err)
	}

	d.SetId(site.ID)

	return setWebAccelActivationResourceData(d, client, site)
}

func resourceSakuraCloudWebAccelActivationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	siteID := d.Id()

	if d.HasChange("enabled") {
		statusString := "disabled"
		if d.Get("enabled").(bool) {
			statusString = "enabled"
		}
		_, err := webaccel.NewOp(client.webaccelClient).UpdateStatus(ctx, siteID, &webaccel.UpdateSiteStatusRequest{Status: statusString})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSakuraCloudWebAccelActivationRead(ctx, d, meta)
}

func resourceSakuraCloudWebAccelActivationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	siteID := d.Get("site_id").(string)

	if _, err = webaccel.NewOp(client.webaccelClient).UpdateStatus(ctx, siteID, &webaccel.UpdateSiteStatusRequest{Status: "disabled"}); err != nil {
		return diag.Errorf("deleting SakuraCloud WebAccel Activation[%s] is failed: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func setWebAccelActivationResourceData(d *schema.ResourceData, client *APIClient, data *webaccel.Site) diag.Diagnostics {
	d.SetId(data.ID)
	d.Set("enabled", data.Status == "enabled") // nolint
	return nil
}
