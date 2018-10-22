# ISOイメージ/CD-ROM(sakuracloud_cdrom)

---

### 設定例

```hcl
#ISOイメージの定義
resource "sakuracloud_cdrom" "cdrom" {
  name = "cdrom01"

  #size = 5 # or 10, デフォルトは5GB

  # ISOイメージをアップロードする場合
  iso_image_file = "test/dummy.iso"
  hash           = md5(file("test/dummy.iso"))

  # 単一ファイルを内包するISOイメージを作成する場合
  #content           = file("test/dummy-upd.json")
  #content_file_name = "config"

  description = "Description"
  tags        = ["tag1", "tag2"]
}

```

### パラメーター

|パラメーター         |必須  |名称                |初期値     |設定値                    |補足                                          |
|-------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `name`            | ◯   | ISOイメージ名           | -        | 文字列                  | - |
| `size`            | -   | ISOイメージサイズ(GB単位) | 5       |  `5`<br />`10`         | - |
| `iso_image_file`  | -   | ISOイメージファイルパス| - | 文字列 | ※注1 |
| `hash`            | -   | ISOイメージファイルのMD5ハッシュ値| - | 文字列 | `iso-image-file`の変更検知用MD5ハッシュ |
| `content`         | -   | ISOイメージコンテント  | - | 文字列 | ※注1 |
| `content_file_name`| -  | 作成されるISOイメージ内のファイル名 | `config` | 文字列 | 現バージョンではボリュームラベルとしてもこの値が利用される |
| `icon_id`         | -   | アイコンID         | - | 文字列 | - |
| `description`     | -   | 説明  | - | 文字列 | - |
| `tags`            | -   | タグ | - | リスト(文字列) | - |
| `zone`            | -   | ゾーン | - | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


#### 注1

`iso_image_file`/`content`はいずれか1つだけ指定可能です。


### 属性

|属性名                | 名称                    | 補足                                        |
|---------------------|------------------------|--------------------------------------------|
| `id`                | ISOイメージID               | -                                          |

