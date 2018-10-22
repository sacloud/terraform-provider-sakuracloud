# オブジェクトストレージ(sakuracloud_bucket_object)

さくらのクラウド オブジェクトストレージ上のオブジェクトを扱うためのリソースです。  

[Note]  
さくらのクラウドではAPIによるバケットの生成がサポートされていません。   
このため、当リソースを利用する場合、あらかじめコントロールパネルからバケットを作成しておく必要があります。  

- このリソースは`import`非対応です。

---

### 設定例

```hcl
resource "sakuracloud_bucket_object" "foobar" {
  #アクセスキー(環境変数SACLOUD_OJS_ACCESS_KEY_ID、またはAWS_ACCESS_KEY_IDでも指定可能)  
  #access_key = ""

  #シークレットキー(環境変数SACLOUD_OJS_SECRET_ACCESS_KEY、またはAWS_SECRET_ACCESS_KEYでも指定可能)
  #secret_key = ""

  #バケット名
  bucket = "your-bucket-name"

  #キー
  key = "your/object/key.txt"

  #ボディ:
  #文字列で指定する場合
  content = "content"

  #ファイルを指定する場合
  #source = "path/your/object/file"
  
  #ETag: オブジェクトストレージ上でオブジェクトが変更されたことを検知したい場合に指定
  #etag = md5(file("path/your/object/file"))
}
```

### パラメーター

|パラメーター         |必須  |名称            |初期値  |設定値                    |補足                                          |
|-------------------|:---:|----------------|:-----:|------------------------|----------------------------------------------|
| `access_key`      | ◯   | アクセスキー     | -    | 文字列                  | 環境変数`SACLOUD_OJS_ACCESS_KEY_ID`、または`AWS_ACCESS_KEY_ID`でも指定可能 |
| `secret_key`      | ◯   | シークレットキー | -     | 文字列                  | 環境変数`SACLOUD_OJS_SECRET_ACCESS_KEY`、または`AWS_SECRET_ACCESS_KEY`でも指定可能 |
| `bucket`          | ◯   | バケット名      | -     | 文字列                  | - |
| `key`             | ◯   | オブジェクトキー | -     | 文字列                  | - |
| `content`         | △   | ボディ(文字列)   | -    | 文字列                  | オブジェクトの内容を文字列で指定する場合に利用する。`source`との併用は不可 |
| `source`          | △   | ボディ(ファイル) | -     | 文字列                 | オブジェクトの内容をファイルで指定する場合に利用する。`content`との併用は不可 |
| `etag`            | -   | ETag           | -    | 文字列          | オブジェクトの内容のMD5ハッシュ、オブジェクトストレージ上でオブジェクトが変更されたことを検知したい場合に指定する |

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

### 属性(データリソースでのみ利用可能)

以下の属性はデータリソースでのみ利用可能です。

|属性名                | 名称                    | 補足                              |
|---------------------|------------------------|--------------------------------------------|
| `body`              | ボディ                  | コンテントタイプが`text/*`または`application/json`の場合のみ利用可能です |


### データリソースの利用例

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