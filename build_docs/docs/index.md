# Terraform for さくらのクラウド

---

## 概要

`Terraform for さくらのクラウド`とは、
[Terraform](https://terraform.io)から[さくらのクラウド](http://cloud.sakura.ad.jp)を操作するためのTerraform用プラグインです。  


## 目次

#### ドキュメント
- [Installation/インストールガイド](installation/)

#### アップグレードガイド:

- [v0.12での変更点](upgrade_to_v012/)
- [v1.1での変更点](upgrade_to_v11/)
- [v1.4での変更点](upgrade_to_v14/)
- [v1.6での変更点](upgrade_to_v16/)      
- [v1.7での変更点](upgrade_to_v17/)      
- [v1.8での変更点](upgrade_to_v18/)      

#### その他:

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
    - [データベース(リードレプリカ)](configuration/resources/database_read_replica/)
    - [NFS](configuration/resources/nfs/)
    - [SIM](configuration/resources/sim/)
    - [モバイルゲートウェイ](configuration/resources/mobile_gateway/)
    - [スタートアップスクリプト](configuration/resources/note/)
    - [公開鍵](configuration/resources/ssh_key/)
    - [公開鍵(生成)](configuration/resources/ssh_key_gen/)
    - [アイコン](configuration/resources/icon/)
    - [専有ホスト](configuration/resources/private_host/)
    - [DNS](configuration/resources/dns/)
    - [IPv4逆引きレコード](configuration/resources/ipv4_ptr/)
    - [GSLB](configuration/resources/gslb/)
    - [シンプル監視](configuration/resources/simple_monitor/)
    - [自動バックアップ](configuration/resources/auto_backup/)
    - [オブジェクトストレージ](configuration/resources/bucket_object/)
- 特殊なリソース
    - [サーバ コネクタ](configuration/resources/server_connector)
- データリソース:
    - [データリソースとは](configuration/resources/data_resource)
    - [サーバ](configuration/resources/data/server)
    - [ディスク](configuration/resources/data/disk)
    - [アーカイブ](configuration/resources/data/archive)
    - [ISOイメージ(CD-ROM)](configuration/resources/data/cdrom)
    - [スイッチ](configuration/resources/data/switch)
    - [ルータ](configuration/resources/data/internet)
    - [サブネット](configuration/resources/data/subnet)
    - [パケットフィルタ](configuration/resources/data/packet_filter)
    - [ブリッジ](configuration/resources/data/bridge)
    - [ロードバランサ](configuration/resources/data/load_balancer)
    - [VPCルータ](configuration/resources/data/vpc_router/)
    - [データベース](configuration/resources/data/database)
    - [NFS](configuration/resources/data/nfs)
    - [スタートアップスクリプト](configuration/resources/data/note)
    - [公開鍵](configuration/resources/data/ssh_key)
    - [アイコン](configuration/resources/data/icon)
    - [専有ホスト](configuration/resources/data/private_host)
    - [DNS](configuration/resources/data/dns)
    - [GSLB](configuration/resources/data/gslb)
    - [シンプル監視](configuration/resources/data/simple_monitor)
    - [オブジェクトストレージ](configuration/resources/data/bucket_object)
  