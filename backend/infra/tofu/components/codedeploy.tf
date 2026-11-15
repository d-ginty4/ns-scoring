locals {
  is_prod = var.environment == "prod"
}

# CodeDeploy Permissions
data "aws_iam_policy_document" "codedeploy_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      identifiers = ["codedeploy.amazonaws.com"]
      type        = "Service"
    }
  }
}

resource "aws_iam_role" "codedeploy_role" {
  name               = "${var.app_name}-${var.environment}-codedeploy-role"
  assume_role_policy = data.aws_iam_policy_document.codedeploy_assume_role.json
}

resource "aws_iam_role_policy_attachment" "codedeploy_permissions" {
  role       = aws_iam_role.codedeploy_role.id
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSCodeDeployRoleForLambda"
}

resource "aws_iam_role_policy" "codedeploy_policy_s3" {
  role = aws_iam_role.codedeploy_role.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "s3:GetObject",
        "s3:GetObjectVersion",
      ]
      Resource = [
        "${aws_s3_bucket.storage.arn}/*",
        "${aws_s3_bucket.storage.arn}",
      ]
    }]
  })
}

resource "aws_codedeploy_app" "app" {
  name             = "${var.app_name}-${var.environment}-apis-lambda"
  compute_platform = "Lambda"
}

resource "aws_codedeploy_deployment_group" "deployment_group" {
  app_name               = aws_codedeploy_app.app.name
  deployment_group_name  = "${var.app_name}-${var.environment}-apis-deployment-group"
  service_role_arn       = aws_iam_role.codedeploy_role.arn
  deployment_config_name = local.is_prod ? "CodeDeployDefault.LambdaLinear10PercentEvery10Minutes" : "CodeDeployDefault.LambdaAllAtOnce"

  deployment_style {
    deployment_type   = "BLUE_GREEN"
    deployment_option = "WITH_TRAFFIC_CONTROL"
  }

  auto_rollback_configuration {
    enabled = local.is_prod
    events  = local.is_prod ? ["DEPLOYMENT_FAILURE", "DEPLOYMENT_STOP_ON_ALARM"] : ["DEPLOYMENT_FAILURE"]
  }

  alarm_configuration {
    enabled = local.is_prod
    alarms = local.is_prod ? [
      aws_cloudwatch_metric_alarm.lambda_errors[0].alarm_name,
      aws_cloudwatch_metric_alarm.lambda_duration[0].alarm_name,
    ] : []
  }
}

resource "aws_cloudwatch_metric_alarm" "lambda_errors" {
  count = local.is_prod ? 1 : 0

  alarm_name          = "${var.app_name}-${var.environment}-lambda-errors"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 1
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = 60
  statistic           = "Sum"
  threshold           = 0
  treat_missing_data  = "notBreaching"

  dimensions = {
    FunctionName = aws_lambda_function.apis.function_name
    Resource     = "${aws_lambda_function.apis.function_name}:${aws_lambda_alias.apis.name}"
  }
}

resource "aws_cloudwatch_metric_alarm" "lambda_duration" {
  count = local.is_prod ? 1 : 0

  alarm_name          = "${var.app_name}-${var.environment}-lambda-duration"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "Duration"
  namespace           = "AWS/Lambda"
  period              = 60
  extended_statistic  = "p99"
  threshold           = 25000 # 25 seconds (function timeout is 30s)
  treat_missing_data  = "notBreaching"

  dimensions = {
    FunctionName = aws_lambda_function.apis.function_name
    Resource     = "${aws_lambda_function.apis.function_name}:${aws_lambda_alias.apis.name}"
  }
}
