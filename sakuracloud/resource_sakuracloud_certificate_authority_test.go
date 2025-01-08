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
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/sacloud/iaas-api-go"
)

func TestAccSakuraCloudCertificateAuthority_basic(t *testing.T) {
	skipIfFakeModeEnabled(t)
	skipIfEnvIsNotSet(t, "SAKURACLOUD_ENABLE_MANAGED_PKI")

	resourceName := "sakuracloud_certificate_authority.foobar"
	rand := randomName()

	var reg iaas.CertificateAuthority
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudCertificateAuthorityDestroy,
			testCheckSakuraCloudIconDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudCertificateAuthority_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudCertificateAuthorityExists(resourceName, &reg),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "subject.0.common_name", "pki.usacloud.jp"),
					resource.TestCheckResourceAttr(resourceName, "subject.0.organization", "usacloud"),
					resource.TestCheckResourceAttr(resourceName, "subject.0.organization_units.0", "ou1"),
					resource.TestCheckResourceAttr(resourceName, "subject.0.organization_units.1", "ou2"),
					resource.TestCheckResourceAttr(resourceName, "subject.0.country", "JP"),
					resource.TestCheckResourceAttrSet(resourceName, "certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "serial_number"),
					resource.TestCheckResourceAttrSet(resourceName, "not_before"),
					resource.TestCheckResourceAttrSet(resourceName, "not_after"),
					resource.TestCheckResourceAttrSet(resourceName, "crl_url"),
					resource.TestCheckResourceAttrPair(
						resourceName, "icon_id",
						"sakuracloud_icon.foobar", "id",
					),

					// client
					resource.TestCheckResourceAttr(resourceName, "client.#", "4"),

					resource.TestCheckResourceAttr(resourceName, "client.0.subject.0.common_name", "client1.usacloud.jp"),
					resource.TestCheckResourceAttr(resourceName, "client.0.url", ""),
					resource.TestCheckResourceAttrSet(resourceName, "client.0.issue_state"),
					resource.TestCheckResourceAttrSet(resourceName, "client.0.certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "client.0.serial_number"),
					resource.TestCheckResourceAttrSet(resourceName, "client.0.not_before"),
					resource.TestCheckResourceAttrSet(resourceName, "client.0.not_after"),

					resource.TestCheckResourceAttr(resourceName, "client.1.subject.0.common_name", "client2.usacloud.jp"),
					resource.TestCheckResourceAttr(resourceName, "client.1.url", ""),
					resource.TestCheckResourceAttrSet(resourceName, "client.1.issue_state"),
					resource.TestCheckResourceAttrSet(resourceName, "client.1.certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "client.1.serial_number"),
					resource.TestCheckResourceAttrSet(resourceName, "client.1.not_before"),
					resource.TestCheckResourceAttrSet(resourceName, "client.1.not_after"),

					resource.TestCheckResourceAttr(resourceName, "client.2.subject.0.common_name", "client3.usacloud.jp"),
					resource.TestCheckResourceAttr(resourceName, "client.2.url", ""),
					resource.TestCheckResourceAttrSet(resourceName, "client.2.issue_state"),
					resource.TestCheckResourceAttr(resourceName, "client.2.certificate", ""),
					resource.TestCheckResourceAttr(resourceName, "client.2.serial_number", ""),
					resource.TestCheckResourceAttr(resourceName, "client.2.not_before", ""),
					resource.TestCheckResourceAttr(resourceName, "client.2.not_after", ""),

					resource.TestCheckResourceAttr(resourceName, "client.3.subject.0.common_name", "client4.usacloud.jp"),
					resource.TestMatchResourceAttr(resourceName, "client.3.url", regexp.MustCompile(".+")),
					resource.TestCheckResourceAttrSet(resourceName, "client.3.issue_state"),
					resource.TestCheckResourceAttr(resourceName, "client.3.certificate", ""),
					resource.TestCheckResourceAttr(resourceName, "client.3.serial_number", ""),
					resource.TestCheckResourceAttr(resourceName, "client.3.not_before", ""),
					resource.TestCheckResourceAttr(resourceName, "client.3.not_after", ""),

					// server
					resource.TestCheckResourceAttr(resourceName, "server.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "server.0.subject.0.common_name", "server1.usacloud.jp"),
					resource.TestCheckResourceAttrSet(resourceName, "server.0.certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "server.0.serial_number"),
					resource.TestCheckResourceAttrSet(resourceName, "server.0.not_before"),
					resource.TestCheckResourceAttrSet(resourceName, "server.0.not_after"),
					resource.TestCheckResourceAttrSet(resourceName, "server.0.issue_state"),

					resource.TestCheckResourceAttr(resourceName, "server.1.subject.0.common_name", "server2.usacloud.jp"),
					resource.TestCheckResourceAttrSet(resourceName, "server.1.certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "server.1.serial_number"),
					resource.TestCheckResourceAttrSet(resourceName, "server.1.not_before"),
					resource.TestCheckResourceAttrSet(resourceName, "server.1.not_after"),
					resource.TestCheckResourceAttrSet(resourceName, "server.1.issue_state"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudCertificateAuthority_update, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudCertificateAuthorityExists(resourceName, &reg),
					resource.TestCheckResourceAttr(resourceName, "name", rand+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "description-upd"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1-upd"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2-upd"),
					resource.TestCheckResourceAttr(resourceName, "icon_id", ""),
					resource.TestCheckResourceAttrSet(resourceName, "certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "serial_number"),
					resource.TestCheckResourceAttrSet(resourceName, "not_before"),
					resource.TestCheckResourceAttrSet(resourceName, "not_after"),
					resource.TestCheckResourceAttrSet(resourceName, "crl_url"),
				),
			},
		},
	})
}

func testCheckSakuraCloudCertificateAuthorityExists(n string, auto_backup *iaas.CertificateAuthority) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no CertificateAuthority ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		regOp := iaas.NewCertificateAuthorityOp(client)

		foundCertificateAuthority, err := regOp.Read(context.Background(), sakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundCertificateAuthority.ID.String() != rs.Primary.ID {
			return errors.New("resource not found")
		}

		*auto_backup = *foundCertificateAuthority

		return nil
	}
}

func testCheckSakuraCloudCertificateAuthorityDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_certificate_authority" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		regOp := iaas.NewCertificateAuthorityOp(client)
		_, err := regOp.Read(context.Background(), sakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("CertificateAuthority still exists")
		}
	}

	return nil
}

var testAccSakuraCloudCertificateAuthority_basic = `
locals {
  dummy_public_key = <<EOT
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3tp9UdXap3VDB31oZof6
URnZN1BiO1cwOSi5cRsss27aeDZbhI03DZwzkJx0V95M6IeumaocQMIPiAkZZr8Z
knzC0UgIwn4H8/9VuC7MZuZyPj5f/2NSXF8V8wYpwg4UK0CVLwWoW2Z9Msws8Ls8
NiSqzngh8thR1vlq4aO7CzJDbt6Sgusu7XxE8CRXfwJ9dNIy/IA8lwkUi+gBYBZb
5DjAfg/RQxuPzvQsxjX84TO3XSkU+++MC0aol0UdkfInCFqTN9p9Ql/xvQlBM7NJ
BKGQ7/WMbktv0UZCQ+TNTyKL619syRuPoSAFMOt9SgnJdsmbjwhitkOdgYj6SJrf
uQIDAQAB
-----END PUBLIC KEY-----
EOT

  dummy_csr = <<EOT
-----BEGIN CERTIFICATE REQUEST-----
MIICgzCCAWsCAQAwPjELMAkGA1UEBhMCSlAxETAPBgNVBAgMCHVzYWNsb3VkMRww
GgYDVQQDDBNwa2ktY3NyLnVzYWNsb3VkLmpwMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEA3tp9UdXap3VDB31oZof6URnZN1BiO1cwOSi5cRsss27aeDZb
hI03DZwzkJx0V95M6IeumaocQMIPiAkZZr8ZknzC0UgIwn4H8/9VuC7MZuZyPj5f
/2NSXF8V8wYpwg4UK0CVLwWoW2Z9Msws8Ls8NiSqzngh8thR1vlq4aO7CzJDbt6S
gusu7XxE8CRXfwJ9dNIy/IA8lwkUi+gBYBZb5DjAfg/RQxuPzvQsxjX84TO3XSkU
+++MC0aol0UdkfInCFqTN9p9Ql/xvQlBM7NJBKGQ7/WMbktv0UZCQ+TNTyKL619s
yRuPoSAFMOt9SgnJdsmbjwhitkOdgYj6SJrfuQIDAQABoAAwDQYJKoZIhvcNAQEL
BQADggEBAFtYrKClAY0gsre2HbddbSek9kCZgK+NugW1irFqyJQ9aBXGTVQwtcI9
HBuA8vRoPEyzRl5Ua60mJK2YhAfG/uSJDgxWi0bK7Op574q9wdZMWc+hmolPX5kL
xEELoOsuwU5FB0azXCgnlmRJT5kbpIanCAKxScEDkJIB5qP/aSW1IjIlLgXh8CMr
vnreokhuEglFsL5CuMb72OlQVVc6E3DIheYLXhF83Pomff672shbm0HbDRWBgsMP
nryNWBxB/JyTkewcSPknZkeT9QSIV/AYwOmcC292T7EtF+fgSc01N5pgZigc/gOi
7S+hqAhb+LnU0WXc2PhwklN+xj+So1g=
-----END CERTIFICATE REQUEST-----
EOT
}

resource "sakuracloud_certificate_authority" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  tags        = ["tag1", "tag2"]
  icon_id     = sakuracloud_icon.foobar.id

  validity_period_hours = 24 * 3650
  subject {
    common_name        = "pki.usacloud.jp"
    country            = "JP"
    organization       = "usacloud"
    organization_units = ["ou1", "ou2"]
  }

  # by public_key
  client {
    subject {
      common_name        = "client1.usacloud.jp"
      country            = "JP"
      organization       = "usacloud"
      organization_units = ["ou1", "ou2"]
    }
    validity_period_hours = 24 * 3650
    public_key            = local.dummy_public_key
  }

  # by csr
  client {
    subject {
      common_name        = "client2.usacloud.jp"
      country            = "JP"
      organization       = "usacloud"
      organization_units = ["ou1", "ou2"]
    }
    validity_period_hours = 24 * 3650
    csr                   = local.dummy_csr
  }

  # by email
  client {
    subject {
      common_name        = "client3.usacloud.jp"
      country            = "JP"
      organization       = "usacloud"
      organization_units = ["ou1", "ou2"]
    }
    validity_period_hours = 24 * 3650
    email                 = "example@example.com"
  }

  # by URL
  client {
    subject {
      common_name        = "client4.usacloud.jp"
      country            = "JP"
      organization       = "usacloud"
      organization_units = ["ou1", "ou2"]
    }
    validity_period_hours = 24 * 3650
  }

  server {
    subject {
      common_name        = "server1.usacloud.jp"
      country            = "JP"
      organization       = "usacloud"
      organization_units = ["ou1", "ou2"]
    }

    subject_alternative_names = ["alt1.usacloud.jp", "alt2.usacloud.jp"]

    validity_period_hours = 24 * 3650
    public_key            = local.dummy_public_key
  }

  server {
    subject {
      common_name        = "server2.usacloud.jp"
      country            = "JP"
      organization       = "usacloud"
      organization_units = ["ou1", "ou2"]
    }

    subject_alternative_names = ["alt1.usacloud.jp", "alt2.usacloud.jp"]

    validity_period_hours = 24 * 3650
    csr                   = local.dummy_csr
  }
}

resource "sakuracloud_icon" "foobar" {
  name          = "{{ .arg0 }}"
  base64content = "iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAIAAADYYG7QAAAABGdBTUEAALGPC/xhBQAAAAFzUkdCAK7OHOkAAAAgY0hSTQAAeiYAAICEAAD6AAAAgOgAAHUwAADqYAAAOpgAABdwnLpRPAAAAAZiS0dEAP8A/wD/oL2nkwAAAAlwSFlzAAALEwAACxMBAJqcGAAACdBJREFUWMPNmHtw1NUVx8+5v9/+9rfJPpJNNslisgmIiCCgDQZR5GWnilUDPlpUqjOB2mp4qGM7tVOn/yCWh4AOVUprHRVB2+lMa0l88Kq10iYpNYPWkdeAmFjyEJPN7v5+v83ec/rH3Q1J2A2Z1hnYvz755ZzzvXPPveeee/GbC24FJmZGIYD5QgPpTBIAAICJLgJAwUQMAIDMfOEBUQchgJmAEC8CINLPThpfFCAG5orhogCBQiAAEyF8PQCATEQyxQzMzFIi4Ojdv86UEVF/f38ymezv7yciANR0zXAZhuHSdR0RRxNHZyJEBERmQvhfAAABIJlMJhIJt9t9TXX11GlTffleQGhvbz/4YeuRw4c13ZWfnycQR9ACQEShAyIxAxEKMXoAIVQ6VCzHcSzLmj937qqVK8aNrYKhv4bGxue3bvu8rc3n9+ualisyMzOltMjYccBqWanKdD5gBgAppZNMJhKJvlgs1heLxWL3fPfutU8/VVhYoGx7e3uJyOVyAcCEyy6bN2d266FDbW3thsuFI0gA4qy589PTOJC7EYEBbNu2ElYg4J9e/Y3p1dWBgN+l67csWKBC/mrbth07dnafOSMQp0y58pEVK2tm1ABAW9vn93zvgYRl5+XlAXMuCbxh3o3MDMyIguE8wADRaJ/H7Vp873119y8JBALDsrN8xcpXX3utoKDQNE1iiEV7ieSzmzYuXrwYAH7z4m83bNocDAZ1Tc8hQThrzjwYxY8BmCjaF/P78n+xZs0Ns64f+Ndnn53yevOLioo2btq8bsOGsvAYn9eHAoFZStnR0aFpWsObfxw/fvzp06fvXnyvZVmmx4M5hHQa3S4DwIRlm4Zr7dNPz7r+OgDo6el5bsuWtxrf6u7u9njygsHC9i/+U1Ia9ubnMzATA7MQIlRS8tnJk3/e1fDoI6vKysoqK8pbP/q323RDdi2hq/0ysHGyAwopU4lEfNXKlWo0Hx069MDSZcePHy8MBk3Tk0ylTnd1+wsKTNMERLUGlLtA1A3jyNEjagIKgsFk0gEM5NCSOst0+wEjAEvHtktKSuoeWAIAX3311f11Szs7OydcPtFwGYDp0sagWhoa7K4G5/f71TfHskEVdHXMn6M16CzLDcRkWfaM6dWm6QGAjZs2t7W1X1JeYRgGMzERMxOnNYa5O8mkrmkzr50JAKlUqq29Le2VQ0sACmYmIvU1OwAmLKt6ejUAyJTcu3dfQTCoaZqUkgEoY0ODvKRMSWbLsjo6O2fPmbuw9nYAOHjw4KdHjhqGoRqgLFpS6oNOE84JRDLVX1FeDgBd3V0pIrfLxZn5GGLMrE40y7YTCcula7W3167++c+UzfNbtzGRK+ObxR1RZyJARPUpNxBzPBYDAE3ThCYkETMjIPMQdwCwbNttGItqb6uqrJo2deqMGTVK8qWXX969+92SsjAi5hRF1BkQKJ3REUDXtE+PHL3ppptCoVBpcXFXVzdJqerFWWNmKaVt2T9YWldf//Dg6rL52efWrV/vCxQYLhdJmV2LmaUUkEkZZGbvXGBm0+P563vvqT/vW7LEcRwnmUxv7wFjZiYyDJdabQCQSsnt27d/6+YFT61Z4/UHBvZadi1mQBRERMwEMAIwkdttNh/8V2trKwB85647a2tv7+npTfb3y6HGKLREIvHKK6+my66ubd/x+p69+0KlZf5AQKV+BC0G0MaURwZGlxMAiam9vf3YsWNL7rsXAL694Oa2tvZPPvnEZRiozBABAIE1XfvggwMfffzxnXcsAoBrZ8zYs3+/pmm6ECNJIKrto4UvueQ8pxiRZduxWKympuauRQsnT56saRoAlIRCbzbsYmYhxGB7TdPcHk9LS3O4LHz1VVcFg8HmpubjJ0643W44/w8FS6kqW1YgKROW5VjWivr6P/3h93V1dYZhKNeD/2zp7elVjfAQLyKP2+0PFG5/NZ242XNm25bNRCNrKUjfy5gIzwXE/mQyEYs98dMnHnrw+yr6hx+2/qOp6djRo43vvGu4XJquZ3X3mO7OL8+cOnUqEolURSpUx53LeDDolDlE+ByQRNG+vlmzZ6vROI69fMWqN954Ix5PBAoLC4PBfK+XMqfSEHdEQJRS2ratyl1KSmLG3FoDoKcXFCIQDQOZTCLAQ8uWKtNlD/5w546dkaqqKq8XERDFQIkb7g6QSqUK/f5wOAwA0WgUiM+u/WxaChBRJxSgzsXhK5+sZDISiVxTUwMAjY2Nu3Y1RMZd6vXmAzCAIOB0uHP2SyqVisViCxcu9Pl8ANDc0oK6xswkxMg7mon0dGHMUqkg6Tjh0lLTdAPABwf+niKZ5zFRtRmQ8RrqyACyv783Gi0vL390eb0qqm+/szvPNNMzNGIFRnUvA0SAzOwNAiLJmU4zHo8DCgAgZgAETtswyX4pk8lkehP0pywrUTV27JaNGyqrKgHgha1bT548WRYOMwDk1hrIna46gbTAUBBCUwcqAFw6frwuRCqV0nUdmFB1MCRtx9E0bWwkEresRDzu9/nm3Th/Vf3DoVAIAJqbmtauXZfv9WpCpBd7Dq00EOGkKdNylCi0EgkhxP4971ZUVJw8ceK2RXd0dX9ZUFCgCaFyYTtOrC/22CMrf/LjH3V0dvX1RSsjEVemUDU3NS1d9uAXHR2lpaVqV4+iMIJWXFKKiEpgCCAKxI6OjuLioutmziwoLBxTFn7r7Xei0WhKSsdxYvF4PJ649Zabn1m/DhC93vxgMKiKuGUlntm46bHHHz/T0xsqKdEEZpYKZ9caJIpXTJmWfuVDofpPBcAMKKLRXoHwl727x106HgAOHDiw5ZcvHD5ymBiCwcJFtbXLM21GQ0ODZVm90ej77/9t3779XV2dBcEifyCgIcLQyCMBMU6cNCX3wQIkqbOzY+LlE373+s6KSER97untdSy7tKx0wHD16tVPPvkkAIDQvV6fz+fNz/emXzyAYVS5yqSsqLh4UM8GwwAFmqZ54sSJXY2NJSUlkyZNAgDTNL1er/Jvb29/uL7+1y++VFQcKg2PCYVCfr/XND1C01QnnytydkDECVdcqdpqtXGGgcqulHTmy+54PH71VdNunD+/sqoSEaPRaEtzy569exO2UxQM5nm9ynpQgrIEPA8w42UTJ6dLEkNWUI0KMTu2E4v3xftiSccGAKHpnrw8v8/vyfPoug4Zv1xxRgOIoDNJQAEMmfo9HNT9DxFN03QbRrCwCNQjHAp1gVc2mQKbM86oAFCA0GDQnSEXqMcGwPQjmND1zGgEAFBmNOeNMzIQSZ0GXvJHuJedPXRkLhiN+2hAVxUdz77yXWDQUdMGFUa40DC4Y/ya5vz/BMEkmVm9dl94QPwvNJB+oilXgHEAAAAldEVYdGRhdGU6Y3JlYXRlADIwMTYtMDItMTBUMjE6MDg6MzMtMDg6MDB4P0OtAAAAJXRFWHRkYXRlOm1vZGlmeQAyMDE2LTAyLTEwVDIxOjA4OjMzLTA4OjAwCWL7EQAAAABJRU5ErkJggg=="
}
`

var testAccSakuraCloudCertificateAuthority_update = `
resource "sakuracloud_certificate_authority" "foobar" {
  name        = "{{ .arg0 }}-upd"
  description = "description-upd"
  tags        = ["tag1-upd", "tag2-upd"]

  validity_period_hours = 24 * 3650
  subject {
    common_name        = "pki.usacloud.jp"
    country            = "JP"
    organization       = "usacloud"
    organization_units = ["ou1", "ou2"]
  }
}
`
