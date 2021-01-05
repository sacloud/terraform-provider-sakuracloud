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

func dataSourceSakuraCloudLocalRouter() *schema.Resource {
	resourceName := "LocalRouter"
	return &schema.Resource{
		Read: dataSourceSakuraCloudLocalRouterRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"switch": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource ID of the Switch",
						},
						"category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The category name of connected services (e.g. `cloud`, `vps`)",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the Zone",
						},
					},
				},
			},
			"network_interface": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The virtual IP address",
						},
						"ip_addresses": schemaDataSourceIPAddresses(resourceName),
						"netmask":      schemaDataSourceNetMask(resourceName),
						"vrid": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Virtual Router Identifier",
						},
					},
				},
			},
			"peer": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"peer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the peer LocalRouter",
						},
						"secret_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The secret key of the peer LocalRouter",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The flag to enable the LocalRouter",
						},
						"description": schemaDataSourceDescription(resourceName),
					},
				},
			},
			"static_route": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CIDR block of destination",
						},
						"next_hop": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address of the next hop",
						},
					},
				},
			},
			"icon_id":     schemaDataSourceIconID(resourceName),
			"description": schemaDataSourceDescription(resourceName),
			"tags":        schemaDataSourceTags(resourceName),
			"secret_keys": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:      schema.TypeString,
					Sensitive: true,
				},
				Computed:    true,
				Description: "A list of secret key used for peering from other LocalRouters",
			},
		},
	}
}

func dataSourceSakuraCloudLocalRouterRead(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	searcher := sacloud.NewLocalRouterOp(client)

	findCondition := &sacloud.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud LocalRouter resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.LocalRouters) == 0 {
		return filterNoResultErr()
	}

	targets := res.LocalRouters
	d.SetId(targets[0].ID.String())
	return setLocalRouterResourceData(ctx, d, client, targets[0])
}
