# SIM(sakuracloud_sim)

---

### 設定例

```hcl
# SIMの定義
resource "sakuracloud_sim" "sim" {
  name        = "example"
  description = "example"
  tags        = ["tags1", "tags2"]

  iccid    = "<SIMに記載されているICCID>"
  passcode = "<SIMに記載されているPasscode>"
  imei     = "<端末識別番号(IMEIロックする場合のみ)>"

  #enabled  = true

  mobile_gateway_id = sakuracloud_mobile_gateway.mgw.id # 接続するモバイルゲートウェイのID
  ipaddress         = "192.168.100.2"                   # SIMに割り当てるIPアドレス        
}

# モバイルゲートウェイの定義
resource "sakuracloud_mobile_gateway" "mgw" {
  name                = "example"
  internet_connection = true
}
```

### パラメーター

|パラメーター           |必須  |名称                |初期値     |設定値                    |補足                                          |
|---------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `name`              | ◯   | 名称              | -        | 文字列                  | - |
| `iccid`             | ◯   | ICCID             | -        | 文字列                  | - |
| `passcode`          | ◯   | パスコード         | -        | 文字列                  | - |
| `imei`              | -   | IMEI(端末識別番号)  | - | 文字列 | - |
| `enabled`           | -   | 有効/無効          | `true` | `true`<br />`false`| - |
| `mobile_gateway_id` | -   | モバイルゲートウェイID  | - | 文字列 | - |
| `ipaddress`         | -   | IPアドレス  | - | 文字列 | - |
| `description`       | -   | 説明  | - | 文字列 | - |
| `tags`              | -   | タグ | - | リスト(文字列) | - |
| `icon_id`           | -   | アイコンID         | - | 文字列| - |
| `zone`              | -   | ゾーン | - | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |

### 属性

|属性名                | 名称                    | 補足                                        |
|---------------------|------------------------|--------------------------------------------|
| `id`                | SIM ID               | -                                          |
