# ============================================================================
# EC2 Bastion Host Module Outputs
# ============================================================================

# ----------------------------------------------------------------------------
# Instance Outputs
# ----------------------------------------------------------------------------

output "instance_id" {
  description = "ID of the bastion host EC2 instance"
  value       = aws_instance.bastion.id
}

output "instance_arn" {
  description = "ARN of the bastion host EC2 instance"
  value       = aws_instance.bastion.arn
}

output "instance_state" {
  description = "State of the bastion host EC2 instance"
  value       = aws_instance.bastion.instance_state
}

output "instance_type" {
  description = "Instance type of the bastion host"
  value       = aws_instance.bastion.instance_type
}

# ----------------------------------------------------------------------------
# Network Outputs
# ----------------------------------------------------------------------------

output "public_ip" {
  description = "Public IP address of the bastion host"
  value       = var.allocate_elastic_ip ? aws_eip.bastion[0].public_ip : aws_instance.bastion.public_ip
}

output "private_ip" {
  description = "Private IP address of the bastion host"
  value       = aws_instance.bastion.private_ip
}

output "elastic_ip_id" {
  description = "Allocation ID of the Elastic IP (if allocated)"
  value       = var.allocate_elastic_ip ? aws_eip.bastion[0].id : ""
}

output "public_dns" {
  description = "Public DNS name of the bastion host"
  value       = aws_instance.bastion.public_dns
}

output "private_dns" {
  description = "Private DNS name of the bastion host"
  value       = aws_instance.bastion.private_dns
}

# ----------------------------------------------------------------------------
# Security Outputs
# ----------------------------------------------------------------------------

output "security_group_id" {
  description = "ID of the bastion host security group"
  value       = aws_security_group.bastion.id
}

output "security_group_name" {
  description = "Name of the bastion host security group"
  value       = aws_security_group.bastion.name
}

output "iam_role_arn" {
  description = "ARN of the bastion host IAM role"
  value       = aws_iam_role.bastion.arn
}

output "iam_role_name" {
  description = "Name of the bastion host IAM role"
  value       = aws_iam_role.bastion.name
}

output "key_pair_name" {
  description = "Name of the SSH key pair"
  value       = aws_key_pair.bastion.key_name
}

# ----------------------------------------------------------------------------
# Connection Information
# ----------------------------------------------------------------------------

output "ssh_command" {
  description = "SSH command to connect to the bastion host"
  value       = "ssh -i ~/.ssh/<your-private-key> ec2-user@${var.allocate_elastic_ip ? aws_eip.bastion[0].public_ip : aws_instance.bastion.public_ip}"
}

output "ssm_command" {
  description = "AWS SSM Session Manager command to connect to the bastion host"
  value       = "aws ssm start-session --target ${aws_instance.bastion.id}"
}
