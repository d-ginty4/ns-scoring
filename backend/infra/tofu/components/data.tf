data "aws_caller_identity" "current" {}

data "aws_cognito_user_pools" "main" {
  name = var.cognito_user_pool_name
}

locals {
  cognito_user_pool_id = tolist(data.aws_cognito_user_pools.main.ids)[0]
}
