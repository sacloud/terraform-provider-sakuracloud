---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_bridge"
sidebar_current: "docs-sakuracloud-resource-networking-bridge"
description: |-
  Provides a SakuraCloud Bridge resource. This can be used to create, update, and delete Bridges.
---

# sakuracloud\_bridge

Provides a SakuraCloud Bridge resource. This can be used to create, update, and delete Bridges.

## Example Usage

```hcl
# Create a new Bridge
resource "sakuracloud_bridge" "foobar" {
  name        = "foobar"
  switch_ids  = [sakuracloud_switch.foobar.id]
  description = "description"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `switch_ids` - (Optional) The ID list of the Switches connected to the Bridge. 
* `description` - (Optional) The description of the resource.
* `zone` - (Optional) The ID of the zone to which the resource belongs.  

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `switch_ids` - The ID list of the Switches connected to the Bridge. 
* `description` - The description of the resource.
* `zone` - The ID of the zone to which the resource belongs.

## Import

Bridges can be imported using the Bridge ID.

```
$ terraform import sakuracloud_bridge.foobar <bridge_id>
```
