---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_enhanced_db"
subcategory: "Lab"
description: |-
  Get information about an existing Enhanced Database.
---

# Data Source: sakuracloud_enhanced_db

Get information about an existing Enhanced Database.

## Example Usage

```hcl
data "sakuracloud_enhanced_db" "foobar" {
  filter {
    names = ["foobar"]
  }
}
```
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
* `operator` - (Optional) The filtering operator. This must be one of following: `partial_match_and`/`exact_match_or`. Default: `partial_match_and`


## Attribute Reference

* `id` - The id of the sakuracloud_enhanced_db.
* `allowed_networks` - A list of CIDR blocks allowed to connect.
* `database_name` - The name of database.
* `database_type` - The type of database.
* `description` - The description of the EnhancedDB.
* `hostname` - The name of database host. This will be built from `database_name` + `tidb-is1.db.sakurausercontent.com`.
* `icon_id` - The icon id attached to the EnhancedDB.
* `max_connections` - The value of max connections setting.
* `name` - The name of the EnhancedDB.
* `region` - The region name.
* `tags` - Any tags assigned to the EnhancedDB.

