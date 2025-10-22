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

variable "region" {
  description = "The AWS region of the ECS service."
}

variable "vpc_id" {
  description = "The ID of the VPC to place the ECS service."
  type        = string
}

variable "subnet_ids" {
  description = "The IDs of the subnets to place the ECS service in (public or private)."
  type        = list(string)
}

variable "assign_public_ip" {
  description = "Assign a public IP to the ECS tasks (required for public subnets, false for private subnets with NAT gateway)."
  type        = bool
  default     = false
}

variable "container_image" {
  description = "The image URI of the container to create the ECS service."
  type        = string
}

variable "container_port" {
  description = "The port of the container that is running in the ECS service."
  type        = number
}

variable "desired_count" {
  description = "The number of ECS services to place on the cluster."
  type        = number
  default     = 1
}

variable "task_cpu" {
  description = "The CPU allocated to the the ECS service."
  type        = number
}

variable "task_memory" {
  description = "The memory allocated to the ECS service."
  type        = number
}

variable "environment_variables" {
  description = "The environment variables set for the container in the ECS service."
  type = list(object({
    name  = string
    value = string
  }))
}

variable "enable_ecs_exec" {
  description = "Allow shell executions within the ECS container."
  type        = string
}

variable "target_group_arn" {
  description = "The ARN of the target group of the ALB."
  type        = string
}

variable "ecs_tasks_security_group_id" {
  description = "The security group ID for ECS tasks (created by alb-ecs module)."
  type        = string
}

variable "enable_container_insights" {
  description = "Enables/disables CloudWatch Container Insights for a specified cluster."
  type        = bool
  default     = true
}

variable "log_retention_days" {
  description = "The number of days to retain CruxBackend API CloudWatch logs."
  type        = number
  default     = 7
}
