resource "sakuracloud_secretmanager" "foobar" {
  name        = "foobar"
  description = "description"
  tags        = ["tag1", "tag2"]
  kms_key_id  = "kms-resource-id" # e.g. sakuracloud_kms.foobar.id
}
