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
	"errors"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

const filterAttrName = "filter"

type filterSchemaOption struct {
	excludeTags bool
}

var (
	filterConfigKeys = []string{
		"filter.0.id",
		"filter.0.names",
		"filter.0.condition",
	}
	filterConfigKeysWithTags = append(filterConfigKeys, "filter.0.tags")
	filteringOperators       = []string{
		filteringOperatorPartialMatchAnd,
		filteringOperatorExactMatchOr,
	}
)

const (
	filteringOperatorPartialMatchAnd = "partial_match_and"
	filteringOperatorExactMatchOr    = "exact_match_or"
)

func filterSchema(opt *filterSchemaOption) *schema.Schema {
	if opt == nil {
		opt = &filterSchemaOption{}
	}
	keys := filterConfigKeysWithTags
	if opt.excludeTags {
		keys = filterConfigKeys
	}
	s := map[string]*schema.Schema{
		"id": {
			Type:         schema.TypeString,
			Optional:     true,
			ExactlyOneOf: keys,
			Description:  "The resource id on SakuraCloud used for filtering",
		},
		"names": {
			Type:         schema.TypeList,
			Optional:     true,
			ExactlyOneOf: keys,
			Elem:         &schema.Schema{Type: schema.TypeString},
			Description:  "The resource names on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition",
		},
		"tags": {
			Type:         schema.TypeSet,
			Optional:     true,
			ExactlyOneOf: keys,
			Elem:         &schema.Schema{Type: schema.TypeString},
			Set:          schema.HashString,
			Description:  "The resource tags on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition",
		},
		"condition": {
			Type:         schema.TypeList,
			Optional:     true,
			ExactlyOneOf: keys,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The name of the target field. This value is case-sensitive",
					},

					"values": {
						Type:        schema.TypeList,
						Required:    true,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Description: "The values of the condition. If multiple values ​​are specified, they combined as AND condition",
					},
					"operator": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(filteringOperators, false)),
						Default:          filteringOperatorPartialMatchAnd,
						Description: desc.Sprintf(
							"The filtering operator. This must be one of following:  \n%s",
							filteringOperators,
						),
					},
				},
			},
			Description: "One or more name/values pairs used for filtering. There are several valid keys, for a full reference, check out finding section in the [SakuraCloud API reference](https://developer.sakura.ad.jp/cloud/api/1.1/)",
		},
	}
	if opt.excludeTags {
		delete(s, "tags")
	}
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: s,
		},
		Description: "One or more values used for filtering, as defined below",
	}
}

var errFilterNoResult = errors.New("Your query returned no results. Please change your filter or selectors and try again")

func filterNoResultErr() diag.Diagnostics {
	if os.Getenv(resource.EnvTfAcc) != "" {
		return nil
	}
	return diag.FromErr(errFilterNoResult)
}

type nameFilterable interface {
	GetName() string
}

func hasNames(target interface{}, cond []string) bool {
	t, ok := target.(nameFilterable)
	if !ok {
		return false
	}
	name := t.GetName()
	for _, c := range cond {
		if !strings.Contains(name, c) {
			return false
		}
	}
	return true
}

type tagFilterable interface {
	HasTag(string) bool
}

func hasTags(target interface{}, cond []string) bool {
	t, ok := target.(tagFilterable)
	if !ok {
		return false
	}
	for _, c := range cond {
		if !t.HasTag(c) {
			return false
		}
	}
	return true
}
