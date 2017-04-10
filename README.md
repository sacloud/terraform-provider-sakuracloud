# Terraform for さくらのクラウド

[![Build Status](https://travis-ci.org/yamamoto-febc/terraform-provider-sakuracloud.svg?branch=master)](https://travis-ci.org/yamamoto-febc/terraform-provider-sakuracloud)

Terraformでさくらのクラウドを操作するためのプラグイン



## クイックスタート

#### 準備

  - Dockerをインストールしておく
  - さくらのクラウドAPIキーを取得しておく

Dockerがない場合は[Installation / インストール](https://yamamoto-febc.github.io/terraform-provider-sakuracloud/installation/)を参考に
TerraformとTerraform for さくらのクラウドを手元のマシンにインストールしてからご利用ください。

さくらのクラウドAPIキーの取得方法は[こちら](https://yamamoto-febc.github.io/terraform-provider-sakuracloud/installation/#api)を参照してください。

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
    os_type = "centos"
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

## ドキュメント

* Terraform for さくらのクラウド ドキュメント
    * https://yamamoto-febc.github.io/terraform-provider-sakuracloud/

## インストール

[リリースページ](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/releases/latest)から最新のバイナリを取得し、
Terraformバイナリと同じディレクトリに展開してください。

詳細は[Installation / インストール](https://yamamoto-febc.github.io/terraform-provider-sakuracloud/installation/)を参照してください。

## 使い方/各リソースの設定方法

Terraform定義ファイル(tfファイル)を作成してご利用ください。

設定ファイルの記載方法は[リファレンス](https://yamamoto-febc.github.io/terraform-provider-sakuracloud/#_2)を参照ください。

さくらのクラウドの以下のリソースをサポートしています。

### サポートしているリソース

  - [サーバー](https://yamamoto-febc.github.io/terraform-provider-sakuracloud/configuration/resources/server/)
  - [ディスク](https://yamamoto-febc.github.io/terraform-provider-sakuracloud//configuration/resources/disk/)
  - [スイッチ](https://yamamoto-febc.github.io/terraform-provider-sakuracloud//configuration/resources/switch/)
  - [ルーター](https://yamamoto-febc.github.io/terraform-provider-sakuracloud//configuration/resources/internet/)
  - [パケットフィルタ](https://yamamoto-febc.github.io/terraform-provider-sakuracloud/configuration/resources/packet_filter/)
  - [ブリッジ](https://yamamoto-febc.github.io/terraform-provider-sakuracloud/configuration/resources/bridge/)
  - [ロードバランサー](https://yamamoto-febc.github.io/terraform-provider-sakuracloud/configuration/resources/load_balancer/)
  - [VPCルーター](https://yamamoto-febc.github.io/terraform-provider-sakuracloud/configuration/resources/vpc_router/)
  - [データベース](https://yamamoto-febc.github.io/terraform-provider-sakuracloud/configuration/resources/database/)
  - [スタートアップスクリプト](https://yamamoto-febc.github.io/terraform-provider-sakuracloud/configuration/resources/note/)
  - [公開鍵](https://yamamoto-febc.github.io/terraform-provider-sakuracloud/configuration/resources/ssh_key/)
  - [DNS](https://yamamoto-febc.github.io/terraform-provider-sakuracloud/configuration/resources/dns/)
  - [GSLB](https://yamamoto-febc.github.io/terraform-provider-sakuracloud/configuration/resources/gslb/)
  - [シンプル監視](https://yamamoto-febc.github.io/terraform-provider-sakuracloud/configuration/resources/simple_monitor/)
  - [自動バックアップ](https://yamamoto-febc.github.io/terraform-provider-sakuracloud/configuration/resources/auto_backup/)
  - [データソース](http://yamamoto-febc.github.io/terraform-provider-sakuracloud/configuration/resources/data_resource/)


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

#### ドキュメント

ドキュメントはGithub Pagesを利用しています。(masterブランチの`docs`ディレクトリ配下)  
静的ファイルの生成は`mkdocs`コマンドで行なっています。  
ドキュメントの追加や修正は`build_docs`ディレクトリ以下のファイルの追加/修正を行なった上で`mkdocs`コマンドでファイル生成してコミットしてください。

    # build_docs配下のファイルからドキュメント生成(docsディレクトリ再生成)
    make build-docs
    
    # ドキュメントのプレビュー用サーバー起動(http://localhost/でプレビュー可能)
    make serve-docs

## License

  This project is published under [Apache 2.0 License](LICENSE).

## Author

  * Kazumichi Yamamoto ([@yamamoto-febc](https://github.com/yamamoto-febc))
