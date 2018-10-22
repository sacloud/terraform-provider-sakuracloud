# サーバ コネクタ(sakuracloud_server_connector)

---

サーバとディスク/パケットフィルタ/ISOイメージなどの他のリソースとの間の接続のみを扱うリソース。

他のリソースがサーバの値を利用する場合などに利用する。  
(例: パケットフィルタにてサーバのIPアドレスを利用する場合など)


### 設定例

```hcl
resource "sakuracloud_server" "sv" {
  name = "sv"
}

resource "sakuracloud_packet_filter" "pf" {
  name = "pf"

  expressions {
    protocol  = "ip"
    source_nw = sakuracloud_server.sv.ipaddress
  }
}

# リソース間の接続のみを扱うリソース = コネクタリソース
resource "sakuracloud_server_connector" "connector" {
  server_id = sakuracloud_server.sv.id

  # パケットフィルタ
  packet_filter_ids = [sakuracloud_packet_filter.pf.id]
  
  # ディスク
  #disks             = [sakuracloud_disk.disk01.id]
  
  # ISOイメージ(CD-ROM)
  #cdrom_id          = data.sakuracloud_cdrom.centos.id
}
```

### パラメーター

|パラメーター          |必須  |名称                |初期値     |設定値 |補足                                          |
|--------------------|:---:|--------------------|:--------:|------|----------------------------------------------|
| `server_id`        | ◯ | サーバID          | -   | 文字列 | - |
| `disks`            | - | ディスクID          | -   | リスト(文字列) | サーバに接続するディスクのID |
| `packet_filter_ids`| - | パケットフィルタID | - | リスト(文字列) | NICに適用するパケットフィルタのIDをリストで指定する。リストの先頭からeth0,eth1の順で適用される |
| `cdrom_id`         | - | CDROM(ISOイメージ)ID | - | 文字列 | - |
| `graceful_shutdown_timeout` | - | シャットダウンまでの待ち時間 | - | 数値(秒数) | シャットダウンが必要な場合の通常シャットダウンするまでの待ち時間(指定の時間まで待ってもシャットダウンしない場合は強制シャットダウンされる) |
| `zone`             | - | ゾーン | - | `is1a`<br />`is1b`<br />`tk1a`<br />`tk1v` | - |


### 属性

|属性名                    | 名称                     | 補足                                        |
|-------------------------|-------------------------|--------------------------------------------|
| `id`                    | ID                      | -                                          |
