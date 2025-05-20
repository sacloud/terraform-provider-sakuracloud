data "sakuracloud_archive" "ubuntu" {
  os_type = "ubuntu"
}

resource "sakuracloud_disk" "foobar" {
  name              = "foobar"
  plan              = "ssd"
  connector         = "virtio"
  size              = 20
  source_archive_id = data.sakuracloud_archive.ubuntu.id
  #distant_from      = ["111111111111"]

  description = "description"
  tags        = ["tag1", "tag2"]
}