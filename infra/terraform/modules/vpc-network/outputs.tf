output "vpc_id" {
  description = "The ID of the VPC network created by this module."
  value       = aws_vpc.network.id
}

output "public_subnet_ids" {
  description = "The IDs of the Public Subnets created by this module"
  value       = [for subnet in aws_subnet.public_subnets : subnet.id]
}

output "private_subnet_ids" {
  description = "The IDs of the Private Subnets created by this module"
  value       = [for subnet in aws_subnet.private_subnets : subnet.id]
}

output "nat_gateway_ids" {
  description = "The IDs of the NAT Gateways created by this module"
  value       = [for nat in aws_nat_gateway.nat_gateways : nat.id]
}

output "nat_gateway_public_ips" {
  description = "The public IP addresses of the NAT Gateways"
  value       = [for eip in aws_eip.nat_gateway_eips : eip.public_ip]
}
