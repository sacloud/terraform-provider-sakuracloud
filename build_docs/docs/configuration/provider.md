# プロバイダ設定

---

### 設定例

```hcl
provider sakuracloud {
  token  = "your API token"
  secret = "your API secret"
  zone   = "target zone"
  
  # retry_max       = 10
  # retry_interval  = 5   # 単位:秒
  # timeout         = 20  # 単位:分
  # api_root_url    = "https://secure.sakura.ad.jp/cloud/zone"  
  # trace           = false
}
```

### パラメーター

以下のパラメーターをサポートしています。

|パラメーター       |必須  |名称                |初期値     |設定値 |説明                                          |
|-----------------|:---:|--------------------|:--------:|------|----------------------------------------------|
|`token`          | ◯   | APIキー<br />(トークン)     | -        |文字列|環境変数`SAKURACLOUD_ACCESS_TOKEN`での指定も可         |
|`secret`         | ◯   | APIキー<br />(シークレット)  | -        |文字列|環境変数`SAKURACLOUD_ACCESS_TOKEN_SECRET`での指定も可  |
|`zone`           | -   | 対象ゾーン           | `is1b`   |`is1b`<br />`tk1a`<br />`tk1v`|環境変数`SAKURACLOUD_ZONE`での指定も可|
|`retry_max`      | -   | リトライ回数         | `10`     | 数値 |環境変数`SAKURACLOUD_RETRY_MAX`での指定も可  |
|`retry_interval` | -   | リトライ時待機時間    | `5`     |数値(秒)|環境変数`SAKURACLOUD_RETRY_INTERVAL`での指定も可  |
|`timeout`        | -   | タイムアウト         | `20`     | 数値(分) |環境変数`SAKURACLOUD_TIMEOUT`での指定も可|
|`api_root_url`   | -   | さくらのクラウドAPI ルートURL | -        |文字列|テストなどのためにAPIのルートAPIを変更したい場合に設定する。<br />末尾にスラッシュを含めないでください。<br />指定しない場合のルートURLは`https://secure.sakura.ad.jp/cloud/zone`<br />環境変数`SAKURACLOUD_API_ROOT_URL`での指定も可  |
|`trace`          | -   | トレースフラグ       | `false`     |`true`<br />`false`|(開発者向け)詳細ログの出力ON/OFFを指定します。 <br />環境変数`SAKURACLOUD_TRACE_MODE`での指定も可|
|`use_marker_tags`| -   | マーカータグ利用有無  | `false`     |`true`<br />`false`| (v1.4で廃止) Terraformで作成されたリソースであるかの判別用に全てのリソースにタグを付与するかを指定<br />環境変数`SAKURACLOUD_USE_MARKER_TAGS`での指定も可|
|`marker_tag_name`| -   | マーカータグ名       | `@terraform`|文字列| (v1.4で廃止) `use_marker_tags`がtrueの場合のタグ名 <br />環境変数`SAKURACLOUD_MARKER_TAG_NAME`での指定も可|

各パラメータとも環境変数での指定が可能です。

`token`と`secret`を環境変数で指定した場合、プロバイダ設定の記述は不要です。

#### 環境変数の指定例

```bash
$ export SAKURACLOUD_ACCESS_TOKEN="取得したAPIトークン"
$ export SAKURACLOUD_ACCESS_TOKEN_SECRET="取得したAPIシークレット"
```
