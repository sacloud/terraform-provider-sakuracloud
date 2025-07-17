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
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/kms-api-go"
	v1 "github.com/sacloud/kms-api-go/apis/v1"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func dataSourceSakuraCloudKMS() *schema.Resource {
	const resourceName = "KMS"
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudKMSRead,

		Schema: map[string]*schema.Schema{
			"name": func() *schema.Schema {
				s := schemaDataSourceName(resourceName)
				s.Optional = true
				return s
			}(),
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: desc.Sprintf("The resource id of the %s", resourceName),
			},
			"description": schemaDataSourceDescription(resourceName),
			"tags":        schemaDataSourceTags(resourceName),
			"key_origin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Key origin of the KMS key. 'generated' or 'imported'",
			},
		},
	}
}

func dataSourceSakuraCloudKMSRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	name := d.Get("name").(string)
	id := d.Get("resource_id").(string)
	if name == "" && id == "" {
		return diag.Errorf("name or resource_id required")
	}

	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	keyOp := kms.NewKeyOp(client.kmsClient)

	var key *v1.Key
	if name != "" {
		keys, err := keyOp.List(ctx)
		if err != nil {
			return diag.Errorf("could not find KMS resource: %s", err)
		}

		key, err = filterKMSByName(keys, name)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		key, err = keyOp.Read(ctx, id)
		if err != nil {
			return filterNoResultErr()
		}
	}

	return setKMSResourceData(d, key)
}

func filterKMSByName(keys v1.Keys, name string) (*v1.Key, error) {
	match := slices.Collect(func(yield func(v1.Key) bool) {
		for _, v := range keys {
			if name != v.Name {
				continue
			}

			if !yield(v) {
				return
			}
		}
	})

	if len(match) == 0 {
		return nil, errFilterNoResult
	}
	if len(match) > 1 {
		return nil, fmt.Errorf("multiple KMS resources found with the same condition. name=%q", name)
	}
	return &match[0], nil
}
