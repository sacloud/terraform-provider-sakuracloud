---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_proxylb"
sidebar_current: "docs-sakuracloud-datasource-proxylb"
description: |-
  Get information on a SakuraCloud ProxyLB.
---

# sakuracloud\_proxylb

Use this data source to retrieve information about a SakuraCloud ProxyLB.

## Example Usage

```hcl
data "sakuracloud_proxylb" "foobar" {
  name_selectors = ["foobar"]
}
```

## Argument Reference

 * `name_selectors` - (Optional) The list of names to filtering.
 * `tag_selectors` - (Optional) The list of tags to filtering.
 * `filter` - (Optional) The map of filter key and value.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - Name of the resource.
* `plan` - The plan of the resource.
* `vip_failover` - The flag of enable VIP Fail-Over.  
* `sticky_session` - The flag of enable Sticky-Session.  
* `bind_ports` - The external listen ports. It contains some attributes to [Bind Ports](#bind-ports).
* `health_check` - The health check rules. It contains some attributes to [Health Check](#health-check).
* `sorry_server` - The pair of IPAddress and port number of sorry-server.
* `servers` - Real servers. It contains some attributes to [Servers](#servers).
* `certificate` - Certificate used to terminate SSL/TSL. It contains some attributes to [Certificate](#certificate).
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon.

### Bind Ports

Attributes for Bind-Ports:

* `proxy_mode` - Proxy protocol.  
* `port` - Port number used in tcp proxy.
* `redirect_to_https` - The flag for enable to redirect to https.
* `support_http2` - The flag for enable to support HTTP/2.

### Health Check

Attributes for Health Check:

* `protocol` - Protocol used in health check.  
* `delay_loop` - Health check access interval (unit:`second`, default:`10`).
* `host_header` - The value of `Host` header used in http/https health check access.
* `path` - The request path used in http health check access.

### Servers

* `ipaddress` - The IP address of the Real-Server.
* `port` - Port number.
* `enabled` - The flag for enable/disable the Real-Server (default:`true`).

### Certificate

* `server_cert` - The server certificate.
* `intermediate_cert` - The intermediate certificate.
* `private_key` - The private key.
* `additional_certificates` - The additional certificates.
