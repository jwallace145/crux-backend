# ============================================================================
# Application Load Balancer Module Variables
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
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: dev, staging, prod."
  }
}

variable "vpc_id" {
  description = "ID of the VPC where the ALB will be created"
  type        = string
}

variable "public_subnet_ids" {
  description = "List of public subnet IDs for the ALB (requires at least 2 for high availability)"
  type        = list(string)

  validation {
    condition     = length(var.public_subnet_ids) >= 2
    error_message = "At least 2 public subnets are required for ALB high availability."
  }
}

variable "container_port" {
  description = "Port number the container listens on"
  type        = number

  validation {
    condition     = var.container_port > 0 && var.container_port <= 65535
    error_message = "Container port must be between 1 and 65535."
  }
}

variable "domain" {
  description = "The domain name of the API."
  type        = string
}

# ----------------------------------------------------------------------------
# Health Check Configuration
# ----------------------------------------------------------------------------

variable "health_check_path" {
  description = "Path for health check requests"
  type        = string
  default     = "/health"
}

variable "health_check_interval" {
  description = "Interval between health checks (in seconds)"
  type        = number
  default     = 30

  validation {
    condition     = var.health_check_interval >= 5 && var.health_check_interval <= 300
    error_message = "Health check interval must be between 5 and 300 seconds."
  }
}

variable "health_check_timeout" {
  description = "Health check timeout (in seconds)"
  type        = number
  default     = 5

  validation {
    condition     = var.health_check_timeout >= 2 && var.health_check_timeout <= 120
    error_message = "Health check timeout must be between 2 and 120 seconds."
  }
}

variable "health_check_healthy_threshold" {
  description = "Number of consecutive successful health checks before marking target as healthy"
  type        = number
  default     = 2

  validation {
    condition     = var.health_check_healthy_threshold >= 2 && var.health_check_healthy_threshold <= 10
    error_message = "Healthy threshold must be between 2 and 10."
  }
}

variable "health_check_unhealthy_threshold" {
  description = "Number of consecutive failed health checks before marking target as unhealthy"
  type        = number
  default     = 3

  validation {
    condition     = var.health_check_unhealthy_threshold >= 2 && var.health_check_unhealthy_threshold <= 10
    error_message = "Unhealthy threshold must be between 2 and 10."
  }
}

variable "health_check_matcher" {
  description = "HTTP status codes to consider healthy (e.g., '200', '200-299')"
  type        = string
  default     = "200"
}

# ----------------------------------------------------------------------------
# HTTPS Configuration
# ----------------------------------------------------------------------------

variable "enable_https" {
  description = "Enable HTTPS listener (requires certificate_arn)"
  type        = bool
  default     = false
}

variable "certificate_arn" {
  description = "ARN of the SSL certificate for HTTPS listener (required if enable_https is true)"
  type        = string
  default     = ""

  validation {
    condition     = var.enable_https == false || (var.enable_https == true && length(var.certificate_arn) > 0)
    error_message = "Certificate ARN must be provided when HTTPS is enabled."
  }
}

variable "ssl_policy" {
  description = "SSL policy for HTTPS listener"
  type        = string
  default     = "ELBSecurityPolicy-TLS13-1-2-2021-06"
}

variable "redirect_http_to_https" {
  description = "Redirect HTTP traffic to HTTPS (only applies when enable_https is true)"
  type        = bool
  default     = true
}

# ----------------------------------------------------------------------------
# ALB Configuration
# ----------------------------------------------------------------------------

variable "internal" {
  description = "Whether the ALB is internal (true) or internet-facing (false)"
  type        = bool
  default     = false
}

variable "enable_deletion_protection" {
  description = "Enable deletion protection for the ALB"
  type        = bool
  default     = false
}

variable "enable_http2" {
  description = "Enable HTTP/2 for the ALB"
  type        = bool
  default     = true
}

variable "idle_timeout" {
  description = "Idle timeout for ALB connections (in seconds)"
  type        = number
  default     = 60

  validation {
    condition     = var.idle_timeout >= 1 && var.idle_timeout <= 4000
    error_message = "Idle timeout must be between 1 and 4000 seconds."
  }
}

variable "allowed_cidr_blocks" {
  description = "List of CIDR blocks allowed to access the ALB"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

# ----------------------------------------------------------------------------
# Access Logs Configuration
# ----------------------------------------------------------------------------

variable "enable_access_logs" {
  description = "Enable access logs for the ALB"
  type        = bool
  default     = false
}

variable "access_logs_bucket" {
  description = "S3 bucket for ALB access logs (required if enable_access_logs is true)"
  type        = string
  default     = ""
}

variable "access_logs_prefix" {
  description = "Prefix for ALB access logs in S3 bucket"
  type        = string
  default     = "alb-logs"
}

# ----------------------------------------------------------------------------
# Target Group Configuration
# ----------------------------------------------------------------------------

variable "deregistration_delay" {
  description = "Time to wait before deregistering a target (in seconds)"
  type        = number
  default     = 30

  validation {
    condition     = var.deregistration_delay >= 0 && var.deregistration_delay <= 3600
    error_message = "Deregistration delay must be between 0 and 3600 seconds."
  }
}

variable "enable_stickiness" {
  description = "Enable sticky sessions for the target group"
  type        = bool
  default     = false
}

variable "stickiness_duration" {
  description = "Duration of sticky session cookie (in seconds)"
  type        = number
  default     = 86400 # 24 hours

  validation {
    condition     = var.stickiness_duration >= 1 && var.stickiness_duration <= 604800
    error_message = "Stickiness duration must be between 1 and 604800 seconds (7 days)."
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
