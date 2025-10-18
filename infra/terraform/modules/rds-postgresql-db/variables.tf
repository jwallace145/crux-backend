variable "environment" {
  description = "The environment of the CruxBackend PostgreSQL database."
  type        = string

  validation {
    condition     = contains(["dev", "stg", "prod"], var.environment)
    error_message = "The environment must be 'dev', 'stg', or 'prod'."
  }
}

variable "db_name" {
  description = "Name of the PostgreSQL database"
  type        = string
}

variable "db_username" {
  description = "Master username for the database"
  type        = string
}

variable "db_password" {
  description = "Master password for the database (use AWS Secrets Manager in production)"
  type        = string
  sensitive   = true
}

variable "instance_class" {
  description = "RDS instance class (t4g.micro is cheapest)"
  type        = string
  default     = "db.t4g.micro"
}

variable "allocated_storage" {
  description = "Allocated storage in GB (minimum 20GB for gp3)"
  type        = number
  default     = 20
}

variable "max_allocated_storage" {
  description = "Maximum storage for autoscaling (0 to disable)"
  type        = number
  default     = 50
}

variable "postgres_version" {
  description = "PostgreSQL engine version"
  type        = string
  default     = "16.3"
}

variable "backup_retention_days" {
  description = "Number of days to retain backups (0-35, 7 recommended)"
  type        = number
  default     = 7
}



variable "vpc_id" {
  description = "The ID of the VPC network that the RDS database will be deployed."
  type        = string
}

variable "subnet_ids" {
  description = "The list of Subnet IDs for the RDS database."
  type        = list(string)
}

variable "allowed_cidr_blocks" {
  description = "CIDR blocks allowed to connect to RDS"
  type        = list(string)
  default     = ["0.0.0.0/0"] # WARNING: Restrict this in production!
}

variable "publicly_accessible" {
  description = "Make RDS publicly accessible (true for side projects, false for production)"
  type        = bool
  default     = true
}

variable "skip_final_snapshot" {
  description = "Skip final snapshot on deletion (true for dev, false for prod)"
  type        = bool
  default     = true
}

variable "deletion_protection" {
  description = "Enable deletion protection (false for dev, true for prod)"
  type        = bool
  default     = false
}

variable "tags" {
  description = "Additional tags for resources"
  type        = map(string)
  default     = {}
}
