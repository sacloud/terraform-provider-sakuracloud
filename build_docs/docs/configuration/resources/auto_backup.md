# 自動バックアップ(sakuracloud_auto_backup)

---

自動バックアップの設定を行うリソースです。

### 設定例

```hcl
#ディスクの定義
resource "sakuracloud_disk" "disk01" {
  name = "disk01"
}

#自動バックアップ
resource "sakuracloud_auto_backup" "foobar" {
  name           = "auto_backup"
  disk_id        = sakuracloud_disk.disk01.id
  weekdays       = ["mon", "tue", "wed"]
  max_backup_num = 1
  description    = "description"
  tags           = ["tag1", "tag2"]
}
```

### パラメーター

|パラメーター       |必須  |名称                |初期値     |設定値                    |補足                                          |
|-----------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `name`          | ◯   | 自動バックアップ名   | -        | 文字列                  | - |
| `disk_id`       | ◯   | ディスクID         | - | 文字列 | - |
| `weekdays`      | ◯   | バックアップ取得曜日 | - | 以下の値のリスト<br />`mon`<br />`tue`<br />`wed`<br />`thu`<br />`fri`<br />`sat`<br />`sun`|- |
| `max_backup_num`| -   | 保持世代数         | 1 | 数値 | `1`から`10`までの整数 |
| `icon_id`       | -   | アイコンID         | - | 文字列 | - |
| `description`   | -   | 説明              | - | 文字列 | - |
| `tags`          | -   | タグ              | - | リスト | - |
| `zone`          | -   | 対象ゾーン          | - | `is1b`<br />`tk1a`<br />`is1a` | - |

### 属性

|属性名                | 名称                    | 補足                                        |
|---------------------|------------------------|--------------------------------------------|
| `id`                | ID                     | -                                          |
