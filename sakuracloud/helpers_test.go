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
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func skipIfFakeModeEnabled(t *testing.T) {
	if isFakeModeEnabled() {
		t.Skip("This test only run if FAKE_MODE environment variable is not set")
	}
}

func isFakeModeEnabled() bool {
	fakeMode := os.Getenv("FAKE_MODE")
	return fakeMode != ""
}

func skipIfEnvIsNotSet(t *testing.T, key ...string) {
	for _, k := range key {
		if os.Getenv(k) == "" {
			t.Skipf("Environment valiable %q is not set", k)
		}
	}
}

func skipIfZoneIsDummy(t *testing.T) {
	if zone := os.Getenv("SAKURACLOUD_ZONE"); zone == "tk1v" {
		t.Skip("This test runs only on non-dummy zone")
	}
}

func testCheckSakuraCloudDataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource is not exists: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("id is not set: %s", n)
		}
		return nil
	}
}

func testCheckSakuraCloudDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("resource still exists: %s", n)
		}
		return nil
	}
}

func randomName() string {
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	return fmt.Sprintf("terraform-acctest-%s", rand)
}

func randomPassword() string {
	return acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
}

func buildConfigWithArgs(config string, args ...string) string {
	data := make(map[string]string)
	for i, v := range args {
		key := fmt.Sprintf("arg%d", i)
		data[key] = v
	}

	buf := bytes.NewBufferString("")
	err := template.Must(template.New("tmpl").Parse(config)).Execute(buf, data)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}
