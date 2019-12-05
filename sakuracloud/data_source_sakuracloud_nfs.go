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
)

func dataSourceSakuraCloudNFS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudNFSRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_mask_len": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_route": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"icon_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		},
	}
}

func dataSourceSakuraCloudNFSRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	searcher := sacloud.NewNFSOp(client)

	findCondition := &sacloud.FindCondition{
		Count: defaultSearchLimit,
	}
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
