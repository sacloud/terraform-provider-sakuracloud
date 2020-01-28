// Copyright 2016-2020 terraform-provider-sakuracloud authors
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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
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
	client := getSacloudAPIClient(d, meta)

	opts := client.Bridge.New()

	opts.Name = d.Get("name").(string)
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}

	bridge, err := client.Bridge.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Bridge resource: %s", err)
	}

	d.SetId(bridge.GetStrID())
	return resourceSakuraCloudBridgeRead(d, meta)
}

func resourceSakuraCloudBridgeRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	bridge, err := client.Bridge.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud Bridge resource: %s", err)
	}

	return setBridgeResourceData(d, client, bridge)
}

func resourceSakuraCloudBridgeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	bridge, err := client.Bridge.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Bridge resource: %s", err)
	}

	if d.HasChange("name") {
		bridge.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			bridge.Description = description.(string)
		} else {
			bridge.Description = ""
		}
	}

	_, err = client.Bridge.Update(bridge.ID, bridge)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Bridge resource: %s", err)
	}

	return resourceSakuraCloudBridgeRead(d, meta)
}

func resourceSakuraCloudBridgeDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	br, err := client.Bridge.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Bridge resource: %s", err)
	}

	if br.Info != nil && br.Info.Switches != nil && len(br.Info.Switches) > 0 {
		for _, s := range br.Info.Switches {
			switchID, _ := s.ID.Int64()
			strSwitchID := s.ID.String()
			sakuraMutexKV.Lock(strSwitchID)
			defer sakuraMutexKV.Unlock(strSwitchID)

			if _, err := client.Switch.Read(switchID); err != nil {
				if err != nil {
					return fmt.Errorf("Error disconnecting Bridge resource: %s", err)
				}
			}
			if _, err := client.Switch.DisconnectFromBridge(switchID); err != nil {
				return fmt.Errorf("Error disconnecting Bridge resource: %s", err)
			}
		}
	}

	_, err = client.Bridge.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Bridge resource: %s", err)
	}
	return nil
}

func setBridgeResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.Bridge) error {
	d.Set("name", data.Name)
	d.Set("description", data.Description)

	var switchIDs []interface{}
	if data.Info != nil && data.Info.Switches != nil && len(data.Info.Switches) > 0 {

		for _, d := range data.Info.Switches {
			swID := d.ID.String()
			sakuraMutexKV.Lock(swID)
			defer sakuraMutexKV.Unlock(swID)

			id, _ := d.ID.Int64()
			if _, err := client.Switch.Read(id); err == nil {
				switchIDs = append(switchIDs, d.ID.String())
			}
		}
	}
	d.Set("switch_ids", switchIDs)

	d.Set("zone", client.Zone)
	return nil
}
