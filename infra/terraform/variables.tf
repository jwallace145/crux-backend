variable "service_name" {
  type        = string
  description = "The service name of the CruxBackend API, Database, and Network."
  default     = "crux"
}

variable "environment" {
  description = "The environment of CruxBackend."
  type        = string

  validation {
    condition     = contains(local.valid_environments, var.environment)
    error_message = "The environment must be 'dev', 'stg', or 'prod'."
  }
}

variable "network" {
  description = "The network configurations of CruxBackend VPC and subnets."

  type = object({
    region         = string
    vpc_cidr_block = string
    nat_gateways = map(object({
      availability_zone = string
      enabled           = bool
    }))
    public_subnets = map(object({
      availability_zone = string
      subnet_cidr_block = string
    }))
    private_subnets = map(object({
      availability_zone = string
      subnet_cidr_block = string
    }))
  })

  validation {
    condition     = contains(local.valid_regions, var.network.region)
    error_message = "The AWS region must be 'us-east-1'."
  }

  validation {
    condition     = can(cidrhost(var.network.vpc_cidr_block, 0))
    error_message = "VPC CIDR block must be a valid IPv4 CIDR notation."
  }

  validation {
    condition = alltrue([
      for subnet in var.network.public_subnets : can(cidrhost(subnet.subnet_cidr_block, 0))
    ])
    error_message = "All public subnet CIDR blocks must be valid IPv4 CIDR notation."
  }

  validation {
    condition = alltrue([
      for subnet in var.network.public_subnets :
      contains(
        lookup(local.valid_availability_zones, var.network.region, []),
        subnet.availability_zone
      )
    ])
    error_message = "All public subnets must use a valid availability zone in the network region."
  }

  validation {
    condition = length(distinct([
      for nat in var.network.nat_gateways : nat.availability_zone
    ])) == length(var.network.nat_gateways)
    error_message = "Only a single NAT gateway can be defined for each availability zone."
  }
}

variable "api" {
  description = "Configuration for the Crux API ECS task and container"
  type = object({
    domain = string
    task = object({
      cpu    = number
      memory = number
    })
    container = object({
      cpu    = number
      memory = number
      image  = string
      port   = number
    })
  })

  validation {
    condition     = contains([256, 512, 1024, 2048, 4096], var.api.task.cpu)
    error_message = "Task CPU must be one of: 256, 512, 1024, 2048, 4096"
  }

  validation {
    condition     = var.api.task.memory >= 512 && var.api.task.memory <= 30720
    error_message = "Task memory must be between 512 and 30720 MB"
  }

  validation {
    condition     = var.api.container.cpu <= var.api.task.cpu
    error_message = "Container CPU cannot exceed task CPU"
  }

  validation {
    condition     = var.api.container.memory <= var.api.task.memory
    error_message = "Container memory cannot exceed task memory"
  }

  validation {
    condition     = var.api.container.port > 0 && var.api.container.port <= 65535
    error_message = "Container port must be between 1 and 65535"
  }

  validation {
    condition     = can(regex("^[0-9]+\\.dkr\\.ecr\\.[a-z0-9-]+\\.amazonaws\\.com/.+:.+$", var.api.container.image))
    error_message = "Container image must be a valid ECR image URL"
  }
}

variable "database" {
  description = "The configurations for the CruxBackend PostgreSQL database."

  type = object({
    postgres_version = string
    instance_class   = string
    multi_az         = bool
    name             = string
    storage = object({
      allocated_storage     = number
      max_allocated_storage = number
    })
  })
}

variable "bastion" {
  description = "Configuration for the bastion host instance"
  type = object({
    enabled                 = bool
    instance_type           = string
    allowed_ssh_cidr_blocks = list(string)
  })

  validation {
    condition     = contains(["t3.micro", "t3.medium"], var.bastion.instance_type)
    error_message = "Instance type must be either t3.micro or t3.medium"
  }

  validation {
    condition = alltrue([
      for cidr in var.bastion.allowed_ssh_cidr_blocks : can(cidrhost(cidr, 0))
    ])
    error_message = "All CIDR blocks must be valid IPv4 CIDR notation"
  }
}

variable "db_user" {
  description = "The name of the DB user for the application."
  type        = string
}

variable "db_password" {
  description = "The password for the DB user for the application."
  type        = string
}

variable "access_token_secret_key" {
  description = "The secret key used for the API access token."
  type        = string
}

variable "refresh_token_secret_key" {
  description = "The secret key used for the API refresh token."
  type        = string
}
