output "secret_id" {
  description = "The ID of the secret (same as ARN)"
  value       = aws_secretsmanager_secret.this.id
}

output "secret_arn" {
  description = "The ARN of the secret"
  value       = aws_secretsmanager_secret.this.arn
}

output "secret_name" {
  description = "The name of the secret"
  value       = aws_secretsmanager_secret.this.name
}

output "secret_version_id" {
  description = "The version ID of the secret value"
  value       = aws_secretsmanager_secret_version.this.version_id
}

output "secret_json" {
  description = "The secret stored as JSON (sensitive)"
  value       = aws_secretsmanager_secret_version.this.secret_string
  sensitive   = true
}
