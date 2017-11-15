# Terraform for さくらのクラウド

[![Build Status](https://travis-ci.org/sacloud/terraform-provider-sakuracloud.svg?branch=master)](https://travis-ci.org/sacloud/terraform-provider-sakuracloud)

Terraformからさくらのクラウドを操作するためのプラグインです。  
このプラグインは`さくらインターネット公認ツール`としてさくらのクラウドユーザコミュニティによって開発されています。

## クイックスタート

さくらのクラウドに以下の環境を構築します。

  - 最新安定版のCentOSを利用
  - ディスク: SSH/20GB, サーバー: 1core/1GBメモリ(デフォルト値のため定義ファイル上では省略)
  - SSH接続時のパスワード/チャレンジレスポンス認証を無効化(公開鍵認証のみに)
  - SSH用の公開鍵はさくらのクラウド上で生成(作成された秘密鍵はローカルマシンへ保存する)

[Installation / インストール](https://sacloud.github.io/terraform-provider-sakuracloud/installation/)を参考に
TerraformとTerraform for さくらのクラウドを手元のマシンにインストールしてください。

インストール後、以下のコマンドを実行することでインフラ構築が行われます。

```bash
#################################################
# さくらのクラウドAPIキーを環境変数に設定
#################################################
export SAKURACLOUD_ACCESS_TOKEN=[さくらのクラウド APIトークン]
export SAKURACLOUD_ACCESS_TOKEN_SECRET=[さくらのクラウド APIシークレット]

#################################################
# Terraform定義ファイル作成
#################################################
mkdir work; cd work
tee sakura.tf <<-'EOF'
# サーバーの管理者パスワードの定義
variable "password" {
  default = "PUT_YOUR_PASSWORD_HERE"
}

# 対象ゾーンを指定
provider sakuracloud {
  zone = "tk1a" # 東京第1ゾーン 
}

# 公開鍵(さくらのクラウド上で生成)
resource "sakuracloud_ssh_key_gen" "key" {
  name = "foobar"

  provisioner "local-exec" {
    command = "echo \"${self.private_key}\" > id_rsa; chmod 0600 id_rsa"
  }

  provisioner "local-exec" {
    when    = "destroy"
    command = "rm -f id_rsa"
  }
}

# パブリックアーカイブ(OS)のID参照用のデータソース定義
data sakuracloud_archive "centos" {
  os_type = "centos"
}

# ディスク定義
resource "sakuracloud_disk" "disk01" {
  name              = "disk01"
  source_archive_id = "${data.sakuracloud_archive.centos.id}"
  ssh_key_ids       = ["${sakuracloud_ssh_key_gen.key.id}"]
  password          = "${var.password}"
  disable_pw_auth   = true
}

# サーバー定義
resource "sakuracloud_server" "server01" {
  name  = "server01"
  disks = ["${sakuracloud_disk.disk01.id}"]
}

# サーバへのSSH接続を表示するアウトプット定義
output "ssh_to_server" {
  value = "ssh -i id_rsa root@${sakuracloud_server.server01.ipaddress}"
}
EOF

#################################################
# インフラ構築(init & apply)
#################################################
terraform init
terraform apply
```

## ドキュメント

* Terraform for さくらのクラウド ドキュメント
    * https://sacloud.github.io/terraform-provider-sakuracloud/

### サポートしているリソース/データソース

#### リソース
  - [サーバー](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/server/)
  - [ディスク](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/disk/)
  - [アーカイブ](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/archive/)
  - [ISOイメージ(CD-ROM)](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/cdrom/)
  - [スイッチ](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/switch/)
  - [ルーター](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/internet/)
  - [サブネット](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/subnet/)
  - [パケットフィルタ](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/packet_filter/)
  - [パケットフィルタ(ルール)](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/packet_filter_rule/)
  - [ブリッジ](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/bridge/)
  - [ロードバランサー](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/load_balancer/)
  - [VPCルーター](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/vpc_router/)
  - [データベース](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/database/)
  - [NFS](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/nfs/)
  - [スタートアップスクリプト](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/note/)
  - [公開鍵](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/ssh_key/)
  - [公開鍵(生成)](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/ssh_key_gen/)
  - [アイコン](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/icon/)
  - [専有ホスト](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/private_host/)
  - [DNS](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/dns/)
  - [GSLB](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/gslb/)
  - [シンプル監視](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/simple_monitor/)
  - [自動バックアップ](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/auto_backup/)
  - [オブジェクトストレージ](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/bucket_object/)
  - [サーバ コネクタ](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/server_connector)

#### データソース
  - [データソース](http://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/data_resource/)

#### サポートしていないリソース

以下のリソースはさくらのクラウド側でAPIが提供されていないため未サポートです。

  - ローカルルータ
  - ウェブアクセラレータ
  - オブジェクトストレージ(バケット作成)
  - ライセンス
  - 割引パスポート
  - クーポン

## Building/Developing

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

    # ドキュメントのプレビュー用サーバー起動(http://localhost/でプレビュー可能)
    make serve-docs
    
    # ドキュメントの検証(textlint)
    make lint-docs
    
    # build_docs配下のファイルからドキュメント生成(docsディレクトリ再生成)
    make build-docs

## License

  This project is published under [Apache 2.0 License](LICENSE).

## Author

  * Kazumichi Yamamoto ([@yamamoto-febc](https://github.com/sacloud))
