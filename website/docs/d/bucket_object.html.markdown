---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_bucket_object"
sidebar_current: "docs-sakuracloud-datasource-bucket-object"
description: |-
  Get information on a SakuraCloud Bucket Object.
---

# sakuracloud\_bucket\_object

Use this data source to retrieve information about a SakuraCloud Bucket Object.

## Example Usage

```hcl
data "sakuracloud_bucket_object" "foobar" {
   bucket = "your-bucket-name"
   key    = "path/to/your/object"
 }

```

## Argument Reference

* `bucket` - (Required) The name of bucket.
* `access_key` - (Required) The access key of bucket. It must be provided, but it can also be sourced from the `SACLOUD_OJS_ACCESS_KEY_ID` or `AWS_ACCESS_KEY_ID` environment variable.
* `secret_key` - (Required) The secret key of bucket. It must be provided, but it can also be sourced from the `SACLOUD_OJS_SECRET_ACCESS_KEY` or `AWS_SECRET_ACCESS_KEY` environment variable.
* `key` - (Required) The key of the bucket object.

## Attributes Reference

* `id` - ID of the resource.
* `content_type` - Content-Type header value of the bucket object.
* `body` - String of the value of the bucket object. Set when Content-Type is `"text/*"` or `"application/json"`.
* `etag` - ETag of the resource.
* `size` - Size of the resource(unit: `byte`).
* `last_modified` - Update date of the resource.
* `http_url` - URL for accessing the object via HTTP(type: `subdomain`).
* `https_url` - URL for accessing the object via HTTPS(type: `subdomain`).
* `http_path_url` - URL for accessing the object via HTTP(type: `path`).
* `http_cache_url` - URL for accessing the object via HTTP(type: `cache`).
* `https_cache_url` - URL for accessing the object via HTTPS(type: `cache`)..

