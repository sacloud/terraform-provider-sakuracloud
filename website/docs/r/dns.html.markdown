---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_dns"
sidebar_current: "docs-sakuracloud-resource-global-dns-zone"
description: |-
  Provides a SakuraCloud DNS zones resource. This can be used to create, update, and delete DNS zones.
---

# sakuracloud\_dns

Provides a SakuraCloud DNS zones resource. This can be used to create, update, and delete DNS zones.

## Example Usage

```hcl
# Create a new DNS zone
resource "sakuracloud_dns" "foobar" {
  zone        = "example.com"
  description = "description"
  tags        = ["foo", "bar"]
}
```

## Argument Reference

The following arguments are supported:

* `zone` - (Required) The name of the zone.
* `records` - (Optional) Records. It contains some attributes to [Records](#records).
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.

### Records

Attributes for Records:

* `name` - (Required) The hostname of target Record. If "@" is specified, it indicates own zone.
* `type` - (Required) The Record type.  
Valid value is one of the following: [ "A" / "AAAA" / "ALIAS" / "CNAME" / "NS" / "MX" / "TXT" / "SRV" / "CAA" ]
* `value` - (Required) The value of the Record. 
* `ttl` - (Optional) The ttl value of the Record (unit:`second`). 
* `priority` - (Optional) The priority used when `type` is `MX` or `SRV`.
* `weight` - (Optional) The weight used when `type` is `SRV`.
* `port` - (Optional) The port number used when `type` is `SRV`. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `zone` - The name of the zone.
* `dns_servers` - List of host names of DNS servers.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.

## Import

DNS can be imported using the DNS ID.

```
$ terraform import sakuracloud_dns.foobar <dns_id>
```
