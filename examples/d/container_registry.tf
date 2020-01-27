data "sakuracloud_container_registry" "foobar" {
  filter {
    names = ["foobar"]
  }
}