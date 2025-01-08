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
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudIcon() *schema.Resource {
	resourceName := "Icon"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudIconCreate,
		ReadContext:   resourceSakuraCloudIconRead,
		UpdateContext: resourceSakuraCloudIconUpdate,
		DeleteContext: resourceSakuraCloudIconDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"base64content"},
				ForceNew:      true,
				Description: desc.Sprintf(
					"The file path to upload to as the Icon. %s",
					desc.Conflicts("base64content"),
				),
			},
			"base64content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
				ForceNew:      true,
				Description: desc.Sprintf(
					"The base64 encoded content to upload to as the Icon. %s",
					desc.Conflicts("source"),
				),
			},
			"tags": schemaResourceTags(resourceName),
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for getting the icon's raw data",
			},
		},
	}
}

func resourceSakuraCloudIconCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	iconOp := iaas.NewIconOp(client)

	req, err := expandIconCreateRequest(d)
	if err != nil {
		return diag.Errorf("creating SakuraCloud Icon is failed: %s", err)
	}
	icon, err := iconOp.Create(ctx, req)
	if err != nil {
		return diag.Errorf("creating SakuraCloud Icon is failed: %s", err)
	}

	d.SetId(icon.ID.String())
	return resourceSakuraCloudIconRead(ctx, d, meta)
}

func resourceSakuraCloudIconRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	iconOp := iaas.NewIconOp(client)
	icon, err := iconOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud Icon[%s]: %s", d.Id(), err)
	}

	return setIconResourceData(ctx, d, client, icon)
}

func resourceSakuraCloudIconUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	iconOp := iaas.NewIconOp(client)
	_, err = iconOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud Icon[%s]: %s", d.Id(), err)
	}

	_, err = iconOp.Update(ctx, sakuraCloudID(d.Id()), expandIconUpdateRequest(d))
	if err != nil {
		return diag.Errorf("updating SakuraCloud Icon[%s] is failed: %s", d.Id(), err)
	}
	return resourceSakuraCloudIconRead(ctx, d, meta)
}

func resourceSakuraCloudIconDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	iconOp := iaas.NewIconOp(client)
	icon, err := iconOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud Icon[%s]: %s", d.Id(), err)
	}

	if err := iconOp.Delete(ctx, icon.ID); err != nil {
		return diag.Errorf("deleting SakuraCloud Icon[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setIconResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *iaas.Icon) diag.Diagnostics {
	d.Set("name", data.Name) // nolint
	d.Set("url", data.URL)   // nolint
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
