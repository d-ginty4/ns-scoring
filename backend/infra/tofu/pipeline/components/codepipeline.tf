data "aws_iam_policy_document" "codepipeline_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["codepipeline.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "codepipeline" {
  name               = "${var.app_name}-${var.environment}-codepipeline-role"
  assume_role_policy = data.aws_iam_policy_document.codepipeline_assume_role.json
}

data "aws_iam_policy_document" "codepipeline_policy" {
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
    effect    = "Allow"
    actions   = ["codestar-connections:UseConnection"]
    resources = [aws_codestarconnections_connection.repo.arn]
  }
  statement {
    effect = "Allow"
    actions = [
      "codebuild:BatchGetBuilds",
      "codebuild:StartBuild",
    ]
    resources = ["*"]
  }
  statement {
    effect = "Allow"
    actions = [
      "codedeploy:CreateDeployment",
      "codedeploy:GetApplication",
      "codedeploy:GetApplicationRevision",
      "codedeploy:GetDeployment",
      "codedeploy:GetDeploymentConfig",
      "codedeploy:RegisterApplicationRevision",
    ]
    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "codepipeline" {
  name   = "${var.app_name}-${var.environment}-codepipeline-policy"
  role   = aws_iam_role.codepipeline.id
  policy = data.aws_iam_policy_document.codepipeline_policy.json
}

resource "aws_codestarconnections_connection" "repo" {
  name          = "${var.app_name}-${var.environment}-github"
  provider_type = "GitHub"
}

resource "aws_codepipeline" "deploy" {
  name          = "${var.app_name}-${var.environment}-deploy"
  role_arn      = aws_iam_role.codepipeline.arn
  pipeline_type = "V2"

  artifact_store {
    location = var.artifact_bucket
    type     = "S3"
  }

  stage {
    name = "Source"

    action {
      name             = "Source"
      category         = "Source"
      owner            = "AWS"
      provider         = "CodeStarSourceConnection"
      version          = "1"
      output_artifacts = ["SourceOutput"]

      configuration = {
        ConnectionArn    = aws_codestarconnections_connection.repo.arn
        FullRepositoryId = var.github_repo
        BranchName       = "main"
        DetectChanges    = false
      }
    }
  }

  stage {
    name = "DeployInfra"

    action {
      name             = "DeployInfra"
      category         = "Build"
      owner            = "AWS"
      provider         = "CodeBuild"
      version          = "1"
      input_artifacts  = ["SourceOutput"]
      output_artifacts = ["InfraOutput"]

      configuration = {
        ProjectName = aws_codebuild_project.infra.name
      }
    }
  }

  stage {
    name = "BuildLambda"

    action {
      name             = "BuildLambda"
      category         = "Build"
      owner            = "AWS"
      provider         = "CodeBuild"
      version          = "1"
      input_artifacts  = ["SourceOutput"]
      output_artifacts = ["LambdaOutput"]

      configuration = {
        ProjectName = aws_codebuild_project.lambda.name
      }
    }
  }

  dynamic "stage" {
    for_each = var.environment == "prod" ? [1] : []
    content {
      name = "Approval"

      action {
        name     = "ManualApproval"
        category = "Approval"
        owner    = "AWS"
        provider = "Manual"
        version  = "1"

        configuration = {
          CustomData = "Approve deployment to ${var.environment}?"
        }
      }
    }
  }

  stage {
    name = "DeployApp"

    action {
      name            = "DeployApp"
      category        = "Deploy"
      owner           = "AWS"
      provider        = "CodeDeploy"
      version         = "1"
      input_artifacts = ["LambdaOutput"]

      configuration = {
        ApplicationName     = var.codedeploy_app_name
        DeploymentGroupName = var.codedeploy_deployment_group_name
      }
    }
  }
}
