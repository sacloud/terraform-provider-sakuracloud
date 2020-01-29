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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	simBuilder "github.com/sacloud/libsacloud/v2/utils/builder/sim"
)

func expandSIMCarrier(d resourceValueGettable) []*sacloud.SIMNetworkOperatorConfig {
	// carriers
	var carriers []*sacloud.SIMNetworkOperatorConfig
	rawCarriers := d.Get("carrier").(*schema.Set).List()
	for _, carrier := range rawCarriers {
		carriers = append(carriers, &sacloud.SIMNetworkOperatorConfig{
			Allow: true,
			Name:  types.SIMOperatorShortNameMap[carrier.(string)].String(),
		})
	}
	return carriers
}

func flattenSIMCarrier(carrierInfo []*sacloud.SIMNetworkOperatorConfig) *schema.Set {
	set := &schema.Set{F: schema.HashString}
	for _, c := range carrierInfo {
		if !c.Allow {
			continue
		}
		for k := range types.SIMOperatorShortNameMap {
			if k == c.Name {
				set.Add(c.Name)
			}
		}
	}
	return set
}

func expandSIMBuilder(d resourceValueGettable, client *APIClient) *simBuilder.Builder {
	return &simBuilder.Builder{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),
		ICCID:       d.Get("iccid").(string),
		PassCode:    d.Get("passcode").(string),
		Activate:    d.Get("enabled").(bool),
		IMEI:        d.Get("imei").(string),
		Carrier:     expandSIMCarrier(d),
		Client:      simBuilder.NewAPIClient(client),
	}
}
