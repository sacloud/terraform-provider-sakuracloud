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
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	envWebAccelSiteName             = "SAKURACLOUD_WEBACCEL_SITE_NAME"
	envWebAccelDomainName           = "SAKURACLOUD_WEBACCEL_DOMAIN_NAME"
	envWebAccelOrigin               = "SAKURACLOUD_WEBACCEL_ORIGIN"
	envObjectStorageEndpoint        = "SAKURASTORAGE_ENDPOINT"
	envObjectStorageRegion          = "SAKURASTORAGE_REGION"
	envObjectStorageBucketName      = "SAKURASTORAGE_BUCKET_NAME"
	envObjectStorageAccessKeyId     = "SAKURASTORAGE_ACCESS_KEY"
	envObjectStorageSecretAccessKey = "SAKURASTORAGE_ACCESS_SECRET"
)

func TestAccSakuraCloudDataSourceWebAccel_ByName(t *testing.T) {
	var siteName string
	if name, ok := os.LookupEnv(envWebAccelSiteName); ok {
		siteName = name
	} else {
		t.Skipf("ENV %q is requilred. skip", envWebAccelSiteName)
		return
	}

	regexpNotEmpty := regexp.MustCompile(".+")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceWebAccelWithName(siteName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudWebAccelDataSourceID("data.sakuracloud_webaccel.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_webaccel.foobar", "name", siteName),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "domain", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "origin", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "request_protocol", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "origin_parameters.0.type", regexpNotEmpty),
					//resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "logging.0.enabled", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "subdomain", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "domain_type", regexpNotEmpty),
					//resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "has_certificate", regexpNotEmpty),
					//resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "host_header", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "status", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "cname_record_value", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "txt_record_value", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "vary_support", regexpNotEmpty),
				),
			},
		},
	})
}

func TestAccSakuraCloudDataSourceWebAccel_ByDomain(t *testing.T) {
	var domainName string
	if name, ok := os.LookupEnv(envWebAccelDomainName); ok {
		domainName = name
	} else {
		t.Skipf("ENV %q is requilred. skip", envWebAccelDomainName)
		return
	}

	regexpNotEmpty := regexp.MustCompile(".+")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceWebAccelWithDomain(domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudWebAccelDataSourceID("data.sakuracloud_webaccel.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_webaccel.foobar", "domain", domainName),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "name", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "origin", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "origin_parameters.0.type", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "subdomain", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "domain_type", regexpNotEmpty),
					//resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "has_certificate", regexpNotEmpty),
					//resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "host_header", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "status", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "cname_record_value", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "txt_record_value", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "vary_support", regexpNotEmpty),
					resource.TestMatchResourceAttr("data.sakuracloud_webaccel.foobar", "normalize_ae", regexpNotEmpty),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudWebAccelDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("WebAccel data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudDataSourceWebAccelWithName(siteName string) string {
	tmpl := `
data "sakuracloud_webaccel" "foobar" {
  name = "%s"
}`
	return fmt.Sprintf(tmpl, siteName)
}

func testAccCheckSakuraCloudDataSourceWebAccelWithDomain(domain string) string {
	tmpl := `
data "sakuracloud_webaccel" "foobar" {
  domain = "%s"
}`
	return fmt.Sprintf(tmpl, domain)
}
