# ブリッジ(sakuracloud_bridge)

---

**全ゾーン共通のグローバルリソースです。**

### 設定例

```hcl
#ブリッジの定義
resource "sakuracloud_bridge" "bridge" {
  name        = "bridge"
  description = "example bridge"
}

# ブリッジに接続するスイッチ(東京第1ゾーン)
resource "sakuracloud_switch" "tk1a" {
  name      = "switch_tk1a"
  bridge_id = sakuracloud_bridge.bridge.id
  zone      = "tk1a"
}

# ブリッジに接続するスイッチ(石狩第2ゾーン)
resource "sakuracloud_switch" "is1b" {
  name      = "switch_is1b"
  bridge_id = sakuracloud_bridge.bridge.id
  zone      = "is1b"
}
```

### パラメーター

|パラメーター         |必須  |名称                |初期値     |設定値                    |補足                                          |
|-------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `name`            | ◯   | ブリッジ名           | -        | 文字列                  | - |
| `description`     | -   | 説明  | - | 文字列 | - |

### 属性

|属性名                | 名称                    | 補足                                        |
|---------------------|------------------------|--------------------------------------------|
| `id`                | ブリッジID               | -                                          |
