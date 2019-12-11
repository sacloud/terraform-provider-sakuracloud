# CHANGELOG: v-nextでのスキーマ変更

## プロバイダー

  - `trace`属性をboolからstringへデータ型変更
  - 環境変数`SAKURACLOUD_TRACE_MODE`から`SAKURACLOUD_TRACE`へ名称変更
  - `fake_mode`属性の追加
  - `fake_store_path`属性の追加

## データソース

- データソース共通

  - `filter`/`name_selectors`/`tag_selectors`の廃止
  - `filters`の導入
  
```hcl
# filtersの例
data sakuracoud_server "example" {
  filters {
    names = ["Ubuntu","server","18"]
    # id = xxxxxxxxxx
    # tags = ["tag1"]
    # conditions {
    #   name = "Name"
    #   values = ["xxxxxxxxxxxxxxxxxx"]
    # }
  }
}
```

- Bridgeデータソース

  - `switch_ids`属性の廃止

- Iconデータソース

  - `body`属性の廃止
  
- Serverデータソース

  - `nw_mask_len`属性をstringからintへデータ型変更

- シンプル監視データソース

  - `health_check`.`status`属性をstringからintへデータ型変更
  - `health_check`.`delay_loop`をトップレベルへ移動
  
- VPCルータデータソース

  - フィールド名変更
    - `interface` -> `interfaces`
    - `dhcp_server` -> `dhcp_servers`
    - `dhcp_static_mapping` -> `dhcp_static_mappings`
    - `firewall` -> `firewalls`
    - `port_forwarding` -> `port_forwardings`
    - `static_route` -> `static_routes`
    - `user` -> `users`
    
  - サイト間VPNの詳細情報属性を除去
  
## リソース

- ブリッジ

  - `switch_ids`属性の廃止

- GSLB 実サーバ

  - `sakuracloud_gslb_server`を廃止
  
- ロードバランサ VIP/実サーバ

  - `sakuracloud_loadbalancer_vip`を廃止  
  - `sakuracloud_loadbalancer_server`を廃止  

- パケットフィルタ ルール

  - `sakuracloud_packet_filter_rule`を`sakuracloud_packet_filter_rules`に変更  
  これまでルールごとに1リソースだったものが複数のリソースを保持するようになった
    
- サーバ    

  - VNC関連項目を除去
  - 表示用IP(Interfaces.Switch.UserIPAddress)の設定を除去
  - NIC/追加NICを統合し`interfaces`を新設
    - `nic`に文字列を指定からオブジェクトを指定するように変更
    - *`nic`*がオブジェクトになることでデフォルト値が設定できなくなる。`interfaces`を明示的に書く必要がある。
    - パケットフィルタとMACアドレスを`interfaces`配下の各要素内に配置
    
#### 変更前

```hcl
resource sakuracloud_server "foobar" {
  name = "foobar"
  
  # nic  = "shared"      # 共有セグメントに接続(デフォルト)
  # nic = "disconnect"   # 切断状態
  # nic = "100000000001" # スイッチに接続(スイッチIDを指定)
  
  # 追加NIC
  additional_nics = ["100000000002", "100000000003"] #スイッチIDのリスト 
  
  # パケットフィルタ
  packet_filter_id = ["200000000001", "200000000002", "200000000003"] 
}
```

#### 変更後

```hcl
resource sakuracloud_server "foobar" {
  name = "foobar"
  interfaces {
    upstream         = "shared"
    packet_filter_id = "200000000001"
  }
  interfaces {
    upstream         = "100000000002" # スイッチID
    packet_filter_id = "200000000002" # パケットフィルタID
  }
  interfaces {
    upstream         = "100000000003" # スイッチID
    packet_filter_id = "200000000003" # パケットフィルタID
  }
}
```
    
- シンプル監視 

  - `health_check`.`status`属性をstringからintへデータ型変更
  - `health_check`.`delay_loop`をトップレベルへ移動
  
- VPCルータ

  - 子リソース(`sakuracloud_vpc_router_xxx`)を廃止
  - ファイアウォールルールのフィールド名変更
    - `source_nw` -> `source_network`
    - `dest_nw` -> `destination_network`
    - `dest_port` -> `destination_port`