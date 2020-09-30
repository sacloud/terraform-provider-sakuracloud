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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudArchiveShare() *schema.Resource {
	resourceName := "ArchiveShare"

	return &schema.Resource{
		Create: resourceSakuraCloudArchiveShareCreate,
		Read:   resourceSakuraCloudArchiveShareRead,
		Delete: resourceSakuraCloudArchiveShareDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"archive_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateSakuracloudIDType,
				Description:  "The id of the archive",
			},
			"share_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The key to use sharing the Archive",
			},
			"zone": schemaResourceZone(resourceName),
		},
	}
}

func resourceSakuraCloudArchiveShareCreate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	archiveOp := sacloud.NewArchiveOp(client)
	archiveID := expandSakuraCloudID(d, "archive_id")

	archive, err := archiveOp.Read(ctx, zone, archiveID)
	if err != nil {
		return fmt.Errorf("sharing SakuraCloud Archive is failed: %s", err)
	}

	// share
	shareInfo, err := archiveOp.Share(ctx, zone, archiveID)
	if err != nil {
		return fmt.Errorf("sharing SakuraCloud Archive is failed: %s", err)
	}

	d.SetId(archive.ID.String())
	d.Set("share_key", shareInfo.SharedKey)
	d.Set("zone", zone)
	return nil
}

func resourceSakuraCloudArchiveShareRead(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	archiveOp := sacloud.NewArchiveOp(client)

	archive, err := archiveOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Archive[%s]: %s", d.Id(), err)
	}

	if !archive.Availability.IsUploading() {
		d.SetId("")
		return nil
	}

	d.SetId(archive.ID.String())
	d.Set("share_key", d.Get("share_key").(string))
	d.Set("zone", zone)
	return nil
}

func resourceSakuraCloudArchiveShareDelete(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	archiveOp := sacloud.NewArchiveOp(client)

	archive, err := archiveOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Archive[%s]: %s", d.Id(), err)
	}

	if !archive.Availability.IsUploading() {
		d.SetId("")
		return nil
	}

	if err := archiveOp.CloseFTP(ctx, zone, archive.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud Archive Share[%s] is failed: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
