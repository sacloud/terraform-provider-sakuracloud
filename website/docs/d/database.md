---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_database"
subcategory: "Appliance"
description: |-
  Get information about an existing Database.
---

# Data Source: sakuracloud_database

Get information about an existing Database.

## Example Usage

```hcl
data "sakuracloud_database" "foobar" {
  filter {
    names = ["foobar"]
  }
}
```
## Argument Reference

* `filter` - (Optional) One or more values used for filtering, as defined below.
* `zone` - (Optional) The name of zone that the Database is in (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.


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
* `operator` - (Optional) The filtering operator. This must be one of following: `partial_match_and`/`exact_match_or`. Default: `partial_match_and`


## Attribute Reference

* `id` - The id of the Database.
* `backup` - A list of `backup` blocks as defined below.
* `continuous_backup` - A list of `continuous_backup` blocks as defined below.
* `database_type` - The type of the database. This will be one of [`mariadb`/`postgres`].
* `database_version` - The version of the database.
* `description` - The description of the Database.
* `disk` - A list of `disk` blocks as defined below.
* `icon_id` - The icon id attached to the Database.
* `monitoring_suite` - A list of `monitoring_suite` blocks as defined below.
* `name` - The name of the Database.
* `network_interface` - A list of `network_interface` blocks as defined below.
* `parameters` - The map for setting RDBMS-specific parameters. Valid keys can be found with the `usacloud database list-parameters` command.
* `password` - The password of default user on the database.
* `plan` - The plan name of the Database. This will be one of [`10g`/`30g`/`90g`/`240g`/`500g`/`1t`].
* `replica_password` - The password of user that processing a replication.
* `replica_user` - The name of user that processing a replication.
* `tags` - Any tags assigned to the Database.
* `username` - The name of default user on the database.

---

A `backup` block exports the following:

* `time` - The time to take backup. This will be formatted with `HH:mm`.
* `weekdays` - The list of name of weekday that doing backup. This will be in [`sun`/`mon`/`tue`/`wed`/`thu`/`fri`/`sat`].

---

A `continuous_backup` block exports the following:

* `connect` - NFS server address for storing backups (e.g., `nfs://192.0.2.1/export`).
* `time` - The time to take backup. This must be formatted with `HH:mm`.
* `weekdays` - The list of name of weekday that doing backup. This will be in [`sun`/`mon`/`tue`/`wed`/`thu`/`fri`/`sat`].

---

A `disk` block exports the following:

* `encryption_algorithm` - The disk encryption algorithm. This must be one of [`none`/`aes256_xts`].
* `kms_key_id` - ID of the KMS key for encryption.

---

A `monitoring_suite` block exports the following:

* `enabled` - Enable sending signals to Monitoring Suite.

---

A `network_interface` block exports the following:

* `gateway` - The IP address of the gateway used by Database.
* `ip_address` - The IP address assigned to the Database.
* `netmask` - The bit length of the subnet assigned to the Database.
* `port` - The number of the listening port.
* `source_ranges` - The range of source IP addresses that allow to access to the Database via network.
* `switch_id` - The id of the switch connected from the Database.


