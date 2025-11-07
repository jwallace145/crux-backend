# ============================================================================
# EC2 Bastion Host Module Variables
# ============================================================================

# ----------------------------------------------------------------------------
# Required Variables
# ----------------------------------------------------------------------------

variable "project_name" {
  description = "Name of the project (e.g., 'crux-project')"
  type        = string

  validation {
    condition     = length(var.project_name) > 0 && length(var.project_name) <= 32
    error_message = "Service name must be between 1 and 32 characters."
  }
}

variable "environment" {
  description = "Environment name (e.g., 'dev', 'staging', 'prod')"
  type        = string

  validation {
    condition     = contains(["dev", "stg", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: dev, stg, staging, prod."
  }
}

variable "vpc_id" {
  description = "ID of the VPC where the bastion host will be created"
  type        = string
}

variable "subnet_id" {
  description = "ID of the public subnet where the bastion host will be launched"
  type        = string
}

variable "allowed_ssh_cidr_blocks" {
  description = "List of CIDR blocks allowed to SSH into the bastion host (e.g., your local IP)"
  type        = list(string)

  validation {
    condition     = length(var.allowed_ssh_cidr_blocks) > 0
    error_message = "At least one CIDR block must be specified for SSH access."
  }
}

variable "ssh_public_key_path" {
  description = "Path to the SSH public key file (e.g., '~/.ssh/id_rsa.pub' or './keys/bastion.pub')"
  type        = string

  validation {
    condition     = length(var.ssh_public_key_path) > 0
    error_message = "SSH public key path must be provided."
  }
}

# ----------------------------------------------------------------------------
# EC2 Configuration
# ----------------------------------------------------------------------------

variable "instance_type" {
  description = "EC2 instance type for the bastion host"
  type        = string
  default     = "t3.micro"

  validation {
    condition     = can(regex("^t[2-4]\\.(nano|micro|small|medium)$", var.instance_type))
    error_message = "Instance type should be a small T-series instance (t2/t3/t4.nano/micro/small/medium)."
  }
}

variable "root_volume_size" {
  description = "Size of the root volume in GB"
  type        = number
  default     = 30

  validation {
    condition     = var.root_volume_size >= 30 && var.root_volume_size <= 50
    error_message = "Root volume size must be between 30 and 50 GB."
  }
}

# ----------------------------------------------------------------------------
# Network Configuration
# ----------------------------------------------------------------------------

variable "allocate_elastic_ip" {
  description = "Allocate an Elastic IP for the bastion host (recommended for consistent SSH access)"
  type        = bool
  default     = true
}

# ----------------------------------------------------------------------------
# Security Configuration
# ----------------------------------------------------------------------------

variable "enable_readonly_access" {
  description = "Attach ReadOnlyAccess IAM policy for AWS CLI debugging"
  type        = bool
  default     = false
}

variable "enable_termination_protection" {
  description = "Enable EC2 termination protection (recommended for production)"
  type        = bool
  default     = false
}

variable "enable_detailed_monitoring" {
  description = "Enable detailed CloudWatch monitoring"
  type        = bool
  default     = false
}

# ----------------------------------------------------------------------------
# User Data
# ----------------------------------------------------------------------------

variable "user_data" {
  description = "Custom user data script (leave empty to use default)"
  type        = string
  default     = ""
}

# ----------------------------------------------------------------------------
# Tags
# ----------------------------------------------------------------------------

variable "tags" {
  description = "Additional tags to apply to all resources"
  type        = map(string)
  default     = {}
}
