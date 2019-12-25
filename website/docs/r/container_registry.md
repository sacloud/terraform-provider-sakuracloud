---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_container_registry"
subcategory: "Lab"
description: |-
  Manages a SakuraCloud Container Registry.
---

# sakuracloud_container_registry

Manages a SakuraCloud Container Registry.

## Argument Reference

* `access_level` - (Required) .
* `description` - (Optional) .
* `icon_id` - (Optional) .
* `name` - (Required) .
* `subdomain_label` - (Required) . Changing this forces a new resource to be created.
* `tags` - (Optional) .
* `user` - (Optional) One or more `user` blocks as defined below.


---

A `user` block supports the following:

* `name` - (Required) .
* `password` - (Required) .


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Container Registry

* `read` -   (Defaults to 5 minutes) Used when reading the Container Registry

* `update` - (Defaults to 5 minutes) Used when updating the Container Registry

* `delete` - (Defaults to 5 minutes) Used when deregistering Container Registry



## Attribute Reference

* `id` - The ID of the Container Registry.
* `fqdn` - .




