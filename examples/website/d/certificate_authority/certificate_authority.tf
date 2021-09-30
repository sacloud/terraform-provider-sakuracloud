data "sakuracloud_certificate_authority" "foobar" {
  filter {
    names = ["foobar"]
  }
}