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

package desc

import (
	"fmt"
	"strings"
)

func QuoteAndJoin(a []string, quote, sep string) string {
	if quote == "" {
		quote = `"`
	}
	var ret []string
	for _, w := range a {
		ret = append(ret, fmt.Sprintf("%s%s%s", quote, w, quote))
	}
	return strings.Join(ret, sep)
}

func QuoteAndJoinInt(a []int, quote, sep string) string {
	if quote == "" {
		quote = `"`
	}
	var ret []string
	for _, w := range a {
		ret = append(ret, fmt.Sprintf("%s%d%s", quote, w, quote))
	}
	return strings.Join(ret, sep)
}

func Sprintf(format string, a ...interface{}) string {
	args := make([]interface{}, len(a))
	for i, a := range a {
		var v interface{}
		switch a := a.(type) {
		case []string:
			v = QuoteAndJoin(a, "`", "/")
		case []int:
			v = QuoteAndJoinInt(a, "`", "/")
		default:
			v = a
		}
		args[i] = v
	}
	return fmt.Sprintf(format, args...)
}

func Range(min, max int) string {
	return fmt.Sprintf("This must be in the range [`%d`-`%d`]", min, max)
}

func Length(min, max int) string {
	return fmt.Sprintf("The length of this value must be in the range [`%d`-`%d`]", min, max)
}

func Conflicts(names ...string) string {
	return Sprintf("This conflicts with [%s]", names)
}

func ResourcePlan(resourceName string, plans interface{}) string {
	return Sprintf(
		"The plan name of the %s. This must be one of [%s]",
		resourceName, plans,
	)
}

func DataSourcePlan(resourceName string, plans interface{}) string {
	return Sprintf(
		"The plan name of the %s. This will be one of [%s]",
		resourceName, plans,
	)
}
