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

package sakuracloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testNameFilterable struct {
	name string
}

func (f *testNameFilterable) GetName() string {
	return f.name
}

type testTagFilterable struct {
	tags []string
}

func (f *testTagFilterable) HasTag(tag string) bool {
	for _, t := range f.tags {
		if t == tag {
			return true
		}
	}
	return false
}

func TestHasNamesFilter(t *testing.T) {
	expects := []struct {
		targetName string
		conds      []string
		hit        bool
	}{
		{
			targetName: "foobar",
			conds:      []string{"bar"},
			hit:        true,
		},
		{
			targetName: "foobar",
			conds:      []string{"foo"},
			hit:        true,
		},
		{
			targetName: "foobar",
			conds:      []string{"foo", "bar"},
			hit:        true,
		},
		{
			targetName: "foobar",
			conds:      []string{"foo1", "bar2"},
			hit:        false,
		},
		{
			targetName: "foobar",
			conds:      []string{"foo1", "bar"},
			hit:        false,
		},
		{
			targetName: "foobar",
			conds:      []string{"foo", "bar1"},
			hit:        false,
		},
	}

	for _, e := range expects {
		target := &testNameFilterable{name: e.targetName}
		assert.Equal(t, e.hit, hasNames(target, e.conds))
	}
}

func TestHasTagsFilter(t *testing.T) {
	expects := []struct {
		sources    []string
		conditions []string
		hit        bool
	}{
		{
			sources:    []string{"tag1"},
			conditions: []string{"tag1"},
			hit:        true,
		},
		{
			sources:    []string{"tag1"},
			conditions: []string{"tag2"},
			hit:        false,
		},
		{
			sources:    []string{"tag1"},
			conditions: []string{"t"},
			hit:        false,
		},
		{
			sources:    []string{"tag1", "tag2"},
			conditions: []string{"tag2"},
			hit:        true,
		},
		{
			sources:    []string{"tag1", "tag2"},
			conditions: []string{"tag1", "t"},
			hit:        false,
		},
	}
	for _, e := range expects {
		target := &testTagFilterable{tags: e.sources}
		assert.Equal(t, e.hit, hasTags(target, e.conditions))
	}
}
