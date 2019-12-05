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

func TestAccResourceSakuraCloudBucketObject(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	randString3 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	key := fmt.Sprintf("%s/%s/%s.txt", randString1, randString2, randString3)
	bucket := os.Getenv("SACLOUD_OJS_ACCESS_KEY_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudBucketObjectConfig_basic(bucket, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBucketObjectExists("sakuracloud_bucket_object.foobar"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar", "key", key),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar", "size", "7"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar", "content_type", "text/plain"),
					// etag = `echo -n content | md5`
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar", "etag", "9a0364b9e99bb480dd25e1f0284c8555"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar",
						"http_url", fmt.Sprintf("http://%s.b.sakurastorage.jp/%s", bucket, key)),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar",
						"https_url", fmt.Sprintf("https://%s.b.sakurastorage.jp/%s", bucket, key)),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar",
						"http_path_url", fmt.Sprintf("http://b.sakurastorage.jp/%s/%s", bucket, key)),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar",
						"https_path_url", fmt.Sprintf("https://b.sakurastorage.jp/%s/%s", bucket, key)),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar",
						"http_cache_url", fmt.Sprintf("http://%s.c.sakurastorage.jp/%s", bucket, key)),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar",
						"https_cache_url", fmt.Sprintf("https://%s.c.sakurastorage.jp/%s", bucket, key)),
				),
			},
			{
				Config: testAccCheckSakuraCloudBucketObjectConfig_update(bucket, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudBucketObjectExists("sakuracloud_bucket_object.foobar"),
					resource.TestCheckResourceAttr(
						"sakuracloud_bucket_object.foobar", "key", key),
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

func testAccCheckSakuraCloudBucketObjectConfig_basic(bucket, key string) string {
	return fmt.Sprintf(`
resource "sakuracloud_bucket_object" "foobar" {
  bucket = "%s"
  key = "%s"
  content = "content"
}`, bucket, key)
}

func testAccCheckSakuraCloudBucketObjectConfig_update(bucket, key string) string {
	return fmt.Sprintf(`
resource "sakuracloud_bucket_object" "foobar" {
  bucket = "%s"
  key = "%s"
  content = "content-upd"
}`, bucket, key)
}
