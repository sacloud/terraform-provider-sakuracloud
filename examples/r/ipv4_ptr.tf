resource sakuracloud_server "server" {
  name = "foobar"
  network_interface {
    upstream = "shared"
  }
}

resource "sakuracloud_ipv4_ptr" "foobar" {
  ip_address     = sakuracloud_server.server.ip_address
  hostname       = "www.example.com"
  retry_max      = 30
  retry_interval = 10
}