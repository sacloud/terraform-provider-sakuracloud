# 公開鍵(生成)(sakuracloud_ssh_key_gen)

---

**全ゾーン共通のグローバルリソースです。**

SSH公開鍵をさくらのクラウド側で生成するためのリソースです。

- このリソースは`import`非対応です。
- このリソースは更新動作に非対応です。属性を変更した場合は常に再生成されます。

### 設定例

```hcl
resource "sakuracloud_ssh_key_gen" "key" {
  name = "foobar"

  # パスフレーズ(オプション、指定する場合は8〜64文字)
  # pass_phrase = "your_pass_phrase"

  description = "Description"
}
```

### パラメーター

|パラメーター         |必須  |名称                |初期値     |設定値                    |補足                                          |
|-------------------|:---:|--------------------|:--------:|------------------------|----------------------------------------------|
| `name`            | ◯   | 公開鍵名           | -        | 文字列                  | 64文字まで|
| `pass_phrase`     | -   | パスフレーズ           | -        | 文字列                  | 空文字、または8〜64文字まで|
| `description`     | -   | 説明  | - | 文字列 | 512文字まで |

### 属性

|属性名                | 名称                    | 補足                                        |
|---------------------|------------------------|--------------------------------------------|
| `id`                | 公開鍵ID                | -                                          |
| `private_key`       | 秘密鍵                  | 生成された秘密鍵                              |
| `public_key`        | 公開鍵                  | -                                       |
| `fingerprint`       | フィンガープリント        | -                                          |

### 利用例

#### サーバのプロビジョニング時、SSH接続用に公開鍵(生成)を利用する

```hcl
#SSH公開鍵
resource sakuracloud_ssh_key_gen "key" {
  name = "foobar"

  provisioner "local-exec" {
    command = "echo \"${self.private_key}\" > id_rsa; chmod 0600 id_rsa"
  }

  provisioner "local-exec" {
    when    = "destroy"
    command = "rm -f id_rsa"
  }
}

#OS(CentOS)
data sakuracloud_archive "centos" {
  os_type = "centos"
}

#ディスクの定義
resource sakuracloud_disk "foobar" {
  name              = "foobar"
  source_archive_id = "${data.sakuracloud_archive.centos.id}"
  password          = "PUT_YOUR_PASSWORD_HERE"

  # 生成した公開鍵のIDを指定
  ssh_key_ids = ["${sakuracloud_ssh_key_gen.key.id}"]

  # SSH接続時のパスワード/チャレンジレスポンス認証を無効化
  disable_pw_auth = true
}

#サーバの定義
resource sakuracloud_server "foobar" {
  name  = "foobar"
  disks = ["${sakuracloud_disk.foobar.id}"]

  # プロビジョニング
  connection {
    user        = "root"
    host        = "${self.ipaddress}"
    private_key = "${sakuracloud_ssh_key_gen.key.private_key}"
  }

  provisioner "remote-exec" {
    inline = [
      "hostname",
    ] # 実行したいコマンドを指定
  }
}

#SSH接続用のアウトプット定義
output "ssh_to_server" {
  value = "ssh -i id_rsa root@${sakuracloud_server.foobar.ipaddress}"
}
```