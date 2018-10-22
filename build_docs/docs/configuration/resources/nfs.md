# NFS(sakuracloud_nfs)

---

### 設定例

```hcl
#NFS定義
resource "sakuracloud_nfs" "foobar" {
  switch_id = sakuracloud_switch.sw.id

  #プラン[100/500/1024(1T)/2048(2T)/4096(4T)]
  #plan = "100"

  ipaddress = "192.168.11.101"

  #default_route = "192.168.11.1"
  nw_mask_len = 24

  name        = "name"
  description = "Description"
  tags        = ["tag1", "tag2"]
}

#NFSの上流スイッチ
resource "sakuracloud_switch" "sw" {
  name = "sw"
}
```

## `sakuracloud_nfs`

### パラメーター

|パラメーター       |必須  |名称           |初期値     |設定値                         |補足                                          |
|-----------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `name`          | ◯   | NFS名 | -        | 文字列                         | - |
| `switch_id`     | ◯   | スイッチID      | -        | 文字列                         | - |
| `plan`          | -   | プラン          |`100`| `100`<br />`500`<br />`1024`<br />`2048`<br />`4096`    |- |
| `ipaddress`     | ◯   | IPアドレス     | -        | 文字列                         | - |
| `nw_mask_len`   | ◯   | ネットマスク     | -        | 数値                          | - |
| `default_route` | -   | ゲートウェイ     | -        | 文字列                        | - |
| `icon_id`       | -   | アイコンID         | - | 文字列 | - |
| `description`   | -   | 説明           | -        | 文字列                         | - |
| `tags`          | -   | タグ           | -        | リスト(文字列)                  | - |
| `graceful_shutdown_timeout` | - | シャットダウンまでの待ち時間 | - | 数値(秒数) | シャットダウンが必要な場合の通常シャットダウンするまでの待ち時間(指定の時間まで待ってもシャットダウンしない場合は強制シャットダウンされる) |
| `zone`            | -   | ゾーン | - | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名          | 名称             | 補足                  |
|---------------|------------------|----------------------|
| `id`            | ID | -                    |

