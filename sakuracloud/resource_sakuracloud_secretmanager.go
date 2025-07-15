// Copyright 2016-2025 terraform-provider-sakuracloud authors
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
	api "github.com/sacloud/api-client-go"
	sm "github.com/sacloud/secretmanager-api-go"
	v1 "github.com/sacloud/secretmanager-api-go/apis/v1"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudSecretManager() *schema.Resource {
	resourceName := "SecretManager"

	return &schema.Resource{
		CreateContext: resourceSakuraCloudSecretManagerCreate,
		ReadContext:   resourceSakuraCloudSecretManagerRead,
		UpdateContext: resourceSakuraCloudSecretManagerUpdate,
		DeleteContext: resourceSakuraCloudSecretManagerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: desc.Sprintf("The name of the %s.", resourceName),
			},
			"kms_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "KMS key ID for the SecretManager vault",
			},
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
		},
	}
}

func resourceSakuraCloudSecretManagerCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	vaultReq, err := expandSecretManagerCreateVault(d)
	if err != nil {
		return diag.Errorf("could not expand SecretManager create vault: %s", err)
	}

	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	vaultOp := sm.NewVaultOp(client.secretmanagerClient)

	createdKey, err := vaultOp.Create(ctx, vaultReq)
	if err != nil {
		return diag.Errorf("create SecretManager vault failed: %s", err)
	}

	vault := v1.Vault{
		ID:          createdKey.ID,
		Name:        createdKey.Name,
		Description: createdKey.Description,
		KmsKeyID:    createdKey.KmsKeyID,
		Tags:        createdKey.Tags,
	}
	return setSecretManagerVaultResourceData(d, &vault)
}

func resourceSakuraCloudSecretManagerRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	vaultOp := sm.NewVaultOp(client.secretmanagerClient)

	vault, err := vaultOp.Read(ctx, d.Id())
	if err != nil {
		if api.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SecretManager[%s] vault: %s", d.Id(), err)
	}

	return setSecretManagerVaultResourceData(d, vault)
}

func resourceSakuraCloudSecretManagerUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	vaultOp := sm.NewVaultOp(client.secretmanagerClient)

	vault, err := vaultOp.Read(ctx, d.Id())
	if err != nil {
		return diag.Errorf("could not read SecretManager[%s] vault: %s", d.Id(), err)
	}

	if _, err = vaultOp.Update(ctx, vault.ID, expandSecretManagerUpdateVault(d, vault)); err != nil {
		return diag.Errorf("could not update SecretManager[%s] vault: %s", d.Id(), err)
	}

	return resourceSakuraCloudSecretManagerRead(ctx, d, meta)
}

func resourceSakuraCloudSecretManagerDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	vaultOp := sm.NewVaultOp(client.secretmanagerClient)

	vault, err := vaultOp.Read(ctx, d.Id())
	if err != nil {
		if api.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SecretManager[%s] vault: %s", d.Id(), err)
	}

	if err := vaultOp.Delete(ctx, vault.ID); err != nil {
		return diag.Errorf("could not delete SecretManager[%s] vault: %s", d.Id(), err)
	}
	return nil
}

func setSecretManagerVaultResourceData(d *schema.ResourceData, data *v1.Vault) diag.Diagnostics {
	d.SetId(data.ID)
	d.Set("name", data.Name)           //nolint:errcheck,gosec
	d.Set("kms_key_id", data.KmsKeyID) //nolint:errcheck,gosec
	if data.Description.IsSet() {
		d.Set("description", data.Description.Value) //nolint:errcheck,gosec
	}
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
