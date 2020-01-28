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

func resourceSakuraCloudPrivateHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudPrivateHostCreate,
		Read:   resourceSakuraCloudPrivateHostRead,
		Update: resourceSakuraCloudPrivateHostUpdate,
		Delete: resourceSakuraCloudPrivateHostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
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
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			powerManageTimeoutKey: powerManageTimeoutParam,
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assigned_core": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"assigned_memory": {
				Type:     schema.TypeInt,
				Computed: true,
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

func resourceSakuraCloudPrivateHostCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	opts := client.PrivateHost.New()

	opts.Name = d.Get("name").(string)
	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.SetIconByID(toSakuraCloudID(iconID.(string)))
	}
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	if rawTags, ok := d.GetOk("tags"); ok {
		if rawTags != nil {
			opts.Tags = expandTags(client, rawTags.([]interface{}))
		}
	}

	plans, err := client.Product.GetProductPrivateHostAPI().FilterBy("Class", "dynamic").Find()
	if err != nil || len(plans.PrivateHostPlans) == 0 {
		return fmt.Errorf("Failed to create SakuraCloud PrivateHost resource: %s", err)
	}
	plan := plans.PrivateHostPlans[0]
	opts.SetPrivateHostPlan(&plan)

	privateHost, err := client.PrivateHost.Create(opts)
	if err != nil {
		return fmt.Errorf("plan is not found on zone[%s]: %s", client.Zone, err)
	}
	d.SetId(privateHost.GetStrID())
	return resourceSakuraCloudPrivateHostRead(d, meta)
}

func resourceSakuraCloudPrivateHostRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	privateHost, err := client.PrivateHost.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud Server PrivateHost resource: %s", err)
	}

	return setPrivateHostResourceData(d, client, privateHost)
}

func resourceSakuraCloudPrivateHostUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	privateHost, err := client.PrivateHost.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud PrivateHost resource: %s", err)
	}

	if d.HasChange("name") {
		privateHost.Name = d.Get("name").(string)
	}
	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			privateHost.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			privateHost.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			privateHost.Description = description.(string)
		} else {
			privateHost.Description = ""
		}
	}
	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			privateHost.Tags = expandTags(client, rawTags)
		} else {
			privateHost.Tags = expandTags(client, []interface{}{})
		}
	}

	_, err = client.PrivateHost.Update(privateHost.ID, privateHost)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud PrivateHost resource: %s", err)
	}

	return resourceSakuraCloudPrivateHostRead(d, meta)
}

func resourceSakuraCloudPrivateHostDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	id := toSakuraCloudID(d.Id())
	// shutdown connected servers
	res, err := client.Server.Find()
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud PrivateHost resource: finding server is failed: %s", err)
	}
	var needRestartServers []int64

	for _, s := range res.Servers {
		if s.PrivateHost != nil && s.PrivateHost.ID == id {

			// still running?
			isNeedRestart := false
			if s.IsUp() {
				isNeedRestart = true
				err := stopServer(client, s.ID, d)
				if err != nil {
					return fmt.Errorf("Error deleting SakuraCloud PrivateHost resource: shutdown server is failed: %s", err)
				}
			}

			s.ClearPrivateHost()
			_, err := client.Server.Update(s.ID, &s)
			if err != nil {
				return fmt.Errorf("Error deleting SakuraCloud PrivateHost resource: clear private-host is failed: %s", err)
			}

			if isNeedRestart {
				needRestartServers = append(needRestartServers, s.ID)
			}
		}
	}

	_, err = client.PrivateHost.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud PrivateHost resource: %s", err)
	}

	for _, serverID := range needRestartServers {
		err := bootServer(client, serverID)
		if err != nil {
			return fmt.Errorf("Error deleting SakuraCloud PrivateHost resource: booting server is failed: %s", err)
		}
	}

	return nil
}

func setPrivateHostResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.PrivateHost) error {
	d.Set("name", data.Name)
	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	d.Set("tags", data.Tags)
	d.Set("hostname", data.GetHostName())
	d.Set("assigned_core", data.GetAssignedCPU())
	d.Set("assigned_memory", data.GetAssignedMemoryGB())
	setPowerManageTimeoutValueToState(d)
	d.Set("zone", client.Zone)
	return nil
}
