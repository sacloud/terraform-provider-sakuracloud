# Terraform for Sakura Cloud

[![Build Status](https://travis-ci.org/sacloud/terraform-provider-sakuracloud.svg?branch=master)](https://travis-ci.org/sacloud/terraform-provider-sakuracloud)
[![Build status](https://ci.appveyor.com/api/projects/status/paynsb52uauq1jl8?svg=true)](https://ci.appveyor.com/project/sacloud-bot/terraform-provider-sakuracloud)
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

Install `Terraform` and `Terraform for Sakura Cloud` on the local machine with reference to [Installation](https://sacloud.github.io/terraform-provider-sakuracloud/installation/).

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
  content     = sakuracloud_ssh_key_gen.key.private_key
  filename = "id_rsa"
}

# Define the data resource for reference to the ID of the public archive (OS)
data "sakuracloud_archive2 "centos" {
  os_type = "centos"
}

# Define Disk
resource "sakuracloud_disk" "disk01" {
  name              = "disk01"
  source_archive_id = data.sakuracloud_archive.centos.id
}

# Define Server
resource "sakuracloud_server" "server01" {
  name  = "server01"
  disks = [sakuracloud_disk.disk01.id]
  
  ssh_key_ids       = [sakuracloud_ssh_key_gen.key.id]
  password          = var.password
  disable_pw_auth   = true
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

- [Terraform for Sakura Cloud Documents](https://sacloud.github.io/terraform-provider-sakuracloud/)

### Supported Resources / Data Resources

#### Resources

- [Server](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/server/)
- [Disk](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/disk/)
- [Archive](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/archive/)
- [ISO Image(CD-ROM)](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/cdrom/)
- [Switch](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/switch/)
- [Router](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/internet/)
- [Subnet](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/subnet/)
- [Packet Filter](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/packet_filter/)
- [Packet Filter(Rule)](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/packet_filter_rule/)
- [Bridge](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/bridge/)
- [Load Balancer](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/load_balancer/)
- [VPC Router](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/vpc_router/)
- [Database](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/database/)
- [NFS](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/nfs/)
- [SIM(Secure Mobile)](http://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/sim/)
- [Mobile Gateway](http://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/mobile_gateway/)
- [Startup Script](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/note/)
- [Public Key](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/ssh_key/)
- [Public Key(Generate)](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/ssh_key_gen/)
- [Icon](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/icon/)
- [Private Host](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/private_host/)
- [DNS](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/dns/)
- [GSLB](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/gslb/)
- [Simple Monitoring](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/simple_monitor/)
- [Auto Backup](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/auto_backup/)
- [Object Storage](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/bucket_object/)
- [Server Connector](https://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/server_connector)

#### Data Resources

- [Data Resources](http://sacloud.github.io/terraform-provider-sakuracloud/configuration/resources/data_resource/)

#### Unsupported Resources

The following resources are unsupported because API is not provided by Sakura Cloud.

- Local Router
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

### Dependent Library

```bash
# List display
govendor list

# Batch updating libraries under vender
govendor fetch +v

# Update library under rvendor from GOPATH
govendor update +v
```

### Documents

The document uses Github Pages. (Under the `docs` directory of the master branch.)

The static file is generated by the `mkdocs` commands.

**To PR for the document, please modify only under the `build_docs` directory.**

**Do not include `docs` directory in PR.**

**The `docs` directory is updated in bulk at the time of release.**

```bash
  # Launch server for document preview.
  # You can preview with `http://localhost/`.
  make serve-docs
  
  # Document validation (textlint)
  make lint-docs
```

#### Preview terraform.io style documents

You can preview terraform.io style documents in english.  
Run following command and open `http://localhost:4567/docs/providers/sakuracloud`.  

```bash
# preview terraform.io style docs
make serve-english-docs 
```

## License

  This project is published under [Apache 2.0 License](LICENSE).

## Author

- [Terraform for Sakura Cloud Authors](AUTHORS)
