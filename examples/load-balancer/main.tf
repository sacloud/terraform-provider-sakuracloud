### 概要
# ロードバランサー利用のサンプルテンプレート
#
# このテンプレートはWebサーバを2台構成してLBにて負荷分散する構成です。
# ホスト名やrootパスワードの設定などは"構成構築"タブ内をご覧ください。
# ※SandboxではVPCルータの設定変更がされないなど正常に完了しません。
#
# -----------------------------------------------
# 変数定義(パスワード、ホスト名など)
# -----------------------------------------------
# サーバパスワード
variable "server_password" {
  default = "your-password"
}

# サーバ1/サーバ2のホスト名
variable "server_hostname" {
  type    = list(string)
  default = ["your-hostname1", "your-hostname2"]
}

# ロードバランサに割り当てるIPアドレスの定義
locals {
  # ロードバランサ自体のIPアドレス
  load_balancer_ip_addresses = [sakuracloud_internet.router.ip_addresses[0]]

  # ロードバランサのVIP
  load_balancer_vip = sakuracloud_internet.router.ip_addresses[1]

  # ソーリーサーバのIPアドレス
  load_balancer_sorry_server_ip_address = sakuracloud_internet.router.ip_addresses[2]

  # スイッチ+ルータのグローバルIPアドレスから実サーバに割り当てるIPアドレスの開始インデックス
  load_balancer_used_ip_address_count = length(local.load_balancer_ip_addresses) + 2 # IPアドレス+VIP+ソーリーサーバ
}

### 構成構築
# -----------------------------------------------
# スイッチ+ルータ
# -----------------------------------------------
resource "sakuracloud_internet" "router" {
  name = "wan_switch"
}

# -----------------------------------------------
# ロードバランサ
# -----------------------------------------------
resource "sakuracloud_load_balancer" "lb" {
  name = "load_balancer"
  plan = "standard"

  network_interface {
    switch_id    = sakuracloud_internet.router.switch_id
    vrid         = 1
    ip_addresses = local.load_balancer_ip_addresses
    netmask      = sakuracloud_internet.router.netmask
    gateway      = sakuracloud_internet.router.gateway
  }

  vip {
    vip          = local.load_balancer_vip
    port         = 80
    delay_loop   = 10
    sorry_server = local.load_balancer_sorry_server_ip_address

    dynamic server {
      for_each = var.server_hostname
      content {
        ip_address = sakuracloud_internet.router.ip_addresses[server.key + local.load_balancer_used_ip_address_count]
        protocol   = "http"
        path       = "/"
        status     = 200
      }
    }
  }
}

# ----------------------------------------------------------
# スタートアップスクリプト(DSR構成のためにループバックアドレス設定)
# パブリックスクリプト"lb-dsr"を参照
# ----------------------------------------------------------
resource "sakuracloud_note" "lb_dsr" {
  name    = "lb_dsr"
  content = <<EOF
PARA1="${local.load_balancer_vip}"
PARA2="net.ipv4.conf.all.arp_ignore = 1"
PARA3="net.ipv4.conf.all.arp_announce = 2"
PARA4="DEVICE=lo:0"
PARA5="IPADDR="$PARA1
PARA6="NETMASK=255.255.255.255"

VERSION=$(rpm -q centos-release --qf %%{VERSION}) || exit 1

case "$VERSION" in
  6 ) ;;
  7 ) firewall-cmd --add-service=http --zone=public --permanent
      firewall-cmd --reload;;
  * ) ;;
esac

cp --backup /etc/sysctl.conf /tmp/ || exit 1

echo $PARA2 >> /etc/sysctl.conf
echo $PARA3 >> /etc/sysctl.conf
sysctl -p 1>/dev/null

cp --backup /etc/sysconfig/network-scripts/ifcfg-lo:0 /tmp/ 2>/dev/null

touch /etc/sysconfig/network-scripts/ifcfg-lo:0
echo $PARA4 > /etc/sysconfig/network-scripts/ifcfg-lo:0
echo $PARA5 >> /etc/sysconfig/network-scripts/ifcfg-lo:0
echo $PARA6 >> /etc/sysconfig/network-scripts/ifcfg-lo:0

ifup lo:0 || exit 1

exit 0
EOF

}

# ----------------------------------------------------------
# サーバーへのWebサーバー(httpd)インストール
# ----------------------------------------------------------
resource "sakuracloud_note" "install_httpd" {
  name    = "install_httpd"
  content = <<EOF
yum install -y httpd || exit 1
echo 'This is a TestPage!!' >> /var/www/html/index.html || exit1
systemctl enable httpd.service || exit 1
systemctl start httpd.service || exit 1
firewall-cmd --add-service=http --zone=public --permanent || exit 1

exit 0
EOF

}

# ----------------------------------------------------------
# サーバーで利用するパブリックアーカイブ(CentOS7)
# ----------------------------------------------------------
data "sakuracloud_archive" "centos" {
  os_type = "centos7" // "ubuntu" を指定するとUbuntuの最新安定版パブリックアーカイブ
}

# ----------------------------------------------------------
# サーバー
# ----------------------------------------------------------
resource "sakuracloud_disk" "disks" {
  name              = "disk${format("%02d", count.index)}" // ディスク名の指定
  plan              = "ssd"                               // プランの指定
  size              = 40                                  // 容量指定(GB)
  source_archive_id = data.sakuracloud_archive.centos.id  // アーカイブの設定
  count             = length(var.server_hostname)
}

resource "sakuracloud_server" "server" {
  name   = "server${format("%02d", count.index)}"
  core   = 2
  memory = 2
  disks  = [sakuracloud_disk.disks[count.index].id]
  network_interface {
    upstream = sakuracloud_internet.router.switch_id
  }

  disk_edit_parameter {
    ip_address = sakuracloud_internet.router.ip_addresses[count.index + local.load_balancer_used_ip_address_count]
    gateway    = sakuracloud_internet.router.gateway
    netmask    = sakuracloud_internet.router.netmask
    hostname   = var.server_hostname[count.index]
    password   = var.server_password
    note_ids   = [sakuracloud_note.lb_dsr.id, sakuracloud_note.install_httpd.id]
  }

  count = length(var.server_hostname)
}

