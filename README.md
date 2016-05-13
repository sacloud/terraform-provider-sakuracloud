# terraform-provider-sakuracloud

Terraform provider for SakuraCloud. - `Terraform for さくらのクラウド`

## Installation

1. Download the plugin from the [releases page](https://github.com/yamamoto-febc/terraform-provider-sakuracloud/releases/latest)
2. Put it in the same directory as the terraform binary. ex:`$GOPATH/bin/`.


## Usage

  - [Provider Configuration](#provider-configuration)
  - Resource Configuration
    - [sakuracloud_server](#resource-configuration-sakuracloud_server)
    - [sakuracloud_disk](#resource-configuration-sakuracloud_disk)
    - [sakuracloud_ssh_key](#resource-configuration-sakuracloud_ssh_key)
    - [sakuracloud_dns](#resource-configuration-sakuracloud_dns)
    - [sakuracloud_gslb](#resource-configuration-sakuracloud_gslb)
    - [sakuracloud_simple_monitor](#resource-configuration-sakuracloud_simple_monitor)
  - [Samples](#samples)

## Provider Configuration

### Example

```
provider "sakuracloud" {
    token = "your API token"
    secret = "your API secret"
    zone = "target zone"
}
```

### Argument Reference

The following arguments are supported:

* `token` - (Required) This is the SakuraCloud API token. This can also be specified
  with the `SAKURACLOUD_ACCESS_TOKEN` shell environment variable.

* `secret` - (Required) This is the SakuraCloud API secret. This can also be specified
  with the `SAKURACLOUD_ACCESS_TOKEN_SECRET` shell environment variable.

* `zone` - (Required) This is the SakuraCloud zone. This can also be specified
  with the `SAKURACLOUD_ZONE` shell environment variable.

* `trace` - (Optional) Flag of trace mode. This can also be specified
  with the `SAKURACLOUD_TRACE_MODE` shell environment variable.

##  Resource Configuration `sakuracloud_server`

Provides a SakuraCloud Server resource. This can be used to create, modify,
and delete Server.

### Example Usage

```
# Create a new Server"
resource "sakuracloud_server" "myserver" {
    name = "myserver"
    disks = ["${sakuracloud_disk.mydisk.id}"]
    switched_interfaces = [""]
    description = "Server from TerraForm for SAKURA CLOUD"
    tags = ["@virtio-net-pci"]
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the server.
* `disks` - (Required) The ID list of the disk to connect server.
* `core` - (Optional) The number of CPU core. default `1`.
* `memory` - (Optional) The size of memory(GB). default `1`.
* `shared_interface` - (Optional) The flag of to create a NIC to connect to a shared segment.
* `switched_interfaces` - (Optional) The ID list of to create a NIC to connect to switch.
   If `""` is specified , it creates a NIC with empty connection.
* `description` - (Optional) The description of the server.
* `tags` - (Optional) The tags of the server.
* `zone` - (Optional) The zone of to create server.

### Attributes Reference

The following attributes are exported:

* `id` - The ID of the server.
* `name` - The name of the server.
* `disks`- The ID list of the disks.
* `core` - The number of the CPU core.
* `memory` - The size(MB) of the memory.
* `shared_interface` - The flag of has NIC to connect to a shared segment.
* `switched_interfaces` - The ID list of the connected switch.
* `description` - The description of the server.
* `tags` - The tags of the server.
* `zone` - The zone of the server.
* `mac_addresses` - The MAC address list of the server.
* `shared_nw_ipaddress` - The IP address that are connected to the shared segment.
* `shared_nw_dns_servers` - The IP address list of server's region on.
* `shared_nw_gateway` - The IP address of default route.
* `shared_nw_address` - The network address of the shared segment.
* `shared_nw_mask_len` - The length of network mask of the shared segment.


## Resource Configuration `sakuracloud_disk`

Provides a SakuraCloud Disk resource. This can be used to create, modify,
and delete Disk .

### Example Usage

```
# Create a new disk with source archive named "Ubuntu Server 14.04"
resource "sakuracloud_disk" "mydisk"{
    name = "mydisk"
    size = 20480
    source_archive_name = "Ubuntu Server 14.04"
    description = "Disk from terraform for SAKURA CLOUD"
    tags = ["hoge1" , "hoge2"]
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the disk.
* `plan` - (Optional) The plan of the disk.
* `connection` - (Optional) The connection of the disk. default is `4`(SSD).
* `size` - (Optional) The size(GB) of the disk. default is `20`.
* `source_archive_id` - (Optional) The ID of source archive.
  Conflicts with `source_archive_name , source_disk_id , source_disk_name`
* `source_archive_name` - (Optional) The name of source archive.
  Conflicts with `source_archive_id , source_disk_id , source_disk_name`
* `source_disk_id` - (Optional) The ID of source disk.
  Conflicts with `source_archive_id , source_archive_name , source_disk_name`
* `source_disk_name` - (Optional) The name of source disk.
  Conflicts with `source_archive_id , source_archive_name , source_disk_id`
* `description` - (Optional) The description of the disk.
* `tags` - (Optional) The tags of the disk.
* `zone` - (Optional) The zone of to create disk.
* `password` - (Optional) The password of the disk.
* `ssh_key_ids` - (Optional) The ID list of SSHKey.
* `disable_pw_auth` - (Optional) The flag that to disable SSH login with password authentication / challenge-response. default id `false`.
* `note_ids` - (Optional) The ID list of Note.

### Attributes Reference

The following attributes are exported:

* `id` - The ID of the disk.
* `name`- The name of the disk.
* `plan` - The plan of the disk.
* `connection` - The connection of the disk.
* `size` - The size(GB) of the disk.
* `source_archive_id` - The ID of source archive.
* `source_archive_name` - The name of source archive.
* `source_disk_id` - The ID of source disk.
* `source_disk_name` - The name of source disk.
* `description` - The description of the disk.
* `tags` - The tags of the disk.
* `zone` - The zone of the disk.
* `password` - The password of the disk.
* `ssh_key_ids` The ID list of SSHKey.
* `disable_pw_auth` - The flag that to disable SSH login with password authentication / challenge-response.
* `note_ids` - The ID list of Note.


## Resource Configuration `sakuracloud_ssh_key`

Provides a SakuraCloud SSHKey resource. This can be used to create, modify,
and delete SSHKey.

### Example Usage

```
resource "sakuracloud_ssh_key" "mykey" {
    name = "mykey"
    public_key = "ssh-rsa XXXXXXXXX....."
    # or
    #public_key = "${file("./id_rsa.pub")}"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SSHKey.
* `public_key` - (Required) The value of the SSHKey.
* `description` - (Optional) The description of the SSHKey.

### Attributes Reference

The following attributes are exported:

* `id` - The ID of the SSHKey.
* `name`- The name of the SSHKey.
* `public_key` - The value of the SSHKey.
* `description` - The description of the SSHKey.
* `fingerprint` - The FingerPrint of the SSHKey.

## Resource Configuration `sakuracloud_dns`

Provides a SakuraCloud DNS resource. This can be used to create,
modify, and delete DNS records.

### Example Usage

```
# Create a new DNS zone and add two A records.
resource "sakuracloud_dns" "dns" {
    zone = "example.com"
    records = {
        name = "test1"
        type = "A"
        value = "192.168.0.1"
    }
    records = {
        name = "test2"
        type = "A"
        value = "192.168.0.2"
    }
}
```

### Argument Reference

The following arguments are supported:

* `zone` - (Required) The DNS target zone name.
* `description` - (Required) The description of DNS.
* `tags` - (Required) The tags of DNS.
* `records` - (Optional) The records of target zone.
  * `name` - (Required) The name of the record.
  * `type` - (Required) The type of the record.
  * `value` - (Required) The value of the record.
  * `ttl` - (Optional) The TTL of the record . default `3600`.
  * `priority` - (Optional) (Only type is MX) The priority of the record.



### Attributes Reference

The following attributes are exported:

* `id` - The ID of the DNS.
* `zone`- The DNS target zone name.
* `dns_servers` - The name servers of the target zone.
* `description` - The description of target zone.
* `tags` - The description of target zone.
* `records` - The records of target zone.
  * `name` - The name of the record.
  * `type` - The type of the record.
  * `value` - The value of the record.
  * `ttl` - The TTL of the record.
  * `priority` - (Only type is MX) The priority of the record.


## Resource Configuration `sakuracloud_gslb`

Provides a SakuraCloud GSLB(Global Site Load Balancing) resource. This can be used to create,
modify, and delete GSLB.

### Example Usage

```
# Create a new GSLB and add two target server.
resource "sakuracloud_gslb" "mygslb" {
    name = "gslb_from_terraform"
    health_check = {
        protocol = "http"
        delay_loop = 10
        host_header = "example.com"
        path = "/"
        status = "200"
    }
    description = "GSLB from terraform for SAKURA CLOUD"
    tags = ["hoge1" , "hoge2" ]
    servers = {
      ipaddress = "192.0.2.1"
    }
    servers = {
      ipaddress = "192.0.2.2"
    }

}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of GSLB.
* `health_check` - (Required) The health_check rule of GSLB.
  * `protocol` - (Required) The protocol to use for health check. Must be in [`http`,`https`,`tcp`,`ping`]
  * `dalay_loop` - (Optional) The delay_loop of health check. Must be between `10` and `60`. default is `10`
  * `host_header` - (Only protocol is `http` or `https`) The host_header to use for health check.
  * `path` - (Only when protocol is `http` or `https`) The request path to use for health check.
  * `status` - (Only when protocol is `http` or `https`) The response code of health check request.
  * `port` - (Only when protocol is `tcp`) The port number to use for health check.
* `weighted` - (Optional)The flag of enabling to weighted balancing. default `false`
* `description` - (Optional) The description of GSLB.
* `tags` - (Required) The tags of GSLB.
* `servers` (Optional) The target servers of GSLB.
  * `ipaddress` - The IPAddress of target server.
  * `enabled` - The flag of enabling to target server.
  * `weight` - (Only when `weighted` is true)The weight of target server.


### Attributes Reference

The following attributes are exported:

* `id` - The ID of the GSLB.
* `name`- The name of GSLB.
* `health_check` - The health_check rule of GSLB.
  * `protocol` - The protocol to use for health check.
  * `dalay_loop` - The delay_loop of health check.
  * `host_header` - The host_header to use for health check.
  * `path` - The request path to use for health check.
  * `status` - The response code of health check request.
  * `port` - The port number to use for health check.
* `weighted` - The flag of enabling to weighted balancing.
* `description` - The description of GSLB.
* `tags` - The tags of GSLB.
* `servers` - The target servers of GSLB.
  * `ipaddress` - The IPAddress of target server.
  * `enabled` - The flag of enabling to target server.
  * `weight` - The weight of target server.
* `FQDN` - The FQDN of GSLB.


## Resource Configuration `sakuracloud_simple_monitor`

Provides a SakuraCloud SimpleMonitor resource. This can be used to create,
modify, and delete SimpleMonitor records.

### Example Usage

```
# Create a new Simple Monitor
resource "sakuracloud_simple_monitor" "mymonitor" {
    target = "${sakuracloud_server.myserver.shared_nw_ipaddress}"
    health_check = {
        protocol = "http"
        delay_loop = 60
        path = "/"
        status = "200"
    }
    notify_email_enabled = true
    notify_slack_enabled = true
    notify_slack_webhook = "https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
}
```

### Argument Reference

The following arguments are supported:

* `target` - (Required) The monitor target IP or domain name.
* `health_check` - (Required) The health_check rule of SimpleMonitor.
  * `protocol` - (Required) The protocol to use for health check. Must be in [`http`,`https`,`tcp`,`ping`,`ssh`,`dns`]
  * `dalay_loop` - (Optional) The delay_loop of health check. Must be between `60` and `3600`. default is `60`
  * `path` - (Only when protocol is `http` or `https`) The request path to use for health check.
  * `status` - (Only when protocol is `http` or `https`) The response code of health check request.
  * `port` - (Only when protocol is `tcp` or `ssh`) The port number to use for health check.
  * `qname` - (Only when protocol is `dns`) The port number to use for health check.
  * `expected_data` - (Only when protocol is `dns`) The port number to use for health check.
* `description` - (Optional) The descripton of SimpleMonitor.
* `tags` - (Optional) The tags of SimpleMonitor.
* `notify_email_enabled` - (Optional) The flag of enabled email.
* `notify_slack_enabled` - (Optional) The flag of enabled slack.
* `notify_slack_webhook` - (Optional) The URL of slack webhook.
* `enabled` - (Optional) The flag of enabled SimpleMonitor.



### Attributes Reference

The following attributes are exported:

* `id` - The ID of the DNS
* `target` - The monitor target IP or domain name.
* `health_check` - The health_check rule of SimpleMonitor.
  * `protocol` - The protocol to use for health check. Must be in [`http`,`https`,`tcp`,`ping`,`ssh`,`dns`]
  * `dalay_loop` - The delay_loop of health check. Must be between `60` and `3600`. default is `60`
  * `path` - (Only when protocol is `http` or `https`) The request path to use for health check.
  * `status` - (Only when protocol is `http` or `https`) The response code of health check request.
  * `port` - (Only when protocol is `tcp` or `ssh`) The port number to use for health check.
  * `qname` - (Only when protocol is `dns`) The port number to use for health check.
  * `expected_data` - (Only when protocol is `dns`) The port number to use for health check.
* `description` - The descripton of SimpleMonitor.
* `tags` - The tags of SimpleMonitor.
* `notify_email_enabled` - The flag of enabled email.
* `notify_slack_enabled` - The flag of enabled slack.
* `notify_slack_webhook` - The URL of slack webhook.
* `enabled` - The flag of enabled SimpleMonitor.


## Samples

```sample.tf
#**************************************************************************************
# TerraForm for さくらのクラウド
#**************************************************************************************
# tfファイルのサンプルです。
#
# あらかじめ環境変数`SAKURACVLOUD_ACCESS_TOKEN`と
# `SAKURACLOUD_ACCESS_TOKEN_SECRET`を設定しておきます。
# アップロードするSSH鍵は`ssh-keygen`などで作成しておきます。
#
# `terraform apply`すると以下の内容でさくらのクラウド上にプロビジョニングが行われます。
#
# 1) 手元のSSH公開鍵をアップロード
# 2) ディスク作成(Ubuntu 14.04をソースアーカイブとしたもの)
# 3) サーバー作成(パスワード認証無効化状態)
# 4) サーバーに割り振られたグローバルIPへのPING監視
# 5) サーバーに割り振られたグローバルIPでDNS(Aレコード)登録
#**************************************************************************************

provider "sakuracloud" {
    zone = "is1a"
}

/************************
 Server
************************/
resource "sakuracloud_server" "myserver" {
    name = "myserver"
    disks = ["${sakuracloud_disk.mydisk.id}"]
    description = "Server from TerraForm for SAKURA CLOUD"
    tags = ["@virtio-net-pci"]
}

/************************
 Disk
************************/
resource "sakuracloud_disk" "mydisk" {
    name = "mydisk"
    source_archive_name = "Ubuntu Server 14.04.4 LTS 64bit"
    description = "Disk from TerraForm for SAKURA CLOUD"
    ssh_key_ids = ["${sakuracloud_ssh_key.mykey.id}"]
    disable_pw_auth = true
}

/************************
 SSHKey
************************/
resource "sakuracloud_ssh_key" "mykey" {
    name = "key"
    public_key = "${file("./id_rsa.pub")}"
}

/************************
 SimpleMonitor
************************/
resource "sakuracloud_simple_monitor" "mymonitor" {
    target = "${sakuracloud_server.myserver.shared_nw_ipaddress}"
    health_check = {
        protocol = "ping"
    }
    description = "SimpleMonitor from terraform for SAKURA CLOUD"
    notify_email_enabled = true
    notify_slack_enabled = true
    notify_slack_webhook = "https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
}

/************************
 DNS
************************/
resource "sakuracloud_dns" "foobar" {
    zone = "example.com"
    description = "DNS from terraform for SAKURA CLOUD"
    records = {
        name = "terraform-sample"
        type = "A"
        value = "${sakuracloud_server.myserver.shared_nw_ipaddress}"
    }
}
```

## Building/Developing

  `godep get $(go list ./... | grep -v vendor)`

  `godep restore`

  `godep go test .`

  `TF_ACC=1 godep go test -v -timeout=60m .` run acceptance tests. (requires ENV vars)

  `godep go build -o path/to/desired/terraform-provider-sakuracloud bin/terraform-provider-sakuracloud/main.go`


## License

  This project is published under [Apache 2.0 License](LICENSE).

## Author

  * Kazumichi Yamamoto ([@yamamoto-febc](https://github.com/yamamoto-febc))
