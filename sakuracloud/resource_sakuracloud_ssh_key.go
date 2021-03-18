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
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudSSHKey() *schema.Resource {
	resourceName := "SSHKey"
	return &schema.Resource{
		Create: resourceSakuraCloudSSHKeyCreate,
		Read:   resourceSakuraCloudSSHKeyRead,
		Update: resourceSakuraCloudSSHKeyUpdate,
		Delete: resourceSakuraCloudSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name":        schemaResourceName(resourceName),
			"description": schemaResourceDescription(resourceName),
			"public_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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

func resourceSakuraCloudSSHKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	sshKeyOp := sacloud.NewSSHKeyOp(client)

	key, err := sshKeyOp.Create(ctx, expandSSHKeyCreateRequest(d))
	if err != nil {
		return fmt.Errorf("creating SSHKey is failed: %s", err)
	}

	d.SetId(key.ID.String())
	return resourceSakuraCloudSSHKeyRead(d, meta)
}

func resourceSakuraCloudSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

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
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	sshKeyOp := sacloud.NewSSHKeyOp(client)

	key, err := sshKeyOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SSHKey[%s]: %s", d.Id(), err)
	}

	_, err = sshKeyOp.Update(ctx, key.ID, expandSSHKeyUpdateRequest(d))
	if err != nil {
		return fmt.Errorf("updating SSHKey[%s] is failed: %s", d.Id(), err)
	}
	return resourceSakuraCloudSSHKeyRead(d, meta)
}

func resourceSakuraCloudSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

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
		return fmt.Errorf("deleting SSHKey[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setSSHKeyResourceData(_ context.Context, d *schema.ResourceData, _ *APIClient, data *sacloud.SSHKey) error {
	d.Set("name", data.Name)               // nolint
	d.Set("public_key", data.PublicKey)    // nolint
	d.Set("fingerprint", data.Fingerprint) // nolint
	d.Set("description", data.Description) // nolint
	return nil
}
