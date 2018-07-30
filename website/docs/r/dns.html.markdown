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
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.

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
