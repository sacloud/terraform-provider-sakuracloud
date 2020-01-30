data "sakuracloud_simple_monitor" "foobar" {
  filter {
    names = ["foobar"]
  }
}