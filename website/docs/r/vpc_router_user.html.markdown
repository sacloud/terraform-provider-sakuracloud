---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router_user"
sidebar_current: "docs-sakuracloud-resource-vpc-user"
description: |-
  Provides a SakuraCloud VPC Router User resource. This can be used to create and delete VPC Router User.
---

# sakuracloud\_vpc\_router\_user

Provides a SakuraCloud VPC Router User resource. This can be used to create and delete VPC Router User.

## Example Usage

```hcl
# Create a new VPC Router(standard)
resource "sakuracloud_vpc_router" "foobar" {
  name = "foobar"
}

# Create a new VPC Router User.
resource "sakuracloud_vpc_router_user" "user" {
  vpc_router_id = sakuracloud_vpc_router.foobar.id
  name          = "<user-name>"
  password      = "<p@ssword>"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_router_id` - (Required) The ID of the Internet resource.
* `name` - (Required) The user name.
* `password` - (Required) The password.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The user name.
* `password` - The password.
* `zone` - The ID of the zone to which the resource belongs.

## Import (not supported)

Import of VPC Router User is not supported.
