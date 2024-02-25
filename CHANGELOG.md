# CHANGELOG

## 2.25.3 (2024/2/26)

- パッケージのダウングレード #1146

## 2.25.2 (2024/2/14)

- Fix: sakuracloud_proxylb: backend_http_keep_aliveとsyslogのデフォルト値/空の場合の挙動を修正 #1144
- private_key: sensitive=true #1142
- sakuracloud_proxylb: プラン変更で再作成される問題を修正 #1141
- sakuracloud_proxylb_acme: BackendHttpKeepAliveの設定が引き継がれない問題を修正 #1135

## 2.25.1 (2023/12/18)

- Update CHANGELOG #1122, 1123
- github.com/sacloud/autoscaler v0.15.5 #1121
- go: bump github.com/sacloud/iaas-service-go from 1.9.2-0.20230808054001-efad52d748d4 to 1.9.2 #1119
- go: bump github.com/sacloud/webaccel-api-go from 1.1.5 to 1.1.6 #1118
- ci: bump actions/setup-go from 4 to 5 #1115
- ドキュメント更新：ELBリクエストヘッダによる振り分けルール #1114 
- Improve CI #1113
- go: bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.27.0 to 2.30.0 #1112

## 2.25.0 (2023/11/20)

- ELBリクエストヘッダによる振り分けルール対応  #1110 
- go: bump github.com/hashicorp/terraform-plugin-go from 0.16.0 to 0.19.0 #1105 by dependabot bot was closed 4 days ago
- ロードバランサの受け入れテストをデフォルトでスキップ #1099 
- go 1.21 #1094 
- Fix: サーバプラン変更テスト時のパラメータ指定誤りを修正 #1093
- AMDプランのテストを環境変数指定時のみに限定 #1092
- AMDプラン #1090
- go: bump golang.org/x/text from 0.10.0 to 0.12.0 #1089
- オートスケールのdays_of_weekをrequiredに変更 #1085

## 2.24.1 (2023/06/09)

- iaas-api-go v1.11.1 - モバイルゲートウェイの電源APIの修正 #1073

## 2.24.0 (2023/05/26)

- エンハンスドDBでのタイプ/リージョン/接続元ネットワーク指定 #1070
- sacloud/autoscaler@v0.15.3 #1071

## 2.23.2 (2023/04/26)

- validate disk_edit_parameter.hostname #1064

## 2.23.1 (2023/04/20)

- go: bump github.com/goccy/go-yaml from 1.10.1 to 1.11.0 #1056
- update dependencies #1060
- sacloud/iaas-service-go@v1.8.1 #1061
- sacloud/iaas-service-go@v1.8.2 #1062

## 2.23.0 (2023-03-30)

- sacloud/autoscaler v0.15.1 #1053
- fix: http_backend_keep_aliveのスキーマ修正 - Computed:true #1052
- go: bump github.com/goccy/go-yaml from 1.10.0 to 1.10.1 #1051
- fix: http_backend_keep_aliveのスキーマ修正 - デフォルト値の除去 #1050
- iaas-api-go v1.9.1 #1049
- go: bump github.com/hashicorp/terraform-plugin-go from 0.14.3 to 0.15.0 #1048 
- go: bump github.com/sacloud/webaccel-api-go from 1.1.4 to 1.1.5 #1047
- iaas-api-go v1.9 #1046
- go: bump github.com/hashicorp/terraform-plugin-sdk/v2 #1041
- ci: bump actions/setup-go from 3 to 4 #1039

## 2.22.2 (2023-03-17)

FEATURES:

- sacloud/autoscaler v0.14.1 #1038

MISC:

- マネージドPKIのテストをON/OFF可能に #1036
- go 1.20 #1037
- go: bump github.com/stretchr/testify from 1.8.1 to 1.8.2 #1034
- go: bump github.com/sacloud/iaas-service-go from 1.6.0 to 1.6.1 #1033
- go: bump github.com/goccy/go-yaml from 1.9.8 to 1.10.0 #1035

## 2.22.1 (2023-02-08)

FIXES:

- fix: sakuracloud_simple_monitor: verify_sniの反映誤りを修正 #1019

MISC:

- go: bump github.com/sacloud/iaas-api-go from 1.8.0 to 1.8.1
- go: bump github.com/sacloud/iaas-api-go from 1.8.1 to 1.8.3


## 2.22.0 (2023-01-19)

FEATURES:

- ELB: source_ips #1011

MISC:

- e2e: is1a -> is1b #1006 , #1007
- copyright: 2023 #1009

IMPROVEMENTS:

- go: bump github.com/goccy/go-yaml from 1.9.7 to 1.9.8 #1004
- go: bump github.com/hashicorp/terraform-plugin-mux from 0.7.0 to 0.8.0 #1005
- go: bump golang.org/x/text from 0.5.0 to 0.6.0 #1008
- go: bump github.com/hashicorp/terraform-plugin-go from 0.14.2 to 0.14.3  #1010

## 2.21.0 (2022-12-20)

FEATURES

- iaas-api-go v1.7.0 #1000
  - オートスケールのトラフィック量トリガー
  - VPCルータでのDHグループ指定

IMPROVEMENTS:

- iaas-api-go v1.7.1 #1003
- go: bump github.com/goccy/go-yaml from 1.9.6 to 1.9.7 #991
- go: bump golang.org/x/text from 0.4.0 to 0.5.0 #993
- go: bump github.com/sacloud/packages-go from 0.0.6 to 0.0.7 #996
- ci: bump goreleaser/goreleaser-action from 3 to 4 #997
 

## 2.20.1 (2022-12-05)

FEATURES

- iaas-api-go v1.6.2(MIRACLE LINUX9対応) #991

FIXES

- iaas-api-go v1.6.1 #989

IMPROVEMENTS:

- SDK/Frameworkで共通利用するfunc類をinternal配下に切り出し #981
- go: bump github.com/goccy/go-yaml from 1.9.5 to 1.9.6 #978
- go: bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.24.0 to 2.24.1 #982
- go: bump github.com/sacloud/autoscaler from 0.11.1 to 0.11.2 #984
- go: bump github.com/hashicorp/terraform-plugin-go from 0.14.1 to 0.14.2 #987


DOCS:

- docs: Terraform Plugin Frameworkの扱いの方針変更 #980

## 2.20.0 (2020-10-26)

FEATURES:

- sakuracloud_database: RDBMSのバージョン指定機能 #975

IMPROVEMENTS:

- go: bump github.com/sacloud/iaas-service-go from 1.3.2 to 1.4.0 #977

## 2.19.4 (2022-10-19)

IMPROVEMENTS:

- SAKURACLOUD_ZONES環境変数をサポート #974
- go: bump github.com/sacloud/iaas-api-go from 1.4.1 to 1.5.1 #970
- go: bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.23.0 to 2.24.0 #969
- go: bump golang.org/x/text from 0.3.8 to 0.4.0 #973

## 2.19.3 (2022-10-12)

FIXES:

- fix: プロファイルからの値がうまく読み込めないことがある問題を修正 #967

## 2.19.2 (2022-10-11)

FEATURES:

  - sacloud/autoscaler v0.11.1 #964
  - sakuracloud_auto_scale: GA #965


## 2.19.1 (2022-10-06)

IMPROVEMENTS:

  - go 1.19 #957
  - sacloud/makefile v0.0.7 #958

## 2.19.0 (2022-09-30)

Note: 
- sakuracloud_archiveのos_typeで指定できる値が変更されています。指定できる値の詳細はsakuracloud_archiveのドキュメントを参照ください。
- このバージョンからCHANGELOGの形式を変更します

FEATURES:
- iaas-api-go v1.4.0 - OSType更新 #955

IMPROVEMENTS:

- go: bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.21.0 to 2.23.0 #951

## 2.18.0 (2022-08-22)

このバージョンからterraform-plugin-muxによるTerraform Plugin SDK/Framework両方への対応が行われています(GH-932)。  
現時点ではFrameworkを利用するリソースはありませんが、今後追加されるリソースはFrameworkを用いて実装されます。

FEATURES:

  - sakuracloud_vpc_router: スケジュールメンテナンス機能 [GH-948] (@yamamoto-febc)
  - sakuracloud_vpc_router: サイト間VPNパラメータ [GH-943] (@yamamoto-febc)
  - 各データソースのfilter.conditionで完全一致/部分一致のどちらで検索するか指定可能にする [GH-942] (@yamamoto-febc)

IMPROVEMENTS:

  - オートスケーラーのテスト改善 [GH-930] (@yamamoto-febc)
  - オートスケーラーのテスト改善 [GH-927] (@yamamoto-febc)
  - エンハンスドロードバランサでの証明書取得待ち処理の改善 [GH-925] (@yamamoto-febc)
  - Let's Encryptのレートリミットエラー回避対応 [GH-928] (@yamamoto-febc)
  - VPCルータ作成直後に値が参照できない問題を修正 [GH-939] (@yamamoto-febc)

DEVELOPMENTS:

  - iaas-service-go@v1.3.1 [GH-946] (@yamamoto-febc)
  - go: bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.20.0 to 2.21.0 [GH-940] (@dependabot)
  - go: bump github.com/sacloud/webaccel-api-go from 1.1.2 to 1.1.3 [GH-936] (@dependabot)
  - crazy-max/ghaction-import-gpg@v5でのパラメータ名変更対応 [GH-935] (@yamamoto-febc)
  - terraform-plugin-mux対応 [GH-932] (@yamamoto-febc)
  - go: bump github.com/sacloud/packages-go from 0.0.3 to 0.0.4 [GH-922] (@dependabot)
  - go: bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.18.0 to 2.19.0 [GH-923] (@dependabot)
  - iaas-api-go@v1.2 [GH-917] (@yamamoto-febc)
  - sacloud/go-template@v0.0.5 [GH-913 , GH-916] (@yamamoto-febc)
  - sacloud/go-template@v0.0.2 [GH-902] (@yamamoto-febc)
  - go: bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.17.0 to 2.18.0 [GH-910] (@dependabot)
  - go: bump github.com/stretchr/testify from 1.7.1 to 1.7.5 [GH-903] (@dependabot)
  - go: bump github.com/sacloud/iaas-service-go from 1.1.2 to 1.1.3 [GH-905] (@dependabot)
  - go: bump github.com/sacloud/api-client-go from 0.1.0 to 0.2.0 [GH-906] (@dependabot)
  - go: bump github.com/goccy/go-yaml from 1.8.9 to 1.9.5 [GH-907] (@dependabot)
  - go: bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.10.1 to 2.17.0 [GH-904] (@dependabot)

DOCS:
  - docs: Terraform Plugin Frameworkへの対応方針 [GH-931] (@yamamoto-febc)

## 2.17.1 (2022-06-06)

FIXES: 
  - sakuracloud_proxylb_acmeでruleが反映されない問題を修正 [GH-901] (@yamamoto-febc)

DOCS:
   - yamlencodeからjsonencodeへ変更 [GH-900] (@yamamoto-febc)

## 2.17.0 (2022-06-03) 

FEATURES:

  - AutoScale [GH-895] (@yamamoto-febc)
  - iaas-api-go/v1.1.2 [GH-896] (@yamamoto-febc)
  - sakuracloud_vpc_router: supports netmask /29 [GH-897] (@yamamoto-febc)

IMPROVEMENTS:

  - iaas-api-go v1.1.1 [GH-894] (@yamamoto-febc)
  - github.com/hashicorp/hc-install v0.3.2 [GH-892] (@yamamoto-febc)
  - switch to iaas-service-go [GH-891] (@yamamoto-febc)

MISC:

  - docs: misc updates [GH-898] (@yamamoto-febc)

## 2.16.2 (2022-02-18)

FEATURES:

  - sakuracloud_dns: upgrade MaxItems to 2000 [GH-889] (@yamamoto-febc)

## 2.16.1 (2022-01-27)

MISC:
  - Upgrade dependencies - github.com/hashicorp/terraform-plugin-sdk v2.10.1 [GH-887] (@yamamoto-febc)
  - libsacloud v2.32.1 - Removed some os-types: freebsd and coreos [GH-886] (@yamamoto-febc)

## 2.16.0 (2021-12-27)

FEATURES:

- vpc_router: dns_forwarding [GH-881] (@yamamoto-febc)
- simple_monitor: retry [GH-880] (@yamamoto-febc)

IMPROVEMENTS:

- fix broken test - TestAccSakuraCloudDataSourceCDROM_basic [GH-883] (@yamamoto-febc)
- Remove centos8 from examples [GH-882] (@yamamoto-febc)
- libsacloud v2.31 [GH-879] (@yamamoto-febc)
- libsacloud v2.30.0 - PostgreSQL 13 [GH-878] (@yamamoto-febc)

Note: `data.sakuracloud_archive#os_type`:  `centos8` has been removed. 

## 2.15.0 (2021-12-08)

FEATURES:

- simple_monitor: verify_sni [GH-877] (@yamamoto-febc)

IMPROVEMENTS:

- Upgrade dependencies - libsacloud [GH-876] (@yamamoto-febc)
- Update dependencies - terraform-plugin-sdk/v2 v2.9.0 [GH-875] (@yamamoto-febc)
- Update dependencies - terraform-plugin-sdk/v2 v2.8.0 [GH-873] (@yamamoto-febc) 
- Go 1.17 [GH-872] (@yamamoto-febc)
- Update dependencies- github.com/sacloud/libsacloud/v2 v2.28.0 [GH-871] (@yamamoto-febc)

## 2.14.2 (2021-10-25)

FIXES:

  * libsacloud v2.27.1 [GH-870] (@yamamoto-febc)

MISC:

  * Skip local router tests in acceptance test [GH-867] (@yamamoto-febc)

## 2.14.1 (2021-10-14)

FEATURES

  - libsacloud v2.27 - miracle linux [GH-866] (@yamamoto-febc)

## 2.14.0 (2021-10-08)

FEATURES

  - GPU plan [GH-865] (@yamamoto-febc)

## 2.13.0 (2021-10-05)

FEATURES

  - Managed PKI [GH-862] (@yamamoto-febc)
  - ELB: Proxy Protocol v2 [GH-857] (@yamamoto-febc)
  - libsacloud v2.25.1 - debian11 [GH-860] (@yamamoto-febc)
  - simple_monitor: ftp/ftps [GH-861] (@yamamoto-febc)

MISC:
  - libsacloud v2.26.0 [GH-864] (@yamamoto-febc)
  - Update docs: user_data [GH-858] (@yamamoto-febc)

## 2.12.0 (2021-08-19)

FEATURES:
  - cloud-init [GH-856] (@yamamoto-febc)

## 2.11.0 (2021-07-30)

FEATURES:
  - Support @previous-id tags for server/internet/proxylb [GH-855] (@yamamoto-febc)

## 2.10.2 (2021-07-26)

FIXES

  - libsacloud v2.21.1 [GH-854] (@yamamoto-febc)

## 2.10.1 (2021-07-19)

FIXES

- libsacloud v2.20.1 [GH-852] (@yamamoto-febc)

## 2.10.0 (2021-07-09)

FEATURES

- Enhanced Database [GH-847] by @yamamoto-febc
- SimpleMonitor: timeout [GH-850] by @yamamoto-febc
- ELB: supports syslog block and ssl_policy field [GH-849] by @yamamoto-febc

MISC

- libsacloud v2.20.0 [GH-851] by @yamamoto-febc
- Update dev env [GH-848] by @yamamoto-febc

## 2.9.3 (2021-06-29)

IMPROVEMENTS

- Terraform Plugin SDK v2.7.0 [GH-840] (@yamamoto-febc)
- Rename flag from -debuggable to -debug [GH-844] (@yamamoto-febc)

FIXES

- Moved default size setting from the schema definition section [GH-843] (@yamamoto-febc)


## 2.9.2 (2021-06-24)

- Fix zone name attribute of DNS [GH-838] (@chibiegg) 

## 2.9.1 (2021-06-21)

- fixed gzip misconfiguration [GH-837] (@yamamoto-febc)

## 2.9.0 (2021-06-15)

- LoadBalancer: Increased VIP limit to 20 [GH-835] (@yamamoto-febc)
- VPCRouter: WireGuard server [GH-834] (@yamamoto-febc)
- ELB: extending rule-based-balancing [GH-831] (@yamamoto-febc)
- libsacloud v2.19.1 [GH-834] (@yamamoto-febc)
- libsacloud v2.19.0 [GH-829] (@yamamoto-febc)

## 2.8.4(2021-05-02)

* Remove bucket_object resources [GH-818] (@yamamoto-febc)
* sakuracloud_proxylb_acme: subject_alt_names [GH-819] (@yamamoto-febc)
* Fixes CI problems [GH-821 , GH-823] (@yamamoto-febc)
  - github.com/hashicorp/terraform-plugin-sdk v2.6.1
* libsacloud v2.18.1 [GH-825] (@yamamoto-febc)

## 2.8.3(2021-04-12)

FEATURES:

  * simple_monitor: contains_string & proxylb: gzip [GH-814] (@yamamoto-febc)
  * libsacloud v2.17 [GH-814] (@yamamoto-febc)

## 2.8.2(2021-04-01)

FEATURES:

  * simple_monitor: http2 [GH-813] (@yamamoto-febc)

## 2.8.1(2021-03-26)

FIXES:

  * Fixed parameter handling for server plan change operation [GH-812] (@yamamoto-febc)

DOCS:

  * Fix packet_filter examples [GH-811] (@tokibi) 

## 2.8.0(2021-03-22)

ENHANCEMENTS:

  * **Terraform Plugin SDK v2** [GH-807] (@yamamoto-febc)

FEATURES:

  * sakuracloud_server: ssh_keys [GH-805] (@yamamoto-febc)
  * sakuracloud_sim: Sentitive:true [GH-802] (@yamamoto-febc)

IMPROVEMENTS:

  * Added some rules for tfproviderlint [GH-809] (@yamamoto-febc)
  * upgrade dependencies - libsacloud to v2.14.1 - added internet plans [GH-806],[GH-799] (@yamamoto-febc)
  * darwin/arm64 [GH-800] (@yamamoto-febc)

## 2.7.1 (2021-02-17)

* libsacloud v2.13 [GH-796] (@yamamoto-febc)

## 2.7.0 (2021-01-16)

* VPCRouter version [GH-795] (@yamamoto-febc)
* Added support for parameters blocks in sakuracloud_database [GH-794] (@yamamoto-febc)
* libsacloud v2.11.0 [GH-793] (@yamamoto-febc)

## 2.6.0 (2021-01-05)

* Update copyright [GH-792] (@yamamoto-febc)
* WebAccelerator [GH-791] (@yamamoto-febc)
* libsacloud v2.9 [GH-790] (@yamamoto-febc)

## 2.5.4 (2020-11-19)

* libsacloud v2.8.10 [GH-786] (@yamamoto-febc)

## 2.5.3 (2020-10-27)

* libsacloud v2.8.6 [GH-783] (@yamamoto-febc)

## 2.5.2 (2020-10-23)

*　libsacloud v2.8.5 [GH-781] (@yamamoto-febc)

## 2.5.1 (2020-10-21)

* libsacloud v2.8.4 [GH-779] (@yamamoto-febc)
* Use d.Id() for building error message [GH-774] (@yamamoto-febc)

## 2.5.0 (2020-09-30)

FEATURES

* ESME [GH-772] (@yamamoto-febc)
* libsacloud v2.8.1 [GH-771] (@yamamoto-febc)

IMPROVEMENTS

* Use d.Id() for building error message [GH-774] (@yamamoto-febc)

MISC

* Remove Dockerfile [GH-773] (@yamamoto-febc)

## 2.4.1 (2020-09-17)

FEATURES

* Add default_zone to provider config [GH-769] (@higebu)

## 2.4.0 (2020-08-20)

FEATURES

* tk1b zone (libsacloud v2.7.0) [GH-767] (@yamamoto-febc)

IMPROVEMENTS

* Skip acc tests for the object storage when env is not set [GH-762] (@yamamoto-febc)

## 2.3.6 (2020-08-11)

FIXES

* Fix error handling at VPCRouter [GH-757] (@yamamoto-febc)
* Use libsacloud v2.6.4 [GH-754] (@yamamoto-febc)
* Fix time comparison method - use time#Equal() [GH-753] (@yamamoto-febc)
* Use libsacloud v2.6.3  [GH-751] (@yamamoto-febc)

DOCS

* docs: Terraform v0.13 [GH-758] (@yamamoto-febc)

## 2.3.5 (2020-06-19)

* Update dependencies - libsacloud v2.6.1 [GH-748] (@yamamoto-febc)

## 2.3.4 (2020-06-16)

* ProxyLB: supports anycast [GH-747] (@yamamoto-febc)

## 2.3.3 (2020-05-18)

Note: Publishing in the Terraform Registry is supported this version and later.
see http://registry.terraform.io/providers/sacloud/sakuracloud/

FEATURES

* Support for publishing in the Terraform Registry [GH-744] (@yamamoto-febc)

FIXES

* Fixes plan changing of ProxyLB [GH-745] (@yamamoto-febc)

## 2.3.2 (2020-05-15)

* This is an experimental release for testing publishing to the Terraform registry. Don't use this in a production environment.

## 2.3.1 (2020-04-24)

FIXES

* Modify how to determine whether to pass disk_edit_parameter to ServerBuilder [GH-737] (@yamamoto-febc)

MISC

* Fix broken CI - install golangci-lint via install script [GH-735] (@yamamoto-febc)

## 2.3.0 (2020-04-20)

* Startup Script Parameters [GH-731] (@yamamoto-febc)
    * libsacloud v2.5.1
* libsacloud v2.5.2 - improve error messages [GH-733] (@yamamoto-febc)    

## 2.2.0 (2020-03-17)

FEATURES

* Add sakuracloud_archive_share resource [GH-728] (@yamamoto-febc)
* Supports transferred/shared archives [GH-727] (@yamamoto-febc)
    * libsacloud v2.4.1

IMPROVEMENTS

* Set ID to state even if got error from builders [GH-726] (@yamamoto-febc)
* libsacloud v2.3.0 - MariaDB 10.4 [GH-724]

## 2.1.2 (2020-03-10)

* Remove deletion waiter [GH-713] (@yamamoto-febc)
* libsacloud v2.1.7 [GH-713] (@yamamoto-febc)
* Go 1.14 [GH-712] (@yamamoto-febc)
* Fix wrong error message [GH-718] (@yamamoto-febc)
* libsacloud v2.1.8 to avoid marshal JSON error at SIM [GH-714] (@yamamoto-febc)
* libsacloud v2.1.9 [GH-723] (@yamamoto-febc)

## 2.1.1 (2020-02-28)

IMPROVEMENS/FIXES

* tfproviderlint v0.10.0 [GH-708] (@yamamoto-febc)
* libsacloud v2.1.4 [GH-708] (@yamamoto-febc)
* Upgrade libsacloud to v2.1.5 [GH-709] (@yamamoto-febc)

## 2.1.0 (2020-02-14)

FEATURES

* Container Registry: VirtualDomain/User permission [GH-704] (@yamamoto-febc)
* PostgreSQL 12.1 [GH-704] (@yamamoto-febc)

IMPROVEMENTS

* Terraform Plugin SDK v1.7.0 [GH-703] (@yamamoto-febc)
* tfproviderlint v0.9.0 [GH-698] (@yamamoto-febc)

## 2.0.1 (2020-02-06)

FIXES

* libsacloud v2.0.2 [GH-697] @yamamoto-febc

IMPROVEMENTS

* terraform-plugin-sdk v1.6.0 and tfproviderlint v0.9.0 [GH-698] @yamamoto-febc
* Use libsacloud v2.0.1 [GH-696] @yamamoto-febc

## 2.0.0 (2020-01-31)

NOTES:

* Initial release of v2.0

