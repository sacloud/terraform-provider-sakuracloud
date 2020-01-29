resource "sakuracloud_mobile_gateway" "foobar" {
  name                = "foobar"
  internet_connection = true
  dns_servers         = data.sakuracloud_zone.zone.dns_servers

  private_network_interface {
    switch_id  = sakuracloud_switch.foobar.id
    ip_address = "192.168.11.101"
    netmask    = 24
  }

  description = "description"
  tags        = ["tag1", "tag2"]

  traffic_control {
    quota                = 256
    band_width_limit     = 64
    enable_email         = true
    enable_slack         = true
    slack_webhook        = "https://hooks.slack.com/services/xxx/xxx/xxx"
    auto_traffic_shaping = true
  }

  static_route {
    prefix   = "192.168.10.0/24"
    next_hop = "192.168.11.1"
  }
  static_route {
    prefix   = "192.168.10.0/25"
    next_hop = "192.168.11.2"
  }
  static_route {
    prefix   = "192.168.10.0/26"
    next_hop = "192.168.11.3"
  }
}

data sakuracloud_zone "zone" {}

resource "sakuracloud_switch" "foobar" {
  name = "foobar"
}