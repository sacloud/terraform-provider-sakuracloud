---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_database"
subcategory: "Appliance"
description: |-
  Manages a SakuraCloud Database.
---

# sakuracloud_database

Manages a SakuraCloud Database.

## Example Usage

```hcl
variable username {}
variable password {}
variable replica_password {}

resource "sakuracloud_database" "foobar" {
  database_type    = "mariadb"
  database_version = "10.11" // optional
  plan             = "30g"

  username = var.username
  password = var.password

  replica_password = var.replica_password

  network_interface {
    switch_id     = sakuracloud_switch.foobar.id
    ip_address    = "192.168.11.11"
    netmask       = 24
    gateway       = "192.168.11.1"
    port          = 3306
    source_ranges = ["192.168.11.0/24", "192.168.12.0/24"]
  }

  backup {
    time     = "00:00"
    weekdays = ["mon", "tue"]
  }

  # continuous_backupを指定するときはdatabase_versionが必須
  # continuous_backup {
  #   days_of_week = ["mon", "tue"]
  #   time         = "01:30"
  #   connect      = "nfs://${sakuracloud_nfs.foobar.network_interface[0].ip_address}/export"
  # }

  parameters = {
    max_connections = 100
  }

  monitoring_suite {
    enabled = true
  }

  disk {
    encryption_algorithm = "aes256_xts"
    kms_key_id           = sakuracloud_kms.foobar.id
  }

  name        = "foobar"
  description = "description"
  tags        = ["tag1", "tag2"]
}

resource "sakuracloud_nfs" "foobar" {
  name = "foobar"
  plan = "ssd"
  size = "100"

  network_interface {
    switch_id   = sakuracloud_switch.foobar.id
    ip_address  = "192.168.11.111"
    netmask     = 24
    gateway     = "192.168.11.1"
  }
}

resource "sakuracloud_switch" "foobar" {
  name = "foobar"
}

resource "sakuracloud_kms" "foobar" {
  name = "foobar"
}
```
## Argument Reference

#### Plan

* `database_type` - (Optional) The type of the database. This must be one of [`mariadb`/`postgres`]. Changing this forces a new resource to be created. Default:`postgres`.
* `database_version` - (Optional) The version of the database.  Changing this forces a new resource to be created.
* `plan` - (Optional) The plan name of the Database. This must be one of [`10g`/`30g`/`90g`/`240g`/`500g`/`1t`]. Changing this forces a new resource to be created. Default:`10g`.

#### User

* `username` - (Required) The name of default user on the database. The length of this value must be in the range [`3`-`20`]. Changing this forces a new resource to be created.
* `password` - (Required) The password of default user on the database.

#### Disk

* `disk` - (Optional) A `disk` block as defined below. Changing this forces a new resource to be created.

---

A `disk` block supports the following:

* `encryption_algorithm` - (Optional) The disk encryption algorithm. This must be one of [`none`/`aes256_xts`].
* `kms_key_id` - (Optional) ID of the KMS key for encryption.


#### Network

* `network_interface` - (Required) An `network_interface` block as defined below.

---

A `network_interface` block supports the following:

* `gateway` - (Required) The IP address of the gateway used by Database.
* `ip_address` - (Required) The IP address to assign to the Database.
* `netmask` - (Required) The bit length of the subnet to assign to the Database. This must be in the range [`8`-`29`].
* `port` - (Optional) The number of the listening port. This must be in the range [`1024`-`65535`].
* `source_ranges` - (Optional) The range of source IP addresses that allow to access to the Database via network.
* `switch_id` - (Required) The id of the switch to which the Database connects.

#### Backup

* `backup` - (Optional) A `backup` block as defined below.
* `continuous_backup` - (Optional) A `continuous_backup` block as defined below. This field can only be specified when `database_version` is provided.
 
---

A `backup` block supports the following:

* `time` - (Optional) The time to take backup. This must be formatted with `HH:mm`.
* `weekdays` - (Optional) A list of weekdays to backed up. The values in the list must be in [`sun`/`mon`/`tue`/`wed`/`thu`/`fri`/`sat`].
 
---

A `continuous_backup` block supports the following:

* `connect` - (Required) NFS server address for storing backups (e.g., `nfs://192.0.2.1/export`).
* `time` - (Required) The time to take backup. This must be formatted with `HH:mm`.
* `days_of_week` - (Required) A list of weekdays to backed up. The values in the list must be in [`sun`/`mon`/`tue`/`wed`/`thu`/`fri`/`sat`].


#### RDBMS Parameters

* `parameters` - (Optional) The map for setting RDBMS-specific parameters. Valid keys can be found with the `usacloud database list-parameters` command.

#### Replication

* `replica_password` - (Optional) The password of user that processing a replication.
* `replica_user` - (Optional) The name of user that processing a replication. Default:`replica`.

#### Monitoring

* `monitoring_suite` - (Optional) An `monitoring_suite` block as defined below.
---

A `monitoring_suite` block supports the following:

* `enabled` - (Optional) Enable sending signals to Monitoring Suite.


#### Common Arguments

* `name` - (Required) The name of the Database. The length of this value must be in the range [`1`-`64`].
* `tags` - (Optional) Any tags to assign to the Database.
* `icon_id` - (Optional) The icon id to attach to the Database.
* `description` - (Optional) The description of the Database. The length of this value must be in the range [`1`-`512`].
* `zone` - (Optional) The name of zone that the Database will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Database
* `update` - (Defaults to 60 minutes) Used when updating the Database
* `delete` - (Defaults to 20 minutes) Used when deleting Database

## Attribute Reference

* `id` - The id of the Database.

