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
		expect     bool
	}{
		{
			targetName: "foobar",
			conds:      []string{"bar"},
			expect:     true,
		},
		{
			targetName: "foobar",
			conds:      []string{"foo"},
			expect:     true,
		},
		{
			targetName: "foobar",
			conds:      []string{"foo", "bar"},
			expect:     true,
		},
		{
			targetName: "foobar",
			conds:      []string{"foo1", "bar2"},
			expect:     false,
		},
		{
			targetName: "foobar",
			conds:      []string{"foo1", "bar"},
			expect:     false,
		},
		{
			targetName: "foobar",
			conds:      []string{"foo", "bar1"},
			expect:     false,
		},
	}

	for _, e := range expects {
		target := &testNameFilterable{name: e.targetName}
		assert.Equal(t, e.expect, hasNames(target, e.conds))
	}

}

func TestHasTagsFilter(t *testing.T) {
	expects := []struct {
		tags   []string
		conds  []string
		expect bool
	}{
		{
			tags:   []string{"tag1"},
			conds:  []string{"tag1"},
			expect: true,
		},
		{
			tags:   []string{"tag1"},
			conds:  []string{"tag2"},
			expect: false,
		},
		{
			tags:   []string{"tag1"},
			conds:  []string{"t"},
			expect: false,
		},
		{
			tags:   []string{"tag1", "tag2"},
			conds:  []string{"tag2"},
			expect: true,
		},
		{
			tags:   []string{"tag1", "tag2"},
			conds:  []string{"tag1", "t"},
			expect: false,
		},
	}
	for _, e := range expects {
		target := &testTagFilterable{tags: e.tags}
		assert.Equal(t, e.expect, hasTags(target, e.conds))
	}
}
