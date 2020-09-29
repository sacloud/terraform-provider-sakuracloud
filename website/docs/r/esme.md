---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_esme"
subcategory: "SMS"
description: |-
  Manages a SakuraCloud sakuracloud_esme.
---

# sakuracloud_esme

Manages a SakuraCloud ESME resource.

## Example Usage

```hcl
resource "sakuracloud_esme" "foobar" {
  name        = "foobar"
  description = "description"
  tags        = ["tag1", "tag2"]
}
```

## Argument Reference

* `name` - (Required) The name of the ESME. The length of this value must be in the range [`1`-`64`].

#### Common Arguments

* `description` - (Optional) The description of the ESME. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the ESME.
* `tags` - (Optional) Any tags to assign to the ESME.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the sakuracloud_esme
* `update` - (Defaults to 5 minutes) Used when updating the sakuracloud_esme
* `delete` - (Defaults to 5 minutes) Used when deleting sakuracloud_esme

## Attribute Reference

* `id` - The id of the sakuracloud_esme.
* `send_message_with_generated_otp_api_url` - The API URL for send SMS with generated OTP.
* `send_message_with_inputted_otp_api_url` - The API URL for send SMS with inputted OTP.

