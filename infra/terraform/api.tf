locals {
  api_name = "${var.service_name}-api"
}

# ============================================================================
# ACM Certificate for HTTPS
# ============================================================================

module "certificate" {
  source = "./modules/acm-certificate"

  service_name = local.api_name
  environment  = var.environment
  domain_name  = var.api.domain
  use_wildcard = true
}

# ============================================================================
# Application Load Balancer with HTTPS
# ============================================================================

module "alb" {
  source = "./modules/alb-ecs"

  service_name      = local.api_name
  environment       = var.environment
  domain            = var.api.domain
  vpc_id            = module.network.vpc_id
  public_subnet_ids = module.network.public_subnet_ids
  container_port    = var.api.container.port
  health_check_path = "/health"

  # HTTPS Configuration
  enable_https           = true
  certificate_arn        = module.certificate.certificate_arn
  redirect_http_to_https = true
}

module "api" {
  source = "./modules/ecs-service"

  # service configs
  service_name = local.api_name
  environment  = var.environment

  # networking configs
  region           = var.network.region
  vpc_id           = module.network.vpc_id
  subnet_ids       = module.network.private_subnet_ids
  assign_public_ip = false # Using private subnets with NAT gateway

  # load balancer configs
  target_group_arn            = module.alb.target_group_arn
  ecs_tasks_security_group_id = module.alb.ecs_tasks_security_group_id

  # task configs
  task_cpu    = var.api.task.cpu
  task_memory = var.api.task.memory

  # container configs
  container_image = var.api.container.image
  container_port  = var.api.container.port
  enable_ecs_exec = true
  environment_variables = [
    { name = "ENVIRONMENT", value = var.environment },
    { name = "PORT", value = var.api.container.port },
    { name = "DB_HOST", value = module.db.address },
    { name = "DB_PORT", value = module.db.port },
    { name = "DB_NAME", value = module.db.database_name },
    { name = "DB_USER", value = module.db.username },
    { name = "DB_PASSWORD", value = "cruxdbpassword" },
    { name = "DB_SSLMODE", value = "require" },
    { name = "ACCESS_TOKEN_SECRET_KEY", value = var.access_token_secret_key },
    { name = "REFRESH_TOKEN_SECRET_KEY", value = var.refresh_token_secret_key }
  ]
}
