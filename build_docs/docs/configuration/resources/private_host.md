# 専有ホスト(sakuracloud_private_host)

---

### 設定例

```hcl
resource "sakuracloud_private_host" "private_host" {
    name = "private_host"
    description = "PrivateHost from terraform for SAKURA CLOUD"
    tags = ["tag1" , "tag2"]
}
```

### パラメーター

|パラメーター         |必須  |名称                |初期値     |設定値                    |補足                                          |
|-------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `name`            | ◯   | 専有ホスト名           | -        | 文字列                  | - |
| `description`     | -   | 説明  | - | 文字列 | - |
| `tags`            | -   | タグ | - | リスト(文字列) | - |
| `zone`            | -   | ゾーン | - | `is1b`<br />`tk1a` | - |

### 属性

|属性名                | 名称                    | 補足                                        |
|---------------------|------------------------|--------------------------------------------|
| `id`                | 専有ホストID               | -                                          |
| `name`              | 専有ホスト名              | -                                          |
| `description`       | 説明                    | -                                          |
| `tags`              | タグ                    | -                                          |
| `zone`              | ゾーン                  | -                                          |
| `assigned_core`     | 割当済みコア数           | -                                          |
| `assigned_memory`   | 割当済みメモリ(GB単位)    | -                                          |
