resource "sakuracloud_vpc_router" "standard" {
  name                = "standard"
  description         = "description"
  tags                = ["tag1", "tag2"]
  internet_connection = true
}

resource "sakuracloud_vpc_router" "premium" {
  name        = "premium"
  description = "description"
  tags        = ["tag1", "tag2"]
  plan        = "premium"

  internet_connection = true

  switch_id    = sakuracloud_internet.foobar.switch_id
  vip          = sakuracloud_internet.foobar.ip_addresses[0]
  ip_addresses = [sakuracloud_internet.foobar.ip_addresses[1], sakuracloud_internet.foobar.ip_addresses[2]]
  aliases      = [sakuracloud_internet.foobar.ip_addresses[3]]
  vrid         = 1

  network_interface {
    index        = 1
    switch_id    = sakuracloud_switch.foobar.id
    vip          = "192.168.11.1"
    ip_addresses = ["192.168.11.2", "192.168.11.3"]
    netmask      = 24
  }

  dhcp_server {
    interface_index = 1

    range_start = "192.168.11.11"
    range_stop  = "192.168.11.20"
    dns_servers = ["8.8.8.8", "8.8.4.4"]
  }

  dhcp_static_mapping {
    ip_address  = "192.168.11.10"
    mac_address = "aa:bb:cc:aa:bb:cc"
  }

  firewall {
    interface_index = 1

    direction = "send"
    expression {
      protocol            = "tcp"
      source_network      = ""
      source_port         = "80"
      destination_network = ""
      destination_port    = ""
      allow               = true
      logging             = true
      description         = "desc"
    }

    expression {
      protocol            = "ip"
      source_network      = ""
      source_port         = ""
      destination_network = ""
      destination_port    = ""
      allow               = false
      logging             = true
      description         = "desc"
    }
  }

  l2tp {
    pre_shared_secret = "example"
    range_start       = "192.168.11.21"
    range_stop        = "192.168.11.30"
  }

  port_forwarding {
    protocol     = "udp"
    public_port  = 10022
    private_ip   = "192.168.11.11"
    private_port = 22
    description  = "desc"
  }

  pptp {
    range_start = "192.168.11.31"
    range_stop  = "192.168.11.40"
  }

  site_to_site_vpn {
    peer              = "192.0.2.1"
    remote_id         = "192.0.2.1"
    pre_shared_secret = "example"
    routes            = ["192.0.2.0/8"]
    local_prefix      = ["192.168.21.0/24"]
  }

  static_nat {
    public_ip   = sakuracloud_internet.foobar.ip_addresses[3]
    private_ip  = "192.168.11.12"
    description = "desc"
  }

  static_route {
    prefix   = "172.16.0.0/16"
    next_hop = "192.168.11.99"
  }

  user {
    name     = "username"
    password = "password"
  }
}

resource "sakuracloud_internet" "foobar" {
  name = "foobar"
}

resource sakuracloud_switch "foobar" {
  name = "foobar"
}