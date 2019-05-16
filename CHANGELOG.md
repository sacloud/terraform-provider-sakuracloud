# Changelog

## 1.12.0 (2019-05-16)

* コア専有プラン #445 (yamamoto-febc)


## 1.11.4 (2019-04-19)

* APIリクエストでのレートリミット設定 #442 (yamamoto-febc)
* ポーリング時のリトライ処理誤り修正 #443 (yamamoto-febc)


## 1.11.3 (2019-04-18)

* スイッチ+ルータでの5000Mbpsプラン #440 (yamamoto-febc)


## 1.11.2 (2019-04-17)

* APIクライアントの分離 #435 (yamamoto-febc)


## 1.11.1 (2019-04-12)

* Terraform provider SDK: v0.12-early7 #429 (yamamoto-febc)
* パラメータ追加: APIリクエストタイムアウト #430 (yamamoto-febc)


## 1.11.0 (2019-04-08)

* VPCルータでネストした子要素の定義を可能にする #392 (yamamoto-febc)
* モバイルゲートウェイでネストした子要素の定義を可能にする #405 (yamamoto-febc)
* スタートアップスクリプトのdescriptionをcomputed:trueに変更 #406 (yamamoto-febc)
* Terraform v0.12.0-beta1 #408 (yamamoto-febc)
* エンハンスドロードバランサ追加 #410 (yamamoto-febc)
* エンハンスドロードバランサ: プラン選択 #414 (yamamoto-febc)
* marker_tag関連コードの除去 #415 (yamamoto-febc)
* ディスクリソースでのディスク修正関連パラメータをRemovedに変更 #416 (yamamoto-febc)
* Terraform v0.12対応 #422 (yamamoto-febc)
* Update docs - v2.0 to v1.11 #423 (yamamoto-febc)
* NFS新プラン対応 #424 (yamamoto-febc)
* エンハンスドロードバランサの正式版対応 #426 (yamamoto-febc)
* リリースファイルの名前からサフィックスを除去 #427 (yamamoto-febc)


## 1.10.0 (2019-02-21)

* Upgrade libsacloud to v1.14 #402 (yamamoto-febc)


## 1.9.1 (2018-12-19)

* Windows Server 2019対応 #399 (yamamoto-febc)


## 1.9.0 (2018-12-13)

* ゾーン情報データリソース追加 #395 (yamamoto-febc)
* SIMマルチキャリア対応 #396 (yamamoto-febc)
* Fix broken acceptance tests #397 (yamamoto-febc)


## 1.8.0 (2018-10-30)

* 表示用IPアドレス display_ipaddress 追加 #380 (yamamoto-febc)


## 1.7.1 (2018-10-29)

* terraform.ioスタイルの英語版ドキュメントのプレビュー #373 (yamamoto-febc)
* Fix ldflags #377 (blp1526)


## 1.7.0 (2018-10-17)

* サーバ/ディスク新プラン対応 #362 (yamamoto-febc)
* モバイルゲートウェイでのトラフィックコントロール機能 #363 (yamamoto-febc)
* libsacloud更新/データベースアプライアンス最新API対応 #364 (yamamoto-febc)
* vendor更新/terraform v0.11.8 #366 (yamamoto-febc)
* 石狩第1ゾーンの利用可能リソース拡大 #367 (yamamoto-febc)


## 1.6.1 (2018-09-13)

* DNSゾーン登録時のパラメータにデフォルト値を明示 #358 (yamamoto-febc)


## 1.6.0 (2018-09-11)

* アプライアンス作成失敗時のリトライ処理 #336 (yamamoto-febc)
* ディスクの修正パラメータをサーバリソース配下に移動 #351 (yamamoto-febc)


## 1.5.0 (2018-09-04)

* パブリックアーカイブ更新 #348 (yamamoto-febc)
* ディスクのストレージ隔離用パラメータ distant_from 追加 #349 (yamamoto-febc)


## 1.4.0 (2018-08-24)

* シンプル監視でのBasic認証 #341 (yamamoto-febc)
* ロードバランサーVIPヘの説明欄追加 #342 (yamamoto-febc)
* marker_tags設定を除去 #343 (yamamoto-febc)


## 1.3.4 (2018-08-09)

* サーバ作成時のプロビジョニング用接続情報のデフォルト値設定 #326 (yamamoto-febc)
* marker_tagsに廃止予定マークを付与 #338 (yamamoto-febc)
* ディスクの作成失敗時のリトライ #339 (yamamoto-febc)


## 1.3.3 (2018-07-24)

* fix typo: ipv4_prt -> ipv4_ptr #322 (yamamoto-febc)


## 1.3.2 (2018-07-19)

* typo storage type SSH to SSD #318 (223n)
* SQL Server 2017パブリックアーカイブ対応  #320 (yamamoto-febc)


## 1.3.1 (2018-07-05)

* VPCルータの設定適用タイミング修正 #313 (yamamoto-febc)


## 1.3.0 (2018-07-05)

* VPCルータでのインターネット接続 有効/無効 設定 #311 (yamamoto-febc)


## 1.2.0 (2018-07-03)

* SIMルート リソースの追加 #308 (yamamoto-febc)
* Update README.md #309 (yamamoto-febc)


## 1.1.3 (2018-06-28)

* シンプル監視でのSNI対応 #306 (yamamoto-febc)


## 1.1.2 (2018-06-05)

* データベースアプライアンスのデフォルトバージョン更新 #304 (yamamoto-febc)


## 1.1.1 (2018-04-13)

* 逆引きレコード登録機能 #300 (yamamoto-febc)


## 1.1.0 (2018-04-13)

* Pickup from v1.1-dev #274 (yamamoto-febc)
* Update sacloud/libsacloud #277 (yamamoto-febc)
* go v1.10対応 #278 (yamamoto-febc)
* Accept-Languageパラメータ追加 #279 (yamamoto-febc)
* SetID呼び出しタイミングの統一 #282 (yamamoto-febc)
* データベースアプライアンス 500GB/1TB プラン追加 #283 (yamamoto-febc)
* TypeSetからTypeListへの変更 #284 (yamamoto-febc)
* リリース時のドキュメント生成を自動化 #286 (yamamoto-febc)
* アーカイブID変更時のディスク再生成抑制についてのドキュメント #287 (yamamoto-febc)
* インストールガイド更新 #288 (yamamoto-febc)
* アップグレードガイドの追加 #289 (yamamoto-febc)
* セキュアモバイル対応 #291 (yamamoto-febc)
* Homebrewスクリプトの更新 #292 (yamamoto-febc)
* Homebrewでのインストールガイド更新 #293 (yamamoto-febc)
* モバイルゲートウェイでのスタティックルート設定機能 #294 (yamamoto-febc)
* サーバの1番目のNICを接続なしで作成 #295 (yamamoto-febc)
* モバイルゲートウェイの起動時の待機処理 #297 (yamamoto-febc)


## 1.0.5 (2018-01-19)

* サーバリソースの変更検知でCustomizeDiffを利用 #255 (yamamoto-febc)
* データベースアプライアンス作成APIパラメータの修正 #265 (yamamoto-febc)
* Terraform v0.11.2対応 #266 (yamamoto-febc)


## 1.0.4 (2017-12-19)

* Revert "サーバリソースの変更検知にCustomizeDiff導入" #253 (yamamoto-febc)


## 1.0.3 (2017-12-15)

* APIクライアントの503エラー時リトライでの指定秒数待機 #250 (yamamoto-febc)


## 1.0.2 (2017-12-14)

* API呼び出しで503発生時のリトライ #248 (yamamoto-febc)


## 1.0.1 (2017-12-07)

* リリーススクリプト調整 #239 (yamamoto-febc)
* TravisCI上でのCI/CD改善 #241 (yamamoto-febc)
* サーバリソースの変更検知にCustomizeDiff導入 #243 (yamamoto-febc)


## 1.0.0 (2017-12-04)

* パブリックアーカイブPlesk除去 #192 (yamamoto-febc)
* ISOイメージリソースの追加 #198 (yamamoto-febc)
* ソースコードレイアウトの修正 #202 (yamamoto-febc)
* アーカイブリソースの追加 #205 (yamamoto-febc)
* マーカータグ付与オプション #206 (yamamoto-febc)
* リソースRead時の404レスポンスハンドリング改善 #207 (yamamoto-febc)
* サーバリソースでの専有ホストID指定 #208 (yamamoto-febc)
* データリソースへのセレクタ属性追加 #210 (yamamoto-febc)
* バージョン情報JSON追加 #211 (yamamoto-febc)
* テストケース追加 #212 (yamamoto-febc)
* VPCルータへの機能追加(DHCPでのDNSサーバ配布/NICごとのファイアウォール) #214 (yamamoto-febc)
* 配布サイト用に静的コンテンツのビルドを実行 #216 (yamamoto-febc)
* 専有ホスト対応 #218 (yamamoto-febc)
* リリースファイルにバージョン情報を付与 #219 (yamamoto-febc)
* 専有サーバでのサーバ登録解除処理 #221 (yamamoto-febc)
* Terraform v0.11対応 #222 (yamamoto-febc)
* リンクしているTerraformライブラリのバージョン情報表示ツール #223 (yamamoto-febc)
* TravisCIからの通知設定追加 #224 (yamamoto-febc)
* テスト用のゾーン検証スキップオプション追加 #225 (yamamoto-febc)
* CI/CD時にBotアカウントを利用 #226 (yamamoto-febc)
* テスト用APIルートURL変更オプションの追加 #227 (yamamoto-febc)
* データリソース関連の改善 #228 (yamamoto-febc)
* コピー待ち処理の改善 #229 (yamamoto-febc)
* AppVeyorでのCI #230 (yamamoto-febc)
* VPCルータ データリソースの追加 #231 (yamamoto-febc)
* AUTHORS出力処理の追加 #232 (yamamoto-febc)
* コンテンツ自体の最大横幅を1200pxに変更 #233 (wate)
* Terraform v0.11.1 #235 (yamamoto-febc)
* サブリソースなどでの404レスポンスハンドリング改善 #236 (yamamoto-febc)
* タグを持つリソースでCustomizeDiffを登録 #237 (yamamoto-febc)


## 0.14.0 (2017-10-14)

* VPCルータでのサイト間VPN接続情報追加 #184 (yamamoto-febc)
* リソース追加:パケットフィルタルール #185 (yamamoto-febc)
* 2段階シャットダウン実装 #186 (yamamoto-febc)
* シンプル監視でのSSLサーバ証明書 有効期限監視 #188 (yamamoto-febc)
* シンプル監視方法での内部ID保持方法の修正 #190 (yamamoto-febc)

## 0.13.1 (2017-10-01)

* スタートアップスクリプトへのclass属性追加 #178 (yamamoto-febc)
* リリースプロセス変更 #180 (yamamoto-febc)
* NFSプラン拡大対応 #181 (yamamoto-febc)


## 0.13.0 (2017-09-07)

* NFSアプライアンス #176 (yamamoto-febc)


## 0.12.0 (2017-09-04)

* NICドライバ(interface_driver)の追加 #173 (yamamoto-febc)
* 属性名の変更 #174 (yamamoto-febc)


## 0.11.0 (2017-07-21)

* ドキュメントのパケットフィルタ設定例を修正 #159 (mapk0y)
* オブジェクトストレージ対応 #161 (yamamoto-febc)
* パケットフィルタ更新時のリソース再生成判定を修正 #162 (yamamoto-febc)
* サンプル.tfファイルのクリーンアップ #163 (yamamoto-febc)
* go get対応 #164 (yamamoto-febc)
* tfファイルのフォーマット #165 (yamamoto-febc)
* アイコン(sakuracloud_icon)追加 #167 (yamamoto-febc)
* サーバ コネクタリソース追加 #168 (yamamoto-febc)
* コード改善 #169 (yamamoto-febc)


## 0.10.4 (2017-07-10)

* textlintルール追加 #153 (yamamoto-febc)
* Linuxbrew対応 #155 (yamamoto-febc)
* クラウドAPI呼び出しへのスロットリング追加 #157 (yamamoto-febc)


## 0.10.3 (2017-06-28)

* textlint導入 #150 (yamamoto-febc)
* Linux i386版での数値型パラメータ欠落への対応 #151 (yamamoto-febc)


## 0.10.2 (2017-06-21)

* ドキュメントのタイポ修正 #142 (s-shinoda)
* ドキュメント内のパラメータの型を修正 #143 (s-shinoda)
* シンプル監視:portのtfstateへの反映処理修正 #145 (yamamoto-febc)
* Windows2016 SQLServer Standard(RDS+Office)パブリックアーカイブ追加 #146 (yamamoto-febc)
* SiteGuardパブリックアーカイブ除去 #147 (yamamoto-febc)


## 0.10.1 (2017-05-29)

* Update database docs #136 (yamamoto-febc)
* TravisCI上でのgolintインストール #138 (yamamoto-febc)
* ディスクリソースでの属性名変更 #140 (yamamoto-febc)


## 0.10.0 (2017-05-24)

* golint導入 #128 (yamamoto-febc)
* データベースアプライアンス更新 #129 (yamamoto-febc)
* ロードバランサでの属性名変更 #130 (yamamoto-febc)
* ドキュメントでの音引き表記揺れ #131 (yamamoto-febc)
* sacloud配下へのリポジトリ移転 #132 (yamamoto-febc)


## 0.9.1 (2017-05-12)

* Pleskパブリックアーカイブ追加 #125 (yamamoto-febc)


## 0.9.0 (2017-05-09)

* データベースアプライアンス 石狩第2ゾーン対応 #113 (yamamoto-febc)
* RancherOS パブリックアーカイブ追加 #114 (yamamoto-febc)
* VPCRouterでの文字列長バリデーション追加 #115 (yamamoto-febc)
* スイッチ+ルータでの追加IPアドレス(sakuracloud_subnet)追加 #117 (yamamoto-febc)
* スイッチ+ルータでのIPv6対応 #118 (yamamoto-febc)
* ネットワーク関連の属性名をシンプル化 #119 (yamamoto-febc)
* パスワード項目に対してsensitiveフラグ有効化 #120 (yamamoto-febc)


## 0.8.1 (2017-04-15)

- サーバーNIC接続変更時の挙動修正 #104 (yamamoto-febc)
- Windows2008パブリックアーカイブ除去 #106 (yamamoto-febc)
- 公開鍵生成機能 #107 (yamamoto-febc)


## 0.8.0 (2017-04-06)

* Don't set network parameters #98 (yamamoto-febc)
* Remove backup_hour param #99 (yamamoto-febc)
* Add database plans #100 (yamamoto-febc)


## 0.7.2 (2017-03-29)

* Using CI pipelines is started (yamamoto-febc)
