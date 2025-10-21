# Primary db endpoint (hostname:port)
output "endpoint" {
  description = "The connection endpoint for the RDS instance (hostname:port)"
  value       = aws_db_instance.this.endpoint
}

# Database hostname only (without port)
output "address" {
  description = "The hostname of the RDS instance (without port)"
  value       = aws_db_instance.this.address
}

# Database port
output "port" {
  description = "The port the database is listening on"
  value       = aws_db_instance.this.port
}

# Database identifier
output "id" {
  description = "The RDS instance identifier"
  value       = aws_db_instance.this.id
}

# Database ARN
output "arn" {
  description = "The ARN of the RDS instance"
  value       = aws_db_instance.this.arn
}

# Database name
output "database_name" {
  description = "The name of the database"
  value       = aws_db_instance.this.db_name
}

# Database username
output "username" {
  description = "The master username for the database"
  value       = aws_db_instance.this.username
  sensitive   = true
}

# Security group ID
output "security_group_id" {
  description = "The ID of the security group attached to the RDS instance"
  value       = aws_security_group.this.id
}

# Subnet group name
output "subnet_group_name" {
  description = "The name of the DB subnet group"
  value       = aws_db_subnet_group.this.name
}

# Resource ID (for monitoring)
output "resource_id" {
  description = "The RDS Resource ID (for CloudWatch, etc.)"
  value       = aws_db_instance.this.resource_id
}

# Availability zone
output "availability_zone" {
  description = "The availability zone of the RDS instance"
  value       = aws_db_instance.this.availability_zone
}

# Connection string helper (useful for local testing)
output "connection_string" {
  description = "PostgreSQL connection string (password not included)"
  value       = "postgresql://${aws_db_instance.this.username}@${aws_db_instance.this.endpoint}/${aws_db_instance.this.db_name}?sslmode=require"
  sensitive   = true
}
