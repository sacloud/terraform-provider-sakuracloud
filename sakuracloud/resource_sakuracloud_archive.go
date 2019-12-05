package sakuracloud

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mitchellh/go-homedir"
	"github.com/sacloud/ftps"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudArchive() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudArchiveCreate,
		Read:   resourceSakuraCloudArchiveRead,
		Update: resourceSakuraCloudArchiveUpdate,
		Delete: resourceSakuraCloudArchiveDelete,
		CustomizeDiff: customdiff.ComputedIf("hash", func(d *schema.ResourceDiff, meta interface{}) bool {
			return d.HasChange("archive_file")
		}),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Computed: true,
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
	client, ctx, zone := getSacloudV2Client(d, meta)
	archiveOp := sacloud.NewArchiveOp(client)

	// prepare create parameters
	req := &sacloud.ArchiveCreateBlankRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		SizeMB:      toSizeMB(d.Get("size").(int)),
		IconID:      types.StringID(d.Get("icon_id").(string)),
		Tags:        expandTagsV2(d.Get("tags").([]interface{})),
	}

	source := d.Get("archive_file").(string)
	path, err := expandHomeDir(source)
	if err != nil {
		return fmt.Errorf("prepare creating SakuraCloud Archive resource is failed: %s", err)
	}

	// create
	archive, ftpServer, err := archiveOp.CreateBlank(ctx, zone, req)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Archive resource is failed: %s", err)
	}

	// upload
	ftpClient := ftps.NewClient(ftpServer.User, ftpServer.Password, ftpServer.HostName)
	if err := ftpClient.Upload(path); err != nil {
		return fmt.Errorf("upload archive_file to SakuraCloud is failed: %s", err)
	}

	// close FTPS connection
	if err := archiveOp.CloseFTP(ctx, zone, archive.ID); err != nil {
		return fmt.Errorf("closing FTPS Connection is failed: %s", err)

	}

	d.SetId(archive.ID.String())
	return resourceSakuraCloudArchiveRead(d, meta)
}

func resourceSakuraCloudArchiveRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	archiveOp := sacloud.NewArchiveOp(client)

	archive, err := archiveOp.Read(ctx, zone, types.StringID(d.Id()))
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
	client, ctx, zone := getSacloudV2Client(d, meta)
	archiveOp := sacloud.NewArchiveOp(client)

	archive, err := archiveOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Archive[%s]: %s", d.Id(), err)
	}

	req := &sacloud.ArchiveUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTagsV2(d.Get("tags").([]interface{})),
		IconID:      types.StringID(d.Get("icon_id").(string)),
	}

	archive, err = archiveOp.Update(ctx, zone, archive.ID, req)
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Archive[%s] is failed: %s", d.Id(), err)
	}

	contentAttrs := []string{"archive_file", "hash"}
	isContentChanged := false
	for _, attr := range contentAttrs {
		if d.HasChange(attr) {
			isContentChanged = true
			break
		}
	}
	if isContentChanged {
		source := d.Get("archive_file").(string)
		path, err := expandHomeDir(source)
		if err != nil {
			return fmt.Errorf("prepare upload SakuraCloud Archive[%s] is failed: %s", archive.ID, err)
		}

		// open FTPS connections
		ftpServer, err := archiveOp.OpenFTP(ctx, zone, archive.ID, &sacloud.OpenFTPRequest{ChangePassword: true})
		if err != nil {
			return fmt.Errorf("Failed to Open FTPS Connection: %s", err)
		}

		// upload
		ftpClient := ftps.NewClient(ftpServer.User, ftpServer.Password, ftpServer.HostName)
		if err := ftpClient.Upload(path); err != nil {
			return fmt.Errorf("upload archive_file to SakuraCloud is failed: %s", err)
		}

		// close FTPS connection
		if err := archiveOp.CloseFTP(ctx, zone, archive.ID); err != nil {
			return fmt.Errorf("closing FTPS Connection is failed: %s", err)
		}
	}

	return resourceSakuraCloudArchiveRead(d, meta)
}

func resourceSakuraCloudArchiveDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	archiveOp := sacloud.NewArchiveOp(client)

	archive, err := archiveOp.Read(ctx, zone, types.StringID(d.Id()))
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
	// NOTE 本来はAPIにてmd5ハッシュを取得できるのが望ましい。現状ではここでファイルを読んで算出する。
	if v, ok := d.GetOk("archive_file"); ok {
		source := v.(string)

		path, err := expandHomeDir(source)
		if err != nil {
			return err
		}
		hash, err := md5CheckSumFromFile(path)
		if err != nil {
			return err
		}
		d.Set("hash", hash)
	}

	if !data.IconID.IsEmpty() {
		d.Set("icon_id", data.IconID.String())
	}
	if err := d.Set("tags", flattenTags(data.Tags)); err != nil {
		return fmt.Errorf("error setting tags: %v", data.Tags)
	}
	d.Set("name", data.Name)
	d.Set("size", data.GetSizeGB())
	d.Set("description", data.Description)
	d.Set("zone", getV2Zone(d, client))
	return nil
}

func expandHomeDir(path string) (string, error) {
	expanded, err := homedir.Expand(path)
	if err != nil {
		return "", fmt.Errorf("expanding homedir in path[%s] is failed: %s", expanded, err)
	}
	// file exists?
	if _, err := os.Stat(expanded); err != nil {
		return "", fmt.Errorf("opening archive_file[%s] is failed: %s", expanded, err)
	}
	return expanded, nil
}

func md5CheckSumFromFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("opening archive_file[%s] is failed: %s", path, err)
	}
	defer f.Close() // nolint

	b := base64.NewEncoder(base64.StdEncoding, f)
	defer b.Close() // nolint

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, f); err != nil {
		return "", fmt.Errorf("encoding to base64 from archive_file[%s] is failed: %s", path, err)
	}

	h := md5.New()
	if _, err := io.Copy(h, &buf); err != nil {
		return "", fmt.Errorf("calculating md5 from archive_file[%s] is failed: %s", path, err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
