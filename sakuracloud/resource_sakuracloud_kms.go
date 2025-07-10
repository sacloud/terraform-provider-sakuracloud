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
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "github.com/sacloud/api-client-go"
	"github.com/sacloud/kms-api-go"
	v1 "github.com/sacloud/kms-api-go/apis/v1"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudKMS() *schema.Resource {
	resourceName := "KMS"

	return &schema.Resource{
		CreateContext: resourceSakuraCloudKMSCreate,
		ReadContext:   resourceSakuraCloudKMSRead,
		UpdateContext: resourceSakuraCloudKMSUpdate,
		DeleteContext: resourceSakuraCloudKMSDelete,
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
			"key_origin": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "generated",
				ValidateDiagFunc: validateWithCustomFunc(func(v string) error {
					if v == "generated" || v == "imported" {
						return nil
					} else {
						return errors.New("key_origin must be 'generated' or 'imported'")
					}
				}),
				Description: "Key origin of the KMS key. 'generated' or 'imported'. Default is 'generated'.",
			},
			"plain_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Plain key for imported KMS key. Required when `key_origin` is 'imported'.",
			},
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
		},
	}
}

func resourceSakuraCloudKMSCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	keyReq, err := expandKMSCreateKey(d)
	if err != nil {
		return diag.Errorf("could not expand KMS create key: %s", err)
	}

	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	keyOp := kms.NewKeyOp(client.kmsClient)

	createdKey, err := keyOp.Create(ctx, keyReq)
	if err != nil {
		return diag.Errorf("create KMS queue failed: %s", err)
	}

	key := v1.Key{
		ID:          createdKey.ID,
		Name:        createdKey.Name,
		Description: createdKey.Description,
		KeyOrigin:   createdKey.KeyOrigin,
		Tags:        createdKey.Tags,
	}
	return setKMSResourceData(d, &key)
}

func resourceSakuraCloudKMSRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	keyOp := kms.NewKeyOp(client.kmsClient)

	key, err := keyOp.Read(ctx, d.Id())
	if err != nil {
		if api.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read KMS[%s] key: %s", d.Id(), err)
	}

	return setKMSResourceData(d, key)
}

func resourceSakuraCloudKMSUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	keyOp := kms.NewKeyOp(client.kmsClient)

	key, err := keyOp.Read(ctx, d.Id())
	if err != nil {
		return diag.Errorf("could not read KMS[%s] key: %s", d.Id(), err)
	}

	if _, err = keyOp.Update(ctx, key.ID, expandKMSUpdateKey(d, key)); err != nil {
		return diag.Errorf("could not update KMS[%s] key: %s", d.Id(), err)
	}

	return resourceSakuraCloudKMSRead(ctx, d, meta)
}

func resourceSakuraCloudKMSDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	keyOp := kms.NewKeyOp(client.kmsClient)

	key, err := keyOp.Read(ctx, d.Id())
	if err != nil {
		if api.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read KMS[%s] key: %s", d.Id(), err)
	}

	if err := keyOp.Delete(ctx, key.ID); err != nil {
		return diag.Errorf("could not delete KMS[%s] key: %s", d.Id(), err)
	}
	return nil
}

func setKMSResourceData(d *schema.ResourceData, data *v1.Key) diag.Diagnostics {
	d.SetId(data.ID)
	d.Set("name", data.Name)                    //nolint:errcheck,gosec
	d.Set("key_origin", string(data.KeyOrigin)) //nolint:errcheck,gosec
	if data.Description.IsSet() {
		d.Set("description", data.Description.Value) //nolint:errcheck,gosec
	}
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
