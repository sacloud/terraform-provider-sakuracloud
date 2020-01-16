---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_bucket_object"
subcategory: "ObjectStorage"
description: |-
  Manages a SakuraCloud Bucket Object.
---

# sakuracloud_bucket_object

Manages a SakuraCloud Bucket Object.

## Example Usage

```hcl
resource "sakuracloud_bucket_object" "foobar" {
  bucket  = "foobar"
  key     = "example.txt"
  content = file("example.txt")
}
```

## Argument Reference

* `access_key` - (Required) The access key for using SakuraCloud Object Storage API.
* `secret_key` - (Required) The secret key for using SakuraCloud Object Storage API.
* `bucket` - (Required) The name of the bucket. Changing this forces a new resource to be created.
* `key` - (Required) The name of the bucket object. Changing this forces a new resource to be created.
* `content` - (Optional) The content to upload to as the bucket object. This conflicts with [`source`].
* `content_type` - (Optional) The content-type of the bucket object.
* `etag` - (Optional) The etag of the bucket object.
* `source` - (Optional) The file path to upload to as the bucket object. This conflicts with [`content`].

## Attribute Reference

* `id` - The id of the Bucket Object.
* `http_cache_url` - The URL for cached access to the bucket object via HTTP.
* `http_path_url` - The URL with path-format for accessing the bucket object via HTTP.
* `http_url` - The URL for accessing the bucket object via HTTP.
* `https_cache_url` - The URL for cached access to the bucket object via HTTPS.
* `https_path_url` - The URL with path-format for accessing the bucket object via HTTPS.
* `https_url` - The URL for accessing the bucket object via HTTPS.
* `last_modified` - The time when the bucket object last modified.
* `size` - The size of the bucket object in bytes.

