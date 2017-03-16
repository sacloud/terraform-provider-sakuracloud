# VPCルーター(sakuracloud_vpc_router)

### 設定例

```tf:VPCルーター設定サンプル.tf

# VPCルーターの上流ルーター(プレミアム以上のプランの場合、ルーターが必須)
resource "sakuracloud_internet" "router1" {
    name = "myinternet1"
}

# VPCルーター配下に接続するスイッチ
resource "sakuracloud_switch" "sw01"{
    name = "sw01"
}

# VPCルーター本体の定義(プレミアム/ハイスペックプランの場合)
resource "sakuracloud_vpc_router" "foobar" {
    name = "vpc_router_sample"
    plan = "premium"
    switch_id = "${sakuracloud_internet.router1.switch_id}"          # 上流のスイッチID
    vip = "${sakuracloud_internet.router1.nw_ipaddresses.0}"         # VIP
    ipaddress1 = "${sakuracloud_internet.router1.nw_ipaddresses.1}"  # 実IP1
    ipaddress2 = "${sakuracloud_internet.router1.nw_ipaddresses.2}"  # 実IP2
    aliases = ["${sakuracloud_internet.router1.nw_ipaddresses.3}"]   # IPエイリアス
    VRID = 1
    syslog_host = "192.168.11.1"                                     # syslog転送先ホスト
}

# VPCルーター本体の定義(スタンダードプランの場合)
#resource "sakuracloud_vpc_router" "foobar" {
#    name = "vpc_router_setting_test"
#    plan = "standard"
#}

# VPCルーター配下のプライベートNIC(プレミアム/ハイスペックプランの場合)
resource "sakuracloud_vpc_router_interface" "eth1"{
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"

    index = 1                                       # NICのインデックス(1〜7)
    switch_id = "${sakuracloud_switch.sw01.id}"     # スイッチのID
    vip = "192.168.11.1"                            # VIP
    ipaddress = ["192.168.11.2" , "192.168.11.3"]   # 実IPリスト
    nw_mask_len = 24
}

# VPCルーター配下のプライベートNIC(スタンダードプランの場合)
#resource "sakuracloud_vpc_router_interface" "eth1"{
#    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
#
#    index = 1                                       # NICのインデックス(1〜7)
#    switch_id = "${sakuracloud_switch.sw01.id}"     # スイッチのID
#    ipaddress = ["192.168.11.2"]                    # 実IPリスト
#    nw_mask_len = 24
#}


# StaticNAT機能(プレミアム/ハイスペックプランの場合のみ利用可能)
resource "sakuracloud_vpc_router_static_nat" "staticNAT1" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_id = "${sakuracloud_vpc_router_interface.eth1.id}" # 対象プライベートIPが属するNICのID

    global_address = "${sakuracloud_internet.router1.nw_ipaddresses.3}"     # グローバル側IPアドレス(VPCルーター本体に割り当てたIPエイリアス)
    private_address = "192.168.11.11"                                       # プライベート側アドレス
}

# ポートフォワーディング
resource "sakuracloud_vpc_router_port_forwarding" "forward1" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_id = "${sakuracloud_vpc_router_interface.eth1.id}" # 対象プライベートIPが属するNICのID

    protocol = "tcp"                  # プロトコル(tcp/udp)
    global_port = 10022               # グローバル側ポート番号
    private_address = "192.168.11.11" # プライベートIPアドレス
    private_port = 22                 # プライベート側ポート番号
}

# ファイアウォール(VPC内部から外部への通信)
resource "sakuracloud_vpc_router_firewall" "send_fw" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    direction = "send"

    # VPC内部のWebサーバーから外部への応答パケットの許可
    expressions = {
        protocol = "tcp"
        source_nw = ""
        source_port = "80"
        dest_nw = ""
        dest_port = ""
        allow = true
    }

    # 全拒否(暗黙Deny)
    expressions = {
        protocol = "ip"
        source_nw = ""
        source_port = ""
        dest_nw = ""
        dest_port = ""
        allow = false
    }
}

# ファイアウォール(VPC外部から内部への通信)
resource "sakuracloud_vpc_router_firewall" "receive_fw" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    direction = "receive"

    # VPC内部のWebサーバーへのパケットを許可
    expressions = {
        protocol = "tcp"
        source_nw = ""
        source_port = ""
        dest_nw = ""
        dest_port = "80"
        allow = true
    }

    # 全拒否(暗黙Deny)
    expressions = {
        protocol = "ip"
        source_nw = ""
        source_port = ""
        dest_nw = ""
        dest_port = ""
        allow = false
    }
}

# DHCPサーバー機能
resource "sakuracloud_vpc_router_dhcp_server" "dhcp" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_index = "${sakuracloud_vpc_router_interface.eth1.index}" # 対象プライベートIPが属するNICのインデックス

    range_start = "192.168.11.151" # IPアドレス動的割り当て範囲(開始)
    range_stop = "192.168.11.200"  # IPアドレス動的割り当て範囲(終了)
}

# DHCPスタティック割り当て
resource "sakuracloud_vpc_router_dhcp_static_mapping" "dhcp_map" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_dhcp_server_id = "${sakuracloud_vpc_router_dhcp_server.dhcp.id}" # DHCPサーバーリソースのID

    macaddress = "aa:bb:cc:aa:bb:cc"  # 対象MACアドレス
    ipaddress = "192.168.11.20"       # 割りあてるIPアドレス
}

# リモートアクセス:PPTPサーバー機能
resource "sakuracloud_vpc_router_pptp" "pptp"{
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_id = "${sakuracloud_vpc_router_interface.eth1.id}"

    range_start = "192.168.11.101" # IPアドレス動的割り当て範囲(開始)
    range_stop = "192.168.11.150"  # IPアドレス動的割り当て範囲(終了)
}

# リモートアクセス:L2TP/IPSecサーバー機能
resource "sakuracloud_vpc_router_l2tp" "l2tp" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_id = "${sakuracloud_vpc_router_interface.eth1.id}"

    pre_shared_secret = "hogehoge" # 事前共有シークレット
    range_start = "192.168.11.51"  # IPアドレス動的割り当て範囲(開始)
    range_stop = "192.168.11.100"  # IPアドレス動的割り当て範囲(終了)

}

# リモートユーザーアカウント
resource "sakuracloud_vpc_router_user" "user1" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"

    name = "username"     # ユーザー名
    password = "password" # パスワード
}

# サイト間VPN
resource "sakuracloud_vpc_router_site_to_site_vpn" "s2s" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    peer = "8.8.8.8"
    remote_id = "8.8.8.8"
    pre_shared_secret = "presharedsecret"
    routes = ["10.0.0.0/8"]
    local_prefix = ["192.168.21.0/24"]
}

# スタティックルート
resource "sakuracloud_vpc_router_static_route" "route1" {
    vpc_router_id = "${sakuracloud_vpc_router.foobar.id}"
    vpc_router_interface_id = "${sakuracloud_vpc_router_interface.eth1.id}"
    prefix = "172.16.0.0/16"
    next_hop = "192.168.11.99"
}

```

## `sakuracloud_vpc_router`

VPCルーター本体を表します。

### パラメーター

|パラメーター       |必須  |名称           |初期値     |設定値                         |補足                                          |
|-----------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `name`          | ◯   | ロードバランサ名 | -        | 文字列                         | - |
| `plan`          | -   | プラン          |`standard`| `standard`<br />`premium`<br />`highspec` | - |
| `switch_id`     | △   | スイッチID      | -        | 文字列                         | プランが`premium`、`highspec`の場合必須 |
| `vip`           | △   | IPアドレス1     | -        | 文字列                         | プランが`premium`、`highspec`の場合必須 |
| `ipaddress1`    | △   | IPアドレス1     | -        | 文字列                         | プランが`premium`、`highspec`の場合必須 |
| `ipaddress2`    | △   | IPアドレス2     | -        | 文字列                         | プランが`premium`、`highspec`の場合必須 |
| `VRID`          | △   | VRID           | -        | 数値                          | プランが`premium`、`highspec`の場合必須 |
| `aliases`       | -   | IPエイリアス    | -        | リスト(文字列)                  | プランが`premium`、`highspec`の場合のみ有効 |
| `syslog_host`   | -   | syslog転送先ホスト| -      | 文字列                         | - |
| `description`   | -   | 説明           | -        | 文字列                         | - |
| `tags`          | -   | タグ           | -        | リスト(文字列)                  | - |
| `zone`          | -   | ゾーン          | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名          | 名称             | 補足                  |
|---------------|------------------|----------------------|
| `id`            | ID             | -                    |
| `name`          | VPCルーター名   | -                    |
| `plan`          | プラン          | -                    |
| `switch_id`     | スイッチID      | -                    |
| `vip`           | VIP            | -                     |
| `ipaddress1`    | IPアドレス1      | -                    |
| `ipaddress2`    | IPアドレス2      | -                    |
| `VRID`          | VRID           | -                     |
| `aliases`       | IPエイリアス      | -                   |
| `syslog_host`   | syslog転送先ホスト | -                   |
| `description`   | 説明             | -                   |
| `tags`          | タグ             | -                  |
| `zone`          | ゾーン           | -                   |
| `global_address`| グローバルIP     | VPCルーター自身のグローバルIP |

## `sakuracloud_vpc_router_interface`

VPCルーターが持つプライベートNICを表します。

1台のVPCルーターにつき7つまでのプライベートNICを登録できます。

また、プライベートNICの上流には(ルーターでは無い)スイッチを接続する必要があります。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-interface/)を参照ください。)

### パラメーター

|パラメーター          |必須  |名称           |初期値     |設定値                         |補足                                          |
|--------------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`    | ◯   | VPCルーターID   | -        | 文字列                         | - |
| `index`            | ◯   | NIC番号        | -        | 数値(1〜7)                     | - |
| `vip`              | △   | VIP            | -        | 文字列                         | プランが`premium`、`highspec`の場合必須 |
| `ipaddress`        | ◯   | IPアドレス      | -        | リスト(文字列)                  | プランが`standard`の場合は1つ、`premium`、`highspec`の場合は2つ指定する |
| `nw_mask_len`      | ◯   | プリフィックス   | -        | 数値(16〜28)                          | - |
| `zone`             | -   | ゾーン          | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名          | 名称             | 補足                  |
|---------------|------------------|----------------------|
| `id`            | ID             | -                    |
| `vpc_router_id` | VPCルーターID   | -                    |
| `index`         | NIC番号        | -                    |
| `vip`           | VIP            | -                    |
| `ipaddress`     | IPアドレス      | -                     |
| `nw_mask_len`   | プリフィックス    | -                    |
| `zone`          | ゾーン           | -                   |


## `sakuracloud_vpc_router_static_nat`

VPCルーターでのスタティックNAT機能を表します。

**このリソースはVPCルーターのプランが`premium`、または`highspec`の場合に利用できます。**

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-nat/)を参照ください。)

### パラメーター

|パラメーター          |必須  |名称           |初期値     |設定値                         |補足                                          |
|--------------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`           | ◯   | VPCルーターID         | -        | 文字列                   | - |
| `vpc_router_interface_id` | ◯   | プライベートNIC ID     | -        | 文字列                   | - |
| `global_address`          | ◯   | グローバル側IPアドレス  | -        | 文字列                   | VPCルーターのIPエイリアスの中のいづれかの値を指定する |
| `private_address`         | ◯   | プライベート側IPアドレス | -        | 文字列                  | - |
| `description`             | -   | 説明             | -        | 文字列                  | - |
| `zone`                    | -   | ゾーン          | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID             | -                    |
| `vpc_router_id`          | VPCルーターID   | -                    |
| `vpc_router_interface_id`| プライベートNIC ID | -                    |
| `global_address`         | グローバル側IPアドレス            | -                    |
| `private_address`        | プライベート側IPアドレス      | -                     |
| `description`            | 説明      | -                     |
| `zone`                   | ゾーン           | -                   |


## `sakuracloud_vpc_router_port_forwarding`

VPCルーターでのポートフォワーディング(Reverse NAT)機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-nat/)を参照ください。)

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`           | ◯   | VPCルーターID         | -        | 文字列                   | - |
| `vpc_router_interface_id` | ◯   | プライベートNIC ID     | -        | 文字列                   | - |
| `protocol`                | ◯   | プロトコル             | -        | `tcp`<br />`udp`       | - |
| `global_port`             | ◯   | グローバル側ポート番号   | -        | 数値(1〜65535)                   | - |
| `private_address`         | ◯   | プライベート側IPアドレス | -        | 文字列                          | - |
| `private_port`            | ◯   | プライベート側ポート番号 | -        | 数値(1〜65535)                  | - |
| `description`             | -   | 説明             | -        | 文字列                  | - |
| `zone`                    | -   | ゾーン                 | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |
| `vpc_router_id`          | VPCルーターID          | -                    |
| `vpc_router_interface_id`| プライベートNIC ID     | -                    |
| `protocol`               | プロトコル              | -                    |
| `global_port`            | グローバル側ポート番号   | -                    |
| `private_address`        | プライベート側IPアドレス  | -                     |
| `private_port`           | プライベート側ポート番号  | -                     |
| `description`            | 説明                   | -                     |
| `zone`                   | ゾーン                 | -                   |

## `sakuracloud_vpc_router_firewall`

VPCルーターでのファイアウォール機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-firewall/)を参照ください。)

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`           | ◯   | VPCルーターID         | -        | 文字列                   | - |
| `direction`               | ◯   | 通信方向 | -        | `send`<br />`receive`               | VPCルーター内から見た通信方向を指定する |
| `expressions`             | ◯   | フィルタルール        | -        | リスト(マップ)           | 詳細は[`expressions`](#expressions)を参照 |
| `zone`                    | -   | ゾーン                 | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |

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
| `vpc_router_id`          | VPCルーターID          | -                    |
| `private_port`           | プライベート側ポート番号  | -                     |
| `expressions`            | フィルタルール    | [`expressions`](#expressions)のリスト |
| `zone`                   | ゾーン                 | -                   |

## `sakuracloud_vpc_router_dhcp_server`

VPCルーターでのDHCPサーバー機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-dhcp/)を参照ください。)

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`              | ◯   | VPCルーターID         | -        | 文字列                   | - |
| `vpc_router_interface_index` | ◯   | プライベートNIC 番号   | -        | 数値                   | - |
| `range_start`                | ◯   | 動的割り当て範囲(開始) | -        | 文字列                          | - |
| `range_stop`                 | ◯   | 動的割り当て範囲(終了) | -        | 文字列                          | - |
| `zone`                    | -   | ゾーン                 | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |
| `vpc_router_id`          | VPCルーターID          | -                    |
| `vpc_router_interface_index`| プライベートNIC 番号     | -                    |
| `range_start`            | 動的割り当て範囲(開始) | -                     |
| `range_stop`             | 動的割り当て範囲(終了)  | -                     |
| `zone`                   | ゾーン                 | -                   |

## `sakuracloud_vpc_router_dhcp_static_mapping`

VPCルーターでのDHCPスタティック割当機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-dhcp/)を参照ください。)

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`             | ◯   | VPCルーターID         | -        | 文字列                   | - |
| `vpc_router_dhcp_server_id` | ◯   | DHCPサーバーID       | -        | 数値                   | - |
| `ipaddress`                 | ◯   | IPアドレス | -        | 文字列                          | - |
| `macaddress`                | ◯   | MACアドレス | -        | 文字列                          | 英字は小文字で入力する |
| `zone`                      | -   | ゾーン                 | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |
| `vpc_router_id`          | VPCルーターID          | -                    |
| `vpc_router_dhcp_server_id`| DHCPサーバーID     | -                    |
| `ipaddress`              | IPアドレス | -                     |
| `macaddress`             | MACアドレス  | -                     |
| `zone`                   | ゾーン                 | -                   |

## `sakuracloud_vpc_router_pptp`

VPCルーターでのPPTPサーバー機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-remoteaccess/)を参照ください。)

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`              | ◯   | VPCルーターID         | -        | 文字列                   | - |
| `vpc_router_interface_id`    | ◯   | プライベートNIC ID     | -        | 文字列                   | - |
| `range_start`                | ◯   | 動的割り当て範囲(開始) | -        | 文字列                          | - |
| `range_stop`                 | ◯   | 動的割り当て範囲(終了) | -        | 文字列                          | - |
| `zone`                    | -   | ゾーン                 | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |
| `vpc_router_id`          | VPCルーターID          | -                    |
| `vpc_router_interface_id`| プライベートNIC ID     | -                    |
| `range_start`            | 動的割り当て範囲(開始) | -                     |
| `range_stop`             | 動的割り当て範囲(終了)  | -                     |
| `zone`                   | ゾーン                 | -                   |

## `sakuracloud_vpc_router_l2tp`

VPCルーターでのL2TP/IPSecサーバー機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-remoteaccess/)を参照ください。)

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`              | ◯   | VPCルーターID         | -        | 文字列                   | - |
| `vpc_router_interface_id`    | ◯   | プライベートNIC ID     | -        | 文字列                   | - |
| `pre_shared_secret`          | ◯   | 事前共有シークレット   | -        | 文字列                          | - |
| `range_start`                | ◯   | 動的割り当て範囲(開始) | -        | 文字列                          | - |
| `range_stop`                 | ◯   | 動的割り当て範囲(終了) | -        | 文字列                          | - |
| `zone`                    | -   | ゾーン                 | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |
| `vpc_router_id`          | VPCルーターID          | -                    |
| `vpc_router_interface_id`| プライベートNIC ID     | -                    |
| `pre_shared_secret`      | 事前共有シークレット   | -                     |
| `range_start`            | 動的割り当て範囲(開始) | -                     |
| `range_stop`             | 動的割り当て範囲(終了)  | -                     |
| `zone`                   | ゾーン                 | -                   |

## `sakuracloud_vpc_router_user`

VPCルーターでのリモートユーザーを表します。

このリソースは100個まで指定することが可能です。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-remoteaccess/)を参照ください。)

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`           | ◯   | VPCルーターID         | -        | 文字列                   | - |
| `name`                    | ◯   | ユーザー名 | -        | 文字列                          | - |
| `password`                | ◯   | パスワード | -        | 文字列                          | - |
| `zone`                    | -   | ゾーン                 | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |
| `vpc_router_id`          | VPCルーターID          | -                    |
| `name`                   | ユーザー名 | -                     |
| `password`               | パスワード  | -                     |
| `zone`                   | ゾーン                 | -                   |


## `sakuracloud_vpc_router_site_to_site_vpn`

VPCルーターでのサイト間VPNを表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-site-to-site-vpn/)を参照ください。)

### パラメーター

|パラメーター            |必須  |名称                 |初期値     |設定値                         |補足                                          |
|----------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`      | ◯   | VPCルーターID         | -        | 文字列                   | - |
| `peer`               | ◯   | 対向IPアドレス | -        | 文字列                          | - |
| `remote_id`          | ◯   | 対向ID | -        | 文字列                          | - |
| `pre_shared_secret`  | ◯   | 事前共有シークレット | -        | 文字列                          | - |
| `routes`             | ◯   | 対向Prefix | -        | リスト(文字列)                          | - |
| `local_prefix`       | ◯   | ローカルPrefix | -        | リスト(文字列)                          | - |
| `zone`               | -   | ゾーン                 | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |
| `vpc_router_id`          | VPCルーターID          | -                    |
| `peer`                   | 対向IPアドレス | -                     |
| `remote_id`              | 対向ID | -                     |
| `pre_shared_secret`      | 事前共有シークレット | -                     |
| `routes`                 | 対向Prefix | -                     |
| `local_prefix`           | ローカルPrefix | -                     |
| `zone`                   | ゾーン                 | -                   |

## `sakuracloud_vpc_router_static_route`

VPCルーターでのスタティックルート機能を表します。

(詳細は[さくらのクラウドのマニュアル](http://cloud-news.sakura.ad.jp/vpc-router/vpc-static-route/)を参照ください。)

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `vpc_router_id`              | ◯   | VPCルーターID         | -        | 文字列                   | - |
| `vpc_router_interface_id`    | ◯   | プライベートNIC ID     | -        | 文字列                   | - |
| `prefix`                    | ◯   | プリフィックス | -        | 文字列                          | - |
| `next_hop`                  | ◯   | ネクストホップ | -        | 文字列                          | - |
| `zone`                      | -   | ゾーン                 | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |
| `vpc_router_id`          | VPCルーターID          | -                    |
| `vpc_router_interface_id`| プライベートNIC ID     | -                    |
| `prefix`                 | プリフィックス | -                     |
| `next_hop`               | ネクストホップ  | -                     |
| `zone`                   | ゾーン                 | -                   |
