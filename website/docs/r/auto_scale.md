---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_auto_scale"
subcategory: "Misc"
description: |-
  Manages a SakuraCloud sakuracloud_auto_scale.
---

# sakuracloud_auto_scale

Manages a SakuraCloud sakuracloud_auto_scale.

## Example Usage

```hcl
locals {
  zone               = "is1a"
  server_name_prefix = "target-server-"
  api_key_id         = "<your-api-key>"
}

resource "sakuracloud_auto_scale" "foobar" {
  name = "example"

  # 監視対象が存在するゾーン
  zones = [local.zone]

  # 設定ファイル
  config = jsonencode({
    resources : [{
      type : "Server",
      selector : {
        names : [sakuracloud_server.foobar.name],
        zones : [local.zone],
      },
    }],
  })

  # APIキーのID
  api_key_id = local.api_key_id

  # しきい値
  cpu_threshold_scaling {
    # 監視対象のサーバ名のプリフィックス
    server_prefix = local.server_name_prefix

    # 性能アップするCPU使用率
    up = 80

    # 性能ダウンするCPU使用率
    down = 20
  }
}

resource "sakuracloud_server" "foobar" {
  name = local.server_name_prefix
  force_shutdown = true
  zone = local.zone
}
```
## Argument Reference

* `api_key_id` - (Required) The id of the API key.. Changing this forces a new resource to be created.
* `config` - (Required) The configuration file for sacloud/autoscaler.
* `cpu_threshold_scaling` - (Required) A `cpu_threshold_scaling` block as defined below.
* `description` - (Optional) The description of the AutoScale. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the AutoScale.
* `name` - (Required) The name of the AutoScale. The length of this value must be in the range [`1`-`64`].
* `router_threshold_scaling` - (Optional) A `router_threshold_scaling` block as defined below.
* `tags` - (Optional) Any tags to assign to the AutoScale.
* `trigger_type` - (Optional) This must be one of [`cpu`/`router`].
* `zones` - (Required) List of zone names where monitored resources are located.

---

A `cpu_threshold_scaling` block supports the following:

* `server_prefix` - (Required) Server name prefix to be monitored. 
* `up` - (Required) Threshold for average CPU utilization to scale up/out. 
* `down` - (Required) Threshold for average CPU utilization to scale down/in.

---

A `router_threshold_scaling` block supports the following:

* `router_prefix` - (Required) Router name prefix to be monitored.
* `direction` - (Required) This must be one of [`in`/`out`].
* `mbps` - (Required) Mbps.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the sakuracloud_auto_scale
* `update` - (Defaults to 5 minutes) Used when updating the sakuracloud_auto_scale
* `delete` - (Defaults to 5 minutes) Used when deleting sakuracloud_auto_scale

## Attribute Reference

* `id` - The id of the sakuracloud_auto_scale.



