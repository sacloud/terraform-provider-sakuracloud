# NFS(sakuracloud_nfs)

---

### 設定例

```hcl
#NFS定義
resource sakuracloud_nfs "foobar" {
  switch_id = "${sakuracloud_switch.sw.id}"

  #プラン(現在は100gのみ)
  #plan = "100g"

  ipaddress = "192.168.11.101"
  #default_route = "192.168.11.1"
  nw_mask_len = 24

  name        = "name"
  description = "Description"
  tags        = ["tag1", "tag2"]

  zone = "tk1v"
}

#NFSの上流スイッチ
resource sakuracloud_switch "sw" {
  name = "sw"
  zone = "tk1v"
}
```

## `sakuracloud_nfs`

### パラメーター

|パラメーター       |必須  |名称           |初期値     |設定値                         |補足                                          |
|-----------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `name`          | ◯   | NFS名 | -        | 文字列                         | - |
| `switch_id`     | ◯   | スイッチID      | -        | 文字列                         | - |
| `plan`          | -   | プラン          |`100g`| `100g`    | 現在は値として`100g`のみ利用可能 |
| `ipaddress`     | ◯   | IPアドレス     | -        | 文字列                         | - |
| `nw_mask_len`   | ◯   | ネットマスク     | -        | 数値                          | - |
| `default_route` | -   | ゲートウェイ     | -        | 文字列                        | - |
| `icon_id`       | -   | アイコンID         | - | 文字列 | - |
| `description`   | -   | 説明           | -        | 文字列                         | - |
| `tags`          | -   | タグ           | -        | リスト(文字列)                  | - |
| `zone`          | -   | ゾーン          | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名          | 名称             | 補足                  |
|---------------|------------------|----------------------|
| `id`            | ID | -                    |

