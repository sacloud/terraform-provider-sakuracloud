package sakuracloud

import (
	"fmt"

	"github.com/sacloud/libsacloud/v2/sacloud/types"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudSSHKeyGen() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudSSHKeyGenCreate,
		Read:   resourceSakuraCloudSSHKeyGenRead,
		Delete: resourceSakuraCloudSSHKeyGenDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
			},
			"pass_phrase": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(8, 64),
			},
			"private_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudSSHKeyGenCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	sshKeyOp := sacloud.NewSSHKeyOp(client)

	key, err := sshKeyOp.Generate(ctx, &sacloud.SSHKeyGenerateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		PassPhrase:  d.Get("pass_phrase").(string),
	})
	if err != nil {
		return fmt.Errorf("generating SSHKey is failed: %s", err)
	}

	d.SetId(key.ID.String())

	// Note: CreateのレスポンスにのみPrivateKeyが含まれる
	return setSSHKeyGenResourceData(d, client, key)
}

func resourceSakuraCloudSSHKeyGenRead(d *schema.ResourceData, meta interface{}) error {
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

	return setSSHKeyGenResourceData(d, client, key)
}

func resourceSakuraCloudSSHKeyGenDelete(d *schema.ResourceData, meta interface{}) error {
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

func setSSHKeyGenResourceData(d *schema.ResourceData, _ *APIClient, data interface{}) error {

	if key, ok := data.(sshKeyType); ok {
		d.Set("name", key.GetName())
		d.Set("public_key", key.GetPublicKey())
		d.Set("fingerprint", key.GetFingerprint())
		d.Set("description", key.GetDescription())

		// has private key?
		if pKey, ok := data.(sshKeyGenType); ok {
			d.Set("private_key", pKey.GetPrivateKey())
		}
	}

	return nil
}

type sshKeyType interface {
	GetName() string
	GetPublicKey() string
	GetFingerprint() string
	GetDescription() string
}
type sshKeyGenType interface {
	GetPrivateKey() string
}
