#terraform {
#  required_providers {
#    tls = {
#      source  = "hashicorp/tls"
#      version = "3.1.0"
#    }
#    sakuracloud = {
#      source  = "sacloud/sakuracloud"
#      version = "2.22.2"
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

