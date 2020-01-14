---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_proxylb_acme"
subcategory: "Global"
description: |-
  Manages a SakuraCloud ProxyLB ACME Setting.
---

# sakuracloud_proxylb_acme

Manages a SakuraCloud ProxyLB ACME Setting.

## Example Usage

```hcl
resource sakuracloud_proxylb_acme "foobar" {
  proxylb_id       = sakuracloud_proxylb.foobar.id
  accept_tos       = true
  common_name      = "www.example.com"
  update_delay_sec = 120
}

data "sakuracloud_proxylb" "foobar" {
  filter {
    names = ["foobar"]
  }
}
```
## Argument Reference

* `accept_tos` - (Required) The flag to accept the current Let's Encrypt terms of service(see: https://letsencrypt.org/repository/). This must be set `true` explicitly. Changing this forces a new resource to be created.
* `common_name` - (Required) The FQDN used by ACME. This must set resolvable value. Changing this forces a new resource to be created.
* `proxylb_id` - (Required) The id of the ProxyLB that set ACME settings to. Changing this forces a new resource to be created.
* `update_delay_sec` - (Optional) The wait time in seconds. This typically used for waiting for a DNS propagation. Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 20 minutes) Used when creating the ProxyLB ACME Setting



* `delete` - (Defaults to 5 minutes) Used when deregistering ProxyLB ACME Setting



## Attribute Reference

* `id` - The id of the ProxyLB ACME Setting.
* `certificate` - A list of `certificate` blocks as defined below.


---

A `certificate` block exports the following:

* `additional_certificate` - A list of `additional_certificate` blocks as defined below.
* `intermediate_cert` - The intermediate certificate for a server.
* `private_key` - The private key for a server.
* `server_cert` - The certificate for a server.

---

A `additional_certificate` block exports the following:

* `intermediate_cert` - The intermediate certificate for a server.
* `private_key` - The private key for a server.
* `server_cert` - The certificate for a server.



