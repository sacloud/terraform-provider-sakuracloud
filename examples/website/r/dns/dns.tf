resource "sakuracloud_dns" "foobar" {
  zone        = "example.com"
  description = "description"
  tags        = ["tag1", "tag2"]
  record {
    name  = "www"
    type  = "A"
    value = "192.168.11.1"
  }
  record {
    name  = "www"
    type  = "A"
    value = "192.168.11.2"
  }
  monitoring_suite {
    enabled = true
  }
}