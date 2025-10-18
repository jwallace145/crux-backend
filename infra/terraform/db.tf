module "db" {
  source                = "./modules/rds-postgresql-db"
  environment           = var.environment
  vpc_id                = module.network.vpc_id
  subnet_ids            = module.network.public_subnet_ids
  db_name               = var.database.name
  db_username           = var.database.username
  db_password           = var.database.password
  instance_class        = var.database.instance_class
  postgres_version      = var.database.postgres_version
  allocated_storage     = var.database.allocated_storage
  max_allocated_storage = var.database.max_allocated_storage
}
