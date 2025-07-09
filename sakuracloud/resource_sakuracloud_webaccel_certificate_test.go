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
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	envWebAccelCertificateCrt    = "SAKURACLOUD_WEBACCEL_CERT_PATH"
	envWebAccelCertificateKey    = "SAKURACLOUD_WEBACCEL_KEY_PATH"
	envWebAccelCertificateCrtUpd = "SAKURACLOUD_WEBACCEL_CERT_PATH_UPD"
	envWebAccelCertificateKeyUpd = "SAKURACLOUD_WEBACCEL_KEY_PATH_UPD"
)

func TestAccResourceSakuraCloudWebAccelCertificate_basic(t *testing.T) {
	envKeys := []string{
		envWebAccelSiteName,
		envWebAccelCertificateCrt,
		envWebAccelCertificateKey,
		envWebAccelCertificateCrtUpd,
		envWebAccelCertificateKeyUpd,
	}
	for _, k := range envKeys {
		if os.Getenv(k) == "" {
			t.Skipf("ENV %q is requilred. skip", k)
			return
		}
	}

	siteName := os.Getenv(envWebAccelSiteName)
	crt := os.Getenv(envWebAccelCertificateCrt)
	key := os.Getenv(envWebAccelCertificateKey)
	crtUpd := os.Getenv(envWebAccelCertificateCrtUpd)
	keyUpd := os.Getenv(envWebAccelCertificateKeyUpd)

	regexpNotEmpty := regexp.MustCompile(".+")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudWebAccelCertificateConfig(siteName, crt, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("sakuracloud_webaccel_certificate.foobar", "id", regexpNotEmpty),
					resource.TestMatchResourceAttr("sakuracloud_webaccel_certificate.foobar", "site_id", regexpNotEmpty),
					resource.TestMatchResourceAttr("sakuracloud_webaccel_certificate.foobar", "not_before", regexpNotEmpty),
					resource.TestMatchResourceAttr("sakuracloud_webaccel_certificate.foobar", "not_after", regexpNotEmpty),
					resource.TestMatchResourceAttr("sakuracloud_webaccel_certificate.foobar", "issuer_common_name", regexpNotEmpty),
					resource.TestMatchResourceAttr("sakuracloud_webaccel_certificate.foobar", "subject_common_name", regexpNotEmpty),
					resource.TestMatchResourceAttr("sakuracloud_webaccel_certificate.foobar", "sha256_fingerprint", regexpNotEmpty),
				),
			},
			{
				Config: testAccCheckSakuraCloudWebAccelCertificateConfig(siteName, crtUpd, keyUpd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("sakuracloud_webaccel_certificate.foobar", "id", regexpNotEmpty),
					resource.TestMatchResourceAttr("sakuracloud_webaccel_certificate.foobar", "site_id", regexpNotEmpty),
					resource.TestMatchResourceAttr("sakuracloud_webaccel_certificate.foobar", "not_before", regexpNotEmpty),
					resource.TestMatchResourceAttr("sakuracloud_webaccel_certificate.foobar", "not_after", regexpNotEmpty),
					resource.TestMatchResourceAttr("sakuracloud_webaccel_certificate.foobar", "issuer_common_name", regexpNotEmpty),
					resource.TestMatchResourceAttr("sakuracloud_webaccel_certificate.foobar", "subject_common_name", regexpNotEmpty),
					resource.TestMatchResourceAttr("sakuracloud_webaccel_certificate.foobar", "sha256_fingerprint", regexpNotEmpty),
				),
			},
		},
	})
}

func TestAccResourceSakuraCloudWebAccelCertificate_LetsEncrypt(t *testing.T) {
	envKeys := []string{
		envWebAccelSiteName,
	}
	for _, k := range envKeys {
		if os.Getenv(k) == "" {
			t.Skipf("ENV %q is requilred. skip", k)
			return
		}
	}

	siteName := os.Getenv(envWebAccelSiteName)
	regexpNotEmpty := regexp.MustCompile(".+")
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckSakuraCloudWebAccelCertificateFreeCertConfig(siteName, false),
				ExpectError: regexp.MustCompile("must be true for the creation"),
				Destroy:     false,
			},
			{
				Config: testAccCheckSakuraCloudWebAccelCertificateFreeCertConfig(siteName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("sakuracloud_webaccel_certificate.foobar", "id", regexpNotEmpty),
					resource.TestMatchResourceAttr("sakuracloud_webaccel_certificate.foobar", "site_id", regexpNotEmpty),
					resource.TestCheckResourceAttr("sakuracloud_webaccel_certificate.foobar", "lets_encrypt", "true"),
				),
				Destroy: false,
			},
			{
				Config:      testAccCheckSakuraCloudWebAccelCertificateFreeCertConfig(siteName, false),
				ExpectError: regexp.MustCompile("must not be false"),
				Destroy:     false,
			},
			{
				Config:  testAccCheckSakuraCloudWebAccelCertificateFreeCertConfig(siteName, true),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudWebAccelCertificateConfig(siteName, crt, key string) string {
	tmpl := `
data sakuracloud_webaccel "site" {
  name = "%s"
}
resource sakuracloud_webaccel_certificate "foobar" {
  site_id           = data.sakuracloud_webaccel.site.id
  certificate_chain = file("%s") 
  private_key       = file("%s")
}
`
	return fmt.Sprintf(tmpl, siteName, crt, key)
}

func testAccCheckSakuraCloudWebAccelCertificateFreeCertConfig(siteName string, enabled bool) string {
	tmpl := `
data sakuracloud_webaccel "site" {
  name = "%s"
}
resource sakuracloud_webaccel_certificate "foobar" {
  site_id      = data.sakuracloud_webaccel.site.id
  lets_encrypt = %s
}
`
	if enabled {
		return fmt.Sprintf(tmpl, siteName, "true")
	} else {
		return fmt.Sprintf(tmpl, siteName, "false")
	}
}
