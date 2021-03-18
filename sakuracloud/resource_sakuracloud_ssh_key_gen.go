// Copyright 2016-2021 terraform-provider-sakuracloud authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sakuracloud

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudSSHKeyGen() *schema.Resource {
	resourceName := "SSHKey"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudSSHKeyGenCreate,
		ReadContext:   resourceSakuraCloudSSHKeyGenRead,
		DeleteContext: resourceSakuraCloudSSHKeyGenDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  descf("The name of the %s. %s", resourceName, descLength(1, 64)),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
				Description:  descf("The description of the %s. %s", resourceName, descLength(1, 512)),
			},
			"pass_phrase": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(8, 64),
				Description: descf(
					"The pass phrase of the private key. %s",
					descLength(8, 64),
				),
			},
			"private_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The body of the private key",
			},
			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The body of the public key",
			},
			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fingerprint of the public key",
			},
		},
	}
}

func resourceSakuraCloudSSHKeyGenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	sshKeyOp := sacloud.NewSSHKeyOp(client)

	key, err := sshKeyOp.Generate(ctx, expandSSHKeyGenerateRequest(d))
	if err != nil {
		return diag.Errorf("generating SSHKey is failed: %s", err)
	}

	d.SetId(key.ID.String())

	// Note: CreateのレスポンスにのみPrivateKeyが含まれる
	return setSSHKeyGenResourceData(d, client, key)
}

func resourceSakuraCloudSSHKeyGenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	sshKeyOp := sacloud.NewSSHKeyOp(client)
	key, err := sshKeyOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SSHKey[%s]: %s", d.Id(), err)
	}

	return setSSHKeyGenResourceData(d, client, key)
}

func resourceSakuraCloudSSHKeyGenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	sshKeyOp := sacloud.NewSSHKeyOp(client)
	key, err := sshKeyOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SSHKey[%s]: %s", d.Id(), err)
	}

	if err := sshKeyOp.Delete(ctx, key.ID); err != nil {
		return diag.Errorf("deleting SSHKey[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setSSHKeyGenResourceData(d *schema.ResourceData, _ *APIClient, data interface{}) diag.Diagnostics {
	if key, ok := data.(sshKeyType); ok {
		d.Set("name", key.GetName())               // nolint
		d.Set("public_key", key.GetPublicKey())    // nolint
		d.Set("fingerprint", key.GetFingerprint()) // nolint
		d.Set("description", key.GetDescription()) // nolint

		if pKey, ok := data.(sshKeyGenType); ok {
			d.Set("private_key", pKey.GetPrivateKey()) // nolint
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
