# ISOイメージ/CD-ROM(sakuracloud_cdrom)

---

### 設定例

```hcl
#ISOイメージの参照
data "sakuracloud_cdrom" "cdrom" {
  name_selectors = ["CentOS", "7.4"]
}

# 参照したISOイメージをサーバに挿入
resource "sakuracloud_server" "server" {
  name     = "foobar"
  cdrom_id = data.sakuracloud_cdrom.cdrom.id
}

```

### パラメーター

|パラメーター         |必須  |名称                |初期値     |設定値                    |補足                                          |
|-------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `name_selectors`  | -   | 検索条件(名称)      | -        | リスト(文字列)           | 複数指定した場合はAND条件  |
| `tag_selectors`   | -   | 検索条件(タグ)      | -        | リスト(文字列)           | 複数指定した場合はAND条件  |
| `filter`          | -   | 検索条件(その他)    | -        | オブジェクト             | APIにそのまま渡されます。検索条件を指定してもAPI側が対応していない場合があります。 |
| `zone`            | -   | ゾーン | - | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                | 名称                    | 補足                                        |
|---------------------|------------------------|--------------------------------------------|
| `id`                | ISOイメージID               | -                                          |
| `name`              | ISOイメージ名             | - |
| `size`              | ISOイメージサイズ(GB単位)  | - |
| `icon_id`           | アイコンID                | - |
| `description`       | 説明  | - |
| `tags`              | タグ | - |
| `zone`              | ゾーン | - |

