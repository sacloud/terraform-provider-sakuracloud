resource "sakuracloud_sim" "foobar" {
  name        = "foobar"
  description = "description"
  tags        = ["tag1", "tag2"]

  iccid    = "your-iccid"
  passcode = "your-password"
  #imei     = "your-imei"
  carrier = ["softbank", "docomo", "kddi"]

  enabled = true
}