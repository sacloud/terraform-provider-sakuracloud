# サーバ(sakuracloud_server)

---

### 設定例

```hcl
resource "sakuracloud_server" "myserver" {
    name = "myserver"
    disks = ["${sakuracloud_disk.mydisk.id}"]
    
    # コア数
    # core = 1
    
    # メモリサイズ(GB)
    # memory = 1

    # 上流のNWとの接続方法
    # nic = "shared"

    # 追加NIC
    # additional_nics = ["${sakuracloud_switch.myswitch.id}"]

    # パケットフィルタ
    # packet_filter_ids = ["${sakuracloud_packet_filter.myfilter.id}"]
    
    # ISOイメージ(CD-ROM)
    # cdrom_id = "${data.sakuracloud_cdrom.mycdrom.id}"

    # ネットワーク設定(nicにスイッチIDが指定されている場合のみ)
    # ipaddress = "192.168.0.101"
    # gateway = "192.168.0.1"
    # nw_mask_len = 24

    description = "Server from TerraForm for SAKURA CLOUD"
    tags = ["@virtio-net-pci"]
}
```

### パラメーター

|パラメーター|必須  |名称                |初期値     |設定値 |補足                                          |
|----------|:---:|--------------------|:--------:|------|----------------------------------------------|
| `name`   | ◯   | サーバ名           | -   | 文字列 | - |
| `disks`  | ◯   | ディスクID          | -   | リスト(文字列) | サーバに接続するディスクのID |
| `core`   | -   | CPUコア数           | 1   | 数値 | 指定可能な値は[こちら](http://cloud.sakura.ad.jp/specification/server-disk/)のプラン一覧を参照ください |
| `memory` | -   | メモリ(GB単位)       | 1  | 数値 | 指定可能な値は[こちら](http://cloud.sakura.ad.jp/specification/server-disk/)のプラン一覧を参照ください |
| `nic` | - | 基本NIC | `shared` | `shared`(共有セグメント)<br />`[switch_id]`(スイッチのID)<br />`""`(接続なし)|eth0の上流NWとの接続方法を指定する。 |
| `additional_nics` | - | 追加NIC | - | リスト(文字列) | 追加で割り当てるNIC。接続するスイッチのID、または空文字を指定する。 |
| `packet_filter_ids`| - | パケットフィルタID | - | リスト(文字列) | NICに適用するパケットフィルタのIDをリストで指定する。リストの先頭からeth0,eth1の順で適用される |
| `cdrom_id` | - | CDROM(ISOイメージ)ID | - | 文字列 | - |
| `private_host_id` | - | 専有ホストID | - | 文字列 | 専有ホストは東京第１ゾーン(tk1a)でのみ利用可能 |
| `ipaddress`| - | 基本NIC-IPアドレス | - | 文字列 | [注1](#注1) |
| `gateway`  | - | 基本NIC-ゲートウェイ | - | 文字列 | [注1](#注1) |
| `nw_mask_len` | - | 基本NIC-サブネットマスク長 | - | 文字列 | [注1](#注1) |
| `description` | - | 説明 | - | 文字列 | - |
| `tags` | - | タグ | - | リスト(文字列) | サーバに付与するタグ。@で始まる特殊タグについては[こちら](http://cloud-news.sakura.ad.jp/special-tags/)を参照 |
| `zone` | - | ゾーン | - | `is1b`<br />`tk1a`<br />`tk1v` | - |

#### 注1

`nic`にスイッチのIDが指定されており、かつ`disks`の最初のパラメーターに
ディスクの修正に対応しているディスクのIDが指定されている場合に有効。
ディスクの修正は主にLinux系パブリックアーカイブを元にしたディスクの場合にサポートされています。

### 属性

|属性名                    | 名称                     | 補足                                        |
|-------------------------|-------------------------|--------------------------------------------|
| `id`                    | ID                      | -                                          |
| `name`                  | サーバ名                | -                                          |
| `disks`                 | ディスクID                | -                                          |
| `core`                  | CPUコア数                 | -                                         |
| `memory`                | メモリ(GB単位)            | -                                          |
| `nic`                   | 基本NIC                  | -                                         |
| `additional_nics`       | 追加NIC                  | -                                         |
| `packet_filter_ids`     | パケットフィルタID         | -                                         |
| `cdrom_id`              | CDROM(ISOイメージ)ID         | -                                         |
| `private_host_id`       | 専有ホストID              | -                                         |
| `description`           | 説明                     | -                                         |
| `tags`                  | タグ                     | -                                         |
| `zone`                  | ゾーン                    | -                                         |
| `macaddresses`          | MACアドレス               | MACアドレスのリスト(NICの個数分のリスト)        |
| `ipaddress`             | 基本NIC-IPアドレス         | eth0のIPアドレス                            |
| `dns_servers`           | 基本NIC-DNSサーバ        | eth0の属するセグメントの推奨ネームサーバのリスト|
| `gateway`               | 基本NIC-ゲートウェイ        | eth0の属するセグメントのゲートウェイIPアドレス   |
| `nw_address`            | 基本NIC-ネットワークアドレス | eth0のIPアドレスのネットワークアドレス          |
| `nw_mask_len`           | 基本NIC-サブネットマスク長   | eth0のIPアドレスのサブネットマスク長           |
