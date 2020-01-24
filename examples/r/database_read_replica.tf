resource "sakuracloud_database_read_replica" "foobar" {
  master_id   = data.sakuracloud_database.master.id
  network_interface {
    ip_address  = "192.168.11.111"
  }
  name        = "foobar"
  description = "description"
  tags        = ["tag1", "tag2"]
}

data sakuracloud_database "master" {
  filter {
    names = ["master-database-name"]
  }
}