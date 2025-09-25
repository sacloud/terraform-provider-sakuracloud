data "sakuracloud_webaccel" "site" {
  name = "your-site-name"
  # or
  # domain = "your-domain"
}

resource "sakuracloud_webaccel_activation" "site_status" {
  site_id = data.sakuracloud_webaccel.site.id
  enabled = true
}
