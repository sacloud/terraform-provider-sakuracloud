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
  config = yamlencode({
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