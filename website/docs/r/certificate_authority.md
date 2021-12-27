---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_certificate_authority"
subcategory: "Lab"
description: |-
  Manages a SakuraCloud sakuracloud_certificate_authority.
---

# sakuracloud_certificate_authority

Manages a SakuraCloud sakuracloud_certificate_authority.

## Example Usage

```hcl
#terraform {
#  required_providers {
#    tls = {
#      source  = "hashicorp/tls"
#      version = "3.1.0"
#    }
#    sakuracloud = {
#      source  = "sacloud/sakuracloud"
#      version = "2.16.0"
#    }
#  }
#}

resource "tls_private_key" "client_key" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

resource "tls_cert_request" "client_csr" {
  key_algorithm   = "ECDSA"
  private_key_pem = tls_private_key.client_key.private_key_pem

  subject {
    common_name  = "client-csr.usacloud.com"
    organization = "usacloud"
  }
}

resource "tls_private_key" "server_key" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

resource "tls_cert_request" "server_csr" {
  key_algorithm   = "ECDSA"
  private_key_pem = tls_private_key.server_key.private_key_pem

  subject {
    common_name  = "server-csr.usacloud.com"
    organization = "usacloud"
  }
}

resource "sakuracloud_certificate_authority" "foobar" {
  name = "foobar"

  validity_period_hours = 24 * 3650

  subject {
    common_name        = "pki.usacloud.jp"
    country            = "JP"
    organization       = "usacloud"
    organization_units = ["ou1", "ou2"]
  }

  # by public_key
  client {
    subject {
      common_name        = "client1.usacloud.jp"
      country            = "JP"
      organization       = "usacloud"
      organization_units = ["ou1", "ou2"]
    }
    validity_period_hours = 24 * 3650
    public_key            = tls_private_key.client_key.public_key_pem
  }

  // by CSR
  client {
    subject {
      common_name        = "client2.usacloud.jp"
      country            = "JP"
      organization       = "usacloud"
      organization_units = ["ou1", "ou2"]
    }
    validity_period_hours = 24 * 3650
    csr                   = tls_cert_request.client_csr.cert_request_pem
  }

  # by email
  client {
    subject {
      common_name        = "client3.usacloud.jp"
      country            = "JP"
      organization       = "usacloud"
      organization_units = ["ou1", "ou2"]
    }
    validity_period_hours = 24 * 3650
    email                 = "example@example.com"
  }

  # by URL
  client {
    subject {
      common_name        = "client4.usacloud.jp"
      country            = "JP"
      organization       = "usacloud"
      organization_units = ["ou1", "ou2"]
    }
    validity_period_hours = 24 * 3650
  }

  # by public key
  server {
    subject {
      common_name        = "server1.usacloud.jp"
      country            = "JP"
      organization       = "usacloud"
      organization_units = ["ou1", "ou2"]
    }

    subject_alternative_names = ["alt1.usacloud.jp", "alt2.usacloud.jp"]

    validity_period_hours = 24 * 3650
    public_key            = tls_private_key.server_key.public_key_pem
  }

  # by CSR
  server {
    subject {
      common_name        = "server2.usacloud.jp"
      country            = "JP"
      organization       = "usacloud"
      organization_units = ["ou1", "ou2"]
    }

    subject_alternative_names = ["alt1.usacloud.jp", "alt2.usacloud.jp"]

    validity_period_hours = 24 * 3650
    csr                   = tls_cert_request.server_csr.cert_request_pem
  }
}


```
## Argument Reference

* `client` - (Optional) One or more `client` blocks as defined below.
* `description` - (Optional) The description of the Certificate Authority. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the Certificate Authority.
* `name` - (Required) The name of the Certificate Authority. The length of this value must be in the range [`1`-`64`].
* `server` - (Optional) One or more `server` blocks as defined below.
* `subject` - (Required) A `subject` block as defined below. Changing this forces a new resource to be created.
* `tags` - (Optional) Any tags to assign to the Certificate Authority.
* `validity_period_hours` - (Required) The number of hours after initial issuing that the certificate will become invalid. Changing this forces a new resource to be created.


---

A `client` block supports the following:

* `csr` - (Optional) Input for issuing a certificate.
* `email` - (Optional) Input for issuing a certificate.
* `hold` - (Optional) Flag to suspend/hold the certificate.
* `public_key` - (Optional) Input for issuing a certificate.
* `subject` - (Required) A `subject` block as defined below.
* `validity_period_hours` - (Required) The number of hours after initial issuing that the certificate will become invalid.

---

A `server` block supports the following:

* `csr` - (Optional) Input for issuing a certificate.
* `hold` - (Optional) Flag to suspend/hold the certificate.
* `public_key` - (Optional) Input for issuing a certificate.
* `subject` - (Required) A `subject` block as defined below.
* `subject_alternative_names` - (Optional) .
* `validity_period_hours` - (Required) The number of hours after initial issuing that the certificate will become invalid.

---

A `subject` block supports the following:

* `common_name` - (Required) .
* `country` - (Required) .
* `organization` - (Required) .
* `organization_units` - (Optional) .


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the sakuracloud_certificate_authority
* `update` - (Defaults to 5 minutes) Used when updating the sakuracloud_certificate_authority
* `delete` - (Defaults to 5 minutes) Used when deleting sakuracloud_certificate_authority



## Attribute Reference

* `id` - The id of the sakuracloud_certificate_authority.
* `certificate` - The body of the CA's certificate in PEM format.
* `crl_url` - The URL of the CRL.
* `not_after` - The date on which the certificate validity period ends, in RFC3339 format.
* `not_before` - The date on which the certificate validity period begins, in RFC3339 format.
* `serial_number` - The body of the CA's certificate in PEM format.


---

A `client` block exports the following:

* `certificate` - The body of the CA's certificate in PEM format.
* `id` - The id of the certificate.
* `issue_state` - Current state of the certificate.
* `not_after` - The date on which the certificate validity period ends, in RFC3339 format.
* `not_before` - The date on which the certificate validity period begins, in RFC3339 format.
* `serial_number` - The body of the CA's certificate in PEM format.
* `url` - The URL for issuing the certificate.

---

A `server` block exports the following:

* `certificate` - The body of the CA's certificate in PEM format.
* `id` - The id of the certificate.
* `issue_state` - Current state of the certificate.
* `not_after` - The date on which the certificate validity period ends, in RFC3339 format.
* `not_before` - The date on which the certificate validity period begins, in RFC3339 format.
* `serial_number` - The body of the CA's certificate in PEM format.


