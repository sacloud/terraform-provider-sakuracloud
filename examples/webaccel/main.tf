resource "sakuracloud_webaccel" "foobar" {
  name             = "hoge"
  domain_type      = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type     = "web"
    origin   = "docs.usacloud.jp"
    protocol = "https"
  }
  origin_guard_token {}
  onetime_url_secrets = [
    "abc-0x123456"
  ]
  vary_support      = true
  default_cache_ttl = 3600
  normalize_ae      = "gzip"
}


resource "sakuracloud_webaccel_acl" "foobar_acl" {
  acl = join("\n", [
    "allow 192.0.1.1/32",
    "allow 192.0.2.5/25",
    "deny all",
  ])
  site_id    = sakuracloud_webaccel.foobar.id
  depends_on = [sakuracloud_webaccel.foobar]
}

resource "sakuracloud_webaccel_activation" "foobar_status" {
  site_id    = sakuracloud_webaccel.foobar.id
  enabled    = true
  depends_on = [sakuracloud_webaccel_acl.foobar_acl]
}
