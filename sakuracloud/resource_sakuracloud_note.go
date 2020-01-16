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

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud/types"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudNote() *schema.Resource {
	resourceName := "Note"
	return &schema.Resource{
		Create: resourceSakuraCloudNoteCreate,
		Read:   resourceSakuraCloudNoteRead,
		Update: resourceSakuraCloudNoteUpdate,
		Delete: resourceSakuraCloudNoteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"content": {
				Type:     schema.TypeString,
				Required: true,
				Description: descf(
					"The content of the %s. This must be specified as a shell script or as a cloud-config",
					resourceName,
				),
			},
			"class": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "shell",
				ValidateFunc: validation.StringInSlice(types.NoteClassStrings, false),
				Description:  descf("The class of the %s. This must be one of %s", resourceName, types.NoteClassStrings),
			},
			"icon_id": schemaResourceIconID(resourceName),
			"tags":    schemaResourceTags(resourceName),
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("The description of the %s. This will be computed from special tags within body of `content`", resourceName),
			},
		},
	}
}

func resourceSakuraCloudNoteCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	noteOp := sacloud.NewNoteOp(client)

	note, err := noteOp.Create(ctx, expandNoteCreateRequest(d))
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Note is failed: %s", err)
	}

	d.SetId(note.ID.String())
	return resourceSakuraCloudNoteRead(d, meta)
}

func resourceSakuraCloudNoteRead(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	noteOp := sacloud.NewNoteOp(client)

	note, err := noteOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Note[%s]: %s", d.Id(), err)
	}

	return setNoteResourceData(ctx, d, client, note)
}

func resourceSakuraCloudNoteUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	noteOp := sacloud.NewNoteOp(client)

	note, err := noteOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Note[%s]: %s", d.Id(), err)
	}

	_, err = noteOp.Update(ctx, note.ID, expandNoteUpdateRequest(d))
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Note[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudNoteRead(d, meta)
}

func resourceSakuraCloudNoteDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	noteOp := sacloud.NewNoteOp(client)

	note, err := noteOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Note[%s]: %s", d.Id(), err)
	}

	if err := noteOp.Delete(ctx, note.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud Note[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setNoteResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Note) error {
	d.Set("name", data.Name)               // nolint
	d.Set("content", data.Content)         // nolint
	d.Set("class", data.Class)             // nolint
	d.Set("icon_id", data.IconID.String()) // nolint
	d.Set("description", data.Description) // nolint
	return d.Set("tags", flattenTags(data.Tags))
}
