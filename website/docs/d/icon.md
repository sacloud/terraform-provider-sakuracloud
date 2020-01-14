---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_icon"
subcategory: "Misc"
description: |-
  Get information about an existing Icon.
---

# Data Source: sakuracloud_icon

Get information about an existing Icon.

## Example Usage

```hcl
data "sakuracloud_icon" "foobar" {
  filter {
    names = ["foobar"]
  }
}
```
## Argument Reference

* `filter` - (Optional) One or more values used for filtering, as defined below.


---

A `filter` block supports the following:

* `condition` - (Optional) One or more name/values pairs used for filtering. There are several valid keys, for a full reference, check out finding section in the [SakuraCloud API reference](https://developer.sakura.ad.jp/cloud/api/1.1/).
* `id` - (Optional) The resource id on SakuraCloud used for filtering.
* `names` - (Optional) The resource names on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition.
* `tags` - (Optional) The resource tags on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition.

---

A `condition` block supports the following:

* `name` - (Required) The name of the target field. This value is case-sensitive.
* `values` - (Required) The values of the condition. If multiple values ​​are specified, they combined as AND condition.


## Attribute Reference

* `id` - The id of the Icon.
* `name` - The name of the Icon.
* `tags` - Any tags assigned to the Icon.
* `url` - The URL for getting the icon's raw data.




