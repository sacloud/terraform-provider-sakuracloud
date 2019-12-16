// Copyright 2016-2019 terraform-provider-sakuracloud authors
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
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudSwitchCreate,
		Read:   resourceSakuraCloudSwitchRead,
		Update: resourceSakuraCloudSwitchUpdate,
		Delete: resourceSakuraCloudSwitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"bridge_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"server_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		},
	}
}

func resourceSakuraCloudSwitchCreate(d *schema.ResourceData, meta interface{}) error {
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	swOp := sacloud.NewSwitchOp(client)

	req := &sacloud.SwitchCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),
	}

	sw, err := swOp.Create(ctx, zone, req)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Switch is failed: %s", err)
	}

	if bridgeID, ok := d.GetOk("bridge_id"); ok {
		brID := bridgeID.(string)
		if brID != "" {
			if err := swOp.ConnectToBridge(ctx, zone, sw.ID, sakuraCloudID(brID)); err != nil {
				return fmt.Errorf("connecting Switch[%s] to Bridge[%s] is failed: %s", sw.ID, brID, err)
			}
		}
	}
	d.SetId(sw.ID.String())
	return resourceSakuraCloudSwitchRead(d, meta)
}

func resourceSakuraCloudSwitchRead(d *schema.ResourceData, meta interface{}) error {
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	swOp := sacloud.NewSwitchOp(client)

	sw, err := swOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Switch[%s] : %s", d.Id(), err)
	}
	return setSwitchResourceData(ctx, d, client, sw)
}

func resourceSakuraCloudSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	swOp := sacloud.NewSwitchOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	sw, err := swOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Switch[%s] : %s", d.Id(), err)
	}

	req := &sacloud.SwitchUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),
	}

	sw, err = swOp.Update(ctx, zone, sw.ID, req)
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Switch[%s] is failed : %s", d.Id(), err)
	}

	if d.HasChange("bridge_id") {
		if bridgeID, ok := d.GetOk("bridge_id"); ok {
			brID := bridgeID.(string)
			if brID == "" && !sw.BridgeID.IsEmpty() {
				if err := swOp.DisconnectFromBridge(ctx, zone, sw.ID); err != nil {
					return fmt.Errorf("disconnecting from Bridge[%s] is failed: %s", sw.BridgeID, err)
				}
			} else {
				if err := swOp.ConnectToBridge(ctx, zone, sw.ID, sakuraCloudID(brID)); err != nil {
					return fmt.Errorf("connecting to Bridge[%s] is failed: %s", brID, err)
				}
			}
		}
	}

	return resourceSakuraCloudSwitchRead(d, meta)
}

func resourceSakuraCloudSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	swOp := sacloud.NewSwitchOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	sw, err := swOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Switch[%s]: %s", d.Id(), err)
	}

	if !sw.BridgeID.IsEmpty() {
		if err := swOp.DisconnectFromBridge(ctx, zone, sw.ID); err != nil {
			return fmt.Errorf("disconnecting Switch[%s] from Bridge[%s] is failed: %s", sw.ID, sw.BridgeID, err)
		}
	}

	if err := waitForDeletionBySwitchID(ctx, client, zone, sw.ID); err != nil {
		return fmt.Errorf("waiting deletion is failed: Switch[%s] still used by others: %s", sw.ID, err)
	}

	if err := swOp.Delete(ctx, zone, sw.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud Switch[%s] is failed: %s", sw.ID, err)
	}
	return nil
}

func setSwitchResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Switch) error {
	zone := getZone(d, client)
	var serverIDs []string
	if data.ServerCount > 0 {
		swOp := sacloud.NewSwitchOp(client)
		searched, err := swOp.GetServers(ctx, zone, data.ID)
		if err != nil {
			return fmt.Errorf("could not find SakuraCloud Servers: switch[%s]", err)
		}
		for _, s := range searched.Servers {
			serverIDs = append(serverIDs, s.ID.String())
		}
	}

	d.Set("name", data.Name)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	d.Set("bridge_id", data.BridgeID.String())
	if err := d.Set("server_ids", serverIDs); err != nil {
		return err
	}
	d.Set("zone", zone)
	return nil
}
