# データリソース

---

データリソース(Data Resource)とは、読み取り専用のリソースです。
すでにさくらのクラウド上に存在するリソースの値を参照するために用います。

以下の例ではディスクのコピー元アーカイブのIDを参照するために
`sakuracloud_archive`データリソースを利用しています。

データリソースを利用することで`sakuracloud_disk`の定義中にアーカイブのIDを直接指定しないように出来ます。

### 利用例

```hcl
data "sakuracloud_archive" "ubuntu" {
  filter {
    name   = "Name"
    values = ["Ubuntu Server"]
  }
  filter {
    name   = "Tags"
    values = ["current-stable", "arch-64bit", "os-linux"]
  }
}

resource "sakuracloud_disk" "disk01" {
  name              = "disk01"
  source_archive_id = data.sakuracloud_archive.ubuntu.id
}
```

### パラメーター(アーカイブのみ)

アーカイブ(`sakuracloud_archive`)リソースでは、`os_type`パラメーターが利用可能です。

```hcl
data "sakuracloud_archive" "ubuntu" {
  os_type = "ubuntu" # Ubuntuの最新安定版パブリックアーカイブ
}

data "sakuracloud_archive" "centos" {
  os_type = "centos" # CentOSの最新安定版パブリックアーカイブ
}

```

詳細は[アーカイブデータリソースのドキュメント](data/archive/)を参照してください。
