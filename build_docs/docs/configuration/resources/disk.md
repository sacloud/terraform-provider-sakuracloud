# ディスク(sakuracloud_disk)

---

### 設定例

```hcl
#ディスクの定義
resource sakuracloud_disk "disk01" {
  name      = "disk01"
  #plan      = "ssd"
  #connector = "virtio"
  #size      = 20

  #コピー元アーカイブ/ディスクの指定
  source_archive_id = "${data.sakuracloud_archive.centos.id}"
  #source_disk_id = "${sakuracloud_disk.disk.id}"

  #ディスクの修正関連
  hostname = "your-host-name"
  password = "your-password"
  
  #SSH公開鍵
  ssh_key_ids = ["${sakuracloud_ssh_key_gen.key.id}"]
  
  #スタートアップスクリプト
  note_ids = ["${sakuracloud_note.note.id}"]
  
  #SSH接続でのパスワード/チャレンジレスポンス認証無効化
  disable_pw_auth = true
 
  description = "Description"
  tags        = ["tag1", "tag2"]
}

#コピー元アーカイブの定義
data sakuracloud_archive "centos" {
  os_type = "centos"
}

#SSH公開鍵
resource sakuracloud_ssh_key_gen "key" {
  name = "foobar"
}

#スタートアップスクリプト
resource sakuracloud_note "note" {
  name    = "note"
  content = "#!/bin/sh ..."
}


```

### パラメーター

|パラメーター         |必須  |名称                |初期値     |設定値                    |補足                                          |
|-------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `name`            | ◯   | ディスク名           | -        | 文字列                  | - |
| `plan`            | -   | ディスクプラン        | `ssd` | `ssd`<br />`hdd` | - |
| `connector`      | -   | ディスク接続          | `virtio` | `virtio`<br />`ide`    | - |
| `size`            | -   | ディスクサイズ(GB単位) | 20       | 数値                    | - |
|`source_archive_id`| -   | コピー元アーカイブID   | -        | 文字列                | [注1](#注1) |
|`source_disk_id`   | -   | コピー元ディスクID   | -        | 文字列                | [注1](#注1) |
| `hostname`        | -   | ホスト名               | - | 文字列 | ディスク修正機能で設定される、ホスト名 [注2](#注2)|
| `password`        | -   | パスワード               | - | 文字列 | ディスク修正機能で設定される、OS管理者パスワード [注2](#注2)|
| `ssh_key_ids`     | -   | SSH公開鍵ID             | - | リスト(文字列) | ディスク修正機能で設定される、SSH認証用の公開鍵ID [注2](#注2)|
| `disable_pw_auth` | -   | パスワードでの認証無効化   | - | `true`<br />`false` | ディスク修正機能で設定される、SSH接続でのパスワード/チャレンジレスポンス認証の無効化 [注2](#注2)|
| `note_ids`        | -   | スタートアップスクリプトID | - | リスト(文字列) | スタートアップスクリプトのID |
| `icon_id`         | -   | アイコンID         | - | 文字列 | - |
| `description`     | -   | 説明  | - | 文字列 | - |
| `tags`            | -   | タグ | - | リスト(文字列) | - |
| `graceful_shutdown_timeout` | - | シャットダウンまでの待ち時間 | - | 数値(秒数) | シャットダウンが必要な場合の通常シャットダウンするまでの待ち時間(指定の時間まで待ってもシャットダウンしない場合は強制シャットダウンされる) |
| `zone`            | -   | ゾーン | - | `is1b`<br />`tk1a`<br />`tk1v` | - |


#### 互換性

- `source_archive_name`/`source_disk_name`パラメータはv0.3.6にて廃止されました。
v0.3.6以降では[DataResource](data_resource.md)を利用してください。

- `plan`パラメータはv0.5.0にて変更されました。
    - 旧:数値(`2`(HDD) or `4`(SSD))を指定
    - 新:文字列(`hdd` or `ssd`)を指定


#### 注1

`source_archive_id`/`source_disk_id`はいずれか1つだけ指定可能です。

#### 注2

  - OSによりディスク修正機能に対応していない場合があります。
  - これらの値は投入専用です。属性においても投入値を表します(さくらのクラウドAPIからは取得できない項目です)。

### 属性

|属性名                | 名称                    | 補足                                        |
|---------------------|------------------------|--------------------------------------------|
| `id`                | ディスクID               | -                                          |
| `server_id`         | サーバID               | 接続されているサーバのID                     |

