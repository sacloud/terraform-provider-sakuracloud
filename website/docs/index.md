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
terraform {
  required_providers {
    sakuracloud = {
      source = "sacloud/sakuracloud"

      # We recommend pinning to the specific version of the SakuraCloud Provider you're using
      # since new versions are released frequently
      version = "2.4.1"
      #version = "~> 2"
    }
  }
}
# Configure the SakuraCloud Provider
provider "sakuracloud" {
  # More information on the authentication methods supported by
  # the SakuraCloud Provider can be found here:
  # https://docs.usacloud.jp/terraform/provider/

  # profile = "..."
}
```

## Authentication Methods

The SakuraCloud provider supports following authentication methods:

- Static credentials
- Shared credentials file
- Environment variables

### Static credentials ###

!> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system.

Static credentials can be provided by adding a `token` and `secret`
in-line in the SakuraCloud provider block:

Usage:

```hcl
provider "sakuracloud" {
  token  = "my-access-token"
  secret = "my-access-secret"
  zone   = "is1a"
}
```
### Environment variables

You can provide your credentials via the `SAKURACLOUD_ACCESS_TOKEN` and
`SAKURACLOUD_ACCESS_TOKEN_SECRET`, environment variables, representing your 
Access Token and your Access Secret, respectively. 

```hcl
provider "sakuracloud" {}
```

Usage:

```sh
$ export SAKURACLOUD_ACCESS_TOKEN="my-access-token"
$ export SAKURACLOUD_ACCESS_TOKEN_SECRET="my-access-secret"
$ export SAKURACLOUD_ZONE="is1a"
$ terraform plan
```

### Shared credentials file

You can use a shared credentials file by specifying `profile` parameter.
A shared credentials file is formatted as JSON, the default location is `$HOME/.usacloud/<profile name>/config.json` on Linux and OS X, or
`"%USERPROFILE%\.usacloud\<profile name>/config.json"` for Windows users.

Example shared credentials file is follows:

```json
{
	"AccessToken": "my-access-token",
	"AccessTokenSecret": "my-access-secret",
	"Zone": "is1a"
}
```

## Argument Reference

* `accept_language` - (Optional) The value of AcceptLanguage header used when calling SakuraCloud API. It can also be sourced from the `SAKURACLOUD_ACCEPT_LANGUAGE` environment variables, or via a shared credentials file if `profile` is specified.
* `api_request_rate_limit` - (Optional) The maximum number of SakuraCloud API calls per second. It can also be sourced from the `SAKURACLOUD_RATE_LIMIT` environment variables, or via a shared credentials file if `profile` is specified. Default:`10`.
* `api_request_timeout` - (Optional) The timeout seconds for each SakuraCloud API call. It can also be sourced from the `SAKURACLOUD_API_REQUEST_TIMEOUT` environment variables, or via a shared credentials file if `profile` is specified. Default:`300`.
* `api_root_url` - (Optional) The root URL of SakuraCloud API. It can also be sourced from the `SAKURACLOUD_API_ROOT_URL` environment variables, or via a shared credentials file if `profile` is specified. Default:`https://secure.sakura.ad.jp/cloud/zone`.
* `default_zone` - (Optional) The name of zone to use as default for global resources. It must be provided, but it can also be sourced from the `SAKURACLOUD_DEFAULT_ZONE` environment variables, or via a shared credentials file if `profile` is specified.
* `fake_mode` - (Optional) The flag to enable fake of SakuraCloud API call. It is for debugging or developping the provider. It can also be sourced from the `FAKE_MODE` environment variables, or via a shared credentials file if `profile` is specified.
* `fake_store_path` - (Optional) The file path used by SakuraCloud API fake driver for storing fake data. It is for debugging or developping the provider. It can also be sourced from the `FAKE_STORE_PATH` environment variables, or via a shared credentials file if `profile` is specified.
* `profile` - (Optional) The profile name of your SakuraCloud account. Default:`default`.
* `retry_max` - (Optional) The maximum number of API call retries used when SakuraCloud API returns status code `423` or `503`. It can also be sourced from the `SAKURACLOUD_RETRY_MAX` environment variables, or via a shared credentials file if `profile` is specified. Default:`100`.
* `retry_wait_max` - (Optional) The maximum wait interval(in seconds) for retrying API call used when SakuraCloud API returns status code `423` or `503`.  It can also be sourced from the `SAKURACLOUD_RETRY_WAIT_MAX` environment variables, or via a shared credentials file if `profile` is specified.
* `retry_wait_min` - (Optional) The minimum wait interval(in seconds) for retrying API call used when SakuraCloud API returns status code `423` or `503`. It can also be sourced from the `SAKURACLOUD_RETRY_WAIT_MAX` environment variables, or via a shared credentials file if `profile` is specified.
* `secret` - (Optional) The API secret of your SakuraCloud account. It must be provided, but it can also be sourced from the `SAKURACLOUD_ACCESS_TOKEN_SECRET` environment variables, or via a shared credentials file if `profile` is specified.
* `token` - (Optional) The API token of your SakuraCloud account. It must be provided, but it can also be sourced from the `SAKURACLOUD_ACCESS_TOKEN` environment variables, or via a shared credentials file if `profile` is specified.
* `trace` - (Optional) The flag to enable output trace log. It can also be sourced from the `SAKURACLOUD_TRACE` environment variables, or via a shared credentials file if `profile` is specified.
* `zone` - (Optional) The name of zone to use as default. It must be provided, but it can also be sourced from the `SAKURACLOUD_ZONE` environment variables, or via a shared credentials file if `profile` is specified.
* `zones` - (Optional) A list of available SakuraCloud zone name. It can also be sourced via a shared credentials file if `profile` is specified. Default:[`is1a`, `is1b`, `tk1a`, `tk1v`].


