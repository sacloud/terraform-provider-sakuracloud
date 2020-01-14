data "sakuracloud_cdrom" "foobar" {
  filter {
    condition {
      name   = "Name"
      values = ["Parted Magic 2013_08_01"]
    }
  }
}