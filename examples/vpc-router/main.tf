### 概要
#
# VPCルータ利用のサンプルテンプレート
#
# VPCルータを設置し、サーバに対してHTTPとHTTPSのみインターネットからの接続を許可します。
# そのほかの通信に関してはL2TP/IPSecにて接続することにより可能となります。
#
# <構築手順>
#   1) tffile編集画面の"変数定義"タブにて以下の値を編集します。
#      - L2TP/IPSec 事前共有キー(vpc_router.l2tp.pre_shared_secret)
#      - L2TP/IPSec ユーザー名(vpc_router.l2tp.username)
#      - L2TP/IPSec パスワード(vpc_router.l2tp.password)
#      - サーバ管理者のパスワード(server.password)
#
#      ※ 上記以外の変数については必要に応じて編集してください。
#
#   2) リソースマネージャー画面にて"計画/反映"を実行します。
#
# リソースの各設定などは各リソースのタブ内をご覧ください。
# ※SandboxではVPCルータの設定変更がされないなど正常に完了しません。
#
### 変数定義(VPCルータ)
variable vpc_router {
  default = {
    name = "vpc-router" # VPCルータ名

    l2tp = {
      pre_shared_secret = "put-your-secret"     # < 変更してください
      username          = "put-your-name"     # < 変更してください
      password          = "put-your-password" # < 変更してください
      range_start       = "192.168.100.251"   # クライアントに割り当てる開始IPアドレス
      range_stop        = "192.168.100.254"   # クライアントに割り当てる終了IPアドレス
    }

    # VPCルータのプライベート側NICのIPアドレス/マスク長の設定
    eth1 = {
      ip_address = "192.168.100.1"
      netmask    = 24
    }

    # サーバへのReverse NAT(port_forwarding)対象
    reverse_nat = [
      {
        // for HTTP
        protocol = "tcp"
        port     = 80
      },
      {
        // for HTTPS
        protocol = "tcp"
        port     = 443
      },
    ]
  }
  type = object({
    name = string,
    l2tp = object({
      pre_shared_secret = string,
      username          = string,
      password          = string,
      range_start       = string,
      range_stop        = string,
    }),
    eth1 = object({
      ip_address = string,
      netmask    = number,
    }),
    reverse_nat = list(object({
      protocol = string,
      port     = number,
    })),
  })
}

### 変数定義(サーバ)
variable server {
  default = {
    name       = "server"                 # サーバ名
    hostname   = "server"                 # ホスト名
    password   = "put-your-root-password" # < 変更してください
    ip_address = "192.168.100.101"        # サーバ IPアドレス
    netmask    = 24                       # サーバ マスク長
    core       = 2                        # サーバ コア数
    memory     = 2                        # サーバ メモリサイズ(GB単位)
    diskname   = "disk"                   # ディスク名
  }

  type = object({
    name       = string,
    hostname   = string,
    password   = string,
    ip_address = string,
    netmask    = number,
    core       = number,
    memory     = number,
    diskname   = string,
  })
}

### 変数定義(スイッチ)
variable switch {
  default = {
    name = "local-sw" # スイッチ名
  }

  type = object({ name = string })
}

### VPCルータ
# VPCルータ
resource "sakuracloud_vpc_router" "vpc" {
  name = var.vpc_router.name # VPCルータ名
  private_network_interface {
    index        = 1
    switch_id    = sakuracloud_switch.sw.id
    ip_addresses = [var.vpc_router.eth1.ip_address] # VPCルータIPアドレスの設定
    netmask      = var.vpc_router.eth1.netmask      # ネットワークマスク
  }
  l2tp {
    pre_shared_secret = var.vpc_router.l2tp.pre_shared_secret
    range_start       = var.vpc_router.l2tp.range_start # IPアドレス動的割り当て範囲(開始)
    range_stop        = var.vpc_router.l2tp.range_stop  # IPアドレス動的割り当て範囲(終了)
  }

  user {
    name     = var.vpc_router.l2tp.username
    password = var.vpc_router.l2tp.password
  }

  dynamic port_forwarding {
    for_each = var.vpc_router.reverse_nat
    content {
      protocol     = port_forwarding.value.protocol
      public_port  = port_forwarding.value.port
      private_ip   = var.server.ip_address
      private_port = port_forwarding.value.port
    }
  }
}

### スイッチ
resource "sakuracloud_switch" "sw" {
  name = var.switch.name
}

### サーバ/ディスク

# コピー元アーカイブ(CentOS7)
data "sakuracloud_archive" "centos" {
  os_type = "centos7"
}

# ディスク
resource "sakuracloud_disk" "disk" {
  name              = var.server.diskname
  source_archive_id = data.sakuracloud_archive.centos.id
}

# サーバ
resource "sakuracloud_server" "server" {
  name   = var.server.name
  core   = var.server.core
  memory = var.server.memory
  disks  = [sakuracloud_disk.disk.id]

  network_interface {
    upstream = sakuracloud_switch.sw.id
  }

  disk_edit_parameter {
    ip_address = var.server.ip_address
    netmask    = var.server.netmask
    gateway    = var.vpc_router.eth1.ip_address
    hostname   = var.server.hostname
    password   = var.server.password
  }
}