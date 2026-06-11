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
	"reflect"
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
)

//nolint:gosec // Test uses dummy credential values
func TestFlattenLocalRouterPeers(t *testing.T) {
	tt := []struct {
		Name           string
		InputData      *iaas.LocalRouter
		InputConfig    resourceValueGettable
		ExpectedOutput []interface{}
	}{
		{
			Name: "restores secret_key from config",
			InputData: &iaas.LocalRouter{
				Peers: []*iaas.LocalRouterPeer{
					{
						ID:          types.StringID("123456789012"),
						Enabled:     true,
						Description: "peer1",
						SecretKey:   "api-secret-key", // API may return this but should not be used
					},
				},
			},
			InputConfig: &resourceMapValue{
				value: map[string]interface{}{
					"peer": []interface{}{
						map[string]interface{}{
							"peer_id":    "123456789012",
							"secret_key": "config-secret-key",
						},
					},
				},
			},
			ExpectedOutput: []interface{}{
				map[string]interface{}{
					"peer_id":     "123456789012",
					"secret_key":  "config-secret-key",
					"enabled":     true,
					"description": "peer1",
				},
			},
		},
		{
			Name: "empty secret_key when peer not in config",
			InputData: &iaas.LocalRouter{
				Peers: []*iaas.LocalRouterPeer{
					{
						ID:          types.StringID("123456789012"),
						Enabled:     true,
						Description: "peer1",
						SecretKey:   "api-secret-key",
					},
				},
			},
			InputConfig: &resourceMapValue{
				value: map[string]interface{}{},
			},
			ExpectedOutput: []interface{}{
				map[string]interface{}{
					"peer_id":     "123456789012",
					"secret_key":  "",
					"enabled":     true,
					"description": "peer1",
				},
			},
		},
		{
			Name: "multiple peers with mixed config presence",
			InputData: &iaas.LocalRouter{
				Peers: []*iaas.LocalRouterPeer{
					{
						ID:          types.StringID("123456789012"),
						Enabled:     true,
						Description: "peer1",
						SecretKey:   "api-secret-1",
					},
					{
						ID:          types.StringID("999999999999"),
						Enabled:     false,
						Description: "peer2",
						SecretKey:   "api-secret-2",
					},
				},
			},
			InputConfig: &resourceMapValue{
				value: map[string]interface{}{
					"peer": []interface{}{
						map[string]interface{}{
							"peer_id":    "123456789012",
							"secret_key": "config-secret-1",
						},
					},
				},
			},
			ExpectedOutput: []interface{}{
				map[string]interface{}{
					"peer_id":     "123456789012",
					"secret_key":  "config-secret-1",
					"enabled":     true,
					"description": "peer1",
				},
				map[string]interface{}{
					"peer_id":     "999999999999",
					"secret_key":  "",
					"enabled":     false,
					"description": "peer2",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			res := flattenLocalRouterPeers(tc.InputData, tc.InputConfig)
			if !reflect.DeepEqual(res, tc.ExpectedOutput) {
				t.Fatalf("FAILED %s: got: %v\nwant: %v", tc.Name, res, tc.ExpectedOutput)
			}
		})
	}
}
