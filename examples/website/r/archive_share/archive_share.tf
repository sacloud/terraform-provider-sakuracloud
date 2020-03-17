resource "sakuracloud_archive" "source" {
  name         = "foobar"
  size         = 20
  archive_file = "test/dummy.raw"
}

resource "sakuracloud_archive_share" "share_info" {
  archive_id = sakuracloud_archive.source.id
}