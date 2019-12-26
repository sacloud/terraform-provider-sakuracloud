---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_proxylb_acme"
subcategory: "Global"
description: |-
  Manages a SakuraCloud ProxyLB ACME Setting.
---

# sakuracloud_proxylb_acme

Manages a SakuraCloud ProxyLB ACME Setting.

## Argument Reference

* `accept_tos` - (Required) If you set this flag to true, you accept the current Let's Encrypt terms of service(see: https://letsencrypt.org/repository/). Changing this forces a new resource to be created.
* `common_name` - (Required) . Changing this forces a new resource to be created.
* `proxylb_id` - (Required) . Changing this forces a new resource to be created.
* `update_delay_sec` - (Optional) . Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 minutes) Used when creating the ProxyLB ACME Setting

* `read` -   (Defaults to 5 minutes) Used when reading the ProxyLB ACME Setting


* `delete` - (Defaults to 5 minutes) Used when deregistering ProxyLB ACME Setting



## Attribute Reference

* `id` - The ID of the ProxyLB ACME Setting.
* `certificate` - A list of `certificate` blocks as defined below.


---

A `certificate` block exports the following:

* `additional_certificate` - A list of `additional_certificate` blocks as defined below.
* `intermediate_cert` - .
* `private_key` - .
* `server_cert` - .

---

A `additional_certificate` block exports the following:

* `intermediate_cert` - .
* `private_key` - .
* `server_cert` - .



