---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_archive"
sidebar_current: "docs-sakuracloud-datasource-archive"
description: |-
  Get information on a SakuraCloud Archive.
---

# sakuracloud\_archive

Use this data source to retrieve information about a SakuraCloud Archive.

## Example Usage

### Using `os_type` parameter

```hcl
data "sakuracloud_archive" "centos" {
  os_type = "centos"
}
```

### Using filter parameters

```hcl
data "sakuracloud_archive" "ubuntu" {
  name_selectors = ["Ubuntu", "LTS"]
  tag_selectors  = ["current-stable", "os-linux"]
}
```

## Argument Reference

 * `os_type` - (Optional) The slug of target public archive. Valid values are in [`os_type` section](#os_type-parameter-reference).
 * `name_selectors` - (Optional) The list of names to filtering.
 * `tag_selectors` - (Optional) The list of tags to filtering.
 * `filter` - (Optional) The map of filter key and value.
 * `zone` - (Optional) The ID of the zone.

## Attributes Reference

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `size` - The size of the resource (unit:`GB`).
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.

## `os_type` Parameter Reference

* `centos` - CentOS 7
* `centos6` - CentOS 6
* `ubuntu` - Ubuntu 
* `debian` - Debian 
* `vyos` - VyOS
* `coreos` - CoreOS
* `rancheros` - RancherOS
* `kusanagi` - Kusanagi (CentOS7)
* `sophos-utm` - Sophos-UTM
* `freebsd` - FreeBSD
* `windows2012` - Windows 2012
* `windows2012-rds` - Windows 2012 (RDS)
* `windows2012-rds-office` - Windows 2012 (RDS + Office)
* `windows2016` - Windows 2016
* `windows2016-rds` - Windows 2016 (RDS)
* `windows2016-rds-office` - Windows 2016 (RDS + Office)
* `windows2016-sql-web` - Windows 2016 SQLServer (Web)
* `windows2016-sql-standard` - Windows 2016 SQLServer (Standard)
* `windows2016-sql2017-standard` - Windows 2016 SQLServer 2017 (Standard)
* `windows2016-sql-standard-all` - Windows 2016 SQLServer (Standard, RDS + Office)
* `windows2016-sql2017-standard-all` - Windows 2016 SQLServer 2017 (Standard, RDS + Office)
