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
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudIcon() *schema.Resource {
	return &schema.Resource{
		Create:        resourceSakuraCloudIconCreate,
		Read:          resourceSakuraCloudIconRead,
		Update:        resourceSakuraCloudIconUpdate,
		Delete:        resourceSakuraCloudIconDelete,
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"base64content"},
				ForceNew:      true,
			},
			"base64content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
				ForceNew:      true,
			},
			"body": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudIconCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	iconOp := sacloud.NewIconOp(client)

	body, err := expandIconBody(d)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Icon is failed: %s", err)
	}

	icon, err := iconOp.Create(ctx, &sacloud.IconCreateRequest{
		Name:  d.Get("name").(string),
		Tags:  expandTagsV2(d.Get("tags").([]interface{})),
		Image: body,
	})
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Icon is failed: %s", err)
	}

	d.SetId(icon.ID.String())
	return resourceSakuraCloudIconRead(d, meta)
}

func resourceSakuraCloudIconRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	iconOp := sacloud.NewIconOp(client)

	icon, err := iconOp.Read(ctx, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Icon: %s", err)
	}

	return setIconResourceData(ctx, d, client, icon)
}

func resourceSakuraCloudIconUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	iconOp := sacloud.NewIconOp(client)

	_, err := iconOp.Read(ctx, types.StringID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Icon: %s", err)
	}

	_, err = iconOp.Update(ctx, types.StringID(d.Id()), &sacloud.IconUpdateRequest{
		Name: d.Get("name").(string),
		Tags: expandTagsV2(d.Get("tags").([]interface{})),
	})
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Icon is failed: %s", err)
	}
	return resourceSakuraCloudIconRead(d, meta)
}

func resourceSakuraCloudIconDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	iconOp := sacloud.NewIconOp(client)

	icon, err := iconOp.Read(ctx, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Icon: %s", err)
	}

	if err := iconOp.Delete(ctx, icon.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud Icon is failed: %s", err)
	}
	return nil
}

func setIconResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Icon) error {
	d.Set("name", data.Name)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	d.Set("url", data.URL)
	return nil
}

func expandIconBody(d resourceValueGettable) (string, error) {
	var body string
	if v, ok := d.GetOk("source"); ok {
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return "", fmt.Errorf("expanding homedir in source (%s) is failed: %s", source, err)
		}
		file, err := os.Open(path)
		if err != nil {
			return "", fmt.Errorf("opening SakuraCloud Icon source(%s) is failed: %s", source, err)
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return "", fmt.Errorf("reading SakuraCloud Icon source file is failed: %s", err)
		}
		body = base64.StdEncoding.EncodeToString(data)
	} else if v, ok := d.GetOk("base64content"); ok {
		body = v.(string)
	} else {
		return "", fmt.Errorf(`"source" or "base64content" field is required`)
	}
	return body, nil
}
