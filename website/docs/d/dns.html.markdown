---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_dns"
sidebar_current: "docs-sakuracloud-datasource-dns"
description: |-
  Get information on a SakuraCloud dns.
---

# sakuracloud\_dns

Use this data source to retrieve information about a SakuraCloud dns.

## Example Usage

```hcl
data sakuracloud_dns "foobar" {
  name_selectors = ["example.com"]
}
```

## Argument Reference

 * `name_selectors` - (Optional) The filter value list of name.
 * `tag_selectors` - (Optional) The filter value list of tags.
 * `filter` - (Optional) The map of filter key and value.

## Attributes Reference

* `id` - The ID of the resource.
* `zone` - The name of the zone.
* `dns_servers` - List of host names of DNS servers.
* `description` - The description of the resource.
* `tags` - The tag list of the resource.
* `icon_id` - The ID of the icon of the resource.
