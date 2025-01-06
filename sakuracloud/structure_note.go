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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/iaas-api-go"
)

func expandNoteCreateRequest(d *schema.ResourceData) *iaas.NoteCreateRequest {
	return &iaas.NoteCreateRequest{
		Name:    d.Get("name").(string),
		Tags:    expandTags(d),
		IconID:  expandSakuraCloudID(d, "icon_id"),
		Class:   d.Get("class").(string),
		Content: d.Get("content").(string),
	}
}

func expandNoteUpdateRequest(d *schema.ResourceData) *iaas.NoteUpdateRequest {
	return &iaas.NoteUpdateRequest{
		Name:    d.Get("name").(string),
		Tags:    expandTags(d),
		IconID:  expandSakuraCloudID(d, "icon_id"),
		Class:   d.Get("class").(string),
		Content: d.Get("content").(string),
	}
}
