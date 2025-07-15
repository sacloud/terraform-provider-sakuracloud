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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sm "github.com/sacloud/secretmanager-api-go/apis/v1"
)

func expandSecretManagerCreateVault(d *schema.ResourceData) (sm.CreateVault, error) {
	req := sm.CreateVault{
		Name:     d.Get("name").(string),
		Tags:     expandTags(d),
		KmsKeyID: d.Get("kms_key_id").(string),
	}

	if desc, ok := d.GetOk("description"); ok {
		req.Description = sm.NewOptString(desc.(string))
	}

	return req, nil
}

func expandSecretManagerUpdateVault(d *schema.ResourceData, before *sm.Vault) sm.Vault {
	req := sm.Vault{
		Name:     d.Get("name").(string),
		KmsKeyID: before.KmsKeyID,
	}

	if _, ok := d.GetOk("tags"); ok {
		req.Tags = expandTags(d)
	}
	if desc, ok := d.GetOk("description"); ok {
		req.Description = sm.NewOptString(desc.(string))
	}

	return req
}
