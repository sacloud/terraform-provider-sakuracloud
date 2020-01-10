---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_ipv4_ptr"
subcategory: "Networking"
description: |-
  Manages a SakuraCloud IPv4 PTR.
---

# sakuracloud_ipv4_ptr

Manages a SakuraCloud IPv4 PTR.

## Argument Reference

* `hostname` - (Required) The value of the PTR record. This must be FQDN.
* `ip_address` - (Required) The IP address to which the PTR record is set.
* `retry_interval` - (Optional) The wait interval(in seconds) for retrying API call used when SakuraCloud API returns any errors. Default:`10`.
* `retry_max` - (Optional) The maximum number of API call retries used when SakuraCloud API returns any errors. Default:`30`.
* `zone` - (Optional) The name of zone that the IPv4 PTR will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the IPv4 PTR


* `update` - (Defaults to 60 minutes) Used when updating the IPv4 PTR

* `delete` - (Defaults to 5 minutes) Used when deregistering IPv4 PTR



## Attribute Reference

* `id` - The id of the IPv4 PTR.




