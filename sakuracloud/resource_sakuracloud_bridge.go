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
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/libsacloud/v2/helper/cleanup"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudBridge() *schema.Resource {
	resourceName := "Bridge"
	return &schema.Resource{
		Create: resourceSakuraCloudBridgeCreate,
		Read:   resourceSakuraCloudBridgeRead,
		Update: resourceSakuraCloudBridgeUpdate,
		Delete: resourceSakuraCloudBridgeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceSakuraCloudBridgeCreate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	bridgeOp := sacloud.NewBridgeOp(client)

	bridge, err := bridgeOp.Create(ctx, zone, &sacloud.BridgeCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	})
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Bridge is failed: %s", err)
	}

	d.SetId(bridge.ID.String())
	return resourceSakuraCloudBridgeRead(d, meta)
}

func resourceSakuraCloudBridgeRead(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

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
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	bridgeOp := sacloud.NewBridgeOp(client)

	bridge, err := bridgeOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Bridge[%s]: %s", d.Id(), err)
	}

	_, err = bridgeOp.Update(ctx, zone, bridge.ID, &sacloud.BridgeUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	})
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Bridge[%s] is failed: %s", d.Id(), err)
	}
	return resourceSakuraCloudBridgeRead(d, meta)
}

func resourceSakuraCloudBridgeDelete(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	bridgeOp := sacloud.NewBridgeOp(client)

	bridge, err := bridgeOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Bridge[%s]: %s", d.Id(), err)
	}

	if err := cleanup.DeleteBridge(ctx, client, zone, client.zones, bridge.ID, client.checkReferencedOption()); err != nil {
		return fmt.Errorf("deleting SakuraCloud Bridge[%s] is failed: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func setBridgeResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Bridge) error {
	d.Set("name", data.Name)               // nolint
	d.Set("description", data.Description) // nolint
	d.Set("zone", getZone(d, client))      // nolint
	return nil
}
