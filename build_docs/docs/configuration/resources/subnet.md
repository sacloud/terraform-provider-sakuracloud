# サブネット(sakuracloud_subnet)

---

ルータ(`sakuracloud_internet`)に追加可能なグローバルIPアドレスブロックを表すリソースです。  

### 設定例

```hcl
# ルータの定義
resource "sakuracloud_internet" "foobar" {
  name = "router"
}

# ルータに追加するグローバルIPアドレスブロック
resource "sakuracloud_subnet" "foobar" {
  # ルータのID
  internet_id = sakuracloud_internet.foobar.id

  # ネクストホップ
  next_hop = sakuracloud_internet.foobar.nw_min_ipaddress
  
  # ネットワークマスク長
  #nw_mask_len = 28
}

```

### パラメーター

|パラメーター         |必須  |名称                |初期値     |設定値                    |補足                                          |
|-------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `internet_id`     | ◯   | ルータID           | -        | 文字列                  | - |
| `nw_mask_len`     | -   | ネットワークマスク長  | `28` | `28`<br />`27`<br />`26` | グローバルIPのプリフィックス(ネットワークマスク長) |
| `next_hop`        | -   | ネクストホップ| - | 文字列 | ネクストホップのIPv4アドレス |
| `zone`            | -   | ゾーン | - | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |

### 属性

|属性名                | 名称                    | 補足                                        |
|---------------------|------------------------|--------------------------------------------|
| `id`                | スイッチID               | -                                          |
| `switch_id`         | スイッチID              | (内部的に)接続されているスイッチID              |
| `nw_address`        | ネットワークアドレス      | 割り当てられたグローバルIPのネットワークアドレス |
| `min_ipaddress`  | 最小IPアドレス           | 割り当てられたグローバルIPアドレスのうち、利用可能な先頭IPアドレス |
| `max_ipaddress`  | 最大IPアドレス           | 割り当てられたグローバルIPアドレスのうち、利用可能な最後尾IPアドレス |
| `ipaddresses`    | IPアドレスリスト         | 割り当てられたグローバルIPアドレスのうち、利用可能なIPアドレスのリスト |
