---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_bucket_object"
subcategory: "ObjectStorage"
description: |-
  Get information about an existing Bucket Object.
---

# Data Source: sakuracloud_bucket_object

Get information about an existing Bucket Object.

## Example Usage

```hcl
data "sakuracloud_bucket_object" "foobar" {
  bucket = "foobar"
  key    = "key.txt"
}
```
## Argument Reference

* `access_key` - (Required) The access key for using SakuraCloud Object Storage API.
* `bucket` - (Required) The name of bucket.
* `key` - (Required) The name of the BucketObject.
* `secret_key` - (Required) The secret key for using SakuraCloud Object Storage API.



## Attribute Reference

* `id` - The id of the Bucket Object.
* `body` - The body of the BucketObject.
* `content_type` - The content type of the BucketObject.
* `etag` - The etag of the BucketObject.
* `http_cache_url` - The URL for cached access to the BucketObject via HTTP.
* `http_path_url` - The URL with path-format for accessing the BucketObject via HTTP.
* `http_url` - The URL for accessing the BucketObject via HTTP.
* `https_cache_url` - The URL for cached access to the BucketObject via HTTPS.
* `https_path_url` - The URL with path-format for accessing the BucketObject via HTTPS.
* `https_url` - The URL for accessing the BucketObject via HTTPS.
* `last_modified` - The time when the BucketObject last modified.
* `size` - The size of the BucketObject in bytes.



