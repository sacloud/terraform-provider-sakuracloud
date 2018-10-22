# データベース リードレプリカ(sakuracloud_database_read_replica)

---

### 設定例

```hcl
# データベース リードレプリカ定義
resource "sakuracloud_database_read_replica" "foobar" {
  # マスター側データベースのID 
  master_id = sakuracloud_database.master.id

  ipaddress1 = "192.168.11.111"

  # IPアドレス以外のネットワーク関連項目が未指定の場合、マスター側から引き継ぐ
  #switch_id     = sakuracloud_switch.sw.id
  #nw_mask_len   = 24
  #default_route = "192.168.11.1"  

  name = "slave"
}

# データベースの定義(マスター側)
resource "sakuracloud_database" "master" {
  database_type = "postgresql"
  plan          = "10g"
  user_name     = "defuser"
  user_password = "DatabasePasswordUser397"

  replica_password = "DatabasePasswordUser397"

  switch_id     = sakuracloud_switch.sw.id
  ipaddress1    = "192.168.11.101"
  nw_mask_len   = 24
  default_route = "192.168.11.1"

  name = "name"
}

#接続するスイッチの定義
resource "sakuracloud_switch" "sw" {
  name = "sw"
}
```

## `sakuracloud_database_read_replica`

データベースアプライアンス(リードレプリカ)を表します。

### パラメーター

|パラメーター       |必須  |名称           |初期値     |設定値                         |補足                                          |
|-----------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `master_id`     | ◯   | マスター側データベースアプライアンスのID   | -        | 文字列                         | - |
| `name`          | ◯   | データベース名   | -        | 文字列                         | - |
| `ipaddress1`    | ◯   | IPアドレス1     | -        | 文字列                         | - |
| `switch_id`     | -   | スイッチID      | - | 文字列                         | 未指定の場合マスター側の値を引き継ぐ |
| `nw_mask_len`   | -   | ネットマスク     | -        | 数値                          | 未指定の場合マスター側の値を引き継ぐ  |
| `default_route` | -   | ゲートウェイ     | -        | 文字列                        | 未指定の場合マスター側の値を引き継ぐ  |
| `icon_id`       | -   | アイコンID         | - | 文字列 | - |
| `description`   | -   | 説明           | -        | 文字列                         | - |
| `tags`          | -   | タグ           | -        | リスト(文字列)                  | - |
| `graceful_shutdown_timeout` | - | シャットダウンまでの待ち時間 | - | 数値(秒数) | シャットダウンが必要な場合の通常シャットダウンするまでの待ち時間(指定の時間まで待ってもシャットダウンしない場合は強制シャットダウンされる) |
| `zone`          | -   | ゾーン          | -        | `tk1a`<br />`is1b`<br />`is1a` | - |

* マスター側データベースアプライアンスはレプリケーションが有効、かつマスター側として構成されている必要があります。

### 属性

|属性名          | 名称             | 補足                  |
|---------------|------------------|----------------------|
| `id`            | データベースID | -                    |
