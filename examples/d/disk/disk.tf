data "sakuracloud_disk" "foobar" {
  filter {
    names = ["foobar"]
  }
}