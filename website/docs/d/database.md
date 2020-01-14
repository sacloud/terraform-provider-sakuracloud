---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_database"
subcategory: "Appliance"
description: |-
  Get information about an existing Database.
---

# Data Source: sakuracloud_database

Get information about an existing Database.

## Argument Reference

* `filter` - (Optional) One or more values used for filtering, as defined below.


---

A `filter` block supports the following:

* `condition` - (Optional) One or more name/values pairs used for filtering. There are several valid keys, for a full reference, check out finding section in the [SakuraCloud API reference](https://developer.sakura.ad.jp/cloud/api/1.1/).
* `id` - (Optional) The resource id on SakuraCloud used for filtering.
* `names` - (Optional) The resource names on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition.
* `tags` - (Optional) The resource tags on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition.

---

A `condition` block supports the following:

* `name` - (Required) The name of the target field. This value is case-sensitive.
* `values` - (Required) The values of the condition. If multiple values ​​are specified, they combined as AND condition.


## Attribute Reference

* `id` - The id of the Database.
* `backup_time` - The time to take backup. This will be formatted with `HH:mm`.
* `backup_weekdays` - The list of name of weekday that doing backup. This will be in [`sun`/`mon`/`tue`/`wed`/`thu`/`fri`/`sat`].
* `database_type` - The type of the database. This will be one of [`mariadb`/`postgres`].
* `description` - The description of the database.
* `gateway` - The IP address of the gateway used by database.
* `icon_id` - The icon id attached to the database.
* `ip_addresses` - The list of IP address assigned to the database.
* `name` - The name of the database.
* `netmask` - The bit length of the subnet assigned to the database.
* `password` - The password of default user on the database.
* `plan` - The plan name of the database. This will be one of [`10g`/`30g`/`90g`/`240g`/`500g`/`1t`].
* `port` - The number of the listening port.
* `replica_password` - The password of user that processing a replication.
* `replica_user` - The name of user that processing a replication.
* `source_ranges` - The range of source IP addresses that allow to access to the database via network.
* `switch_id` - The id of the switch connected from the database.
* `tags` - Any tags assigned to the database.
* `username` - The name of default user on the database.
* `zone` - The name of zone that the database is in (e.g. `is1a`, `tk1a`).




