data sakuracloud_webaccel "site" {
  name = "your-site-name"
  # or
  # domain = "your-domain"
}

resource sakuracloud_webaccel_acl "acl" {
  site_id = data.sakuracloud_webaccel.site.id

  acl = join("\n", [
    "deny 192.0.2.5/25",
    "deny 198.51.100.0",
    "allow all",
  ])
}