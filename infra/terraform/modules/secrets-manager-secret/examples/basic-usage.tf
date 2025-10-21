# Example: Store JWT secrets in AWS Secrets Manager
# This example shows how to use the secrets-manager-secret module
# to store JWT token secrets for the Crux API

# Store JWT secrets as a single JSON object
module "jwt_secrets" {
  source = "../"

  secret_name        = "crux-api/jwt-secrets-dev"
  secret_description = "JWT token secrets for Crux API (Development)"

  # Pass multiple secrets as a map
  secrets = {
    ACCESS_TOKEN_SECRET_KEY  = var.access_token_secret_key
    REFRESH_TOKEN_SECRET_KEY = var.refresh_token_secret_key
  }

  # Recovery window for accidental deletion
  recovery_window_in_days = 30

  tags = {
    Environment = "dev"
    Application = "crux-api"
    ManagedBy   = "Terraform"
    Purpose     = "JWT Authentication"
  }
}

# Store db credentials as a single JSON object
module "database_secrets" {
  source = "../"

  secret_name        = "crux-api/database-credentials-dev"
  secret_description = "Database credentials for Crux API (Development)"

  secrets = {
    DB_HOST     = "cruxdb-dev.us-east-1.rds.amazonaws.com"
    DB_PORT     = "5432"
    DB_NAME     = "cruxdb"
    DB_USER     = "cruxadmin"
    DB_PASSWORD = var.database_password
  }

  recovery_window_in_days = 30

  tags = {
    Environment = "dev"
    Application = "crux-api"
    ManagedBy   = "Terraform"
    Purpose     = "Database Access"
  }
}

# Example outputs
output "jwt_secrets_arn" {
  description = "ARN of the JWT secrets in Secrets Manager"
  value       = module.jwt_secrets.secret_arn
}

output "database_secrets_arn" {
  description = "ARN of the database secrets in Secrets Manager"
  value       = module.database_secrets.secret_arn
}

# Example: Reference in IAM policy
data "aws_iam_policy_document" "ecs_task_secrets_access" {
  statement {
    effect = "Allow"
    actions = [
      "secretsmanager:GetSecretValue"
    ]
    resources = [
      module.jwt_secrets.secret_arn,
      module.database_secrets.secret_arn
    ]
  }

  # Also allow decryption if using KMS
  statement {
    effect = "Allow"
    actions = [
      "kms:Decrypt"
    ]
    resources = ["*"]
    condition {
      test     = "StringEquals"
      variable = "kms:ViaService"
      values   = ["secretsmanager.us-east-1.amazonaws.com"]
    }
  }
}

# Variables used in this example
variable "access_token_secret_key" {
  description = "Secret key for JWT access tokens"
  type        = string
  sensitive   = true
}

variable "refresh_token_secret_key" {
  description = "Secret key for JWT refresh tokens"
  type        = string
  sensitive   = true
}

variable "database_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}
