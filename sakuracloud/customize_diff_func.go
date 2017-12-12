package sakuracloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"reflect"
	"regexp"
	"sort"
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

func ignoreIfMatchedChangeDiff(key string, regForOld *regexp.Regexp, regForNew *regexp.Regexp) schema.CustomizeDiffFunc {
	return func(d *schema.ResourceDiff, meta interface{}) error {
		if d.HasChange(key) {
			o, n := d.GetChange(key)
			if o != nil && n != nil {

				check := func(v interface{}, r *regexp.Regexp) bool {
					if r == nil {
						return true
					}
					return r.MatchString(v.(string))
				}

				if check(o, regForOld) && check(n, regForNew) {
					d.Clear(key)
				}
			}
		}
		return nil
	}
}
