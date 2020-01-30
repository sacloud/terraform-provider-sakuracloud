resource "sakuracloud_archive" "foobar" {
  name         = "foobar"
  description  = "description"
  tags         = ["tag1", "tag2"]
  size         = 20
  archive_file = "test/dummy.raw"
}