data "sakuracloud_nfs" "foobar" {
  filter {
    names = ["foobar"]
  }
}