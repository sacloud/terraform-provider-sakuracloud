---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_sim"
sidebar_current: "docs-sakuracloud-resource-secure-mobile-sim"
description: |-
  Provides a SakuraCloud SIM resource. This can be used to create, update, and delete SIMs.
---

# sakuracloud\_sim

Provides a SakuraCloud SIM resource. This can be used to create, update, and delete SIMs.

## Example Usage

```hcl
# Create a new SIM
resource "sakuracloud_sim" "foobar" {
  name     = "foobar"
  iccid    = "<your-iccid>"
  passcode = "<your-passcode>"

  # imei     = "<imei>"
  # enabled  = true

  # connect to the Mobile Gateway 
  mobile_gateway_id = sakuracloud_mobile_gateway.foobar.id
  ipaddress         = "192.168.2.1"

  description = "description"
  tags        = ["foo", "bar"]
}

# Create a new Mobile Gateway
resource "sakuracloud_mobile_gateway" "foobar" {
  name = "foobar"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `iccid` - (Required) The ICCID of the SIM.  
* `passcode` - (Required) The Passcode of the SIM.  
* `imei` - (Optional) The IMEI of the device that allows communication.
* `enabled` - (Optional) The flag of enable/disable the Server.
* `mobile_gateway_id` - (Optional) The ID of the Mobile Gateway to which the SIM belongs.
* `ipaddress` - (Optional) The IP address of the SIM. Used when connect to mobile gateway.
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `iccid` - The ICCID of the SIM. 
* `ipaddress` - The IP address of the SIM. Used when connected with mobile gateway.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon.

## Import

SIMs can be imported using the SIM ID.

```
$ terraform import sakuracloud_sim.foobar <sim_id>
```
