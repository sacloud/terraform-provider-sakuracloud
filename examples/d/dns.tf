data "sakuracloud_dns" "foobar" {
  filter {
    names = ["foobar"]
  }
}