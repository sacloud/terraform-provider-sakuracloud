resource "sakuracloud_disk" "foobar" {
  name = "foobar"
}
resource "sakuracloud_auto_backup" "foobar" {
  name           = "foobar"
  disk_id        = sakuracloud_disk.foobar.id
  weekdays       = ["mon", "tue", "wed", "thu", "fri", "sat", "sun"]
  max_backup_num = 5
  description    = "description"
  tags           = ["tag1", "tag2"]
}