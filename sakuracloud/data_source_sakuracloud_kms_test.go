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
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	v1 "github.com/sacloud/kms-api-go/apis/v1"
)

func TestAccSakuraCloudDataSourceKMS_basic(t *testing.T) {
	skipIfFakeModeEnabled(t)

	resourceName := "data.sakuracloud_kms.foobar"
	rand := randomName()

	var key v1.Key
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceKMS_byName, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudKMSExists("sakuracloud_kms.foobar", &key),
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "key_origin", "generated"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceKMS_byResourceId, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudKMSExists("sakuracloud_kms.foobar", &key),
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "key_origin", "generated"),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceKMS_byName = `
resource "sakuracloud_kms" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  tags        = ["tag1", "tag2"]
}

data "sakuracloud_kms" "foobar" {
  name = "{{ .arg0 }}"

  # NOTE: resourceを先に作らせてから参照するために依存関係を明示
  depends_on = [
    sakuracloud_kms.foobar
  ]
}`

var testAccSakuraCloudDataSourceKMS_byResourceId = `
resource "sakuracloud_kms" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  tags        = ["tag1", "tag2"]
}

data "sakuracloud_kms" "foobar" {
  resource_id = sakuracloud_kms.foobar.id

  # NOTE: resourceを先に作らせてから参照するために依存関係を明示
  depends_on = [
    sakuracloud_kms.foobar
  ]
}`

func TestFilterKMSByName(t *testing.T) {
	t.Parallel()

	keys := []v1.Key{
		{
			Name: "test-key1",
			Tags: []string{"tag1"},
		},
		{
			Name: "test-key2",
			Tags: []string{"tag1", "tag2"},
		},
	}

	testCases := []struct {
		name    string
		keyName string
		want    *v1.Key
		wantErr bool
	}{
		{
			name:    "found by name",
			keyName: "test-key1",
			want:    &keys[0],
		},
		{
			name:    "not found",
			keyName: "not-exist",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := filterKMSByName(keys, tc.keyName)
			if tc.wantErr && err == nil {
				t.Errorf("filterKMSByName wants error but got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("filterKMSByName error = %v", err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("filterKMSByName got = %v, want %v", got, tc.want)
			}
		})
	}
}
