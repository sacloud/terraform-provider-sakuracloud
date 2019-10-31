# ウェブアクセラレータ証明書(sakuradloud_webaccel_certificate)

---

### 設定例

```hcl
# サイト情報の参照用
data sakuracloud_webaccel "site" {
  name = "example"
}

# 証明書
resource sakuracloud_webaccel_certificate "example" {
  site_id           = data.sakuracloud_webaccel.site.id
  certificate_chain = file("crt")
  key               = file("key")
}
```

### パラメーター

| パラメーター              | 必須    | 名称                   | 初期値        | 設定値                      | 補足                                             |
|-- ----------------- | :---: | -------------------- | :--------: | ------------------------ | -------------------------------------------- --|
| `site_id`           | ◯     | ウェブアクセラレータのサイトID                | -          | 文字列                      | -                                              |
| `certificate_chain` | -     | 証明書               | -          | 文字列                      | 中間証明書がある場合はサーバ証明書、中間証明書の順番に連結したものを指定                                              |
| `key`               | -     | 秘密鍵               | -          | 文字列                      | -                                              |

### 属性

| 属性名                   | 名称                       | 補足                                           |
|-- ------------------- | ------------------------ | ------------------------------------------ --|
| `id`                  | ID                   | サイトIDと同一                                            |
| `serial_number`       | シリアルナンバー                 | -                                             |
| `not_before`          | -                        | RFC3339形式                                             |
| `not_after`           | -                        | RFC3339形式                                             |
| `issuer_common_name`  | Issuerコモンネーム             |                                              |
| `subject_common_name` | Subjectコモンネーム            |                                              |
| `dns_names`           | DNS名                     |                                              |
| `sha256_fingerprint`  | フィンガープリント                |                                              |
