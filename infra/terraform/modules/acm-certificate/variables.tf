# ============================================================================
# ACM Certificate Module Variables
# ============================================================================

# ----------------------------------------------------------------------------
# Required Variables
# ----------------------------------------------------------------------------

variable "service_name" {
  description = "Name of the service (e.g., 'crux-backend')"
  type        = string

  validation {
    condition     = length(var.service_name) > 0 && length(var.service_name) <= 32
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

variable "domain_name" {
  description = "The primary domain name for the certificate (e.g., 'cruxproject.io')"
  type        = string

  validation {
    condition     = can(regex("^[a-z0-9][a-z0-9-\\.]*[a-z0-9]$", var.domain_name))
    error_message = "Domain name must be a valid domain format."
  }
}

# ----------------------------------------------------------------------------
# Optional Variables
# ----------------------------------------------------------------------------

variable "use_wildcard" {
  description = "Create a wildcard certificate (*.domain_name) instead of exact match"
  type        = bool
  default     = true
}

variable "subject_alternative_names" {
  description = "Additional domain names to include in the certificate (SANs)"
  type        = list(string)
  default     = []

  validation {
    condition = alltrue([
      for name in var.subject_alternative_names : can(regex("^[a-z0-9*][a-z0-9-\\.]*[a-z0-9]$", name))
    ])
    error_message = "All SANs must be valid domain formats."
  }
}

# ----------------------------------------------------------------------------
# Tags
# ----------------------------------------------------------------------------

variable "tags" {
  description = "Additional tags to apply to all resources"
  type        = map(string)
  default     = {}
}
