resource "sakuracloud_server" "foobar" {
  name        = "foobar"
  disks       = [sakuracloud_disk.foobar.id]
  description = "description"
  tags        = ["tag1", "tag2"]

  network_interface {
    upstream         = "shared"
    packet_filter_id = data.sakuracloud_packet_filter.foobar.id
  }

  disk_edit_parameter {
    hostname = "hostname"
    password = "password"
    disable_pw_auth = true

    # ssh_key_ids     = ["<ID>", "<ID>"]
    # note {
    #  id         = "<ID>"
    #  api_key_id = "<ID>"
    #  variables = {
    #    foo = "bar"
    #  }
    # }
  }
}

data "sakuracloud_packet_filter" "foobar" {
  filter {
    names = ["foobar"]
  }
}

data "sakuracloud_archive" "ubuntu" {
  os_type = "ubuntu2004"
}

resource "sakuracloud_disk" "foobar" {
  name              = "{{ .arg0 }}"
  source_archive_id = data.sakuracloud_archive.ubuntu.id
}