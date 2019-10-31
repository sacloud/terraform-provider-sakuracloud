---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_switch"
sidebar_current: "docs-sakuracloud-resource-webaccel-certificate"
description: |-
  Provides a SakuraCloud WebAccel Certificate resource. This can be used to create, update, and delete Certificate.
---

# sakuracloud\_webaccel\_certificate

Provides a SakuraCloud WebAccel Certificate resource. This can be used to create, update, and delete Certificate.

## Example Usage

```hcl
data sakuracloud_webaccel "site" {
  name = "example"
}

resource sakuracloud_webaccel_certificate "example" {
  site_id           = data.sakuracloud_webaccel.site.id
  certificate_chain = file("crt")
  private_key               = file("key")
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) The id of the target site on WebAccel.
* `certificate_chain` - (Required) The contents of certificate.
* `private_key` - (Required) The content of private key.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `serial_number` - .
* `not_before` - .
* `not_after` - .
* `issuer_common_name` - .
* `subject_common_name` - .
* `dns_names` - .
* `sha256_fingerprint` - .

## Import

Importing the WebAccel Certificate is not supported.
