resource "sakuracloud_note" "foobar" {
  name    = "foobar"
  content = file("startup-script.sh")
}