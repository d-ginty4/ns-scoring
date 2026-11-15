locals {
  allowed_origins = [
    var.subdomain != "" ? "https://${var.subdomain}" : "https://${var.cloudfront_domain_name}",
    var.subdomain != "" ? "https://${var.subdomain}/" : "https://${var.cloudfront_domain_name}/",
  ]
}

resource "aws_apigatewayv2_api" "gateway" {
  name          = "${var.app_name}-${var.environment}-gateway"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = local.allowed_origins
    allow_headers = ["authorization", "content-type"]
    allow_methods = ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
  }
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.gateway.id
  name        = "$default"
  auto_deploy = true
}
