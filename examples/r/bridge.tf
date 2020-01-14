resource "sakuracloud_switch" "is1a" {
  name        = "is1a"
  description = "description"
  bridge_id   = sakuracloud_bridge.foobar.id
  zone        = "is1a"
}

resource "sakuracloud_switch" "is1b" {
  name        = "is1b"
  description = "description"
  bridge_id   = sakuracloud_bridge.foobar.id
  zone        = "is1b"
}

resource "sakuracloud_bridge" "foobar" {
  name        = "foobar"
  description = "description"
}