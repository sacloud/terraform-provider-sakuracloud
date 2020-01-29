resource "sakuracloud_bucket_object" "foobar" {
  bucket  = "foobar"
  key     = "example.txt"
  content = file("example.txt")
}