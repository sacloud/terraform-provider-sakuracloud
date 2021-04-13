# オブジェクトストレージ対応

## 概要

2021/04に正式サービス化された新オブジェクトストレージへの対応を行う。

## 背景/前提条件

さくらのクラウドでは2021/04以前までにもオブジェクトストレージサービスを展開していた。(旧オブジェクトストレージ)  
このプロバイダーではv2.8.3時点では旧オブジェクトストレージの一部機能にのみ対応している。  

このため、新オブジェクトストレージへの対応が早急に必要となっている。

## 対応方針/方法

### 新オブジェクトストレージの特徴

[https://manual.sakura.ad.jp/cloud/objectstorage/about.html](https://manual.sakura.ad.jp/cloud/objectstorage/about.html)

実装にあたり注意すべき点を以下に転載する。

### 料金

- 月額課金と従量課金の2段階の課金
- 月額課金は以下項目に対するもの。超過分については従量課金となる。
  - ストレージ容量: 100GiB
  - 転送量: 10GiB
  - リクエスト数: 100,000リクエスト
  - オブジェクト数: 100,000オブジェクト
- リクエスト数/転送量については以下のルールが適用される
  - レスポンスのステータスコードが`100`/`200`/`204`/`206`の場合のみカウントされる  
  - さくらインターネットのIPアドレスレンジ以外からのGETリクエストのみ転送量がカウントされる  

### サービス仕様の留意点

- 最大同時アクセス数: 10
- 認証: Signature v4対応
- キャッシュサーバの有無: 無し
- 初期状態ではリソース数に一定の上限あり(緩和はサポート経由となる)

!!! Note  
    - 同時アクセス制限の実装が必要  
    - 従来のSignature v4未対応という制約がなくなったためライブラリの選定に幅が出来た
    - リソース数上限に対する実装での対応が必要
  
### 実装されているAPI

[https://manual.sakura.ad.jp/cloud/objectstorage/api.html](https://manual.sakura.ad.jp/cloud/objectstorage/api.html)

- バケットの作成はAPI経由では行えない
- APIでの操作には通常のAPIキーに加えバケットの作成時に表示されるアクセスキーが必要

### 実装方針案

大まかに2つの方針がある。

- awsなどの他のプロバイダーにまかせ自前実装しない
- aws sdkなどを用いて自前実装

#### AWSプロバイダーを利用する場合

以下のように利用する。

```tf
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  alias      = "sacloud"
  region     = var.s3_region
  access_key = var.s3_access_key
  secret_key = var.s3_secret_key
  endpoints {
    s3 = var.s3_endpoint
  }

  skip_credentials_validation = true
  skip_region_validation      = true
  skip_requesting_account_id  = true
  skip_metadata_api_check     = true
}

#================================================

variable "s3_access_key" {}
variable "s3_secret_key" {}

variable "s3_region" {
  default = "jp-north-1"
}

variable "s3_endpoint" {
  default = "https://s3.isk01.sakurastorage.jp"
}

#================================================

# 作成済みのバケットを参照
data "aws_s3_bucket" "example" {
  provider = aws.sacloud

  bucket   = "aws-provider-test"
}

# オブジェクトの作成
resource "aws_s3_bucket_object" "object" {
  provider = aws.sacloud

  bucket = data.aws_s3_bucket.example.id
  key    = "example.txt"
  source = "example.txt"
  etag   = filemd5("example.txt")
  acl    = "public-read"
}
```

#### メリット/デメリット

- AWSプロバイダーをそのまま利用可能なためメンテナンスコストが下げられる
- API経由でバケット作成が行えないため、バケットに対する操作が行えない
    - ACL
    - バージョニング
    - CORS
    
