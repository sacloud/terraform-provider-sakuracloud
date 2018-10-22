---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_private_host"
sidebar_current: "docs-sakuracloud-resource-computing-private-host"
description: |-
  Provides a SakuraCloud Private Host resource. This can be used to create, update, and delete Private Hosts.
---

# sakuracloud\_private\_host

Provides a SakuraCloud Private Host resource. This can be used to create, update, and delete Private Hosts.

## Example Usage

```hcl
# Create a new Private Host
resource "sakuracloud_private_host" "foobar" {
  name        = "foobar"
  description = "description"
  tags        = ["foo", "bar"]
}

# Add server on Private Host
resource "sakuracloud_server" "foobar" {
  name            = "example"
  private_host_id = sakuracloud_private_host.foobar.id
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.
* `graceful_shutdown_timeout` - (Optional) The wait time (seconds) to do graceful shutdown the server connected to the resource.
* `zone` - (Optional) The ID of the zone to which the resource belongs.  
Valid value is one of the following: ["is1b" / "tk1a"]

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `hostname` - The HostName of the resource.
* `assigned_core` - The number of cores assigned to the Server.
* `assigned_memory` - The size of memory allocated to the Server (unit:`GB`).
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `zone` - The ID of the zone to which the resource belongs.

## Import

Private Hosts can be imported using the Private Host ID.

```
$ terraform import sakuracloud_private_host.foobar <private_host_id>
```
