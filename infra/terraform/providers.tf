terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    datadog = {
      source  = "datadog/datadog"
      version = "~> 3.0"
    }
  }

  # Backend configuration will be provided via -backend-config
  backend "s3" {}
}

provider "aws" {
  region = var.network.region

  default_tags {
    tags = {
      Project : var.project_name
      Environment : var.environment
      ManagedBy : "Terraform"
    }
  }
}
