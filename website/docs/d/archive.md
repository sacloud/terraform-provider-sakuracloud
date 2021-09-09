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
  os_type = "centos8"
}
```
## Argument Reference

* `filter` - (Optional) One or more values used for filtering, as defined below.
* `os_type` - (Optional) The criteria used to filter SakuraCloud archives. This must be one of following:  
  - **CentOS**: [`centos`/`centos8`/`centos8stream`/`centos7`]  
  - **Alt RHEL/CentOS**: [`almalinux`/`rockylinux`]
  - **Ubuntu**: [`ubuntu`/`ubuntu2004`/`ubuntu1804`]
  - **Debian**: [`debian`/`debian10`/`debian11`]  
  - **CoreOS/ContainerLinux**: `coreos`  
  - **RancherOS**: `rancheros`  
  - **k3OS**: `k3os`  
  - **FreeBSD**: `freebsd`  
  - **Kusanagi**: `kusanagi`  
  - **Windows2016**: [`windows2016`/`windows2016-rds`/`windows2016-rds-office`]  
  - **Windows2016+SQLServer**:  [`windows2016-sql-web`/`windows2016-sql-standard`/`windows2016-sql-standard-all`]  
  - **Windows2016+SQLServer2017**: [`windows2016-sql2017-standard`/`windows2016-sql2017-enterprise`/`windows2016-sql2017-standard-all`]  
  - **Windows2019**: [`windows2019`/`windows2019-rds`/`windows2019-rds-office2016`/`windows2019-rds-office2019`]  
  - **Windows2019+SQLServer2017**: [`windows2019-sql2017-web`/`windows2019-sql2017-standard`/`windows2019-sql2017-enterprise`/`windows2019-sql2017-standard-all`]  
  - **Windows2019+SQLServer2019**: [`windows2019-sql2019-web`/`windows2019-sql2019-standard`/`windows2019-sql2019-enterprise`/`windows2019-sql2019-standard-all`]  
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
* `values` - (Required) The values of the condition. If multiple values are specified, they combined as AND condition.


## Attribute Reference

* `id` - The id of the Archive.
* `description` - The description of the Archive.
* `icon_id` - The icon id attached to the Archive.
* `name` - The name of the Archive.
* `size` - The size of Archive in GiB.
* `tags` - Any tags assigned to the Archive.



