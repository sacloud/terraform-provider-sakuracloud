# ブリッジ(sakuracloud_bridge)

---

**全ゾーン共通のグローバルリソースです。**

### 設定例

```hcl
#ブリッジの参照
data "sakuracloud_bridge" "br01" {
  name_selectors = ["br01"]
}

# ブリッジに接続するスイッチ(東京第1ゾーン)
resource "sakuracloud_switch" "sw_tk1a" {
  name      = "switch_tk1a"
  bridge_id = data.sakuracloud_bridge.br01.id
  zone      = "tk1a"
}

# ブリッジに接続するスイッチ(石狩第2ゾーン)
resource "sakuracloud_switch" "sw_is1b" {
  name      = "switch_is1b"
  bridge_id = data.sakuracloud_bridge.br01.id
  zone      = "is1b"
}

```

### パラメーター

|パラメーター         |必須  |名称                |初期値     |設定値                    |補足                                          |
|-------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `name_selectors`  | -   | 検索条件(名称)      | -        | リスト(文字列)           | 複数指定した場合はAND条件  |
| `filter`          | -   | 検索条件(その他)    | -        | オブジェクト             | APIにそのまま渡されます。検索条件を指定してもAPI側が対応していない場合があります。 |

### 属性

|属性名                | 名称                    | 補足                                        |
|---------------------|------------------------|--------------------------------------------|
| `id`                | ブリッジID               | -                                          |
| `name`              | ブリッジ名           | - |
| `description`       | 説明  | - |
