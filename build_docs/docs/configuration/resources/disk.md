# ディスク(sakuracloud_disk)

---

### 設定例

```hcl
data sakuracloud_archive "centos" {
    os_type = "centos"
}
resource "sakuracloud_disk" "disk01"{
    name = "disk01"
    source_archive_id = "${data.sakuracloud_archive.centos.id}"
    ssh_key_ids = ["${sakuracloud_ssh_key.key.id}"]
    disable_pw_auth = true
}

```

### パラメーター

|パラメーター         |必須  |名称                |初期値     |設定値                    |補足                                          |
|-------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `name`            | ◯   | ディスク名           | -        | 文字列                  | - |
| `plan`            | -   | ディスクプラン        | `ssd` | `ssd`<br />`hdd` | - |
| `connection`      | -   | ディスク接続          | `virtio` | `virtio`<br />`ide`    | - |
| `size`            | -   | ディスクサイズ(GB単位) | 20       | 数値                    | - |
|`source_archive_id`| -   | コピー元アーカイブID   | -        | 文字列                | [注1](#注1) |
|`source_disk_id`   | -   | コピー元ディスクID   | -        | 文字列                | [注1](#注1) |
| `hostname`        | -   | ホスト名               | - | 文字列 | ディスク修正機能で設定される、ホスト名 [注2](#注2)|
| `password`        | -   | パスワード               | - | 文字列 | ディスク修正機能で設定される、OS管理者パスワード [注2](#注2)|
| `ssh_key_ids`     | -   | SSH公開鍵ID             | - | リスト(文字列) | ディスク修正機能で設定される、SSH認証用の公開鍵ID [注2](#注2)|
| `disable_pw_auth` | -   | パスワードでの認証無効化   | - | `true`<br />`false` | ディスク修正機能で設定される、SSH接続でのパスワード/チャレンジレスポンス認証の無効化 [注2](#注2)|
| `note_ids`        | -   | スタートアップスクリプトID | - | リスト(文字列) | スタートアップスクリプトのID |
| `description`     | -   | 説明  | - | 文字列 | - |
| `tags`            | -   | タグ | - | リスト(文字列) | - |
| `zone`            | -   | ゾーン | - | `is1b`<br />`tk1a`<br />`tk1v` | - |


#### 互換性

- `source_archive_name`/`source_disk_name`パラメータはv0.3.6にて廃止されました。
v0.3.6以降では[DataResource](data_resource.md)を利用してください。

- `plan`パラメータはv0.5.0にて変更されました。
    - 旧:数値(`2`(HDD) or `4`(SSD))を指定
    - 新:文字列(`hdd` or `ssd`)を指定


#### 注1

`source_archive_id`/`source_disk_id`はいずれか一つだけ指定可能です。

#### 注2

  - OSによりディスク修正機能に対応していない場合があります。
  - これらの値は投入専用です。属性においても投入値を表します(さくらのクラウドAPIからは取得できない項目です)。

### 属性

|属性名                | 名称                    | 補足                                        |
|---------------------|------------------------|--------------------------------------------|
| `id`                | ディスクID               | -                                          |
| `name`              | ディスク名               | -                                          |
| `plan`              | ディスクプラン            | -                                          |
| `connection`        | ディスク接続             | -                                          |
| `size`              | ディスクサイズ(GB単位)    | -                                          |
|`source_archive_id`  | コピー元アーカイブID      | -                                          |
|`source_disk_id`     | コピー元ディスクID        | -                                          |
| `hostname`          | ホスト名                | -                                          |
| `password`          | パスワード               | -                                          |
| `ssh_key_ids`       | SSH公開鍵ID             | -                                          |
| `disable_pw_auth`   | パスワードでの認証無効化   | -                                          |
| `note_ids`          | スタートアップスクリプトID | -                                          |
| `description`       | 説明                    | -                                          |
| `tags`              | タグ                    | -                                          |
| `zone`              | ゾーン                  | -                                          |
| `server_id`         | サーバーID               | 接続されているサーバーのID                     |

