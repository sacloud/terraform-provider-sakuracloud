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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func dataSourceSakuraCloudPrivateHost() *schema.Resource {
	resourceName := "PrivateHost"
	return &schema.Resource{
		Read: dataSourceSakuraCloudPrivateHostRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"class":        schemaDataSourceClass(resourceName, []string{types.PrivateHostClassDynamic, types.PrivateHostClassWindows}),
			"icon_id":      schemaDataSourceIconID(resourceName),
			"description":  schemaDataSourceDescription(resourceName),
			"tags":         schemaDataSourceTags(resourceName),
			"hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The hostname of the private host",
			},
			"assigned_core": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of CPUs assigned to servers on the private host",
			},
			"assigned_memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total size of memory assigned to servers on the private host",
			},
			"zone": schemaDataSourceZone(resourceName),
		},
	}
}

func dataSourceSakuraCloudPrivateHostRead(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	searcher := sacloud.NewPrivateHostOp(client)

	findCondition := &sacloud.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, zone, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud PrivateHost resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.PrivateHosts) == 0 {
		return filterNoResultErr()
	}

	targets := res.PrivateHosts
	d.SetId(targets[0].ID.String())
	return setPrivateHostResourceData(ctx, d, client, targets[0])
}
