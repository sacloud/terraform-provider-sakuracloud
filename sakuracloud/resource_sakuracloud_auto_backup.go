package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudAutoBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudAutoBackupCreate,
		Read:   resourceSakuraCloudAutoBackupRead,
		Update: resourceSakuraCloudAutoBackupUpdate,
		Delete: resourceSakuraCloudAutoBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"weekdays": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"max_backup_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 10),
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

func resourceSakuraCloudAutoBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	autoBackupOp := sacloud.NewAutoBackupOp(client)

	if err := validateBackupWeekdays(d, "weekdays"); err != nil {
		return err
	}

	req := &sacloud.AutoBackupCreateRequest{
		Name:                    d.Get("name").(string),
		Description:             d.Get("description").(string),
		Tags:                    expandTagsV2(d.Get("tags").([]interface{})),
		DiskID:                  extractSakuraID(d, "disk_id"),
		MaximumNumberOfArchives: d.Get("max_backup_num").(int),
		BackupSpanWeekdays:      expandBackupWeekdays(d.Get("weekdays").([]interface{})),
		IconID:                  extractSakuraID(d, "icon_id"),
	}

	autoBackup, err := autoBackupOp.Create(ctx, zone, req)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud AutoBackup resource is failed: %s", err)
	}

	d.SetId(autoBackup.ID.String())
	return resourceSakuraCloudAutoBackupRead(d, meta)
}

func resourceSakuraCloudAutoBackupRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	autoBackupOp := sacloud.NewAutoBackupOp(client)

	autoBackup, err := autoBackupOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not find SakuraCloud AutoBackup resource: %s", err)
	}
	return setAutoBackupResourceData(d, client, autoBackup)
}

func setAutoBackupResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.AutoBackup) error {
	d.Set("name", data.Name)
	d.Set("disk_id", data.DiskID.String())
	if err := d.Set("weekdays", flattenBackupWeekdays(data.BackupSpanWeekdays)); err != nil {
		return fmt.Errorf("error setting weekdays: %v", data.BackupSpanWeekdays)
	}
	d.Set("max_backup_num", data.MaximumNumberOfArchives)
	if !data.IconID.IsEmpty() {
		d.Set("icon_id", data.IconID.String())
	}
	d.Set("description", data.Description)
	if err := d.Set("tags", flattenTags(data.Tags)); err != nil {
		return fmt.Errorf("error setting tags: %v", data.Tags)
	}
	d.Set("zone", getV2Zone(d, client))
	return nil
}

func resourceSakuraCloudAutoBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	autoBackupOp := sacloud.NewAutoBackupOp(client)

	autoBackup, err := autoBackupOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud AutoBackup[%s]: %s", d.Id(), err)
	}

	if err := validateBackupWeekdays(d, "weekdays"); err != nil {
		return err
	}

	req := &sacloud.AutoBackupUpdateRequest{
		Name:                    d.Get("name").(string),
		Description:             d.Get("description").(string),
		Tags:                    expandTagsV2(d.Get("tags").([]interface{})),
		MaximumNumberOfArchives: d.Get("max_backup_num").(int),
		BackupSpanWeekdays:      expandBackupWeekdays(d.Get("weekdays").([]interface{})),
		IconID:                  extractSakuraID(d, "icon_id"),
		SettingsHash:            autoBackup.SettingsHash,
	}

	autoBackup, err = autoBackupOp.Update(ctx, zone, autoBackup.ID, req)
	if err != nil {
		return fmt.Errorf("updating SakuraCloud AutoBackup[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudAutoBackupRead(d, meta)
}

func resourceSakuraCloudAutoBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	autoBackupOp := sacloud.NewAutoBackupOp(client)

	autoBackup, err := autoBackupOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud AutoBackup[%s]: %s", d.Id(), err)
	}

	if err := autoBackupOp.Delete(ctx, zone, autoBackup.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud AutoBackup[%s] is failed: %s", autoBackup.ID, err)
	}

	d.SetId("")
	return nil
}
