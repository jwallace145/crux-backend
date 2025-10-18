variable "service_name" {
  type        = string
  description = "The name of the VPC network to be created."
  default     = "crux-backend"
}

variable "environment" {
  description = "The environment of CruxBackend."
  type        = string

  validation {
    condition     = contains(["dev", "stg", "prod"], var.environment)
    error_message = "The environment must be 'dev', 'stg', or 'prod'."
  }
}

variable "region" {
  description = "The AWS region of the VPC network to be created."
  type        = string

  validation {
    condition     = contains(local.valid_regions, var.region)
    error_message = "The AWS region must be 'us-east-1'."
  }
}

variable "vpc_cidr_block" {
  type        = string
  description = "The CIDR block of the VPC network to be created."
}

variable "public_subnets" {
  description = "The public subnets to create in the VPC network with access to the external web."

  type = map(object({
    availability_zone = string
    subnet_cidr_block = string
  }))
}

variable "enable_dns_hostnames" {
  type        = bool
  description = "Enable DNS hostnames in VPC network."
  default     = true
}

variable "enable_dns_support" {
  type        = bool
  description = "Enable DNS support in VPC network."
  default     = true
}
