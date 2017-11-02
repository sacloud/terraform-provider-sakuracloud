# Terraform for さくらのクラウド

---

## 概要

`Terraform for さくらのクラウド`とは、
[Terraform](https://terraform.io)から[さくらのクラウド](http://cloud.sakura.ad.jp)を操作するためのTerraform用プラグインです。  


## 目次

#### ドキュメント
- [Installation/インストールガイド](installation/)
- [Articles/ブログ/記事/読み物](articles/)

####  Reference/設定リファレンス
- 全体/共通項目
    - [プロバイダ](configuration/provider/)
- リソース
    - [サーバ](configuration/resources/server/)
    - [ディスク](configuration/resources/disk/)
    - [アーカイブ](configuration/resources/archive/)
    - [ISOイメージ(CD-ROM)](configuration/resources/cdrom/)
    - [スイッチ](configuration/resources/switch/)
    - [ルータ](configuration/resources/internet/)
    - [サブネット](configuration/resources/subnet/)
    - [パケットフィルタ](configuration/resources/packet_filter/)
    - [パケットフィルタ(ルール)](configuration/resources/packet_filter_rule/)
    - [ブリッジ](configuration/resources/bridge/)
    - [ロードバランサ](configuration/resources/load_balancer/)
    - [VPCルータ](configuration/resources/vpc_router/)
    - [データベース](configuration/resources/database/)
    - [NFS](configuration/resources/nfs/)
    - [スタートアップスクリプト](configuration/resources/note/)
    - [公開鍵](configuration/resources/ssh_key/)
    - [公開鍵(生成)](configuration/resources/ssh_key_gen/)
    - [アイコン](configuration/resources/icon/)
    - [DNS](configuration/resources/dns/)
    - [GSLB](configuration/resources/gslb/)
    - [シンプル監視](configuration/resources/simple_monitor/)
    - [自動バックアップ](configuration/resources/auto_backup/)
    - [オブジェクトストレージ](configuration/resources/bucket_object/)
- 特殊なリソース
    - [サーバ コネクタ](configuration/resources/server_connector)
- データソース:
    - [データソース](configuration/resources/data_resource/)
  