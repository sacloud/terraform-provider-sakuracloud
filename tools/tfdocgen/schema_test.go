// Copyright 2016-2021 terraform-provider-sakuracloud authors
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

package tfdocgen

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestNewSchema(t *testing.T) {
	in := map[string]*schema.Schema{
		"arg1": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"attr1": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"slice1": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"nest1_both": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"arg2": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"attr2": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"nest2": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"arg3": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"attr3": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		"nest1_arg": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"arg2": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"nest1_attr": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"attr2": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}

	expect := &Schema{
		Arguments: []Argument{
			{Name: "arg1", Optional: true},
			{Name: "nest1_arg", Optional: true, Description: "One or more `nest1_arg` blocks as defined below"},
			{Name: "nest1_attr", Optional: true, Description: "One or more `nest1_attr` blocks as defined below"},
			{Name: "nest1_both", Optional: true, Description: "A `nest1_both` block as defined below"},
			{Name: "slice1", Optional: true},
		},
		Attributes: []Attribute{
			{Name: "attr1"},
		},
		ArgumentBlocks: []ArgumentBlock{
			{
				Name:    "nest1_arg",
				Parents: []string{"nest1_arg"},
				Arguments: []Argument{
					{Name: "arg2", Optional: true},
				},
			},
			{
				Name:    "nest1_both",
				Parents: []string{"nest1_both"},
				Arguments: []Argument{
					{Name: "arg2", Optional: true},
					{Name: "nest2", Optional: true, Description: "One or more `nest2` blocks as defined below"},
				},
			},
			{
				Name:    "nest2",
				Parents: []string{"nest1_both", "nest2"},
				Arguments: []Argument{
					{Name: "arg3", Optional: true},
				},
			},
		},
		AttributeBlocks: []AttributeBlock{
			{
				Name:    "nest1_attr",
				Parents: []string{"nest1_attr"},
				Attributes: []Attribute{
					{Name: "attr2"},
				},
			},
			{
				Name:    "nest1_both",
				Parents: []string{"nest1_both"},
				Attributes: []Attribute{
					{Name: "attr2"},
				},
			},
			{
				Name:    "nest2",
				Parents: []string{"nest1_both", "nest2"},
				Attributes: []Attribute{
					{Name: "attr3"},
				},
			},
		},
	}

	got := NewSchema(in)
	if !reflect.DeepEqual(got, expect) {
		t.Errorf("unexpected Schema: \nexpect: %+v \ngot   : %+v", got, expect)
	}
}
