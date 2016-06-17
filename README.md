# Terraform for さくらのクラウド

Terraformでさくらのクラウドを操作するためのプラグイン

## クイックスタート

```bash
$ TERRAFORM_URL="https://releases.hashicorp.com/terraform/0.6.16/terraform_0.6.16_darwin_amd64.zip"
$ TERRAFORM_SAKURA_URL="https://github.com/yamamoto-febc/terraform-provider-sakuracloud/releases/download/v0.3.4-RC/terraform-provider-sakuracloud_linux-amd64.zip"

#################################################
# Terraform インストール
#################################################
$ mkdir -p ~/terraform ; cd ~/terraform
$ curl -L $TERRAFORM_URL > terraform.zip
$ unzip terraform.zip
$ export PATH=$PATH:~/terraform/

#################################################
# Terraform for さくらのクラウド インストール
#################################################
$ curl -L $TERRAFORM_SAKURA_URL > terraform-provider-sakuracloud.zip
$ unzip terraform-provider-sakuracloud.zip

#################################################
# APIキー設定(コントロールパネルで取得しておいてください)
#################################################
$ export SAKURACLOUD_ACCESS_TOKEN="さくらのクラウド APIトークン"
$ export SAKURACLOUD_ACCESS_TOKEN_SECRET="さくらのクラウド APIシークレット"

#################################################
# Terraform定義ファイル作成
#################################################
$ mkdir work; cd work
$ ssh-keygen -C "" -P "" -f id_rsa   # サーバーへのSSH用キーペア生成
$ tee sakura.tf <<-'EOF'

# 公開鍵
resource "sakuracloud_ssh_key" "key"{
    name = "sshkey"
    public_key = "${file("id_rsa.pub")}"
}

# ディスク
resource "sakuracloud_disk" "disk01"{
    name = "disk01"
    source_archive_name = "CentOS 7.2 64bit"
    ssh_key_ids = ["${sakuracloud_ssh_key.key.id}"]
    disable_pw_auth = true
    zone = "is1b"
}

#サーバー
resource "sakuracloud_server" "server01" {
    name = "server01"
    disks = ["${sakuracloud_disk.disk01.id}"]
    tags = ["@virtio-net-pci"]
    zone = "is1b"
}

EOF

#################################################
# Terraformでインフラ作成
#################################################
$ terraform apply
```

## インストール

[リリースページ](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/releases/latest)から最新のバイナリを取得し、
Terraformバイナリと同じディレクトリに展開してください。

詳細は[Installation / インストール](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/wiki/Installation)を参照してください。

## 使い方/各リソースの設定方法

Terraform定義ファイル(tfファイル)を作成してご利用ください。
設定ファイルの記載方法は[Wikiページ](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/wiki)を参照ください。
さくらのクラウドの以下のリソースをサポートしています。

### サポートしているリソース

  - [サーバー](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/wiki/Configuration-Resource-Server)
  - [ディスク](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/wiki/Configuration-Resource-Disk)
  - [スイッチ](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/wiki/Configuration-Resource-Switch)
  - [ルーター](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/wiki/Configuration-Resource-Internet)
  - [パケットフィルタ](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/wiki/Configuration-Resource-PacketFilter)
  - [ブリッジ](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/wiki/Configuration-Resource-Bridge)
  - [スタートアップスクリプト](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/wiki/Configuration-Resource-Note)
  - [公開鍵](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/wiki/Configuration-Resource-SSHKey)
  - [DNS](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/wiki/Configuration-Resource-DNS)
  - [GSLB](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/wiki/Configuration-Resource-GSLB)
  - [シンプル監視](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/wiki/Configuration-Resource-SimpleMonitor)


## Building/Developing

  `godep get $(go list ./... | grep -v vendor)`

  `godep restore`

  `godep go test .`

  `TF_ACC=1 godep go test -v -timeout=60m .` run acceptance tests. (requires ENV vars)

  `godep go build -o path/to/desired/terraform-provider-sakuracloud bin/terraform-provider-sakuracloud/main.go`


## License

  This project is published under [Apache 2.0 License](LICENSE).

## Author

  * Kazumichi Yamamoto ([@yamamoto-febc](https://github.com/yamamoto-febc))
