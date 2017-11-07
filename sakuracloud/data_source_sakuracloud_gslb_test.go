package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"testing"
)

func TestAccSakuraCloudGSLBDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudGSLBDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceGSLBBase,
				Check:  testAccCheckSakuraCloudGSLBDataSourceID("sakuracloud_gslb.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceGSLBConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBDataSourceID("data.sakuracloud_gslb.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "sorry_server", "8.8.8.8"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "health_check.0.protocol", "http"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "health_check.0.delay_loop", "10"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "health_check.0.host_header", "terraform.io"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_gslb.foobar", "tags.2", "tag3"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceGSLBConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBDataSourceID("data.sakuracloud_gslb.foobar"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceGSLBConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBDataSourceNotExists("data.sakuracloud_gslb.foobar"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceGSLBConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudGSLBDataSourceNotExists("data.sakuracloud_gslb.foobar"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudGSLBDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find GSLB data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("GSLB data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudGSLBDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found GSLB data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudGSLBDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_gslb" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.GSLB.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("GSLB still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceGSLBBase = `
resource "sakuracloud_gslb" "foobar" {
    name = "name_test"
    health_check = {
        protocol = "http"
        delay_loop = 10
        host_header = "terraform.io"
        path = "/"
        status = "200"
    }
    sorry_server = "8.8.8.8"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}`

var testAccCheckSakuraCloudDataSourceGSLBConfig = `
resource "sakuracloud_gslb" "foobar" {
    name = "name_test"
    health_check = {
        protocol = "http"
        delay_loop = 10
        host_header = "terraform.io"
        path = "/"
        status = "200"
    }
    sorry_server = "8.8.8.8"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_gslb" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
}`

var testAccCheckSakuraCloudDataSourceGSLBConfig_With_Tag = `
resource "sakuracloud_gslb" "foobar" {
    name = "name_test"
    health_check = {
        protocol = "http"
        delay_loop = 10
        host_header = "terraform.io"
        path = "/"
        status = "200"
    }
    sorry_server = "8.8.8.8"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_gslb" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
}`

var testAccCheckSakuraCloudDataSourceGSLBConfig_With_NotExists_Tag = `
resource "sakuracloud_gslb" "foobar" {
    name = "name_test"
    health_check = {
        protocol = "http"
        delay_loop = 10
        host_header = "terraform.io"
        path = "/"
        status = "200"
    }
    sorry_server = "8.8.8.8"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_gslb" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
}`

var testAccCheckSakuraCloudDataSourceGSLBConfig_NotExists = `
resource "sakuracloud_gslb" "foobar" {
    name = "name_test"
    health_check = {
        protocol = "http"
        delay_loop = 10
        host_header = "terraform.io"
        path = "/"
        status = "200"
    }
    sorry_server = "8.8.8.8"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_gslb" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`
