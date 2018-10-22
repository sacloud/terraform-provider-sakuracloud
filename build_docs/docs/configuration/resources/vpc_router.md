# VPCルータ(sakuracloud_vpc_router)

---

### 設定例

```hcl
# VPCルータの上流ルータ(プレミアム以上のプランの場合、ルータが必須)
resource "sakuracloud_internet" "router1" {
  name = "myinternet1"
}

# VPCルータ配下に接続するスイッチ
resource "sakuracloud_switch" "sw01" {
  name = "sw01"
}

# VPCルータ本体の定義(プレミアム/ハイスペックプランの場合)
resource "sakuracloud_vpc_router" "foobar" {
  name                = "vpc_router_sample"
  plan                = "premium"
  switch_id           = sakuracloud_internet.router1.switch_id             # 上流のスイッチID
  vip                 = sakuracloud_internet.router1.ipaddresses[0]        # VIP
  ipaddress1          = sakuracloud_internet.router1.ipaddresses[1]        # 実IP1
  ipaddress2          = sakuracloud_internet.router1.ipaddresses[2]        # 実IP2
  aliases             = [sakuracloud_internet.router1.ipaddresses[3]]      # IPエイリアス
  vrid                = 1
  syslog_host         = "192.168.11.1" # syslog転送先ホスト
  internet_connection = true           # インターネット接続 有効/無効
}

# VPCルータ本体の定義(スタンダードプランの場合)
#resource sakuracloud_vpc_router "foobar" {
#    name = "vpc_router_setting_test"
#    plan = "standard"
#}

# VPCルータ配下のプライベートNIC(プレミアム/ハイスペックプランの場合)
resource "sakuracloud_vpc_router_interface" "eth1" {
  vpc_router_id = sakuracloud_vpc_router.foobar.id

  index       = 1                                # NICのインデックス(1〜7)
  switch_id   = sakuracloud_switch.sw01.id       # スイッチのID
  vip         = "192.168.11.1"                   # VIP
  ipaddress   = ["192.168.11.2", "192.168.11.3"] # 実IPリスト
  nw_mask_len = 24
}

# VPCルータ配下のプライベートNIC(スタンダードプランの場合)
#resource "sakuracloud_vpc_router_interface" "eth1"{
#    vpc_router_id = sakuracloud_vpc_router.foobar.id
#
#    index = 1                                       # NICのインデックス(1〜7)
#    switch_id = sakuracloud_switch.sw01.id          # スイッチのID
#    ipaddress = ["192.168.11.2"]                    # 実IPリスト
#    nw_mask_len = 24
#}

# StaticNAT機能(プレミアム/ハイスペックプランの場合のみ利用可能)
resource "sakuracloud_vpc_router_static_nat" "staticNAT1" {
  vpc_router_id           = sakuracloud_vpc_router.foobar.id
  vpc_router_interface_id = sakuracloud_vpc_router_interface.eth1.id # 対象プライベートIPが属するNICのID

  global_address  = sakuracloud_internet.router1.ipaddresses[3] # グローバル側IPアドレス(VPCルータ本体に割り当てたIPエイリアス)
  private_address = "192.168.11.11"                             # プライベート側アドレス
}

# ポートフォワーディング
resource "sakuracloud_vpc_router_port_forwarding" "forward1" {
  vpc_router_id           = sakuracloud_vpc_router.foobar.id
  vpc_router_interface_id = sakuracloud_vpc_router_interface.eth1.id # 対象プライベートIPが属するNICのID

  protocol        = "tcp"           # プロトコル(tcp/udp)
  global_port     = 10022           # グローバル側ポート番号
  private_address = "192.168.11.11" # プライベートIPアドレス
  private_port    = 22              # プライベート側ポート番号
}

# ファイアウォール(VPC内部から外部への通信)
resource "sakuracloud_vpc_router_firewall" "send_fw" {
  vpc_router_id = sakuracloud_vpc_router.foobar.id

  # vpc_router_interface_index = 0 # 対象インターフェースのインデックス(グローバル含む)
  direction = "send"

  # VPC内部のWebサーバから外部への応答パケットの許可
  expressions {
    protocol    = "tcp"
    source_nw   = ""
    source_port = "80"
    dest_nw     = ""
    dest_port   = ""
    allow       = true
  }

  # 全拒否(暗黙Deny)
  expressions {
    protocol    = "ip"
    source_nw   = ""
    source_port = ""
    dest_nw     = ""
    dest_port   = ""
    allow       = false
  }
}

# ファイアウォール(VPC外部から内部への通信)
resource "sakuracloud_vpc_router_firewall" "receive_fw" {
  vpc_router_id = sakuracloud_vpc_router.foobar.id

  # vpc_router_interface_index = 0 # 対象インターフェースのインデックス(グローバル含む)
  direction = "receive"

  # VPC内部のWebサーバへのパケットを許可
  expressions {
    protocol    = "tcp"
    source_nw   = ""
    source_port = ""
    dest_nw     = ""
    dest_port   = "80"
    allow       = true
  }

  # 全拒否(暗黙Deny)
  expressions {
    protocol    = "ip"
    source_nw   = ""
    source_port = ""
    dest_nw     = ""
    dest_port   = ""
    allow       = false
  }
}

# DHCPサーバ機能
resource "sakuracloud_vpc_router_dhcp_server" "dhcp" {
  vpc_router_id              = sakuracloud_vpc_router.foobar.id
  vpc_router_interface_index = sakuracloud_vpc_router_interface.eth1.index # 対象プライベートIPが属するNICのインデックス

  range_start = "192.168.11.151" # IPアドレス動的割り当て範囲(開始)
  range_stop  = "192.168.11.200" # IPアドレス動的割り当て範囲(終了)
  
  # dns_servers = ["8.8.4.4", "8.8.8.8"] # 配布するDNSサーバIPアドレスのリスト
}


# DHCPスタティック割り当て
resource "sakuracloud_vpc_router_dhcp_static_mapping" "dhcp_map" {
  vpc_router_id             = sakuracloud_vpc_router.foobar.id
  vpc_router_dhcp_server_id = sakuracloud_vpc_router_dhcp_server.dhcp.id # DHCPサーバリソースのID

  macaddress = "aa:bb:cc:aa:bb:cc" # 対象MACアドレス
  ipaddress  = "192.168.11.20"     # 割りあてるIPアドレス
}

# リモートアクセス:PPTPサーバ機能
resource "sakuracloud_vpc_router_pptp" "pptp" {
  vpc_router_id           = sakuracloud_vpc_router.foobar.id
  vpc_router_interface_id = sakuracloud_vpc_router_interface.eth1.id

  range_start = "192.168.11.101" # IPアドレス動的割り当て範囲(開始)
  range_stop  = "192.168.11.150" # IPアドレス動的割り当て範囲(終了)
}

# リモートアクセス:L2TP/IPSecサーバ機能
resource "sakuracloud_vpc_router_l2tp" "l2tp" {
  vpc_router_id           = sakuracloud_vpc_router.foobar.id
  vpc_router_interface_id = sakuracloud_vpc_router_interface.eth1.id

  pre_shared_secret = "hogehoge"       # 事前共有シークレット
  range_start       = "192.168.11.51"  # IPアドレス動的割り当て範囲(開始)
  range_stop        = "192.168.11.100" # IPアドレス動的割り当て範囲(終了)
}

# リモートユーザーアカウント
resource "sakuracloud_vpc_router_user" "user1" {
  vpc_router_id = sakuracloud_vpc_router.foobar.id

  name     = "username" # ユーザー名
  password = "password" # パスワード
}

# サイト間VPN
resource "sakuracloud_vpc_router_site_to_site_vpn" "s2s" {
  vpc_router_id     = sakuracloud_vpc_router.foobar.id
  peer              = "8.8.8.8"
  remote_id         = "8.8.8.8"
  pre_shared_secret = "presharedsecret"
  routes            = ["10.0.0.0/8"]
  local_prefix      = ["192.168.21.0/24"]
}

# スタティックルート
resource "sakuracloud_vpc_router_static_route" "route1" {
  vpc_router_id           = sakuracloud_vpc_router.foobar.id
  vpc_router_interface_id = sakuracloud_vpc_router_interface.eth1.id
  prefix                  = "172.16.0.0/16"
  next_hop                = "192.168.11.99"
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

## `sakuracloud_vpc_router_interface`

VPCルータが持つプライベートNICを表します。

1台のVPCルータにつき7つまでのプライベートNICを登録できます。

また、プライベートNICの上流には(ルータでは無い)スイッチを接続する必要があります。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-interface/)を参照ください)

### パラメーター

|パラメーター          |必須  |名称           |初期値     |設定値                         |補足                                          |
|--------------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`    | ◯   | VPCルータID   | -        | 文字列                         | - |
| `index`            | ◯   | NIC番号        | -        | 数値(1〜7)                     | - |
| `vip`              | △   | VIP            | -        | 文字列                         | プランが`premium`、`highspec`の場合必須 |
| `ipaddress`        | ◯   | IPアドレス      | -        | リスト(文字列)                  | プランが`standard`の場合は1つ、`premium`、`highspec`の場合は2つ指定する |
| `nw_mask_len`      | ◯   | プリフィックス   | -        | 数値(16〜28)                          | - |
| `graceful_shutdown_timeout` | - | シャットダウンまでの待ち時間 | - | 数値(秒数) | シャットダウンが必要な場合の通常シャットダウンするまでの待ち時間(指定の時間まで待ってもシャットダウンしない場合は強制シャットダウンされる) |
| `zone`          | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名          | 名称             | 補足                  |
|---------------|------------------|----------------------|
| `id`            | ID             | -                    |


## `sakuracloud_vpc_router_static_nat`

VPCルータでのスタティックNAT機能を表します。

**このリソースはVPCルータのプランが`premium`、または`highspec`の場合に利用できます。**

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-nat/)を参照ください)

### パラメーター

|パラメーター          |必須  |名称           |初期値     |設定値                         |補足                                          |
|--------------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`           | ◯   | VPCルータID         | -        | 文字列                   | - |
| `vpc_router_interface_id` | ◯   | プライベートNIC ID     | -        | 文字列                   | - |
| `global_address`          | ◯   | グローバル側IPアドレス  | -        | 文字列                   | VPCルータのIPエイリアスの中のいづれかの値を指定する |
| `private_address`         | ◯   | プライベート側IPアドレス | -        | 文字列                  | - |
| `description`             | -   | 説明             | -        | 文字列                  | - |
| `zone`          | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID             | -                    |

## `sakuracloud_vpc_router_port_forwarding`

VPCルータでのポートフォワーディング(Reverse NAT)機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-nat/)を参照ください)

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`           | ◯   | VPCルータID         | -        | 文字列                   | - |
| `vpc_router_interface_id` | ◯   | プライベートNIC ID     | -        | 文字列                   | - |
| `protocol`                | ◯   | プロトコル             | -        | `tcp`<br />`udp`       | - |
| `global_port`             | ◯   | グローバル側ポート番号   | -        | 数値(1〜65535)                   | - |
| `private_address`         | ◯   | プライベート側IPアドレス | -        | 文字列                          | - |
| `private_port`            | ◯   | プライベート側ポート番号 | -        | 数値(1〜65535)                  | - |
| `description`             | -   | 説明             | -        | 文字列                  | - |
| `zone`          | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |

### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |

## `sakuracloud_vpc_router_firewall`

VPCルータでのファイアウォール機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-firewall/)を参照ください)

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`           | ◯   | VPCルータID         | -        | 文字列                   | - |
| `vpc_router_interface_index`| -   | 対象インターフェースのインデックス| 0        | 数値(`0`-`7`)| - |
| `direction`               | ◯   | 通信方向 | -        | `send`<br />`receive`               | VPCルータ内から見た通信方向を指定する |
| `expressions`             | ◯   | フィルタルール        | -        | リスト(マップ)           | 詳細は[`expressions`](#expressions)を参照 |
| `zone`          | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |

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


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |

## `sakuracloud_vpc_router_dhcp_server`

VPCルータでのDHCPサーバ機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-dhcp/)を参照ください)

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`              | ◯   | VPCルータID         | -        | 文字列                   | - |
| `vpc_router_interface_index` | ◯   | プライベートNIC 番号   | -        | 数値                   | - |
| `range_start`                | ◯   | 動的割り当て範囲(開始) | -        | 文字列                          | - |
| `range_stop`                 | ◯   | 動的割り当て範囲(終了) | -        | 文字列                          | - |
| `dns_servers`                | -   | DNSサーバーIPアドレス | -        | リスト(文字列)                   | 省略した場合はゾーンごとのデフォルトDNSサーバが割り当てられる |
| `zone`          | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |

## `sakuracloud_vpc_router_dhcp_static_mapping`

VPCルータでのDHCPスタティック割当機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-dhcp/)を参照ください)

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`             | ◯   | VPCルータID         | -        | 文字列                   | - |
| `vpc_router_dhcp_server_id` | ◯   | DHCPサーバID       | -        | 文字列                   | - |
| `ipaddress`                 | ◯   | IPアドレス | -        | 文字列                          | - |
| `macaddress`                | ◯   | MACアドレス | -        | 文字列                          | 英字は小文字で入力する |
| `zone`          | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |

## `sakuracloud_vpc_router_pptp`

VPCルータでのPPTPサーバ機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-remoteaccess/)を参照ください)

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`              | ◯   | VPCルータID         | -        | 文字列                   | - |
| `vpc_router_interface_id`    | ◯   | プライベートNIC ID     | -        | 文字列                   | - |
| `range_start`                | ◯   | 動的割り当て範囲(開始) | -        | 文字列                          | - |
| `range_stop`                 | ◯   | 動的割り当て範囲(終了) | -        | 文字列                          | - |
| `zone`          | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |

## `sakuracloud_vpc_router_l2tp`

VPCルータでのL2TP/IPSecサーバ機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-remoteaccess/)を参照ください)

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`              | ◯   | VPCルータID         | -        | 文字列                   | - |
| `vpc_router_interface_id`    | ◯   | プライベートNIC ID     | -        | 文字列                   | - |
| `pre_shared_secret`          | ◯   | 事前共有シークレット   | -        | 文字列                          | - |
| `range_start`                | ◯   | 動的割り当て範囲(開始) | -        | 文字列                          | - |
| `range_stop`                 | ◯   | 動的割り当て範囲(終了) | -        | 文字列                          | - |
| `zone`          | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |

## `sakuracloud_vpc_router_user`

VPCルータでのリモートユーザーを表します。

このリソースは100個まで指定することが可能です。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-remoteaccess/)を参照ください)

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`           | ◯   | VPCルータID         | -        | 文字列                   | - |
| `name`                    | ◯   | ユーザー名 | -        | 文字列                          | - |
| `password`                | ◯   | パスワード | -        | 文字列                          | - |
| `zone`          | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |


## `sakuracloud_vpc_router_site_to_site_vpn`

VPCルータでのサイト間VPNを表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-site-to-site-vpn/)を参照ください)

### パラメーター

|パラメーター            |必須  |名称                 |初期値     |設定値                         |補足                                          |
|----------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`      | ◯   | VPCルータID         | -        | 文字列                   | - |
| `peer`               | ◯   | 対向IPアドレス | -        | 文字列                          | - |
| `remote_id`          | ◯   | 対向ID | -        | 文字列                          | - |
| `pre_shared_secret`  | ◯   | 事前共有シークレット | -        | 文字列                          | - |
| `routes`             | ◯   | 対向Prefix | -        | リスト(文字列)                          | - |
| `local_prefix`       | ◯   | ローカルPrefix | -        | リスト(文字列)                          | - |
| `zone`          | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


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


## `sakuracloud_vpc_router_static_route`

VPCルータでのスタティックルート機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-static-route/)を参照ください)

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`              | ◯   | VPCルータID         | -        | 文字列                   | - |
| `vpc_router_interface_id`    | ◯   | プライベートNIC ID     | -        | 文字列                   | - |
| `prefix`                    | ◯   | プリフィックス | -        | 文字列                          | - |
| `next_hop`                  | ◯   | ネクストホップ | -        | 文字列                          | - |
| `zone`          | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |
