resource "sakuracloud_kms" "foobar" {
  name        = "foobar"
  description = "description"
  tags        = ["tag1", "tag2"]
  # key_origin  = "imported" # Optional, default is "generated"
  # plain_key   = "AfL5zzjD4RgeFQm3vvAADwPNrurNUc616877wsa8v4w="
}
