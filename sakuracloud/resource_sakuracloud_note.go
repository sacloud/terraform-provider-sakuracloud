package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
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
				Optional: true,
			},
			"class": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(sacloud.NoteClassShell),
				ValidateFunc: validation.StringInSlice([]string{
					string(sacloud.NoteClassShell),
					string(sacloud.NoteClassYAMLCloudConfig),
				}, false),
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSakuraCloudNoteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	opts := client.Note.New()

	opts.Name = d.Get("name").(string)
	opts.Content = d.Get("content").(string)

	if class, ok := d.GetOk("class"); ok {
		s := class.(string)
		if s == "" {
			s = string(sacloud.NoteClassShell)
		}
		opts.SetClassByStr(s)
	}

	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.SetIconByID(toSakuraCloudID(iconID.(string)))
	}
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	if rawTags, ok := d.GetOk("tags"); ok {
		if rawTags != nil {
			opts.Tags = expandTags(client, rawTags.([]interface{}))
		}
	}

	note, err := client.Note.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Note resource: %s", err)
	}

	d.SetId(note.GetStrID())
	return resourceSakuraCloudNoteRead(d, meta)
}

func resourceSakuraCloudNoteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	note, err := client.Note.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud Note resource: %s", err)
	}

	return setNoteResourceData(d, client, note)
}

func resourceSakuraCloudNoteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	note, err := client.Note.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Note resource: %s", err)
	}

	if d.HasChange("name") {
		note.Name = d.Get("name").(string)
	}
	if d.HasChange("class") {
		if class, ok := d.GetOk("class"); ok {
			s := class.(string)
			if s == "" {
				s = string(sacloud.NoteClassShell)
			}
			note.SetClassByStr(s)
		}
	}
	if d.HasChange("content") {
		note.Content = d.Get("content").(string)
	}
	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			note.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			note.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			note.Description = description.(string)
		} else {
			note.Description = ""
		}
	}

	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			note.Tags = expandTags(client, rawTags)
		} else {
			note.Tags = expandTags(client, []interface{}{})
		}
	}

	note, err = client.Note.Update(note.ID, note)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Note resource: %s", err)
	}
	return resourceSakuraCloudNoteRead(d, meta)
}

func resourceSakuraCloudNoteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	_, err := client.Note.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Note resource: %s", err)
	}

	return nil
}

func setNoteResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.Note) error {
	d.Set("name", data.Name)
	d.Set("content", data.Content)
	d.Set("class", data.Class)
	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	d.Set("tags", realTags(client, data.Tags))
	return nil
}
