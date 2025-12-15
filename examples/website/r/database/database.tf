variable username {}
variable password {}
variable replica_password {}

resource "sakuracloud_database" "foobar" {
  database_type    = "mariadb"
  database_version = "10.11" // optional
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

  # continuous_backupを指定するときはdatabase_versionが必須
  # continuous_backup {
  #   days_of_week = ["mon", "tue"]
  #   time         = "01:30"
  #   connect      = "nfs://${sakuracloud_nfs.foobar.network_interface[0].ip_address}/export"
  # }

  parameters = {
    max_connections = 100
  }

  monitoring_suite {
    enabled = true
  }

  disk {
    encryption_algorithm = "aes256_xts"
    kms_key_id           = sakuracloud_kms.foobar.id
  }

  name        = "foobar"
  description = "description"
  tags        = ["tag1", "tag2"]
}

resource "sakuracloud_nfs" "foobar" {
  name = "foobar"
  plan = "ssd"
  size = "100"

  network_interface {
    switch_id   = sakuracloud_switch.foobar.id
    ip_address  = "192.168.11.111"
    netmask     = 24
    gateway     = "192.168.11.1"
  }
}

resource "sakuracloud_switch" "foobar" {
  name = "foobar"
}

resource "sakuracloud_kms" "foobar" {
  name = "foobar"
}