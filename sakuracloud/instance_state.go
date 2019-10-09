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
