---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_cdrom"
subcategory: "Storage"
description: |-
  Get information about an existing CD-ROM.
---

# Data Source: sakuracloud_cdrom

Get information about an existing CD-ROM.

## Example Usage

```hcl
data "sakuracloud_cdrom" "foobar" {
  filter {
    condition {
      name   = "Name"
      values = ["Parted Magic 2013_08_01"]
    }
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

* `id` - The id of the CD-ROM.
* `description` - The description of the CD-ROM.
* `icon_id` - The icon id attached to the CD-ROM.
* `name` - The name of the CD-ROM.
* `size` - The size of CD-ROM in GiB.
* `tags` - Any tags assigned to the CD-ROM.
* `zone` - The name of zone that the CD-ROM is in (e.g. `is1a`, `tk1a`).




