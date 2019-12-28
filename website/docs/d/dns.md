---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_dns"
subcategory: "Global"
description: |-
  Get information about an existing DNS.
---

# Data Source: sakuracloud_dns

Get information about an existing DNS.

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


## Attribute Reference

* `id` - The id of the DNS.
* `description` - The description of the DNS.
* `dns_servers` - The list of IP addresses of DNS server that manage this zone.
* `icon_id` - The icon id attached to the DNS.
* `record` - A list of `record` blocks as defined below.
* `tags` - Any tags assigned to the DNS.
* `zone` - The name of managed domain.


---

A `record` block exports the following:

* `name` - The name of the DNS Record.
* `port` - The number of port.
* `priority` - The priority of target DNS Record.
* `ttl` - The number of the TTL.
* `type` - The type of DNS Record. This will be one of [`A`/`AAAA`/`ALIAS`/`CNAME`/`NS`/`MX`/`TXT`/`SRV`/`CAA`/`PTR`].
* `value` - The value of the DNS Record.
* `weight` - The weight of target DNS Record.



