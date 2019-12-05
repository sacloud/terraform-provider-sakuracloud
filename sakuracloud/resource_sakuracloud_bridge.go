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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudBridge() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudBridgeCreate,
		Read:   resourceSakuraCloudBridgeRead,
		Update: resourceSakuraCloudBridgeUpdate,
		Delete: resourceSakuraCloudBridgeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"switch_ids": {
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
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a"}),
			},
		},
	}
}

func resourceSakuraCloudBridgeCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	bridgeOp := sacloud.NewBridgeOp(client)

	req := &sacloud.BridgeCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	bridge, err := bridgeOp.Create(ctx, zone, req)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Bridge is failed: %s", err)
	}

	d.SetId(bridge.ID.String())
	return resourceSakuraCloudBridgeRead(d, meta)
}

func resourceSakuraCloudBridgeRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	bridgeOp := sacloud.NewBridgeOp(client)
	bridge, err := bridgeOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Bridge[%s]: %s", d.Id(), err)
	}
	return setBridgeResourceData(ctx, d, client, bridge)
}

func resourceSakuraCloudBridgeUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	bridgeOp := sacloud.NewBridgeOp(client)

	bridge, err := bridgeOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Bridge[%s]: %s", d.Id(), err)
	}

	req := &sacloud.BridgeUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	bridge, err = bridgeOp.Update(ctx, zone, bridge.ID, req)
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Bridge[%s] is failed: %s", d.Id(), err)
	}
	return resourceSakuraCloudBridgeRead(d, meta)
}

func resourceSakuraCloudBridgeDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	bridgeOp := sacloud.NewBridgeOp(client)

	bridge, err := bridgeOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Bridge[%s]: %s", d.Id(), err)
	}

	if err := bridgeOp.Delete(ctx, zone, bridge.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud AutoBackup[%s] is failed: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func setBridgeResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Bridge) error {
	d.Set("name", data.Name)
	d.Set("description", data.Description)

	swOp := sacloud.NewSwitchOp(client)
	var switchIDs []interface{}
	for _, d := range data.BridgeInfo {
		if _, err := swOp.Read(ctx, d.ZoneName, d.ID); err == nil {
			switchIDs = append(switchIDs, d.ID.String())
		}
	}
	if err := d.Set("switch_ids", switchIDs); err != nil {
		return fmt.Errorf("error setting switch_ids: %v", switchIDs)
	}

	d.Set("zone", getZone(d, client))
	return nil
}
