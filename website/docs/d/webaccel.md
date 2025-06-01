---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_webaccel"
subcategory: "WebAccelerator"
description: |-
  Get information about an existing WebAccelerator site.
---

# Data Source: sakuracloud_webaccel

Get information about an existing sakuracloud_webaccel.

## Argument Reference

Either of following arguments should be specified:

* `domain` - (Optional) domain name of the site.
* `name` - (Optional) nickname of the site.

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
* `s3_access_key_id` - blank attribute.
* `s3_secret_access_key` - blank attribute.
* `s3_doc_index` - (Optional) Whether the document indexing is enabled or not. `false` by default.

---

A `logging` block provides following attributes:

* `s3_bucket_name` - The bucket name for the log destination. In present, a bucket on the sakura object storage is supported.
* `s3_access_key_id` - blank attribute.
* `s3_secret_access_key` - blank attribute.
* `enabled` - Whether the access logging is enabled or not. `false` by default.
