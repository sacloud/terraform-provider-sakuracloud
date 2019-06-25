# エンハンスドロードバランサ Let's Encrypt設定(sakuracloud_proxylb_acme)

---

**全ゾーン共通のグローバルリソースです。**

エンハンスドロードバランサにてLet's Encryptで証明書を取得するための設定を行うリソースです。

### 設定例

```hcl
resource "sakuracloud_proxylb_acme" "cert" {
  proxylb_id = sakuracloud_proxylb.foobar.id
  accept_tos = true
  common_name = "foobar.example.com"
  update_delay_sec = 120
}

resource "sakuracloud_proxylb" "foobar" {
  name         = "foobar"
  plan         = 1000 
  vip_failover = true # default: false

  bind_ports {
    proxy_mode        = "http"
    port              = 80
  }
  bind_ports {
    proxy_mode    = "https"
    port          = 443
  }
  
  servers {
    ipaddress = "133.242.0.3"
    port = 80
  }
  servers {
    ipaddress = "133.242.0.4"
    port = 80
  }
}
```

## `sakuracloud_proxylb_acme`

### パラメーター

|パラメーター         |必須  |名称           |初期値     |設定値                    |補足                                          |
|-------------------|:---:|---------------|:--------:|------------------------|----------------------------------------------|
| `proxylb_id`            | ◯   | エンハンスドロードバランサID        | -        | 文字列                  | - |
| `accept_tos`            | ◯   | 利用規約への同意 | -        | Let's Encryptの[利用規約](https://letsencrypt.org/repository/)への同意     | `true`の場合のみ証明書の発行を行う |
| `common_name`    | ◯   | コモンネーム | -        | -     | 証明書発行対象となるFQDN |
| `update_delay_sec`      | -   | 更新待ち秒数  | -        | 0                  | エンハンスドロードバランサへのLet's Encrpt設定投入までの待ち時間 ([注1](#注1))    |

#### 注1 更新待ち秒数について

エンハンスドロードバランサでLet's Encryptでの証明書取得を行うには、`common_name`で指定したFQDNがエンハンスドロードバランサのVIPまたはFQDNを指している必要があります。  
参考: [さくらのクラウドマニュアル - Let's Encrypt証明書自動インストール・更新機能](https://manual.sakura.ad.jp/cloud/appliance/enhanced-lb/#let-s-encrypt)

TerraformにてDNSレコードの登録を行う場合、レコード作成直後にコモンネームの解決をできない場合があります。  
この項目はそのような場合のために待ち時間を指定するものです。

### 属性

|属性名          | 名称             | 補足                                        |
|---------------|-----------------|--------------------------------------------|
| `id`          | ID              | エンハンスドロードバランサのIDが設定される                                          |
| `certificate`  | 証明書 | 発行された証明書(sakuracloud_proxylbにて参照できる値と同じ)    |

### `certificate`

|パラメーター  |名称          |初期値   |設定値                 |補足                                          |
|------------|--------------|:------:|---------------------|----------------------------------------------|
| `server_cert`      | サーバ証明書 | -      | 文字列               | -|
| `intermediate_cert`| 中間証明書   | - | 文字列 | - |
| `private_key`      | 秘密鍵      | - | 文字列 | - |
| `additional_certificates`| 追加証明書      | - | リスト | 詳細は[`certificate`](#certificate)を参照 |
