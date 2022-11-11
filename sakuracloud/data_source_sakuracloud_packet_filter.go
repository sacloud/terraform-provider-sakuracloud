// Copyright 2016-2022 terraform-provider-sakuracloud authors
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
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func dataSourceSakuraCloudPacketFilter() *schema.Resource {
	resourceName := "PacketFilter"
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudPacketFilterRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{excludeTags: true}),
			"name":         schemaDataSourceName(resourceName),
			"description":  schemaDataSourceDescription(resourceName),
			"expression": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
							Description: desc.Sprintf(
								"The protocol used for filtering. This will be one of [%s]",
								types.PacketFilterProtocolStrings,
							),
						},
						"source_network": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A source IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`)",
						},
						"source_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A source port number or port range used for filtering (e.g. `1024`, `1024-2048`)",
						},
						"destination_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A destination port number or port range used for filtering (e.g. `1024`, `1024-2048`)",
						},
						"allow": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The flag to allow the packet through the filter",
						},
						"description": schemaDataSourceDescription("expression"),
					},
				},
			},
			"zone": schemaDataSourceZone(resourceName),
		},
	}
}

func dataSourceSakuraCloudPacketFilterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	searcher := iaas.NewPacketFilterOp(client)

	findCondition := &iaas.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, zone, findCondition)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud PacketFilter resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.PacketFilters) == 0 {
		return filterNoResultErr()
	}

	targets := res.PacketFilters
	d.SetId(targets[0].ID.String())
	return setPacketFilterResourceData(ctx, d, client, targets[0])
}
