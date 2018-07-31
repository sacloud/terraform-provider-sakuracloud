package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

func TestAccResourceSakuraCloudBucketObject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudBucketObjectConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBucketObjectExists("sakuracloud_bucket_object.foobar"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar", "key", "foo/bar/test.txt"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar", "size", "7"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar", "content_type", "text/plain"),
					// etag = `echo -n content | md5`
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar", "etag", "9a0364b9e99bb480dd25e1f0284c8555"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar",
						"http_url", "http://terraform-for-sakuracloud-test.b.sakurastorage.jp/foo/bar/test.txt"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar",
						"https_url", "https://terraform-for-sakuracloud-test.b.sakurastorage.jp/foo/bar/test.txt"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar",
						"http_path_url", "http://b.sakurastorage.jp/terraform-for-sakuracloud-test/foo/bar/test.txt"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar",
						"https_path_url", "https://b.sakurastorage.jp/terraform-for-sakuracloud-test/foo/bar/test.txt"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar",
						"http_cache_url", "http://terraform-for-sakuracloud-test.c.sakurastorage.jp/foo/bar/test.txt"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar",
						"https_cache_url", "https://terraform-for-sakuracloud-test.c.sakurastorage.jp/foo/bar/test.txt"),
				),
			},
			{
				Config: testAccCheckSakuraCloudBucketObjectConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBucketObjectExists("sakuracloud_bucket_object.foobar"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar", "key", "foo/bar/test.txt"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar", "size", "11"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar", "etag", "63438b4e5a535fd413b24cdc3e380f3d"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudBucketObjectExists(n string) resource.TestCheckFunc {
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

func testAccCheckSakuraCloudBucketObjectDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_bucket_object" {
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

const testAccCheckSakuraCloudBucketObjectConfig_basic = `
resource "sakuracloud_bucket_object" "foobar" {
  key = "foo/bar/test.txt"
  bucket = "terraform-for-sakuracloud-test"
  content = "content"
}`

const testAccCheckSakuraCloudBucketObjectConfig_update = `
resource "sakuracloud_bucket_object" "foobar" {
  key = "foo/bar/test.txt"
  bucket = "terraform-for-sakuracloud-test"
  content = "content-upd"
}`
