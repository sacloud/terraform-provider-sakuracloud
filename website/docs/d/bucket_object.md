---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_bucket_object"
subcategory: "ObjectStorage"
description: |-
  Get information about an existing Bucket Object.
---

# Data Source: sakuracloud_bucket_object

Get information about an existing Bucket Object.

## Argument Reference

* `access_key` - (Required) The access key for using SakuraCloud Object Storage API.
* `bucket` - (Required) The name of bucket.
* `key` - (Required) The name of the bucket object.
* `secret_key` - (Required) The secret key for using SakuraCloud Object Storage API.



## Attribute Reference

* `id` - The id of the Bucket Object.
* `body` - The body of the bucket object.
* `content_type` - The content type of the bucket object.
* `etag` - The etag of the bucket object.
* `http_cache_url` - The URL for cached access to the bucket object via HTTP.
* `http_path_url` - The URL with path-format for accessing the bucket object via HTTP.
* `http_url` - The URL for accessing the bucket object via HTTP.
* `https_cache_url` - The URL for cached access to the bucket object via HTTPS.
* `https_path_url` - The URL with path-format for accessing the bucket object via HTTPS.
* `https_url` - The URL for accessing the bucket object via HTTPS.
* `last_modified` - The time when the bucket object last modified.
* `size` - The size of the bucket object in bytes.




