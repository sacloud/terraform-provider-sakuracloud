package sakuracloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSakuraCloudBucketObjectDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudNoteDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceBucketObject,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.sakuracloud_bucket_object.foobar", "key", "foo/bar/test.txt"),
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
						"http_url", "http://terraform-for-sakuracloud-test.b.sakurastorage.jp/foo/bar/test.txt"),
					resource.TestCheckResourceAttr(
						"data.sakuracloud_bucket_object.foobar",
						"https_url", "https://terraform-for-sakuracloud-test.b.sakurastorage.jp/foo/bar/test.txt"),
					resource.TestCheckResourceAttr(
						"data.sakuracloud_bucket_object.foobar",
						"http_path_url", "http://b.sakurastorage.jp/terraform-for-sakuracloud-test/foo/bar/test.txt"),
					resource.TestCheckResourceAttr(
						"data.sakuracloud_bucket_object.foobar",
						"http_cache_url", "http://terraform-for-sakuracloud-test.c.sakurastorage.jp/foo/bar/test.txt"),
					resource.TestCheckResourceAttr(
						"data.sakuracloud_bucket_object.foobar",
						"https_cache_url", "https://terraform-for-sakuracloud-test.c.sakurastorage.jp/foo/bar/test.txt"),
				),
			},
		},
	})
}

var testAccCheckSakuraCloudDataSourceBucketObject = `
resource "sakuracloud_bucket_object" "foobar" {
  bucket  = "terraform-for-sakuracloud-test"
  key     = "foo/bar/test.txt"
  content = "content"
}

data "sakuracloud_bucket_object" "foobar" {
  bucket = "${sakuracloud_bucket_object.foobar.bucket}"
  key    = "${sakuracloud_bucket_object.foobar.key}"
}
`
