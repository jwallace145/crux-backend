variable "service_name" {
  type        = string
  description = "The service name of the CruxBackend API, Database, and Network."
  default     = "crux"
}

variable "environment" {
  description = "The environment of CruxBackend."
  type        = string

  validation {
    condition     = contains(["dev", "stg", "prod"], var.environment)
    error_message = "The environment must be 'dev', 'stg', or 'prod'."
  }
}
