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
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/ostype"
	"github.com/sacloud/libsacloud/v2/utils/archive"
)

func dataSourceSakuraCloudArchive() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudArchiveRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"os_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringInSlice(ostype.OSTypeShortNames, false),
				ConflictsWith: []string{"filters"},
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
	client, ctx, zone := getSacloudV2Client(d, meta)
	searcher := sacloud.NewArchiveOp(client)

	var data *sacloud.Archive
	if osType, ok := d.GetOk("os_type"); ok {
		strOSType := osType.(string)
		if strOSType != "" {
			res, err := archive.FindByOSType(ctx, searcher, zone, ostype.StrToOSType(strOSType))
			if err != nil {
				return err
			}
			data = res
		}
	} else {
		findCondition := &sacloud.FindCondition{
			Count: defaultSearchLimit,
		}
		if rawFilter, ok := d.GetOk(filterAttrName); ok {
			findCondition.Filter = expandSearchFilter(rawFilter)
		}

		res, err := searcher.Find(ctx, zone, findCondition)
		if err != nil {
			return fmt.Errorf("could not find SakuraCloud Archive resource: %s", err)
		}
		if res == nil || res.Count == 0 {
			return filterNoResultErr()
		}

		targets := res.Archives
		if len(targets) == 0 {
			return filterNoResultErr()
		}
		data = targets[0]
	}

	if data != nil {
		d.SetId(data.ID.String())
		return setResourceData(d, map[string]interface{}{
			"name":        data.Name,
			"size":        data.GetSizeGB(),
			"icon_id":     data.IconID.String(),
			"description": data.Description,
			"tags":        data.Tags,
			"zone":        getV2Zone(d, client),
		})
	}
	return nil
}
