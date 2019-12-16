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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	internetBuilder "github.com/sacloud/libsacloud/v2/utils/builder/internet"
)

func expandInternetBuilder(d *schema.ResourceData, client *APIClient) *internetBuilder.Builder {
	return &internetBuilder.Builder{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Tags:           expandTags(d),
		IconID:         expandSakuraCloudID(d, "icon_id"),
		NetworkMaskLen: d.Get("netmask").(int),
		BandWidthMbps:  d.Get("band_width").(int),
		EnableIPv6:     d.Get("enable_ipv6").(bool),
		Client:         internetBuilder.NewAPIClient(client),
	}
}
