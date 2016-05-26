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
    - [sakuracloud_switch](#resource-configuration-sakuracloud_switch)
    - [sakuracloud_internet](#resource-configuration-sakuracloud_internet)
    - [sakuracloud_packet_filter](#resource-configuration-sakuracloud_packet_filter)
    - [sakuracloud_bridge](#resource-configuration-sakuracloud_bridge)
    - [sakuracloud_note](#resource-configuration-sakuracloud_note)
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

* `zone` - (Optional) This is the SakuraCloud zone. This can also be specified
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
    additional_interfaces = [""]
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
* `base_interface` - (Optional) The ID of to create a NIC to connect to a shared segment.
   When `shared`,it connect to the shared segment.
   When `switch_id` , it connect to the switch+router.
   When `""` , it creates a NIC with empty connection.
* `additional_interfaces` - (Optional) The ID list of to create a NIC to connect to switch.
   If `""` is specified , it creates a NIC with empty connection.
* `packet_filter_ids` - (Optional) The ID list of the packet filter.
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
* `base_interface` - The ID of NIC to connect to a shared segment.
* `additional_interfaces` - The ID list of the connected switch.
* `description` - The description of the server.
* `tags` - The tags of the server.
* `zone` - The zone of the server.
* `packet_filter_ids` - The ID list of the packet filter.
* `mac_addresses` - The MAC address list of the server.
* `base_nw_ipaddress` - The IP address that are connected to the shared segment.
* `base_nw_dns_servers` - The IP address list of server's region on.
* `base_nw_gateway` - The IP address of default route.
* `base_nw_address` - The network address of the shared segment.
* `base_nw_mask_len` - The length of network mask of the shared segment.


## Resource Configuration `sakuracloud_disk`

Provides a SakuraCloud Disk resource. This can be used to create, modify,
and delete Disk .

### Example Usage

```
# Create a new disk with source archive named "Ubuntu Server 14.04"
resource "sakuracloud_disk" "mydisk"{
    name = "mydisk"
    size = 20
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

## Resource Configuration `sakuracloud_switch`

Provides a SakuraCloud Switch resource. This can be used to create, modify,
and delete Switch.

### Example Usage

```
resource "sakuracloud_switch" "myswitch" {
    name = "mykey"
    description = "Switch from terraform for SAKURA CLOUD"
    tags = ["hoge1" , "hoge2"]
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the switch.
* `description` - (Optional) The description of the switch.
* `tags` - (Optional) The tags of the switch.
* `bridge_id` - (Optional) The ID of the bridge.
* `zone` - (Optional) The zone of the switch.

### Attributes Reference

The following attributes are exported:

* `id` - The ID of the switch.
* `name`- The name of the switch.
* `description` - The description of the switch.
* `bridge_id` - The ID of the bridge.
* `tags` - The tags of the switch.
* `zone` - The zone of the switch.
* `server_ids` - The ID list of connected server.

## Resource Configuration `sakuracloud_internet`

Provides a SakuraCloud Intenet(router) resource. This can be used to create, modify,
and delete Internet(router).

### Example Usage

```
resource "sakuracloud_internet" "myrouter" {
    name = "myrouter"
    description = "Switch from terraform for SAKURA CLOUD"
    tags = ["hoge1" , "hoge2"]
    nw_mask_len = 28
    band_width = 100
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the internet.
* `description` - (Optional) The description of the internet.
* `tags` - (Optional) The tags of the internet.
* `zone` - (Optional) The zone of the internet.
* `nw_mask_len` - (Optional) The length of network address maks.
* `band_width` - (Optional) The bandwitch of the internet.


### Attributes Reference

The following attributes are exported:

* `id` - The ID of the internet.
* `name`- The name of the internet.
* `description` - The description of the internet.
* `tags` - The tags of the internet.
* `zone` - The zone of the internet.
* `switch_id` - The ID of connected switch.
* `server_ids` - The ID list of connected servers.
* `nw_address` - The ipaddress of network.
* `nw_gateway` - The ipaddress of gateway.
* `nw_min_ipaddress` - The min ipaddress of alocated to the internet.
* `nw_max_ipaddress` - The max ipaddress of alocated to the internet.
* `nw_ipaddresses` - The ipaddress list of alocated to the internet.

## Resource Configuration `sakuracloud_packet_filter`

Provides a SakuraCloud PacketFilter resource. This can be used to create, modify,
and delete PacketFilter.

### Example Usage

```
resource "sakuracloud_packet_filter" "myfilter" {
    name = "myfilter"
    description = "PacketFilter from terraform for SAKURA CLOUD"
    expressions = {
        protocol = "tcp"
        source_nw = "192.168.2.0/24"
        source_port = "0-65535"
        dest_port = "80"
        allow = true
    }
    expressions = {
        protocol = "ip"
        source_nw = "0.0.0.0"
        allow = false
        description = "Deny all"
    }
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the packet filter.
* `description` - (Optional) The description of the packet filter.
* `zone` - (Optional) The zone of the packet filter.
* `expressions` - (Required) The expression list of filter.
  * `protocol` - (Required) The protocol of the expression. Following values is allowed [`tcp`,`udp`,`icmp`,`fragment`,`ip`].
  * `source_nw` - (Required) The source network address of the expression.
  * `source_port` - (Required) The source port of the expression.
  * `dest_port` - (Required) The destination port of the expression.
  * `allow` - (Required) The allow flag of athe expression.
  * `description` - (Required) The description of the expression.

### Attributes Reference

The following attributes are exported:

* `id` - The ID of the packet filter.
* `name`- The name of the packet filter.
* `description` - The description of the packet filter.
* `zone` - The zone of the packet filter.
* `expressions` - The expression list of filter.
  * `protocol` - The protocol of the expression.
  * `source_nw` - The source network address of the expression.
  * `source_port` - The source port of the expression.
  * `dest_port` - The destination port of the expression.
  * `allow` - The allow flag of athe expression.
  * `description` - The description of the expression.


## Resource Configuration `sakuracloud_bridge`

Provides a SakuraCloud Bridge resource. This can be used to create, modify,
and delete Bridge.

### Example Usage

```
resource "sakuracloud_bridge" "mybridge" {
    name = "mybridge"
    description = "BRIDGE from terraform for SAKURA CLOUD"
    zone = "is1a"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the bridge.
* `description` - (Optional) The description of the bridge.
* `zone` - (Optional) The zone of the bridge.


### Attributes Reference

The following attributes are exported:

* `id` - The ID of the bridge.
* `name`- The name of the bridge.
* `description` - The description of the bridge.
* `zone` - The zone of the bridge.
* `switch_ids` - The ID list of connected switches.


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

## Resource Configuration `sakuracloud_note`

Provides a SakuraCloud Note resource. This can be used to create, modify,
and delete Note.

### Example Usage

```
resource "sakuracloud_note" "mynote" {
    name = "mynote"
    content = "#!/bin/sh ,,,,"
    # or
    #content = "${file("./example.sh")}"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the note.
* `content` - (Required) The value of the note.
* `description` - (Optional) The description of the note.
* `tags` - (Optional) The tags of the note.

### Attributes Reference

The following attributes are exported:

* `id` - The ID of the note.
* `name`- The name of the note.
* `content` - The content of the note.
* `description` - The description of the note.
* `tags` - The tags of the note.



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
    target = "${sakuracloud_server.myserver.base_nw_ipaddress}"
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
    target = "${sakuracloud_server.myserver.base_nw_ipaddress}"
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
        value = "${sakuracloud_server.myserver.base_nw_ipaddress}"
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
