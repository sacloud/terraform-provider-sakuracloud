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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudNote() *schema.Resource {
	resourceName := "Note"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudNoteCreate,
		ReadContext:   resourceSakuraCloudNoteRead,
		UpdateContext: resourceSakuraCloudNoteUpdate,
		DeleteContext: resourceSakuraCloudNoteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceSakuraCloudNoteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	noteOp := sacloud.NewNoteOp(client)
	note, err := noteOp.Create(ctx, expandNoteCreateRequest(d))
	if err != nil {
		return diag.Errorf("creating SakuraCloud Note is failed: %s", err)
	}

	d.SetId(note.ID.String())
	return resourceSakuraCloudNoteRead(ctx, d, meta)
}

func resourceSakuraCloudNoteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	noteOp := sacloud.NewNoteOp(client)
	note, err := noteOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud Note[%s]: %s", d.Id(), err)
	}

	return setNoteResourceData(ctx, d, client, note)
}

func resourceSakuraCloudNoteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	noteOp := sacloud.NewNoteOp(client)
	note, err := noteOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud Note[%s]: %s", d.Id(), err)
	}

	_, err = noteOp.Update(ctx, note.ID, expandNoteUpdateRequest(d))
	if err != nil {
		return diag.Errorf("updating SakuraCloud Note[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudNoteRead(ctx, d, meta)
}

func resourceSakuraCloudNoteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	noteOp := sacloud.NewNoteOp(client)
	note, err := noteOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud Note[%s]: %s", d.Id(), err)
	}

	if err := noteOp.Delete(ctx, note.ID); err != nil {
		return diag.Errorf("deleting SakuraCloud Note[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setNoteResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Note) diag.Diagnostics {
	d.Set("name", data.Name)               // nolint
	d.Set("content", data.Content)         // nolint
	d.Set("class", data.Class)             // nolint
	d.Set("icon_id", data.IconID.String()) // nolint
	d.Set("description", data.Description) // nolint
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
