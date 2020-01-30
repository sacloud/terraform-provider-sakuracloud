resource "sakuracloud_gslb" "foobar" {
  name = "example"

  health_check {
    protocol    = "http"
    delay_loop  = 10
    host_header = "example.com"
    path        = "/"
    status      = "200"
  }

  sorry_server = "192.2.0.1"

  server {
    ip_address = "192.2.0.11"
    weight     = 1
    enabled    = true
  }
  server {
    ip_address = "192.2.0.12"
    weight     = 1
    enabled    = true
  }

  description = "description"
  tags        = ["tag1", "tag2"]
}