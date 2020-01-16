---
layout: "sakuracloud"
page_title: "Provider: SakuraCloud"
description: |-
  The SakuraCloud Provider is used to interact with the many resources supported by its APIs.
---

# SakuraCloud Provider

The SakuraCloud Provider is used to interact with the many resources supported by its APIs.

## Example Usage

```hcl
# Configure the SakuraCloud Provider
provider "sakuracloud" {
  # We recommend pinning to the specific version of the SakuraCloud Provider you're using
  # since new versions are released frequently
  version = "=2.0.0"

  # More information on the authentication methods supported by
  # the SakuraCloud Provider can be found here:
  # https://docs.usacloud.jp/terraform/configuration/provider/

  # profile = "..."
}
```
## Argument Reference

* `accept_language` - (Optional) The value of AcceptLanguage header used when calling SakuraCloud API. It can also be sourced from the `SAKURACLOUD_ACCEPT_LANGUAGE` environment variables, or via a shared credentials file if `profile` is specified..
* `api_request_rate_limit` - (Optional) The maximum number of SakuraCloud API calls per second. It can also be sourced from the `SAKURACLOUD_RATE_LIMIT` environment variables, or via a shared credentials file if `profile` is specified. Default:`%!s(int=10)`.
* `api_request_timeout` - (Optional) The timeout seconds for each SakuraCloud API call. It can also be sourced from the `SAKURACLOUD_API_REQUEST_TIMEOUT` environment variables, or via a shared credentials file if `profile` is specified. Default:`%!s(int=300)`.
* `api_root_url` - (Optional) The root URL of SakuraCloud API. It can also be sourced from the `SAKURACLOUD_API_ROOT_URL` environment variables, or via a shared credentials file if `profile` is specified. Default:`https://secure.sakura.ad.jp/cloud/zone`.
* `fake_mode` - (Optional) The flag to enable fake of SakuraCloud API call. It is for debugging or developping the provider. It can also be sourced from the `FAKE_MODE` environment variables, or via a shared credentials file if `profile` is specified..
* `fake_store_path` - (Optional) The file path used by SakuraCloud API fake driver for storing fake data. It is for debugging or developping the provider. It can also be sourced from the `FAKE_STORE_PATH` environment variables, or via a shared credentials file if `profile` is specified..
* `profile` - (Optional) The profile name of your SakuraCloud account. Default:`default`.
* `retry_max` - (Optional) The maximum number of API call retries used when SakuraCloud API returns status code `429` or `503`. It can also be sourced from the `SAKURACLOUD_RETRY_MAX` environment variables, or via a shared credentials file if `profile` is specified. Default:`100`.
* `retry_wait_max` - (Optional) The maximum wait interval(in seconds) for retrying API call used when SakuraCloud API returns status code `429` or `503`.  It can also be sourced from the `SAKURACLOUD_RETRY_WAIT_MAX` environment variables, or via a shared credentials file if `profile` is specified..
* `retry_wait_min` - (Optional) The minimum wait interval(in seconds) for retrying API call used when SakuraCloud API returns status code `429` or `503`. It can also be sourced from the `SAKURACLOUD_RETRY_WAIT_MAX` environment variables, or via a shared credentials file if `profile` is specified..
* `secret` - (Optional) The API secret of your SakuraCloud account. It must be provided, but it can also be sourced from the `SAKURACLOUD_ACCESS_TOKEN_SECRET` environment variables, or via a shared credentials file if `profile` is specified..
* `token` - (Optional) The API token of your SakuraCloud account. It must be provided, but it can also be sourced from the `SAKURACLOUD_ACCESS_TOKEN` environment variables, or via a shared credentials file if `profile` is specified..
* `trace` - (Optional) The flag to enable output trace log. It can also be sourced from the `SAKURACLOUD_TRACE` environment variables, or via a shared credentials file if `profile` is specified..
* `zone` - (Optional) The name of zone to use as default. It must be provided, but it can also be sourced from the `SAKURACLOUD_ZONE` environment variables, or via a shared credentials file if `profile` is specified..
* `zones` - (Optional) A list of available SakuraCloud zone name. It can also be sourced via a shared credentials file if `profile` is specified. Default:[`is1a`, `is1b`, `tk1a`, `tk1v`].



