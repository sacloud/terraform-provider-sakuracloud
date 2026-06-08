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
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
)

func TestFlattenSimpleMonitorHealthCheck(t *testing.T) {
	tt := []struct {
		Name           string
		InputMonitor   *iaas.SimpleMonitor
		InputConfig    resourceValueGettable
		ExpectPassword string
		ExpectContains string
		ExpectPath     string
	}{
		{
			Name: "restores password from config for http protocol",
			InputMonitor: &iaas.SimpleMonitor{
				HealthCheck: &iaas.SimpleMonitorHealthCheck{
					Protocol:          types.SimpleMonitorProtocols.HTTP,
					Path:              "/",
					Status:            types.StringNumber(200),
					ContainsString:    "ok",
					Host:              "usacloud.jp",
					Port:              types.StringNumber(80),
					BasicAuthUsername: "foo",
					BasicAuthPassword: "api-password",
				},
			},
			InputConfig: &resourceMapValue{
				value: map[string]interface{}{
					"health_check": []interface{}{
						map[string]interface{}{
							"password": "config-password",
						},
					},
				},
			},
			ExpectPassword: "config-password",
			ExpectContains: "ok",
			ExpectPath:     "/",
		},
		{
			Name: "restores password from config for https protocol",
			InputMonitor: &iaas.SimpleMonitor{
				HealthCheck: &iaas.SimpleMonitorHealthCheck{
					Protocol:          types.SimpleMonitorProtocols.HTTPS,
					Path:              "/",
					Status:            types.StringNumber(200),
					ContainsString:    "ok",
					Host:              "usacloud.jp",
					Port:              types.StringNumber(443),
					SNI:               types.StringFlag(true),
					BasicAuthUsername: "foo",
					BasicAuthPassword: "api-password",
					HTTP2:             types.StringFlag(true),
				},
			},
			InputConfig: &resourceMapValue{
				value: map[string]interface{}{
					"health_check": []interface{}{
						map[string]interface{}{
							"password": "config-password",
						},
					},
				},
			},
			ExpectPassword: "config-password",
			ExpectContains: "ok",
			ExpectPath:     "/",
		},
		{
			Name: "empty password when not in config",
			InputMonitor: &iaas.SimpleMonitor{
				HealthCheck: &iaas.SimpleMonitorHealthCheck{
					Protocol:          types.SimpleMonitorProtocols.HTTP,
					Path:              "/",
					Status:            types.StringNumber(200),
					ContainsString:    "",
					Host:              "",
					Port:              types.StringNumber(0),
					BasicAuthUsername: "foo",
					BasicAuthPassword: "api-password",
				},
			},
			InputConfig: &resourceMapValue{
				value: map[string]interface{}{},
			},
			ExpectPassword: "",
			ExpectPath:     "/",
		},
		{
			Name: "tcp protocol does not include password",
			InputMonitor: &iaas.SimpleMonitor{
				HealthCheck: &iaas.SimpleMonitorHealthCheck{
					Protocol: types.SimpleMonitorProtocols.TCP,
					Port:     types.StringNumber(22),
				},
			},
			InputConfig: &resourceMapValue{
				value: map[string]interface{}{},
			},
			ExpectPassword: "NO_PASSWORD_KEY",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			res := flattenSimpleMonitorHealthCheck(tc.InputMonitor, tc.InputConfig)
			if len(res) != 1 {
				t.Fatalf("expected 1 element got: %d", len(res))
			}
			got := res[0].(map[string]interface{})

			// Verify password restoration behavior
			if tc.ExpectPassword == "NO_PASSWORD_KEY" {
				if _, ok := got["password"]; ok {
					t.Fatalf("password key should not exist for %s protocol", tc.InputMonitor.HealthCheck.Protocol)
				}
			} else {
				if got["password"] != tc.ExpectPassword {
					t.Fatalf("password mismatch: got: %v want: %v", got["password"], tc.ExpectPassword)
				}
			}

			// Verify other fields are correctly set
			if tc.ExpectPath != "" {
				if got["path"] != tc.ExpectPath {
					t.Fatalf("path mismatch: got: %v want: %v", got["path"], tc.ExpectPath)
				}
			}
			if got["protocol"] != tc.InputMonitor.HealthCheck.Protocol {
				t.Fatalf("protocol mismatch: got: %v want: %v", got["protocol"], tc.InputMonitor.HealthCheck.Protocol)
			}
		})
	}
}
