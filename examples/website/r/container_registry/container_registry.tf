variable users {
  type = list(object({
    name       = string
    password   = string
    permission = string
  }))
  default = [
    {
      name       = "user1"
      password   = "password1"
      permission = "all"
    },
    {
      name       = "user2"
      password   = "password2"
      permission = "readwrite"
    }
  ]
}

resource "sakuracloud_container_registry" "foobar" {
  name            = "foobar"
  subdomain_label = "your-subdomain-label"
  access_level    = "readwrite" # this must be one of ["readwrite"/"readonly"/"none"]

  description = "description"
  tags        = ["tag1", "tag2"]

  dynamic user {
    for_each = var.users
    content {
      name       = user.value.name
      password   = user.value.password
      permission = user.value.permission
    }
  }
}