---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_container_registry"
subcategory: "Lab"
description: |-
  Get information about an existing Container Registry.
---

# Data Source: sakuracloud_container_registry

Get information about an existing Container Registry.

## Example Usage

```hcl
data "sakuracloud_container_registry" "foobar" {
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
* `operator` - (Optional) The filtering operator. This must be one of following: `partial_match_and`/`exact_match_or`. Default: `partial_match_and`


## Attribute Reference

* `id` - The id of the Container Registry.
* `access_level` - The level of access that allow to users. This will be one of [`readwrite`/`readonly`/`none`].
* `description` - The description of the ContainerRegistry.
* `fqdn` - The FQDN for accessing the container registry. FQDN is built from `subdomain_label` + `.sakuracr.jp`.
* `icon_id` - The icon id attached to the ContainerRegistry.
* `name` - The name of the ContainerRegistry.
* `subdomain_label` - The label at the lowest of the FQDN used when be accessed from users.
* `tags` - Any tags assigned to the ContainerRegistry.
* `user` - A list of `user` blocks as defined below.
* `virtual_domain` - The alias for accessing the container registry.


---

A `user` block exports the following:

* `name` - The user name used to authenticate remote access.
* `permission` - The level of access that allow to the user. This will be one of [`all`/`readwrite`/`readonly`].

