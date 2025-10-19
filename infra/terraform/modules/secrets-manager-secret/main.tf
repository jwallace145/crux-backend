# AWS Secrets Manager Secret
# Stores multiple key-value pairs as a JSON object

resource "aws_secretsmanager_secret" "this" {
  name        = var.secret_name
  description = var.secret_description

  # KMS encryption
  kms_key_id = var.kms_key_id

  # Recovery window configuration
  recovery_window_in_days = var.recovery_window_in_days

  # Tags
  tags = merge(
    var.tags,
    {
      Name       = var.secret_name
      ManagedBy  = "Terraform"
      SecretType = "JSON"
    }
  )
}

# Store the secrets as JSON
resource "aws_secretsmanager_secret_version" "this" {
  secret_id = aws_secretsmanager_secret.this.id

  # Convert the map of secrets to JSON
  secret_string = jsonencode(var.secrets)
}

# Optional: Configure automatic rotation
resource "aws_secretsmanager_secret_rotation" "this" {
  count = var.enable_rotation ? 1 : 0

  secret_id           = aws_secretsmanager_secret.this.id
  rotation_lambda_arn = var.rotation_lambda_arn

  rotation_rules {
    automatically_after_days = var.rotation_days
  }

  depends_on = [aws_secretsmanager_secret_version.this]
}
