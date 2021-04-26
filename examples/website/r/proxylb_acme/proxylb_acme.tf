resource sakuracloud_proxylb_acme "foobar" {
  proxylb_id        = sakuracloud_proxylb.foobar.id
  accept_tos        = true
  common_name       = "www.example.com"
  subject_alt_names = ["www1.example.com"]
  update_delay_sec = 120
}

data "sakuracloud_proxylb" "foobar" {
  filter {
    names = ["foobar"]
  }
}