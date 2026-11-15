resource "aws_apigatewayv2_integration" "apis_integration" {
  api_id                 = aws_apigatewayv2_api.gateway.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_alias.apis.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "root_route" {
  api_id    = aws_apigatewayv2_api.gateway.id
  route_key = "ANY /"
  target    = "integrations/${aws_apigatewayv2_integration.apis_integration.id}"
}

resource "aws_apigatewayv2_route" "all_route" {
  api_id    = aws_apigatewayv2_api.gateway.id
  route_key = "ANY /{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.apis_integration.id}"
}
