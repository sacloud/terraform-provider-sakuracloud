---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_bucket_object"
sidebar_current: "docs-sakuracloud-resource-objectstorage-bucket-object"
description: |-
  Provides a SakuraCloud Bucket Object resource. This can be used to create, update, and delete Bucket Objects.
---

# sakuracloud\_bucket\_object

Provides a SakuraCloud Bucket Object resource. This can be used to create, update, and delete Bucket Objects.

~> **NOTE on Bucket:**  Sakura Cloud does not support bucket creation by API.
Buckets should be created on the control panel.

## Example Usage

```hcl
# Create a new Bucket Object
resource "sakuracloud_bucket_object" "foobar" {
  bucket = "your-bucket-name"
  key    = "path/to/your/object"
  
  source = "path/to/your/source/file"
  # or
  #content     = "your-content-body"
  
  content_type = "application/json"
}
```


## Argument Reference

* `bucket` - (Required) The name of bucket.
* `access_key` - (Required) The access key of bucket. It must be provided, but it can also be sourced from the `SACLOUD_OJS_ACCESS_KEY_ID` or `AWS_ACCESS_KEY_ID` environment variable.
* `secret_key` - (Required) The secret key of bucket. It must be provided, but it can also be sourced from the `SACLOUD_OJS_SECRET_ACCESS_KEY` or `AWS_SECRET_ACCESS_KEY` environment variable.
* `key` - (Required) The key of the bucket object.
* `source` - (Optional) Source file path of value of the bucket object.
* `content` - (Optional) String of the value of the bucket object. 
* `content_type` - (Optional) Content-Type header value of the bucket object.

## Attributes Reference

* `id` - ID of the resource.
* `content_type` - Content-Type header value of the bucket object.
* `body` - String of the value of the bucket object. Set when Content-Type is `"text/*"` or `"application/json"`.
* `etag` - ETag of the resource.
* `size` - Size of the resource (unit:`byte`).
* `last_modified` - Update date of the resource.
* `http_url` - URL for accessing the object via HTTP (type:`subdomain`).
* `https_url` - URL for accessing the object via HTTPS (type:`subdomain`).
* `http_path_url` - URL for accessing the object via HTTP (type:`path`).
* `http_cache_url` - URL for accessing the object via HTTP (type:`cache`).
* `https_cache_url` - URL for accessing the object via HTTPS (type:`cache`)..


## Import (not supported)

Import of Bucket Object is not supported.
