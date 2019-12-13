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
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

func TestAccSakuraCloudBucketObject_basic(t *testing.T) {
	skipIfFakeModeEnabled(t)

	resourceName := "sakuracloud_bucket_object.foobar"

	rand1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	rand2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	rand3 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	key := fmt.Sprintf("%s/%s/%s.txt", rand1, rand2, rand3)
	bucket := os.Getenv("SACLOUD_OJS_ACCESS_KEY_ID")

	httpURL := fmt.Sprintf("http://%s.b.sakurastorage.jp/%s", bucket, key)
	httpsURL := fmt.Sprintf("https://%s.b.sakurastorage.jp/%s", bucket, key)
	httpPathURL := fmt.Sprintf("http://b.sakurastorage.jp/%s/%s", bucket, key)
	httpsPathURL := fmt.Sprintf("https://b.sakurastorage.jp/%s/%s", bucket, key)
	httpCacheURL := fmt.Sprintf("http://%s.c.sakurastorage.jp/%s", bucket, key)
	httpsCacheURL := fmt.Sprintf("https://%s.c.sakurastorage.jp/%s", bucket, key)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckSakuraCloudBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudBucketObject_basic, bucket, key),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudBucketObjectExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "size", "7"),
					resource.TestCheckResourceAttr(resourceName, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(resourceName, "etag", "9a0364b9e99bb480dd25e1f0284c8555"), // etag = `echo -n content | md5`
					resource.TestCheckResourceAttr(resourceName, "http_url", httpURL),
					resource.TestCheckResourceAttr(resourceName, "https_url", httpsURL),
					resource.TestCheckResourceAttr(resourceName, "http_path_url", httpPathURL),
					resource.TestCheckResourceAttr(resourceName, "https_path_url", httpsPathURL),
					resource.TestCheckResourceAttr(resourceName, "http_cache_url", httpCacheURL),
					resource.TestCheckResourceAttr(resourceName, "https_cache_url", httpsCacheURL),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudBucketObject_update, bucket, key),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudBucketObjectExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "size", "11"),
					resource.TestCheckResourceAttr(resourceName, "etag", "63438b4e5a535fd413b24cdc3e380f3d"),
				),
			},
		},
	})
}

func testCheckSakuraCloudBucketObjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No BucketObject ID is set")
		}

		return nil
	}
}

func testCheckSakuraCloudBucketObjectDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_bucket_object" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		accessKey := rs.Primary.Attributes["access_key"]
		secretKey := rs.Primary.Attributes["secret_key"]
		strBucket := rs.Primary.Attributes["bucket"]

		auth, err := aws.GetAuth(accessKey, secretKey)
		if err != nil {
			return err
		}
		client := s3.New(auth, aws.Region{
			Name:       "us-west-2",
			S3Endpoint: "https://b.sakurastorage.jp",
		})
		bucket := client.Bucket(strBucket)

		_, err = bucket.GetKey(rs.Primary.ID)
		if err == nil {
			return errors.New("BucketObject still exists")
		}
	}

	return nil
}

var testAccSakuraCloudBucketObject_basic = `
resource "sakuracloud_bucket_object" "foobar" {
  bucket  = "{{ .arg0 }}"
  key     = "{{ .arg1 }}"
  content = "content"
}`

var testAccSakuraCloudBucketObject_update = `
resource "sakuracloud_bucket_object" "foobar" {
  bucket  = "{{ .arg0 }}"
  key     = "{{ .arg1 }}"
  content = "content-upd"
}`
