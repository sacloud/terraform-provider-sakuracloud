# ディスク(sakuracloud_disk)

---

### 設定例

```hcl
#ディスクの定義
resource "sakuracloud_disk" "disk01" {
  name = "disk01"

  #plan      = "ssd"
  #connector = "virtio"
  #size      = 20

  #コピー元アーカイブ/ディスクの指定
  source_archive_id = data.sakuracloud_archive.centos.id

  #source_disk_id = sakuracloud_disk.disk.id

  #コピー元アーカイブのIDが変更された場合のリソース再生成を防ぎたい場合
  #lifecycle {
  #  # コピー元がアーカイブの場合の例
  #  ignore_changes = ["source_archive_id"] 
  #}

  # 指定ディスクと別ストレージに格納したい場合
  #distant_from = ["<your-disk-id>"]

  description = "Description"
  tags        = ["tag1", "tag2"]
}

#コピー元アーカイブの定義
data "sakuracloud_archive" "centos" {
  os_type = "centos"
}
```

### パラメーター

|パラメーター         |必須  |名称                |初期値     |設定値                    |補足                                          |
|-------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `name`            | ◯   | ディスク名           | -        | 文字列                  | - |
| `plan`            | -   | ディスクプラン        | `ssd` | `ssd`<br />`hdd` | - |
| `connector`      | -   | ディスク接続          | `virtio` | `virtio`<br />`ide`    | - |
| `size`            | -   | ディスクサイズ(GB単位) | 20       | 数値                    | - |
| `distant_from`    | -   | ストレージ隔離対象ディスクID | -       | リスト(文字列)                    | ストレージを隔離したいディスクのIDのリスト |
|`source_archive_id`| -   | コピー元アーカイブID   | -        | 文字列                | [注1](#注1) |
|`source_disk_id`   | -   | コピー元ディスクID   | -        | 文字列                | [注1](#注1) |
| `hostname`        | -   | ホスト名               | - | 文字列 | (非推奨) ディスク修正機能で設定される、ホスト名 [注2](#注2)|
| `password`        | -   | パスワード               | - | 文字列 | (非推奨) ディスク修正機能で設定される、OS管理者パスワード [注2](#注2)|
| `ssh_key_ids`     | -   | SSH公開鍵ID             | - | リスト(文字列) | (非推奨) ディスク修正機能で設定される、SSH認証用の公開鍵ID [注2](#注2)|
| `disable_pw_auth` | -   | パスワードでの認証無効化   | - | `true`<br />`false` | (非推奨) ディスク修正機能で設定される、SSH接続でのパスワード/チャレンジレスポンス認証の無効化 [注2](#注2)|
| `note_ids`        | -   | スタートアップスクリプトID | - | リスト(文字列) | (非推奨) ディスク修正機能で設置される、スタートアップスクリプトのID [注2](#注2)|
| `icon_id`         | -   | アイコンID         | - | 文字列 | - |
| `description`     | -   | 説明  | - | 文字列 | - |
| `tags`            | -   | タグ | - | リスト(文字列) | - |
| `graceful_shutdown_timeout` | - | シャットダウンまでの待ち時間 | - | 数値(秒数) | シャットダウンが必要な場合の通常シャットダウンするまでの待ち時間(指定の時間まで待ってもシャットダウンしない場合は強制シャットダウンされる) |
| `zone`            | -   | ゾーン | - | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


#### 互換性

- `source_archive_name`/`source_disk_name`パラメータはv0.3.6にて廃止されました。
v0.3.6以降では[DataResource](data_resource.md)を利用してください。

- `plan`パラメータはv0.5.0にて変更されました。
    - 旧:数値(`2`(HDD) or `4`(SSD))を指定
    - 新:文字列(`hdd` or `ssd`)を指定


#### 注1 コピー元アーカイブ/ディスクの指定 / アーカイブID変更時のリソース再生成の抑制

`source_archive_id`/`source_disk_id`はいずれか1つだけ指定可能です。

また、`source_archive_id`を以下のようにアーカイブデータソースを利用して指定している場合、
さくらのクラウド側でのアーカイブ更新時にアーカイブIDも変更となる場合があります。

```hcl
data sakuracloud_archive "centos" {
  os_type = "centos" # 最新安定版のCentOS
}

resource sakuracloud_disk "disk" {
   # ...
   source_archive_disk = "${data.sakuracloud_archive.centos.id}" # アーカイブデータソースを利用してID指定
}
```

アーカイブIDが変更された後に`terraform apply`を実行すると、ディスクの再生成が行われます。  
この挙動は以下のように記述することで抑制可能です。

```hcl
  lifecycle {
    ignore_changes = ["source_archive_id"] 
  }
```

これは、Terraformの[メタパラメータ](https://www.terraform.io/docs/configuration/resources.html#meta-parameters)と呼ばれるもので、
標準のTerraformの挙動を上書きします。 

この機能を利用する場合、以下の点に留意ください。

  - `source_archive_disk`を手動で変更しても反映されない(`terraform taint`などで手動でリソース再生性が必要)
  - コピー元となるアーカイブが変更されているため、次回ディスク生成時に現在の構成と同じにならない可能性がある

#### 注2 ディスク修正機能の項目

これらの項目は[サーバリソース](server.md)に移動されました。  
過去のバージョンとの互換性維持のために現在でも利用可能ですが、Terraform for さくらのクラウドの将来のバージョンでは削除される予定です。

また、これらの値をサーバリソース/ディスクリソースの両方に記載した場合の動作は不定です。
混乱を避けるためにいずれか一方にのみ記載するようにしてください。

  - OSによりディスク修正機能に対応していない場合があります。
  - これらの値は投入専用です。属性においても投入値を表します(さくらのクラウドAPIからは取得できない項目です)。

### 属性

|属性名                | 名称                    | 補足                                        |
|---------------------|------------------------|--------------------------------------------|
| `id`                | ディスクID               | -                                          |
| `server_id`         | サーバID               | 接続されているサーバのID                     |

