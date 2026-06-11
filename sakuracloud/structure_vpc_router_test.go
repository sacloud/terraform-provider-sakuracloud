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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
)

func TestFlattenVPCRouterL2TP(t *testing.T) {
	tt := []struct {
		Name           string
		InputVPCRouter *iaas.VPCRouter
		InputConfig    resourceValueGettable
		ExpectedOutput []interface{}
	}{
		{
			Name: "restores pre_shared_secret from config",
			InputVPCRouter: &iaas.VPCRouter{
				Settings: &iaas.VPCRouterSetting{
					L2TPIPsecServerEnabled: types.StringTrue,
					L2TPIPsecServer: &iaas.VPCRouterL2TPIPsecServer{
						PreSharedSecret: "api-secret",
						RangeStart:      "192.168.11.21",
						RangeStop:       "192.168.11.30",
					},
				},
			},
			InputConfig: &resourceMapValue{
				value: map[string]interface{}{
					"l2tp": []interface{}{
						map[string]interface{}{
							"pre_shared_secret": "config-secret",
						},
					},
				},
			},
			ExpectedOutput: []interface{}{
				map[string]interface{}{
					"pre_shared_secret": "config-secret",
					"range_start":       "192.168.11.21",
					"range_stop":        "192.168.11.30",
				},
			},
		},
		{
			Name: "empty pre_shared_secret when not in config",
			InputVPCRouter: &iaas.VPCRouter{
				Settings: &iaas.VPCRouterSetting{
					L2TPIPsecServerEnabled: types.StringTrue,
					L2TPIPsecServer: &iaas.VPCRouterL2TPIPsecServer{
						PreSharedSecret: "api-secret",
						RangeStart:      "192.168.11.21",
						RangeStop:       "192.168.11.30",
					},
				},
			},
			InputConfig: &resourceMapValue{
				value: map[string]interface{}{},
			},
			ExpectedOutput: []interface{}{
				map[string]interface{}{
					"pre_shared_secret": "",
					"range_start":       "192.168.11.21",
					"range_stop":        "192.168.11.30",
				},
			},
		},
		{
			Name: "disabled L2TP returns empty",
			InputVPCRouter: &iaas.VPCRouter{
				Settings: &iaas.VPCRouterSetting{
					L2TPIPsecServerEnabled: types.StringFalse,
				},
			},
			InputConfig: &resourceMapValue{
				value: map[string]interface{}{},
			},
			ExpectedOutput: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			res := flattenVPCRouterL2TP(tc.InputVPCRouter, tc.InputConfig)
			if len(res) != len(tc.ExpectedOutput) {
				t.Fatalf("FAILED %s: length mismatch got: %d want: %d", tc.Name, len(res), len(tc.ExpectedOutput))
			}
			for i := range res {
				got := res[i].(map[string]interface{})
				want := tc.ExpectedOutput[i].(map[string]interface{})
				for k, v := range want {
					if got[k] != v {
						t.Fatalf("FAILED %s[%d].%s: got: %v want: %v", tc.Name, i, k, got[k], v)
					}
				}
			}
		})
	}
}

func TestFlattenVPCRouterSiteToSiteConfig(t *testing.T) {
	tt := []struct {
		Name           string
		InputVPCRouter *iaas.VPCRouter
		InputConfig    resourceValueGettable
		ExpectEnabled  bool
		ExpectPeer     string
		ExpectSecret   string
	}{
		{
			Name: "restores pre_shared_secret from config by peer",
			InputVPCRouter: &iaas.VPCRouter{
				Settings: &iaas.VPCRouterSetting{
					SiteToSiteIPsecVPN: &iaas.VPCRouterSiteToSiteIPsecVPN{
						Config: []*iaas.VPCRouterSiteToSiteIPsecVPNConfig{
							{
								Peer:            "8.8.8.8",
								PreSharedSecret: "api-secret",
								RemoteID:        "8.8.8.8",
								Routes:          []string{"10.0.0.0/8"},
								LocalPrefix:     []string{"192.168.21.0/24"},
							},
						},
					},
				},
			},
			InputConfig: &resourceMapValue{
				value: map[string]interface{}{
					"site_to_site_vpn": []interface{}{
						map[string]interface{}{
							"peer":              "8.8.8.8",
							"pre_shared_secret": "config-secret",
						},
					},
				},
			},
			ExpectEnabled: true,
			ExpectPeer:    "8.8.8.8",
			ExpectSecret:  "config-secret",
		},
		{
			Name: "empty pre_shared_secret when peer not in config",
			InputVPCRouter: &iaas.VPCRouter{
				Settings: &iaas.VPCRouterSetting{
					SiteToSiteIPsecVPN: &iaas.VPCRouterSiteToSiteIPsecVPN{
						Config: []*iaas.VPCRouterSiteToSiteIPsecVPNConfig{
							{
								Peer:            "8.8.8.8",
								PreSharedSecret: "api-secret",
								RemoteID:        "8.8.8.8",
								Routes:          []string{"10.0.0.0/8"},
								LocalPrefix:     []string{"192.168.21.0/24"},
							},
						},
					},
				},
			},
			InputConfig: &resourceMapValue{
				value: map[string]interface{}{},
			},
			ExpectEnabled: true,
			ExpectPeer:    "8.8.8.8",
			ExpectSecret:  "",
		},
		{
			Name: "no site_to_site_vpn settings returns empty",
			InputVPCRouter: &iaas.VPCRouter{
				Settings: &iaas.VPCRouterSetting{},
			},
			InputConfig: &resourceMapValue{
				value: map[string]interface{}{},
			},
			ExpectEnabled: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			res := flattenVPCRouterSiteToSiteConfig(tc.InputVPCRouter, tc.InputConfig)
			if !tc.ExpectEnabled {
				if len(res) != 0 {
					t.Fatalf("FAILED %s: expected empty got: %d elements", tc.Name, len(res))
				}
				return
			}
			if len(res) != 1 {
				t.Fatalf("FAILED %s: expected 1 element got: %d", tc.Name, len(res))
			}
			got := res[0].(map[string]interface{})
			if got["peer"] != tc.ExpectPeer {
				t.Fatalf("FAILED %s peer: got: %v want: %v", tc.Name, got["peer"], tc.ExpectPeer)
			}
			if got["pre_shared_secret"] != tc.ExpectSecret {
				t.Fatalf("FAILED %s pre_shared_secret: got: %v want: %v", tc.Name, got["pre_shared_secret"], tc.ExpectSecret)
			}
			if got["remote_id"] != "8.8.8.8" {
				t.Fatalf("FAILED %s remote_id: got: %v want: %v", tc.Name, got["remote_id"], "8.8.8.8")
			}
			for _, key := range []string{"routes", "local_prefix"} {
				if _, ok := got[key]; !ok {
					t.Fatalf("FAILED %s: missing key %s", tc.Name, key)
				}
				set, ok := got[key].(*schema.Set)
				if !ok {
					t.Fatalf("FAILED %s: %s is not *schema.Set", tc.Name, key)
				}
				if set.Len() != 1 {
					t.Fatalf("FAILED %s: %s set length mismatch", tc.Name, key)
				}
			}
		})
	}
}
