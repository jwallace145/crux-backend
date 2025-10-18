output "vpc_id" {
  description = "The ID of the VPC network created by this module."
  value       = aws_vpc.network.id
}

output "public_subnet_ids" {
  description = "The IDs of the Public Subnets created by this module"
  value       = [for subnet in aws_subnet.public_subnets : subnet.id]
}
