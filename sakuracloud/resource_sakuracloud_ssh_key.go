// Copyright 2016-2019 terraform-provider-sakuracloud authors
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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
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
	client, ctx, _ := getSacloudClient(d, meta)
	sshKeyOp := sacloud.NewSSHKeyOp(client)

	key, err := sshKeyOp.Create(ctx, expandSSHKeyCreateRequest(d))
	if err != nil {
		return fmt.Errorf("creating SSHKey is failed: %s", err)
	}

	d.SetId(key.ID.String())
	return resourceSakuraCloudSSHKeyRead(d, meta)
}

func resourceSakuraCloudSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	sshKeyOp := sacloud.NewSSHKeyOp(client)

	key, err := sshKeyOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SSHKey[%s]: %s", d.Id(), err)
	}
	return setSSHKeyResourceData(ctx, d, client, key)
}

func resourceSakuraCloudSSHKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	sshKeyOp := sacloud.NewSSHKeyOp(client)

	key, err := sshKeyOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SSHKey[%s]: %s", d.Id(), err)
	}

	_, err = sshKeyOp.Update(ctx, key.ID, expandSSHKeyUpdateRequest(d))
	if err != nil {
		return fmt.Errorf("updating SSHKey[%s] is failed: %s", key.ID, err)
	}
	return resourceSakuraCloudSSHKeyRead(d, meta)
}

func resourceSakuraCloudSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	sshKeyOp := sacloud.NewSSHKeyOp(client)

	key, err := sshKeyOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SSHKey[%s]: %s", d.Id(), err)
	}

	if err := sshKeyOp.Delete(ctx, key.ID); err != nil {
		return fmt.Errorf("deleting SSHKey[%s] is failed: %s", key.ID, err)
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
