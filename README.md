# Terraform Provider for SakuraCloud

![Test Status](https://github.com/sacloud/terraform-provider-sakuracloud/workflows/Tests/badge.svg)
[![Discord](https://img.shields.io/badge/Discord-SAKURA%20Users-blue)](https://discord.gg/yUEDN8hbMf)

- Terraform Website: https://terraform.io
- Terraform Registry: https://registry.terraform.io/providers/sacloud/sakuracloud/latest
- Documentation: https://registry.terraform.io/providers/sacloud/sakuracloud/latest/docs
- Documentation(ja): https://docs.usacloud.jp/terraform
- Discord: https://discord.gg/yUEDN8hbMf

## Usage Example

```hcl
# Configure the SakuraCloud Provider
terraform {
  required_providers {
    sakuracloud = {
      source = "sacloud/sakuracloud"

      # We recommend pinning to the specific version of the SakuraCloud Provider you're using
      # since new versions are released frequently
      version = "2.25.2"
      #version = "~> 2"
    }
  }
}
provider "sakuracloud" {
  # More information on the authentication methods supported by
  # the SakuraCloud Provider can be found here:
  # https://docs.usacloud.jp/terraform/provider/

  # profile = "..."
}

variable password {}

data "sakuracloud_archive" "ubuntu" {
  os_type = "ubuntu2004"
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

 `terraform-proivder-sakuracloud` Copyright (C) 2016-2023 terraform-provider-sakuraclou authors.
 
  This project is published under [Apache 2.0 License](LICENSE).
