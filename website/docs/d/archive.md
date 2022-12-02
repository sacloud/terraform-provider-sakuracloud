---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_archive"
subcategory: "Storage"
description: |-
  Get information about an existing Archive.
---

# Data Source: sakuracloud_archive

Get information about an existing Archive.

## Example Usage

```hcl
data "sakuracloud_archive" "foobar" {
  os_type = "ubuntu"
}
```
## Argument Reference

* `filter` - (Optional) One or more values used for filtering, as defined below.
* `os_type` - (Optional) The criteria used to filter SakuraCloud archives. This must be one of following:  
  - **CentOS**: [`centos`/`centos7`]  
  - **AlmaLinux**: [`almalinux`/`almalinux8`/`almalinux9`]
  - **RockyLinux**: [`rockylinux`/`rockylinux8`/`rockylinux9`]
  - **MIRACLE LINUX**[`miracle`/`miraclelinux`/`miracle8`/`miraclelinux8`/`miracle9`/`miraclelinux9`]
  - **Ubuntu**: [`ubuntu`/`ubuntu2204`/`ubuntu2004`/`ubuntu1804`]
  - **Debian**: [`debian`/`debian10`/`debian11`]  
  - **Kusanagi**: `kusanagi`  
* `zone` - (Optional) The name of zone that the Archive is in (e.g. `is1a`, `tk1a`).

---

A `filter` block supports the following:

* `condition` - (Optional) One or more name/values pairs used for filtering. There are several valid keys, for a full reference, check out finding section in the [SakuraCloud API reference](https://developer.sakura.ad.jp/cloud/api/1.1/).
* `id` - (Optional) The resource id on SakuraCloud used for filtering.
* `names` - (Optional) The resource names on SakuraCloud used for filtering. If multiple values are specified, they combined as AND condition.
* `tags` - (Optional) The resource tags on SakuraCloud used for filtering. If multiple values are specified, they combined as AND condition.

---

A `condition` block supports the following:

* `name` - (Required) The name of the target field. This value is case-sensitive.
* `values` - (Required) The values of the condition. If multiple values ​​are specified, they combined as AND condition.
* `operator` - (Optional) The filtering operator. This must be one of following: `partial_match_and`/`exact_match_or`. Default: `partial_match_and`


## Attribute Reference

* `id` - The id of the Archive.
* `description` - The description of the Archive.
* `icon_id` - The icon id attached to the Archive.
* `name` - The name of the Archive.
* `size` - The size of Archive in GiB.
* `tags` - Any tags assigned to the Archive.



