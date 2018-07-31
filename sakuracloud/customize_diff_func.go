package sakuracloud

import (
	"reflect"
	"sort"

	"github.com/hashicorp/terraform/helper/schema"
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
	client := getSacloudAPIClient(d, meta)
	if d.HasChange("tags") {
		o, n := d.GetChange("tags")
		if o != nil && n != nil {
			os := realTags(client, expandTags(client, o.([]interface{})))
			ns := realTags(client, expandTags(client, n.([]interface{})))

			sort.Strings(os)
			sort.Strings(ns)
			if reflect.DeepEqual(os, ns) {
				return d.Clear("tags")
			}
		}
	}
	return nil
}
