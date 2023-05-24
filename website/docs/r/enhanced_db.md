---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_enhanced_db"
subcategory: "Lab"
description: |-
  Manages a SakuraCloud sakuracloud_enhanced_db.
---

# sakuracloud_enhanced_db

Manages a SakuraCloud sakuracloud_enhanced_db.

## Example Usage

```hcl
resource "sakuracloud_enhanced_db" "foobar" {
  name            = "example"
  database_name   = "example"
  password        = "your-password"

  description = "..."
  tags        = ["...", "..."]
}
```

## Argument Reference

* `allowed_networks` - (Optional) A list of CIDR blocks allowed to connect.
* `database_name` - (Required) The name of database. Changing this forces a new resource to be created.
* `database_type` - (Required) The type of database. This must be one of [`tidb`/`mariadb`]. Changing this forces a new resource to be created.
* `description` - (Optional) The description of the Enhanced Database. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the Enhanced Database.
* `name` - (Required) The name of the Enhanced Database. The length of this value must be in the range [`1`-`64`].
* `password` - (Required) The password of database.
* `region` - (Required) The name of region that the database is in. This must be one of [`is1`/`tk1`]. Changing this forces a new resource to be created.
* `tags` - (Optional) Any tags to assign to the Enhanced Database.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the sakuracloud_enhanced_db
* `update` - (Defaults to 5 minutes) Used when updating the sakuracloud_enhanced_db
* `delete` - (Defaults to 5 minutes) Used when deleting sakuracloud_enhanced_db


## Attribute Reference

* `id` - The id of the sakuracloud_enhanced_db.
* `hostname` - The name of database host. This will be built from `database_name` + `tidb-is1.db.sakurausercontent.com`.
* `max_connections` - The value of max connections setting.



