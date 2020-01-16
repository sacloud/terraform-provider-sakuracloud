---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_proxylb"
subcategory: "Global"
description: |-
  Get information about an existing ProxyLB.
---

# Data Source: sakuracloud_proxylb

Get information about an existing ProxyLB.

## Example Usage

```hcl
data "sakuracloud_proxylb" "foobar" {
  filter {
    names = ["foobar"]
  }
}
```
## Argument Reference

* `filter` - (Optional) One or more values used for filtering, as defined below.


---

A `certificate` block supports the following:

* `additional_certificate` - (Optional) One or more `additional_certificate` blocks as defined below.

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


## Attribute Reference

* `id` - The id of the ProxyLB.
* `bind_port` - A list of `bind_port` blocks as defined below.
* `certificate` - A list of `certificate` blocks as defined below.
* `description` - The description of the ProxyLB.
* `fqdn` - The FQDN for accessing to the ProxyLB. This is typically used as value of CNAME record.
* `health_check` - A list of `health_check` blocks as defined below.
* `icon_id` - The icon id attached to the ProxyLB.
* `name` - The name of the ProxyLB.
* `plan` - The plan name of the ProxyLB. This will be one of [`100`/`500`/`1000`/`5000`/`10000`/`50000`/`100000`].
* `proxy_networks` - A list of CIDR block used by the ProxyLB to access the server.
* `region` - The name of region that the proxy LB is in. This will be one of [`tk1`/`is1`].
* `rule` - A list of `rule` blocks as defined below.
* `server` - A list of `server` blocks as defined below.
* `sorry_server` - A list of `sorry_server` blocks as defined below.
* `sticky_session` - The flag to enable sticky session.
* `tags` - Any tags assigned to the ProxyLB.
* `timeout` - The timeout duration in seconds.
* `vip` - The virtual IP address assigned to the ProxyLB.
* `vip_failover` - The flag to enable VIP fail-over.


---

A `bind_port` block exports the following:

* `port` - The number of listening port.
* `proxy_mode` - The proxy mode. This will be one of [`http`/`https`/`tcp`].
* `redirect_to_https` - The flag to enable redirection from http to https. This flag is used only when `proxy_mode` is `http`.
* `response_header` - A list of `response_header` blocks as defined below.
* `support_http2` - The flag to enable HTTP/2. This flag is used only when `proxy_mode` is `https`.

---

A `response_header` block exports the following:

* `header` - The field name of HTTP header added to response by the ProxyLB.
* `value` - The field value of HTTP header added to response by the ProxyLB.

---

A `certificate` block exports the following:

* `intermediate_cert` - The intermediate certificate for a server.
* `private_key` - The private key for a server.
* `server_cert` - The certificate for a server.

---

A `additional_certificate` block exports the following:

* `intermediate_cert` - The intermediate certificate for a server.
* `private_key` - The private key for a server.
* `server_cert` - The certificate for a server.

---

A `health_check` block exports the following:

* `delay_loop` - The interval in seconds between checks.
* `host_header` - The value of host header send when checking by HTTP.
* `path` - The path used when checking by HTTP.
* `port` - The port number used when checking by TCP.
* `protocol` - The protocol used for health checks. This will be one of [`http`/`tcp`].

---

A `rule` block exports the following:

* `group` - The name of load balancing group. When proxyLB received request which matched to `host` and `path`, proxyLB forwards the request to servers that having same group name.
* `host` - The value of HTTP host header that is used as condition of rule-based balancing.
* `path` - The request path that is used as condition of rule-based balancing.

---

A `server` block exports the following:

* `enabled` - The flag to enable as destination of load balancing.
* `group` - The name of load balancing group. This is used when using rule-based load balancing.
* `ip_address` - The IP address of the destination server.
* `port` - The port number of the destination server.

---

A `sorry_server` block exports the following:

* `ip_address` - The IP address of the SorryServer. This will be used when all servers are down.
* `port` - The port number of the SorryServer. This will be used when all servers are down.


