// Copyright 2016-2022 terraform-provider-sakuracloud authors
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/enhanceddb/builder"
)

func expandEnhancedDBBuilder(d *schema.ResourceData, client *APIClient, settingsHash string) *builder.Builder {
	return &builder.Builder{
		ID:           types.StringID(d.Id()),
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		Tags:         expandTags(d),
		IconID:       expandSakuraCloudID(d, "icon_id"),
		DatabaseName: d.Get("database_name").(string),
		Password:     d.Get("password").(string),
		SettingsHash: settingsHash,
		Client:       iaas.NewEnhancedDBOp(client),
	}
}
