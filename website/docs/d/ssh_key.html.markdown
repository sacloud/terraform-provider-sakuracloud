---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_ssh_key"
sidebar_current: "docs-sakuracloud-datasource-ssh-key"
description: |-
  Get information on a SakuraCloud SSH Key.
---

# sakuracloud\_ssh_key

Use this data source to retrieve information about a SakuraCloud SSH Key.

## Example Usage

```hcl
data "sakuracloud_ssh_key" "foobar" {
  name_selectors = ["foobar"]
}
```

## Argument Reference

 * `name_selectors` - (Optional) The list of names to filtering.
 * `filter` - (Optional) The map of filter key and value.

## Attributes Reference

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `description` - The description of the resource.
* `public_key` - The body of the public key. 
* `finger_print` - The fingerprint of the public key.
