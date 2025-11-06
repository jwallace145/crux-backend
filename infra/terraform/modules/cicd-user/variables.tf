variable "project_name" {
  type        = string
  description = "The name of the project the CICD user is meant to build."
}

variable "environment" {
  description = "The environment of the CICD user."
  type        = string

  validation {
    condition     = contains(["dev", "stg", "prod"], var.environment)
    error_message = "The environment must be 'dev', 'stg', or 'prod'."
  }
}
