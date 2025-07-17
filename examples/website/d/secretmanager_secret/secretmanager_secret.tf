data "sakuracloud_secretmanager_secret" "foobar" {
  name     = "foobar"
  vault_id = "secretmanager-vault-id"
}