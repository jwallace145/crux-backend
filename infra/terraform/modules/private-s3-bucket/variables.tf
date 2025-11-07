variable "project_name" {
  description = "Name of the project."
  type        = string
}

variable "environment" {
  description = "The environment of the project."
  type        = string

  validation {
    condition     = contains(["dev", "stg", "prod"], var.environment)
    error_message = "The environment must be 'dev', 'stg', or 'prod'."
  }
}

variable "enable_versioning" {
  description = "Enable versioning for the S3 bucket"
  type        = bool
  default     = true
}

variable "enable_encryption" {
  description = "Enable server-side encryption for the S3 bucket"
  type        = bool
  default     = true
}

variable "lifecycle_rules" {
  description = "List of lifecycle rules for the bucket"
  type = list(object({
    id                            = string
    enabled                       = bool
    prefix                        = optional(string)
    expiration_days               = optional(number)
    noncurrent_version_expiration = optional(number)
    transition_days               = optional(number)
    transition_storage_class      = optional(string)
  }))
  default = []
}

variable "read_access" {
  description = <<-EOT
    Map of read access permissions. Each entry specifies principals and prefixes they can read.
    Example:
    {
      "app-role" = {
        principal = "arn:aws:iam::123456789012:role/AppRole"
        prefixes  = ["data/", "logs/"]
      }
    }
  EOT
  type = map(object({
    principal = string
    prefixes  = list(string)
  }))
  default = {}
}

variable "write_access" {
  description = <<-EOT
    Map of write access permissions. Each entry specifies principals and prefixes they can write to.
    Example:
    {
      "data-processor" = {
        principal = "arn:aws:iam::123456789012:role/DataProcessor"
        prefixes  = ["uploads/", "temp/"]
      }
    }
  EOT
  type = map(object({
    principal = string
    prefixes  = list(string)
  }))
  default = {}
}

variable "delete_access" {
  description = <<-EOT
    Map of delete access permissions. Each entry specifies principals and prefixes they can delete from.
    Example:
    {
      "admin-role" = {
        principal = "arn:aws:iam::123456789012:role/AdminRole"
        prefixes  = ["temp/", "archive/"]
      }
    }
  EOT
  type = map(object({
    principal = string
    prefixes  = list(string)
  }))
  default = {}
}

variable "full_access_principals" {
  description = "List of principals that have full access to the bucket (all actions on all objects)"
  type        = list(string)
  default     = []
}

variable "block_public_access" {
  description = "Enable all S3 Block Public Access settings"
  type        = bool
  default     = true
}

variable "force_destroy" {
  description = "Allow bucket to be destroyed even if it contains objects"
  type        = bool
  default     = false
}
