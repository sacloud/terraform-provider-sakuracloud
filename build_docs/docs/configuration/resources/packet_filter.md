# パケットフィルタ(sakuracloud_packet_filter)

---

### 設定例

```hcl
resource "sakuracloud_packet_filter" "myfilter" {
  name = "myfilter"

  expressions {
    protocol    = "tcp"
    source_nw   = "192.168.2.0/24"
    source_port = "0-65535"
    dest_port   = "80"
    allow       = true
  }

  expressions {
    protocol    = "ip"
    source_nw   = "0.0.0.0/0"
    allow       = false
    description = "Deny all"
  }

  description = "PacketFilter from terraform for SAKURA CLOUD"
}
```

### パラメーター

|パラメーター         |必須  |名称                |初期値     |設定値                    |補足                                          |
|-------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `name`            | ◯   | スイッチ名           | -        | 文字列                  | - |
| `expressions`     | ◯   | フィルタルール        | -        | リスト(マップ)           | 詳細は[`expressions`](#expressions)を参照 |
| `description`     | -   | 説明  | - | 文字列 | - |
| `tags`            | -   | タグ | - | リスト(文字列) | - |
| `zone`            | -   | ゾーン | - | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |

**注意**  

同一のパケットフィルタに対し`expressions`属性と`sakuracloud_packet_filter_rule`リソースの併用はできません。

#### `expressions`

|パラメーター     |必須  |名称             |初期値     |設定値                    |補足                                          |
|---------------|:---:|----------------|:--------:|------------------------|----------------------------------------------|
| `protocol`    | ◯   | プロトコル       | -        | `tcp`<br />`udp`<br />`icmp`<br />`fragment`<br />`ip`| - |
| `source_nw`   | -   | 送信元ネットワーク | -       | `xxx.xxx.xxx.xxx`(IP)<br />`xxx.xxx.xxx.xxx/nn`(ネットワーク)<br />`xxx.xxx.xxx.xxx/yyy.yyy.yyy.yyy`(アドレス範囲)  | 空欄の場合はANY |
| `source_port` | -   | 送信元ポート      | -       | `0`〜`65535`の整数<br />`xx-yy`(範囲指定)<br />`0xPPPP/0xMMMM`(16進範囲指定) | 空欄の場合はANY |
| `dest_port`   | -   | 宛先ポート       | -        | `0`〜`65535`の整数<br />`xx-yy`(範囲指定)<br />`0xPPPP/0xMMMM`(16進範囲指定) | 空欄の場合はANY |
| `allow`       | -   | アクション       | `true`        | `true`<br />`false` | `true`の場合ALLOW動作<br />`false`の場合DENY動作 |
| `description` | -   | 説明            | -        | 文字列 | - |


### 属性

|属性名          | 名称             | 補足                                        |
|---------------|-----------------|--------------------------------------------|
| `id`          | ID              | -                                          |
