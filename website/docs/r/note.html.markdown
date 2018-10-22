---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_note"
sidebar_current: "docs-sakuracloud-resource-misc-note"
description: |-
  Provides a SakuraCloud Note (Startup-Script) resource. This can be used to create, update, and delete Notes.
---

# sakuracloud\_note

Provides a SakuraCloud Note (Startup-Script) resource. This can be used to create, update, and delete Notes.

## Example Usage

```hcl
# Create a new Note
resource "sakuracloud_note" "foobar" {
  name  = "foobar"
  class = "shell"

  content     = <<-EOS
  #!/bin/sh

  : your-script-here
  EOS
  
  # for RancherOS(cloud_config)
  # class   = "yaml_cloud_config"
  # content = file("path/to/your/content")

  description = "description"
  tags        = ["foo", "bar"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `class` - (Required) The content body of the Note.  
Valid value is one of the following: [ "shell" (default) / "yaml_cloud_config" ]
* `content` - (Required) The content body of the Note.
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `class` - The name of the note class.
* `content` - The body of the note. 
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.

## Import

Notes can be imported using the Note ID.

```
$ terraform import sakuracloud_note.foobar <note_id>
```
