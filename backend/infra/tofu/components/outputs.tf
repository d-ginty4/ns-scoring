output "apis_lambda_function_name" {
  value = aws_lambda_function.apis.function_name
}

output "apis_lambda_alias_name" {
  value = aws_lambda_alias.apis.name
}

output "codedeploy_app_name" {
  value = aws_codedeploy_app.app.name
}

output "codedeploy_deployment_group_name" {
  value = aws_codedeploy_deployment_group.deployment_group.deployment_group_name
}

output "storage_bucket" {
  value = aws_s3_bucket.storage.bucket
}

output "api_gateway_url" {
  value = aws_apigatewayv2_api.gateway.api_endpoint
}
