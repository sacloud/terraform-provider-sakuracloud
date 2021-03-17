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
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mitchellh/go-homedir"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func expandIconBody(d resourceValueGettable) (string, error) {
	var body string
	if v, ok := d.GetOk("source"); ok {
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return "", fmt.Errorf("expanding homedir in source (%s) is failed: %s", source, err)
		}
		file, err := os.Open(path)
		if err != nil {
			return "", fmt.Errorf("opening SakuraCloud Icon source(%s) is failed: %s", source, err)
		}
		data, err := io.ReadAll(file)
		if err != nil {
			return "", fmt.Errorf("reading SakuraCloud Icon source file is failed: %s", err)
		}
		body = base64.StdEncoding.EncodeToString(data)
	} else if v, ok := d.GetOk("base64content"); ok {
		body = v.(string)
	} else {
		return "", fmt.Errorf(`"source" or "base64content" field is required`)
	}
	return body, nil
}

func expandIconCreateRequest(d *schema.ResourceData) (*sacloud.IconCreateRequest, error) {
	body, err := expandIconBody(d)
	if err != nil {
		return nil, fmt.Errorf("creating SakuraCloud Icon is failed: %s", err)
	}
	return &sacloud.IconCreateRequest{
		Name:  d.Get("name").(string),
		Tags:  expandTags(d),
		Image: body,
	}, nil
}

func expandIconUpdateRequest(d *schema.ResourceData) *sacloud.IconUpdateRequest {
	return &sacloud.IconUpdateRequest{
		Name: d.Get("name").(string),
		Tags: expandTags(d),
	}
}
