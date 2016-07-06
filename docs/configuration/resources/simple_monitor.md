# シンプル監視(sakuracloud_simple_monitor)

### 設定例

```
resource "sakuracloud_simple_monitor" "mymonitor" {
    target = "${sakuracloud_server.myserver.base_nw_ipaddress}"
    health_check = {
        protocol = "http"
        delay_loop = 60
        path = "/"
        status = "200"
    }
    notify_email_enabled = true
    notify_slack_enabled = true
    notify_slack_webhook = "https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
}
```

### パラメーター

|パラメーター             |必須  |名称                |初期値     |設定値                    |補足                                      |
|-----------------------|:---:|--------------------|:--------:|------------------------|------------------------------------------|
| `target`              | ◯   | 監視対象名(IPアドレス) | -    | 文字列                  | 監視対象のFQDNまたはIPアドレス |
| `health_check`        | ◯   | 監視方法          | -       | マップ           | 詳細は[`health_check`](#health_check)を参照 |
| `description`         | -   | 説明             | -       | 文字列 | - |
| `tags`                | -   | タグ             | -       | リスト(文字列) | - |
| `notify_email_enabled`| -   | Eメール通知有効    | `true`  | `true`<br />`false` | - |
| `notify_slack_enabled`| -   | Slack通知有効     | `false` | `true`<br />`false` | - |
| `notify_slack_webhook`| -   | Slack WebhookURL | -       | 文字列 | - |
| `enabled`             | -   | 有効              | `true` | `true`<br />`false` | - |

#### `health_check`

|パラメーター      |必須  |名称                |初期値     |設定値                    |補足                                          |
|----------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `protocol`     | ◯   | プロトコル        | -        | `http`<br />`https`<br />`ping`<br />`tcp`<br />`dns`<br />`ssh`<br />`smtp`<br />`pop3`<br />`snmp`| - |
| `dalay_loop`   | -   | チェック間隔(秒)        | `60`        | 数値                  | `60`〜`3600` |
| `path`         | △   | パス  | - | 文字列 | プロトコルが`http`または`https`の場合のみ有効かつ必須 |
| `host_header`  | △   | HOSTヘッダ  | - | 文字列 | プロトコルが`http`または`https`の場合のみ有効 |
| `status`       | △   | レスポンスコード | - | 文字列 | プロトコルが`http`または`https`の場合のみ有効かつ必須 |
| `port`         | △   | ポート番号 | - | 数値 | プロトコルが`tcp`,`ssh`,`smtp`,`pop3`の場合のみ有効かつ必須 |
| `qname`        | △   | 問合せFQDN | - | 数値 | プロトコルが`dns`の場合のみ有効かつ必須 |
| `expected_data`| △   | 期待値 | - | 数値 | プロトコルが`dns`,`snmp`の場合のみ有効<br />`dns`の場合、省略すると、何らかのAレコードの応答があるかのチェックとなる<br />`snmp`の場合は必須 |
| `community`    | △   | コミュニティ名 | - | 文字列 | プロトコルが`snmp`の場合のみ有効かつ必須 |
| `snmp_version` | △   | SNMPバージョン | - | `1`<br />`2c` | プロトコルが`snmp`の場合のみ有効かつ必須 |
| `oid`          | △   | OID | - | 文字列 | プロトコルが`snmp`の場合のみ有効かつ必須 |


### 属性

|属性名          | 名称             | 補足                                        |
|---------------|-----------------|--------------------------------------------|
| `id`                   | ID              | -                                          |
| `target`               | 監視対象名(IPアドレス)| -                                          |
| `health_check`         | 監視方法          | 詳細は[`health_check`](#health_check)を参照 |
| `description`          | 説明             | -                                          |
| `tags`                 | タグ             | -                                          |
| `notify_email_enabled` | Eメール通知有効    | -                                          |
| `notify_slack_enabled` | Slack通知有効     | -                                          |
| `notify_slack_webhook` | Slack WebhookURL| -                                          |
| `enabled`              | 有効             | -                                          |
