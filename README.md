# Terraform Provider for SakuraCloud

![Test Status](https://github.com/sacloud/terraform-provider-sakuracloud/workflows/Tests/badge.svg)
[![Slack](https://slack.usacloud.jp/badge.svg)](https://slack.usacloud.jp/)

- Terraform Website: https://terraform.io
- Documentation: https://docs.usacloud.jp/terraform
- Slack Workspace for Users: https://slack.usacloud.jp

## Usage Example

```hcl
# Configure the SakuraCloud Provider
provider "sakuracloud" {
  # We recommend pinning to the specific version of the SakuraCloud Provider you're using
  # since new versions are released frequently
  version = "=2.0.1"

  # More information on the authentication methods supported by
  # the SakuraCloud Provider can be found here:
  # https://docs.usacloud.jp/terraform/provider/

  # profile = "..."
}

variable password {}

data "sakuracloud_archive" "ubuntu" {
  os_type = "ubuntu"
}

resource "sakuracloud_disk" "example" {
  name              = "example"
  source_archive_id = data.sakuracloud_archive.ubuntu.id

  # If you want to prevent re-creation of the disk
  # when archive id is changed, please uncomment this.
  # lifecycle {
  #   ignore_changes = [
  #     source_archive_id,
  #   ]
  # }
}

resource "sakuracloud_server" "example" {
  name        = "example"
  disks       = [sakuracloud_disk.example.id]
  core        = 1
  memory      = 2
  description = "description"
  tags        = ["app=web", "stage=staging"]

  network_interface {
    upstream = "shared"
  }

  disk_edit_parameter {
    hostname        = "example"
    password        = var.password
  }
}
```

## Requirements

- [Terraform](https://terraform.io) v0.12+

## License

 `terraform-proivder-sakuracloud` Copyright (C) 2016-2020 terraform-provider-sakuraclou authors.
 
  This project is published under [Apache 2.0 License](LICENSE).
