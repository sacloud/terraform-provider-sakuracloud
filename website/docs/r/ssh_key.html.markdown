---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_ssh_key"
sidebar_current: "docs-sakuracloud-resource-misc-ssh-key"
description: |-
  Provides a SakuraCloud SSH Key resource. This can be used to create, update, and delete SSH Keys.
---

# sakuracloud\_ssh_key

Provides a SakuraCloud SSH Key resource. This can be used to create, update, and delete SSH Keys.

## Example Usage

```hcl
# Create a new SSH Key(from file)
resource "sakuracloud_ssh_key" "foobar" {
  name       = "foobar"
  public_key = file("~/.ssh/id_rsa")
}

# Create a new SSH Key(from tls_private_key resource)
resource "tls_private_key" "foobar" {
  algorithm = "RSA"
}

resource "sakuracloud_ssh_key" "foobar2" {
  name       = "foobar2"
  public_key = tls_private_key.foobar.public_key_openssh
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `public_key` - (Required) The body of the public key. 
* `description` - (Optional) The description of the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `description` - The description of the resource.
* `public_key` - The body of the public key. 
* `finger_print` - The fingerprint of the public key.

## Import

SSH Keys can be imported using the SSH Key ID.

```
$ terraform import sakuracloud_ssh_key.foobar <ssh_key_id>
```
