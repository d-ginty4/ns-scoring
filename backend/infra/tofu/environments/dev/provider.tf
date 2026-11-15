terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.70.0"
    }
  }
}

provider "aws" {
  region  = var.region
  profile = var.profile != "" ? var.profile : null

  default_tags {
    tags = {
      Project = var.app_name
    }
  }
}

