# エンハンスドロードバランサ(sakuracloud_proxylb)

---

**全ゾーン共通のグローバルリソースです。**

### 設定例

```hcl
resource "sakuracloud_proxylb" "foobar" {
  name         = "foobar"
  plan         = 1000 
  vip_failover = true # default: false

  health_check {
    protocol    = "http"
    path        = "/"
    host_header = "example.com"
    delay_loop  = 10
  }
  
  bind_ports {
    proxy_mode        = "http"
    port              = 80
    redirect_to_https = true
  }
  bind_ports {
    proxy_mode    = "https"
    port          = 443
    support_http2 = true
  }
  
  sorry_server {
    ipaddress = "192.2.0.1"
    port      = 80
  }

  servers {
    ipaddress = "133.242.0.3"
    port = 80
  }
  servers {
    ipaddress = "133.242.0.4"
    port = 80
  }

  certificate {
    server_cert = file("server.crt")
    private_key = file("server.key")    
    # intermediate_cert = file("intermediate.crt")
    
    additional_certificates {
      server_cert = file("server2.crt")
      private_key = file("server2.key")    
      # intermediate_cert = file("intermediate2.crt")
    }
    
    additional_certificates {
      server_cert = file("server3.crt")
      private_key = file("server3.key")    
      # intermediate_cert = file("intermediate3.crt")
    }
  }
}
```

## `sakuracloud_proxylb`

### パラメーター

|パラメーター         |必須  |名称           |初期値     |設定値                    |補足                                          |
|-------------------|:---:|---------------|:--------:|------------------------|----------------------------------------------|
| `name`            | ◯   | エンハンスドロードバランサ名        | -        | 文字列                  | - |
| `plan`            | -   | プラン        | `1000`        | `1000`<br />`5000`<br />`10000`<br />`50000`<br />`100000`     | - |
| `vip_failover`    | -   | VIPフェイルオーバ | `false`        | `true` or `false`     | - |
| `bind_ports`      | ◯   | 待ち受けポート  | -        | リスト                  | 詳細は[`bind_ports`](#bind_ports)を参照    |
| `health_check`    | ◯   | ヘルスチェック  | -        | マップ                  | 詳細は[`health_check`](#health_check)を参照    |
| `sorry_server`     | -   | ソーリーサーバ  | -      | マップ| 詳細は[`sorry_server`](#sorry_server)を参照 |
| `servers`     | -   | 実サーバ  | -      | リスト | 詳細は[`servers`](#servers)を参照 |
| `certificate`     | -   | SSL証明書 | -      | マップ | 詳細は[`certificate`](#certificate)を参照 |
| `icon_id`         | -   | アイコンID         | - | 文字列 | - |
| `description`     | -   | 説明  | -      | 文字列 | - |
| `tags`            | -   | タグ | -      | リスト(文字列) | - |

### `bind_ports`

この要素は最大2つまで指定可能です。  

|パラメーター     |必須  |名称                |初期値     |設定値                    |補足                                          |
|---------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `proxy_mode`    | ◯   | プロキシ方式 | -        | `http`<br />`https`| - |
| `port`        | ◯  | ポート番号 | - | 数値 | - |
| `redirect_to_https`  | -  | HTTPSへのリダイレクト | - | bool | `proxy_mode`が`http`の場合のみ有効 |
| `support_http2`      | -  | HTTP/2のサポート | - | bool | `proxy_mode`が`https`の場合のみ有効  |

### `health_check`

|パラメーター     |必須  |名称                |初期値     |設定値                    |補足                                          |
|---------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `protocol`    | ◯   | プロトコル        | -        | `http`<br />`tcp`| - |
| `delay_loop`  | -   | チェック間隔(秒)        | `10`        | 数値                  | `10`〜`60` |
| `host_header` | -   | Hostヘッダ  | - | 文字列 | プロトコルが`http`の場合のみ有効 |
| `path`        | △   | パス  | - | 文字列 | プロトコルが`http`の場合のみ有効かつ必須 |

### `sorry_server`

|パラメーター     |必須  |名称                |初期値     |設定値                    |補足                                          |
|---------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `ipaddress`   | ◯  | IPアドレス | -        | 文字列 | - |
| `port`        | ◯  | ポート番号 | - | 数値 | - |

### `servers`

この要素は最大40個まで指定可能です。  

|パラメーター  |必須  |名称          |初期値   |設定値                 |補足                                          |
|------------|:---:|--------------|:------:|---------------------|----------------------------------------------|
| `ipaddress`| ◯   | IPアドレス     | -      | 文字列               | 実サーバのIPアドレス|
| `port`     | ◯  | ポート番号 | - | 数値 | - |
| `enabled`  | -   | 有効          | `true` | `true`<br />`false` | - |

### `certificate`

|パラメーター  |必須  |名称          |初期値   |設定値                 |補足                                          |
|------------|:---:|--------------|:------:|---------------------|----------------------------------------------|
| `server_cert`      | ◯  | サーバ証明書 | -      | 文字列               | -|
| `intermediate_cert`| -  | 中間証明書   | - | 文字列 | - |
| `private_key`      | ◯  | 秘密鍵      | - | 文字列 | - |
| `additional_certificates`| -  | 追加証明書      | - | リスト | 詳細は[`certificate`](#certificate)を参照 |


### 属性

|属性名          | 名称             | 補足                                        |
|---------------|-----------------|--------------------------------------------|
| `id`          | ID              | -                                          |
| `vip`        | VIP       | ロードバランサに割り当てられたグローバルIP(`vip_failover`が`false`の場合のみ有効)    |
| `fqdn`        | VIP       | ロードバランサに割り当てられたグローバルIP(`vip_failover`が`true`の場合のみ有効)    |
| `proxy_networks`  | プロキシ元ネットワーク | プロキシ元IP(CIDR)のリスト    |

