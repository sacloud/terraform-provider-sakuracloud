# オブジェクトストレージ(sakuracloud_bucket_object)

---

### 設定例

```hcl
data "sakuracloud_bucket_object" "foobar" {
  #アクセスキー(環境変数SACLOUD_OJS_ACCESS_KEY_ID、またはAWS_ACCESS_KEY_IDでも指定可能)  
  #access_key = ""

  #シークレットキー(環境変数SACLOUD_OJS_SECRET_ACCESS_KEY、またはAWS_SECRET_ACCESS_KEYでも指定可能)
  #secret_key = ""

  #バケット名
  bucket = "your-bucket-name"

  #キー
  key = "your/object/key.txt"
}
```

### パラメーター

|パラメーター         |必須  |名称            |初期値  |設定値                    |補足                                          |
|-------------------|:---:|----------------|:-----:|------------------------|----------------------------------------------|
| `access_key`      | ◯   | アクセスキー     | -    | 文字列                  | 環境変数`SACLOUD_OJS_ACCESS_KEY_ID`、または`AWS_ACCESS_KEY_ID`でも指定可能 |
| `secret_key`      | ◯   | シークレットキー | -     | 文字列                  | 環境変数`SACLOUD_OJS_SECRET_ACCESS_KEY`、または`AWS_SECRET_ACCESS_KEY`でも指定可能 |
| `bucket`          | ◯   | バケット名      | -     | 文字列                  | - |
| `key`             | ◯   | オブジェクトキー | -     | 文字列                  | - |

### 属性

|属性名                | 名称                    | 補足                              |
|---------------------|------------------------|--------------------------------------------|
| `id`                | オブジェクトキー           | `key`と同じ値                    |
| `content_type`      | コンテントタイプ           |                                 |
| `size`              | サイズ(byte)              |                                |
| `last_modified`     | 更新日時                   |                                |
| `http_url`          | バーチャルホスト型URL(HTTP) |                                |
| `https_url`         | バーチャルホスト型URL(HTTPS)|                                |
| `http_path_url`     | パス型URL(HTTP)           |                                |
| `https_path_url`    | パス型URL(HTTPS)          |                                |
| `http_cache_url`    | キャッシュ配信URL(HTTP)    |                                |
| `https_cache_url`   | キャッシュ配信URL(HTTPS)   |                                |
| `body`              | ボディ                  | コンテントタイプが`text/*`または`application/json`の場合のみ利用可能です |
