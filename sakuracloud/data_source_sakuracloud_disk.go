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

func dataSourceSakuraCloudDisk() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudDiskRead,

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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connector": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Computed: true, //ReadOnly
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

func dataSourceSakuraCloudDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	//filters
	if rawFilter, filterOk := d.GetOk("filter"); filterOk {
		filters := expandFilters(rawFilter)
		for key, f := range filters {
			client.Disk.FilterBy(key, f)
		}
	}

	res, err := client.Disk.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Disk resource: %s", err)
	}
	if res == nil || res.Count == 0 {
		return filterNoResultErr()
	}
	var data *sacloud.Disk
	targets := res.Disks

	if rawNameSelector, ok := d.GetOk("name_selectors"); ok {
		selectors := expandStringList(rawNameSelector.([]interface{}))
		var filtered []sacloud.Disk
		for _, a := range targets {
			if hasNames(&a, selectors) {
				filtered = append(filtered, a)
			}
		}
		targets = filtered
	}
	if rawTagSelector, ok := d.GetOk("tag_selectors"); ok {
		selectors := expandStringList(rawTagSelector.([]interface{}))
		var filtered []sacloud.Disk
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
	return setDiskResourceData(d, client, data)
}
