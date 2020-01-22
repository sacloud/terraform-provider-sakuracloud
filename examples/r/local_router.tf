resource "sakuracloud_local_router" "foobar" {
  name        = "example"
  description = "descriptio"
  tags        = ["tag1", "tag2"]

  switch {
    code     = sakuracloud_switch.foobar.id
    category = "cloud"
    zone_id  = "is1a"
  }

  network_interface {
    vip          = "192.168.11.1"
    ip_addresses = ["192.168.11.11", "192.168.11.12"]
    netmask      = 24
    vrid         = 101
  }

  static_route {
    prefix   = "10.0.0.0/24"
    next_hop = "192.168.11.2"
  }
  static_route {
    prefix   = "172.16.0.0/16"
    next_hop = "192.168.11.3"
  }

  peer {
    peer_id     = data.sakuracloud_local_router.peer.id
    secret_key  = data.sakuracloud_local_router.secret_keys[0]
    description = "description"
  }
}

resource "sakuracloud_switch" "foobar" {
  name = "example"
}

data "sakuracloud_local_router" "peer" {
  filter {
    names = ["peer"]
  }
}
