data "sakuracloud_proxylb" "foobar" {
  filter {
    names = ["foobar"]
  }
}