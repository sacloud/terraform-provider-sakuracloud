---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_proxylb_acme"
sidebar_current: "docs-sakuracloud-resource-global-proxylb-acme"
description: |-
  Provides a SakuraCloud ProxyLB resource. This can be used to create, update, and delete ProxyLBs.
---

# sakuracloud\_proxylb

Provides a SakuraCloud ProxyLB(Enhanced-LoadBalancer) resource. This can be used to create, update, and delete ProxyLBs.

## Example Usage

```hcl
resource "sakuracloud_proxylb_acme" "cert" {
  proxylb_id = sakuracloud_proxylb.foobar.id
  accept_tos = true
  common_name = "foobar.example.com"
  update_delay_sec = 120
}

resource "sakuracloud_proxylb" "foobar" {
  name         = "foobar"
  plan         = 1000 
  vip_failover = true # default: false

  bind_ports {
    proxy_mode        = "http"
    port              = 80
  }
  bind_ports {
    proxy_mode    = "https"
    port          = 443
  }
  
  servers {
    ipaddress = "133.242.0.3"
    port = 80
  }
  servers {
    ipaddress = "133.242.0.4"
    port = 80
  }
}
```

## Argument Reference

The following arguments are supported:

* `proxylb_id` - (Required) The ID of target ProxyLB resource.  
* `accept_tos` - (Required) The flag for accept Let's Encrypt's [Terms of Service](https://letsencrypt.org/repository/).  
* `common_name` - (Require) The FQDN of target domain.  
* `update_delay_sec` - (Optional) The wait time for update settings.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `proxylb_id` - The ID of target ProxyLB resource.  
* `common_name` - The FQDN of target domain.  
* `certificate` - Certificate used to terminate SSL/TSL. It contains some attributes to [Certificate](#certificate).

### Certificate

* `server_cert` - The server certificate.
* `intermediate_cert` - The intermediate certificate.
* `private_key` - The private key.
* `additional_certificates` - Additional certificates.

## Import

ProxyLB ACME can be imported using the ProxyLB ID.

```
$ terraform import sakuracloud_proxylb_acme.foobar <proxylb_id>
```
