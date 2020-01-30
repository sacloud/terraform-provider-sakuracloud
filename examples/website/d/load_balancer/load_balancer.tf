data "sakuracloud_load_balancer" "foobar" {
  filter {
    names = ["foobar"]
  }
}