data "sakuracloud_vpc_router" "foobar" {
  filter {
    names = ["foobar"]
  }
}