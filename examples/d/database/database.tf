data "sakuracloud_database" "foobar" {
  filter {
    names = ["foobar"]
  }
}