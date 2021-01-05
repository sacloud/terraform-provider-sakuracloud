// Copyright 2016-2021 terraform-provider-sakuracloud authors
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
	"github.com/sacloud/libsacloud/v2/pkg/size"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func expandCDROMContentHash(d *schema.ResourceData) string {
	// NOTE 本来はAPIにてmd5ハッシュを取得できるのが望ましい。現状ではここでファイルを読んで算出する。
	if v, ok := d.GetOk("iso_image_file"); ok {
		source := v.(string)

		path, err := expandHomeDir(source)
		if err != nil {
			return ""
		}
		hash, err := md5CheckSumFromFile(path)
		if err != nil {
			return ""
		}
		return hash
	}
	return ""
}

func expandCDROMCreateRequest(d *schema.ResourceData) *sacloud.CDROMCreateRequest {
	return &sacloud.CDROMCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		SizeMB:      d.Get("size").(int) * size.GiB,
		IconID:      expandSakuraCloudID(d, "icon_id"),
		Tags:        expandTags(d),
	}
}

func expandCDROMUpdateRequest(d *schema.ResourceData) *sacloud.CDROMUpdateRequest {
	return &sacloud.CDROMUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),
	}
}
