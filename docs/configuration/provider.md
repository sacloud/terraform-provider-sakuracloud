# プロバイダ設定

### 設定例
```
provider "sakuracloud" {
    token = "your API token"
    secret = "your API secret"
    zone = "target zone"
}
```

### パラメーター

以下のパラメーターをサポートしています。

|パラメーター|必須  |名称                |初期値     |設定値 |説明                                          |
|----------|:---:|--------------------|:--------:|------|----------------------------------------------|
|`token`   | ◯  |APIキー<br />(トークン)     | -        |文字列|環境変数`SAKURACLOUD_ACCESS_TOKEN`での指定も可         |
|`secret`  | ◯  |APIキー<br />(シークレット)  | -        |文字列|環境変数`SAKURACLOUD_ACCESS_TOKEN_SECRET`での指定も可  |
|`zone`    | -   | 対象ゾーン           | `is1b`   |`is1b`<br />`tk1a`<br />`tk1v`|環境変数`SAKURACLOUD_ZONE`での指定も可|
|`timeout` | -   | タイムアウト         | `20`     | 数値(分) |環境変数`SAKURACLOUD_TIMEOUT`での指定も可|
|`trace`   | -   | トレースフラグ       | `false`     |`true`<br />`false`|(開発者向け)詳細ログの出力ON/OFFを指定します。 <br />環境変数`SAKURACLOUD_TRACE_MODE`での指定も可|

各パラメータとも環境変数での指定が可能です。

`token`と`secret`を環境変数で指定した場合、プロバイダ設定の記述は不要です。

#### 環境変数の指定例

```bash
$ export SAKURACLOUD_ACCESS_TOKEN="取得したAPIトークン"
$ export SAKURACLOUD_ACCESS_TOKEN_SECRET="取得したAPIシークレット"
```