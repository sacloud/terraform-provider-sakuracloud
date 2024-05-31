---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_webaccel_acl"
subcategory: "WebAccelerator"
description: |-
  Manages a SakuraCloud WebAccelerator Site ACL.
---

# sakuracloud_webaccel_acl

Manages a SakuraCloud sakuracloud_webaccel_acl.

## Example Usage

```hcl
data sakuracloud_webaccel "site" {
  name = "your-site-name"
  # or
  # domain = "your-domain"
}

resource sakuracloud_webaccel_acl "acl" {
  site_id = data.sakuracloud_webaccel.site.id

  acl = join("\n", [
    "deny 192.0.2.5/25",
    "deny 198.51.100.0",
    "allow all",
  ])
}
```

## Argument Reference

* `acl` - (Required) .
* `site_id` - (Required) .



## Attribute Reference

* `id` - The id of the sakuracloud_webaccel_acl.



