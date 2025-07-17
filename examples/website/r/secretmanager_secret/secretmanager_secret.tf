resource "sakuracloud_secretmanager_secret" "foobar" {
  name     = "foobar"
  value    = "secret value!"
  vault_id = "secretmanager-resource-id" # e.g. sakuracloud_secretmanager.foobar.id
}