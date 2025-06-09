// Copyright 2016-2025 terraform-provider-sakuracloud authors
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
)

func dataSourceSakuraCloudZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudZoneRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the zone (e.g. `is1a`,`tk1a`)",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the zone",
			},
			"description": schemaDataSourceDescription("zone"),
			"region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the region that the zone belongs",
			},
			"region_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the region that the zone belongs",
			},
			"dns_servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of IP address of DNS server in the zone",
			},
		},
	}
}

func dataSourceSakuraCloudZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zoneSlug, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	zoneOp := iaas.NewZoneOp(client)

	if v, ok := d.GetOk("name"); ok {
		zoneSlug = v.(string)
	}

	res, err := zoneOp.Find(ctx, &iaas.FindCondition{})
	if err != nil {
		return diag.Errorf("could not find SakuraCloud Zone resource: %s", err)
	}
	if res == nil || len(res.Zones) == 0 {
		return filterNoResultErr()
	}
	var data *iaas.Zone

	for _, z := range res.Zones {
		if z.Name == zoneSlug {
			data = z
			break
		}
	}
	if data == nil {
		return filterNoResultErr()
	}

	d.SetId(data.ID.String())
	d.Set("name", data.Name)                    //nolint
	d.Set("zone_id", data.ID.String())          //nolint
	d.Set("description", data.Description)      //nolint
	d.Set("region_id", data.Region.ID.String()) //nolint
	d.Set("region_name", data.Region.Name)      //nolint
	return diag.FromErr(d.Set("dns_servers", data.Region.NameServers))
}
