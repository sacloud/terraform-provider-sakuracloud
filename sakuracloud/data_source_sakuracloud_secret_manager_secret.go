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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sm "github.com/sacloud/secretmanager-api-go"
	v1 "github.com/sacloud/secretmanager-api-go/apis/v1"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func dataSourceSakuraCloudSecretManagerSecret() *schema.Resource {
	const resourceName = "SecretManagerSecret"
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudSecretManagerSecretRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: desc.Sprintf("The name of the %s", resourceName),
			},
			"vault_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: desc.Sprintf("The secret manager's vault id of the %s", resourceName),
			},
			"version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Target version to unveil stored secret. Without this parameter, latest version is used",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unveiled result of stored secret",
				Sensitive:   true,
			},
		},
	}
}

func dataSourceSakuraCloudSecretManagerSecretRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	id := d.Get("vault_id").(string)
	secretOp := sm.NewSecretOp(client.secretmanagerClient, id)

	unveilReq := v1.Unveil{Name: name}
	if ver, ok := d.GetOk("version"); ok {
		unveilReq.Version = v1.NewOptNilInt(ver.(int))
	}
	unveil, err := secretOp.Unveil(ctx, unveilReq)
	if err != nil {
		return diag.FromErr(err)
	}

	return setSecretManagerUnveilResourceData(d, unveil)
}

func setSecretManagerUnveilResourceData(d *schema.ResourceData, data *v1.Unveil) diag.Diagnostics {
	d.SetId(data.Name)
	d.Set("name", data.Name)   //nolint:errcheck,gosec
	d.Set("value", data.Value) //nolint:errcheck,gosec
	if data.Version.IsSet() {
		d.Set("version", data.Version.Value) //nolint:errcheck,gosec
	}
	return nil
}
