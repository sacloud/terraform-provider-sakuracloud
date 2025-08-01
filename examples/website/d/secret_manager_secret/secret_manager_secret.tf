data "sakuracloud_secret_manager_secret" "foobar" {
  name     = "foobar"
  vault_id = "secret_manager-resource-id"
}