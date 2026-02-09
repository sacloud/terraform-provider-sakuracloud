# Terraform Provider for SakuraCloud

![Test Status](https://github.com/sacloud/terraform-provider-sakuracloud/workflows/Tests/badge.svg)
[![Discord](https://img.shields.io/badge/Discord-SAKURA%20Users-blue)](https://discord.gg/yUEDN8hbMf)

- Terraform Website: https://terraform.io
- Terraform Registry: https://registry.terraform.io/providers/sacloud/sakuracloud/latest
- Documentation: https://registry.terraform.io/providers/sacloud/sakuracloud/latest/docs
- Documentation(ja): https://docs.usacloud.jp/terraform
- Discord: https://discord.gg/yUEDN8hbMf

> [!IMPORTANT]
> さくらのクラウド向けTerraformプロバイダーはv3が最新となります。  
> v3: https://github.com/sacloud/terraform-provider-sakura
> 
> v2はメンテナンスモードに移行しており、新規の機能追加や大きな変更はまずv3で行われます。  
> これからさくらのクラウド向けTerraformプロバイダーを利用される場合はv3をご利用ください。  
> 
> 詳細は以下のクラウドニュースを参照してください。  
> https://cloud.sakura.ad.jp/news/2025/12/25/terraform-provider-sakura-v3/

バージョン情報:

| バージョン                                                      | ステータス | 説明                  |
|------------------------------------------------------------|-----------|---------------------|
| [v3](https://github.com/sacloud/terraform-provider-sakura) | 現行バージョン | 最新の機能追加・変更はここで行われます |
| v2                                                         | メンテナンスモード | 積極的な機能開発は行いません      |
| v1                                                         | EOL | サポート終了              |

## Usage Example

```hcl
# Configure the SakuraCloud Provider
terraform {
  required_providers {
    sakuracloud = {
      source = "sacloud/sakuracloud"

      # We recommend pinning to the specific version of the SakuraCloud Provider you're using
      # since new versions are released frequently
      version = "2.34.2"
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

 `terraform-provider-sakuracloud` Copyright (C) 2016-2023 terraform-provider-sakuracloud authors.

  This project is published under [Apache 2.0 License](LICENSE).
