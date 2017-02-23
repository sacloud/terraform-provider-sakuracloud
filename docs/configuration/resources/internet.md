# ルーター(sakuracloud_internet)

### 設定例

```
resource "sakuracloud_internet" "myrouter" {
    name = "myrouter"
    description = "Router from terraform for SAKURA CLOUD"
    tags = ["hoge1" , "hoge2"]
    nw_mask_len = 28
    band_width = 100
}
```

### パラメーター

|パラメーター         |必須  |名称                |初期値     |設定値                    |補足                                          |
|-------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `name`            | ◯   | ルーター名           | -        | 文字列                  | - |
| `nw_mask_len`     | -   | ネットワークマスク長  | `28` | `28`<br />`27`<br />`26` | グローバルIPのプリフィックス(ネットワークマスク長) |
| `band_width`      | -   | 帯域幅(Mbps単位)  | `100` | `100`<br />`250`<br />`500`<br />`1000`<br />`1500`<br />`2000`<br />`2500`<br />`3000` | - |
| `description`     | -   | 説明  | - | 文字列 | - |
| `tags`            | -   | タグ | - | リスト(文字列) | - |
| `zone`            | -   | ゾーン | - | `is1b`<br />`tk1a`<br />`tk1v` | - |

### 属性

|属性名                | 名称                    | 補足                                        |
|---------------------|------------------------|--------------------------------------------|
| `id`                | ルーターID               | -                                          |
| `name`              | ルーター名               | -                                          |
| `description`       | 説明                    | -                                          |
| `tags`              | タグ                    | -                                          |
| `zone`              | ゾーン                  | -                                          |
| `server_ids`         | サーバーID              | 接続されているサーバーのID(リスト)             |
| `switch_id`          | スイッチID              | (内部的に)接続されているスイッチID              |
| `nw_address`         | ネットワークアドレス      | ルーターに割り当てられたグローバルIPのネットワークアドレス |
| `nw_gateway`         | ゲートウェイ             | ルーターに割り当てられたセグメントのゲートウェイIPアドレス |
| `nw_min_ipaddress`   | 最小IPアドレス           | ルーターに割り当てられたグローバルIPアドレスのうち、利用可能な先頭IPアドレス [注1](#ルーター-sakuracloud_internet_属性_注1) |
| `nw_max_ipaddress`   | 最大IPアドレス           | ルーターに割り当てられたグローバルIPアドレスのうち、利用可能な最後尾IPアドレス [注1](#ルーター-sakuracloud_internet_属性_注1) |
| `nw_ipaddresses`     | IPアドレスリスト         | ルーターに割り当てられたグローバルIPアドレスのうち、利用可能なIPアドレスのリスト [注1](#ルーター-sakuracloud_internet_属性_注1)|

#### 注1

割り当てられたグローバルIPのうち、先頭の4個はネットワークアドレスやルータ用のため、最後尾の1個はブロードキャストアドレスのため使用できません。

詳細は[こちら](http://cloud-news.sakura.ad.jp/faq_top/faq/#H004)を参照ください。

`nw_min_ipaddress`と`nw_max_ipaddress`、`nw_ipaddresses`には利用可能なIPアドレスが設定されています。
