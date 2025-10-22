# ============================================================================
# Crux Backend Infrastructure Outputs
# ============================================================================

# ============================================================================
# Network Information
# ============================================================================

output "vpc_id" {
  description = "The ID of the VPC"
  value       = module.network.vpc_id
}

output "public_subnet_ids" {
  description = "List of public subnet IDs"
  value       = module.network.public_subnet_ids
}

output "private_subnet_ids" {
  description = "List of private subnet IDs"
  value       = module.network.private_subnet_ids
}

output "nat_gateway_ids" {
  description = "List of NAT Gateway IDs"
  value       = module.network.nat_gateway_ids
}

output "nat_gateway_public_ips" {
  description = "List of NAT Gateway public IP addresses"
  value       = module.network.nat_gateway_public_ips
}

# ============================================================================
# API / Application Information
# ============================================================================

output "api_url" {
  description = "The HTTPS URL to access the API"
  value       = module.alb.alb_https_url
}

output "api_alb_dns_name" {
  description = "The DNS name of the Application Load Balancer (for CNAME records)"
  value       = module.alb.alb_dns_name
}

output "api_alb_zone_id" {
  description = "The canonical hosted zone ID of the ALB (for Route53 alias records)"
  value       = module.alb.alb_zone_id
}

output "api_domain" {
  description = "The custom domain configured for the API"
  value       = var.api.domain
}

output "api_cluster_name" {
  description = "The name of the ECS cluster running the API"
  value       = module.api.cluster_name
}

output "api_cluster_arn" {
  description = "The ARN of the ECS cluster"
  value       = module.api.cluster_arn
}

output "api_service_name" {
  description = "The name of the ECS service"
  value       = module.api.service_name
}

output "api_task_definition_arn" {
  description = "The ARN of the current task definition"
  value       = module.api.task_definition_arn
}

output "api_log_group_name" {
  description = "The CloudWatch log group name for API logs"
  value       = module.api.log_group_name
}

output "api_container_port" {
  description = "The port the API container listens on"
  value       = var.api.container.port
}

# ============================================================================
# SSL Certificate Information
# ============================================================================

output "ssl_certificate_arn" {
  description = "The ARN of the SSL certificate for HTTPS"
  value       = module.certificate.certificate_arn
}

output "ssl_certificate_domain" {
  description = "The domain name of the SSL certificate"
  value       = module.certificate.certificate_domain_name
}

output "ssl_certificate_status" {
  description = "The status of the SSL certificate"
  value       = module.certificate.certificate_status
}

# ============================================================================
# Database Information
# ============================================================================

output "db_endpoint" {
  description = "The database connection endpoint (hostname:port)"
  value       = module.db.endpoint
}

output "db_address" {
  description = "The database hostname (without port)"
  value       = module.db.address
}

output "db_port" {
  description = "The database port"
  value       = module.db.port
}

output "db_name" {
  description = "The name of the database"
  value       = module.db.database_name
}

output "db_username" {
  description = "The master username for the database"
  value       = module.db.username
  sensitive   = true
}

output "db_connection_string" {
  description = "PostgreSQL connection string (password not included)"
  value       = module.db.connection_string
  sensitive   = true
}

output "db_availability_zone" {
  description = "The availability zone where the database is located"
  value       = module.db.availability_zone
}

output "db_instance_id" {
  description = "The RDS instance identifier"
  value       = module.db.id
}

output "db_arn" {
  description = "The ARN of the RDS instance"
  value       = module.db.arn
}

# ============================================================================
# Container Registry Information
# ============================================================================

output "ecr_repository_url" {
  description = "The URL of the ECR repository for the API"
  value       = module.crux_api_ecr.repository_url
}

output "ecr_repository_name" {
  description = "The name of the ECR repository"
  value       = module.crux_api_ecr.repository_name
}

output "ecr_repository_arn" {
  description = "The ARN of the ECR repository"
  value       = module.crux_api_ecr.repository_arn
}

# ============================================================================
# Secrets Manager Information
# ============================================================================

output "db_secrets_arn" {
  description = "The ARN of the database secrets in AWS Secrets Manager"
  value       = module.db_secrets.secret_arn
}

output "jwt_secrets_arn" {
  description = "The ARN of the JWT secrets in AWS Secrets Manager"
  value       = module.jwt_secrets.secret_arn
}

# ============================================================================
# Quick Reference / Connection Information
# ============================================================================

output "deployment_summary" {
  description = "Quick reference summary of the deployment"
  value = {
    environment      = var.environment
    region           = var.network.region
    api_url          = module.alb.alb_https_url
    api_domain       = var.api.domain
    db_endpoint      = module.db.endpoint
    ecs_cluster      = module.api.cluster_name
    ecs_service      = module.api.service_name
    logs             = "https://console.aws.amazon.com/cloudwatch/home?region=${var.network.region}#logsV2:log-groups/log-group/${module.api.log_group_name}"
    ecr_push_command = "aws ecr get-login-password --region ${var.network.region} | docker login --username AWS --password-stdin ${module.crux_api_ecr.repository_url}"
  }
}
