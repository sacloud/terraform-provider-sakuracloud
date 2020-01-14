---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_bridge"
subcategory: "Networking"
description: |-
  Get information about an existing Bridge.
---

# Data Source: sakuracloud_bridge

Get information about an existing Bridge.

## Example Usage

```hcl
data "sakuracloud_bridge" "foobar" {
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

---

A `condition` block supports the following:

* `name` - (Required) The name of the target field. This value is case-sensitive.
* `values` - (Required) The values of the condition. If multiple values ​​are specified, they combined as AND condition.


## Attribute Reference

* `id` - The id of the Bridge.
* `description` - The description of the Bridge.
* `name` - The name of the Bridge.
* `zone` - The name of zone that the Bridge is in (e.g. `is1a`, `tk1a`).




