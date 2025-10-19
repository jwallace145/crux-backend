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
  db_name          = var.database.name

  # Networking configs
  vpc_id     = module.network.vpc_id
  subnet_ids = module.network.public_subnet_ids

  # Database user details
  db_username = var.database.user.username
  db_password = var.database.user.password

  # Storage configs
  allocated_storage     = var.database.storage.allocated_storage
  max_allocated_storage = var.database.storage.max_allocated_storage
}
