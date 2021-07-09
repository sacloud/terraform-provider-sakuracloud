resource "sakuracloud_proxylb" "foobar" {
  name           = "foobar"
  plan           = 100
  vip_failover   = true
  sticky_session = true
  gzip           = true
  timeout        = 10
  region         = "is1"

  health_check {
    protocol    = "http"
    delay_loop  = 10
    host_header = "example.com"
    path        = "/"
  }

  sorry_server {
    ip_address = "192.0.2.1"
    port       = 80
  }

  syslog {
    server = "192.0.2.1"
    port   = 514
  }

  bind_port {
    proxy_mode = "http"
    port       = 80
    response_header {
      header = "Cache-Control"
      value  = "public, max-age=10"
    }
  }

  server {
    ip_address = sakuracloud_server.foobar.ip_address
    port       = 80
    group      = "group1"
  }

  rule {
    action = "forward"
    host   = "www.example.com"
    path   = "/"
    group  = "group1"
  }
  rule {
    action               = "redirect"
    host                 = "www2.example.com"
    path                 = "/"
    group                = "group1"
    redirect_location    = "https://redirect.example.com"
    redirect_status_code = "301"
  }
  rule {
    action               = "fixed"
    host                 = "www3.example.com"
    path                 = "/"
    group                = "group1"
    fixed_status_code    = "200"
    fixed_content_type   = "text/plain"
    fixed_message_body   = "body"
  }

  description = "description"
  tags        = ["tag1", "tag2"]
}

resource sakuracloud_server "foobar" {
  name = "foobar"
  network_interface {
    upstream = "shared"
  }
}