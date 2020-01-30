data "sakuracloud_server" "foobar" {
  filter {
    names = ["foobar"]
  }
}