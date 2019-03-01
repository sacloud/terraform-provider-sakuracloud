# ロードバランサ(sakuracloud_load_balancer)

---

### 設定例

```hcl
#ロードバランサの定義
resource "sakuracloud_load_balancer" "foobar" {
  switch_id = sakuracloud_switch.sw.id
  vrid      = 1

  #冗長化構成の有無
  high_availability = false

  #プラン(standard or highspec)
  plan = "standard"

  ipaddress1 = "192.168.11.101"

  #ipaddress2 = "192.168.11.101"
  #default_route = "192.168.11.1"
  nw_mask_len = 24

  name        = "name"
  description = "Description"
  tags        = ["tag1", "tag2"]
  
  # VIPの定義
  vips {
    vip          = "192.168.11.201"
    port         = 80
    delay_loop   = 10
    sorry_server = "192.168.11.21"
    
    # VIP配下の実サーバの定義
    servers {
      ipaddress      = "192.168.11.51"
      check_protocol = "http"
      check_path     = "/ping.html"
      check_status   = 200
    }
    servers {
      ipaddress      = "192.168.11.52"
      check_protocol = "http"
      check_path     = "/ping.html"
      check_status   = 200
    }
  }
}

# 互換性維持のため以下の書き方も可能
#ロードバランサVIPの定義
#resource "sakuracloud_load_balancer_vip" "vip1" {
#  load_balancer_id = sakuracloud_load_balancer.foobar.id
#  vip              = "192.168.11.201"
#  port             = 80
#
#  #delay_loop       = 10
#  #sorry_server     = "192.168.11.11"
#  #description      = "description"
#}

#ロードバランサVIP配下の実サーバの定義(1台目)
#resource "sakuracloud_load_balancer_server" "server01" {
#  load_balancer_vip_id = sakuracloud_load_balancer_vip.vip1.id
#  ipaddress            = "192.168.11.51"
#  check_protocol       = "http"
#  check_path           = "/"
#  check_status         = "200"
#
#  #enabled              = true
#}

#ロードバランサVIP配下の実サーバの定義(2台目)
#resource "sakuracloud_load_balancer_server" "server02" {
#  load_balancer_vip_id = sakuracloud_load_balancer_vip.vip1.id
#  ipaddress            = "192.168.11.52"
#  check_protocol       = "http"
#  check_path           = "/"
#  check_status         = "200"
#
#  #enabled              = true
#}

#ロードバランサ上流のスイッチ
resource "sakuracloud_switch" "sw" {
  name = "sw"
}
```

## `sakuracloud_load_balancer`

ロードバランサ本体を表します。

### パラメーター

|パラメーター       |必須  |名称           |初期値     |設定値                         |補足                                          |
|-----------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `name`          | ◯   | ロードバランサ名 | -        | 文字列                         | - |
| `switch_id`     | ◯   | スイッチID      | -        | 文字列                         | - |
| `vrid`          | ◯   | VRID           | -        | 数値                          | - |
| `high_availability`     | -   | 冗長化          | `false`  | `true`<br />`false`           | - |
| `plan`          | -   | プラン          |`standard`| `standard`<br />`highspec`    | - |
| `ipaddress1`    | ◯   | IPアドレス1     | -        | 文字列                         | - |
| `ipaddress2`    | △   | IPアドレス2     | -        | 文字列                         | 冗長化構成の場合必須 |
| `nw_mask_len`   | ◯   | ネットマスク     | -        | 数値                          | - |
| `default_route` | -   | ゲートウェイ     | -        | 文字列                        | - |
| `vips`    | -   | VIP | -        | リスト | 詳細は[`vips`](#vips)を参照    |
| `icon_id`       | -   | アイコンID         | - | 文字列 | - |
| `description`   | -   | 説明           | -        | 文字列                         | - |
| `tags`          | -   | タグ           | -        | リスト(文字列)                  | - |
| `graceful_shutdown_timeout` | - | シャットダウンまでの待ち時間 | - | 数値(秒数) | シャットダウンが必要な場合の通常シャットダウンするまでの待ち時間(指定の時間まで待ってもシャットダウンしない場合は強制シャットダウンされる) |
| `zone`          | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |

### VIPs

この要素は最大10個まで指定可能です。  

|パラメーター          |必須  |名称           |初期値     |設定値                         |補足                                          |
|--------------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `vip`              | ◯   | VIPアドレス     | -        | 文字列                         | - |
| `port`             | ◯   | ポート番号      | -        | 数値                          | - |
| `delay_loop`       | -   | チェック間隔秒数  | `10`    | `10`〜`2147483647`の整数           | - |
| `sorry_server`     | -   | ソーリーサーバ  | -        | 文字列     | VIPに紐づく実サーバがすべてダウンした場合、<br />すべてのアクセスを指定したサーバに誘導します |
| `description`      | -   | 説明           | -        | 文字列                         | - |
| `servers`    | -   | 実サーバ| -        | リスト | 詳細は[`servers`](#servers)を参照    |

### `servers`

この要素は最大40個まで指定可能です。

|パラメーター             |必須  |名称           |初期値     |設定値                         |補足                                          |
|------------------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `ipaddress`            | ◯   | IPアドレス      | -        | 文字列            | - |
| `check_protocol`       | ◯   | チェック方法     | -        | `ping`<br />`tcp`<br />`http`<br />`https` | - |
| `check_path`           | △   | チェック対象パス  | -       | 文字列           | チェック方法が`http`、`https`の場合必須 |
| `check_status`         | △   | チェック期待値   | -        | 文字列           | 期待するレスポンスコード<br />チェック方法が`http`、`https`の場合必須 |
| `enabled`              | -   | 有効/無効       | `true`    | `true`<br />`false`   | - |


### 属性

|属性名          | 名称             | 補足                  |
|---------------|------------------|----------------------|
| `id`            | ロードバランサID | -                    |
| `vip_ids`       | VIP IDリスト     | ロードバランサ配下のVIPのIDリスト   |

**注意**  

同一のロードバランサに対し`vips`属性と`sakuracloud_load_balancer_vip`リソースの併用はできません。  
また、同一のロードバランサVIPに対し`servers`属性と`sakuracloud_load_balancer_server`リソースの併用はできません。


## `sakuracloud_load_balancer_vip`

ロードバランサが持つVIPを表します。

1台のロードバランサにつき4つまでのVIPを登録できます。
(詳細は[さくらのクラウドのマニュアル](https://help.sakura.ad.jp/app/answers/detail/a_id/2517)を参照ください)

### パラメーター

|パラメーター          |必須  |名称           |初期値     |設定値                         |補足                                          |
|--------------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `load_balancer_id` | ◯   | ロードバランサID | -        | 文字列                         | - |
| `vip`              | ◯   | VIPアドレス     | -        | 文字列                         | - |
| `port`             | ◯   | ポート番号      | -        | 数値                          | - |
| `delay_loop`       | -   | チェック間隔秒数  | `10`    | `10`〜`2147483647`の整数           | - |
| `sorry_server`     | -   | ソーリーサーバ  | -        | 文字列     | VIPに紐づく実サーバがすべてダウンした場合、<br />すべてのアクセスを指定したサーバに誘導します |
| `description`      | -   | 説明           | -        | 文字列                         | - |
| `zone`             | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名          | 名称             | 補足                  |
|---------------|------------------|----------------------|
| `id`               | ID             | -                    |
| `servers`          | 実サーバIDリスト           | 配下の実サーバのIDリスト   |

## `sakuracloud_load_balancer_server`

ロードバランサが持つVIP配下の実サーバを表します。

1つのVIPにつき、40台までの実サーバを登録できます。
(詳細は[さくらのクラウドのマニュアル](https://help.sakura.ad.jp/app/answers/detail/a_id/2517)を参照ください)

### パラメーター

|パラメーター             |必須  |名称           |初期値     |設定値                         |補足                                          |
|------------------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `load_balancer_vip_id` | ◯   | VIP ID         | -        | 文字列            | - |
| `ipaddress`            | ◯   | IPアドレス      | -        | 文字列            | - |
| `check_protocol`       | ◯   | チェック方法     | -        | `ping`<br />`tcp`<br />`http`<br />`https` | - |
| `check_path`           | △   | チェック対象パス  | -       | 文字列           | チェック方法が`http`、`https`の場合必須 |
| `check_status`         | △   | チェック期待値   | -        | 文字列           | 期待するレスポンスコード<br />チェック方法が`http`、`https`の場合必須 |
| `enabled`              | -   | 有効/無効       | `true`    | `true`<br />`false`   | - |
| `zone`                 | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名          | 名称             | 補足                  |
|---------------|------------------|----------------------|
| `id`               | ID             | -                    |
