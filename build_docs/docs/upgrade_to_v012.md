# v0.12での変更点

---

v0.12ではリソースの属性名に対し後方互換性のない変更が行われています。  
これまでのバージョンをご利用いただいていた場合、tfファイルの修正などの対応が必要になる場合があります。  
以下の属性名の変更内容をご確認いただき必要に応じてtfファイルの修正を行ってください。

## 属性名の変更内容

以下の変更が行われています。

- これまで非推奨となっていた属性
- これまで名称に大文字が使われていた属性  
(terraform v0.10以降、属性名に大文字を利用することは非推奨となりました)

具体的な対象の属性は以下の通りです。  
tfファイルにて属性名(旧)が利用されていた場合、属性名(新)へ置き換える必要があります。 

|対象リソース                               |属性名(旧)                | 属性名(新)           | 説明           | 
|-----------------------------------------|-------------------------|---------------------|---------------|
|サーバ(`sakuracloud_server`)              | `base_interface`        | `nic`               | 基本NIC |
|                                         | `additional_interfaces` | `additional_nics`   | 追加NIC |
|                                         | `base_nw_ipaddress`     | `ipaddress`         | IPアドレス |
|                                         | `base_nw_dns_servers`   | `dns_servers`       | DNSサーバ |
|                                         | `base_nw_gateway`       | `gateway`           | ゲートウェイ |
|                                         | `base_nw_address`       | `nw_address`        | ネットワークアドレス |
|                                         | `base_nw_mask_len`      | `nw_mask_len`       | ネットワークマスク長 |
|スイッチ+ルータ(`sakuracloud_internet`)    | `nw_gateway`            | `gateway`           | ゲートウェイ |
|                                         | `nw_min_ipaddress`      | `min_ipaddress`     | 割り当て可能な最小IPアドレス |
|                                         | `nw_max_ipaddress`      | `max_ipaddress`     | 割り当て可能な最大IPアドレス |
|                                         | `nw_ipaddresses`        | `ipaddresses`       | IPアドレス |
|ロードバランサ(`sakuracloud_load_balancer`)| `is_double`             | `high_availability` | 冗長化構成の有無 |
|                                         | `VRID`                  | `vrid`              | VRID |
|VPCルータ(`sakuracloud_vpc_router`)       | `VIRD`                  | `vrid`              | VRID |
|GSLB(`sakuracloud_gslb`)                 | `FQDN`                  | `fqdn`              | FQDN |

