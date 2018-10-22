---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_database"
sidebar_current: "docs-sakuracloud-resource-appliance-database"
description: |-
  Provides a SakuraCloud Database resource. This can be used to create, update, and delete Databases.
---

# sakuracloud\_database

Provides a SakuraCloud Database resource. This can be used to create, update, and delete Databases.

## Example Usage

```hcl
# Create a new Database
resource "sakuracloud_database" "foobar" {
  name            = "foobar"
  database_type   = "mariadb"
  plan            = "30g"
  user_name       = "user"
  user_password   = "p@ssword"
  allow_networks  = ["192.168.2.0/24"]
  port            = 33061
  backup_time     = "00:00"
  backup_weekdays = ["mon", "tue", "wed"]

  switch_id     = sakuracloud_switch.foobar.id
  ipaddress1    = "192.168.11.101"
  nw_mask_len   = 24
  default_route = "192.168.11.1"

  description = "description"
  tags        = ["foo", "bar"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `database_type` - (Optional) The Database type.  
Valid value is one of the following: [ "postgresql" (default) / "mariadb"]
* `plan` - (Optional) The plan (size) of the Database.   
Valid value is one of the following: [ "10g" (default) / "30g" / "90g" / "240g" / "500g" / "1t" ]
* `user_name` - (Required) The username to access database.
* `user_password` - (Required) The password to access database.
* `allow_networks` - (Optional) The network address list that allowed connections to the database.
* `port` - (Optional) The number of the port on which the database is listening (default:`5432`).
* `backup_time` - (Optional) The time to perform backup (format:`HH:mm`).
* `backup_weekdays` - (Optional) Day of the week to get backup.  
Valid values are the following: ["mon", "tue", "wed", "thu", "fri", "sat", "sun"]
* `switch_id` - (Required) The ID of the switch connected to the database.
* `ipaddress1` - (Required) The IP address of the database.
* `nw_mask_len` - (Required) The network mask length of the database.
* `default_route` - (Required) The default route IP address of the database.
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.
* `graceful_shutdown_timeout` - (Optional) The wait time (seconds) to do graceful shutdown the Database.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `plan` - The name of the resource plan.
* `user_name` - The username to access database.
* `user_password` - The password to access database.
* `allow_networks` - The network address list that allowed connections to the database.
* `port` - The number of the port on which the database is listening.
* `backup_time` - The time to perform backup.
* `switch_id` - The ID of the switch connected to the database.
* `ipaddress1` - The IP address of the database.
* `nw_mask_len` - The network mask length of the database.
* `default_route` - The default route IP address of the database.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.

## Import

Databases can be imported using the Database ID.

```
$ terraform import sakuracloud_database.foobar <database_id>
```
