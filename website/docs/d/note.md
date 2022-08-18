---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_note"
subcategory: "Misc"
description: |-
  Get information about an existing Note.
---

# Data Source: sakuracloud_note

Get information about an existing Note.

## Example Usage

```hcl
data "sakuracloud_note" "foobar" {
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

* `id` - The id of the Note.
* `class` - The class of the Note. This will be one of [`shell`/`yaml_cloud_config`].
* `content` - The content of the Note.
* `description` - The description of the Note.
* `icon_id` - The icon id attached to the Note.
* `name` - The name of the Note.
* `tags` - Any tags assigned to the Note.



