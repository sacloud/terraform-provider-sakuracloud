---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_database"
sidebar_current: "docs-sakuracloud-datasource-database"
description: |-
  Get information on a SakuraCloud Database.
---

# sakuracloud\_database

Use this data source to retrieve information about a SakuraCloud Database.

## Example Usage

```hcl
data "sakuracloud_database" "foobar" {
  name_selectors = ["foobar"]
}
```

## Argument Reference

 * `name_selectors` - (Optional) The list of names to filtering.
 * `tag_selectors` - (Optional) The list of tags to filtering.
 * `filter` - (Optional) The map of filter key and value.
 * `zone` - (Optional) The ID of the zone.

## Attributes Reference

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `plan` - The name of the resource plan.
* `user_name` - The username to access database.
* `user_password` - The password to access database.
* `allow_networks` - The network address list that allowed connections to the database.
* `port` - The number of the port on which the database is listening.
* `backup_time` - The time to perform backup.
* `backup_weekdays` - Day of the week to get backup.  
* `switch_id` - The ID of the switch connected to the database.
* `ipaddress1` - The IP address of the database.
* `nw_mask_len` - The network mask length of the database.
* `default_route` - The default route IP address of the database.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.


