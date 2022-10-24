variable username {}
variable password {}
variable replica_password {}

resource "sakuracloud_database" "foobar" {
  database_type    = "mariadb"
  # database_version = "13" // optional
  plan             = "30g"

  username = var.username
  password = var.password

  replica_password = var.replica_password

  network_interface {
    switch_id     = sakuracloud_switch.foobar.id
    ip_address    = "192.168.11.11"
    netmask       = 24
    gateway       = "192.168.11.1"
    port          = 3306
    source_ranges = ["192.168.11.0/24", "192.168.12.0/24"]
  }

  backup {
    time     = "00:00"
    weekdays = ["mon", "tue"]
  }

  parameters = {
    max_connections = 100
  }

  name        = "foobar"
  description = "description"
  tags        = ["tag1", "tag2"]
}

resource "sakuracloud_switch" "foobar" {
  name = "foobar"
}