resource "sakuracloud_secret_manager_secret" "foobar" {
  name     = "foobar"
  value    = "secret value!"
  vault_id = "secret_manager-resource-id" # e.g. sakuracloud_secret_manager.foobar.id
}