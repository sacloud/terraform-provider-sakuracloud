# データソース

データソース(Data Resource)とは、読み取り専用のリソースです。
すでにさくらのクラウド上に存在するリソースの値を参照するために用います。

以下の例ではディスクのコピー元アーカイブのIDを参照するために
`sakuracloud_archive`データソースを利用しています。

データソースを利用することで`sakuracloud_disk`の定義中にアーカイブのIDを直接指定しないように出来ます。

### 利用例

```
data sakuracloud_archive "ubuntu" {
    filter = {
        name = "Name"
        values = ["Ubuntu Server"]
    }
    filter = {
        name = "Tags"
        values = ["current-stable","arch-64bit","os-linux"]
    }
}

resource sakuracloud_disk "disk01"{
    name = "disk01"
    source_archive_id = "${data.sakuracloud_archive.ubuntu.id}"
}
```

### パラメーター(全データソース共通)

|パラメーター|必須  |名称                |初期値     |設定値 |補足                                          |
|----------|:---:|---------------------|:--------:|------|----------------------------------------------|
| `filter` | -   | 検索条件             | -        | -    | `name`に検索対象の属性名、`values`に検索値をリストで指定します。<br />`values`に複数の値が設定されている場合、AND条件となります。 |
| `zone`   | -   | ゾーン               | -        | `is1b`<br />`tk1a`<br />`tk1v` | - |

`filter`の`name`属性に指定可能な値は[さくらのクラウド APIドキュメント](http://developer.sakura.ad.jp/cloud/api/1.1/)を参照ください。

### サポートしているデータソース

データソースそれぞれの属性はリソースごとに異なります。
詳細は各リソースのドキュメントを参照してください。

|データソース                   | 名称                    | 補足                                        |
|------------------------------|------------------------|--------------------------------------------|
| `sakuracloud_archive`        | アーカイブ               | -                                          |
| `sakuracloud_bridge`         | ブリッジ                | -                                          |
| `sakuracloud_cdrom`          | ISOイメージ             | -                                          |
| `sakuracloud_disk`           | ディスク                | -                                          |
| `sakuracloud_dns`            | DNS                    | -                                          |
| `sakuracloud_gslb`           | GSLB                   | -                                          |
| `sakuracloud_internet`       | ルーター                | -                                          |
| `sakuracloud_note`           | スタートアップスクリプト   | -                                          |
| `sakuracloud_packet_filter`  | パケットフィルタ         | -                                          |
| `sakuracloud_server`         | サーバー                | -                                          |
| `sakuracloud_simple_monitor` | シンプル監視            | -                                          |
| `sakuracloud_ssh_key`        | 公開鍵                 | -                                          |
| `sakuracloud_switch`         | スイッチ                | -                                          |
