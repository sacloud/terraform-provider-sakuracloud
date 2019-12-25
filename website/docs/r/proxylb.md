---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_proxylb"
subcategory: "Global"
description: |-
  Manages a SakuraCloud ProxyLB.
---

# sakuracloud_proxylb

Manages a SakuraCloud ProxyLB.

## Argument Reference

* `bind_port` - (Required) One or more `bind_port` blocks as defined below.
* `certificate` - (Optional) An `certificate` block as defined below.
* `description` - (Optional) .
* `health_check` - (Required) A `health_check` block as defined below.
* `icon_id` - (Optional) .
* `name` - (Required) .
* `plan` - (Optional) . Defaults to `100`.
* `region` - (Optional) . Changing this forces a new resource to be created. Defaults to `is1`.
* `rule` - (Optional) One or more `rule` blocks as defined below.
* `server` - (Optional) One or more `server` blocks as defined below.
* `sorry_server` - (Optional) A `sorry_server` block as defined below.
* `sticky_session` - (Optional) .
* `tags` - (Optional) .
* `timeout` - (Optional) . Defaults to `10`.
* `vip_failover` - (Optional) . Changing this forces a new resource to be created.


---

A `bind_port` block supports the following:

* `port` - (Optional) .
* `proxy_mode` - (Required) .
* `redirect_to_https` - (Optional) .
* `response_header` - (Optional) One or more `response_header` blocks as defined below.
* `support_http2` - (Optional) .

---

A `response_header` block supports the following:

* `header` - (Required) .
* `value` - (Required) .

---

A `certificate` block supports the following:

* `additional_certificate` - (Optional) One or more `additional_certificate` blocks as defined below.
* `intermediate_cert` - (Optional) .
* `private_key` - (Optional) .
* `server_cert` - (Optional) .

---

A `additional_certificate` block supports the following:

* `intermediate_cert` - (Optional) .
* `private_key` - (Required) .
* `server_cert` - (Required) .

---

A `health_check` block supports the following:

* `delay_loop` - (Optional) .
* `host_header` - (Optional) .
* `path` - (Optional) .
* `protocol` - (Required) .

---

A `rule` block supports the following:

* `group` - (Optional) .
* `host` - (Optional) .
* `path` - (Optional) .

---

A `server` block supports the following:

* `enabled` - (Optional) .
* `group` - (Optional) .
* `ip_address` - (Required) .
* `port` - (Required) .

---

A `sorry_server` block supports the following:

* `ip_address` - (Required) .
* `port` - (Optional) .


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the ProxyLB

* `read` -   (Defaults to 5 minutes) Used when reading the ProxyLB

* `update` - (Defaults to 5 minutes) Used when updating the ProxyLB

* `delete` - (Defaults to 5 minutes) Used when deregistering ProxyLB



## Attribute Reference

* `id` - The ID of the ProxyLB.
* `fqdn` - .
* `proxy_networks` - .
* `vip` - .




