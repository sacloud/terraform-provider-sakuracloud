package sakuracloud

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	envWebAccelCertificateCrt    = "SAKURACLOUD_WEBACCEL_CERT_PATH"
	envWebAccelCertificateKey    = "SAKURACLOUD_WEBACCEL_KEY_PATH"
	envWebAccelCertificateCrtUpd = "SAKURACLOUD_WEBACCEL_CERT_PATH_UPD"
	envWebAccelCertificateKeyUpd = "SAKURACLOUD_WEBACCEL_KEY_PATH_UPD"
)

func TestAccResourceSakuraCloudWebAccelCertificate(t *testing.T) {
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
