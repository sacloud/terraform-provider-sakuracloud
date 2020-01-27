---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_server_vnc_info"
subcategory: "Compute"
description: |-
  Get information about VNC for connecting to an existing Server.
---

# Data Source: sakuracloud_server_vnc_info

Get information about VNC for connecting to an existing Server.

## Example Usage

```hcl
data "sakuracloud_server_vnc_info" "foobar" {
  server_id = sakuracloud_server.foobar.id
}
```

## Argument Reference

* `server_id` - (Required) The id of the Server.
* `zone` - (Optional) The name of zone that the Server is in (e.g. `is1a`, `tk1a`).

## Attribute Reference

* `id` - The id of the Server VNC Information.
* `host` - The host name for connecting by VNC.
* `password` - The password for connecting by VNC.
* `port` - The port number for connecting by VNC.

