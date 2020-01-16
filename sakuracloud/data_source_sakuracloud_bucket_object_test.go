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
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSakuraCloudBucketObjectDataSource_Basic(t *testing.T) {
	skipIfFakeModeEnabled(t)

	resourceName := "data.sakuracloud_bucket_object.foobar"

	rand1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	rand2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	rand3 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	key := fmt.Sprintf("%s/%s/terraform-acctest-%s.txt", rand1, rand2, rand3)
	bucket := os.Getenv("SACLOUD_OJS_ACCESS_KEY_ID")

	httpURL := fmt.Sprintf("http://%s.b.sakurastorage.jp/%s", bucket, key)
	httpsURL := fmt.Sprintf("https://%s.b.sakurastorage.jp/%s", bucket, key)
	httpPathURL := fmt.Sprintf("http://b.sakurastorage.jp/%s/%s", bucket, key)
	httpsPathURL := fmt.Sprintf("https://b.sakurastorage.jp/%s/%s", bucket, key)
	httpCacheURL := fmt.Sprintf("http://%s.c.sakurastorage.jp/%s", bucket, key)
	httpsCacheURL := fmt.Sprintf("https://%s.c.sakurastorage.jp/%s", bucket, key)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceBucketObject, bucket, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "size", "7"),
					resource.TestCheckResourceAttr(resourceName, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(resourceName, "body", "content"),
					resource.TestCheckResourceAttr(resourceName, "etag", "9a0364b9e99bb480dd25e1f0284c8555"),
					resource.TestCheckResourceAttr(resourceName, "http_url", httpURL),
					resource.TestCheckResourceAttr(resourceName, "https_url", httpsURL),
					resource.TestCheckResourceAttr(resourceName, "http_path_url", httpPathURL),
					resource.TestCheckResourceAttr(resourceName, "https_path_url", httpsPathURL),
					resource.TestCheckResourceAttr(resourceName, "http_cache_url", httpCacheURL),
					resource.TestCheckResourceAttr(resourceName, "https_cache_url", httpsCacheURL),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceBucketObject = `
resource "sakuracloud_bucket_object" "foobar" {
  bucket  = "{{ .arg0 }}"
  key     = "{{ .arg1 }}"
  content = "content"
}

data "sakuracloud_bucket_object" "foobar" {
  bucket = sakuracloud_bucket_object.foobar.bucket
  key    = sakuracloud_bucket_object.foobar.key
}
`
