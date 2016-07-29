# サーバー(sakuracloud_server)

### 設定例

```
resource "sakuracloud_server" "myserver" {
    name = "myserver"
    disks = ["${sakuracloud_disk.mydisk.id}"]
    additional_interfaces = [""]
    description = "Server from TerraForm for SAKURA CLOUD"
    tags = ["@virtio-net-pci"]
}
```

### パラメーター

|パラメーター|必須  |名称                |初期値     |設定値 |補足                                          |
|----------|:---:|--------------------|:--------:|------|----------------------------------------------|
| `name`   | ◯   | サーバー名           | -   | 文字列 | - |
| `disks`  | ◯   | ディスクID          | -   | リスト(文字列) | サーバーに接続するディスクのID |
| `core`   | -   | CPUコア数           | 1   | 数値 | 指定可能な値は[こちら](http://cloud.sakura.ad.jp/specification/server-disk/)のプラン一覧を参照ください |
| `memory` | -   | メモリ(GB単位)       | 1  | 数値 | 指定可能な値は[こちら](http://cloud.sakura.ad.jp/specification/server-disk/)のプラン一覧を参照ください |
| `base_interface` | - | 基本NIC | `shared` | `shared`(共有セグメント)<br />`switch_id`(スイッチ)<br />`""`(接続なし)|eth0の上流NWとの接続方法を指定する。 |
| `additional_interfaces` | - | 追加NIC | - | リスト(文字列) | 追加で割り当てるNIC。接続するスイッチのID、または空文字を指定する。 |
| `packet_filter_ids`| - | パケットフィルタID | - | リスト(文字列) | NICに適用するパケットフィルタのIDをリストで指定する。リストの先頭からeth0,eth1の順で適用される |
| `description` | - | 説明 | - | 文字列 | - |
| `tags` | - | タグ | - | リスト(文字列) | サーバーに付与するタグ。@で始まる特殊タグについては[こちら](http://cloud-news.sakura.ad.jp/special-tags/)を参照 |
| `zone` | - | ゾーン | - | `is1b`<br />`tk1a`<br />`tk1v` | - |

### 属性

|属性名                    | 名称                     | 補足                                        |
|-------------------------|-------------------------|--------------------------------------------|
| `id`                    | ID                      | -                                          |
| `name`                  | サーバー名                | -                                          |
| `disks`                 | ディスクID                | -                                          |
| `core`                  | CPUコア数                 | -                                         |
| `memory`                | メモリ(GB単位)            | -                                          |
| `base_interface`        | 基本NIC                  | -                                         |
| `additional_interfaces` | 追加NIC                  | -                                         |
| `packet_filter_ids`     | パケットフィルタID         | -                                         |
| `description`           | 説明                     | -                                         |
| `tags`                  | タグ                     | -                                         |
| `zone`                  | ゾーン                    | -                                         |
| `macaddresses`         | MACアドレス               | MACアドレスのリスト(NICの個数分のリスト)        |
| `base_nw_ipaddress`     | 基本NIC-IPアドレス         | eth0のIPアドレス                            |
| `base_nw_dns_servers`   | 基本NIC-DNSサーバー        | eth0の属するセグメントの推奨ネームサーバーのリスト|
| `base_nw_gateway`       | 基本NIC-ゲートウェイ        | eth0の属するセグメントのゲートウェイIPアドレス   |
| `base_nw_address`       | 基本NIC-ネットワークアドレス | eth0のIPアドレスのネットワークアドレス          |
| `base_nw_mask_len`      | 基本NIC-サブネットマスク長   | eth0のIPアドレスのサブネットマスク長           |
