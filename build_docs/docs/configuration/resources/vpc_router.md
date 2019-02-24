# VPCルータ(sakuracloud_vpc_router)

---

NOTE: VPCルータの各設定を子リソースとして定義したい場合(Terraform v0.11までの方法)は[VPCルータ(子リソース)](./vpc_router_old)を参照してください。

### 設定例

```hcl
# VPCルータの上流ルータ(プレミアム以上のプランの場合、ルータが必須)
resource "sakuracloud_internet" "router1" {
    name = "myinternet1"
}

# VPCルータ配下に接続するスイッチ
resource sakuracloud_switch "sw" {
  name = "name_before"
}

# VPCルータ本体の定義(スタンダードプランの場合)
#resource sakuracloud_vpc_router "foobar" {
#    name = "vpc_router_setting_test"
#    plan = "standard"
#}

# VPCルータ本体の定義(プレミアム/ハイスペックプランの場合)
resource "sakuracloud_vpc_router" "foobar" {
  name        = "example"
  description = "example"
  tags        = ["tag1" , "tag2"]
  plan        = "premium"

  internet_connection = true

  switch_id  = sakuracloud_internet.router1.switch_id
  vip        = sakuracloud_internet.router1.ipaddresses[0]
  ipaddress1 = sakuracloud_internet.router1.ipaddresses[1]
  ipaddress2 = sakuracloud_internet.router1.ipaddresses[2]
  aliases    = [ sakuracloud_internet.router1.ipaddresses[3] ]
  vrid       = 1

  # プライベートNICの定義(複数定義可能)
  interface {
    switch_id   = sakuracloud_switch.sw.id
    vip         = "192.168.11.1"
    ipaddress   = ["192.168.11.2" , "192.168.11.3"]
    nw_mask_len = 24 
  }

  # ポートフォワード
  port_forwarding {
    protocol        = "udp"
    global_port     = 10022
    private_address = "192.168.11.11"
    private_port    = 22
    description     = "desc"
  }

  # スタティックNAT(プレミアム/ハイスペックプランのみ)
  static_nat {
    global_address  = sakuracloud_internet.router1.ipaddresses[3]
    private_address = "192.168.11.12"
    description     = "desc"
  }

  # ファイアウォール
  firewall {
    vpc_router_interface_index = 1

    direction = "send"
    expressions {
        protocol    = "tcp"
        source_nw   = ""
        source_port = "80"
        dest_nw     = ""
        dest_port   = ""
        allow       = true
        logging     = true
        description = "desc"
    }

    expressions {
        protocol    = "ip"
        source_nw   = ""
        source_port = ""
        dest_nw     = ""
        dest_port   = ""
        allow       = false
        logging     = true
        description = "desc"
    }
  }
  
  # DHCPサーバ
  dhcp_server {
    vpc_router_interface_index = 1

    range_start = "192.168.11.11"
    range_stop  = "192.168.11.20"
    dns_servers = ["8.8.8.8", "8.8.4.4"]
  }

  # DHCP スタティックマッピング
  dhcp_static_mapping {
    ipaddress  = "192.168.11.10"
    macaddress = "aa:bb:cc:aa:bb:cc"
  }


  # L2TP/IPsec
  l2tp {
    pre_shared_secret = "example"
    range_start       = "192.168.11.21"
    range_stop        = "192.168.11.30"
  }

  # PPTP
  pptp {
    range_start = "192.168.11.31"
    range_stop  = "192.168.11.40"
  }

  # リモートアクセスユーザー
  user {
    name     = "username"
    password = "password"
  }
 
  # サイト間VPN 
  site_to_site_vpn {
    peer              = "8.8.8.8"
    remote_id         = "8.8.8.8"
    pre_shared_secret = "example"
    routes            = ["10.0.0.0/8"]
    local_prefix      = ["192.168.21.0/24"]
  }

  # スタティックルート
  static_route {
    prefix   = "172.16.0.0/16"
    next_hop = "192.168.11.99"
  }


```

## `sakuracloud_vpc_router`

VPCルータ本体を表します。

### パラメーター

|パラメーター       |必須  |名称           |初期値     |設定値                         |補足                                          |
|-----------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `name`          | ◯   | ロードバランサ名 | -        | 文字列                         | - |
| `plan`          | -   | プラン          |`standard`| `standard`<br />`premium`<br />`highspec` | - |
| `switch_id`     | △   | スイッチID      | -        | 文字列                         | プランが`premium`、`highspec`の場合必須 |
| `vip`           | △   | IPアドレス1     | -        | 文字列                         | プランが`premium`、`highspec`の場合必須 |
| `ipaddress1`    | △   | IPアドレス1     | -        | 文字列                         | プランが`premium`、`highspec`の場合必須 |
| `ipaddress2`    | △   | IPアドレス2     | -        | 文字列                         | プランが`premium`、`highspec`の場合必須 |
| `vrid`          | △   | VRID           | -        | 数値                          | プランが`premium`、`highspec`の場合必須 |
| `aliases`       | -   | IPエイリアス    | -        | リスト(文字列)                  | プランが`premium`、`highspec`の場合のみ有効 |
| `syslog_host`   | -   | syslog転送先ホスト| -      | 文字列                         | - |
| `internet_connection` | -   | インターネット接続  | `true` | `true`<br />`false`| - |
| `interface` | -   | プライベートNIC | リスト | [interface](#interface)を参照 | - |
| `port_forwarding` | -   | ポートフォワーディング | リスト | [port_forwarding](#port_forwarding)を参照 | - |
| `static_nat` | -   | スタティックNAT | リスト | [static_nat](#static_nat)を参照 | - |
| `firewall` | -   | ファイアウォール | リスト | [firewall](#firewall)を参照 | - |
| `dhcp_server` | -   | DHCPサーバ | リスト | [dhcp_server](#dhcp_server)を参照 | - |
| `dhcp_static_mapping` | -   | DHCPスタティックマッピング | リスト | [dhcp_static_mapping](#dhcp_static_mapping)を参照 | - |
| `l2pt` | -   | L2TP/IPsec | リスト | [l2tp](#l2tp)を参照 | - |
| `pptp` | -   | PPTP | リスト | [pptp](#pptp)を参照 | - |
| `user` | -   | リモートアクセスユーザー | リスト | [user](#user)を参照 | - |
| `site_to_site_vpn` | -   | サイト間VPN | リスト | [site_to_site_vpn](#site_to_site_vpn)を参照 | - |
| `static_route` | -   | スタティックルート | リスト | [static_route](#static_route)を参照 | - |
| `icon_id`       | -   | アイコンID         | - | 文字列| - |
| `description`   | -   | 説明           | -        | 文字列                         | - |
| `tags`          | -   | タグ           | -        | リスト(文字列)                  | - |
| `graceful_shutdown_timeout` | - | シャットダウンまでの待ち時間 | - | 数値(秒数) | シャットダウンが必要な場合の通常シャットダウンするまでの待ち時間(指定の時間まで待ってもシャットダウンしない場合は強制シャットダウンされる) |
| `zone`          | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名          | 名称             | 補足                  |
|---------------|------------------|----------------------|
| `id`            | ID             | -                    |
| `global_address`| グローバルIP     | VPCルータ自身のグローバルIP |

## `interface`

VPCルータが持つプライベートNICを表します。

1台のVPCルータにつき7つまでのプライベートNICを登録できます。

また、プライベートNICの上流には(ルータでは無い)スイッチを接続する必要があります。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-interface/)を参照ください)

|パラメーター          |必須  |名称           |初期値     |設定値                         |補足                                          |
|--------------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `vip`              | △   | VIP            | -        | 文字列                         | プランが`premium`、`highspec`の場合必須 |
| `ipaddress`        | ◯   | IPアドレス      | -        | リスト(文字列)                  | プランが`standard`の場合は1つ、`premium`、`highspec`の場合は2つ指定する |
| `nw_mask_len`      | ◯   | プリフィックス   | -        | 数値(16〜28)                          | - |


## `static_nat`

VPCルータでのスタティックNAT機能を表します。

**このリソースはVPCルータのプランが`premium`、または`highspec`の場合に利用できます。**

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-nat/)を参照ください)

|パラメーター          |必須  |名称           |初期値     |設定値                         |補足                                          |
|--------------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `global_address`          | ◯   | グローバル側IPアドレス  | -        | 文字列                   | VPCルータのIPエイリアスの中のいづれかの値を指定する |
| `private_address`         | ◯   | プライベート側IPアドレス | -        | 文字列                  | - |
| `description`             | -   | 説明             | -        | 文字列                  | - |

## `port_forwarding`

VPCルータでのポートフォワーディング(Reverse NAT)機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-nat/)を参照ください)

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `protocol`                | ◯   | プロトコル             | -        | `tcp`<br />`udp`       | - |
| `global_port`             | ◯   | グローバル側ポート番号   | -        | 数値(1〜65535)                   | - |
| `private_address`         | ◯   | プライベート側IPアドレス | -        | 文字列                          | - |
| `private_port`            | ◯   | プライベート側ポート番号 | -        | 数値(1〜65535)                  | - |
| `description`             | -   | 説明             | -        | 文字列                  | - |

## `firewall`

VPCルータでのファイアウォール機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-firewall/)を参照ください)

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_interface_index`| -   | 対象インターフェースのインデックス| 0        | 数値(`0`-`7`)| - |
| `direction`               | ◯   | 通信方向 | -        | `send`<br />`receive`               | VPCルータから見た通信方向を指定する |
| `expressions`             | ◯   | フィルタルール        | -        | リスト(マップ)           | 詳細は[`expressions`](#expressions)を参照 |

#### `expressions`

|パラメーター     |必須  |名称             |初期値     |設定値                    |補足                                          |
|---------------|:---:|----------------|:--------:|------------------------|----------------------------------------------|
| `protocol`    | ◯   | プロトコル       | -        | `tcp`<br />`udp`<br />`icmp`<br />`ip`| - |
| `source_nw`   | ◯   | 送信元ネットワーク | -       | `xxx.xxx.xxx.xxx`(IP)<br />`xxx.xxx.xxx.xxx/nn`(ネットワーク)<br />`xxx.xxx.xxx.xxx/yyy.yyy.yyy.yyy`(アドレス範囲)  | 空欄の場合はANY |
| `source_port` | ◯   | 送信元ポート      | -       | `0`〜`65535`の整数<br />`xx-yy`(範囲指定)<br />`0xPPPP/0xMMMM`(16進範囲指定) | 空欄の場合はANY |
| `dest_nw`     | ◯   | 送信元ネットワーク | -       | `xxx.xxx.xxx.xxx`(IP)<br />`xxx.xxx.xxx.xxx/nn`(ネットワーク)<br />`xxx.xxx.xxx.xxx/yyy.yyy.yyy.yyy`(アドレス範囲)  | 空欄の場合はANY |
| `dest_port`   | ◯   | 宛先ポート        | -        | `0`〜`65535`の整数<br />`xx-yy`(範囲指定)<br />`0xPPPP/0xMMMM`(16進範囲指定) | 空欄の場合はANY |
| `allow`       | ◯   | アクション        | -        | `true`<br />`false` | `true`の場合ALLOW動作<br />`false`の場合DENY動作 |
| `logging`     | -   | ログ出力         | -        | `true`<br />`false`    | - |
| `description` | -   | 説明             | -        | 文字列                  | - |

## `dhcp_server`

VPCルータでのDHCPサーバ機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-dhcp/)を参照ください)

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_interface_index` | ◯   | プライベートNIC 番号   | -        | 数値                   | - |
| `range_start`                | ◯   | 動的割り当て範囲(開始) | -        | 文字列                          | - |
| `range_stop`                 | ◯   | 動的割り当て範囲(終了) | -        | 文字列                          | - |
| `dns_servers`                | -   | DNSサーバーIPアドレス | -        | リスト(文字列)                   | 省略した場合はゾーンごとのデフォルトDNSサーバが割り当てられる |

## `dhcp_static_mapping`

VPCルータでのDHCPスタティック割当機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-dhcp/)を参照ください)

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `ipaddress`                 | ◯   | IPアドレス | -        | 文字列                          | - |
| `macaddress`                | ◯   | MACアドレス | -        | 文字列                          | 英字は小文字で入力する |

## `pptp`

VPCルータでのPPTPサーバ機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-remoteaccess/)を参照ください)

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `range_start`                | ◯   | 動的割り当て範囲(開始) | -        | 文字列                          | - |
| `range_stop`                 | ◯   | 動的割り当て範囲(終了) | -        | 文字列                          | - |

## `l2tp`

VPCルータでのL2TP/IPSecサーバ機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-remoteaccess/)を参照ください)

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `pre_shared_secret`          | ◯   | 事前共有シークレット   | -        | 文字列                          | - |
| `range_start`                | ◯   | 動的割り当て範囲(開始) | -        | 文字列                          | - |
| `range_stop`                 | ◯   | 動的割り当て範囲(終了) | -        | 文字列                          | - |

## `user`

VPCルータでのリモートユーザーを表します。

このリソースは100個まで指定可能です。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-remoteaccess/)を参照ください)

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `name`                    | ◯   | ユーザー名 | -        | 文字列                          | - |
| `password`                | ◯   | パスワード | -        | 文字列                          | - |

## `site_to_site_vpn`

VPCルータでのサイト間VPNを表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-site-to-site-vpn/)を参照ください)

|パラメーター            |必須  |名称                 |初期値     |設定値                         |補足                                          |
|----------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `peer`               | ◯   | 対向IPアドレス | -        | 文字列                          | - |
| `remote_id`          | ◯   | 対向ID | -        | 文字列                          | - |
| `pre_shared_secret`  | ◯   | 事前共有シークレット | -        | 文字列                          | - |
| `routes`             | ◯   | 対向Prefix | -        | リスト(文字列)                          | - |
| `local_prefix`       | ◯   | ローカルPrefix | -        | リスト(文字列)                          | - |

### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                          | ID                    | -                    |
| `esp_authentication_protocol` | -            | -               |
| `esp_dh_group`                | -            | -               |
| `esp_encryption_protocol`     | -            | -               |
| `esp_lifetime`                | -            | -               |
| `esp_mode`                    | -            | -               |
| `esp_perfect_forward_secrecy` | -            | -               |
| `ike_authentication_protocol` | -            | -               |
| `ike_encryption_protocol`     | -            | -               |
| `ike_lifetime`                | -            | -               |
| `ike_mode`                    | -            | -               |
| `ike_perfect_forward_secrecy` | -            | -               |
| `ike_pre_shared_secret`       | -            | -               |
| `peer_id`                     | -            | -               |
| `peer_inside_networks`        | -            | -               |
| `peer_outside_ipaddress`      | -            | -               |
| `vpc_router_inside_networks`  | -            | -               |
| `vpc_router_outside_ipaddress`| -            | -               |

## `static_route`

VPCルータでのスタティックルート機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-static-route/)を参照ください)

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `prefix`                    | ◯   | プリフィックス | -        | 文字列                          | - |
| `next_hop`                  | ◯   | ネクストホップ | -        | 文字列                          | - |
