data "sakuracloud_private_host" "foobar" {
  filter {
    names = ["foobar"]
  }
}