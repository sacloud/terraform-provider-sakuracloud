# GSLB(sakuracloud_gslb / server)

---

**全ゾーン共通のグローバルリソースです。**

### 設定例

```hcl
data "sakuracloud_gslb" "gslb" {
 name_selectors = ["foobar"]
}

#GSLB配下のサーバ1
resource "sakuracloud_gslb_server" "gslb_server01" {
  gslb_id   = data.sakuracloud_gslb.gslb.id
  ipaddress = "192.0.2.1"
  #weight    = 1
  #enabled   = true
}


#GSLB配下のサーバ2
resource "sakuracloud_gslb_server" "gslb_server02" {
  gslb_id   = data.sakuracloud_gslb.gslb.id
  ipaddress = "192.0.2.2"
  #weight    = 1
  #enabled   = true
}
```

## `sakuracloud_gslb`

### パラメーター

|パラメーター         |必須  |名称           |初期値     |設定値                    |補足                                          |
|-------------------|:---:|---------------|:--------:|------------------------|----------------------------------------------|
| `name_selectors`  | -   | 検索条件(名称)      | -        | リスト(文字列)           | 複数指定した場合はAND条件  |
| `tag_selectors`   | -   | 検索条件(タグ)      | -        | リスト(文字列)           | 複数指定した場合はAND条件  |
| `filter`          | -   | 検索条件(その他)    | -        | オブジェクト             | APIにそのまま渡されます。検索条件を指定してもAPI側が対応していない場合があります。 |

### 属性

|属性名          | 名称             | 補足                                        |
|---------------|-----------------|--------------------------------------------|
| `id`          | ID              | -                                          |
| `fqdn`        | GSLB-FQDN       | GSLB作成時に割り当てられるFQDN<br />ロードバランシングしたいホスト名をFQDNのCNAMEとしてDNS登録する    |
| `name`        | GSLB名        | - |
| `health_check`| ヘルスチェック  | - |
| `weighted`    | 重み付け応答    | -|
| `sorry_server`| ソーリーサーバ  | - |
| `icon_id`     | アイコンID     | - |
| `description` | 説明  | -      | - |
| `tags`        | タグ | - |

