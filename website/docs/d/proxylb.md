---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_proxylb"
subcategory: "Global"
description: |-
  Get information about an existing ProxyLB.
---

# Data Source: sakuracloud_proxylb

Get information about an existing ProxyLB.

## Argument Reference

* `filter` - (Optional) A `filter` block as defined below.


---

A `certificate` block supports the following:

* `additional_certificate` - (Optional) One or more `additional_certificate` blocks as defined below.

---

A `filter` block supports the following:

* `condition` - (Optional) One or more `condition` blocks as defined below.
* `id` - (Optional) .
* `names` - (Optional) .
* `tags` - (Optional) .

---

A `condition` block supports the following:

* `name` - (Required) .
* `values` - (Required) .


## Attribute Reference

* `id` - The ID of the ProxyLB.
* `bind_port` - A list of `bind_port` blocks as defined below.
* `certificate` - A list of `certificate` blocks as defined below.
* `description` - .
* `fqdn` - .
* `health_check` - A list of `health_check` blocks as defined below.
* `icon_id` - .
* `name` - .
* `plan` - .
* `proxy_networks` - .
* `region` - .
* `rule` - A list of `rule` blocks as defined below.
* `server` - A list of `server` blocks as defined below.
* `sorry_server` - A list of `sorry_server` blocks as defined below.
* `sticky_session` - .
* `tags` - .
* `timeout` - .
* `vip` - .
* `vip_failover` - .


---

A `bind_port` block exports the following:

* `port` - .
* `proxy_mode` - .
* `redirect_to_https` - .
* `response_header` - A list of `response_header` blocks as defined below.
* `support_http2` - .

---

A `response_header` block exports the following:

* `header` - .
* `value` - .

---

A `certificate` block exports the following:

* `intermediate_cert` - .
* `private_key` - .
* `server_cert` - .

---

A `additional_certificate` block exports the following:

* `intermediate_cert` - .
* `private_key` - .
* `server_cert` - .

---

A `health_check` block exports the following:

* `delay_loop` - .
* `host_header` - .
* `path` - .
* `port` - .
* `protocol` - .

---

A `rule` block exports the following:

* `group` - .
* `host` - .
* `path` - .

---

A `server` block exports the following:

* `enabled` - .
* `group` - .
* `ip_address` - .
* `port` - .

---

A `sorry_server` block exports the following:

* `ip_address` - .
* `port` - .



