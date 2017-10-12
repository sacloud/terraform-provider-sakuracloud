# サーバ(sakuracloud_server)

---

### 設定例

```hcl
resource sakuracloud_server "myserver" {
  name  = "myserver"
  disks = ["${sakuracloud_disk.mydisk.id}"]

  #コア数
  #core = 1

  #メモリサイズ(GB)
  #memory = 1

  #NICドライバ(virtio or e1000)
  #interface_driver = "virtio"

  #上流のNWとの接続方法
  #nic = "shared"

  #追加NIC
  #additional_nics = ["${sakuracloud_switch.myswitch.id}"]

  #パケットフィルタ
  #packet_filter_ids = ["${sakuracloud_packet_filter.myfilter.id}"]

  #ISOイメージ(CD-ROM)
  #cdrom_id = "${data.sakuracloud_cdrom.mycdrom.id}"

  #ネットワーク設定(nicにスイッチIDが指定されている場合のみ)
  #ipaddress   = "192.168.0.101"
  #gateway     = "192.168.0.1"
  #nw_mask_len = 24

  description = "Server from TerraForm for SAKURA CLOUD"
}
```

### パラメーター

|パラメーター|必須  |名称                |初期値     |設定値 |補足                                          |
|----------|:---:|--------------------|:--------:|------|----------------------------------------------|
| `name`   | ◯   | サーバ名           | -   | 文字列 | - |
| `disks`  | ◯   | ディスクID          | -   | リスト(文字列) | サーバに接続するディスクのID |
| `core`   | -   | CPUコア数           | 1   | 数値 | 指定可能な値は[こちら](http://cloud.sakura.ad.jp/specification/server-disk/)のプラン一覧を参照ください |
| `memory` | -   | メモリ(GB単位)       | 1  | 数値 | 指定可能な値は[こちら](http://cloud.sakura.ad.jp/specification/server-disk/)のプラン一覧を参照ください |
| `interface_driver` | -   | NICドライバ       | `virtio`  | `virtio`<br />`e1000` | - |
| `nic` | - | 基本NIC | `shared` | `shared`(共有セグメント)<br />`[switch_id]`(スイッチのID)<br />`""`(接続なし)|eth0の上流NWとの接続方法を指定する。 |
| `additional_nics` | - | 追加NIC | - | リスト(文字列) | 追加で割り当てるNIC。接続するスイッチのID、または空文字を指定する。 |
| `packet_filter_ids`| - | パケットフィルタID | - | リスト(文字列) | NICに適用するパケットフィルタのIDをリストで指定する。リストの先頭からeth0,eth1の順で適用される |
| `icon_id`       | -   | アイコンID         | - | 文字列| - |
| `description` | - | 説明 | - | 文字列 | - |
| `cdrom_id` | - | CDROM(ISOイメージ)ID | - | 文字列 | - |
| `ipaddress`| - | 基本NIC:IPアドレス | - | 文字列 | [注1](#注1) |
| `gateway`  | - | 基本NIC:ゲートウェイ | - | 文字列 | [注1](#注1) |
| `nw_mask_len` | - | 基本NIC:サブネットマスク長 | - | 文字列 | [注1](#注1) |
| `tags` | - | タグ | - | リスト(文字列) | サーバに付与するタグ。@で始まる特殊タグについては[こちら](http://cloud-news.sakura.ad.jp/special-tags/)を参照 |
| `graceful_shutdown_timeout` | - | シャットダウンまでの待ち時間 | - | 数値(秒数) | シャットダウンが必要な場合の通常シャットダウンするまでの待ち時間(指定の時間まで待ってもシャットダウンしない場合は強制シャットダウンされる) |
| `zone` | - | ゾーン | - | `is1b`<br />`tk1a`<br />`tk1v` | - |

#### 注1

`nic`にスイッチのIDが指定されており、かつ`disks`の最初のパラメーターに
ディスクの修正に対応しているディスクのIDが指定されている場合に有効。
ディスクの修正は主にLinux系パブリックアーカイブを元にしたディスクの場合にサポートされています。

### 属性

|属性名                    | 名称                     | 補足                                        |
|-------------------------|-------------------------|--------------------------------------------|
| `id`                    | ID                      | -                                          |
| `macaddresses`          | MACアドレス               | MACアドレスのリスト(NICの個数分のリスト)        |
| `dns_servers`           | 基本NIC:DNSサーバ        | eth0の属するセグメントの推奨ネームサーバのリスト|
| `nw_address`            | 基本NIC:ネットワークアドレス | eth0のIPアドレスのネットワークアドレス          |
