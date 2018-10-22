# DNS(sakuracloud_dns / record)

---

**全ゾーン共通のグローバルリソースです。**

### 設定例

```hcl
# DNSゾーン参照
data "sakuracloud_dns" "dns" {
  name_selectors = ["example.com"]
}

# Aレコード1(test1.example.com)
resource "sakuracloud_dns_record" "record01" {
  dns_id = data.sakuracloud_dns.dns.id
  name   = "test"
  type   = "A"
  value  = "192.168.0.1"
}

# Aレコード2(test.example.com)
resource "sakuracloud_dns_record" "record02" {
  dns_id = data.sakuracloud_dns.dns.id
  name   = "test"
  type   = "A"
  value  = "192.168.0.2"
}
```

## `sakuracloud_dns`

### パラメーター

|パラメーター         |必須  |名称                |初期値     |設定値                    |補足                                          |
|-------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `name_selectors`  | -   | 検索条件(名称)      | -        | リスト(文字列)           | 複数指定した場合はAND条件  |
| `tag_selectors`   | -   | 検索条件(タグ)      | -        | リスト(文字列)           | 複数指定した場合はAND条件  |
| `filter`          | -   | 検索条件(その他)    | -        | オブジェクト             | APIにそのまま渡されます。検索条件を指定してもAPI側が対応していない場合があります。 |


### 属性

|属性名          | 名称             | 補足                                        |
|---------------|-----------------|--------------------------------------------|
| `id`          | ID              | -                                          |
| `dns_servers` | DNSサーバ       | 対象DNSゾーンの委譲先となるネームサーバのリスト  || `zone`            | ◯   | 対象DNSゾーン        | -        | 文字列                  | - |
| `icon_id`     | アイコンID         | - |
| `description` | 説明  | - |
| `tags`        | タグ | - |
