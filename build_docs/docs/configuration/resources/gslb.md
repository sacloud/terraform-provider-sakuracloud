# GSLB(sakuracloud_gslb / server)

---

**全ゾーン共通のグローバルリソースです。**

`sakuracloud_gslb`がGSLB設定を、`sakuracloud_gslb_server`が対象サーバを表しています。

### 設定例

```hcl
resource "sakuracloud_gslb" "gslb" {
  name = "gslb_from_terraform"

  health_check {
    protocol    = "http"
    delay_loop  = 10
    host_header = "example.com"
    path        = "/"
    status      = "200"
  }

  #port       = 80

  #weighted     = false
  #sorry_server = "192.0.2.254"
  description = "GSLB from terraform for SAKURA CLOUD"
  tags        = ["tag1", "tag2"]
}

#GSLB配下のサーバ1
resource "sakuracloud_gslb_server" "gslb_server01" {
  gslb_id   = sakuracloud_gslb.gslb.id
  ipaddress = "192.0.2.1"
  #weight    = 1
  #enabled   = true
}


#GSLB配下のサーバ2
resource "sakuracloud_gslb_server" "gslb_server02" {
  gslb_id   = sakuracloud_gslb.gslb.id
  ipaddress = "192.0.2.2"
  #weight    = 1
  #enabled   = true
}

```

## `sakuracloud_gslb`

### パラメーター

|パラメーター         |必須  |名称           |初期値     |設定値                    |補足                                          |
|-------------------|:---:|---------------|:--------:|------------------------|----------------------------------------------|
| `name`            | ◯   | GSLB名        | -        | 文字列                  | - |
| `health_check`    | ◯   | ヘルスチェック  | -        | マップ                  | 詳細は[`health_check`](#health_check)を参照    |
| `weighted`        | -   | 重み付け応答    | `false` | `true`<br />`false` | `true`:有効<br />`false`:無効 |
| `sorry_server`     | -   | ソーリーサーバ  | -      | 文字列 | - |
| `icon_id`         | -   | アイコンID         | - | 文字列 | - |
| `description`     | -   | 説明  | -      | 文字列 | - |
| `tags`            | -   | タグ | -      | リスト(文字列) | - |

### `health_check`

|パラメーター     |必須  |名称                |初期値     |設定値                    |補足                                          |
|---------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `protocol`    | ◯   | プロトコル        | -        | `http`<br />`https`<br />`ping`<br />`tcp`| - |
| `delay_loop`  | -   | チェック間隔(秒)        | `10`        | 数値                  | `10`〜`60` |
| `host_header` | -   | Hostヘッダ  | - | 文字列 | プロトコルが`http`または`https`の場合のみ有効 |
| `path`        | △   | パス  | - | 文字列 | プロトコルが`http`または`https`の場合のみ有効かつ必須 |
| `status`      | △   | レスポンスコード | - | 文字列 | プロトコルが`http`または`https`の場合のみ有効かつ必須 |
| `port`        | △   | ポート番号 | - | 数値 | プロトコルが`tcp`の場合のみ有効かつ必須 |

### 属性

|属性名          | 名称             | 補足                                        |
|---------------|-----------------|--------------------------------------------|
| `id`          | ID              | -                                          |
| `fqdn`        | GSLB-FQDN       | GSLB作成時に割り当てられるFQDN<br />ロードバランシングしたいホスト名をFQDNのCNAMEとしてDNS登録する    |



## `sakuracloud_gslb_server`

### パラメーター

|パラメーター  |必須  |名称          |初期値   |設定値                 |補足                                          |
|------------|:---:|--------------|:------:|---------------------|----------------------------------------------|
| `gslb_id`  | ◯   | GSLB-ID      | -      | 文字列                | 対象GSLBのID |
| `ipaddress`| ◯   | IPアドレス     | -      | 文字列               | 監視対象サーバのIPアドレス|
| `enabled`  | -   | 有効          | `true` | `true`<br />`false` | - |
| `weight`   | -   | 重み          | `1`    | 数値                 | 重み付け応答が有効な場合のみ有効。`1`〜`10000`|


### 属性

|属性名       | 名称             | 補足 |
|------------|-----------------|------|
| `id`       | ID              | -  |
