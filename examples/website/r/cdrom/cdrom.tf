resource "sakuracloud_cdrom" "foobar" {
  name           = "foobar"
  size           = 5
  iso_image_file = "example.iso"
  description    = "description"
  tags           = ["tag1", "tag2"]
}