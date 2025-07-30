---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_webaccel_certificate"
subcategory: "WebAccelerator"
description: |-
  Manages a SakuraCloud WebAccelerator Certificate.
---

# sakuracloud_webaccel_certificate

Manages a SakuraCloud sakuracloud_webaccel_certificate.

## Example Usage

```hcl
data sakuracloud_webaccel "site" {
  name = "your-site-name"
  # or
  # domain = "your-domain"
}

resource sakuracloud_webaccel_certificate "foobar" {
  site_id           = data.sakuracloud_webaccel.site.id
  certificate_chain = file("path/to/your/certificate/chain")
  private_key       = file("path/to/your/private/key")
}
```

## Argument Reference

* `certificate_chain` - (Optional) Certificate chain for the site (mutually exclusive with `lets_encrypt`, used with `private_key`).
* `private_key` - (Optional) Private key for the site (mutually exclusive with `lets_encrypt`, used with `certificate_chain`).
* `lets_encrypt` - (Optional) `true` for enabling lets_encrypt certificate auto renewal  (mutually exclusive with `certificate_chain` or `private_key`).
* `site_id` - (Required) .

## Attribute Reference

* `id` - The id of the sakuracloud_webaccel_certificate.
* `dns_names` - .
* `issuer_common_name` - .
* `not_after` - .
* `not_before` - .
* `serial_number` - .
* `sha256_fingerprint` - .
* `subject_common_name` - .
* `lets_encrypt` - .
