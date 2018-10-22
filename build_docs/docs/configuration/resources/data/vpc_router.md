# VPCルータ(sakuracloud_vpc_router)

---

### 設定例

```hcl
data "sakuracloud_vpc_router" "foobar" {
  name_selectors = ["foobar"]
}
```

## `sakuracloud_vpc_router`

VPCルータ本体を表します。

### パラメーター

|パラメーター       |必須  |名称           |初期値     |設定値                         |補足                                          |
|-----------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `name_selectors`  | -   | 検索条件(名称)      | -        | リスト(文字列)           | 複数指定した場合はAND条件  |
| `tag_selectors`   | -   | 検索条件(タグ)      | -        | リスト(文字列)           | 複数指定した場合はAND条件  |
| `filter`          | -   | 検索条件(その他)    | -        | オブジェクト             | APIにそのまま渡されます。検索条件を指定してもAPI側が対応していない場合があります。 |
| `zone`          | -   | ゾーン          | -        | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名          | 名称             | 補足                  |
|---------------|------------------|----------------------|
| `id`            | ID             | -                    |
| `name`          | ロードバランサ名 | - |
| `plan`          | プラン          | - |
| `switch_id`     | スイッチID      | - |
| `vip`           | IPアドレス1     | - |
| `ipaddress1`    | IPアドレス1     | - |
| `ipaddress2`    | IPアドレス2     | - |
| `vrid`          | VRID           | - |
| `aliases`       | IPエイリアス    | - |
| `syslog_host`   | syslog転送先ホスト| - |
| `internet_connection` | インターネット接続 | - |
| `icon_id`       | アイコンID         | - |
| `description`   | 説明           | - |
| `tags`          | タグ           | - |
| `global_address`| グローバルIP     | - |
