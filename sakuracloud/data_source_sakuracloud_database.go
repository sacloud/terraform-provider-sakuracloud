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
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func dataSourceSakuraCloudDatabase() *schema.Resource {
	resourceName := "Database"

	return &schema.Resource{
		Read: dataSourceSakuraCloudDatabaseRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"database_type": {
				Type:     schema.TypeString,
				Computed: true,
				Description: descf(
					"The type of the database. This will be one of [%s]",
					types.RDBMSTypeStrings,
				),
			},
			"plan": schemaDataSourcePlan(resourceName, types.DatabasePlanStrings),
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of default user on the database",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The password of default user on the database",
			},
			"replica_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of user that processing a replication",
			},
			"replica_password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The password of user that processing a replication",
			},
			"network_interface": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch_id":     schemaDataSourceSwitchID(resourceName),
						"ip_address":    schemaDataSourceIPAddress(resourceName),
						"netmask":       schemaDataSourceNetMask(resourceName),
						"gateway":       schemaDataSourceGateway(resourceName),
						"port":          schemaDataSourcePort(),
						"source_ranges": schemaDataSourceSourceRanges(resourceName),
					},
				},
			},
			"backup": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"weekdays": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
							Description: descf(
								"The list of name of weekday that doing backup. This will be in [%s]",
								types.BackupWeekdayStrings,
							),
						},
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time to take backup. This will be formatted with `HH:mm`",
						},
					},
				},
			},
			"parameters": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "The map for setting RDMBS-specific parameters. Valid keys can be found with the `usacloud database list-parameters` command",
			},
			"icon_id":     schemaDataSourceIconID(resourceName),
			"description": schemaDataSourceDescription(resourceName),
			"tags":        schemaDataSourceTags(resourceName),
			"zone":        schemaDataSourceZone(resourceName),
		},
	}
}

func dataSourceSakuraCloudDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	searcher := sacloud.NewDatabaseOp(client)

	findCondition := &sacloud.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, zone, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud Database resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.Databases) == 0 {
		return filterNoResultErr()
	}

	targets := res.Databases
	d.SetId(targets[0].ID.String())
	return setDatabaseResourceData(ctx, d, client, targets[0], zone)
}
