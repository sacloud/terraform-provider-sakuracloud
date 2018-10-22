---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_dns"
sidebar_current: "docs-sakuracloud-datasource-dns"
description: |-
  Get information on a SakuraCloud DNS.
---

# sakuracloud\_dns

Use this data source to retrieve information about a SakuraCloud DNS.

## Example Usage

```hcl
data "sakuracloud_dns" "foobar" {
  name_selectors = ["example.com"]
}
```

## Argument Reference

 * `name_selectors` - (Optional) The list of names to filtering.
 * `tag_selectors` - (Optional) The list of tags to filtering.
 * `filter` - (Optional) The map of filter key and value.

## Attributes Reference

* `id` - The ID of the resource.
* `zone` - The name of the zone.
* `dns_servers` - List of host names of DNS servers.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
