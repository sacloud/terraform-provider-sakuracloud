# v1.6での変更点

---

## ディスクの修正関連パラメータの移動

v1.6ではディスクの修正関連パラメータがディスクリソース(`sakuracloud_disk`)からサーバリソース(`sakuracloud_server`)へと移動されました。

これまでのバージョンとの互換性維持のため引き続きディスクリソースでのパラメータ指定を行うことも可能ですが、
将来のバージョンではこれらのパラメータを削除する予定です。


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
  source_archive_id = "${data.sakuracloud_archive.ubuntu.id}"
 
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
  disks = ["${sakuracloud_disk.foobar.id}"]
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
  source_archive_id = "${data.sakuracloud_archive.ubuntu.id}"
}

# サーバリソース
resource "sakuracloud_server" "foobar" {
  name = "myserver"
  disks = ["${sakuracloud_disk.foobar.id}"]
  
  # ディスクの修正関連のパラメータ 
  hostname = "myserver"
  password = "p@ssw0rd"
  ssh_key_ids = ["100000000000", "200000000000"]
  disable_pw_auth = true
  note_ids = ["100000000000", "200000000000"]

}
```

