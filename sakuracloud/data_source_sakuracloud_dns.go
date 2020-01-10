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

func dataSourceSakuraCloudDNS() *schema.Resource {
	resourceName := "DNS"
	return &schema.Resource{
		Read: dataSourceSakuraCloudDNSRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of managed domain",
			},
			"dns_servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of IP address of DNS server that manage this zone",
			},
			"icon_id":     schemaDataSourceIconID(resourceName),
			"description": schemaDataSourceDescription(resourceName),
			"tags":        schemaDataSourceTags(resourceName),
			"record": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": schemaDataSourceName("DNS Record"),
						"type": {
							Type:     schema.TypeString,
							Computed: true,
							Description: descf(
								"The type of DNS Record. This will be one of [%s]",
								types.DNSRecordTypesStrings(),
							),
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value of the DNS Record",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of the TTL",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The priority of target DNS Record",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The weight of target DNS Record",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of port",
						},
					},
				},
			},
		},
	}
}

func dataSourceSakuraCloudDNSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	searcher := sacloud.NewDNSOp(client)

	findCondition := &sacloud.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud Disk resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.DNS) == 0 {
		return filterNoResultErr()
	}

	targets := res.DNS
	d.SetId(targets[0].ID.String())
	return setDNSResourceData(ctx, d, client, targets[0])
}
