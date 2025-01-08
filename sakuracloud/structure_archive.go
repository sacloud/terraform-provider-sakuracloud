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
	"io"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	archiveUtil "github.com/sacloud/iaas-service-go/archive/builder"
)

func expandArchiveBuilder(d *schema.ResourceData, zone string, client *APIClient) (archiveUtil.Builder, func(), error) {
	var reader io.ReadCloser
	source := d.Get("archive_file").(string)
	if source != "" {
		sourcePath, err := expandHomeDir(source)
		if err != nil {
			return nil, nil, err
		}
		f, err := os.Open(sourcePath)
		if err != nil {
			return nil, nil, err
		}
		reader = f
	}

	sourceArchiveZone := stringOrDefault(d, "source_archive_zone")
	if sourceArchiveZone != "" {
		if _, errs := validation.StringInSlice(client.zones, false)(sourceArchiveZone, "source_archive_zone"); len(errs) > 0 {
			return nil, nil, errs[0]
		}
		if zone == sourceArchiveZone {
			sourceArchiveZone = ""
		}
	}
	sizeGB := intOrDefault(d, "size")
	if sizeGB == 0 {
		sizeGB = 20
	}

	// Note: APIとしてはディスクやアーカイブをソースとした場合Sizeの指定はできないが、
	//       archiveUtil.Director側でAPIに渡すパラメータを制御しているためここでは常に渡して問題ない
	director := &archiveUtil.Director{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		Tags:              expandTags(d),
		IconID:            expandSakuraCloudID(d, "icon_id"),
		SizeGB:            sizeGB,
		SourceReader:      reader,
		SourceDiskID:      expandSakuraCloudID(d, "source_disk_id"),
		SourceArchiveID:   expandSakuraCloudID(d, "source_archive_id"),
		SourceArchiveZone: sourceArchiveZone,
		SourceSharedKey:   types.ArchiveShareKey(stringOrDefault(d, "source_shared_key")),
		Client:            archiveUtil.NewAPIClient(client),
	}
	return director.Builder(), func() {
		if reader != nil {
			reader.Close() // nolint
		}
	}, nil
}

func expandArchiveHash(d *schema.ResourceData) string {
	source := d.Get("archive_file").(string)
	if source == "" {
		return ""
	}

	path, err := expandHomeDir(source)
	if err != nil {
		return ""
	}

	// NOTE 本来はAPIにてmd5ハッシュを取得できるのが望ましい。現状ではここでファイルを読んで算出する。
	hash, err := md5CheckSumFromFile(path)
	if err != nil {
		return ""
	}
	return hash
}

func expandArchiveUpdateRequest(d *schema.ResourceData) *iaas.ArchiveUpdateRequest {
	return &iaas.ArchiveUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),
	}
}
