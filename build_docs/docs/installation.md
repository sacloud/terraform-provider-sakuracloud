# Installation / インストール

---

## Terraformのセットアップ

- 1) こちらの[公式サイト](https://www.terraform.io/downloads.html)からzipファイルのダウンロードを行います。
- 2) 適当なディレクトリに展開します。
- 3) 2)のディレクトリにパスを通します。

以下はMacでの例です。展開先ディレクトリは`~/terraform`として記載しています。

#### terraformインストール

```bash
# ~/terraformディレクトリ作成
$ mkdir -p ~/terraform ; cd ~/terraform
# ダウンロード
$ curl -sL https://releases.hashicorp.com/terraform/0.9.5/terraform_0.9.5_darwin_amd64.zip > terraform.zip
# 展開
$ unzip terraform.zip
# パスを通す
$ export PATH=$PATH:~/terraform/

```

### 動作確認

`terraform`コマンドを実行してみましょう。
以下のような表示がされればOKです。

#### terraform動作確認 

```bash
$ terraform
usage: terraform [--version] [--help] <command> [<args>]

Available commands are:
    apply       Builds or changes infrastructure
    destroy     Destroy Terraform-managed infrastructure
    get         Download and install modules for the configuration
    graph       Create a visual graph of Terraform resources
    init        Initializes Terraform configuration from a module
    output      Read an output from a state file
    plan        Generate and show an execution plan
    push        Upload this Terraform module to Atlas to run
    refresh     Update local state file against real resources
    remote      Configure remote state storage
    show        Inspect Terraform state or plan
    taint       Manually mark a resource for recreation
    validate    Validates the Terraform files
    version     Prints the Terraform version
```

## Terraform for さくらのクラウドのセットアップ

- 1) 以下のリンクから、プラットフォームに応じたzipファイルをダウンロードします
    - Windows: 
        - [i386](https://terraform.b.sakurastorage.jp/downloads/terraform-provider-sakuracloud_windows-386.zip)
        - [amd64](https://terraform.b.sakurastorage.jp/downloads/terraform-provider-sakuracloud_windows-amd64.zip)
    - Linux:
        - [i386](https://terraform.b.sakurastorage.jp/downloads/terraform-provider-sakuracloud_linux-386.zip)
        - [amd64](https://terraform.b.sakurastorage.jp/downloads/terraform-provider-sakuracloud_linux-amd64.zip)
    - Darwin(macOS)
        - [i386](https://terraform.b.sakurastorage.jp/downloads/terraform-provider-sakuracloud_darwin-386.zip)
        - [amd64](https://terraform.b.sakurastorage.jp/downloads/terraform-provider-sakuracloud_darwin-amd64.zip)
- 2) terraformバイナリと同じディレクトリに展開します

> Note:  
Terraform for さくらのクラウドのバイナリーファイルは[GitHub上のリリースページ](https://github.com/sacloud/terraform-provider-sakuracloud/releases/latest)でも配布しています。

#### terraform for さくらのクラウド インストール

```bash
$ cd ~/terraform
# ダウンロード
$ curl -sL https://terraform.b.sakurastorage.jp/downloads/terraform-provider-sakuracloud_darwin-amd64.zip > terraform-provider-sakuracloud.zip
# 展開
$ unzip terraform-provider-sakuracloud.zip

```


## さくらのクラウドAPIキーの取得

さくらのクラウドのコントロールパネルにログインしAPIキーを発行します。
以下を参考に実施してください。APIキーを発行したら、`ACCESS_TOKEN`と`ACCESS_TOKEN_SECRET`を控えておきましょう。

#### さくらのクラウド コントロールパネルへのログイン

![ログイン.png](images/login.png "ログイン.png")

#### さくらのクラウド(IaaS)を選択

![01_コンパネ.png](images/apikey01.png "01_コンパネ.png")

#### APIキー発行画面へ移動

![02_APIキー.png](images/apikey02.png "02_APIキー.png")

#### APIキーの発行

![03_APIキー.png](images/apikey03.png "03_APIキー.png")

#### 発行されたAPIキーの確認

![04_APIキー.png](images/apikey04.png "04_APIキー.png")

## Dockerでの実行

TerraformとTerraform for さくらのクラウドを同梱したDockerイメージを用意しています。

[Terraform for さくらのクラウド Dockerイメージ](https://hub.docker.com/r/sacloud/terraform/)

以下のように実行します。

#### Dockerでの実行
```bash
$ docker run -it --rm \
         -e SAKURACLOUD_ACCESS_TOKEN=[さくらのクラウド APIトークン] \
         -e SAKURACLOUD_ACCESS_TOKEN_SECRET=[さくらのクラウド APIシークレット] \
         -v $PWD:/work \
         sacloud/terraform apply
```

#### docker-composeでの実行
```bash
# あらかじめ以下コマンドで必要な設定ファイルをダウンロード/編集しておく
# curl -LO https://github.com/sacloud/terraform-for-sakuracloud-docker/raw/master/docker-compose.yml
# curl -L https://github.com/sacloud/terraform-for-sakuracloud-docker/raw/master/env-sample > .env

$ docker-compose run --rm terraform apply
```

## Homebrewでのインストール

macOSでHomebrewをご利用中の場合、以下のコマンドでTerraform for さくらのクラウドをインストール可能です。

    brew tap sacloud/homebrew-terraform-provider-sakuracloud ; brew install terraform-provider-sakuracloud

注: brewコマンドでのインストール後、プラグインの有効化のために画面の指示に従って`~/.terraformrc`に以下のような記述を行なってください。

    providers {
        sakuracloud = "terraform-provider-sakuracloud"
    }