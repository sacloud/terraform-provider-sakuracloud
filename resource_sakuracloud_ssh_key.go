package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
)

func resourceSakuraCloudSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudSSHKeyCreate,
		Read:   resourceSakuraCloudSSHKeyRead,
		Update: resourceSakuraCloudSSHKeyUpdate,
		Delete: resourceSakuraCloudSSHKeyDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudSSHKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

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

	d.SetId(key.ID)
	return resourceSakuraCloudSSHKeyRead(d, meta)
}

func resourceSakuraCloudSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	key, err := client.SSHKey.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud SSHKey resource: %s", err)
	}

	d.Set("name", key.Name)
	d.Set("public_key", key.PublicKey)
	d.Set("fingerprint", key.Fingerprint)
	d.Set("description", key.Description)

	return nil
}

func resourceSakuraCloudSSHKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	key, err := client.SSHKey.Read(d.Id())
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
	d.SetId(key.ID)

	return resourceSakuraCloudSSHKeyRead(d, meta)
}

func resourceSakuraCloudSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	_, err := client.SSHKey.Delete(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud SSHKey resource: %s", err)
	}

	return nil
}
