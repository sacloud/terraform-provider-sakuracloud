---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_ssh_key"
subcategory: "Misc"
description: |-
  Get information about an existing SSH Key.
---

# Data Source: sakuracloud_ssh_key

Get information about an existing SSH Key.

## Example Usage

```hcl
data "sakuracloud_ssh_key" "foobar" {
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

* `id` - The id of the SSH Key.
* `description` - The description of the SSHKey.
* `fingerprint` - The fingerprint of public key.
* `name` - The name of the SSHKey.
* `public_key` - The value of public key.




