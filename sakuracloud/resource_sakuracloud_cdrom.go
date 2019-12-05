package sakuracloud

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/terraform/helper/customdiff"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/go-homedir"
	"github.com/sacloud/ftps"
	"github.com/sacloud/iso9660wrap"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudCDROM() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudCDROMCreate,
		Read:   resourceSakuraCloudCDROMRead,
		Update: resourceSakuraCloudCDROMUpdate,
		Delete: resourceSakuraCloudCDROMDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: customdiff.ComputedIf("hash", func(d *schema.ResourceDiff, meta interface{}) bool {
			return d.HasChange("iso_image_file") || d.HasChange("content")
		}),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  5,
			},
			"iso_image_file": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content"},
			},
			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"iso_image_file"},
			},
			"content_file_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       cdromDefaultISOLabel,
				ConflictsWith: []string{"iso_image_file"},
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
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		},
	}
}

const (
	cdromDefaultISOLabel = "config"
)

func resourceSakuraCloudCDROMCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	cdromOp := sacloud.NewCDROMOp(client)

	// prepare create parameters
	req := &sacloud.CDROMCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		SizeMB:      toSizeMB(d.Get("size").(int)),
		IconID:      types.StringID(d.Get("icon_id").(string)),
		Tags:        expandTagsV2(d.Get("tags").([]interface{})),
	}

	filePath, isTemporal, err := prepareContentFile(d)
	if isTemporal {
		defer os.Remove(filePath)
	}
	if err != nil {
		return err
	}

	cdrom, ftpServer, err := cdromOp.Create(ctx, zone, req)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud CDROM resource is failed: %s", err)
	}

	// upload
	ftpClient := ftps.NewClient(ftpServer.User, ftpServer.Password, ftpServer.HostName)
	if err := ftpClient.Upload(filePath); err != nil {
		return fmt.Errorf("upload CD-ROM contents to SakuraCloud is failed: %s", err)
	}

	// close
	if err := cdromOp.CloseFTP(ctx, zone, cdrom.ID); err != nil {
		return fmt.Errorf("closing FTPS Connection is failed: %s", err)
	}

	d.SetId(cdrom.ID.String())
	return resourceSakuraCloudCDROMRead(d, meta)
}

func resourceSakuraCloudCDROMRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	cdromOp := sacloud.NewCDROMOp(client)

	cdrom, err := cdromOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud CDROM[%s]: %s", d.Id(), err)
	}
	return setCDROMResourceData(ctx, d, client, cdrom)
}

func resourceSakuraCloudCDROMUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	cdromOp := sacloud.NewCDROMOp(client)

	cdrom, err := cdromOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud CDROM[%s]: %s", d.Id(), err)
	}

	req := &sacloud.CDROMUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTagsV2(d.Get("tags").([]interface{})),
		IconID:      types.StringID(d.Get("icon_id").(string)),
	}

	cdrom, err = cdromOp.Update(ctx, zone, cdrom.ID, req)
	if err != nil {
		return fmt.Errorf("updating SakuraCloud CDROM[%s] is failed: %s", d.Id(), err)
	}

	contentAttrs := []string{"iso_image_file", "content", "content_file_name", "hash"}
	isContentChanged := false
	for _, attr := range contentAttrs {
		if d.HasChange(attr) {
			isContentChanged = true
			break
		}
	}
	if isContentChanged {
		// eject
		ejectedServerIDs, err := ejectCDROMFromAllServers(ctx, d, client, cdrom.ID)
		if err != nil {
			return fmt.Errorf("could not eject CDROM[%s] from Server: %s", cdrom.ID, err)
		}

		filePath, isTemporal, err := prepareContentFile(d)
		if isTemporal {
			defer os.Remove(filePath) // nolint
		}
		if err != nil {
			return err
		}
		ftpServer, err := cdromOp.OpenFTP(ctx, zone, cdrom.ID, &sacloud.OpenFTPRequest{ChangePassword: false})
		if err != nil {
			return fmt.Errorf("opening FTPS connection to CDROM[%s] is failed: %s", cdrom.ID, err)
		}
		// upload
		ftpClient := ftps.NewClient(ftpServer.User, ftpServer.Password, ftpServer.HostName)
		if err := ftpClient.Upload(filePath); err != nil {
			return fmt.Errorf("upload CD-ROM contents to SakuraCloud is failed: %s", err)
		}
		// close
		if err := cdromOp.CloseFTP(ctx, zone, cdrom.ID); err != nil {
			return fmt.Errorf("closing FTPS Connection is failed: %s", err)
		}

		// re-insert CDROM
		if err := insertCDROMToAllServers(ctx, d, client, cdrom.ID, ejectedServerIDs); err != nil {
			return fmt.Errorf("could not insert CDROM[%s] to Server: %s", cdrom.ID, err)
		}
	}

	return resourceSakuraCloudCDROMRead(d, meta)
}

func resourceSakuraCloudCDROMDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	cdromOp := sacloud.NewCDROMOp(client)

	cdrom, err := cdromOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud CDROM[%s]: %s", d.Id(), err)
	}

	// eject
	if _, err := ejectCDROMFromAllServers(ctx, d, client, cdrom.ID); err != nil {
		return fmt.Errorf("could not eject CDROM[%s] from Server: %s", cdrom.ID, err)
	}

	if err := cdromOp.Delete(ctx, zone, cdrom.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud CDROM[%s] is failed: %s", cdrom.ID, err)
	}

	d.SetId("")
	return nil
}

func setCDROMResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.CDROM) error {
	// NOTE 本来はAPIにてmd5ハッシュを取得できるのが望ましい。現状ではここでファイルを読んで算出する。
	if v, ok := d.GetOk("iso_image_file"); ok {
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

	d.Set("name", data.Name)
	d.Set("size", data.GetSizeGB())
	if !data.IconID.IsEmpty() {
		d.Set("icon_id", data.IconID.String())
	}
	d.Set("description", data.Description)
	if err := d.Set("tags", flattenTags(data.Tags)); err != nil {
		return fmt.Errorf("error settings tags: %v", data.Tags)
	}
	d.Set("zone", getV2Zone(d, client))
	return nil
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
		tmpFile, err := ioutil.TempFile("", "tf-sakuracloud-cdrom")
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
	serverOp := sacloud.NewServerOp(client)
	zone := getV2Zone(d, client)
	searched, err := serverOp.Find(ctx, zone, &sacloud.FindCondition{Count: defaultSearchLimit})
	if err != nil {
		return nil, err
	}
	var ejectedIDs []types.ID
	for _, server := range searched.Servers {
		if server.CDROMID == cdromID {
			if err := serverOp.EjectCDROM(ctx, zone, server.ID, &sacloud.EjectCDROMRequest{ID: cdromID}); err != nil {
				return nil, err
			}
			ejectedIDs = append(ejectedIDs, server.ID)
		}
	}
	return ejectedIDs, nil
}

func insertCDROMToAllServers(ctx context.Context, d *schema.ResourceData, client *APIClient, cdromID types.ID, serverIDs []types.ID) error {
	serverOp := sacloud.NewServerOp(client)
	zone := getV2Zone(d, client)

	for _, id := range serverIDs {
		if err := serverOp.InsertCDROM(ctx, zone, id, &sacloud.InsertCDROMRequest{ID: cdromID}); err != nil {
			return err
		}
	}
	return nil
}
