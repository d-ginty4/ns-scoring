output "apis_lambda_function_name" {
  value = module.components.apis_lambda_function_name
}

output "apis_lambda_alias_name" {
  value = module.components.apis_lambda_alias_name
}

output "codedeploy_app_name" {
  value = module.components.codedeploy_app_name
}

output "codedeploy_deployment_group_name" {
  value = module.components.codedeploy_deployment_group_name
}

output "storage_bucket" {
  value = module.components.storage_bucket
}

output "api_gateway_url" {
  value = module.components.api_gateway_url
}
