---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_webaccel_activation"
subcategory: "WebAccelerator"
description: |-
  Manages a SakuraCloud WebAccelerator Site Status.
---

# sakuracloud_webaccel_activation

Manages a SakuraCloud Web Accelerator Site Status.

## Example Usage

```hcl
data "sakuracloud_webaccel" "site" {
  name = "your-site-name"
  # or
  # domain = "your-domain"
}

resource "sakuracloud_webaccel_activation" "site_status" {
  site_id    = data.sakuracloud_webaccel.site.id
  enabled    = true
}
```

## Argument Reference

* `site_id` - (Required) The target webaccel site.
* `enabled` - (Required) Whether the site is activated or not.

## Attribute Reference

* `site_id` - The target webaccel site.
* `enabled` - Whether the site is activated or not.
