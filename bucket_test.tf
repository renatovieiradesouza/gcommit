terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

variable "bucket_name" {
  type        = string
  description = "Nome unico do bucket S3 para teste."
  default     = "trocar-para-um-bucket-unico"
}

resource "aws_s3_bucket" "test_bucket" {
  bucket = var.bucket_name

  tags = {
    Name       = "gcommit-test-bucket"
    ManagedBy  = "terraform"
    Repository = "gcommit"
  }
}
