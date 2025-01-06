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
)

type dummyResourceValueChangeHandler struct {
	oldState resourceValueGettable
	newState resourceValueGettable
}

func (d *dummyResourceValueChangeHandler) HasChanges(keys ...string) bool {
	for _, key := range keys {
		old := d.oldState.Get(key)
		new := d.newState.Get(key)
		if !reflect.DeepEqual(old, new) {
			return true
		}
	}
	return false
}

func (d *dummyResourceValueChangeHandler) GetChange(key string) (interface{}, interface{}) {
	return d.oldState.Get(key), d.newState.Get(key)
}

func TestStructureServer_isDiskEditParameterChanged(t *testing.T) {
	cases := []struct {
		msg    string
		in     *dummyResourceValueChangeHandler
		expect bool
	}{
		{
			msg: "nil",
			in: &dummyResourceValueChangeHandler{
				oldState: &resourceMapValue{},
				newState: &resourceMapValue{},
			},
			expect: false,
		},
		{
			msg: "added: disks",
			in: &dummyResourceValueChangeHandler{
				oldState: &resourceMapValue{},
				newState: &resourceMapValue{
					value: map[string]interface{}{
						"disks": []interface{}{"1"},
					},
				},
			},
			expect: true,
		},
		{
			msg: "added: disk_edit_parameter",
			in: &dummyResourceValueChangeHandler{
				oldState: &resourceMapValue{},
				newState: &resourceMapValue{
					value: map[string]interface{}{
						"disk_edit_parameter": []interface{}{
							map[string]interface{}{},
						},
					},
				},
			},
			expect: true,
		},
		{
			msg: "added: network_interface",
			in: &dummyResourceValueChangeHandler{
				oldState: &resourceMapValue{},
				newState: &resourceMapValue{
					value: map[string]interface{}{
						"network_interface": []interface{}{
							map[string]interface{}{},
						},
					},
				},
			},
			expect: true,
		},
		{
			msg: "updated: no changes",
			in: &dummyResourceValueChangeHandler{
				oldState: &resourceMapValue{
					value: map[string]interface{}{
						"disks": []interface{}{"1"},
						"disk_edit_parameter": []interface{}{
							map[string]interface{}{
								"password": "password",
							},
						},
						"network_interface": []interface{}{
							map[string]interface{}{
								"upstream": "shared",
							},
						},
					},
				},
				newState: &resourceMapValue{
					value: map[string]interface{}{
						"disks": []interface{}{"1"},
						"disk_edit_parameter": []interface{}{
							map[string]interface{}{
								"password": "password",
							},
						},
						"network_interface": []interface{}{
							map[string]interface{}{
								"upstream": "shared",
							},
						},
					},
				},
			},
			expect: false,
		},
		{
			msg: "updated: disks",
			in: &dummyResourceValueChangeHandler{
				oldState: &resourceMapValue{
					value: map[string]interface{}{
						"disks": []interface{}{"1"},
						"disk_edit_parameter": []interface{}{
							map[string]interface{}{
								"password": "password",
							},
						},
						"network_interface": []interface{}{
							map[string]interface{}{
								"upstream": "shared",
							},
						},
					},
				},
				newState: &resourceMapValue{
					value: map[string]interface{}{
						"disks": []interface{}{"2"},
						"disk_edit_parameter": []interface{}{
							map[string]interface{}{
								"password": "password",
							},
						},
						"network_interface": []interface{}{
							map[string]interface{}{
								"upstream": "shared",
							},
						},
					},
				},
			},
			expect: true,
		},
		{
			msg: "updated: disk_edit_parameter",
			in: &dummyResourceValueChangeHandler{
				oldState: &resourceMapValue{
					value: map[string]interface{}{
						"disks": []interface{}{"1"},
						"disk_edit_parameter": []interface{}{
							map[string]interface{}{
								"password": "password",
							},
						},
						"network_interface": []interface{}{
							map[string]interface{}{
								"upstream": "shared",
							},
						},
					},
				},
				newState: &resourceMapValue{
					value: map[string]interface{}{
						"disks": []interface{}{"1"},
						"disk_edit_parameter": []interface{}{
							map[string]interface{}{
								"password": "password-upd",
							},
						},
						"network_interface": []interface{}{
							map[string]interface{}{
								"upstream": "shared",
							},
						},
					},
				},
			},
			expect: true,
		},
		{
			msg: "updated: network_interface.upstream",
			in: &dummyResourceValueChangeHandler{
				oldState: &resourceMapValue{
					value: map[string]interface{}{
						"disks": []interface{}{"1"},
						"disk_edit_parameter": []interface{}{
							map[string]interface{}{
								"password": "password",
							},
						},
						"network_interface": []interface{}{
							map[string]interface{}{
								"upstream": "1",
							},
						},
					},
				},
				newState: &resourceMapValue{
					value: map[string]interface{}{
						"disks": []interface{}{"1"},
						"disk_edit_parameter": []interface{}{
							map[string]interface{}{
								"password": "password",
							},
						},
						"network_interface": []interface{}{
							map[string]interface{}{
								"upstream": "2",
							},
						},
					},
				},
			},
			expect: true,
		},
		{
			msg: "updated: network_interface.other",
			in: &dummyResourceValueChangeHandler{
				oldState: &resourceMapValue{
					value: map[string]interface{}{
						"disks": []interface{}{"1"},
						"disk_edit_parameter": []interface{}{
							map[string]interface{}{
								"password": "password",
							},
						},
						"network_interface": []interface{}{
							map[string]interface{}{
								"upstream":         "1",
								"user_ip_address":  "192.168.0.1",
								"packet_filter_id": "1",
							},
						},
					},
				},
				newState: &resourceMapValue{
					value: map[string]interface{}{
						"disks": []interface{}{"1"},
						"disk_edit_parameter": []interface{}{
							map[string]interface{}{
								"password": "password",
							},
						},
						"network_interface": []interface{}{
							map[string]interface{}{
								"upstream":         "1",
								"user_ip_address":  "192.168.0.2",
								"packet_filter_id": "2",
							},
						},
					},
				},
			},
			expect: false,
		},
		{
			msg: "deleted: disks",
			in: &dummyResourceValueChangeHandler{
				oldState: &resourceMapValue{
					value: map[string]interface{}{
						"disks": []interface{}{"1"},
						"disk_edit_parameter": []interface{}{
							map[string]interface{}{},
						},
						"network_interface": []interface{}{
							map[string]interface{}{
								"upstream": "shared",
							},
						},
					},
				},
				newState: &resourceMapValue{
					value: map[string]interface{}{
						"disk_edit_parameter": []interface{}{
							map[string]interface{}{},
						},
						"network_interface": []interface{}{
							map[string]interface{}{
								"upstream": "shared",
							},
						},
					},
				},
			},
			expect: true,
		},
		{
			msg: "deleted: network_interface",
			in: &dummyResourceValueChangeHandler{
				oldState: &resourceMapValue{
					value: map[string]interface{}{
						"disks": []interface{}{"1"},
						"disk_edit_parameter": []interface{}{
							map[string]interface{}{},
						},
						"network_interface": []interface{}{
							map[string]interface{}{
								"upstream": "shared",
							},
						},
					},
				},
				newState: &resourceMapValue{
					value: map[string]interface{}{
						"disks": []interface{}{"1"},
						"disk_edit_parameter": []interface{}{
							map[string]interface{}{},
						},
					},
				},
			},
			expect: true,
		},
	}

	for _, tc := range cases {
		got := isDiskEditParameterChanged(tc.in)
		if got != tc.expect {
			t.Fatalf("got unexpected state: pattern: %s expected: %t actual: %t", tc.msg, tc.expect, got)
		}
	}
}
