# サーバ(sakuracloud_server)

---

### 設定例

```hcl
resource "sakuracloud_server" "myserver" {
  name  = "myserver"
  disks = [sakuracloud_disk.mydisk.id]

  #コア数
  #core = 1

  #メモリサイズ(GB)
  #memory = 1

  #NICドライバ(virtio or e1000)
  #interface_driver = "virtio"

  #上流のNWとの接続方法
  #nic = "shared"

  #追加NIC
  #additional_nics = [sakuracloud_switch.myswitch.id]

  #パケットフィルタ
  #packet_filter_ids = [sakuracloud_packet_filter.myfilter.id]

  #ISOイメージ(CD-ROM)
  #cdrom_id = data.sakuracloud_cdrom.mycdrom.id

  #ネットワーク設定(nicにスイッチIDが指定されている場合のみ)
  #ipaddress   = "192.168.0.101"
  #gateway     = "192.168.0.1"
  #nw_mask_len = 24

  description = "Server from TerraForm for SAKURA CLOUD"
  tags        = ["tag1", "tag2"]

  #==========================
  #ディスクの修正関連
  #==========================
  hostname = "myserver" #ホスト名
  password = "p@ssw0rd" #パスワード
  
  #SSH公開鍵
  #ssh_key_ids = [sakuracloud_ssh_key_gen.key.id]
  
  #スタートアップスクリプト
  #note_ids = [sakuracloud_note.note.id]
  
  #SSH接続でのパスワード/チャレンジレスポンス認証無効化
  #disable_pw_auth = true
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
| `nic` | - | 基本NIC | `shared` | `shared`(共有セグメント)<br />`[switch_id]`(スイッチのID)<br />`"disconnect"`(接続なし)|eth0の上流NWとの接続方法を指定する。 |
| `additional_nics` | - | 追加NIC | - | リスト(文字列) | 追加で割り当てるNIC。接続するスイッチのID、または空文字を指定する。 |
| `packet_filter_ids`| - | パケットフィルタID | - | リスト(文字列) | NICに適用するパケットフィルタのIDをリストで指定する。リストの先頭からeth0,eth1の順で適用される |
| `icon_id`       | -   | アイコンID         | - | 文字列| - |
| `description` | - | 説明 | - | 文字列 | - |
| `cdrom_id` | - | CDROM(ISOイメージ)ID | - | 文字列 | - |
| `ipaddress`| - | 基本NIC:IPアドレス | - | 文字列 | ディスク修正機能で設定される、IPアドレス [注1](#注1)|
| `gateway`  | - | 基本NIC:ゲートウェイ | - | 文字列 | ディスク修正機能で設定される、ゲートウェイアドレス [注1](#注1) |
| `nw_mask_len` | - | 基本NIC:サブネットマスク長 | - | 文字列 | ディスク修正機能で設定される、サブネットマスク長 [注1](#注1) |
| `hostname`        | -   | ホスト名               | - | 文字列 | ディスク修正機能で設定される、ホスト名 [注1](#注1) |
| `password`        | -   | パスワード               | - | 文字列 | ディスク修正機能で設定される、OS管理者パスワード [注1](#注1)|
| `ssh_key_ids`     | -   | SSH公開鍵ID             | - | リスト(文字列) | ディスク修正機能で設定される、SSH認証用の公開鍵ID [注1](#注1)|
| `disable_pw_auth` | -   | パスワードでの認証無効化   | - | `true`<br />`false` | ディスク修正機能で設定される、SSH接続でのパスワード/チャレンジレスポンス認証の無効化 [注1](#注1)|
| `note_ids`        | -   | スタートアップスクリプトID | - | リスト(文字列) | ディスク修正機能で設置される、スタートアップスクリプトのID [注1](#注1)|
| `private_host_id` | - | 専有ホストID | - | 文字列 | 専有ホストは東京第1ゾーン(tk1a)と石狩第2ゾーン(is1b)でのみ利用可能 |
| `tags` | - | タグ | - | リスト(文字列) | サーバに付与するタグ。@で始まる特殊タグについては[こちら](http://cloud-news.sakura.ad.jp/special-tags/)を参照 |
| `graceful_shutdown_timeout` | - | シャットダウンまでの待ち時間 | - | 数値(秒数) | シャットダウンが必要な場合の通常シャットダウンするまでの待ち時間(指定の時間まで待ってもシャットダウンしない場合は強制シャットダウンされる) |
| `zone` | - | ゾーン | - | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |

#### 注1 ディスク修正機能関連の項目

- サーバにディスクが接続されている場合のみ有効です。
- サーバに接続されたディスクのうち、最初のディスクのみが対象となります。
- OS(ディスクのコピー元アーカイブ)によってはディスク修正機能に対応していない場合があります。
- これらの値は投入専用です。属性においても投入値を表します(さくらのクラウドAPIからは取得できない項目です)。
- IPアドレス/ゲートウェイ/サブネットマスク長については`nic`にスイッチのIDが指定されている場合にのみ設定されます。
  
これらの値をサーバリソース/ディスクリソースの両方に記載した場合の動作は不定です。
混乱を避けるためにいずれか一方にのみ記載するようにしてください。

### 属性

|属性名                    | 名称                     | 補足                                        |
|-------------------------|-------------------------|--------------------------------------------|
| `id`                    | ID                      | -                                          |
| `macaddresses`          | MACアドレス               | MACアドレスのリスト(NICの個数分のリスト)        |
| `dns_servers`           | 基本NIC:DNSサーバ        | eth0の属するセグメントの推奨ネームサーバのリスト|
| `nw_address`            | 基本NIC:ネットワークアドレス | eth0のIPアドレスのネットワークアドレス          |
| `private_host_name`     | 専有ホスト名 | -          |
