resource "sakuracloud_webaccel" "foobar" {
  name             = "hoge"
  domain_type      = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type     = "web"
    origin   = "docs.usacloud.jp"
    protocol = "https"
  }
  onetime_url_secrets = [
    "abc-0x123456"
  ]
  vary_support      = true
  default_cache_ttl = 3600
  normalize_ae      = "gzip"
}

