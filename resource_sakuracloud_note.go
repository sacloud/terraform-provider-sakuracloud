package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
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

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSakuraCloudNoteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	opts := client.Note.New()

	opts.Name = d.Get("name").(string)
	opts.Content = d.Get("content").(string)
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	if rawTags, ok := d.GetOk("tags"); ok {
		if rawTags != nil {
			opts.Tags = expandStringList(rawTags.([]interface{}))
		}
	}

	note, err := client.Note.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Note resource: %s", err)
	}

	d.SetId(note.ID)
	return resourceSakuraCloudNoteRead(d, meta)
}

func resourceSakuraCloudNoteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	note, err := client.Note.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Note resource: %s", err)
	}

	d.Set("name", note.Name)
	d.Set("content", note.Content)
	d.Set("description", note.Description)
	d.Set("tags", note.Tags)

	return nil
}

func resourceSakuraCloudNoteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	note, err := client.Note.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Note resource: %s", err)
	}

	if d.HasChange("name") {
		note.Name = d.Get("name").(string)
	}
	if d.HasChange("content") {
		note.Content = d.Get("content").(string)
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
			note.Tags = expandStringList(rawTags)
		} else {
			note.Tags = []string{}
		}
	}

	note, err = client.Note.Update(note.ID, note)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Note resource: %s", err)
	}
	d.SetId(note.ID)

	return resourceSakuraCloudNoteRead(d, meta)
}

func resourceSakuraCloudNoteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	_, err := client.Note.Delete(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Note resource: %s", err)
	}

	return nil
}
