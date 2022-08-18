---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_esme"
subcategory: "SMS"
description: |-
  Get information about an existing ESME.
---

# Data Source: sakuracloud_esme

Get information about an existing ESME.

## Example Usage

```hcl
data "sakuracloud_esme" "foobar" {
  filter {
    names = ["foobar"]
  }
}
```

## Argument Reference

* `filter` - (Optional) One or more values used for filtering, as defined below.

---

A `filter` block supports the following:

* `condition` - (Optional) One or more name/values pairs used for filtering. There are several valid keys, for a full reference, check out finding section in the [SakuraCloud API reference](https://developer.sakura.ad.jp/cloud/api/1.1/).
* `id` - (Optional) The resource id on SakuraCloud used for filtering.
* `names` - (Optional) The resource names on SakuraCloud used for filtering. If multiple values are specified, they combined as AND condition.
* `tags` - (Optional) The resource tags on SakuraCloud used for filtering. If multiple values are specified, they combined as AND condition.

---

A `condition` block supports the following:

* `name` - (Required) The name of the target field. This value is case-sensitive.
* `values` - (Required) The values of the condition. If multiple values are specified, they combined as AND condition.
* `operator` - (Optional) The filtering operator. This must be one of following: `partial_match_and`/`exact_match_or`. Default: `partial_match_and`


## Attribute Reference

* `id` - The id of the ESME.
* `description` - The description of the ESME.
* `icon_id` - The icon id attached to the ESME.
* `name` - The name of the ESME.
* `send_message_with_generated_otp_api_url` - The API URL for send SMS with generated OTP.
* `send_message_with_inputted_otp_api_url` - The API URL for send SMS with inputted OTP.
* `tags` - Any tags assigned to the ESME.



