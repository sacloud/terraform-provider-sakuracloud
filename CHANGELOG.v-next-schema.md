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

- PacketFilterルール

  - `sakuracloud_packet_filter_rule`を`sakuracloud_packet_filter_rules`に変更  
  これまでルールごとに1リソースだったものが複数のリソースを保持するようになった
    
- シンプル監視 

  - `health_check`.`status`属性をstringからintへデータ型変更
  - `health_check`.`delay_loop`をトップレベルへ移動
  
- VPCルータ

  - 子リソース(`sakuracloud_vpc_router_xxx`)を廃止
  - ファイアウォールルールのフィールド名変更
    - `source_nw` -> `source_network`
    - `dest_nw` -> `destination_network`
    - `dest_port` -> `destination_port`