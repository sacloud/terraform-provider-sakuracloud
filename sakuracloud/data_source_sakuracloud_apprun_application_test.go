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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSakuraCloudDataSourceApprunApplication_basic(t *testing.T) {
	resourceName := "data.sakuracloud_apprun_application.foobar"
	rand := randomName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceApprunApplication_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
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
		},
	})
}

func TestAccSakuraCloudDataSourceApprunApplication_withCRUser(t *testing.T) {
	resourceName := "data.sakuracloud_apprun_application.foobar"
	rand := randomName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceApprunApplication_withCRUser, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
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
				),
			},
		},
	})
}

func TestAccSakuraCloudDataSourceApprunApplication_withProbe(t *testing.T) {
	resourceName := "sakuracloud_apprun_application.foobar"
	rand := randomName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceApprunApplication_withProbe, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
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

func TestAccSakuraCloudDataSourceApprunApplication_withTraffic(t *testing.T) {
	resourceName := "sakuracloud_apprun_application.foobar"
	rand := randomName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceApprunApplication_withTraffic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
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
		},
	})
}

var testAccSakuraCloudDataSourceApprunApplication_basic = `
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

data "sakuracloud_apprun_application" "foobar" {
  name = sakuracloud_apprun_application.foobar.name
}
`

var testAccSakuraCloudDataSourceApprunApplication_withCRUser = `
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

data "sakuracloud_apprun_application" "foobar" {
  name = sakuracloud_apprun_application.foobar.name
}
`

var testAccSakuraCloudDataSourceApprunApplication_withProbe = `
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

data "sakuracloud_apprun_application" "foobar" {
  name = sakuracloud_apprun_application.foobar.name
}
`

var testAccSakuraCloudDataSourceApprunApplication_withTraffic = `
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

data "sakuracloud_apprun_application" "foobar" {
  name = sakuracloud_apprun_application.foobar.name
}
`
