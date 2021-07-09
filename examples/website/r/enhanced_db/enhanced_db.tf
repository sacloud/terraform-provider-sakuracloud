resource "sakuracloud_enhanced_db" "foobar" {
  name            = "example"
  database_name   = "example"
  password        = "your-password"

  description = "..."
  tags        = ["...", "..."]
}