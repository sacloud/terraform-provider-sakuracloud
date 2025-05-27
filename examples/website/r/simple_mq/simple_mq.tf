resource "sakuracloud_simple_mq" "foobar" {
  name        = "foobar"
  description = "description"
  tags        = ["tag1", "tag2"]

  visibility_timeout_seconds = 30
  expire_seconds             = 345600
}
