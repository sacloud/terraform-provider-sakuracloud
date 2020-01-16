// Copyright 2016-2020 terraform-provider-sakuracloud authors
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
	"github.com/sacloud/libsacloud/sacloud"
)

func dataSourceSakuraCloudDNS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudDNSRead,

		Schema: map[string]*schema.Schema{
			"name_selectors": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tag_selectors": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dns_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSakuraCloudDNSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	//filters
	if rawFilter, filterOk := d.GetOk("filter"); filterOk {
		filters := expandFilters(rawFilter)
		for key, f := range filters {
			client.DNS.FilterBy(key, f)
		}
	}

	res, err := client.DNS.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud DNS resource: %s", err)
	}
	if res == nil || res.Count == 0 {
		return filterNoResultErr()
	}
	var data *sacloud.DNS
	targets := res.CommonServiceDNSItems

	if rawNameSelector, ok := d.GetOk("name_selectors"); ok {
		selectors := expandStringList(rawNameSelector.([]interface{}))
		var filtered []sacloud.DNS
		for _, a := range targets {
			if hasNames(&a, selectors) {
				filtered = append(filtered, a)
			}
		}
		targets = filtered
	}
	if rawTagSelector, ok := d.GetOk("tag_selectors"); ok {
		selectors := expandStringList(rawTagSelector.([]interface{}))
		var filtered []sacloud.DNS
		for _, a := range targets {
			if hasTags(&a, selectors) {
				filtered = append(filtered, a)
			}
		}
		targets = filtered
	}

	if len(targets) == 0 {
		return filterNoResultErr()
	}
	data = &targets[0]

	d.SetId(data.GetStrID())
	return setDNSResourceData(d, client, data)
}
