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
	"slices"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sm "github.com/sacloud/secretmanager-api-go"
	v1 "github.com/sacloud/secretmanager-api-go/apis/v1"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudSecretManagerSecret() *schema.Resource {
	resourceName := "SecretManagerSecret"

	return &schema.Resource{
		CreateContext: resourceSakuraCloudSecretManagerSecretCreate,
		ReadContext:   resourceSakuraCloudSecretManagerSecretRead,
		//UpdateContext: resourceSakuraCloudSecretManagerSecretUpdate,
		UpdateContext: resourceSakuraCloudSecretManagerSecretCreate,
		DeleteContext: resourceSakuraCloudSecretManagerSecretDelete,
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
			"vault_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: desc.Sprintf("The secret manager's vault id of the %s", resourceName),
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Version of secret value. This value is incremented by create/update",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Secret value",
				Sensitive:   true,
			},
		},
	}
}

func resourceSakuraCloudSecretManagerSecretCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("vault_id").(string)
	secretOp := sm.NewSecretOp(client.secretmanagerClient, id)

	createdSec, err := secretOp.Create(ctx, expandSecretManagerCreateSecret(d))
	if err != nil {
		return diag.Errorf("could not create SecretManagerSecret secret: %s", err)
	}

	secret := v1.Secret{
		Name:          createdSec.Name,
		LatestVersion: createdSec.LatestVersion,
	}
	return setSecretManagerSecretResourceData(d, &secret)
}

func resourceSakuraCloudSecretManagerSecretRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	id := d.Get("vault_id").(string)
	secretOp := sm.NewSecretOp(client.secretmanagerClient, id)

	secret, err := filterSecretManagerSecretByName(d, ctx, secretOp, name)
	if err != nil {
		return diag.Errorf("could not read SecretManagerSecret[%s] secret: %s", name, err)
	}

	return setSecretManagerSecretResourceData(d, secret)
}

func resourceSakuraCloudSecretManagerSecretDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	id := d.Get("vault_id").(string)
	secretOp := sm.NewSecretOp(client.secretmanagerClient, id)

	_, err = filterSecretManagerSecretByName(d, ctx, secretOp, name)
	if err != nil {
		return diag.Errorf("could not read SecretManagerSecret[%s] secret: %s", name, err)
	}

	if err := secretOp.Delete(ctx, v1.DeleteSecret{Name: name}); err != nil {
		return diag.Errorf("could not delete SecretManagerSecret[%s] secret: %s", name, err)
	}

	return nil
}

func filterSecretManagerSecretByName(d *schema.ResourceData, ctx context.Context, secretOp sm.SecretAPI, name string) (*v1.Secret, error) {
	secrets, err := secretOp.List(ctx)
	if err != nil {
		return nil, err
	}

	match := slices.Collect(func(yield func(v1.Secret) bool) {
		for _, v := range secrets {
			if name != v.Name {
				continue
			}
			if !yield(v) {
				return
			}
		}
	})

	if len(match) == 0 {
		d.SetId("")
		return nil, errFilterNoResult
	}

	return &match[0], nil
}

func setSecretManagerSecretResourceData(d *schema.ResourceData, data *v1.Secret) diag.Diagnostics {
	d.SetId(data.Name)
	d.Set("name", data.Name)             //nolint:errcheck,gosec
	d.Set("version", data.LatestVersion) //nolint:errcheck,gosec

	return nil
}
