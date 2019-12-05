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
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudNote() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudNoteCreate,
		Read:   resourceSakuraCloudNoteRead,
		Update: resourceSakuraCloudNoteUpdate,
		Delete: resourceSakuraCloudNoteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
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
				Computed: true,
			},
			"class": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "shell",
				ValidateFunc: validation.StringInSlice([]string{
					"shell",
					"yaml_cloud_config",
				}, false),
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSakuraCloudNoteCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	noteOp := sacloud.NewNoteOp(client)

	note, err := noteOp.Create(ctx, &sacloud.NoteCreateRequest{
		Name:    d.Get("name").(string),
		Tags:    expandTags(d),
		IconID:  expandSakuraCloudID(d, "icon_id"),
		Class:   d.Get("class").(string),
		Content: d.Get("content").(string),
	})
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Note is failed: %s", err)
	}

	d.SetId(note.ID.String())
	return resourceSakuraCloudNoteRead(d, meta)
}

func resourceSakuraCloudNoteRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	noteOp := sacloud.NewNoteOp(client)

	note, err := noteOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Note: %s", err)
	}

	return setNoteResourceData(ctx, d, client, note)
}

func resourceSakuraCloudNoteUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	noteOp := sacloud.NewNoteOp(client)

	note, err := noteOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Note: %s", err)
	}

	_, err = noteOp.Update(ctx, note.ID, &sacloud.NoteUpdateRequest{
		Name:    d.Get("name").(string),
		Tags:    expandTags(d),
		IconID:  expandSakuraCloudID(d, "icon_id"),
		Class:   d.Get("class").(string),
		Content: d.Get("content").(string),
	})
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Note is failed: %s", err)
	}

	return resourceSakuraCloudNoteRead(d, meta)
}

func resourceSakuraCloudNoteDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	noteOp := sacloud.NewNoteOp(client)

	note, err := noteOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Note: %s", err)
	}

	if err := noteOp.Delete(ctx, note.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud Note is failed: %s", err)
	}
	return nil
}

func setNoteResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Note) error {
	d.Set("name", data.Name)
	d.Set("content", data.Content)
	d.Set("class", data.Class)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	return nil
}
