---
layout: "sakuracloud"
page_title: "Provider: SakuraCloud"
sidebar_current: "docs-sakuracloud-index"
description: |-
  The SakuraCloud provider is used to interact with Sakura Cloud (IaaS).
  The provider needs to be configured with the proper credentials before it can be used.
---

# SakuraCloud Provider

The SakuraCloud provider is used to interact with Sakura Cloud (IaaS).
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the SakuraCloud provider
provider "sakuracloud" {
  token  = "<your API token>"
  secret = "<your API secret>"
  zone   = "<target zone>" 
}
```

## Argument Reference

The following arguments are supported:

* `token` - (Required) The SakuraCloud API access token. It must be provided, but it can also be sourced from the `SAKURACLOUD_ACCESS_TOKEN` environment variable.
* `secret` - (Required) The SakuraCloud API access token secret. It must be provided, but it can also be sourced from the `SAKURACLOUD_ACCESS_TOKEN_SECRET` environment variable.
* `zone` - (Optional) The Default target zone of API operations. It can also be sourced from the `SAKURACLOUD_ZONE` environment variable. Default value is `is1b`.
* `accept_language` - (Optional) The value of `Accept-Language` header to be set at API call. It can also be sourced from the `SAKURACLOUD_ACCEPT_LANGUAGE` environment variable.
* `api_root_url` - (Optional) The root URL of API call destination. It can also be sourced from the `SAKURACLOUD_API_ROOT_URL` environment variable.
* `retry_max` - (Optional) The number of retries when an error (status=`503`) occurs in the API call. It can also be sourced from the `SAKURACLOUD_RETRY_MAX` environment variable. Default value is `10`.
* `retry_interval` - (Optional) The retry interval (seconds) when an error (status=`503`) occurs in the API call. It can also be sourced from the `SAKURACLOUD_RETRY_INTERVAL` environment variable. Default value is `5`.
* `timeout` - (Optional) The status change wait time in API call (minutes). It can also be sourced from the `SAKURACLOUD_TIMEOUT` environment variable. Default value is `20`.
* `trace` - (Optional) The flag of output logs at API call. It can also be sourced from the `SAKURACLOUD_TRACE_MODE` environment variable. 
