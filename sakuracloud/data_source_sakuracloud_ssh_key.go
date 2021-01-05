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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func dataSourceSakuraCloudSSHKey() *schema.Resource {
	resourceName := "SSHKey"
	return &schema.Resource{
		Read: dataSourceSakuraCloudSSHKeyRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{excludeTags: true}),
			"name":         schemaDataSourceName(resourceName),
			"description":  schemaDataSourceDescription(resourceName),
			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value of public key",
			},
			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fingerprint of public key",
			},
		},
	}
}

func dataSourceSakuraCloudSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	searcher := sacloud.NewSSHKeyOp(client)

	findCondition := &sacloud.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud SSHKey resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.SSHKeys) == 0 {
		return filterNoResultErr()
	}

	targets := res.SSHKeys
	d.SetId(targets[0].ID.String())
	return setSSHKeyResourceData(ctx, d, client, targets[0])
}
