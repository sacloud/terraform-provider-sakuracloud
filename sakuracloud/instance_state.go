// Copyright 2016-2020 terraform-provider-sakuracloud authors
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
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// StringSliceFromState returns string slice from *terraform.InstanceState.Attributes
func StringSliceFromState(is *terraform.InstanceState, key string) []string {
	res := []string{}

	strCnt, ok := is.Attributes[fmt.Sprintf("%s.#", key)]
	if ok {
		count, err := strconv.Atoi(strCnt)
		if err != nil {
			return res
		}
		for i := 0; i < count; i++ {
			if v, ok := is.Attributes[fmt.Sprintf("%s.%d", key, i)]; ok {
				res = append(res, v)
			}
		}
	}

	return res
}
