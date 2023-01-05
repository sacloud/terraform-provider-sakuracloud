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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	registryBuilder "github.com/sacloud/iaas-service-go/containerregistry/builder"
)

func expandContainerRegistryBuilder(d *schema.ResourceData, client *APIClient, settingsHash string) *registryBuilder.Builder {
	return &registryBuilder.Builder{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Tags:           expandTags(d),
		IconID:         expandSakuraCloudID(d, "icon_id"),
		AccessLevel:    types.EContainerRegistryAccessLevel(d.Get("access_level").(string)),
		VirtualDomain:  stringOrDefault(d, "virtual_domain"),
		SubDomainLabel: d.Get("subdomain_label").(string),
		Users:          expandContainerRegistryUsers(d),
		SettingsHash:   settingsHash,
		Client:         iaas.NewContainerRegistryOp(client),
	}
}

func expandContainerRegistryUsers(d *schema.ResourceData) []*registryBuilder.User {
	var results []*registryBuilder.User
	users := d.Get("user").([]interface{})
	for _, raw := range users {
		d := mapToResourceData(raw.(map[string]interface{}))
		results = append(results, &registryBuilder.User{
			UserName:   stringOrDefault(d, "name"),
			Password:   stringOrDefault(d, "password"),
			Permission: types.EContainerRegistryPermission(stringOrDefault(d, "permission")),
		})
	}
	return results
}

func flattenContainerRegistryUsers(d *schema.ResourceData, users []*iaas.ContainerRegistryUser, includePassword bool) []interface{} {
	inputs := expandContainerRegistryUsers(d)

	var results []interface{}
	for _, user := range users {
		v := map[string]interface{}{
			"name":       user.UserName,
			"permission": user.Permission,
		}
		if includePassword {
			password := ""
			for _, i := range inputs {
				if i.UserName == user.UserName {
					password = i.Password
					break
				}
			}
			v["password"] = password
		}
		results = append(results, v)
	}
	return results
}
