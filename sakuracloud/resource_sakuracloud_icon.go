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
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudIcon() *schema.Resource {
	resourceName := "Icon"
	return &schema.Resource{
		Create: resourceSakuraCloudIconCreate,
		Read:   resourceSakuraCloudIconRead,
		Update: resourceSakuraCloudIconUpdate,
		Delete: resourceSakuraCloudIconDelete,

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
				Description: descf(
					"The file path to upload to as the Icon. %s",
					descConflicts("base64content"),
				),
			},
			"base64content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
				ForceNew:      true,
				Description: descf(
					"The base64 encoded content to upload to as the Icon. %s",
					descConflicts("source"),
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

func resourceSakuraCloudIconCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	iconOp := sacloud.NewIconOp(client)

	req, err := expandIconCreateRequest(d)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Icon is failed: %s", err)
	}
	icon, err := iconOp.Create(ctx, req)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Icon is failed: %s", err)
	}

	d.SetId(icon.ID.String())
	return resourceSakuraCloudIconRead(d, meta)
}

func resourceSakuraCloudIconRead(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	iconOp := sacloud.NewIconOp(client)

	icon, err := iconOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Icon[%s]: %s", d.Id(), err)
	}

	return setIconResourceData(ctx, d, client, icon)
}

func resourceSakuraCloudIconUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	iconOp := sacloud.NewIconOp(client)

	_, err = iconOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Icon[%s]: %s", d.Id(), err)
	}

	_, err = iconOp.Update(ctx, sakuraCloudID(d.Id()), expandIconUpdateRequest(d))
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Icon[%s] is failed: %s", d.Id(), err)
	}
	return resourceSakuraCloudIconRead(d, meta)
}

func resourceSakuraCloudIconDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	iconOp := sacloud.NewIconOp(client)

	icon, err := iconOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Icon[%s]: %s", d.Id(), err)
	}

	if err := iconOp.Delete(ctx, icon.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud Icon[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setIconResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Icon) error {
	d.Set("name", data.Name) // nolint
	d.Set("url", data.URL)   // nolint
	return d.Set("tags", flattenTags(data.Tags))
}
