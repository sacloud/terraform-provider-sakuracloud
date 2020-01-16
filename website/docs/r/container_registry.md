---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_container_registry"
subcategory: "Lab"
description: |-
  Manages a SakuraCloud Container Registry.
---

# sakuracloud_container_registry

Manages a SakuraCloud Container Registry.

## Example Usage

```hcl
variable users {
  type = list(object({
    name     = string
    password = string
  }))
  default = [
    {
      name     = "user1"
      password = "password1"
    },
    {
      name     = "user2"
      password = "password2"
    }
  ]
}

resource "sakuracloud_container_registry" "foobar" {
  name            = "foobar"
  subdomain_label = "your-subdomain-label"
  access_level    = "readwrite" # this must be one of ["readwrite"/"readonly"/"none"]

  description = "description"
  tags        = ["tag1", "tag2"]

  dynamic user {
    for_each = var.users
    content {
      name     = user.value.name
      password = user.value.password
    }
  }
}
```
## Argument Reference

* `name` - (Required) The name of the Container Registry. The length of this value must be in the range [`1`-`64`].
* `access_level` - (Required) The level of access that allow to users. This must be one of [`readwrite`/`readonly`/`none`].
* `subdomain_label` - (Required) The label at the lowest of the FQDN used when be accessed from users. The length of this value must be in the range [`1`-`64`]. Changing this forces a new resource to be created.
* `user` - (Optional) One or more `user` blocks as defined below.

#### Common Arguments

* `description` - (Optional) The description of the Container Registry. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the Container Registry.
* `tags` - (Optional) Any tags to assign to the Container Registry.

---

A `user` block supports the following:

* `name` - (Required) The user name used to authenticate remote access.
* `password` - (Required) The password used to authenticate remote access.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Container Registry
* `update` - (Defaults to 5 minutes) Used when updating the Container Registry
* `delete` - (Defaults to 5 minutes) Used when deleting Container Registry

## Attribute Reference

* `id` - The id of the Container Registry.
* `fqdn` - The FQDN for accessing the Container Registry. FQDN is built from `subdomain_label` + `.sakuracr.jp`.



