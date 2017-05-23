# データソース

---

データソース(Data Resource)とは、読み取り専用のリソースです。
すでにさくらのクラウド上に存在するリソースの値を参照するために用います。

以下の例ではディスクのコピー元アーカイブのIDを参照するために
`sakuracloud_archive`データソースを利用しています。

データソースを利用することで`sakuracloud_disk`の定義中にアーカイブのIDを直接指定しないように出来ます。

### 利用例

```hcl
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

### パラメーター(アーカイブのみ)

アーカイブ(`sakuracloud_archive`)リソースでは、`os_type`パラメーターが利用可能です。

```hcl
data sakuracloud_archive "ubuntu" {
    os_type = "ubuntu" # Ubuntuの最新安定版パブリックアーカイブ
}

data sakuracloud_archive "centos" {
    os_type = "centos" # CentOSの最新安定版パブリックアーカイブ
}

```

`os_type`に指定可能な値は以下の通りです。

|値|詳細                                          |
|---------------------------|--------------------|
| `centos`                  | CentOS 7|
| `ubuntu`                  | Ubuntu 16.04|
| `debian`                  | Debian |
| `vyos`                    | VyOS|
| `coreos`                  | CoreOS|
| `rancheros`               | RancherOS|
| `kusanagi`                | Kusanagi(CentOS7)|
| `site-guard`              | SiteGuard(CentOS7)|
| `plesk`                   | Plesk(CentOS7)|
| `freebsd`                 | FreeBSD|
| `windows2012`             | Windows 2012|
| `windows2012-rds`         | Windows 2012(RDS)|
| `windows2012-rds-office`  | Windows 2012(RDS + Office)|
| `windows2016`             | Windows 2016|
| `windows2016-rds`         | Windows 2016(RDS)|
| `windows2016-rds-office`  | Windows 2016(RDS + Office)|
| `windows2016-sql-web`     | Windows 2016 SQLServer(Web)|
| `windows2016-sql-standard`| Windows 2016 SQLServer(Standard)|

`os_type`に対応していないパブリックアーカイブについては`filter`パラメーターをご利用ください。

### サポートしているデータソース

データソースそれぞれの属性はリソースごとに異なります。
詳細は各リソースのドキュメントを参照してください。

|データソース                   | 名称                    | 補足                                        |
|------------------------------|------------------------|--------------------------------------------|
| `sakuracloud_archive`        | アーカイブ               | -                                          |
| `sakuracloud_bridge`         | ブリッジ                | -                                          |
| `sakuracloud_cdrom`          | ISOイメージ             | -                                          |
| `sakuracloud_database`       | データベース            | -                                          |
| `sakuracloud_disk`           | ディスク                | -                                          |
| `sakuracloud_dns`            | DNS                    | -                                          |
| `sakuracloud_gslb`           | GSLB                   | -                                          |
| `sakuracloud_internet`       | ルータ                | -                                          |
| `sakuracloud_note`           | スタートアップスクリプト   | -                                          |
| `sakuracloud_packet_filter`  | パケットフィルタ         | -                                          |
| `sakuracloud_server`         | サーバ                | -                                          |
| `sakuracloud_simple_monitor` | シンプル監視            | -                                          |
| `sakuracloud_ssh_key`        | 公開鍵                 | -                                          |
| `sakuracloud_subnet`         | サブネット              | -                                          |
| `sakuracloud_switch`         | スイッチ                | -                                          |
