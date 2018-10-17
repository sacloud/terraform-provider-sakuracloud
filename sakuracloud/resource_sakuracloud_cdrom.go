package sakuracloud

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/go-homedir"
	"github.com/sacloud/ftps"
	"github.com/sacloud/iso9660wrap"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
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
		CustomizeDiff: hasTagResourceCustomizeDiff,
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

var (
	cdromDefaultISOLabel = "config"
)

func resourceSakuraCloudCDROMCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	opts := client.CDROM.New()

	opts.Name = d.Get("name").(string)

	filePath, isTemporal, err := prepareContentFile(d)
	if isTemporal {
		defer os.Remove(filePath)
	}
	if err != nil {
		return err
	}

	opts.SizeMB = toSizeMB(d.Get("size").(int))
	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.SetIconByID(toSakuraCloudID(iconID.(string)))
	}
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		opts.Tags = expandTags(client, rawTags)
	}

	cdrom, ftpServer, err := client.CDROM.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud CDROM resource: %s", err)
	}

	// upload
	ftpClient := ftps.NewClient(ftpServer.User, ftpServer.Password, ftpServer.HostName)
	if err := ftpClient.Upload(filePath); err != nil {
		return fmt.Errorf("Failed to upload SakuraCloud CDROM resource: %s", err)
	}

	// close
	if _, err := client.CDROM.CloseFTP(cdrom.ID); err != nil {
		return fmt.Errorf("Failed to Close FTPS Connection from CDROM resource: %s", err)

	}

	d.SetId(cdrom.GetStrID())
	return resourceSakuraCloudCDROMRead(d, meta)
}

func resourceSakuraCloudCDROMRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	cdrom, err := client.CDROM.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud CDROM resource: %s", err)
	}

	return setCDROMResourceData(d, client, cdrom)
}

func resourceSakuraCloudCDROMUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	cdrom, err := client.CDROM.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud CDROM resource: %s", err)
	}
	if d.HasChange("name") {
		cdrom.Name = d.Get("name").(string)
	}
	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			cdrom.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			cdrom.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			cdrom.Description = description.(string)
		} else {
			cdrom.Description = ""
		}
	}
	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			cdrom.Tags = expandTags(client, rawTags)
		} else {
			cdrom.Tags = expandTags(client, []interface{}{})
		}
	}
	cdrom, err = client.CDROM.Update(cdrom.ID, cdrom)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud CDROM resource: %s", err)
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
		allServers, err := client.GetServerAPI().Find()
		if err != nil {
			return fmt.Errorf("Couldn't find SakuraCloud server resource: %s", err)
		}
		var serverIDs []int64
		for _, s := range allServers.Servers {
			if s.Instance.CDROM != nil && s.Instance.CDROM.ID == cdrom.ID {
				// eject
				if _, err := client.GetServerAPI().EjectCDROM(s.ID, cdrom.ID); err != nil {
					return fmt.Errorf("Couldn't eject CDROM from Server: %s", err)
				}
				serverIDs = append(serverIDs, s.ID)
			}
		}

		filePath, isTemporal, err := prepareContentFile(d)
		if isTemporal {
			defer os.Remove(filePath)
		}
		if err != nil {
			return err
		}
		ftpServer, err := client.CDROM.OpenFTP(cdrom.ID, false)
		if err != nil {
			return fmt.Errorf("Failed to Open FTPS Connection to CDROM resource: %s", err)
		}
		// upload
		ftpClient := ftps.NewClient(ftpServer.User, ftpServer.Password, ftpServer.HostName)
		if err := ftpClient.Upload(filePath); err != nil {
			return fmt.Errorf("Failed to upload SakuraCloud CDROM resource: %s", err)
		}
		// close
		if _, err := client.CDROM.CloseFTP(cdrom.ID); err != nil {
			return fmt.Errorf("Failed to Close FTPS Connection from CDROM resource: %s", err)

		}

		// re-insert CDROM
		for _, serverID := range serverIDs {
			if _, err := client.GetServerAPI().InsertCDROM(serverID, cdrom.ID); err != nil {
				return fmt.Errorf("Couldn't insert CDROM from Server: %s", err)
			}
		}
	}

	return resourceSakuraCloudCDROMRead(d, meta)
}

func resourceSakuraCloudCDROMDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	cdrom, err := client.CDROM.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud CDROM resource: %s", err)
	}

	// eject
	allServers, err := client.GetServerAPI().Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud server resource: %s", err)
	}
	for _, s := range allServers.Servers {
		if s.Instance.CDROM != nil && s.Instance.CDROM.ID == cdrom.ID {
			// eject
			if _, err := client.GetServerAPI().EjectCDROM(s.ID, cdrom.ID); err != nil {
				return fmt.Errorf("Couldn't eject CDROM from Server: %s", err)
			}
		}
	}

	_, err = client.CDROM.Delete(toSakuraCloudID(d.Id()))

	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud CDROM resource: %s", err)
	}

	return nil
}

func setCDROMResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.CDROM) error {

	d.Set("name", data.Name)
	d.Set("size", toSizeGB(data.SizeMB))
	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	d.Set("tags", data.Tags)

	// NOTE 本来はAPIにてmd5ハッシュを取得できるのが望ましい
	if v, ok := d.GetOk("iso_image_file"); ok {
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return fmt.Errorf("Error expanding homedir in source (%s): %s", source, err)
		}
		// file exists?
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("Error opening iso_image_file (%s): %s", source, err)
		}

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("Error opening iso_image_file (%s): %s", source, err)
		}
		defer f.Close()
		h := md5.New()
		if _, err := io.Copy(h, f); err != nil {
			return fmt.Errorf("Error calculate md5 from iso_image_file (%s): %s", source, err)
		}

		d.Set("hash", h.Sum(nil))
	}

	d.Set("zone", client.Zone)
	return nil
}

func prepareContentFile(d *schema.ResourceData) (string, bool, error) {
	isTemporal := false
	var filePath string

	if v, ok := d.GetOk("iso_image_file"); ok {
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return "", isTemporal, fmt.Errorf("Error expanding homedir in source (%s): %s", source, err)
		}
		// file exists?
		if _, err := os.Stat(path); err != nil {
			return "", isTemporal, fmt.Errorf("Error opening iso_image_file (%s): %s", source, err)
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
			return "", isTemporal, fmt.Errorf("Error creating temp-file : %s", err)
		}
		defer tmpFile.Close()
		filePath = tmpFile.Name()
		err = writeISOFile(filePath, []byte(content), label)
		if err != nil {
			return "", isTemporal, fmt.Errorf("Error writing temp-file : %s", err)
		}

	} else {
		return "", isTemporal, fmt.Errorf("Must specify \"iso_image_file\" or \"content\" field")
	}
	return filePath, isTemporal, nil
}

func writeISOFile(path string, content []byte, label string) error {
	outfh, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outfh.Close()

	return iso9660wrap.WriteBuffer(outfh, content, label)
}
