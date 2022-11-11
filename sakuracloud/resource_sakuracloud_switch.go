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
	"github.com/sacloud/iaas-api-go/helper/cleanup"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudSwitch() *schema.Resource {
	resourceName := "Switch"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudSwitchCreate,
		ReadContext:   resourceSakuraCloudSwitchRead,
		UpdateContext: resourceSakuraCloudSwitchUpdate,
		DeleteContext: resourceSakuraCloudSwitchDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name":        schemaResourceName(resourceName),
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"bridge_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
				Description:      desc.Sprintf("The bridge id attached to the %s", resourceName),
			},
			"server_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of server id connected to the switch",
			},
			"zone": schemaResourceZone(resourceName),
		},
	}
}

func resourceSakuraCloudSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	swOp := iaas.NewSwitchOp(client)
	req := &iaas.SwitchCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),
	}

	sw, err := swOp.Create(ctx, zone, req)
	if err != nil {
		return diag.Errorf("creating SakuraCloud Switch is failed: %s", err)
	}

	if bridgeID, ok := d.GetOk("bridge_id"); ok {
		brID := bridgeID.(string)
		if brID != "" {
			if err := swOp.ConnectToBridge(ctx, zone, sw.ID, sakuraCloudID(brID)); err != nil {
				return diag.Errorf("connecting Switch[%s] to Bridge[%s] is failed: %s", sw.ID, brID, err)
			}
		}
	}
	d.SetId(sw.ID.String())
	return resourceSakuraCloudSwitchRead(ctx, d, meta)
}

func resourceSakuraCloudSwitchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	swOp := iaas.NewSwitchOp(client)
	sw, err := swOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud Switch[%s] : %s", d.Id(), err)
	}
	return setSwitchResourceData(ctx, d, client, sw)
}

func resourceSakuraCloudSwitchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	swOp := iaas.NewSwitchOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	sw, err := swOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud Switch[%s] : %s", d.Id(), err)
	}

	req := &iaas.SwitchUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),
	}

	sw, err = swOp.Update(ctx, zone, sw.ID, req)
	if err != nil {
		return diag.Errorf("updating SakuraCloud Switch[%s] is failed : %s", d.Id(), err)
	}

	if d.HasChange("bridge_id") {
		if bridgeID, ok := d.GetOk("bridge_id"); ok {
			brID := bridgeID.(string)
			if brID == "" && !sw.BridgeID.IsEmpty() {
				if err := swOp.DisconnectFromBridge(ctx, zone, sw.ID); err != nil {
					return diag.Errorf("disconnecting from Bridge[%s] is failed: %s", sw.BridgeID, err)
				}
			} else {
				if err := swOp.ConnectToBridge(ctx, zone, sw.ID, sakuraCloudID(brID)); err != nil {
					return diag.Errorf("connecting to Bridge[%s] is failed: %s", brID, err)
				}
			}
		}
	}

	return resourceSakuraCloudSwitchRead(ctx, d, meta)
}

func resourceSakuraCloudSwitchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	swOp := iaas.NewSwitchOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	sw, err := swOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud Switch[%s]: %s", d.Id(), err)
	}

	if !sw.BridgeID.IsEmpty() {
		if err := swOp.DisconnectFromBridge(ctx, zone, sw.ID); err != nil {
			return diag.Errorf("disconnecting Switch[%s] from Bridge[%s] is failed: %s", sw.ID, sw.BridgeID, err)
		}
	}

	if err := cleanup.DeleteSwitch(ctx, client, zone, sw.ID, client.checkReferencedOption()); err != nil {
		return diag.Errorf("deleting SakuraCloud Switch[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setSwitchResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *iaas.Switch) diag.Diagnostics {
	zone := getZone(d, client)
	var serverIDs []string
	if data.ServerCount > 0 {
		swOp := iaas.NewSwitchOp(client)
		searched, err := swOp.GetServers(ctx, zone, data.ID)
		if err != nil {
			return diag.Errorf("could not find SakuraCloud Servers: switch[%s]", err)
		}
		for _, s := range searched.Servers {
			serverIDs = append(serverIDs, s.ID.String())
		}
	}

	d.Set("name", data.Name)                   // nolint
	d.Set("icon_id", data.IconID.String())     // nolint
	d.Set("description", data.Description)     // nolint
	d.Set("bridge_id", data.BridgeID.String()) // nolint
	d.Set("zone", zone)                        // nolint
	if err := d.Set("server_ids", serverIDs); err != nil {
		return diag.FromErr(err)
	}
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
