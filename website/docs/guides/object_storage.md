---
layout: "sakuracloud"
page_title: "Guide: Object Storage"
subcategory: "Guides"
description: |-
How to Manage Object Storage in SAKURA Cloud with Terraform.
---

# Object Storage

This provider does not support Sakura's Cloud object storage.
You can use the AWS provider to manage object storage instead.

```terraform
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  alias      = "sacloud"
  region     = var.s3_region
  access_key = var.s3_access_key
  secret_key = var.s3_secret_key
  endpoints {
    s3 = var.s3_endpoint
  }

  skip_credentials_validation = true
  skip_region_validation      = true
  skip_requesting_account_id  = true
  skip_metadata_api_check     = true
}

#================================================

variable "s3_access_key" {}
variable "s3_secret_key" {}

variable "s3_region" {
  default = "jp-north-1"
}

variable "s3_endpoint" {
  default = "https://s3.isk01.sakurastorage.jp"
}

#================================================

data "aws_s3_bucket" "example" {
  provider = aws.sacloud

  bucket   = "your-bucket-name" # set your bucket name
}

resource "aws_s3_bucket_object" "object" {
  provider = aws.sacloud

  bucket = data.aws_s3_bucket.example.id
  key    = "example.txt"
  source = "example.txt"
  etag   = filemd5("example.txt")
  acl    = "public-read"
}
```