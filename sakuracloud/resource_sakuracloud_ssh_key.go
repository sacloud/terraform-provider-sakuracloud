package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
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
	client, ctx, _ := getSacloudV2Client(d, meta)
	sshKeyOp := sacloud.NewSSHKeyOp(client)

	key, err := sshKeyOp.Create(ctx, &sacloud.SSHKeyCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		PublicKey:   d.Get("public_key").(string),
	})
	if err != nil {
		return fmt.Errorf("creating SSHKey is failed: %s", err)
	}

	d.SetId(key.ID.String())
	return resourceSakuraCloudSSHKeyRead(d, meta)
}

func resourceSakuraCloudSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	sshKeyOp := sacloud.NewSSHKeyOp(client)

	key, err := sshKeyOp.Read(ctx, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SSHKey: %s", err)
	}
	return setSSHKeyResourceData(ctx, d, client, key)
}

func resourceSakuraCloudSSHKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	sshKeyOp := sacloud.NewSSHKeyOp(client)

	key, err := sshKeyOp.Read(ctx, types.StringID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SSHKey: %s", err)
	}

	_, err = sshKeyOp.Update(ctx, key.ID, &sacloud.SSHKeyUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	})
	if err != nil {
		return fmt.Errorf("updating SSHKey is failed: %s", err)
	}
	return resourceSakuraCloudSSHKeyRead(d, meta)
}

func resourceSakuraCloudSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	sshKeyOp := sacloud.NewSSHKeyOp(client)

	key, err := sshKeyOp.Read(ctx, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SSHKey: %s", err)
	}

	if err := sshKeyOp.Delete(ctx, key.ID); err != nil {
		return fmt.Errorf("deleting SSHKey is failed: %s", err)
	}
	return nil
}

func setSSHKeyResourceData(_ context.Context, d *schema.ResourceData, _ *APIClient, data *sacloud.SSHKey) error {
	d.Set("name", data.Name)
	d.Set("public_key", data.PublicKey)
	d.Set("fingerprint", data.Fingerprint)
	d.Set("description", data.Description)
	return nil
}
