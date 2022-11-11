// Copyright 2016-2022 terraform-provider-sakuracloud authors
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
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/cleanup"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iso9660wrap"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/ftps"
)

func resourceSakuraCloudCDROM() *schema.Resource {
	resourceName := "CD-ROM"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudCDROMCreate,
		ReadContext:   resourceSakuraCloudCDROMRead,
		UpdateContext: resourceSakuraCloudCDROMUpdate,
		DeleteContext: resourceSakuraCloudCDROMDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: customdiff.ComputedIf("hash", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
			return d.HasChange("iso_image_file") || d.HasChange("content")
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(24 * time.Hour),
			Update: schema.DefaultTimeout(24 * time.Hour),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"size": schemaResourceSize(resourceName, 5, []int{5, 10}...),
			"iso_image_file": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content"},
				Description: desc.Sprintf(
					"The file path to upload to as the CD-ROM. %s",
					desc.Conflicts("content"),
				),
			},
			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"iso_image_file"},
				Description: desc.Sprintf(
					"The content to upload to as the CD-ROM. %s",
					desc.Conflicts("iso_image_file"),
				),
			},
			"content_file_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       cdromDefaultISOLabel,
				ConflictsWith: []string{"iso_image_file"},
				Description: desc.Sprintf(
					"The name of content file to upload to as the CD-ROM. This is only used when `content` is specified. %s",
					desc.Conflicts("iso_image_file"),
				),
			},
			"hash": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The md5 checksum calculated from the base64 encoded file body",
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"zone":        schemaResourceZone(resourceName),
		},
	}
}

const (
	cdromDefaultISOLabel = "config"
)

func resourceSakuraCloudCDROMCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	cdromOp := iaas.NewCDROMOp(client)

	cdrom, ftpServer, err := cdromOp.Create(ctx, zone, expandCDROMCreateRequest(d))
	if err != nil {
		return diag.Errorf("creating SakuraCloud CDROM is failed: %s", err)
	}

	// upload
	err = uploadCDROMFile(&uploadCDROMContext{
		Context:   ctx,
		zone:      zone,
		id:        cdrom.ID,
		client:    client,
		ftpServer: ftpServer,
	}, d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cdrom.ID.String())
	return resourceSakuraCloudCDROMRead(ctx, d, meta)
}

func resourceSakuraCloudCDROMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	cdromOp := iaas.NewCDROMOp(client)

	cdrom, err := cdromOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud CDROM[%s]: %s", d.Id(), err)
	}
	return setCDROMResourceData(ctx, d, client, cdrom)
}

func resourceSakuraCloudCDROMUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	cdromOp := iaas.NewCDROMOp(client)

	cdrom, err := cdromOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud CDROM[%s]: %s", d.Id(), err)
	}

	cdrom, err = cdromOp.Update(ctx, zone, cdrom.ID, expandCDROMUpdateRequest(d))
	if err != nil {
		return diag.Errorf("updating SakuraCloud CDROM[%s] is failed: %s", d.Id(), err)
	}

	if isCDROMContentChanged(d) {
		err = uploadCDROMFile(&uploadCDROMContext{
			Context:   ctx,
			zone:      zone,
			id:        cdrom.ID,
			client:    client,
			ftpServer: nil,
		}, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSakuraCloudCDROMRead(ctx, d, meta)
}

func resourceSakuraCloudCDROMDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	cdromOp := iaas.NewCDROMOp(client)

	cdrom, err := cdromOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud CDROM[%s]: %s", d.Id(), err)
	}

	if err := cleanup.DeleteCDROM(ctx, client, zone, cdrom.ID, client.checkReferencedOption()); err != nil {
		return diag.Errorf("deleting SakuraCloud CDROM[%s] is failed: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func setCDROMResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *iaas.CDROM) diag.Diagnostics {
	d.Set("hash", expandCDROMContentHash(d)) // nolint
	d.Set("name", data.Name)                 // nolint
	d.Set("size", data.GetSizeGB())          // nolint
	d.Set("icon_id", data.IconID.String())   // nolint
	d.Set("description", data.Description)   // nolint
	d.Set("zone", getZone(d, client))        // nolint
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}

type uploadCDROMContext struct {
	context.Context
	zone      string
	id        types.ID
	client    *APIClient
	ftpServer *iaas.FTPServer
}

func uploadCDROMFile(ctx *uploadCDROMContext, d *schema.ResourceData) error {
	cdromOp := iaas.NewCDROMOp(ctx.client)

	filePath, isTemporal, err := prepareContentFile(d)
	if isTemporal {
		defer os.Remove(filePath)
	}
	if err != nil {
		return err
	}

	// eject
	ejectedServerIDs, err := ejectCDROMFromAllServers(ctx, d, ctx.client, ctx.id)
	if err != nil {
		return fmt.Errorf("could not eject CDROM[%s] from Server: %s", ctx.id, err)
	}

	ftpServer := ctx.ftpServer
	if ftpServer == nil {
		fs, err := cdromOp.OpenFTP(ctx, ctx.zone, ctx.id, &iaas.OpenFTPRequest{ChangePassword: false})
		if err != nil {
			return fmt.Errorf("opening FTPS connection to CDROM[%s] is failed: %s", ctx.id, err)
		}
		ftpServer = fs
	}

	// upload
	if err := ftps.UploadFile(ctx, ftpServer.User, ftpServer.Password, ftpServer.HostName, filePath); err != nil {
		return fmt.Errorf("upload CD-ROM contents is failed: %s", err)
	}

	// close
	if err := cdromOp.CloseFTP(ctx, ctx.zone, ctx.id); err != nil {
		return fmt.Errorf("closing FTPS Connection is failed: %s", err)
	}

	// re-insert CDROM
	if err := insertCDROMToAllServers(ctx, d, ctx.client, ctx.id, ejectedServerIDs); err != nil {
		return fmt.Errorf("could not insert CDROM[%s] to Server: %s", ctx.id, err)
	}

	return nil
}

func isCDROMContentChanged(d *schema.ResourceData) bool {
	contentAttrs := []string{"iso_image_file", "content", "content_file_name", "hash"}
	isContentChanged := false
	for _, attr := range contentAttrs {
		if d.HasChange(attr) {
			isContentChanged = true
			break
		}
	}
	return isContentChanged
}

func prepareContentFile(d *schema.ResourceData) (string, bool, error) {
	isTemporal := false
	var filePath string

	if v, ok := d.GetOk("iso_image_file"); ok {
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return "", false, fmt.Errorf("error expanding homedir in source (%s): %s", source, err)
		}
		// file exists?
		if _, err := os.Stat(path); err != nil {
			return "", false, fmt.Errorf("error opening iso_image_file (%s): %s", source, err)
		}
		filePath = path
	} else if v, ok := d.GetOk("content"); ok {
		isTemporal = true
		content := v.(string)
		label := cdromDefaultISOLabel
		if l, ok := d.GetOk("content_file_name"); ok {
			s := l.(string)
			if s != "" {
				label = s
			}
		}

		// create iso9660 format file
		tmpFile, err := os.CreateTemp("", "tf-sakuracloud-cdrom")
		if err != nil {
			return "", isTemporal, fmt.Errorf("error creating temp-file : %s", err)
		}
		defer tmpFile.Close() // nolint
		filePath = tmpFile.Name()
		err = writeISOFile(filePath, []byte(content), label)
		if err != nil {
			return "", isTemporal, fmt.Errorf("error writing temp-file : %s", err)
		}
	} else {
		return "", isTemporal, fmt.Errorf("must specify \"iso_image_file\" or \"content\" field")
	}
	return filePath, isTemporal, nil
}

func writeISOFile(path string, content []byte, label string) error {
	outfh, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outfh.Close() // nolint

	return iso9660wrap.WriteBuffer(outfh, content, label)
}

func ejectCDROMFromAllServers(ctx context.Context, d *schema.ResourceData, client *APIClient, cdromID types.ID) ([]types.ID, error) {
	serverOp := iaas.NewServerOp(client)
	zone := getZone(d, client)
	searched, err := serverOp.Find(ctx, zone, &iaas.FindCondition{})
	if err != nil {
		return nil, err
	}
	var ejectedIDs []types.ID
	for _, server := range searched.Servers {
		if server.CDROMID == cdromID {
			if err := serverOp.EjectCDROM(ctx, zone, server.ID, &iaas.EjectCDROMRequest{ID: cdromID}); err != nil {
				return nil, err
			}
			ejectedIDs = append(ejectedIDs, server.ID)
		}
	}
	return ejectedIDs, nil
}

func insertCDROMToAllServers(ctx context.Context, d *schema.ResourceData, client *APIClient, cdromID types.ID, serverIDs []types.ID) error {
	serverOp := iaas.NewServerOp(client)
	zone := getZone(d, client)

	for _, id := range serverIDs {
		if err := serverOp.InsertCDROM(ctx, zone, id, &iaas.InsertCDROMRequest{ID: cdromID}); err != nil {
			return err
		}
	}
	return nil
}
