# パケットフィルタルール(sakuracloud_packet_filter_rule)

---

パケットフィルタでのフィルタリングルールを表すリソースです。  
  
**このリソースはパケットフィルタリソース(sakuracloud_packet_filter)の`expressions`属性との共存はできません**


### 設定例

```hcl

resource "sakuracloud_packet_filter" "myfilter" {
  name        = "myfilter"
  description = "PacketFilter from terraform for SAKURA CLOUD"
}

resource "sakuracloud_packet_filter_rule" "rule1" {
  packet_filter_id = sakuracloud_packet_filter.myfilter.id

  protocol    = "tcp"
  source_nw   = "192.168.2.0/24"
  source_port = "0-65535"
  dest_port   = "80"
  allow       = true
  order       = 0
}

resource "sakuracloud_packet_filter_rule" "rule2" {
  packet_filter_id = sakuracloud_packet_filter.myfilter.id

  protocol    = "ip"
  allow       = false
  description = "Deny all"
  order       = 1
}

# 以下のようにcount構文と併用することで動的にルールの追加が行えます
#resource sakuracloud_packet_filter_rule "rules" {
#  packet_filter_id = sakuracloud_packet_filter.myfilter.id
#  protocol         = var.protocol[count.index]
#  source_nw        = var.source_nw[count.index]
#  source_port      = var.source_port[count.index]
#  dest_port        = var.dest_port[count.index]
#  allow            = var.allow[count.index]
#
#  order = count.index
#  count = 10
#}

```

### パラメーター

|パラメーター     |必須  |名称             |初期値     |設定値                    |補足                                          |
|---------------|:---:|----------------|:--------:|------------------------|----------------------------------------------|
| `packet_filter_id`| ◯   | パケットフィルタID | -        | 文字列| - |
| `protocol`    | ◯   | プロトコル       | -        | `tcp`<br />`udp`<br />`icmp`<br />`fragment`<br />`ip`| - |
| `source_nw`   | -   | 送信元ネットワーク | -       | `xxx.xxx.xxx.xxx`(IP)<br />`xxx.xxx.xxx.xxx/nn`(ネットワーク)<br />`xxx.xxx.xxx.xxx/yyy.yyy.yyy.yyy`(アドレス範囲)  | 空欄の場合はANY |
| `source_port` | -   | 送信元ポート      | -       | `0`〜`65535`の整数<br />`xx-yy`(範囲指定)<br />`0xPPPP/0xMMMM`(16進範囲指定) | 空欄の場合はANY |
| `dest_port`   | -   | 宛先ポート       | -        | `0`〜`65535`の整数<br />`xx-yy`(範囲指定)<br />`0xPPPP/0xMMMM`(16進範囲指定) | 空欄の場合はANY |
| `allow`       | -   | アクション       | `true`        | `true`<br />`false` | `true`の場合ALLOW動作<br />`false`の場合DENY動作 |
| `description` | -   | 説明            | -        | 文字列 | - |
| `order`       | -   | 並び順           | -        | 数値| 同一のパケットフィルタに対し同一のorderを持つルールを複数適用した場合の並び順は不定となります。 |
| `zone`            | -   | ゾーン | - | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |

**注意点**

同一のパケットフィルタに対し、以下の値が同一のルールを複数適用することはできません。  
(以下属性の値によってルールの同一性判定を行なっているため)

- `protocol`(プロトコル)   
- `source_nw`(送信元ネットワーク)
- `source_port`(送信元ポート)
- `dest_port`(宛先ポート)
- `allow`(アクション)
- `description`(説明)

### 属性

|属性名          | 名称             | 補足                                        |
|---------------|-----------------|--------------------------------------------|
| `id`          | ID              | -                                          |
