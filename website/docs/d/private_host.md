---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_private_host"
subcategory: "Compute"
description: |-
  Get information about an existing Private Host.
---

# Data Source: sakuracloud_private_host

Get information about an existing Private Host.

## Example Usage

```hcl
data "sakuracloud_private_host" "foobar" {
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

* `id` - The id of the Private Host.
* `assigned_core` - The total number of CPUs assigned to servers on the private host.
* `assigned_memory` - The total size of memory assigned to servers on the private host.
* `class` - The class of the PrivateHost. This will be one of [`dynamic`/`ms_windows`].
* `description` - The description of the PrivateHost.
* `hostname` - The hostname of the private host.
* `icon_id` - The icon id attached to the PrivateHost.
* `name` - The name of the PrivateHost.
* `tags` - Any tags assigned to the PrivateHost.
* `zone` - The name of zone that the PrivateHost is in (e.g. `is1a`, `tk1a`).



