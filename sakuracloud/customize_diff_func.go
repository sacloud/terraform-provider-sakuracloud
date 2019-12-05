package sakuracloud

import (
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func composeCustomizeDiff(funcs ...schema.CustomizeDiffFunc) schema.CustomizeDiffFunc {
	return func(d *schema.ResourceDiff, meta interface{}) error {
		for _, f := range funcs {
			err := f(d, meta)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func hasTagResourceCustomizeDiff(d *schema.ResourceDiff, meta interface{}) error {
	if d.HasChange("tags") {
		o, n := d.GetChange("tags")
		if o != nil && n != nil {
			os := expandStringList(o.([]interface{}))
			ns := expandStringList(n.([]interface{}))

			sort.Strings(os)
			sort.Strings(ns)
			if reflect.DeepEqual(os, ns) {
				return d.Clear("tags")
			}
		}
	}
	return nil
}
