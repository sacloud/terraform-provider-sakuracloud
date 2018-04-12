# モバイルゲートウェイ(sakuracloud_mobile_gateway)

---

### 設定例

```hcl
# モバイルゲートウェイの定義
resource "sakuracloud_mobile_gateway" "mgw" {
    name = "example"
    
	internet_connection = true
	
    switch_id           = "${sakuracloud_switch.sw.id}" # スイッチのID
    private_ipaddress   = "192.168.11.101"              # プライベート側IPアドレス
    private_nw_mask_len = 24                            # プライベート側ネットワークマスク長
   
    #dns_server1 = "8.8.8.8" # DNSサーバ1
    #dns_server2 = "8.8.4.4" # DNSサーバ2
    
    description = "example"
    tags = ["tags1" , "tags2"]
}

# スタティックルートの定義
resource "sakuracloud_mobile_gateway_static_route" "route01" {
    mobile_gateway_id = "${sakuracloud_mobile_gateway.mgw.id}"
    
    prefix   = "192.168.10.0/24"
    next_hop = "192.168.11.201"
}

# スタティックルートの定義
resource "sakuracloud_mobile_gateway_static_route" "route02" {
    mobile_gateway_id = "${sakuracloud_mobile_gateway.mgw.id}"
    
    prefix   = "192.168.10.0/26"
    next_hop = "192.168.11.202"
}

# モバイルゲートウェイに接続するスイッチの定義
resource "sakuracloud_switch" "sw" {
    name = "sw"
}
```

### パラメーター

|パラメーター         |必須  |名称                |初期値     |設定値                    |補足                                          |
|-------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `name`            | ◯   | 名称           | -        | 文字列                  | - |
| `internet_connection` | -   | インターネット接続  | `false` | `true`<br />`false`| - |
| `switch_id`       | -   | スイッチID  | - | 文字列 | - |
| `private_ipaddress`       | -   | プライベート側IPアドレス  | - | 文字列 | - |
| `private_nw_mask_len`       | -   | プライベート側ネットワークマスク長さ  | - | 数値(`8`〜`29`)| - |
| `dns_server1`       | -   | DNSサーバIP1アドレス  | - | 文字列 | - |
| `dns_server2`       | -   | DNSサーバ2IPアドレス  | - | 文字列 | - |
| `icon_id`         | -   | アイコンID         | - | 文字列| - |
| `description`     | -   | 説明  | - | 文字列 | - |
| `tags`            | -   | タグ | - | リスト(文字列) | - |
| `graceful_shutdown_timeout` | - | シャットダウンまでの待ち時間 | - | 数値(秒数) | シャットダウンが必要な場合の通常シャットダウンするまでの待ち時間(指定の時間まで待ってもシャットダウンしない場合は強制シャットダウンされる) |
| `zone`            | -   | ゾーン | - | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |

### 属性

|属性名                | 名称                    | 補足                                        |
|---------------------|------------------------|--------------------------------------------|
| `id`                | スイッチID               | -                                          |
| `public_ipaddress`  | パブリック側IPアドレス     | -                                          |
| `public_nw_mask_len`| パブリック側ネットワークマスク長 | -                                          |
| `sim_ids` | 接続されているSIMのIDリスト | -                                          |

## スタティックルート(sakuracloud_mobile_gateway_static_route)

モバイルゲートウェイに登録するスタティックルートを表します。

### パラメーター

|パラメーター                 |必須  |名称                 |初期値     |設定値                         |補足                                          |
|---------------------------|:---:|----------------------|:--------:|-------------------------------|----------------------------------------------|
| `mobile_gateway_id`   | ◯   | モバイルゲートウェイID         | -        | 文字列                   | - |
| `prefix`                    | ◯   | プリフィックス | -        | 文字列(`x.x.x.x/n`形式)                          | - |
| `next_hop`                  | ◯   | ネクストホップ | -        | 文字列                          | - |
| `zone`          | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                     | 名称             | 補足                  |
|--------------------------|------------------|----------------------|
| `id`                     | ID                    | -                    |
