---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_auto_scale"
subcategory: "Misc"
description: |-
  Get information about an existing sakuracloud_auto_scale.
---

# Data Source: sakuracloud_auto_scale

Get information about an existing sakuracloud_auto_scale.

## Example Usage

```hcl
data "sakuracloud_auto_scale" "foobar" {
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
* `names` - (Optional) The resource names on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition.
* `tags` - (Optional) The resource tags on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition.

---

A `condition` block supports the following:

* `name` - (Required) The name of the target field. This value is case-sensitive.
* `values` - (Required) The values of the condition. If multiple values ​​are specified, they combined as AND condition.
* `operator` - (Optional) The filtering operator. This must be one of following: `partial_match_and`/`exact_match_or`. Default: `partial_match_and`

## Attribute Reference

* `id` - The id of the sakuracloud_auto_scale.
* `api_key_id` - The id of the API key.
* `config` - The configuration file for sacloud/autoscaler.
* `cpu_threshold_scaling` - A list of `cpu_threshold_scaling` blocks as defined below.
* `description` - The description of the AutoScale.
* `icon_id` - The icon id attached to the AutoScale.
* `name` - The name of the AutoScale.
* `router_threshold_scaling` - A list of `router_threshold_scaling` blocks as defined below.
* `tags` - Any tags assigned to the AutoScale.
* `trigger_type` - This must be one of [`cpu`/`router`].
* `zones` - List of zone names where monitored resources are located.

---

A `cpu_threshold_scaling` block exports the following:

* `server_prefix` - Server name prefix to be monitored.
* `up` - Threshold for average CPU utilization to scale up/out.

---

A `router_threshold_scaling` block exports the following:

* `direction` - This must be one of [`in`/`out`].
* `mbps` - Mbps.
* `router_prefix` - Router name prefix to be monitored.


