// Copyright 2016-2023 terraform-provider-sakuracloud authors
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

func dataSourceSakuraCloudESME() *schema.Resource {
	resourceName := "ESME"
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudESMERead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"icon_id":      schemaDataSourceIconID(resourceName),
			"description":  schemaDataSourceDescription(resourceName),
			"tags":         schemaDataSourceTags(resourceName),
			"send_message_with_generated_otp_api_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The API URL for send SMS with generated OTP",
			},
			"send_message_with_inputted_otp_api_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The API URL for send SMS with inputted OTP",
			},
		},
	}
}

func dataSourceSakuraCloudESMERead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	searcher := iaas.NewESMEOp(client)

	findCondition := &iaas.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud ESME: %s", err)
	}
	if res == nil || res.Count == 0 {
		return filterNoResultErr()
	}

	targets := res.ESME
	if len(targets) == 0 {
		return filterNoResultErr()
	}

	d.SetId(targets[0].ID.String())
	return setESMEResourceData(d, client, targets[0])
}
