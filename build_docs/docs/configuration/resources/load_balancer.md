# ロードバランサー(sakuracloud_load_balancer)

---

### 設定例

```hcl
# ロードバランサー上流のスイッチ
resource "sakuracloud_switch" "sw" {
    name = "sw"
    zone = "tk1v"
}

# ロードバランサーの定義
resource "sakuracloud_load_balancer" "foobar" {
    switch_id = "${sakuracloud_switch.sw.id}"
    VRID = 1
    ipaddress1 = "192.168.11.101"
    nw_mask_len = 24
    name = "name"
    zone = "tk1v"
}

# ロードバランサーVIPの定義
resource "sakuracloud_load_balancer_vip" "vip1" {
    load_balancer_id = "${sakuracloud_load_balancer.foobar.id}"
    vip = "192.168.11.201"
    port = 80
    zone = "tk1v"
}

# ロードバランサーVIP配下の実サーバーの定義(1台目)
resource "sakuracloud_load_balancer_server" "server01"{
    load_balancer_vip_id = "${sakuracloud_load_balancer_vip.vip1.id}"
    ipaddress = "192.168.11.51"
    check_protocol = "http"
    check_path = "/"
    check_status = "200"
    zone = "tk1v"
}
# ロードバランサーVIP配下の実サーバーの定義(2台目)
resource "sakuracloud_load_balancer_server" "server02"{
    load_balancer_vip_id = "${sakuracloud_load_balancer_vip.vip1.id}"
    ipaddress = "192.168.11.52"
    check_protocol = "http"
    check_path = "/"
    check_status = "200"
    zone = "tk1v"
}
```

## `sakuracloud_load_balancer`

ロードバランサー本体を表します。

### パラメーター

|パラメーター       |必須  |名称           |初期値     |設定値                         |補足                                          |
|-----------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `name`          | ◯   | ロードバランサ名 | -        | 文字列                         | - |
| `switch_id`     | ◯   | スイッチID      | -        | 文字列                         | - |
| `VRID`          | ◯   | VRID           | -        | 数値                          | - |
| `high_availability`     | -   | 冗長化          | `false`  | `true`<br />`false`           | - |
| `plan`          | -   | プラン          |`standard`| `standard`<br />`highspec`    | - |
| `ipaddress1`    | ◯   | IPアドレス1     | -        | 文字列                         | - |
| `ipaddress2`    | △   | IPアドレス2     | -        | 文字列                         | 冗長化構成の場合必須 |
| `nw_mask_len`   | ◯   | ネットマスク     | -        | 数値                          | - |
| `default_route` | -   | ゲートウェイ     | -        | 文字列                        | - |
| `description`   | -   | 説明           | -        | 文字列                         | - |
| `tags`          | -   | タグ           | -        | リスト(文字列)                  | - |
| `zone`          | -   | ゾーン          | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名          | 名称             | 補足                  |
|---------------|------------------|----------------------|
| `id`            | ロードバランサID | -                    |
| `name`          | ロードバランサ名 | -                    |
| `switch_id`     | スイッチID      | -                    |
| `VRID`          | VRID           | -                     |
| `high_availability`     | 冗長化          | -                    |
| `plan`          | プラン          | -                    |
| `ipaddress1`    | IPアドレス1      | -                    |
| `ipaddress2`    | IPアドレス2      | -                    |
| `nw_mask_len`   | ネットマスク      | -                   |
| `default_route` | ゲートウェイ      | -                   |
| `description`   | 説明             | -                   |
| `tags`          | タグ             | -                  |
| `zone`          | ゾーン           | -                   |
| `vip_ids`       | VIP IDリスト     | ロードバランサー配下のVIPのIDリスト   |

## `sakuracloud_load_balancer_vip`

ロードバランサーが持つVIPを表します。

1台のロードバランサーにつき4つまでのVIPを登録できます。
(詳細は[さくらのクラウドのマニュアル](https://help.sakura.ad.jp/app/answers/detail/a_id/2517)を参照ください。)

### パラメーター

|パラメーター          |必須  |名称           |初期値     |設定値                         |補足                                          |
|--------------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `load_balancer_id` | ◯   | ロードバランサID | -        | 文字列                         | - |
| `vip`              | ◯   | VIPアドレス     | -        | 文字列                         | - |
| `port`             | ◯   | ポート番号      | -        | 数値                          | - |
| `delay_loop`       | -   | チェック間隔秒数  | `10`    | `10`〜`2147483647`の整数           | - |
| `sorry_server`     | -   | ソーリーサーバー  | -        | 文字列     | VIPに紐づく実サーバがすべてダウンした場合、<br />すべてのアクセスを指定したサーバに誘導します |
| `zone`             | -   | ゾーン          | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名          | 名称             | 補足                  |
|---------------|------------------|----------------------|
| `id`               | ID             | -                    |
| `load_balancer_id` | ロードバランサID | -                    |
| `vip`              | VIPアドレス      | -                    |
| `port`             | ポート番号           | -                     |
| `delay_loop`       | チェック間隔秒数          | -                    |
| `sorry_server`     | ソーリーサーバー          | -                    |
| `zone`             | ゾーン           | -                   |
| `servers`          | 実サーバーIDリスト           | 配下の実サーバーのIDリスト   |

## `sakuracloud_load_balancer_server`

ロードバランサーが持つVIP配下の実サーバーを表します。

1つのVIPにつき、40台までの実サーバーを登録できます。
(詳細は[さくらのクラウドのマニュアル](https://help.sakura.ad.jp/app/answers/detail/a_id/2517)を参照ください。)

### パラメーター

|パラメーター             |必須  |名称           |初期値     |設定値                         |補足                                          |
|------------------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `load_balancer_vip_id` | ◯   | VIP ID         | -        | 文字列            | - |
| `ipaddress`            | ◯   | IPアドレス      | -        | 文字列            | - |
| `check_protocol`       | ◯   | チェック方法     | -        | `ping`<br />`tcp`<br />`http`<br />`https` | - |
| `check_path`           | △   | チェック対象パス  | -       | 文字列           | チェック方法が`http`、`https`の場合必須 |
| `check_status`         | △   | チェック期待値   | -        | 文字列           | 期待するレスポンスコード<br />チェック方法が`http`、`https`の場合必須 |
| `enabled`              | -   | 有効/無効       | `true`    | `true`<br />`false`   | - |
| `zone`                 | -   | ゾーン          | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名          | 名称             | 補足                  |
|---------------|------------------|----------------------|
| `id`               | ID             | -                    |
| `load_balancer_vip_id` | VIP ID | -                    |
| `ipaddress`              | IPアドレス      | -                    |
| `check_protocol`             | チェック方法           | -                     |
| `check_path`       | チェック対象パス          | -                    |
| `check_status`     | チェック期待値          | -                    |
| `enabled`       | 有効/無効| -                    |
| `zone`             | ゾーン           | -                   |
