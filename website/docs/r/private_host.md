---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_private_host"
subcategory: "Compute"
description: |-
  Manages a SakuraCloud Private Host.
---

# sakuracloud_private_host

Manages a SakuraCloud Private Host.

## Example Usage

```hcl
resource "sakuracloud_private_host" "foobar" {
  name        = "foobar"
  description = "description"
  tags        = ["tag1", "tag2"]
}
```
## Argument Reference

* `class` - (Optional) The class of the PrivateHost. This will be one of [`dynamic`/`ms_windows`]. Default:`dynamic`.
* `description` - (Optional) The description of the PrivateHost. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the PrivateHost.
* `name` - (Required) The name of the PrivateHost. The length of this value must be in the range [`1`-`64`].
* `tags` - (Optional) Any tags to assign to the PrivateHost.
* `zone` - (Optional) The name of zone that the PrivateHost will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Private Host


* `update` - (Defaults to 5 minutes) Used when updating the Private Host

* `delete` - (Defaults to 20 minutes) Used when deregistering Private Host



## Attribute Reference

* `id` - The id of the Private Host.
* `assigned_core` - The total number of CPUs assigned to servers on the private host.
* `assigned_memory` - The total size of memory assigned to servers on the private host.
* `hostname` - The hostname of the private host.



