// Copyright 2016-2022 terraform-provider-sakuracloud authors
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
	"fmt"
	"strings"
)

func quoteAndJoin(a []string, quote, sep string) string {
	if quote == "" {
		quote = `"`
	}
	var ret []string
	for _, w := range a {
		ret = append(ret, fmt.Sprintf("%s%s%s", quote, w, quote))
	}
	return strings.Join(ret, sep)
}

func quoteAndJoinInt(a []int, quote, sep string) string {
	if quote == "" {
		quote = `"`
	}
	var ret []string
	for _, w := range a {
		ret = append(ret, fmt.Sprintf("%s%d%s", quote, w, quote))
	}
	return strings.Join(ret, sep)
}

func descf(format string, a ...interface{}) string {
	args := make([]interface{}, len(a))
	for i, a := range a {
		var v interface{}
		switch a := a.(type) {
		case []string:
			v = quoteAndJoin(a, "`", "/")
		case []int:
			v = quoteAndJoinInt(a, "`", "/")
		default:
			v = a
		}
		args[i] = v
	}
	return fmt.Sprintf(format, args...)
}

func descRange(min, max int) string {
	return fmt.Sprintf("This must be in the range [`%d`-`%d`]", min, max)
}

func descLength(min, max int) string {
	return fmt.Sprintf("The length of this value must be in the range [`%d`-`%d`]", min, max)
}

func descConflicts(names ...string) string {
	return descf("This conflicts with [%s]", names)
}

func descResourcePlan(resourceName string, plans interface{}) string {
	return descf(
		"The plan name of the %s. This must be one of [%s]",
		resourceName, plans,
	)
}

func descDataSourcePlan(resourceName string, plans interface{}) string {
	return descf(
		"The plan name of the %s. This will be one of [%s]",
		resourceName, plans,
	)
}
