data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      identifiers = ["lambda.amazonaws.com"]
      type        = "Service"
    }
  }
}

resource "aws_iam_role" "lambda_role" {
  name               = "${var.app_name}-${var.environment}-lambda-execution-role"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
}

resource "aws_iam_role_policy_attachment" "basic_permissions" {
  role       = aws_iam_role.lambda_role.id
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy" "lambda_policy_dynamo" {
  role = aws_iam_role.lambda_role.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:Query",
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem"
        ]
        Resource = [
          aws_dynamodb_table.main.arn,
          "${aws_dynamodb_table.main.arn}/index/*",
          aws_dynamodb_table.users.arn,
          "${aws_dynamodb_table.users.arn}/index/*"
        ]
      },
    ]
  })
}

resource "aws_lambda_function" "apis" {
  function_name = "${var.app_name}-${var.environment}-apis-controller"
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  role          = aws_iam_role.lambda_role.arn
  memory_size   = 128
  timeout       = 30
  architectures = ["arm64"]
  s3_bucket     = aws_s3_bucket.storage.bucket
  s3_key        = "apis.zip"
  publish       = true
  environment {
    variables = {
      APP_ENV          = var.environment
      COGNITO_POOL_ID  = local.cognito_user_pool_id
    }
  }
}

resource "aws_lambda_alias" "apis" {
  name             = "live"
  function_name    = aws_lambda_function.apis.arn
  function_version = aws_lambda_function.apis.version

  lifecycle {
    ignore_changes = [function_version]
  }
}

resource "aws_lambda_permission" "gateway-apis-permission" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_alias.apis.function_name
  qualifier     = aws_lambda_alias.apis.name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"
}
