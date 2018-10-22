---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_switch"
sidebar_current: "docs-sakuracloud-resource-networking-switch"
description: |-
  Provides a SakuraCloud Switch resource. This can be used to create, update, and delete Switches.
---

# sakuracloud\_switch

Provides a SakuraCloud Switch resource. This can be used to create, update, and delete Switches.

## Example Usage

```hcl
# Create a new Switch
resource "sakuracloud_switch" "foobar" {
  name        = "foobar"
  description = "description"
  tags        = ["foo", "bar"]

  # If you want to connect to the bridge, please uncomment here.
  #bridge_id = sakuracloud_bridge.br.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `bridge_id` - (Optional) The ID of the Bridge to connect to the Switch.
* `icon_id` - (Optional) The ID of the icon.
* `graceful_shutdown_timeout` - (Optional) The wait time (seconds) to do graceful shutdown the server connected to the resource.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `server_ids` - The ID list of the servers connected to the switch.
* `bridge_id` - The ID of the bridge connected to the switch.
* `icon_id` - The ID of the icon.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `zone` - The ID of the zone to which the resource belongs.

## Import

Switches can be imported using the Switch ID.

```
$ terraform import sakuracloud_switch.foobar <switch_id>
```
