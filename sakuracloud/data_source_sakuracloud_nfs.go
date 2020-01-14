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

func dataSourceSakuraCloudNFS() *schema.Resource {
	resourceName := "NFS"
	return &schema.Resource{
		Read: dataSourceSakuraCloudNFSRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"plan":         schemaDataSourcePlan(resourceName, types.NFSPlanStrings),
			"size":         schemaDataSourceSize(resourceName),
			"switch_id":    schemaDataSourceSwitchID(resourceName),
			"ip_address":   schemaDataSourceIPAddress(resourceName),
			"netmask":      schemaDataSourceNetMask(resourceName),
			"gateway":      schemaDataSourceGateway(resourceName),
			"icon_id":      schemaDataSourceIconID(resourceName),
			"description":  schemaDataSourceDescription(resourceName),
			"tags":         schemaDataSourceTags(resourceName),
			"zone":         schemaDataSourceZone(resourceName),
		},
	}
}

func dataSourceSakuraCloudNFSRead(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	searcher := sacloud.NewNFSOp(client)

	findCondition := &sacloud.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, zone, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud NFS resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.NFS) == 0 {
		return filterNoResultErr()
	}

	targets := res.NFS
	d.SetId(targets[0].ID.String())
	return setNFSResourceData(ctx, d, client, targets[0])
}
