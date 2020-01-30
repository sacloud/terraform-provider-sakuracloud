### 概要
# サーバ作成及びシンプル監視のサンプルテンプレート
# このテンプレートはサーバを1台作成してシンプル監視を使用してpingにて
# 死活監視を行います。
# ホスト名やrootパスワードの設定などは"変数定義"タブ内をご覧ください。
# ※Sandboxではシンプル監視のリソースは作成されません。

# -----------------------------------------------
# 変数定義(パスワード、ホスト名など)
# -----------------------------------------------
# サーバ
variable "servers" {
  type = list(object({
    name     = string,
    core     = number,
    memory   = number,
    password = string,
    hostname = string,
    diskname = string,
    disksize = number,
  }))
  default = [
    {
      name     = "server01"      // サーバ名の設定
      core     = 2               // CPUコア数の指定
      memory   = 2               // メモリ容量の指定(GB)
      password = "your-password" // rootパスワードの設定
      hostname = "your-hostname" // ホスト名の設定
      diskname = "disk01"        // ディスク名の設定
      disksize = 40              // ディスク容量指定(GB)
    },
    # 複数台作成したい場合は以下のように指定
    # {
    #   name     = "server02"      // サーバ名の設定
    #   core     = 2               // CPUコア数の指定
    #   memory   = 2               // メモリ容量の指定(GB)
    #   password = "your-password" // rootパスワードの設定
    #   hostname = "your-hostname" // ホスト名の設定
    #   diskname = "disk02"        // ディスク名の設定
    #   disksize = 40              // ディスク容量指定(GB)
    # },
  ]
}

### 構成構築
# ----------------------------------------------------------
# サーバーで利用するパブリックアーカイブ(CentOS7)
# ----------------------------------------------------------
data "sakuracloud_archive" "centos" {
  os_type = "centos7" // "ubuntu" を指定するとUbuntuの最新安定版パブリックアーカイブ
}

# ----------------------------------------------------------
# サーバー
# ----------------------------------------------------------
# ディスク作成
resource "sakuracloud_disk" "disks" {
  name              = var.servers[count.index].diskname
  plan              = "ssd"
  size              = var.servers[count.index].disksize
  source_archive_id = data.sakuracloud_archive.centos.id // アーカイブの設定

  count = length(var.servers)
}

# VM作成
resource "sakuracloud_server" "servers" {
  name   = var.servers[count.index].name
  disks  = [sakuracloud_disk.disks[count.index].id]
  core   = var.servers[count.index].core
  memory = var.servers[count.index].memory

  network_interface {
    upstream = "shared"
  }

  disk_edit_parameter {
    hostname = var.servers[count.index].hostname
    password = var.servers[count.index].password
  }

  count = length(var.servers)
}

# ----------------------------------------------------------
# シンプル監視
# ----------------------------------------------------------
# ping監視の例
resource "sakuracloud_simple_monitor" "mymonitor" {
  target = sakuracloud_server.servers[count.index].ip_address

  health_check {
    protocol = "ping"
  }

  notify_email_enabled = true
  enabled              = true

  count = length(var.servers)
}