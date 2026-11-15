variable "app_name" {
  type    = string
}

variable "environment" {
  type = string
}

variable "region" {
  type = string
}

variable "github_repo" {
  type    = string
}

variable "artifact_bucket" {
  type        = string
  description = "S3 bucket for Lambda deployment artifacts (app stack)"
}

variable "lambda_function_name" {
  type = string
}

variable "lambda_alias_name" {
  type    = string
  default = "live"
}

variable "codedeploy_app_name" {
  type = string
}

variable "codedeploy_deployment_group_name" {
  type = string
}
