# シンプル監視(sakuracloud_simple_monitor)

---

### 設定例

```hcl
data "sakuracloud_simple_monitor" "mymonitor" {
  name_selectors = ["foobar"]
}
```

### パラメーター

|パラメーター             |必須  |名称                |初期値     |設定値                    |補足                                      |
|-----------------------|:---:|--------------------|:--------:|------------------------|------------------------------------------|
| `name_selectors`  | -   | 検索条件(名称)      | -        | リスト(文字列)           | 複数指定した場合はAND条件  |
| `tag_selectors`   | -   | 検索条件(タグ)      | -        | リスト(文字列)           | 複数指定した場合はAND条件  |
| `filter`          | -   | 検索条件(その他)    | -        | オブジェクト             | APIにそのまま渡されます。検索条件を指定してもAPI側が対応していない場合があります。 |

### 属性

|属性名          | 名称             | 補足                                        |
|---------------|-----------------|--------------------------------------------|
| `id`                  | ID              | -                                          |
| `target`              | 監視対象名(IPアドレス) | - | 
| `health_check`        | 監視方法          | - | 
| `icon_id`             | アイコンID         | - | 
| `description`         | 説明             | - | 
| `tags`                | タグ             | - | 
| `notify_email_enabled`| Eメール通知有効    | - | 
| `notify_email_html`   | HTMLメール有効    | - | 
| `notify_slack_enabled`| Slack通知有効     | - | 
| `notify_slack_webhook`| Slack WebhookURL | - | 
| `enabled`             | 有効              | - | 

