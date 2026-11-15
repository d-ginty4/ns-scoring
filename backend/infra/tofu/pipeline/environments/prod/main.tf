terraform {
  backend "s3" {
    bucket       = "ns-scoring-terraform-backend-590184097018-eu-west-1"
    key          = "prod/pipeline.tfstate"
    region       = "eu-west-1"
    use_lockfile = true
  }
}

# Reference to the main stack
data "terraform_remote_state" "main" {
  backend = "s3"
  config = {
    bucket = "ns-scoring-terraform-backend-590184097018-eu-west-1"
    key    = "prod/backend.tfstate"
    region = "eu-west-1"
  }
}

module "pipeline" {
  source = "../../components"

  region      = var.region
  environment = var.environment
  app_name    = var.app_name

  github_repo = var.github_repo

  artifact_bucket                  = data.terraform_remote_state.main.outputs.storage_bucket
  lambda_function_name             = data.terraform_remote_state.main.outputs.apis_lambda_function_name
  lambda_alias_name                = data.terraform_remote_state.main.outputs.apis_lambda_alias_name
  codedeploy_app_name              = data.terraform_remote_state.main.outputs.codedeploy_app_name
  codedeploy_deployment_group_name = data.terraform_remote_state.main.outputs.codedeploy_deployment_group_name
}
