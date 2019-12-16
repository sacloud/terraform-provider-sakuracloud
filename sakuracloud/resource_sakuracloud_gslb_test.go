// Copyright 2016-2019 terraform-provider-sakuracloud authors
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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func TestAccSakuraCloudGSLB_basic(t *testing.T) {
	resourceName := "sakuracloud_gslb.foobar"
	rand := randomName()

	var gslb sacloud.GSLB
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudGSLBDestroy,
			testCheckSakuraCloudIconDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudGSLB_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudGSLBExists(resourceName, &gslb),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "sorry_server", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.delay_loop", "10"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.host_header", "usacloud.jp"),
					resource.TestCheckResourceAttr(resourceName, "server.0.ip_address", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "server.0.weight", "1"),
					resource.TestCheckResourceAttr(resourceName, "server.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "server.1.ip_address", "2.2.2.2"),
					resource.TestCheckResourceAttr(resourceName, "server.1.weight", "2"),
					resource.TestCheckResourceAttr(resourceName, "server.1.enabled", "false"),
					resource.TestCheckResourceAttrPair(
						resourceName, "icon_id",
						"sakuracloud_icon.foobar", "id",
					),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudGSLB_update, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudGSLBExists(resourceName, &gslb),
					resource.TestCheckResourceAttr(resourceName, "name", rand+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "description-upd"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1-upd"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2-upd"),
					resource.TestCheckResourceAttr(resourceName, "sorry_server", "8.8.4.4"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.protocol", "https"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.delay_loop", "20"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.host_header", "usacloud.jp-upd"),
					resource.TestCheckResourceAttr(resourceName, "server.0.ip_address", "2.2.2.2"),
					resource.TestCheckResourceAttr(resourceName, "server.0.weight", "2"),
					resource.TestCheckResourceAttr(resourceName, "server.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "server.1.ip_address", "3.3.3.3"),
					resource.TestCheckResourceAttr(resourceName, "server.1.weight", "3"),
					resource.TestCheckResourceAttr(resourceName, "server.1.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "icon_id", ""),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudGSLB_serverDeleted, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudGSLBExists(resourceName, &gslb),
					resource.TestCheckResourceAttr(resourceName, "name", rand+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "description-upd"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1-upd"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2-upd"),
					resource.TestCheckResourceAttr(resourceName, "sorry_server", "8.8.4.4"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.protocol", "https"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.delay_loop", "20"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.host_header", "usacloud.jp-upd"),
					resource.TestCheckResourceAttr(resourceName, "server.0.ip_address", "2.2.2.2"),
					resource.TestCheckResourceAttr(resourceName, "server.0.weight", "2"),
					resource.TestCheckResourceAttr(resourceName, "server.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "server.1.ip_address", "3.3.3.3"),
					resource.TestCheckResourceAttr(resourceName, "server.1.weight", "3"),
					resource.TestCheckResourceAttr(resourceName, "server.1.enabled", "false"),
				),
			},
		},
	})
}

func testCheckSakuraCloudGSLBExists(n string, gslb *sacloud.GSLB) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no GSLB ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		gslbOp := sacloud.NewGSLBOp(client)
		ctx := context.Background()

		foundGSLB, err := gslbOp.Read(ctx, sakuraCloudID(rs.Primary.ID))
		if err != nil {
			return err
		}

		if foundGSLB.ID.String() != rs.Primary.ID {
			return fmt.Errorf("not found GSLB: %s", rs.Primary.ID)
		}
		*gslb = *foundGSLB
		return nil
	}
}

func testCheckSakuraCloudGSLBDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	gslbOp := sacloud.NewGSLBOp(client)
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_gslb" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		_, err := gslbOp.Read(ctx, sakuraCloudID(rs.Primary.ID))
		if err == nil {
			return errors.New("GSLB still exists")
		}
	}

	return nil
}

func TestAccImportSakuraCloudGSLB(t *testing.T) {
	rand := randomName()
	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("expected 1 state: %#v", s)
		}
		expects := map[string]string{
			"name":                       rand,
			"health_check.0.protocol":    "http",
			"health_check.0.delay_loop":  "10",
			"health_check.0.host_header": "usacloud.jp",
			"health_check.0.path":        "/",
			"health_check.0.status":      "200",
			"weighted":                   "false",
			"sorry_server":               "8.8.8.8",
			"description":                "description",
			"tags.0":                     "tag1",
			"tags.1":                     "tag2",
		}

		if err := compareStateMulti(s[0], expects); err != nil {
			return err
		}
		return stateNotEmptyMulti(s[0], "icon_id", "fqdn")
	}

	resourceName := "sakuracloud_gslb.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudGSLBDestroy,
			testCheckSakuraCloudIconDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudGSLB_basic, rand),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateCheck:  checkFn,
				ImportStateVerify: true,
			},
		},
	})
}

var testAccSakuraCloudGSLB_basic = `
resource "sakuracloud_gslb" "foobar" {
  name = "{{ .arg0 }}"
  health_check {
    protocol    = "http"
    delay_loop  = 10
    host_header = "usacloud.jp"
    path        = "/"
    status      = "200"
  }
  sorry_server = "8.8.8.8"
  server {
    ip_address = "1.1.1.1"
    weight     = 1
    enabled    = true
  }
  server {
    ip_address = "2.2.2.2"
    weight     = 2
    enabled    = false
  }
  description = "description"
  tags        = ["tag1", "tag2"]
  icon_id     = sakuracloud_icon.foobar.id
}

resource "sakuracloud_icon" "foobar" {
  name = "{{ .arg0 }}"
  base64content = "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"
}
`

var testAccSakuraCloudGSLB_update = `
resource "sakuracloud_gslb" "foobar" {
  name = "{{ .arg0 }}-upd"
  health_check {
    protocol    = "https"
    delay_loop  = 20
    host_header = "usacloud.jp-upd"
    path        = "/"
    status      = "200"
  }
  server {
    ip_address = "2.2.2.2"
    weight     = 2
    enabled    = true
  }
  server {
    ip_address = "3.3.3.3"
    weight     = 3
    enabled    = false
  }
  sorry_server = "8.8.4.4"
  description  = "description-upd"
  tags         = ["tag1-upd", "tag2-upd"]
}`

var testAccSakuraCloudGSLB_serverDeleted = `
resource "sakuracloud_gslb" "foobar" {
  name = "{{ .arg0 }}-upd"
  health_check {
    protocol    = "https"
    delay_loop  = 20
    host_header = "usacloud.jp-upd"
    path        = "/"
    status      = "200"
  }
  sorry_server = "8.8.4.4"
  description  = "description-upd"
  tags         = ["tag1-upd", "tag2-upd"]
}`
