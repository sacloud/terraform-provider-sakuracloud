---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_sim"
subcategory: "SecureMobile"
description: |-
  Manages a SakuraCloud SIM.
---

# sakuracloud_sim

Manages a SakuraCloud SIM.

## Example Usage

```hcl
resource "sakuracloud_sim" "foobar" {
  name        = "foobar"
  description = "description"
  tags        = ["tag1", "tag2"]

  iccid    = "your-iccid"
  passcode = "your-password"
  #imei     = "your-imei"
  carrier = ["softbank", "docomo", "kddi"]

  enabled = true
}
```

## Argument Reference

* `name` - (Required) The name of the SIM. The length of this value must be in the range [`1`-`64`].
* `enabled` - (Optional) The flag to enable the SIM. Default:`true`.

#### SIM Registration

* `iccid` - (Required) ICCID(Integrated Circuit Card ID) assigned to the SIM. Changing this forces a new resource to be created.
* `passcode` - (Required) The passcord to authenticate the SIM. Changing this forces a new resource to be created.
* `carrier` - (Required) A list of a communication company. Each element must be one of `kddi`/`docomo`/`softbank`.

#### IMEI Locking

* `imei` - (Optional) The id of the device to restrict devices that can use the SIM.

#### Common Arguments

* `description` - (Optional) The description of the SIM. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the SIM.
* `tags` - (Optional) Any tags to assign to the SIM.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the SIM
* `update` - (Defaults to 5 minutes) Used when updating the SIM
* `delete` - (Defaults to 5 minutes) Used when deleting SIM

## Attribute Reference

* `id` - The id of the SIM.
* `ip_address` - The IP address assigned to the SIM.
* `mobile_gateway_id` - The id of the MobileGateway which the SIM is assigned.

