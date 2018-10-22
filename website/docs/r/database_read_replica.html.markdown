---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_database"
sidebar_current: "docs-sakuracloud-resource-appliance-read-replica"
description: |-
  Provides a SakuraCloud Database(ReadReplica) resource. This can be used to create, update, and delete Databases.
---

# sakuracloud\_database\_read\_replica

Provides a SakuraCloud Database(ReadReplica) resource. This can be used to create, update, and delete Databases.

## Example Usage

```hcl
# Create a new Database(ReadReplica)
resource "sakuracloud_database_read_replica" "foobar" {
  name          = "slave"
  master_id     = sakuracloud_database.foobar.id

  ipaddress1    = "192.168.11.111"
  #switch_id     = sakuracloud_switch.sw.id
  #nw_mask_len   = 24
  #default_route = "192.168.11.1"
}

# Create a new Database(Master)
resource "sakuracloud_database" "foobar" {
  name             = "foobar"
  database_type    = "mariadb"
  plan             = "30g"
  user_name        = "user"
  user_password    = "p@ssw0rd"
  replica_password = "p@ssw0rd"

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

* `master_id` - (Required) The ID of the master Database Appliance.
* `name` - (Required) The name of the resource.
* `ipaddress1` - (Required) The IP address of the database.
* `switch_id` - (Optional) The ID of the switch connected to the database.
* `nw_mask_len` - (Optional) The network mask length of the database.
* `default_route` - (Optional) The default route IP address of the database.
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.
* `graceful_shutdown_timeout` - (Optional) The wait time (seconds) to do graceful shutdown the Database.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `master_id` - The ID of the master Database Appliance.
* `name` - The name of the resource.
* `ipaddress1` - The IP address of the database.
* `switch_id` - The ID of the switch connected to the database.
* `nw_mask_len` - The network mask length of the database.
* `default_route` - The default route IP address of the database.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.

## Import

Databases can be imported using the Database ID.

```
$ terraform import sakuracloud_database_read_replica.foobar <database_id>
```
