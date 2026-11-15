resource "aws_s3_bucket" "storage" {
  bucket = "${var.app_name}-${var.environment}-${data.aws_caller_identity.current.account_id}-storage-${var.region}"
}

resource "aws_s3_bucket_versioning" "storage" {
  bucket = aws_s3_bucket.storage.id
  versioning_configuration {
    status = "Enabled"
  }
}
