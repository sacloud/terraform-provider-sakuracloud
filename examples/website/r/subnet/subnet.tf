resource sakuracloud_internet "foobar" {
  name = "foobar"
}

resource "sakuracloud_subnet" "foobar" {
  internet_id = sakuracloud_internet.foobar.id
  next_hop    = sakuracloud_internet.foobar.min_ip_address
}