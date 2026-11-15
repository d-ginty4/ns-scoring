data "aws_iam_policy_document" "codebuild_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["codebuild.amazonaws.com"]
    }
  }
}

# Infrastructure Build
resource "aws_iam_role" "codebuild_infra" {
  name               = "${var.app_name}-${var.environment}-codebuild-infra-role"
  assume_role_policy = data.aws_iam_policy_document.codebuild_assume_role.json
}

data "aws_iam_policy_document" "codebuild_infra_policy" {
  statement {
    sid    = "CloudWatchLogs"
    effect = "Allow"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = [
      "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:/aws/codebuild/${var.app_name}-${var.environment}-infra",
      "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:/aws/codebuild/${var.app_name}-${var.environment}-infra:*",
    ]
  }
  statement {
    effect = "Allow"
    actions = [
      "s3:GetObject",
      "s3:GetObjectVersion",
      "s3:GetBucketVersioning",
      "s3:PutObject",
    ]
    resources = [
      "arn:aws:s3:::${var.artifact_bucket}",
      "arn:aws:s3:::${var.artifact_bucket}/*",
    ]
  }
  statement {
    sid    = "TerraformState"
    effect = "Allow"
    actions = [
      "s3:GetObject",
      "s3:HeadObject",
      "s3:PutObject",
      "s3:DeleteObject",
      "s3:ListBucket",
    ]
    resources = [
      "arn:aws:s3:::${var.app_name}-terraform-backend-${data.aws_caller_identity.current.account_id}-${var.region}",
      "arn:aws:s3:::${var.app_name}-terraform-backend-${data.aws_caller_identity.current.account_id}-${var.region}/${var.environment}/*",
    ]
  }
  statement {
    sid    = "OpenTofuPermissions"
    effect = "Allow"
    actions = [
      "s3:*",
      "dynamodb:*",
      "apigateway:*",
      "lambda:*",
      "iam:*",
      "codedeploy:*",
      "cloudwatch:*",
      "cognito-idp:*",
    ]
    resources = [
      "*",
    ]
  }
}

resource "aws_iam_role_policy" "codebuild_infra" {
  name   = "${var.app_name}-${var.environment}-codebuild-infra-policy"
  role   = aws_iam_role.codebuild_infra.id
  policy = data.aws_iam_policy_document.codebuild_infra_policy.json
}

resource "aws_codebuild_project" "infra" {
  name          = "${var.app_name}-${var.environment}-infra"
  description   = "Apply OpenTofu app infrastructure"
  service_role  = aws_iam_role.codebuild_infra.arn
  build_timeout = 30

  artifacts {
    type = "CODEPIPELINE"
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "aws/codebuild/amazonlinux2-x86_64-standard:5.0"
    type                        = "LINUX_CONTAINER"
    image_pull_credentials_type = "CODEBUILD"

    environment_variable {
      name  = "ENVIRONMENT"
      value = var.environment
    }
  }

  logs_config {
    cloudwatch_logs {
      group_name = "/aws/codebuild/${var.app_name}-${var.environment}-infra"
    }
  }

  source {
    type      = "CODEPIPELINE"
    buildspec = "infra/tofu/pipeline/buildspec/infra.yml"
  }
}

# Lambda Build
resource "aws_iam_role" "codebuild_lambda" {
  name               = "${var.app_name}-${var.environment}-codebuild-lambda-role"
  assume_role_policy = data.aws_iam_policy_document.codebuild_assume_role.json
}

data "aws_iam_policy_document" "codebuild_lambda_policy" {
  statement {
    sid    = "CloudWatchLogs"
    effect = "Allow"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = [
      "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:/aws/codebuild/${var.app_name}-${var.environment}-lambda",
      "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:/aws/codebuild/${var.app_name}-${var.environment}-lambda:*",
    ]
  }
  statement {
    effect = "Allow"
    actions = [
      "s3:GetObject",
      "s3:GetObjectVersion",
      "s3:GetBucketVersioning",
      "s3:PutObject",
    ]
    resources = [
      "arn:aws:s3:::${var.artifact_bucket}",
      "arn:aws:s3:::${var.artifact_bucket}/*",
    ]
  }
  statement {
    sid    = "LambdaDeploy"
    effect = "Allow"
    actions = [
      "lambda:UpdateFunctionCode",
      "lambda:GetAlias",
      "lambda:PublishVersion",
    ]
    resources = [
      "arn:aws:lambda:${var.region}:${data.aws_caller_identity.current.account_id}:function:${var.lambda_function_name}",
      "arn:aws:lambda:${var.region}:${data.aws_caller_identity.current.account_id}:function:${var.lambda_function_name}:*",
    ]
  }
}

resource "aws_iam_role_policy" "codebuild_lambda" {
  name   = "${var.app_name}-${var.environment}-codebuild-lambda-policy"
  role   = aws_iam_role.codebuild_lambda.id
  policy = data.aws_iam_policy_document.codebuild_lambda_policy.json
}

resource "aws_codebuild_project" "lambda" {
  name          = "${var.app_name}-${var.environment}-lambda"
  description   = "Build and publish Lambda deployment artifact"
  service_role  = aws_iam_role.codebuild_lambda.arn
  build_timeout = 15

  artifacts {
    type = "CODEPIPELINE"
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "aws/codebuild/amazonlinux2-x86_64-standard:5.0"
    type                        = "LINUX_CONTAINER"
    image_pull_credentials_type = "CODEBUILD"

    environment_variable {
      name  = "ARTIFACT_BUCKET"
      value = var.artifact_bucket
    }
    environment_variable {
      name  = "LAMBDA_FUNCTION_NAME"
      value = var.lambda_function_name
    }
    environment_variable {
      name  = "LAMBDA_ALIAS_NAME"
      value = var.lambda_alias_name
    }
    environment_variable {
      name  = "CODEDEPLOY_APP_NAME"
      value = var.codedeploy_app_name
    }
    environment_variable {
      name  = "CODEDEPLOY_DG_NAME"
      value = var.codedeploy_deployment_group_name
    }
  }

  logs_config {
    cloudwatch_logs {
      group_name = "/aws/codebuild/${var.app_name}-${var.environment}-lambda"
    }
  }

  source {
    type      = "CODEPIPELINE"
    buildspec = "infra/tofu/pipeline/buildspec/lambda.yml"
  }
}
