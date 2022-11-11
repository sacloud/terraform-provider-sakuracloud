# terraform-plugin-framework対応

- URL: https://github.com/sacloud/terraform-plugin-sakuracloud/pull/931
- Author: @yamamoto-febc

## 概要/背景

Terraformのプラグイン開発には長らく[Terraform Plugin SDK(v2)](https://github.com/hashicorp/terraform-plugin-sdk)が用いられてきたが、
近年より新しい仕組みとして[Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)が登場した。  

- Terraformドキュメント: [Terraform Plugin SDK(v2)](https://www.terraform.io/plugin/sdkv2)
- Terraformドキュメント: [Terraform Plugin Framework](https://www.terraform.io/plugin/framework)

このドキュメントではterraform-provider-sakuracloudがSDKとFrameworkにどう対応していくかを示す。

### Terraformドキュメントの推奨事項

SDKとFrameworkのどちらを使うべきかは以下のドキュメントに記載されている。

- Terraformドキュメント: [Which SDK Should I Use?](https://www.terraform.io/plugin/which-sdk)

このドキュメントによると、新しいリソース/データソースを作る際にはFrameworkを使うことが推奨されている。
> If you maintain a large existing provider, we recommend that you begin migrating from SDKv2 to the framework by developing new resources and data sources with the framework.
 
また、Frameworkにしかない機能を使いたい場合も当然ながらFrameworkを使うことが推奨される。

[terraform-plugin-mux](https://www.terraform.io/plugin/mux)を用いることでSDK/Frameworkの両方を同時に使うという方法も取れるため、移行するのであれば段階的な移行が推奨されている。

### 他のプロバイダーの対応状況

AWSプロバイダーではterraform-plugin-muxを用いた対応が始まっている。
https://github.com/hashicorp/terraform-provider-aws/pull/25606

他のメジャーなプロバイダーではまだ対応は始まっていないが、AWSプロバイダーが対応を始めたのであれば将来的に他のプロバイダーも対応するものと思われる。

## 対応方針

- terraform-plugin-muxを用いてSDK/Framework両方に対応する
- 新しいリソース/データソースの実装時には基本的にFrameworkを用いる
- 実験的に一部の小さいリソース/データソースをFrameworkへマイグレーションする

## やること/やらないこと

### やること

- terraform-plugin-muxの導入
- 一部の小さなリソース/データソースをFrameworkへマイグレーション

### やらないこと

- 既存リソース全てをFrameworkへマイグレーション

## 改訂履歴

- 2022/7/26: 初版作成
- 2022/11/11: 一部リソースのみFrameworkへマイグレーションする方針に変更
