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
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSakuraCloudResourceWebAccel_WebOrigin(t *testing.T) {
	skipIfFakeModeEnabled(t)

	envKeys := []string{
		envWebAccelOrigin,
	}
	for _, k := range envKeys {
		if os.Getenv(k) == "" {
			t.Skipf("ENV %q is requilred. skip", k)
			return
		}
	}

	siteName := "your-site-name"
	// domainName := os.Getenv(envWebAccelDomainName)
	origin := os.Getenv(envWebAccelOrigin)
	regexpNotEmpty := regexp.MustCompile(".+")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudWebAccelWebOriginConfigBasic(siteName, origin),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "name", siteName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "request_protocol", "https-redirect"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.type", "web"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.origin", origin),
					resource.TestMatchResourceAttr("sakuracloud_webaccel.foobar", "cname_record_value", regexpNotEmpty),
					resource.TestMatchResourceAttr("sakuracloud_webaccel.foobar", "txt_record_value", regexpNotEmpty),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "normalize_ae", "br+gzip"),
				),
			},
		},
	})
}

func TestAccSakuraCloudResourceWebAccel_WebOriginWithOneTimeUrlSecrets(t *testing.T) {
	skipIfFakeModeEnabled(t)

	envKeys := []string{
		envWebAccelOrigin,
	}
	for _, k := range envKeys {
		if os.Getenv(k) == "" {
			t.Skipf("ENV %q is requilred. skip", k)
			return
		}
	}

	siteName := "your-site-name"
	origin := os.Getenv(envWebAccelOrigin)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudWebAccelWebOriginConfigWithOneTimeUrlSecrets(siteName, origin),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "name", siteName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.type", "web"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.origin", origin),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "onetime_url_secrets.0", "sample-secret"),
				),
			},
		},
	})
}

func TestAccSakuraCloudResourceWebAccel_WebOriginWithCORS(t *testing.T) {
	skipIfFakeModeEnabled(t)

	envKeys := []string{
		envWebAccelOrigin,
	}
	for _, k := range envKeys {
		if os.Getenv(k) == "" {
			t.Skipf("ENV %q is requilred. skip", k)
			return
		}
	}

	siteName := "your-site-name"
	// domainName := os.Getenv(envWebAccelDomainName)
	origin := os.Getenv(envWebAccelOrigin)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudWebAccelWebOriginConfigWithCors(siteName, origin),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "name", siteName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.type", "web"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.origin", origin),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "cors_rules.0.allow_all", "false"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "cors_rules.0.allowed_origins.0", "https://apps.example.com"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "onetime_url_secrets.0", "sample-secret"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "normalize_ae", "gzip"),
				),
			},
		},
	})
}

func TestAccSakuraCloudResourceWebAccel_Update(t *testing.T) {
	skipIfFakeModeEnabled(t)

	envKeys := []string{
		envWebAccelOrigin,
		envObjectStorageEndpoint,
		envObjectStorageRegion,
		envObjectStorageBucketName,
		envObjectStorageAccessKeyId,
		envObjectStorageSecretAccessKey,
	}
	for _, k := range envKeys {
		if os.Getenv(k) == "" {
			t.Skipf("ENV %q is requilred. skip", k)
			return
		}
	}

	siteName := "your-site-name"
	// domainName := os.Getenv(envWebAccelDomainName)
	origin := os.Getenv(envWebAccelOrigin)
	endpoint, _ := strings.CutPrefix(os.Getenv(envObjectStorageEndpoint), "https://")
	region := os.Getenv(envObjectStorageRegion)
	bucketName := os.Getenv(envObjectStorageBucketName)
	accessKey := os.Getenv(envObjectStorageAccessKeyId)
	secretKey := os.Getenv(envObjectStorageSecretAccessKey)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudWebAccelWebOriginConfigBasic(siteName, origin),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "name", siteName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.type", "web"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.origin", origin),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "request_protocol", "https-redirect"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "vary_support", "true"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "normalize_ae", "br+gzip"),
				),
			},
			{
				Config: testAccCheckSakuraCloudWebAccelWebOriginConfigWithCors(siteName, origin),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "name", siteName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.type", "web"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.origin", origin),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "request_protocol", "http+https"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "cors_rules.0.allow_all", "false"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "cors_rules.0.allowed_origins.0", "https://apps.example.com"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "normalize_ae", "gzip"),
				),
			},
			{
				Config: testAccCheckSakuraCloudWebAccelWebOriginLoggingConfig(siteName, origin, bucketName, accessKey, secretKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "name", siteName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.type", "web"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.origin", origin),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "logging.0.s3_bucket_name", bucketName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "vary_support", "true"),
				),
			},
			{
				Config: testAccCheckSakuraCloudWebAccelBucketOriginConfig(siteName, endpoint, region, bucketName, accessKey, secretKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "name", siteName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.type", "bucket"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.s3_endpoint", endpoint),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.s3_region", region),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.s3_bucket_name", bucketName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.s3_access_key_id", accessKey),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.s3_secret_access_key", secretKey),
				),
			},
		},
	})
}

func TestAccSakuraCloudResourceWebAccel_BucketOrigin(t *testing.T) {
	skipIfFakeModeEnabled(t)

	envKeys := []string{
		envWebAccelOrigin,
		envObjectStorageEndpoint,
		envObjectStorageRegion,
		envObjectStorageBucketName,
		envObjectStorageAccessKeyId,
		envObjectStorageSecretAccessKey,
	}
	for _, k := range envKeys {
		if os.Getenv(k) == "" {
			t.Skipf("ENV %q is requilred. skip", k)
			return
		}
	}

	siteName := "your-site-name"
	// domainName := os.Getenv(envWebAccelDomainName)
	endpoint, _ := strings.CutPrefix(os.Getenv(envObjectStorageEndpoint), "https://")
	region := os.Getenv(envObjectStorageRegion)
	bucketName := os.Getenv(envObjectStorageBucketName)
	accessKey := os.Getenv(envObjectStorageAccessKeyId)
	secretKey := os.Getenv(envObjectStorageSecretAccessKey)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudWebAccelBucketOriginConfig(siteName, endpoint, region, bucketName, accessKey, secretKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "name", siteName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.type", "bucket"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.s3_endpoint", endpoint),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.s3_region", region),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.s3_bucket_name", bucketName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.s3_access_key_id", accessKey),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.s3_secret_access_key", secretKey),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.s3_doc_index", "true"),
				),
			},
		},
	})
}

func TestAccSakuraCloudResourceWebAccel_Logging(t *testing.T) {
	skipIfFakeModeEnabled(t)

	envKeys := []string{
		envWebAccelOrigin,
		envObjectStorageBucketName,
		envObjectStorageAccessKeyId,
		envObjectStorageSecretAccessKey,
	}
	for _, k := range envKeys {
		if os.Getenv(k) == "" {
			t.Skipf("ENV %q is requilred. skip", k)
			return
		}
	}

	siteName := "your-site-name"
	// domainName := os.Getenv(envWebAccelDomainName)
	origin := os.Getenv(envWebAccelOrigin)
	bucketName := os.Getenv(envObjectStorageBucketName)
	accessKey := os.Getenv(envObjectStorageAccessKeyId)
	secretKey := os.Getenv(envObjectStorageSecretAccessKey)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudWebAccelWebOriginLoggingConfig(siteName, origin, bucketName, accessKey, secretKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "name", siteName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.type", "web"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.origin", origin),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "logging.0.s3_bucket_name", bucketName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "vary_support", "true"),
				),
			},
		},
	})
}

func TestAccSakuraCloudResourceWebAccel_WebOriginWithOriginGuard(t *testing.T) {
	skipIfFakeModeEnabled(t)

	envKeys := []string{
		envWebAccelOrigin,
	}
	for _, k := range envKeys {
		if os.Getenv(k) == "" {
			t.Skipf("ENV %q is requilred. skip", k)
			return
		}
	}

	siteName := "your-site-name"
	origin := os.Getenv(envWebAccelOrigin)
	regexpNotEmpty := regexp.MustCompile(".+")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudWebAccelWebOriginConfigWithOriginGuard(siteName, origin, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "name", siteName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.type", "web"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.origin", origin),
					resource.TestMatchResourceAttr("sakuracloud_webaccel.foobar", "origin_guard_token.0.token", regexpNotEmpty),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_guard_token.0.rotate", "false"),
				),
			},
			{
				Config: testAccCheckSakuraCloudWebAccelWebOriginConfigWithOriginGuard(siteName, origin, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "name", siteName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.type", "web"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.origin", origin),
					resource.TestMatchResourceAttr("sakuracloud_webaccel.foobar", "origin_guard_token.0.token", regexpNotEmpty),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_guard_token.0.rotate", "true"),
				),
			},
			{
				Config: testAccCheckSakuraCloudWebAccelWebOriginConfigWithOriginGuard(siteName, origin, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "name", siteName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.type", "web"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.origin", origin),
					resource.TestMatchResourceAttr("sakuracloud_webaccel.foobar", "origin_guard_token.0.token", regexpNotEmpty),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_guard_token.0.rotate", "false"),
					resource.TestCheckResourceAttrSet("sakuracloud_webaccel.foobar", "origin_guard_token.0.token"),
				),
			},
			{
				Config:      testAccCheckSakuraCloudWebAccelWebOriginConfigBasic(siteName, origin),
				ExpectError: regexpNotEmpty,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "name", siteName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.type", "web"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.origin", origin),
					// Note: It fails the absence of a token, but it does not check the removal of the token in real.
					// FIXME: So you require the manual check of the token absence to ensure the valid behavior.
					resource.TestCheckResourceAttrSet("sakuracloud_webaccel.foobar", "origin_guard_token.0.token"),
				),
			},
		},
	})
}

func TestAccSakuraCloudResourceWebAccel_InvalidConfigurations(t *testing.T) {
	if os.Getenv(envWebAccelOrigin) == "" {
		t.Skipf("ENV %q is requilred. skip", envWebAccelOrigin)
		return
	}
	origin := os.Getenv(envWebAccelOrigin)
	for name, tc := range testAccCheckSakuraCloudWebAccelInvalidConfigs(origin) {
		t.Logf("test for invalid configuration: %s", name)
		resource.Test(t, resource.TestCase{
			ProviderFactories: testAccProviderFactories,
			CheckDestroy: func(*terraform.State) error {
				return nil
			},
			Steps: []resource.TestStep{
				{
					Config: tc,
					ExpectError: func() *regexp.Regexp {
						if name == "valid" {
							return nil
						}
						return regexp.MustCompile(".")
					}(),
				},
			},
		})
	}
}

func testAccCheckSakuraCloudWebAccelWebOriginConfigBasic(siteName string, origin string) string {
	tmpl := `
resource sakuracloud_webaccel "foobar" {
  name = "%s"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "web"
    origin = "%s"
    host_header = "%s"
    protocol = "https"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "br+gzip"
}
`
	return fmt.Sprintf(tmpl, siteName, origin, origin)
}

func testAccCheckSakuraCloudWebAccelWebOriginConfigWithOneTimeUrlSecrets(siteName string, origin string) string {
	tmpl := `
resource sakuracloud_webaccel "foobar" {
  name = "%s"
  domain_type = "subdomain"
  request_protocol = "http+https"
  origin_parameters {
    type = "web"
    origin = "%s"
    host_header = "%s"
    protocol = "https"
  }
  onetime_url_secrets = [
    "sample-secret"
  ]
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "gzip"
}
`
	return fmt.Sprintf(tmpl, siteName, origin, origin)
}

func testAccCheckSakuraCloudWebAccelWebOriginConfigWithCors(siteName string, origin string) string {
	tmpl := `
resource sakuracloud_webaccel "foobar" {
  name = "%s"
  domain_type = "subdomain"
  request_protocol = "http+https"
  origin_parameters {
    type = "web"
    origin = "%s"
    host_header = "%s"
    protocol = "https"
  }
  cors_rules {
    allow_all = false
    allowed_origins = [
       "https://apps.example.com",
       "https://platform.example.com"
    ]
  }
  onetime_url_secrets = [
    "sample-secret"
  ]
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "gzip"
}
`
	return fmt.Sprintf(tmpl, siteName, origin, origin)
}

func testAccCheckSakuraCloudWebAccelWebOriginConfigWithOriginGuard(siteName string, origin string, hasRotateField bool) string {
	tmpl := `
resource sakuracloud_webaccel "foobar" {
  name = "%s"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "web"
    origin = "%s"
    host_header = "%s"
    protocol = "https"
  }
  origin_guard_token {%s}
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "br+gzip"
}
`
	if hasRotateField {
		return fmt.Sprintf(tmpl, siteName, origin, origin, "\n    rotate = true\n")
	} else {
		return fmt.Sprintf(tmpl, siteName, origin, origin, "")
	}
}

func testAccCheckSakuraCloudWebAccelBucketOriginConfig(siteName string, s3Endpoint string, region string, bucketName string, accessKey string, accessSecret string) string {
	tmpl := `
resource sakuracloud_webaccel "foobar" {
  name = "%s"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "bucket"
    s3_endpoint = "%s"
    s3_region = "%s"
    s3_bucket_name = "%s"
    s3_access_key_id = "%s"
    s3_secret_access_key = "%s"
    s3_doc_index = true
  }
  default_cache_ttl = 3600
  normalize_ae = "br+gzip"
}
`
	return fmt.Sprintf(tmpl, siteName, s3Endpoint, region, bucketName, accessKey, accessSecret)
}

func testAccCheckSakuraCloudWebAccelWebOriginLoggingConfig(siteName string, origin string, bucketName string, accessKey string, secretKey string) string {
	tmpl := `
resource sakuracloud_webaccel "foobar" {
  name = "%s"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "web"
    origin = "%s"
    host_header = "%s"
    protocol = "https"
  }
  logging {
    enabled = true
    s3_bucket_name = "%s"
    s3_access_key_id = "%s"
    s3_secret_access_key = "%s"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "br+gzip"
}
`
	return fmt.Sprintf(tmpl, siteName, origin, origin, bucketName, accessKey, secretKey)
}

func testAccCheckSakuraCloudWebAccelInvalidConfigs(origin string) map[string]string {
	confUnknownArgument := `
resource sakuracloud_webaccel "foobar" {
  invalid = true
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "web"
    origin = "%s"
    host_header = "dummy.example.com"
    protocol = "https"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "br+gzip"
}
`

	confInvalidDomainType := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "INVALID"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "web"
    origin = "%s"
    host_header = "dummy.example.com"
    protocol = "https"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "br+gzip"
}
`

	confInvalidRequestProtocol := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "http"
  origin_parameters {
    type = "web"
    origin = "%s"
    host_header = "dummy.example.com"
    protocol = "https"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "br+gzip"
}
`
	confWithoutOriginParameters := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "br+gzip"
}
`

	confInvalidOriginType := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "INVALID"
    origin = "%s"
    host_header = "dummy.example.com"
    protocol = "https"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "br+gzip"
}
`

	confLackingWebOriginParameters := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "web"
    host_header = "dummy.example.com"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "br+gzip"
}
`

	confMismatchedOriginParameters := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "bucket"
    host_header = "dummy.example.com"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "br+gzip"
}
`

	// config without the S3 s3_endpoint parameter
	confLackingBucketOriginParameters := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "bucket"
    s3_region = "jp-sample-1"
    s3_access_key_id = "sample"
    s3_secret_access_key = "sample"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "br+gzip"
}
`

	confInvalidNormalizeAE := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "web"
    origin = "%s"
    host_header = "dummy.example.com"
    protocol = "https"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "INVALID"
}
`
	// config without the S3 secret access key for logging
	confMissingLoggingParameters := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "web"
    origin = "docs.usacloud.jp"
    protocol = "https"
  }
  logging {
    s3_bucket_name = "example-bucket"
    s3_access_key_id = "sample"
  }
}
`
	// allow_all and allowed_origins should not be specified together
	confInvalidCorsConfiguration := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "web"
    origin = "%s"
    host_header = "dummy.example.com"
    protocol = "https"
  }
  cors_rules {
    allow_all = true
    allowed_origins = [
      "https://www2.example.com",
      "https://app.example.com"
    ]
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "br+gzip"
}
`

	valid := `
	resource sakuracloud_webaccel "foobar" {
	name = "dummy"
	domain_type = "subdomain"
	request_protocol = "https-redirect"
	origin_parameters {
	  type = "web"
	  origin = "%s"
	  protocol = "https"
	}
	cors_rules {
	  allowed_origins = [
	    "https://www2.example.com",
	    "https://app.example.com"
	  ]
	}
	vary_support = true
	default_cache_ttl = 3600
	normalize_ae = "br+gzip"
	}
	`

	tt := map[string]string{
		"unknown-argument":                         confUnknownArgument,
		"invalid-request-protocol":                 confInvalidRequestProtocol,
		"invalid-domain-type":                      confInvalidDomainType,
		"no-origin-params":                         confWithoutOriginParameters,
		"invalid-origin-type":                      confInvalidOriginType,
		"lacking-web-origin-params":                confLackingWebOriginParameters,
		"mismatched-bucket-origin-type-and-params": confMismatchedOriginParameters,
		"lacking-bucket-origin-params":             confLackingBucketOriginParameters,
		"invalid-compression":                      confInvalidNormalizeAE,
		"missing-logging-bucket-secret":            confMissingLoggingParameters,
		"invalid-cors-configuration":               confInvalidCorsConfiguration,
		"valid":                                    valid,
	}
	for k, v := range tt {
		if strings.Contains(v, "%s") {
			tt[k] = fmt.Sprintf(tt[k], origin)
		}
	}

	return tt
}
