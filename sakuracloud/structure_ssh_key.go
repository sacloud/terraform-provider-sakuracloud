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
)

func expandSSHKeyCreateRequest(d *schema.ResourceData) *iaas.SSHKeyCreateRequest {
	return &iaas.SSHKeyCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		PublicKey:   d.Get("public_key").(string),
	}
}

func expandSSHKeyUpdateRequest(d *schema.ResourceData) *iaas.SSHKeyUpdateRequest {
	return &iaas.SSHKeyUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
}

func expandSSHKeyGenerateRequest(d *schema.ResourceData) *iaas.SSHKeyGenerateRequest {
	return &iaas.SSHKeyGenerateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		PassPhrase:  d.Get("pass_phrase").(string),
	}
}
