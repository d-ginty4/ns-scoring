// TODO: Define tables
resource "aws_dynamodb_table" "main" {
  name         = "${var.app_name}-${var.environment}-main"
  billing_mode = "PAY_PER_REQUEST"

  hash_key  = "postId"
  range_key = "meta"

  attribute {
    name = "postId"
    type = "S"
  }

  attribute {
    name = "meta"
    type = "S"
  }
}

resource "aws_dynamodb_table" "users" {
  name         = "${var.app_name}-${var.environment}-users"
  billing_mode = "PAY_PER_REQUEST"

  hash_key = "id"

  attribute {
    name = "id"
    type = "S"
  }
}
