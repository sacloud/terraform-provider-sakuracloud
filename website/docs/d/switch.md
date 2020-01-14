---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_switch"
subcategory: "Networking"
description: |-
  Get information about an existing Switch.
---

# Data Source: sakuracloud_switch

Get information about an existing Switch.

## Example Usage

```hcl
data "sakuracloud_switch" "foobar" {
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

* `id` - The id of the Switch.
* `bridge_id` - The bridge id attached to the Switch.
* `description` - The description of the Switch.
* `icon_id` - The icon id attached to the Switch.
* `name` - The name of the Switch.
* `server_ids` - A list of server id connected to the Switch.
* `tags` - Any tags assigned to the Switch.
* `zone` - The name of zone that the Switch is in (e.g. `is1a`, `tk1a`).




