---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_note"
subcategory: "Misc"
description: |-
  Manages a SakuraCloud Note.
---

# sakuracloud_note

Manages a SakuraCloud Note.

## Example Usage

```hcl
resource "sakuracloud_note" "foobar" {
  name    = "foobar"
  content = file("startup-script.sh")
}
```
## Argument Reference

* `class` - (Optional) The class of the Note. This must be one of `shell`/`yaml_cloud_config`. Default:`shell`.
* `content` - (Required) The content of the Note. This must be specified as a shell script or as a cloud-config.
* `icon_id` - (Optional) The icon id to attach to the Note.
* `name` - (Required) The name of the Note. The length of this value must be in the range [`1`-`64`].
* `tags` - (Optional) Any tags to assign to the Note.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Note


* `update` - (Defaults to 5 minutes) Used when updating the Note

* `delete` - (Defaults to 5 minutes) Used when deregistering Note



## Attribute Reference

* `id` - The id of the Note.
* `description` - The description of the Note. This will be computed from special tags within body of `content`.




