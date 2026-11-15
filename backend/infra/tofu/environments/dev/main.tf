terraform {
  backend "s3" {
    bucket       = "ns-scoring-terraform-backend-590184097018-eu-west-1"
    key          = "dev/backend.tfstate"
    region       = "eu-west-1"
    use_lockfile = true
  }
}

data "terraform_remote_state" "frontend" {
  backend = "s3"
  config = {
    bucket = "ns-scoring-terraform-frontend-590184097018-eu-west-1"
    key    = "dev/frontend.tfstate"
    region = "eu-west-1"
  }
}

module "components" {
  source = "../../components"

  region                 = var.region
  environment            = var.environment
  app_name               = var.app_name
  cognito_user_pool_name = data.terraform_remote_state.frontend.outputs.cognito_user_pool_name
  cloudfront_domain_name = data.terraform_remote_state.frontend.outputs.cloudfront_domain_name
}
