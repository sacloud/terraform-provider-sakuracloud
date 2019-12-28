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

var (
	resourceDescriptions = map[string]string{
		"zone": "The name of zone that the %s will be created. (e.g. `is1a`,`tk1a`)",
	}
)
