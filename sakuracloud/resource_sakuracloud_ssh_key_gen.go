// Copyright 2016-2020 terraform-provider-sakuracloud authors
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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/api"
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
	client := meta.(*APIClient)

	name := d.Get("name").(string)
	passPhrase := ""
	if p, ok := d.GetOk("pass_phrase"); ok {
		passPhrase = p.(string)
	}

	description := ""
	if d, ok := d.GetOk("description"); ok {
		description = d.(string)
	}

	key, err := client.SSHKey.Generate(name, passPhrase, description)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud SSHKey resource: %s", err)
	}

	d.SetId(key.GetStrID())
	return setSSHKeyGenResourceData(d, client, key)
}

func resourceSakuraCloudSSHKeyGenRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	key, err := client.SSHKey.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud SSHKey resource: %s", err)
	}

	return setSSHKeyGenResourceData(d, client, key)
}

func resourceSakuraCloudSSHKeyGenDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	_, err := client.SSHKey.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud SSHKey resource: %s", err)
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
	GetStrID() string
	GetName() string
	GetPublicKey() string
	GetFingerprint() string
	GetDescription() string
}
type sshKeyGenType interface {
	GetPrivateKey() string
}
