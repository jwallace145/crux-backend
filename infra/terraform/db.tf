# ================================
# Crux Project PostgreSQL Database
# ================================

module "db" {
  source = "./modules/rds-postgresql-db"

  # Environment
  environment = var.environment

  # Database details
  instance_class   = var.database.instance_class
  postgres_version = var.database.postgres_version
  multi_az         = var.database.multi_az
  db_name          = var.database.name

  # Networking configs
  vpc_id     = module.network.vpc_id
  subnet_ids = module.network.private_subnet_ids
  allowed_security_group_ids = concat(
    [module.alb.ecs_tasks_security_group_id],
    var.bastion.enabled ? [module.bastion[0].security_group_id] : []
  )
  publicly_accessible = false # DB in private subnets

  # Database user details
  db_username = var.db_user
  db_password = var.db_password

  # Storage configs
  allocated_storage     = var.database.storage.allocated_storage
  max_allocated_storage = var.database.storage.max_allocated_storage
}
