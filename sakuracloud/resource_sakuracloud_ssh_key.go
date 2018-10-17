package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudSSHKeyCreate,
		Read:   resourceSakuraCloudSSHKeyRead,
		Update: resourceSakuraCloudSSHKeyUpdate,
		Delete: resourceSakuraCloudSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"public_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudSSHKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	opts := client.SSHKey.New()

	opts.Name = d.Get("name").(string)
	opts.PublicKey = d.Get("public_key").(string)
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}

	key, err := client.SSHKey.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud SSHKey resource: %s", err)
	}

	d.SetId(key.GetStrID())
	return resourceSakuraCloudSSHKeyRead(d, meta)
}

func resourceSakuraCloudSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	key, err := client.SSHKey.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud SSHKey resource: %s", err)
	}

	return setSSHKeyResourceData(d, client, key)
}

func resourceSakuraCloudSSHKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	key, err := client.SSHKey.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud SSHKey resource: %s", err)
	}

	if d.HasChange("name") {
		key.Name = d.Get("name").(string)
	}
	if d.HasChange("public_key") {
		key.Name = d.Get("public_key").(string)
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			key.Description = description.(string)
		} else {
			key.Description = ""
		}
	}

	key, err = client.SSHKey.Update(key.ID, key)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud SSHKey resource: %s", err)
	}
	return resourceSakuraCloudSSHKeyRead(d, meta)
}

func resourceSakuraCloudSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	_, err := client.SSHKey.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud SSHKey resource: %s", err)
	}

	return nil
}

func setSSHKeyResourceData(d *schema.ResourceData, _ *APIClient, data *sacloud.SSHKey) error {
	d.Set("name", data.Name)
	d.Set("public_key", data.PublicKey)
	d.Set("fingerprint", data.Fingerprint)
	d.Set("description", data.Description)
	return nil
}
