---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_certificate_authority"
subcategory: "Lab"
description: |-
  Get information about an existing sakuracloud_certificate_authority.
---

# Data Source: sakuracloud_certificate_authority

Get information about an existing sakuracloud_certificate_authority.

## Example Usage

```hcl
data "sakuracloud_certificate_authority" "foobar" {
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

---

A `server` block supports the following:

* `hold` - (Optional) Flag to suspend/hold the certificate.


## Attribute Reference

* `id` - The id of the sakuracloud_certificate_authority.
* `certificate` - The body of the CA's certificate in PEM format.
* `client` - A list of `client` blocks as defined below.
* `crl_url` - The URL of the CRL.
* `description` - The description of the CertificateAuthority.
* `icon_id` - The icon id attached to the CertificateAuthority.
* `name` - The name of the CertificateAuthority.
* `not_after` - The date on which the certificate validity period ends, in RFC3339 format.
* `not_before` - The date on which the certificate validity period begins, in RFC3339 format.
* `serial_number` - The body of the CA's certificate in PEM format.
* `server` - A list of `server` blocks as defined below.
* `subject_string` - .
* `tags` - Any tags assigned to the CertificateAuthority.


---

A `client` block exports the following:

* `certificate` - The body of the CA's certificate in PEM format.
* `hold` - Flag to suspend/hold the certificate.
* `id` - The id of the certificate.
* `issue_state` - Current state of the certificate.
* `not_after` - The date on which the certificate validity period ends, in RFC3339 format.
* `not_before` - The date on which the certificate validity period begins, in RFC3339 format.
* `serial_number` - The body of the CA's certificate in PEM format.
* `subject_string` - .
* `url` - The URL for issuing the certificate.

---

A `server` block exports the following:

* `certificate` - The body of the CA's certificate in PEM format.
* `id` - The id of the certificate.
* `issue_state` - Current state of the certificate.
* `not_after` - The date on which the certificate validity period ends, in RFC3339 format.
* `not_before` - The date on which the certificate validity period begins, in RFC3339 format.
* `serial_number` - The body of the CA's certificate in PEM format.
* `subject_alternative_names` - .
* `subject_string` - .


