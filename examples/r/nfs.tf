resource "sakuracloud_nfs" "foobar" {
  name        = "foobar"
  switch_id   = sakuracloud_switch.foobar.id
  plan        = "ssd"
  size        = "500"
  ip_address  = "192.168.11.101"
  netmask     = 24
  gateway     = "192.168.11.1"
  description = "description"
  tags        = ["tag1", "tag2"]
}

resource "sakuracloud_switch" "foobar" {
  name = "foobar"
}