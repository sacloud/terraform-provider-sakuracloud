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

func resourceSakuraCloudWebAccelACL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSakuraCloudWebAccelACLCreate,
		ReadContext:   resourceSakuraCloudWebAccelACLRead,
		UpdateContext: resourceSakuraCloudWebAccelACLUpdate,
		DeleteContext: resourceSakuraCloudWebAccelACLDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"acl": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSakuraCloudWebAccelACLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	siteID := d.Get("site_id").(string)

	_, err = webaccel.NewOp(client.webaccelClient).UpsertACL(ctx, siteID, d.Get("acl").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(siteID)
	return resourceSakuraCloudWebAccelACLRead(ctx, d, meta)
}

func resourceSakuraCloudWebAccelACLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	siteID := d.Id()

	acl, err := webaccel.NewOp(client.webaccelClient).ReadACL(ctx, siteID)
	if err != nil {
		return diag.Errorf("could not read SakuraCloud WebAccel ACL[%s]: %s", d.Id(), err)
	}

	if acl.ACL == "" {
		d.SetId("")
		return nil
	}

	return setWebAccelACLResourceData(d, client, acl)
}

func resourceSakuraCloudWebAccelACLUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	siteID := d.Id()

	if d.HasChanges("acl") {
		_, err := webaccel.NewOp(client.webaccelClient).UpsertACL(ctx, siteID, d.Get("acl").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSakuraCloudWebAccelACLRead(ctx, d, meta)
}

func resourceSakuraCloudWebAccelACLDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	siteID := d.Get("site_id").(string)

	if err := webaccel.NewOp(client.webaccelClient).DeleteACL(ctx, siteID); err != nil {
		return diag.Errorf("deleting SakuraCloud WebAccel ACL[%s] is failed: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func setWebAccelACLResourceData(d *schema.ResourceData, client *APIClient, data *webaccel.ACLResult) diag.Diagnostics {
	d.Set("site_id", d.Id()) // nolint
	d.Set("acl", data.ACL)   // nolint
	return nil
}
