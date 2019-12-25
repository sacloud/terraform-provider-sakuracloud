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

* `accept_language` - (Optional) .
* `api_request_rate_limit` - (Optional) .
* `api_request_timeout` - (Optional) .
* `api_root_url` - (Optional) .
* `fake_mode` - (Optional) .
* `fake_store_path` - (Optional) .
* `profile` - (Optional) Your SakuraCloud Profile Name.
* `retry_max` - (Optional) .
* `retry_wait_max` - (Optional) .
* `retry_wait_min` - (Optional) .
* `secret` - (Optional) Your SakuraCloud APIKey(secret).
* `token` - (Optional) Your SakuraCloud APIKey(token).
* `trace` - (Optional) .
* `zone` - (Optional) Target SakuraCloud Zone(is1a | is1b | tk1a | tk1v).
* `zones` - (Optional) Available SakuraCloud Zones(default: [is1a, is1b, tk1a, tk1v]).



## Attribute Reference

* `id` - The ID of the SakuraCloud.




