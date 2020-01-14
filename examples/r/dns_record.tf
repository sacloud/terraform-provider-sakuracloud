resource "sakuracloud_dns" "foobar" {
  zone = "example.com"
}

resource "sakuracloud_dns_record" "record1" {
  dns_id = sakuracloud_dns.foobar.id
  name   = "www"
  type   = "A"
  value  = "192.168.0.1"
}

resource "sakuracloud_dns_record" "record2" {
  dns_id = sakuracloud_dns.foobar.id
  name   = "www"
  type   = "A"
  value  = "192.168.0.2"
}