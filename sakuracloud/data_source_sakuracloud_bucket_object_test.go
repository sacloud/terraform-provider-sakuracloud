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
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSakuraCloudBucketObjectDataSource_Basic(t *testing.T) {

	randString1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	randString3 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	key := fmt.Sprintf("%s/%s/%s.txt", randString1, randString2, randString3)
	bucket := os.Getenv("SACLOUD_OJS_ACCESS_KEY_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudNoteDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceBucketObject(bucket, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.sakuracloud_bucket_object.foobar", "key", key),
					resource.TestCheckResourceAttr(
						"data.sakuracloud_bucket_object.foobar", "size", "7"),
					resource.TestCheckResourceAttr(
						"data.sakuracloud_bucket_object.foobar", "content_type", "text/plain"),
					resource.TestCheckResourceAttr(
						"data.sakuracloud_bucket_object.foobar", "body", "content"),
					resource.TestCheckResourceAttr(
						"data.sakuracloud_bucket_object.foobar", "etag", "9a0364b9e99bb480dd25e1f0284c8555"),
					resource.TestCheckResourceAttr(
						"data.sakuracloud_bucket_object.foobar",
						"http_url", fmt.Sprintf("http://%s.b.sakurastorage.jp/%s", bucket, key)),
					resource.TestCheckResourceAttr(
						"data.sakuracloud_bucket_object.foobar",
						"https_url", fmt.Sprintf("https://%s.b.sakurastorage.jp/%s", bucket, key)),
					resource.TestCheckResourceAttr(
						"data.sakuracloud_bucket_object.foobar",
						"http_path_url", fmt.Sprintf("http://b.sakurastorage.jp/%s/%s", bucket, key)),
					resource.TestCheckResourceAttr(
						"data.sakuracloud_bucket_object.foobar",
						"https_path_url", fmt.Sprintf("https://b.sakurastorage.jp/%s/%s", bucket, key)),
					resource.TestCheckResourceAttr(
						"data.sakuracloud_bucket_object.foobar",
						"http_cache_url", fmt.Sprintf("http://%s.c.sakurastorage.jp/%s", bucket, key)),
					resource.TestCheckResourceAttr(
						"data.sakuracloud_bucket_object.foobar",
						"https_cache_url", fmt.Sprintf("https://%s.c.sakurastorage.jp/%s", bucket, key)),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudDataSourceBucketObject(bucket, key string) string {
	return fmt.Sprintf(`
resource "sakuracloud_bucket_object" "foobar" {
  bucket  = "%s"
  key     = "%s"
  content = "content"
}

data "sakuracloud_bucket_object" "foobar" {
  bucket = "${sakuracloud_bucket_object.foobar.bucket}"
  key    = "${sakuracloud_bucket_object.foobar.key}"
}
`, bucket, key)
}
