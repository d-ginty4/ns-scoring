variable "region" {
  type = string
}

variable "environment" {
  type = string
}

variable "app_name" {
  type = string
}

variable "cognito_user_pool_name" {
  type = string
}

variable "cloudfront_domain_name" {
  type = string
}

variable "subdomain" {
  type    = string
  default = ""
}
