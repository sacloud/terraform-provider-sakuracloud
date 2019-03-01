# IPv4逆引きレコード(sakuracloud_ipv4_ptr)

---

### 設定例

```hcl
# 逆引きレコードの定義
resource "sakuracloud_ipv4_ptr" "foobar" {
  ipaddress = sakuracloud_server.server.ipaddress
  hostname  = "ptr.example.com"
}

# 対象ゾーンの参照
data "sakuracloud_dns" "dns" {
  name_selectors = ["example.com"]
}

# Aレコードの定義
resource "sakuracloud_dns_record" "record01" {
  dns_id = data.sakuracloud_dns.dns.id
  name   = "ptr"
  type   = "A"
  value  = sakuracloud_server.server.ipaddress
}

# Aレコード/PTRレコード登録対象のサーバ
resource "sakuracloud_server" "server" {
  name = "example"
}

```

逆引きレコードを登録するには対応するAレコードがあらかじめ登録されている必要があります。  
詳細は[さくらのクラウド マニュアル](https://manual.sakura.ad.jp/cloud/server/reverse-hostname.html)を参照ください。


### パラメーター

|パラメーター         |必須  |名称                |初期値     |設定値                    |補足                                          |
|-------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
|`ipaddress`       | ◯   | IPアドレス           | -        | 文字列                  | - |
|`hostname`        | ◯   | ホスト名           | -        | 文字列                  | - |
|`retry_max`      | -   | リトライ回数         | `30`     | 数値 | -  |
|`retry_interval` | -   | リトライ時待機時間    | `10`     |数値(秒)| - |
| `zone`            | -   | ゾーン | - | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | 指定のIPアドレスを持つリソースが所属するゾーン |

### 属性

|属性名                | 名称                    | 補足                                        |
|---------------------|------------------------|--------------------------------------------|
| `id`                | ID               | IPアドレスと同値                                          |
