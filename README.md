# Terraform for さくらのクラウド

[![Join the chat at https://gitter.im/terraform-provider-sakuracloud/Lobby](https://badges.gitter.im/terraform-provider-sakuracloud/Lobby.svg)](https://gitter.im/terraform-provider-sakuracloud/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

[![Build Status](https://travis-ci.org/yamamoto-febc/terraform-provider-sakuracloud.svg?branch=master)](https://travis-ci.org/yamamoto-febc/terraform-provider-sakuracloud)

Terraformでさくらのクラウドを操作するためのプラグイン



## クイックスタート

#### 準備

  - Dockerをインストールしておく
  - さくらのクラウドAPIキーを取得しておく

Dockerがない場合は[Installation / インストール](docs/installation.md)を参考に
TerraformとTerraform for さくらのクラウドを手元のマシンにインストールしてからご利用ください。

さくらのクラウドAPIキーの取得方法は[こちら](docs/installation.md#さくらのクラウドapiキーの取得)を参照してください。

```bash
#################################################
# Terraform定義ファイル作成
#################################################
$ mkdir ~/work; cd ~/work
$ ssh-keygen -C "" -P "" -f id_rsa   # サーバーへのSSH用キーペア生成
$ tee sakura.tf <<-'EOF'
resource "sakuracloud_ssh_key" "key"{
    name = "sshkey"
    public_key = "${file("id_rsa.pub")}"
}

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
    zone = "is1b"
}

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
$ docker run -it --rm \
         -e SAKURACLOUD_ACCESS_TOKEN=[さくらのクラウド APIトークン] \
         -e SAKURACLOUD_ACCESS_TOKEN_SECRET=[さくらのクラウド APIシークレット] \
         -v $PWD:/work \
         sacloud/terraform apply
```

## インストール

[リリースページ](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/releases/latest)から最新のバイナリを取得し、
Terraformバイナリと同じディレクトリに展開してください。

詳細は[Installation / インストール](docs/installation.md)を参照してください。

## 使い方/各リソースの設定方法

Terraform定義ファイル(tfファイル)を作成してご利用ください。

設定ファイルの記載方法は[リファレンス](docs/configuration.md)を参照ください。

さくらのクラウドの以下のリソースをサポートしています。

### サポートしているリソース

  - [サーバー](docs/configuration/resources/server.md)
  - [ディスク](docs/configuration/resources/disk.md)
  - [スイッチ](docs/configuration/resources/switch.md)
  - [ルーター](docs/configuration/resources/internet.md)
  - [パケットフィルタ](docs/configuration/resources/packet_filter.md)
  - [ブリッジ](docs/configuration/resources/bridge.md)
  - [ロードバランサー](docs/configuration/resources/load_balancer.md)
  - [VPCルーター](docs/configuration/resources/vpc_router.md)
  - [データベース](docs/configuration/resources/database.md)
  - [スタートアップスクリプト](docs/configuration/resources/note.md)
  - [公開鍵](docs/configuration/resources/ssh_key.md)
  - [DNS](docs/configuration/resources/dns.md)
  - [GSLB](docs/configuration/resources/gslb.md)
  - [シンプル監視](docs/configuration/resources/simple_monitor.md)
  - [自動バックアップ](docs/configuration/resources/auto_backup.md)


## Building/Developing

#### ソースコード
    
[Terraform](https://github.com/hashicorp/terraform)本体でのプロバイダ配置に合わせ、`builtin`ディレクトリ配下にソースを配置しています。
    
     builtin/
       ├── bins
       │     └── provider-sakuracloud  # Terraformプラグインエントリーポイント(mainパッケージ)
       └── providers
             └── sakuracloud           # さくらのクラウド用プロバイダ/リソースなどソース一式

#### ビルド

    make build
    
#### ビルド(クロスコンパイル)

    make build-x
    
#### ビルド(Docker上でのビルド)

    make docker-build
    
#### テスト

    make test
    
#### 受入テスト(実際のさくらのクラウドAPI呼び出しを伴うテスト)

    make testacc
    
#### 依存ライブラリ

    # 一覧表示
    govendor list
    
    # vendor配下のライブラリを一括更新
    govendor fetch +v

    # vendor配下のライブラリをGOPATH上から更新
    govendor update +v

## License

  This project is published under [Apache 2.0 License](LICENSE).

## Author

  * Kazumichi Yamamoto ([@yamamoto-febc](https://github.com/yamamoto-febc))
