# サーバ(sakuracloud_server)

---

### 設定例

```hcl
data "sakuracloud_server" "myserver" {
  name_selectors = ["foobar"]
}
```

### パラメーター

|パラメーター|必須  |名称                |初期値     |設定値 |補足                                          |
|----------|:---:|--------------------|:--------:|------|----------------------------------------------|
| `name_selectors`  | -   | 検索条件(名称)      | -        | リスト(文字列)           | 複数指定した場合はAND条件  |
| `tag_selectors`   | -   | 検索条件(タグ)      | -        | リスト(文字列)           | 複数指定した場合はAND条件  |
| `filter`          | -   | 検索条件(その他)    | -        | オブジェクト             | APIにそのまま渡されます。検索条件を指定してもAPI側が対応していない場合があります。 |
| `zone` | - | ゾーン | - | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |

### 属性

|属性名                    | 名称                     | 補足                                        |
|-------------------------|-------------------------|--------------------------------------------|
| `id`                    | ID                      | -                                          |
| `name`   | サーバ名           | - |
| `disks`  | ディスクID          | - | 
| `core`   | CPUコア数           | - | 
| `memory` | メモリ(GB単位)       | - | 
| `interface_driver`  | NICドライバ  | - |
| `nic`    | 基本NIC | - |
| `additional_nics` | 追加NIC | - |
| `packet_filter_ids`| パケットフィルタID | - |
| `icon_id`     | アイコンID         | - |
| `description` | 説明 | - |
| `cdrom_id` | CDROM(ISOイメージ)ID | - |
| `ipaddress`| 基本NIC:IPアドレス | - |
| `gateway`  | 基本NIC:ゲートウェイ | - |
| `nw_mask_len` | 基本NIC:サブネットマスク長 | - |
| `private_host_id` | 専有ホストID | - | 
| `tags` | タグ | - | 
| `zone` | ゾーン | - | 
| `macaddresses`          | MACアドレス               | MACアドレスのリスト(NICの個数分のリスト)        |
| `dns_servers`           | 基本NIC:DNSサーバ        | eth0の属するセグメントの推奨ネームサーバのリスト|
| `nw_address`            | 基本NIC:ネットワークアドレス | eth0のIPアドレスのネットワークアドレス          |
| `private_host_name`     | 専有ホスト名 | -          |
