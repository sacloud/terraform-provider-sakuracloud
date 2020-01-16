---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_nfs"
subcategory: "Appliance"
description: |-
  Get information about an existing NFS.
---

# Data Source: sakuracloud_nfs

Get information about an existing NFS.

## Example Usage

```hcl
data "sakuracloud_nfs" "foobar" {
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

* `id` - The id of the NFS.
* `description` - The description of the NFS.
* `gateway` - The IP address of the gateway used by NFS.
* `icon_id` - The icon id attached to the NFS.
* `ip_address` - The IP address assigned to the NFS.
* `name` - The name of the NFS.
* `netmask` - The bit length of the subnet assigned to the NFS.
* `plan` - The plan name of the NFS. This will be one of [`hdd`/`ssd`].
* `size` - The size of NFS in GiB.
* `switch_id` - The id of the switch connected from the NFS.
* `tags` - Any tags assigned to the NFS.
* `zone` - The name of zone that the NFS is in (e.g. `is1a`, `tk1a`).



