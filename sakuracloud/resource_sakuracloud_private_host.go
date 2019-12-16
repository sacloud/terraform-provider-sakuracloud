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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	phOp := sacloud.NewPrivateHostOp(client)

	planID, err := expandPrivateHostPlanID(ctx, d, client, zone)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud PrivateHost is failed: %s", err)
	}

	ph, err := phOp.Create(ctx, zone, expandPrivateHostCreateRequest(d, planID))
	if err != nil {
		return fmt.Errorf("creating SakuraCloud PrivateHost is failed: %s", err)
	}

	d.SetId(ph.ID.String())
	return resourceSakuraCloudPrivateHostRead(d, meta)
}

func resourceSakuraCloudPrivateHostRead(d *schema.ResourceData, meta interface{}) error {
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	phOp := sacloud.NewPrivateHostOp(client)

	ph, err := phOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud PrivateHost[%s]: %s", d.Id(), err)
	}
	return setPrivateHostResourceData(ctx, d, client, ph)
}

func resourceSakuraCloudPrivateHostUpdate(d *schema.ResourceData, meta interface{}) error {
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	phOp := sacloud.NewPrivateHostOp(client)

	ph, err := phOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud PrivateHost[%s]: %s", d.Id(), err)
	}

	_, err = phOp.Update(ctx, zone, ph.ID, expandPrivateHostUpdateRequest(d))
	if err != nil {
		return fmt.Errorf("updating SakuraCloud PrivateHost[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudPrivateHostRead(d, meta)
}

func resourceSakuraCloudPrivateHostDelete(d *schema.ResourceData, meta interface{}) error {
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	phOp := sacloud.NewPrivateHostOp(client)

	ph, err := phOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud PrivateHost[%s]: %s", d.Id(), err)
	}

	if err := waitForDeletionByPrivateHostID(ctx, client, zone, ph.ID); err != nil {
		return fmt.Errorf("waiting deletion is failed: PrivateHost[%s] still used by Server: %s", ph.ID, err)
	}

	if err := phOp.Delete(ctx, zone, ph.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud PrivateHost[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setPrivateHostResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.PrivateHost) error {
	d.Set("name", data.Name)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	d.Set("hostname", data.GetHostName())
	d.Set("assigned_core", data.GetAssignedCPU())
	d.Set("assigned_memory", data.GetAssignedMemoryGB())
	d.Set("zone", getZone(d, client))
	return d.Set("tags", data.Tags)
}
