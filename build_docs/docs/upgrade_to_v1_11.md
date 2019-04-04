# v1.11での変更点

---

## Terraform v0.12対応

このバージョンからTerraform v0.12に対応しています。  
Terraform v0.12ではtfファイル記法などに大幅な変更が行われました。  
変更内容の詳細については[公式ドキュメント](https://terraform.io/docs/)を参照してください。

このドキュメントに記載しているtfファイルの例はTerraform v0.12対応となっています。

Terraform v0.12以前のバージョンをご利用中の場合、Terraform v0.12を用いて以下のコマンドを実行することでtfファイルのマイグレーションが行えます。

```bash
# terraform v0.12以前のバージョン向けのtfファイルをv0.12向けに書き換え
$ terraform 0.12upgrade
```

## marker_tags機能の除去
   
v1.4でmarker_tags機能のソースが除去されました。 

## ディスクの修正関連パラメータの移動

v1.6でディスクの修正関連パラメータがディスクリソース(`sakuracloud_disk`)からサーバリソース(`sakuracloud_server`)へと移動されました。
v1.11まではディスクリソースでのパラメータ指定を行うことも可能でしたが、このバージョンではこれらのパラメータは指定できなくなりました。  

### 対応方法

これまでディスクリソースで以下のパラメータを指定していた場合、tfファイルの書き換えが必要となります。

- パスワード(`password`)
- ホスト名(`hostname`)
- SSH接続時パスワード/チャレンジレスポンス認証の無効化フラグ(`disable_pw_auth`)
- スタートアップスクリプト(`note_ids`)
- 公開鍵(`ssh_key_ids`)

これらのパラメータをディスクリソースからサーバリソースへ移動させてください。

#### 対応前のtfファイルの例

```hcl
data "sakuracloud_archive" "ubuntu" {
  os_type = "ubuntu"
}

# ディスクリソース
resource "sakuracloud_disk" "foobar" {
  name = "mydisk"
  source_archive_id = data.sakuracloud_archive.ubuntu.id
 
  # ディスクの修正関連のパラメータ 
  hostname = "myserver"
  password = "p@ssw0rd"
  ssh_key_ids = ["100000000000", "200000000000"]
  disable_pw_auth = true
  note_ids = ["100000000000", "200000000000"]

}

# サーバリソース
resource "sakuracloud_server" "foobar" {
  name = "myserver"
  disks = [sakuracloud_disk.foobar.id]
}
```

#### 対応後のtfファイルの例

```hcl
data "sakuracloud_archive" "ubuntu" {
  os_type = "ubuntu"
}

# ディスクリソース
resource "sakuracloud_disk" "foobar" {
  name = "mydisk"
  source_archive_id = data.sakuracloud_archive.ubuntu.id
}

# サーバリソース
resource "sakuracloud_server" "foobar" {
  name = "myserver"
  disks = [sakuracloud_disk.foobar.id]
  
  # ディスクの修正関連のパラメータ 
  hostname = "myserver"
  password = "p@ssw0rd"
  ssh_key_ids = ["100000000000", "200000000000"]
  disable_pw_auth = true
  note_ids = ["100000000000", "200000000000"]

}
```

## NFS SSDプラン対応

NFSにSSDプランが導入された際にパラメータ名が変更となりました。
v1.11以前でNFSアプライアンスを利用していた場合はtfファイルの変更が必要となります。  

旧: `plan`: NFSのサイズをGB単位で指定
新: `plan`で`hdd` or `ssd`を指定、サイズは`size`パラメータで指定

### 対応方法

- tfファイルで`plan`を指定していた場合は`size`に置き換えてください。

#### 対応前のtfファイルの例

```hcl
resource "sakuracloud_nfs" "foobar" {
  switch_id   = "${sakuracloud_switch.foobar.id}"
  plan        = "100"
  
  # ...
}
```

#### 対応後のtfファイルの例

```hcl
resource "sakuracloud_nfs" "foobar" {
  switch_id   = "${sakuracloud_switch.foobar.id}"
  size        = "100" # planをsizeに置き換える
  
  # ...
}
```

