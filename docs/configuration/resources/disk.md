# ディスク(sakuracloud_disk)

### 設定例

```
data sakuracloud_archive "centos" {
    filter = {
        name   = "Tags"
        values = ["current-stable", "arch-64bit", "distro-centos"]
    }
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
| `plan`            | -   | ディスクプラン        | `4`(SSD) | `2`(HDD)<br />`4`(SSD) | - |
| `connection`      | -   | ディスク接続          | `virtio` | `virtio`<br />`ide`    | - |
| `size`            | -   | ディスクサイズ(GB単位) | 20       | 数値                    | - |
|`source_archive_id`| -   | コピー元アーカイブID   | -        | 文字列                | [注1](#注1) |
|`source_disk_id`   | -   | コピー元ディスクID   | -        | 文字列                | [注1](#注1) |
| `password`        | -   | パスワード               | - | 文字列 | ディスク修正機能で設定されるOS管理者パスワード [注2](#注2)|
| `ssh_key_ids`     | -   | SSH公開鍵ID             | - | リスト(文字列) | ディスク修正機能で設定されるSSH認証用の公開鍵ID [注2](#注2)|
| `disable_pw_auth` | -   | パスワードでの認証無効化   | - | `true`<br />`false` | SSH接続でのパスワード/チャレンジレスポンス認証の無効化 [注2](#注2)|
| `user_ip_address` | -   | IPアドレス | - | 文字列 | eth0のIPアドレス [注2](#注2)|
| `default_route`   | -   | デフォルトゲートウェイ | - | 文字列 | eth0のデフォルトゲートウェイ（user_ip_addressと同時指定時のみ有効） [注2](#注2)|
| `network_mask_len`| -   | ネットワークマスク長 | - | 数値 | eth0のネットワークマスク長（user_ip_addressと同時指定時のみ有効） [注2](#注2)|
| `note_ids`        | -   | スタートアップスクリプトID | - | リスト(文字列) | スタートアップスクリプトのID |
| `description`     | -   | 説明  | - | 文字列 | - |
| `tags`            | -   | タグ | - | リスト(文字列) | - |
| `zone`            | -   | ゾーン | - | `is1b`<br />`tk1a`<br />`tk1v` | - |

#### 互換性

`source_archive_name`/`source_disk_name`パラメータはv0.3.6にて廃止されました。
v0.3.6以降では[DataResource](data_resource.md)を利用してください。

#### 注1

`source_archive_id`/`source_disk_id`はいずれか一つだけ指定可能です。

#### 注2

OSによりディスク修正機能に対応していない場合があります。

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
| `password`          | パスワード               | -                                          |
| `ssh_key_ids`       | SSH公開鍵ID             | -                                          |
| `disable_pw_auth`   | パスワードでの認証無効化   | -                                          |
| `user_ip_address`   | IPアドレス                 | -                                          |
| `default_route`     | デフォルトゲートウェイ     | -                                          |
| `network_mask_len`  | ネットワークマスク長       | -                                          |
| `note_ids`          | スタートアップスクリプトID | -                                          |
| `description`       | 説明                    | -                                          |
| `tags`              | タグ                    | -                                          |
| `zone`              | ゾーン                  | -                                          |
| `server_id`         | サーバーID               | 接続されているサーバーのID                     |

