---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_database_read_replica"
subcategory: "Appliance"
description: |-
  Manages a SakuraCloud Database Read Replica.
---

# sakuracloud_database_read_replica

Manages a SakuraCloud Database Read Replica.

## Example Usage

```hcl
resource "sakuracloud_kms" "foobar" {
  name = "foobar"
}

resource "sakuracloud_database_read_replica" "foobar" {
  master_id   = data.sakuracloud_database.master.id

  network_interface {
    ip_address  = "192.168.11.111"
  }

  disk {
    encryption_algorithm = "aes256_xts"
    kms_key_id           = sakuracloud_kms.foobar.id
  }

  name        = "foobar"
  description = "description"
  tags        = ["tag1", "tag2"]
}

data sakuracloud_database "master" {
  filter {
    names = ["master-database-name"]
  }
}

```

## Argument Reference

* `name` - (Required) The name of the read-replica database. The length of this value must be in the range [`1`-`64`].
* `master_id` - (Required) The id of the replication master database. Changing this forces a new resource to be created.

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

* `ip_address` - (Required) The IP address to assign to the read-replica database.
* `gateway` - (Optional) The IP address of the gateway used by read-replica database. If `gateway` isn't specified, it will be set to the same value of the master database.
* `netmask` - (Optional) The bit length of the subnet to assign to the read-replica database. This must be in the range [`8`-`29`]. If `netmask` isn't specified, it will be set to the same value of the master database.
* `source_ranges` - (Optional) The range of source IP addresses that allow to access to the read-replica database via network.
* `switch_id` - (Optional) The id of the switch to which the read-replica database connects. If `switch_id` isn't specified, it will be set to the same value of the master database.

#### Common Arguments

* `description` - (Optional) The description of the read-replica database. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the read-replica database.
* `tags` - (Optional) Any tags to assign to the read-replica database.
* `zone` - (Optional) The name of zone that the read-replica database will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Database Read Replica
* `update` - (Defaults to 60 minutes) Used when updating the Database Read Replica
* `delete` - (Defaults to 20 minutes) Used when deleting Database Read Replica

## Attribute Reference

* `id` - The id of the Database Read Replica.

