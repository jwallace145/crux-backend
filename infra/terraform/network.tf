# ===================
# CruxProject Network
# ===================

module "network" {
  source = "./modules/vpc-network"

  project_name    = var.project_name
  environment     = var.environment
  region          = var.network.region
  vpc_cidr_block  = var.network.vpc_cidr_block
  public_subnets  = var.network.public_subnets
  private_subnets = var.network.private_subnets
  nat_gateways    = var.network.nat_gateways
}
