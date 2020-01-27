---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_zone"
subcategory: "Data Sources"
description: |-
  Get information on a SakuraCloud WebAccelerator Site.
---

# sakuracloud\_webaccel

Use this data source to retrieve information about a SakuraCloud WebAccelerator Site.

## Example Usage

```hcl
data sakuracloud_webaccel "example" {
  domain = "www.example.com"
  # or 
  # name = "example"
}
```

## Argument Reference

 * `name` - (Optional) The name of site.
 * `domain` - (Optional) The name of domain.

## Attributes Reference

* `id` - The ID of the resource.
* `site_id` - The site ID.
* `origin` - .
* `subdomain` - .
* `domain_type` - .
* `has_certificate` - .
* `host_header` - .
* `status` - .
* `cname_record_value` - .
* `txt_record_value` - .
