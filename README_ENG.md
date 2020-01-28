# Terraform for Sakura Cloud

![Test Status](https://github.com/sacloud/terraform-provider-sakuracloud/workflows/Tests/badge.svg)
[![Slack](https://slack.usacloud.jp/badge.svg)](https://slack.usacloud.jp/)

It is a plugin for operating `Sakura Cloud` from Terraform.

This plugin is developed by Sakura Cloud user communication as Sakura Internet Inc. official recognition tool.

([Japanese version](README.md))

## Quick start

We will build infrastructure on Sakura Cloud with following configuration.

- Use the latest stable version of CentOS.
- Disk: SSD/20GB, Server: 1core/1GB memory
  - It is omitted from the definition file because it is the default value.
- Disable password / Challenge response authentication when connecting to server via SSH. (Allow only public key authentication.)
- The public key for SSH connection is generated on `Sakura Cloud`.
  - Save the created key on the local machine.

Install `Terraform` and `Terraform for Sakura Cloud` on the local machine with reference to [Installation](https://docs.usacloud.jp/terraform-v1/installation/).

After installation, infrastructure building will be done by executing the following commands.

```bash
#################################################
# Set the API key of "Sakura Cloud" as an environment variable.
#################################################
export SAKURACLOUD_ACCESS_TOKEN=[Sakura Cloud API Token]
export SAKURACLOUD_ACCESS_TOKEN_SECRET=[Sakura Cloud API Secret]

#################################################
# Create Terraform definition file.
#################################################
mkdir work; cd work
tee sakura.tf <<-'EOF'
# Define server administrator password
variable "password" {
  default = "PUT_YOUR_PASSWORD_HERE"
}

# Set the target zone
provider sakuracloud {
  zone = "tk1a" # Tokyo No.1 Zone
}

# Create the public key on Sakura Cloud
resource "sakuracloud_ssh_key_gen" "key" {
  name = "foobar"
}

# Store the private key to local machine
resource "local_file" "private_key" {
  content  = "${sakuracloud_ssh_key_gen.key.private_key}"
  filename = "id_rsa"
}

# Define the data resource for reference to the ID of the public archive (OS)
data "sakuracloud_archive" "centos" {
  os_type = "centos"
}

# Define Disk
resource "sakuracloud_disk" "disk01" {
  name              = "disk01"
  source_archive_id = "${data.sakuracloud_archive.centos.id}"
}

# Define Server
resource "sakuracloud_server" "server01" {
  name  = "server01"
  disks = ["${sakuracloud_disk.disk01.id}"]

  ssh_key_ids     = ["${sakuracloud_ssh_key_gen.key.id}"]
  password        = "${var.password}"
  disable_pw_auth = true
}

# Define output to display SSH commands
output "ssh_to_server" {
  value = "ssh -i id_rsa root@${sakuracloud_server.server01.ipaddress}"
}
EOF

#################################################
# Build infrastructure ( init & apply )
#################################################
terraform init
terraform apply
```

## Document

- [Terraform for Sakura Cloud Documents](https://docs.usacloud.jp/terraform-v1/)

#### Unsupported Resources

The following resources are unsupported because API is not provided by Sakura Cloud.

- Local Router(terraform-provider-sakuracloud v2 supports this)
- Resources Manager
- Web Accelerator
- Object Strage (Create Bucket)
- License
- Discount Passport
- Coupon

## Building/Developing

### Build

  ```bash
  make build
  ```

### Build (Cross compiling)

  ```bash
  make build-x
  ```

### Build (Build on Docker)

  ```bash
  make docker-build
  ```

### Test

  ```bash
  make test
  ```

### Acceptance Test (Test with real Sakura Cloud API call)

  ```bash
  make testacc
  ```

#### Preview terraform.io style documents

You can preview terraform.io style documents in english.  
Run following command and open `http://localhost:4567/docs/providers/sakuracloud`.  

```bash
# preview terraform.io style docs
make website
```

## License

 `terraform-proivder-sakuracloud` Copyright (C) 2016-2020 terraform-provider-sakuracloud authors.

  This project is published under [Apache 2.0 License](LICENSE.txt).
