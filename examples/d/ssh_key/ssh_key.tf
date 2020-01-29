data "sakuracloud_ssh_key" "foobar" {
  filter {
    names = ["foobar"]
  }
}