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

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudArchive() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudArchiveCreate,
		Read:   resourceSakuraCloudArchiveRead,
		Update: resourceSakuraCloudArchiveUpdate,
		Delete: resourceSakuraCloudArchiveDelete,
		CustomizeDiff: customdiff.All(
			customdiff.ComputedIf("hash", func(d *schema.ResourceDiff, meta interface{}) bool {
				return d.HasChange("archive_file")
			}),
			hasTagResourceCustomizeDiff,
		),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(24 * time.Hour),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(24 * time.Hour),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      20,
				ValidateFunc: validation.IntInSlice(types.ValidArchiveSizes),
			},
			"archive_file": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hash": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validation.StringInSlice([]string{"is1a", "is1b", "tk1a", "tk1v"}, false),
			},
		},
	}
}

func resourceSakuraCloudArchiveCreate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	archiveOp := sacloud.NewArchiveOp(client)

	archive, ftpServer, err := archiveOp.CreateBlank(ctx, zone, expandArchiveCreateRequest(d))
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Archive is failed: %s", err)
	}

	// upload
	if err := uploadArchiveFile(ctx, archiveOp, zone, archive.ID, d.Get("archive_file").(string), ftpServer); err != nil {
		return err
	}

	d.SetId(archive.ID.String())
	return resourceSakuraCloudArchiveRead(d, meta)
}

func resourceSakuraCloudArchiveRead(d *schema.ResourceData, meta interface{}) error {
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
	return setArchiveResourceData(d, client, archive)
}

func resourceSakuraCloudArchiveUpdate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	archiveOp := sacloud.NewArchiveOp(client)

	archive, err := archiveOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Archive[%s]: %s", d.Id(), err)
	}

	archive, err = archiveOp.Update(ctx, zone, archive.ID, expandArchiveUpdateRequest(d))
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Archive[%s] is failed: %s", archive.ID, err)
	}

	if isArchiveContentChanged(d) {
		// upload
		if err := uploadArchiveFile(ctx, archiveOp, zone, archive.ID, d.Get("archive_file").(string), nil); err != nil {
			return err
		}
	}

	return resourceSakuraCloudArchiveRead(d, meta)
}

func resourceSakuraCloudArchiveDelete(d *schema.ResourceData, meta interface{}) error {
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

	if err := archiveOp.Delete(ctx, zone, archive.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud Archive[%s] is failed: %s", archive.ID, err)
	}

	d.SetId("")
	return nil
}

func setArchiveResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.Archive) error {
	d.Set("hash", expandArchiveHash(d))
	d.Set("icon_id", data.IconID.String())
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	d.Set("name", data.Name)
	d.Set("size", data.GetSizeGB())
	d.Set("description", data.Description)
	d.Set("zone", getZone(d, client))
	return nil
}

func uploadArchiveFile(ctx context.Context, archiveOp sacloud.ArchiveAPI, zone string, id types.ID, filePath string, ftpServer *sacloud.FTPServer) error {
	path, err := expandHomeDir(filePath)
	if err != nil {
		return fmt.Errorf("preparing SakuraCloud Archive creation is failed: %s", err)
	}

	if ftpServer == nil {
		// open FTPS connections
		fs, err := archiveOp.OpenFTP(ctx, zone, id, &sacloud.OpenFTPRequest{ChangePassword: true})
		if err != nil {
			return fmt.Errorf("opening FTPS connection is failed: %s", err)
		}
		ftpServer = fs
	}

	// upload
	if err := uploadFileViaFTPS(ctx, ftpServer.User, ftpServer.Password, ftpServer.HostName, path); err != nil {
		return fmt.Errorf("uploading file to SakuraCloud is failed: %s", err)
	}

	// close FTPS connection
	if err := archiveOp.CloseFTP(ctx, zone, id); err != nil {
		return fmt.Errorf("closing FTPS Connection is failed: %s", err)
	}

	return nil
}

func isArchiveContentChanged(d *schema.ResourceData) bool {
	contentAttrs := []string{"archive_file", "hash"}
	isContentChanged := false
	for _, attr := range contentAttrs {
		if d.HasChange(attr) {
			isContentChanged = true
			break
		}
	}
	return isContentChanged
}
