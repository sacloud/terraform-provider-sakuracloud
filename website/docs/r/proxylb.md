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
* `description` - (Optional) The description of the ProxyLB. The length of this value must be in the range [`1`-`512`].
* `health_check` - (Required) A `health_check` block as defined below.
* `icon_id` - (Optional) The icon id to attach to the ProxyLB.
* `name` - (Required) The name of the ProxyLB. The length of this value must be in the range [`1`-`64`].
* `plan` - (Optional) The plan name of the ProxyLB. This must be one of [`100`/`500`/`1000`/`5000`/`10000`/`50000`/`100000`]. Changing this forces a new resource to be created. Default:`100`.
* `region` - (Optional) The name of region that the proxy LB is in. This must be one of [`tk1`/`is1`]. Changing this forces a new resource to be created. Default:`is1`.
* `rule` - (Optional) One or more `rule` blocks as defined below.
* `server` - (Optional) One or more `server` blocks as defined below.
* `sorry_server` - (Optional) A `sorry_server` block as defined below.
* `sticky_session` - (Optional) The flag to enable sticky session.
* `tags` - (Optional) Any tags to assign to the ProxyLB.
* `timeout` - (Optional) The timeout duration in seconds. Default:`10`.
* `vip_failover` - (Optional) The flag to enable VIP fail-over. Changing this forces a new resource to be created.


---

A `bind_port` block supports the following:

* `port` - (Optional) The number of listening port.
* `proxy_mode` - (Required) The proxy mode. This must be one of [`http`/`https`/`tcp`].
* `redirect_to_https` - (Optional) The flag to enable redirection from http to https. This flag is used only when `proxy_mode` is `http`.
* `response_header` - (Optional) One or more `response_header` blocks as defined below.
* `support_http2` - (Optional) The flag to enable HTTP/2. This flag is used only when `proxy_mode` is `https`.

---

A `response_header` block supports the following:

* `header` - (Required) The field name of HTTP header added to response by the ProxyLB.
* `value` - (Required) The field value of HTTP header added to response by the ProxyLB.

---

A `certificate` block supports the following:

* `additional_certificate` - (Optional) One or more `additional_certificate` blocks as defined below.
* `intermediate_cert` - (Optional) The intermediate certificate for a server.
* `private_key` - (Optional) The private key for a server.
* `server_cert` - (Optional) The certificate for a server.

---

A `additional_certificate` block supports the following:

* `intermediate_cert` - (Optional) The intermediate certificate for a server.
* `private_key` - (Required) The private key for a server.
* `server_cert` - (Required) The certificate for a server.

---

A `health_check` block supports the following:

* `delay_loop` - (Optional) The interval in seconds between checks. This must be in the range [`10`-`60`].
* `host_header` - (Optional) The value of host header send when checking by HTTP.
* `path` - (Optional) The path used when checking by HTTP.
* `port` - (Optional) The port number used when checking by TCP.
* `protocol` - (Required) The protocol used for health checks. This must be one of [`http`/`tcp`].

---

A `rule` block supports the following:

* `group` - (Optional) The name of load balancing group. When proxyLB received request which matched to `host` and `path`, proxyLB forwards the request to servers that having same group name. The length of this value must be in the range [`1`-`10`].
* `host` - (Optional) The value of HTTP host header that is used as condition of rule-based balancing.
* `path` - (Optional) The request path that is used as condition of rule-based balancing.

---

A `server` block supports the following:

* `enabled` - (Optional) The flag to enable as destination of load balancing.
* `group` - (Optional) The name of load balancing group. This is used when using rule-based load balancing. The length of this value must be in the range [`1`-`10`].
* `ip_address` - (Required) The IP address of the destination server.
* `port` - (Required) The port number of the destination server. This must be in the range [`1`-`65535`].

---

A `sorry_server` block supports the following:

* `ip_address` - (Required) The IP address of the SorryServer. This will be used when all servers are down.
* `port` - (Optional) The port number of the SorryServer. This will be used when all servers are down.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the ProxyLB


* `update` - (Defaults to 5 minutes) Used when updating the ProxyLB

* `delete` - (Defaults to 5 minutes) Used when deregistering ProxyLB



## Attribute Reference

* `id` - The id of the ProxyLB.
* `fqdn` - The FQDN for accessing to the ProxyLB. This is typically used as value of CNAME record.
* `proxy_networks` - A list of CIDR block used by the ProxyLB to access the server.
* `vip` - The virtual IP address assigned to the ProxyLB.




