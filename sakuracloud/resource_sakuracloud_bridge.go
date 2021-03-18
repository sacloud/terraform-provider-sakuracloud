// Copyright 2016-2021 terraform-provider-sakuracloud authors
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
	"github.com/sacloud/libsacloud/v2/helper/cleanup"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudBridge() *schema.Resource {
	resourceName := "Bridge"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudBridgeCreate,
		ReadContext:   resourceSakuraCloudBridgeRead,
		UpdateContext: resourceSakuraCloudBridgeUpdate,
		DeleteContext: resourceSakuraCloudBridgeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name":        schemaResourceName(resourceName),
			"description": schemaResourceDescription(resourceName),
			"zone":        schemaResourceZone(resourceName),
		},
	}
}

func resourceSakuraCloudBridgeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	bridgeOp := sacloud.NewBridgeOp(client)
	bridge, err := bridgeOp.Create(ctx, zone, &sacloud.BridgeCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	})
	if err != nil {
		return diag.Errorf("creating SakuraCloud Bridge is failed: %s", err)
	}

	d.SetId(bridge.ID.String())
	return resourceSakuraCloudBridgeRead(ctx, d, meta)
}

func resourceSakuraCloudBridgeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	bridgeOp := sacloud.NewBridgeOp(client)

	bridge, err := bridgeOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud Bridge[%s]: %s", d.Id(), err)
	}
	return setBridgeResourceData(ctx, d, client, bridge)
}

func resourceSakuraCloudBridgeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	bridgeOp := sacloud.NewBridgeOp(client)

	bridge, err := bridgeOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud Bridge[%s]: %s", d.Id(), err)
	}

	_, err = bridgeOp.Update(ctx, zone, bridge.ID, &sacloud.BridgeUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	})
	if err != nil {
		return diag.Errorf("updating SakuraCloud Bridge[%s] is failed: %s", d.Id(), err)
	}
	return resourceSakuraCloudBridgeRead(ctx, d, meta)
}

func resourceSakuraCloudBridgeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	bridgeOp := sacloud.NewBridgeOp(client)

	bridge, err := bridgeOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud Bridge[%s]: %s", d.Id(), err)
	}

	if err := cleanup.DeleteBridge(ctx, client, zone, client.zones, bridge.ID, client.checkReferencedOption()); err != nil {
		return diag.Errorf("deleting SakuraCloud Bridge[%s] is failed: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func setBridgeResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Bridge) diag.Diagnostics {
	d.Set("name", data.Name)               // nolint
	d.Set("description", data.Description) // nolint
	d.Set("zone", getZone(d, client))      // nolint
	return nil
}
