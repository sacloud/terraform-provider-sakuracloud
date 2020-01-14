resource "sakuracloud_packet_filter" "foobar" {
  name        = "foobar"
  description = "description"
}

resource "sakuracloud_packet_filter_rules" "rules" {
  packet_filter_id = sakuracloud_packet_filter.foobar.id

  expression {
    protocol  = "tcp"
    dest_port = "22"
  }

  expression {
    protocol  = "tcp"
    dest_port = "80"
  }

  expression {
    protocol  = "tcp"
    dest_port = "443"
  }

  expression {
    protocol = "icmp"
  }

  expression {
    protocol = "fragment"
  }

  expression {
    protocol    = "udp"
    source_port = "123"
  }

  expression {
    protocol  = "tcp"
    dest_port = "32768-61000"
  }

  expression {
    protocol  = "udp"
    dest_port = "32768-61000"
  }

  expression {
    protocol    = "ip"
    allow       = false
    description = "Deny ALL"
  }
}