# ルータ(sakuracloud_internet)

---

### 設定例

```hcl
resource "sakuracloud_internet" "router" {
  name = "router"

  #ネットワークマスク超
  #nw_mask_len = 28

  #帯域幅(Mbps単位)
  #band_width = 100

  #IPv6有効化
  #enable_ipv6 = false

  description = "Router from terraform for SAKURA CLOUD"
  tags        = ["tag1", "tag2"]
}
```

### パラメーター

|パラメーター         |必須  |名称                |初期値     |設定値                    |補足                                          |
|-------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `name`            | ◯   | ルータ名           | -        | 文字列                  | - |
| `nw_mask_len`     | -   | ネットワークマスク長  | `28` | `28`<br />`27`<br />`26` | グローバルIPのプリフィックス(ネットワークマスク長) |
| `band_width`      | -   | 帯域幅(Mbps単位)  | `100` | `100`<br />`250`<br />`500`<br />`1000`<br />`1500`<br />`2000`<br />`2500`<br />`3000` | - |
| `enable_ipv6`     | -   | IPv6有効化  | - | `true`<br />`false`| - |
| `icon_id`         | -   | アイコンID         | - | 文字列 | - |
| `description`     | -   | 説明  | - | 文字列 | - |
| `tags`            | -   | タグ | - | リスト(文字列) | - |
| `graceful_shutdown_timeout` | - | シャットダウンまでの待ち時間 | - | 数値(秒数) | シャットダウンが必要な場合の通常シャットダウンするまでの待ち時間(指定の時間まで待ってもシャットダウンしない場合は強制シャットダウンされる) |
| `zone`            | -   | ゾーン | - | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |

### 属性

|属性名                | 名称                    | 補足                                        |
|---------------------|------------------------|--------------------------------------------|
| `id`                | ルータID               | -                                          |
| `server_ids`         | サーバID              | 接続されているサーバのID(リスト)             |
| `switch_id`          | スイッチID              | (内部的に)接続されているスイッチID              |
| `nw_address`         | ネットワークアドレス      | ルータに割り当てられたグローバルIPのネットワークアドレス |
| `gateway`         | ゲートウェイ             | ルータに割り当てられたセグメントのゲートウェイIPアドレス |
| `min_ipaddress`   | 最小IPアドレス           | ルータに割り当てられたグローバルIPアドレスのうち、利用可能な先頭IPアドレス [注1](#ルータ-sakuracloud_internet_属性_注1) |
| `max_ipaddress`   | 最大IPアドレス           | ルータに割り当てられたグローバルIPアドレスのうち、利用可能な最後尾IPアドレス [注1](#ルータ-sakuracloud_internet_属性_注1) |
| `ipaddresses`     | IPアドレスリスト         | ルータに割り当てられたグローバルIPアドレスのうち、利用可能なIPアドレスのリスト [注1](#ルータ-sakuracloud_internet_属性_注1)|
| `ipv6_prefix`        | IPv6アドレスプレフィックス| -              |
| `ipv6_prefix_len`    | IPv6アドレスプレフィックス長 | -             |
| `ipv6_nw_address`    | IPv6ネットワークアドレス     | -             |

#### 注1

割り当てられたグローバルIPのうち、先頭の4個はネットワークアドレスやルータ用のため、最後尾の1個はブロードキャストアドレスのため使用できません。

詳細は[こちら](http://cloud-news.sakura.ad.jp/faq_top/faq/#H004)を参照ください。

`min_ipaddress`と`max_ipaddress`、`ipaddresses`には利用可能なIPアドレスが設定されています。
