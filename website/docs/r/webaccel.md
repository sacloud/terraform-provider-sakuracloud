---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_webaccel"
subcategory: "WebAccelerator"
description: |-
  Manages a SakuraCloud WebAccelerator Site.
---

# sakuracloud_webaccel

Manages a SakuraCloud Web Accelerator Site.

~> **Note:** This resource cannot activate the site immediately. Use `webaccel_activation` resource to activate the site.

## Example Usage

```hcl
resource "sakuracloud_webaccel" "foobar" {
  name             = "hoge"
  domain_type      = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type     = "web"
    origin   = "docs.usacloud.jp"
    protocol = "https"
  }
  logging {
    s3_bucket_name = "example-bucket"
    s3_access_key_id = "xxxxxxxxxxxxxxx"
    s3_secret_access_key = "xxxxxxxxxxxxxxxxxxxxxxx"
    enabled = true
  }
  onetime_url_secrets = [
    "abc-0x123456"
  ]
  vary_support      = true
  default_cache_ttl = 3600
  normalize_ae      = "gzip"
}
```

## Argument Reference

* `name` - (Required) The nickname of the site.
* `domain_type` - (Required) Domain type for the site. One of `own_domain` or `subdomain`.
* `request_protocol` - (Required) Request protocol for the site. One of `http+https`, `https` or `https-redirect`.
* `onetime_url_secrets` - (Optional) The list of onetime URL secrets. Three or more secrets cannot be accepted.
* `vary_support` - (Optional) Whether the site supports VARY header or not.
* `default_cache_ttl` - (Optional) The default TTL for the content cache in seconds.
* `normalize_ae` - Compression target of accept-encoding normalization. One of `gzip` or `br+gzip`.

---

An `origin_parameters` block supports following parameters:

#### Core Origin Arguments

* `type` - (Required) The origin type. Either of `web` or `bucket`.

#### Web Origin Arguments (acceptable for `web` origin)

* `origin` - (Required) The origin hostname or IPv4 address.
* `protocol` - (Required) Origin access protocol. Either of `http` or `https`.
* `host_header` - (Optional) HTTP Host header for the origin access.

#### Bucket Origin Arguments (acceptable for `bucket` origin)

* `s3_endpoint` - (Required) S3 endpoint without protocol scheme. Specify `s3.isk01.sakurastorage.jp` for the sakura object storage.
* `s3_region` - (Required) S3 region. Specify `jp-north-1` for the sakura object storage.
* `s3_bucket_name` - The origin bucket name. Bucket prefix is not supported at now.
* `s3_access_key_id` - The S3 access key ID for the bucket.
* `s3_secret_access_key` - The S3 secret access key for the bucket.
* `s3_doc_index` - (Optional) Whether the document indexing is enabled or not. `false` by default.

---

A `logging` block supports following parameters:

* `s3_bucket_name` - The bucket name for the log destination. In present, a bucket on the sakura object storage is supported.
* `s3_access_key_id` - The S3 access key ID for the bucket.
* `s3_secret_access_key` - The S3 secret access key for the bucket.
* `enabled` - Whether the access logging is enabled or not. `false` by default.

## Attribute Reference

* `id` - The id of the sakuracloud_webaccel.
* `cname_record_value` - CNAME record value used to point to the site.
* `subdomain` - FQDN of the site.
* `txt_record_value` - TXT record value which is used to verify the domain ownership.
* `request_protocol` - (Required) Request protocol for the site. One of `http+https`, `https` or `https-redirect`.
* `vary_support` - Whether the site supports VARY header or not.
* `default_cache_ttl` - The default TTL for the content cache in seconds.
* `normalize_ae` - Compression target of accept-encoding normalization. One of `gzip` or `br+gzip`.

---

An `origin_parameters` block provides following attributes:

* `type` - (Required) The origin type. Either of `web` or `bucket`.
* `origin` - (Required) The origin hostname or IPv4 address.
* `protocol` - (Required) Origin access protocol. Either of `http` or `https`.
* `host_header` - (Optional) HTTP Host header for the origin access.
* `s3_endpoint` - (Required) S3 endpoint without protocol scheme. Specify `s3.isk01.sakurastorage.jp` for the sakura object storage.
* `s3_region` - (Required) S3 region. Specify `jp-north-1` for the sakura object storage.
* `s3_bucket_name` - The origin bucket name. Bucket prefix is not supported at now.
* `s3_access_key_id` - The S3 access key ID for the bucket.
* `s3_secret_access_key` - The S3 secret access key for the bucket.
* `s3_doc_index` - (Optional) Whether the document indexing is enabled or not. `false` by default.

---

A `logging` block provides following attributes:

* `s3_bucket_name` - The bucket name for the log destination. In present, a bucket on the sakura object storage is supported.
* `s3_access_key_id` - The S3 access key ID for the bucket.
* `s3_secret_access_key` - The S3 secret access key for the bucket.
* `enabled` - Whether the access logging is enabled or not. `false` by default.
