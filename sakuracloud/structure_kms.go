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
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/sacloud/kms-api-go/apis/v1"
)

func expandKMSCreateKey(d *schema.ResourceData) (kms.CreateKey, error) {
	keyOrig := d.Get("key_origin").(string)
	var req kms.CreateKey
	if keyOrig == "generated" {
		req = kms.CreateKey{
			Name:      d.Get("name").(string),
			Tags:      expandTags(d),
			KeyOrigin: kms.KeyOriginEnumGenerated,
		}
	} else {
		plainKey := d.Get("plain_key").(string)
		if plainKey == "" {
			return kms.CreateKey{}, errors.New("plain_key is required when key_origin is 'imported'")
		}
		req = kms.CreateKey{
			Name:      d.Get("name").(string),
			Tags:      expandTags(d),
			KeyOrigin: kms.KeyOriginEnumImported,
			PlainKey:  kms.NewOptString(plainKey),
		}
	}

	if desc, ok := d.GetOk("description"); ok {
		req.Description = kms.NewOptString(desc.(string))
	}

	return req, nil
}

func expandKMSUpdateKey(d *schema.ResourceData, before *kms.Key) kms.Key {
	req := kms.Key{
		Name: d.Get("name").(string),
	}

	if _, ok := d.GetOk("tags"); ok {
		req.Tags = expandTags(d)
	}
	if desc, ok := d.GetOk("description"); ok {
		req.Description = kms.NewOptString(desc.(string))
	}

	return req
}
