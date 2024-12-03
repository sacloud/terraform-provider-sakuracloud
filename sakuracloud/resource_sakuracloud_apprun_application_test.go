// Copyright 2016-2023 The sacloud/terraform-provider-sakuracloud Authors
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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/sacloud/apprun-api-go"
	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

func TestAccSakuraCloudApprunApplication_basic(t *testing.T) {
	resourceName := "sakuracloud_apprun_application.foobar"
	rand := randomName()

	var application v1.Application
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testCheckSakuraCloudApprunApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudApprunApplication_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudApprunApplicationExists(resourceName, &application),
					testCheckSakuraCloudApprunApplicationAttributes(&application),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "timeout_seconds", "90"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "min_scale", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_scale", "1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.name", "compo1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.max_cpu", "0.1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.max_memory", "256Mi"),
					resource.TestCheckResourceAttr(resourceName, "components.0.deploy_source.0.container_registry.0.image", "apprun-test.sakuracr.jp/test1:latest"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudApprunApplication_update, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudApprunApplicationExists(resourceName, &application),
					testCheckSakuraCloudApprunApplicationAttributes(&application),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "timeout_seconds", "90"),
					resource.TestCheckResourceAttr(resourceName, "port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "min_scale", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_scale", "2"),
					resource.TestCheckResourceAttr(resourceName, "components.0.name", "compo1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.max_cpu", "0.2"),
					resource.TestCheckResourceAttr(resourceName, "components.0.max_memory", "512Mi"),
					resource.TestCheckResourceAttr(resourceName, "components.0.deploy_source.0.container_registry.0.image", "apprun-test.sakuracr.jp/test1:tag1"),
				),
			},
		},
	})
}

func TestAccSakuraCloudApprunApplication_withCRUser(t *testing.T) {
	resourceName := "sakuracloud_apprun_application.foobar"
	rand := randomName()

	var application v1.Application
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testCheckSakuraCloudApprunApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudApprunApplication_withCRUser, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudApprunApplicationExists(resourceName, &application),
					testCheckSakuraCloudApprunApplicationAttributes(&application),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "timeout_seconds", "90"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "min_scale", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_scale", "1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.name", "compo1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.max_cpu", "0.1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.max_memory", "256Mi"),
					resource.TestCheckResourceAttr(resourceName, "components.0.deploy_source.0.container_registry.0.image", "apprun-test.sakuracr.jp/test1:latest"),
					resource.TestCheckResourceAttr(resourceName, "components.0.deploy_source.0.container_registry.0.server", "apprun-test.sakuracr.jp"),
					resource.TestCheckResourceAttr(resourceName, "components.0.deploy_source.0.container_registry.0.username", "user"),
					resource.TestCheckResourceAttr(resourceName, "components.0.deploy_source.0.container_registry.0.password", "password"),
				),
			},
		},
	})
}

func TestAccSakuraCloudApprunApplication_withEnv(t *testing.T) {
	resourceName := "sakuracloud_apprun_application.foobar"
	rand := randomName()

	var application v1.Application
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testCheckSakuraCloudApprunApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudApprunApplication_withEnv, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudApprunApplicationExists(resourceName, &application),
					testCheckSakuraCloudApprunApplicationAttributes(&application),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "timeout_seconds", "90"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "min_scale", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_scale", "1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.name", "compo1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.max_cpu", "0.1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.max_memory", "256Mi"),
					resource.TestCheckResourceAttr(resourceName, "components.0.deploy_source.0.container_registry.0.image", "apprun-test.sakuracr.jp/test1:latest"),
					resource.TestCheckResourceAttr(resourceName, "components.0.env.0.key", "key"),
					resource.TestCheckResourceAttr(resourceName, "components.0.env.0.value", "value"),
					resource.TestCheckResourceAttr(resourceName, "components.0.env.1.key", "key2"),
					resource.TestCheckResourceAttr(resourceName, "components.0.env.1.value", "value2"),
				),
			},
		},
	})
}

func TestAccSakuraCloudApprunApplication_withProbe(t *testing.T) {
	resourceName := "sakuracloud_apprun_application.foobar"
	rand := randomName()

	var application v1.Application
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testCheckSakuraCloudApprunApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudApprunApplication_withProbe, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudApprunApplicationExists(resourceName, &application),
					testCheckSakuraCloudApprunApplicationAttributes(&application),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "timeout_seconds", "90"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "min_scale", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_scale", "1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.name", "compo1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.max_cpu", "0.1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.max_memory", "256Mi"),
					resource.TestCheckResourceAttr(resourceName, "components.0.deploy_source.0.container_registry.0.image", "apprun-test.sakuracr.jp/test1:latest"),
					resource.TestCheckResourceAttr(resourceName, "components.0.probe.0.http_get.0.path", "/"),
					resource.TestCheckResourceAttr(resourceName, "components.0.probe.0.http_get.0.port", "80"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudApprunApplication_withProbeUpdate, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudApprunApplicationExists(resourceName, &application),
					testCheckSakuraCloudApprunApplicationAttributes(&application),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "timeout_seconds", "90"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "min_scale", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_scale", "1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.name", "compo1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.max_cpu", "0.1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.max_memory", "256Mi"),
					resource.TestCheckResourceAttr(resourceName, "components.0.deploy_source.0.container_registry.0.image", "apprun-test.sakuracr.jp/test1:latest"),
					resource.TestCheckResourceAttr(resourceName, "components.0.probe.0.http_get.0.path", "/"),
					resource.TestCheckResourceAttr(resourceName, "components.0.probe.0.http_get.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "components.0.probe.0.http_get.0.headers.0.name", "name1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.probe.0.http_get.0.headers.0.value", "value1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.probe.0.http_get.0.headers.1.name", "name2"),
					resource.TestCheckResourceAttr(resourceName, "components.0.probe.0.http_get.0.headers.1.value", "value2"),
				),
			},
		},
	})
}

func TestAccSakuraCloudApprunApplication_withTraffic(t *testing.T) {
	resourceName := "sakuracloud_apprun_application.foobar"
	rand := randomName()

	var application v1.Application
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testCheckSakuraCloudApprunApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudApprunApplication_withTraffic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudApprunApplicationExists(resourceName, &application),
					testCheckSakuraCloudApprunApplicationAttributes(&application),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "timeout_seconds", "90"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "min_scale", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_scale", "1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.name", "compo1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.max_cpu", "0.1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.max_memory", "256Mi"),
					resource.TestCheckResourceAttr(resourceName, "components.0.deploy_source.0.container_registry.0.image", "apprun-test.sakuracr.jp/test1:latest"),
					resource.TestCheckResourceAttr(resourceName, "traffics.0.version_index", "0"),
					resource.TestCheckResourceAttr(resourceName, "traffics.0.percent", "100"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudApprunApplication_withTrafficUpdate, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudApprunApplicationExists(resourceName, &application),
					testCheckSakuraCloudApprunApplicationAttributes(&application),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "timeout_seconds", "10"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "min_scale", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_scale", "1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.name", "compo1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.max_cpu", "0.1"),
					resource.TestCheckResourceAttr(resourceName, "components.0.max_memory", "256Mi"),
					resource.TestCheckResourceAttr(resourceName, "components.0.deploy_source.0.container_registry.0.image", "apprun-test.sakuracr.jp/test1:latest"),
					resource.TestCheckResourceAttr(resourceName, "traffics.0.version_index", "0"),
					resource.TestCheckResourceAttr(resourceName, "traffics.0.percent", "1"),
					resource.TestCheckResourceAttr(resourceName, "traffics.1.version_index", "1"),
					resource.TestCheckResourceAttr(resourceName, "traffics.1.percent", "99"),
				),
			},
		},
	})
}

func testCheckSakuraCloudApprunApplicationExists(n string, application *v1.Application) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no AppRun Application ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		appOp := apprun.NewApplicationOp(client.apprunClient)

		found, err := appOp.Read(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if *found.Id != rs.Primary.ID {
			return fmt.Errorf("not found AppRun Application: %s", rs.Primary.ID)
		}

		*application = *found
		return nil
	}
}

func testCheckSakuraCloudApprunApplicationAttributes(application *v1.Application) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(*application.Components) == 0 {
			return errors.New("unexpected application components: components is nil")
		}

		c := (*application.Components)[0]
		if c.DeploySource.ContainerRegistry == nil {
			return errors.New("unexpected application components: container_registry is nil")
		}

		return nil
	}
}

func testCheckSakuraCloudApprunApplicationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	appOp := apprun.NewApplicationOp(client.apprunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_apprun_application" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		_, err := appOp.Read(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("still exists AppRun Application:%s", rs.Primary.ID)
		}
	}

	return nil
}

const testAccSakuraCloudApprunApplication_basic = `
resource "sakuracloud_apprun_application" "foobar" {
  name            = "{{ .arg0 }}"
  timeout_seconds = 90
  port            = 80
  min_scale       = 0
  max_scale       = 1
  components {
    name       = "compo1"
    max_cpu    = "0.1"
    max_memory = "256Mi"
    deploy_source {
      container_registry {
        image    = "apprun-test.sakuracr.jp/test1:latest"
      }
    }
  }
}
`

const testAccSakuraCloudApprunApplication_update = `
resource "sakuracloud_apprun_application" "foobar" {
  name            = "{{ .arg0 }}"
  timeout_seconds = 90
  port            = 8080
  min_scale       = 0
  max_scale       = 2
  components {
    name       = "compo1"
    max_cpu    = "0.2"
    max_memory = "512Mi"
    deploy_source {
      container_registry {
        image    = "apprun-test.sakuracr.jp/test1:tag1"
      }
    }
  }
}
`

const testAccSakuraCloudApprunApplication_withCRUser = `
resource "sakuracloud_apprun_application" "foobar" {
  name            = "{{ .arg0 }}"
  timeout_seconds = 90
  port            = 80
  min_scale       = 0
  max_scale       = 1
  components {
    name       = "compo1"
    max_cpu    = "0.1"
    max_memory = "256Mi"
    deploy_source {
      container_registry {
        image    = "apprun-test.sakuracr.jp/test1:latest"
        server   = "apprun-test.sakuracr.jp"
        username = "user"
        password = "password"
      }
    }
  }
}
`

const testAccSakuraCloudApprunApplication_withEnv = `
resource "sakuracloud_apprun_application" "foobar" {
  name            = "{{ .arg0 }}"
  timeout_seconds = 90
  port            = 80
  min_scale       = 0
  max_scale       = 1
  components {
    name       = "compo1"
    max_cpu    = "0.1"
    max_memory = "256Mi"
    deploy_source {
      container_registry {
        image    = "apprun-test.sakuracr.jp/test1:latest"
      }
    }
    env {
      key   = "key"
      value = "value"
    }
    env {
      key   = "key2"
      value = "value2"
    }
  }
}
`

const testAccSakuraCloudApprunApplication_withProbe = `
resource "sakuracloud_apprun_application" "foobar" {
  name            = "{{ .arg0 }}"
  timeout_seconds = 90
  port            = 80
  min_scale       = 0
  max_scale       = 1
  components {
    name       = "compo1"
    max_cpu    = "0.1"
    max_memory = "256Mi"
    deploy_source {
      container_registry {
        image    = "apprun-test.sakuracr.jp/test1:latest"
      }
    }
    probe {
      http_get {
        path = "/"
        port = 80
      }
    }
  }
}
`

const testAccSakuraCloudApprunApplication_withProbeUpdate = `
resource "sakuracloud_apprun_application" "foobar" {
  name            = "{{ .arg0 }}"
  timeout_seconds = 90
  port            = 80
  min_scale       = 0
  max_scale       = 1
  components {
    name       = "compo1"
    max_cpu    = "0.1"
    max_memory = "256Mi"
    deploy_source {
      container_registry {
        image    = "apprun-test.sakuracr.jp/test1:latest"
      }
    }
    probe {
      http_get {
        path = "/"
        port = 80
        headers {
          name  = "name1"
          value = "value1"
        }
        headers {
          name  = "name2"
          value = "value2"
        }
      }
    }
  }
}
`

const testAccSakuraCloudApprunApplication_withTraffic = `
resource "sakuracloud_apprun_application" "foobar" {
  name            = "{{ .arg0 }}"
  timeout_seconds = 90
  port            = 80
  min_scale       = 0
  max_scale       = 1
  components {
    name       = "compo1"
    max_cpu    = "0.1"
    max_memory = "256Mi"
    deploy_source {
      container_registry {
        image    = "apprun-test.sakuracr.jp/test1:latest"
      }
    }
  }
  traffics {
    version_index = 0
    percent       = 100
  }
}
`

const testAccSakuraCloudApprunApplication_withTrafficUpdate = `
resource "sakuracloud_apprun_application" "foobar" {
  name            = "{{ .arg0 }}"
  timeout_seconds = 10
  port            = 80
  min_scale       = 0
  max_scale       = 1
  components {
    name       = "compo1"
    max_cpu    = "0.1"
    max_memory = "256Mi"
    deploy_source {
      container_registry {
        image    = "apprun-test.sakuracr.jp/test1:latest"
      }
    }
  }
  traffics {
    version_index = 0
    percent       = 1
  }
  traffics {
    version_index = 1
    percent       = 99
  }
}
`
