package sakuracloud

import (
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/go-homedir"

	"github.com/sacloud/libsacloud/sacloud"
	"io/ioutil"
	"os"
)

func resourceSakuraCloudIcon() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudIconCreate,
		Read:   resourceSakuraCloudIconRead,
		Update: resourceSakuraCloudIconUpdate,
		Delete: resourceSakuraCloudIconDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"base64content"},
				ForceNew:      true,
			},
			"base64content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
				ForceNew:      true,
			},
			"body": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudIconCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	opts := client.Icon.New()

	opts.Name = d.Get("name").(string)

	var body string
	if v, ok := d.GetOk("source"); ok {
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return fmt.Errorf("Error expanding homedir in source (%s): %s", source, err)
		}
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("Error opening SakuraCloud Icon source (%s): %s", source, err)
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return fmt.Errorf("Error opening SakuraCloud Icon source : %s", err)
		}
		body = base64.StdEncoding.EncodeToString(data)
	} else if v, ok := d.GetOk("base64content"); ok {
		body = v.(string)
	} else {
		return fmt.Errorf("Must specify \"source\" or \"base64content\" field")
	}
	opts.SetImage(body)

	if rawTags, ok := d.GetOk("tags"); ok {
		if rawTags != nil {
			opts.Tags = expandTags(client, rawTags.([]interface{}))
		}
	}

	icon, err := client.Icon.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Icon resource: %s", err)
	}

	d.SetId(icon.GetStrID())
	return resourceSakuraCloudIconRead(d, meta)
}

func resourceSakuraCloudIconRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	icon, err := client.Icon.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Icon resource: %s", err)
	}

	return setIconResourceData(d, client, icon)
}

func resourceSakuraCloudIconUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	icon, err := client.Icon.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Icon resource: %s", err)
	}

	if d.HasChange("name") {
		icon.Name = d.Get("name").(string)
	}

	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			icon.Tags = expandTags(client, rawTags)
		} else {
			icon.Tags = expandTags(client, []interface{}{})
		}
	}

	icon, err = client.Icon.Update(icon.ID, icon)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Icon resource: %s", err)
	}
	d.SetId(icon.GetStrID())

	return resourceSakuraCloudIconRead(d, meta)
}

func resourceSakuraCloudIconDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	_, err := client.Icon.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Icon resource: %s", err)
	}

	return nil
}

func setIconResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.Icon) error {

	d.Set("name", data.Name)

	body, err := client.Icon.GetImage(data.ID, "small")
	if err != nil {
		return fmt.Errorf("Error reading SakuraCloud Icon Resource: %s", err)
	}

	d.Set("body", body)
	d.Set("tags", realTags(client, data.Tags))
	d.Set("url", data.URL)

	d.SetId(data.GetStrID())
	return nil
}
