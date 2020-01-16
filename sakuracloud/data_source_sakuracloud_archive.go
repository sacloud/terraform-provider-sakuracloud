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
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/sacloud/ostype"
)

func dataSourceSakuraCloudArchive() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudArchiveRead,

		Schema: map[string]*schema.Schema{
			"os_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice(ostype.OSTypeShortNames, false),
				ConflictsWith: []string{"filter", "name_selectors", "tag_selectors"},
			},
			"name_selectors": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"os_type"},
			},
			"tag_selectors": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"os_type"},
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
				ConflictsWith: []string{"os_type"},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
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

func dataSourceSakuraCloudArchiveRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	var data *sacloud.Archive
	if osType, ok := d.GetOk("os_type"); ok {
		strOSType := osType.(string)
		if strOSType != "" {

			res, err := client.Archive.FindByOSType(strToOSType(strOSType))
			if err != nil {
				return filterNoResultErr()
			}
			data = res
		}
	} else {

		//filters
		if rawFilter, filterOk := d.GetOk("filter"); filterOk {
			filters := expandFilters(rawFilter)
			for key, f := range filters {
				client.Archive.FilterBy(key, f)
			}
		}

		res, err := client.Archive.Find()
		if err != nil {
			return fmt.Errorf("Couldn't find SakuraCloud Archive resource: %s", err)
		}
		if res == nil || res.Count == 0 {
			return filterNoResultErr()
		}

		targets := res.Archives

		if rawNameSelector, ok := d.GetOk("name_selectors"); ok {
			selectors := expandStringList(rawNameSelector.([]interface{}))
			var filtered []sacloud.Archive
			for _, a := range targets {
				if hasNames(&a, selectors) {
					filtered = append(filtered, a)
				}
			}
			targets = filtered
		}
		if rawTagSelector, ok := d.GetOk("tag_selectors"); ok {
			selectors := expandStringList(rawTagSelector.([]interface{}))
			var filtered []sacloud.Archive
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
	}

	if data != nil {

		d.SetId(data.GetStrID())
		d.Set("name", data.Name)
		d.Set("size", toSizeGB(data.SizeMB))
		d.Set("icon_id", data.GetIconStrID())
		d.Set("description", data.Description)
		d.Set("tags", data.Tags)

		d.Set("zone", client.Zone)
	}

	return nil
}

func strToOSType(strType string) ostype.ArchiveOSTypes {
	return ostype.StrToOSType(strType)
}
