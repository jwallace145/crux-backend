# ============================================================================
# Application Load Balancer Module Outputs
# ============================================================================

# ----------------------------------------------------------------------------
# ALB Outputs
# ----------------------------------------------------------------------------

output "alb_dns_name" {
  description = "DNS name of the Application Load Balancer - use this as a CNAME target for custom domains"
  value       = aws_lb.main.dns_name
}

output "alb_arn" {
  description = "ARN of the Application Load Balancer"
  value       = aws_lb.main.arn
}

output "alb_zone_id" {
  description = "Canonical hosted zone ID of the ALB (for Route53 alias records)"
  value       = aws_lb.main.zone_id
}

output "alb_id" {
  description = "ID of the Application Load Balancer"
  value       = aws_lb.main.id
}

# ----------------------------------------------------------------------------
# Target Group Outputs
# ----------------------------------------------------------------------------

output "target_group_arn" {
  description = "ARN of the target group for ECS service integration"
  value       = aws_lb_target_group.ecs.arn
}

output "target_group_name" {
  description = "Name of the target group"
  value       = aws_lb_target_group.ecs.name
}

# ----------------------------------------------------------------------------
# Security Group Outputs
# ----------------------------------------------------------------------------

output "alb_security_group_id" {
  description = "ID of the ALB security group"
  value       = aws_security_group.alb.id
}

output "ecs_tasks_security_group_id" {
  description = "ID of the ECS tasks security group - use this for ECS service"
  value       = aws_security_group.ecs_tasks.id
}

# ----------------------------------------------------------------------------
# Listener Outputs
# ----------------------------------------------------------------------------

output "http_listener_arn" {
  description = "ARN of the HTTP listener"
  value       = aws_lb_listener.http.arn
}

output "https_listener_arn" {
  description = "ARN of the HTTPS listener (empty if HTTPS is not enabled)"
  value       = var.enable_https ? aws_lb_listener.https[0].arn : ""
}

# ----------------------------------------------------------------------------
# Connection Information
# ----------------------------------------------------------------------------

output "alb_url" {
  description = "Full HTTP URL of the Application Load Balancer"
  value       = "http://${aws_lb.main.dns_name}"
}

output "alb_https_url" {
  description = "Full HTTPS URL of the Application Load Balancer (if HTTPS is enabled)"
  value       = var.enable_https ? "https://${aws_lb.main.dns_name}" : ""
}
