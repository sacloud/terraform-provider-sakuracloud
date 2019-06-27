---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_proxylb"
sidebar_current: "docs-sakuracloud-resource-global-proxylb"
description: |-
  Provides a SakuraCloud ProxyLB resource. This can be used to create, update, and delete ProxyLBs.
---

# sakuracloud\_proxylb

Provides a SakuraCloud ProxyLB(Enhanced-LoadBalancer) resource. This can be used to create, update, and delete ProxyLBs.

## Example Usage

```hcl
resource "sakuracloud_proxylb" "foobar" {
  name         = "foobar"
  plan         = 1000 
  vip_failover = false

  health_check {
    protocol    = "http"
    path        = "/"
    host_header = "example.com"
    delay_loop  = 10
  }

  bind_ports {
    proxy_mode    = "https"
    port          = 443
    support_http2 = true
  }

  sorry_server {
    ipaddress         = "192.2.0.1"
    port              = 80
    redirect_to_https = true
  }

  servers {
    ipaddress = "133.242.0.3"
    port = 80
  }
  servers {
    ipaddress = "133.242.0.4"
    port = 80
  }

  certificate {
    server_cert = file("server.crt")
    private_key = file("server.key")    
    # intermediate_cert = file("intermediate.crt")
    
    additional_certificates {
      server_cert = file("server2.crt")
      private_key = file("server2.key")    
      # intermediate_cert = file("intermediate2.crt")
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.  
* `plan` - (Optional) The plan of the resource.
Valid value is one of the following: [ 1000 (default) / 5000 / 10000 / 50000 / 100000 ]  
* `vip_failover` - (Optional) The flag of enable VIP Fail-Over.  
* `sticky_session` - (Optional) The flag of enable Sticky-Session.  
* `bind_ports` - (Required) The external listen ports. It contains some attributes to [Bind Ports](#bind-ports).
* `health_check` - (Required) The health check rules. It contains some attributes to [Health Check](#health-check).
* `sorry_server` - (Optional) The pair of IPAddress and port number of sorry-server.
* `servers` - (Optional) Real servers. It contains some attributes to [Servers](#servers).
* `certificate` - (Optional) Certificate used to terminate SSL/TSL. It contains some attributes to [Certificate](#certificate).
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.

### Bind Ports

Attributes for Bind-Ports:

* `proxy_mode` - (Required) Proxy protocol.  
Valid value is one of the following: [ "http" / "https"]
* `port` - (Required) Port number used in tcp proxy.
* `redirect_to_https` - (Optional) The flag for enable to redirect to https.
* `support_http2` - (Optional) The flag for enable to support HTTP/2.


### Health Check

Attributes for Health Check:

* `protocol` - (Required) Protocol used in health check.  
Valid value is one of the following: [ "http" / "tcp" ]
* `delay_loop` - (Optional) Health check access interval (unit:`second`, default:`10`).
* `host_header` - (Optional) The value of `Host` header used in http/https health check access.
* `path` - (Optional) The request path used in http health check access.

### Servers

* `ipaddress` - (Required) The IP address of the Real-Server.
* `port` - (Required) Port number.
* `enabled` - (Optional) The flag for enable/disable the Real-Server (default:`true`).

### Certificate

* `server_cert` - (Required) The server certificate.
* `intermediate_cert` - (Optional) The intermediate certificate.
* `private_key` - (Optional) The private key.
* `additional_certificates` - (Optional) Additional certificates.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - Name of the resource.
* `plan` - The plan of the resource.
* `vip_failover` - The flag of enable VIP Fail-Over.  
* `bind_ports` - The external listen ports. It contains some attributes to [Bind Ports](#bind-ports).
* `health_check` - The health check rules. It contains some attributes to [Health Check](#health-check).
* `sorry_server` - The pair of IPAddress and port number of sorry-server.
* `servers` - Real servers. It contains some attributes to [Servers](#servers).
* `certificate` - Certificate used to terminate SSL/TSL. It contains some attributes to [Certificate](#certificate).
* `vip` - The VirtualIPAddress that was assigned. This attribute only valid when `vip_failover` is set to `false`.  
* `fqdn` - The FQDN that was assigned. This attribute only valid when `vip_failover` is set to `true`.  
* `proxy_networks` - ProxyLB network address.   
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon.

## Import

ProxyLBs can be imported using the ProxyLB ID.

```
$ terraform import sakuracloud_proxylb.foobar <proxylb_id>
```
