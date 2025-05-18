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
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSakuraCloudResourceWebAccel_WebOrigin(t *testing.T) {
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
	//domainName := os.Getenv(envWebAccelDomainName)
	origin := os.Getenv(envWebAccelOrigin)

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
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.host", origin),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "vary_support", "true"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "normalize_ae", "brotli"),
				),
			},
		},
	})
}

func TestAccSakuraCloudResourceWebAccel_WebOriginWithCORS(t *testing.T) {
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
	//domainName := os.Getenv(envWebAccelDomainName)
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
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.host", origin),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "cors_rules.0.allow_all", "false"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "cors_rules.0.allowed_origins.0", "https://apps.example.com"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "onetime_url_secrets.0", "sample-secret"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "normalize_ae", "gzip"),
				),
			},
		},
	})
}

func TestAccSakuraCloudResourceWebAccel_WebOriginUpdate(t *testing.T) {
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
	//domainName := os.Getenv(envWebAccelDomainName)
	origin := os.Getenv(envWebAccelOrigin)

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
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.host", origin),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "vary_support", "true"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "normalize_ae", "brotli"),
				),
			},
			{
				Config: testAccCheckSakuraCloudWebAccelWebOriginConfigWithCors(siteName, origin),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "name", siteName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.type", "web"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.host", origin),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "cors_rules.0.allow_all", "false"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "cors_rules.0.allowed_origins.0", "https://apps.example.com"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "normalize_ae", "gzip"),
				),
			},
		},
	})
}

func TestAccSakuraCloudResourceWebAccel_BucketOrigin(t *testing.T) {
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
	//domainName := os.Getenv(envWebAccelDomainName)
	endpoint, _ := strings.CutPrefix(os.Getenv(envObjectStorageEndpoint), "https://")
	log.Println("endpoint", endpoint)
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
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.endpoint", endpoint),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.region", region),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.bucket_name", bucketName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.access_key_id", accessKey),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.secret_access_key", secretKey),
				),
			},
		},
	})
}

func TestAccSakuraCloudResourceWebAccel_Logging(t *testing.T) {
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
	//domainName := os.Getenv(envWebAccelDomainName)
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
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.host", origin),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "logging.0.bucket_name", bucketName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "vary_support", "true"),
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
    host = "%s"
    host_header = "%s"
    protocol = "https"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "brotli"
}
`
	return fmt.Sprintf(tmpl, siteName, origin, origin)
}

func testAccCheckSakuraCloudWebAccelWebOriginConfigWithCors(siteName string, origin string) string {
	tmpl := `
resource sakuracloud_webaccel "foobar" {
  name = "%s"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "web"
    host = "%s"
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

func testAccCheckSakuraCloudWebAccelBucketOriginConfig(siteName string, endpoint string, region string, bucketName string, accessKey string, accessSecret string) string {
	tmpl := `
resource sakuracloud_webaccel "foobar" {
  name = "%s"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "bucket"
    endpoint = "%s"
    region = "%s"
    bucket_name = "%s"
    access_key_id = "%s"
    secret_access_key = "%s"
  }
  default_cache_ttl = 3600
  normalize_ae = "brotli"
}
`
	return fmt.Sprintf(tmpl, siteName, endpoint, region, bucketName, accessKey, accessSecret)
}

func testAccCheckSakuraCloudWebAccelWebOriginLoggingConfig(siteName string, origin string, bucketName string, accessKey string, secretKey string) string {
	tmpl := `
resource sakuracloud_webaccel "foobar" {
  name = "%s"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "web"
    host = "%s"
    host_header = "%s"
    protocol = "https"
  }
  logging {
    enabled = true
    bucket_name = "%s"
    access_key_id = "%s"
    secret_access_key = "%s"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "brotli"
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
    host = "%s"
    host_header = "dummy.example.com"
    protocol = "https"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "brotli"
}
`

	confInvalidDomainType := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "INVALID"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "web"
    host = "%s"
    host_header = "dummy.example.com"
    protocol = "https"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "brotli"
}
`

	confInvalidRequestProtocol := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "http"
  origin_parameters {
    type = "web"
    host = "%s"
    host_header = "dummy.example.com"
    protocol = "https"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "brotli"
}
`
	confWithoutOriginParameters := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "brotli"
}
`

	confInvalidOriginType := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "INVALID"
    host = "%s"
    host_header = "dummy.example.com"
    protocol = "https"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "brotli"
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
  normalize_ae = "brotli"
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
  normalize_ae = "brotli"
}
`

	//config without the S3 endpoint parameter
	confLackingBucketOriginParameters := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "bucket"
    region = "jp-sample-1"
    access_key_id = "sample"
    secret_access_key = "sample"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "brotli"
}
`

	confInvalidNormalizeAE := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "web"
    host = "%s"
    host_header = "dummy.example.com"
    protocol = "https"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "INVALID"
}
`

	//config without the S3 secret access key for logging
	confMissingLoggingParameters := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "web"
    host = "%s"
    host_header = "dummy.example.com"
    protocol = "https"
  }
  logging {
    bucket_name = "example-bucket"
    access_key_id = "sample"
  }
  vary_support = true
}
`

	//allow_all and allowed_origins should not be specified together
	confInvalidCorsConfiguration := `
resource sakuracloud_webaccel "foobar" {
  name = "dummy"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "web"
    host = "%s"
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
  normalize_ae = "brotli"
}
`

	//valid := `
	//resource sakuracloud_webaccel "foobar" {
	// name = "dummy"
	// domain_type = "subdomain"
	// request_protocol = "https-redirect"
	// origin_parameters {
	//   type = "web"
	//   host = "%s"
	//   host_header = "dummy.example.com"
	//   protocol = "https"
	// }
	// cors_rules {
	//   allowed_origins = [
	//     "https://www2.example.com",
	//     "https://app.example.com"
	//   ]
	// }
	// vary_support = true
	// default_cache_ttl = 3600
	// normalize_ae = "brotli"
	//}
	//`

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
		//"valid":                                    valid,
	}
	for k, v := range tt {
		if strings.Contains(v, "%s") {
			tt[k] = fmt.Sprintf(tt[k], origin)
		}
	}

	return tt
}
