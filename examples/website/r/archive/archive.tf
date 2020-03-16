# from archive/disk
resource "sakuracloud_archive" "from-archive-or-disk" {
  name         = "foobar"
  description  = "description"
  tags         = ["tag1", "tag2"]

  source_archive_id   = 123456789012
  source_archive_zone = "tk1a"
  # source_disk_id    = 123456789012
}

# from shared archive
resource "sakuracloud_archive" "from-shared-archive" {
  name         = "foobar"
  description  = "description"
  tags         = ["tag1", "tag2"]

  source_shared_key = "is1a:123456789012:xxx"
}


# from local file
resource "sakuracloud_archive" "foobar" {
  name         = "foobar"
  description  = "description"
  tags         = ["tag1", "tag2"]
  size         = 20
  archive_file = "test/dummy.raw"
}