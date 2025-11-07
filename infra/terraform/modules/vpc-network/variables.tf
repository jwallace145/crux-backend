variable "project_name" {
  type        = string
  description = "The project name for tagging resources"
}

variable "environment" {
  type        = string
  description = "The environment (dev, stg, prod)"
}

variable "region" {
  type        = string
  description = "The AWS region for the VPC"
}

variable "vpc_cidr_block" {
  type        = string
  description = "The CIDR block for the VPC"
}

variable "enable_dns_hostnames" {
  type        = bool
  description = "Enable DNS hostnames in the VPC"
  default     = true
}

variable "enable_dns_support" {
  type        = bool
  description = "Enable DNS support in the VPC"
  default     = true
}

variable "public_subnets" {
  description = "Map of public subnet configurations"
  type = map(object({
    availability_zone = string
    subnet_cidr_block = string
  }))
}

variable "private_subnets" {
  description = "Map of private subnet configurations"
  type = map(object({
    availability_zone = string
    subnet_cidr_block = string
  }))
  default = {}
}

variable "nat_gateways" {
  description = "Map of NAT gateway configurations. Each enabled NAT gateway must have a corresponding public subnet in the same AZ."
  type = map(object({
    availability_zone = string
    enabled           = bool
  }))
  default = {}

  validation {
    condition = alltrue([
      for key, nat in var.nat_gateways :
      !nat.enabled || anytrue([
        for subnet in var.public_subnets :
        subnet.availability_zone == nat.availability_zone
      ])
    ])
    error_message = "Each enabled NAT gateway must have at least one corresponding public subnet in the same availability zone."
  }

  validation {
    condition = alltrue([
      for key, nat in var.nat_gateways :
      contains(lookup(local.valid_availability_zones, var.region, []), nat.availability_zone)
    ])
    error_message = "All NAT gateways must be in valid availability zones for the specified region."
  }
}
