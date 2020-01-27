resource "sakuracloud_internet" "foobar" {
  name = "foobar"

  netmask     = 28
  band_width  = 100
  enable_ipv6 = false

  description = "description"
  tags        = ["tag1", "tag2"]
}