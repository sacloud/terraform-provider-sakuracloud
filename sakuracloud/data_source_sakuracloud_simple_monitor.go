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

func dataSourceSakuraCloudSimpleMonitor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudSimpleMonitorRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"target": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delay_loop": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"health_check": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_header": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sni": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"password": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"qname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"excepcted_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"community": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snmp_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remaining_days": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
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
			"notify_email_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"notify_email_html": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"notify_slack_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"notify_slack_webhook": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSakuraCloudSimpleMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	searcher := sacloud.NewSimpleMonitorOp(client)

	findCondition := &sacloud.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud SimpleMonitor resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.SimpleMonitors) == 0 {
		return filterNoResultErr()
	}

	targets := res.SimpleMonitors
	d.SetId(targets[0].ID.String())
	return setSimpleMonitorResourceData(ctx, d, client, targets[0])
}
