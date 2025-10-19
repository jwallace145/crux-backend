variable "secret_name" {
  description = "The name of the secret in AWS Secrets Manager"
  type        = string

  validation {
    condition     = length(var.secret_name) > 0 && length(var.secret_name) <= 512
    error_message = "Secret name must be between 1 and 512 characters"
  }

  validation {
    condition     = can(regex("^[a-zA-Z0-9/_+=.@-]+$", var.secret_name))
    error_message = "Secret name can only contain alphanumeric characters and the following: /_+=.@-"
  }
}

variable "secret_description" {
  description = "Description of the secret"
  type        = string
  default     = ""
}

variable "secrets" {
  description = "Map of secret key-value pairs to store as JSON in AWS Secrets Manager"
  type        = map(string)
  sensitive   = true

  validation {
    condition     = length(var.secrets) > 0
    error_message = "At least one secret key-value pair must be provided"
  }
}

variable "recovery_window_in_days" {
  description = "Number of days to retain the secret before permanent deletion (0 for immediate deletion, 7-30 for recovery window)"
  type        = number
  default     = 30

  validation {
    condition     = var.recovery_window_in_days == 0 || (var.recovery_window_in_days >= 7 && var.recovery_window_in_days <= 30)
    error_message = "Recovery window must be 0 (immediate deletion) or between 7 and 30 days"
  }
}

variable "tags" {
  description = "Tags to apply to the secret"
  type        = map(string)
  default     = {}
}

variable "kms_key_id" {
  description = "ARN or ID of the AWS KMS key to encrypt the secret. If not specified, uses the default AWS managed key"
  type        = string
  default     = null
}

variable "enable_rotation" {
  description = "Enable automatic rotation for this secret"
  type        = bool
  default     = false
}

variable "rotation_lambda_arn" {
  description = "ARN of the Lambda function that can rotate the secret. Required if enable_rotation is true"
  type        = string
  default     = null
}

variable "rotation_days" {
  description = "Number of days between automatic rotations. Required if enable_rotation is true"
  type        = number
  default     = 30

  validation {
    condition     = var.rotation_days >= 1
    error_message = "Rotation days must be at least 1"
  }
}
