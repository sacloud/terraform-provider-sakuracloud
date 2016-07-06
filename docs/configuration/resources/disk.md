# ディスク(sakuracloud_disk)

### 設定例

```
resource "sakuracloud_disk" "mydisk"{
    name = "mydisk"
    size = 20
    source_archive_name = "CentOS 7.2 64bit"
    description = "Disk from terraform for SAKURA CLOUD"
    tags = ["hoge1" , "hoge2"]
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
|`source_archive_name`| - | コピー元アーカイブ名   | -        | 文字列              | 指定の名前に部分一致するアーカイブのうち、ヒットした先頭1件が利用されます。<br />[注1](#注1) |
|`source_disk_id`   | -   | コピー元ディスクID   | -        | 文字列                | [注1](#注1) |
|`source_disk_name` | -   | コピー元ディスク名   | -        | 文字列                | 指定の名前に部分一致するディスクのうち、ヒットした先頭1件が利用されます。<br />[注1](#注1) |
| `password`        | -   | パスワード               | - | 文字列 | ディスク修正機能で設定されるOS管理者パスワード [注2](#注2)|
| `ssh_key_ids`     | -   | SSH公開鍵ID             | - | リスト(文字列) | ディスク修正機能で設定されるSSH認証用の公開鍵ID [注2](#注2)|
| `disable_pw_auth` | -   | パスワードでの認証無効化   | - | `true`<br />`false` | SSH接続でのパスワード/チャレンジレスポンス認証の無効化 [注2](#注2)|
| `note_ids`        | -   | スタートアップスクリプトID | - | リスト(文字列) | スタートアップスクリプトのID |
| `description`     | -   | 説明  | - | 文字列 | - |
| `tags`            | -   | タグ | - | リスト(文字列) | - |
| `zone`            | -   | ゾーン | - | `is1b`<br />`tk1a`<br />`tk1v` | - |

#### 注1

`source_archive_id`/`source_archive_name`/`source_disk_id`/`source_disk_name`はいずれか一つだけ指定可能です。

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
|`source_archive_name`| コピー元アーカイブ名      | -                                          |
|`source_disk_id`     | コピー元ディスクID        | -                                          |
|`source_disk_name`   | コピー元ディスク名        | -                                          |
| `password`          | パスワード               | -                                          |
| `ssh_key_ids`       | SSH公開鍵ID             | -                                          |
| `disable_pw_auth`   | パスワードでの認証無効化   | -                                          |
| `note_ids`          | スタートアップスクリプトID | -                                          |
| `description`       | 説明                    | -                                          |
| `tags`              | タグ                    | -                                          |
| `zone`              | ゾーン                  | -                                          |
| `server_id`         | サーバーID               | 接続されているサーバーのID                     |

